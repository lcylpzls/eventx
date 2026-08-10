package ecosystem_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/lcylpzls/cachex"
	"github.com/lcylpzls/eventx"
	cachexadapter "github.com/lcylpzls/eventx/adapters/cachex"
	filexadapter "github.com/lcylpzls/eventx/adapters/filex"
	jobxadapter "github.com/lcylpzls/eventx/adapters/jobx"
	"github.com/lcylpzls/filex"
	"github.com/lcylpzls/jobx"
	"github.com/lcylpzls/testx"
)

// TestEcosystemLink 端到端联动：
// filex 上传 → eventx 事件 → cachex 失效 + jobx 异步处理。
func TestEcosystemLink(t *testing.T) {
	bus, err := eventx.New()
	testx.RequireNoError(t, err)
	defer bus.Close()

	cache, err := cachex.New(cachex.WithEventHook(cachexadapter.Hook(bus)))
	testx.RequireNoError(t, err)
	defer cache.Close()

	dispatcher, err := jobx.NewDispatcher(jobx.WithEventHook(jobxadapter.Hook(bus)))
	testx.RequireNoError(t, err)
	defer dispatcher.Shutdown(context.Background())
	if err := dispatcher.Handle("thumbnail", func(ctx context.Context, job jobx.Job) error {
		return nil
	}); err != nil {
		t.Fatalf("Handle 失败：%v", err)
	}

	store, err := filex.New(filex.Config{
		DataDir:   t.TempDir(),
		EventHook: filexadapter.Hook(bus),
	})
	testx.RequireNoError(t, err)
	defer store.Close()
	if _, err := store.CreateBucket(context.Background(), "bucket"); err != nil {
		t.Fatalf("CreateBucket 失败：%v", err)
	}

	// 项目层组装联动：对象上传 → 缓存失效 + 异步任务。
	_, err = eventx.SubscribeTyped[filex.ObjectEvent](bus, "filex.object.put",
		func(ctx context.Context, topic string, e filex.ObjectEvent) error {
			cache.Delete(e.Bucket + "/" + e.Key)
			_, err := dispatcher.Submit(ctx, "thumbnail", []byte(e.Bucket+"/"+e.Key))
			return err
		})
	testx.RequireNoError(t, err)
	jobDone := make(chan struct{}, 1)
	_, err = bus.Subscribe("jobx.task.completed", func(ctx context.Context, e eventx.Event) error {
		jobDone <- struct{}{}
		return nil
	})
	testx.RequireNoError(t, err)

	// 预置缓存值，上传对象后应被事件联动清除。
	cache.Set("bucket/k", "旧值")
	if _, ok := cache.Get("bucket/k"); !ok {
		t.Fatal("预置缓存值缺失")
	}
	if _, err := store.Put(context.Background(), "bucket", "k", strings.NewReader("数据"), filex.PutOptions{}); err != nil {
		t.Fatalf("Put 失败：%v", err)
	}
	if _, ok := cache.Get("bucket/k"); ok {
		t.Fatal("事件联动应清除缓存")
	}

	select {
	case <-jobDone:
	case <-time.After(2 * time.Second):
		t.Fatal("异步任务未完成")
	}
}
