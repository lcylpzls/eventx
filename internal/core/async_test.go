package core

import (
	"context"
	"errors"
	"github.com/lcylpzls/testx"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPublishAsync(t *testing.T) {
	bus := newBus(t)
	done := make(chan any, 1)
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		done <- e.Payload
		return nil
	})
	if err := bus.PublishAsync(context.Background(), "t", "异步载荷"); err != nil {
		t.Fatalf("PublishAsync 失败：%v", err)
	}
	select {
	case got := <-done:
		if got != "异步载荷" {
			t.Fatalf("载荷不匹配：%v", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("异步任务未执行")
	}
}

func TestPublishAsyncQueueFull(t *testing.T) {
	bus := newBus(t, WithWorkers(1), WithQueueSize(1))
	release := make(chan struct{})
	started := make(chan struct{}, 1)
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		started <- struct{}{}
		<-release
		return nil
	})
	if err := bus.PublishAsync(context.Background(), "t", 1); err != nil {
		t.Fatalf("首次投递失败：%v", err)
	}
	<-started
	if err := bus.PublishAsync(context.Background(), "t", 2); err != nil {
		t.Fatalf("队列未满投递失败：%v", err)
	}
	err := bus.PublishAsync(context.Background(), "t", 3)
	assertErrCode(t, err, CodeQueueFull)
	close(release)
}

func TestPublishAsyncErrorHandler(t *testing.T) {
	bus := newBus(t, WithErrorHandler(func(err error) {}))
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		return errors.New("异步失败")
	})
	if err := bus.PublishAsync(context.Background(), "t", nil); err != nil {
		t.Fatalf("PublishAsync 失败：%v", err)
	}
	// 等待 worker 处理完成（错误被回调消费，不阻塞）。
	time.Sleep(50 * time.Millisecond)
}

func TestPublishAsyncErrorCallback(t *testing.T) {
	errCh := make(chan error, 1)
	bus := newBus(t, WithErrorHandler(func(err error) {
		errCh <- err
	}))
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		return errors.New("处理失败")
	})
	if err := bus.PublishAsync(context.Background(), "t", nil); err != nil {
		t.Fatalf("PublishAsync 失败：%v", err)
	}
	select {
	case err := <-errCh:
		if err == nil || err.Error() != "处理失败" {
			t.Fatalf("错误回调不匹配：%v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("错误回调未触发")
	}
}

func TestPublishAsyncPanicRecovered(t *testing.T) {
	errCh := make(chan error, 1)
	bus := newBus(t, WithErrorHandler(func(err error) {
		errCh <- err
	}))
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		panic("异步崩溃")
	})
	if err := bus.PublishAsync(context.Background(), "t", nil); err != nil {
		t.Fatalf("PublishAsync 失败：%v", err)
	}
	select {
	case err := <-errCh:
		assertErrCode(t, err, CodeHandlerPanic)
	case <-time.After(2 * time.Second):
		t.Fatal("panic 未恢复上报")
	}
}

func TestCloseDrainsAsync(t *testing.T) {
	bus := newBus(t, WithWorkers(1), WithQueueSize(16))
	var count atomic.Int32
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		time.Sleep(time.Millisecond)
		count.Add(1)
		return nil
	})
	for i := 0; i < 10; i++ {
		if err := bus.PublishAsync(context.Background(), "t", i); err != nil {
			t.Fatalf("投递失败：%v", err)
		}
	}
	if err := bus.Close(); err != nil {
		t.Fatalf("Close 失败：%v", err)
	}
	if got := count.Load(); got != 10 {
		t.Fatalf("优雅关闭应处理全部在途任务：%d", got)
	}
}

func TestPublishAsyncAfterClose(t *testing.T) {
	bus := newBus(t)
	if err := bus.Close(); err != nil {
		t.Fatalf("Close 失败：%v", err)
	}
	err := bus.PublishAsync(context.Background(), "t", nil)
	assertErrCode(t, err, CodeBusClosed)
}

func TestPublishAsyncInvalidTopic(t *testing.T) {
	bus := newBus(t)
	err := bus.PublishAsync(context.Background(), "", nil)
	assertErrCode(t, err, CodeInvalidTopic)
}

func TestNewInvalidOptions(t *testing.T) {
	_, err := New(WithWorkers(0))
	assertErrCode(t, err, CodeInvalidOption)
	_, err = New(WithQueueSize(0))
	assertErrCode(t, err, CodeInvalidOption)
}

func TestNewNilOption(t *testing.T) {
	bus, err := New(nil)
	testx.RequireNoError(t, err)
	_ = bus
}

func TestPublishAsyncConcurrent(t *testing.T) {
	bus := newBus(t, WithWorkers(4), WithQueueSize(64))
	var mu sync.Mutex
	count := 0
	var published atomic.Int32
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		mu.Lock()
		count++
		mu.Unlock()
		return nil
	})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				if err := bus.PublishAsync(context.Background(), "t", j); err == nil {
					published.Add(1)
				}
			}
		}()
	}
	wg.Wait()
	_ = bus.Close()
	if got := int32(count); got != published.Load() {
		t.Fatalf("成功投递的任务应全部被处理：count=%d published=%d", count, published.Load())
	}
	if published.Load() == 0 {
		t.Fatal("测试应至少成功投递一个任务")
	}
}

func TestPublishAsyncSkipsUnsubscribedInFlight(t *testing.T) {
	bus := newBus(t)
	sub := &subscription{
		id:      1,
		topic:   "t",
		handler: func(ctx context.Context, e Event) error { t.Fatal("已取消订阅不应执行"); return nil },
		bus:     bus,
	}
	sub.closed.Store(true)
	bus.mu.Lock()
	bus.subs["t"] = []*subscription{sub}
	bus.mu.Unlock()
	done := make(chan struct{})
	go func() {
		bus.async.process(bus, asyncTask{ctx: context.Background(), topic: "t"})
		close(done)
	}()
	<-done
}
