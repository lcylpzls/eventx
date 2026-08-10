package eventx

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/lcylpzls/errx"
)

func TestPublishBasic(t *testing.T) {
	bus := New()
	var got Event
	sub, err := bus.Subscribe("orders.created", func(ctx context.Context, e Event) error {
		got = e
		return nil
	})
	if err != nil {
		t.Fatalf("Subscribe 失败：%v", err)
	}
	if sub.ID() == 0 || sub.Topic() != "orders.created" {
		t.Fatalf("订阅访问器不匹配：id=%d topic=%q", sub.ID(), sub.Topic())
	}
	if err := bus.Publish(context.Background(), "orders.created", "10086"); err != nil {
		t.Fatalf("Publish 失败：%v", err)
	}
	if got.Topic != "orders.created" || got.Payload != "10086" {
		t.Fatalf("事件不匹配：%+v", got)
	}
}

func TestSubscribeOrder(t *testing.T) {
	bus := New()
	var order []string
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		order = append(order, "一")
		return nil
	})
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		order = append(order, "二")
		return nil
	})
	if err := bus.Publish(context.Background(), "t", nil); err != nil {
		t.Fatalf("Publish 失败：%v", err)
	}
	if strings.Join(order, ",") != "一,二" {
		t.Fatalf("订阅顺序不匹配：%v", order)
	}
}

func TestPublishAggregateErrors(t *testing.T) {
	bus := New()
	errA := errors.New("错误A")
	errB := errors.New("错误B")
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		return errA
	})
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		return errB
	})
	err := bus.Publish(context.Background(), "t", nil)
	if err == nil {
		t.Fatal("应返回聚合错误")
	}
	if !errors.Is(err, errA) || !errors.Is(err, errB) {
		t.Fatalf("聚合错误应包含全部订阅者错误：%v", err)
	}
}

func TestHandlerPanic(t *testing.T) {
	bus := New()
	called := false
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		panic("订阅者崩溃")
	})
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		called = true
		return nil
	})
	err := bus.Publish(context.Background(), "t", nil)
	assertErrCode(t, err, CodeHandlerPanic)
	if !called {
		t.Fatal("panic 恢复后其他订阅者仍应执行")
	}
}

func TestUnsubscribe(t *testing.T) {
	bus := New()
	count := 0
	sub, _ := bus.Subscribe("t", func(ctx context.Context, e Event) error {
		count++
		return nil
	})
	if err := sub.Unsubscribe(); err != nil {
		t.Fatalf("Unsubscribe 失败：%v", err)
	}
	if err := sub.Unsubscribe(); err != nil {
		t.Fatalf("重复 Unsubscribe 应幂等：%v", err)
	}
	if err := bus.Publish(context.Background(), "t", nil); err != nil {
		t.Fatalf("Publish 失败：%v", err)
	}
	if count != 0 {
		t.Fatalf("取消后不应收到事件：%d", count)
	}
}

func TestInvalidTopic(t *testing.T) {
	long := strings.Repeat("a", 257)
	for _, topic := range []string{"", long, "a\x01b", "orders.*"} {
		_, err := New().Subscribe(topic, func(ctx context.Context, e Event) error { return nil })
		assertErrCode(t, err, CodeInvalidTopic)
		if err := New().Publish(context.Background(), topic, nil); err == nil {
			t.Fatalf("非法主题 %q 发布应报错", topic)
		}
	}
}

func TestNilHandler(t *testing.T) {
	_, err := New().Subscribe("t", nil)
	assertErrCode(t, err, CodeInvalidHandler)
}

func TestClose(t *testing.T) {
	bus := New()
	sub, _ := bus.Subscribe("t", func(ctx context.Context, e Event) error { return nil })
	if err := bus.Close(); err != nil {
		t.Fatalf("Close 失败：%v", err)
	}
	if err := bus.Close(); err != nil {
		t.Fatalf("重复 Close 应幂等：%v", err)
	}
	_, err := bus.Subscribe("t", func(ctx context.Context, e Event) error { return nil })
	assertErrCode(t, err, CodeBusClosed)
	err = bus.Publish(context.Background(), "t", nil)
	assertErrCode(t, err, CodeBusClosed)
	err = sub.Unsubscribe()
	assertErrCode(t, err, CodeSubscriptionNotFound)
}

func TestPublishNoSubscriber(t *testing.T) {
	if err := New().Publish(context.Background(), "t", nil); err != nil {
		t.Fatalf("无订阅者发布应返回 nil：%v", err)
	}
}

func TestPublishSkipsUnsubscribedInFlight(t *testing.T) {
	bus := New()
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
	if err := bus.Publish(context.Background(), "t", nil); err != nil {
		t.Fatalf("Publish 失败：%v", err)
	}
}

func TestConcurrent(t *testing.T) {
	bus := New()
	var mu sync.Mutex
	count := 0
	sub, _ := bus.Subscribe("t", func(ctx context.Context, e Event) error {
		mu.Lock()
		count++
		mu.Unlock()
		return nil
	})
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				_ = bus.Publish(context.Background(), "t", j)
			}
		}()
	}
	wg.Wait()
	_ = sub.Unsubscribe()
	if count != 16*50 {
		t.Fatalf("订阅计数不匹配：%d", count)
	}
}

func assertErrCode(t *testing.T, err error, want errx.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("期望错误 %q，得到 nil", want)
	}
	if !errx.Is(err, want) {
		t.Fatalf("期望错误码 %q，得到 %v", want, err)
	}
}
