package core

import (
	"context"
	"github.com/lcylpzls/testx"
	"testing"
	"time"
)

type orderEvent struct {
	ID string
}

func TestSubscribeTyped(t *testing.T) {
	bus := newBus(t)
	got := make(chan orderEvent, 1)
	sub, err := SubscribeTyped(bus, "orders.created", func(ctx context.Context, topic string, payload orderEvent) error {
		if topic != "orders.created" {
			t.Fatalf("主题不匹配：%q", topic)
		}
		got <- payload
		return nil
	})
	testx.RequireNoError(t, err)
	defer sub.Unsubscribe()
	if err := bus.Publish(context.Background(), "orders.created", orderEvent{ID: "10086"}); err != nil {
		t.Fatalf("Publish 失败：%v", err)
	}
	select {
	case p := <-got:
		if p.ID != "10086" {
			t.Fatalf("类型化载荷不匹配：%+v", p)
		}
	case <-time.After(time.Second):
		t.Fatal("类型化订阅未收到事件")
	}
}

func TestSubscribeTypedNilHandler(t *testing.T) {
	bus := newBus(t)
	_, err := SubscribeTyped[int](bus, "t", nil)
	assertErrCode(t, err, CodeInvalidHandler)
}

func TestSubscribeTypedNilBus(t *testing.T) {
	_, err := SubscribeTyped[int](nil, "t", func(ctx context.Context, topic string, payload int) error {
		return nil
	})
	assertErrCode(t, err, CodeInvalidOption)
}

func TestSubscribeTypedMismatch(t *testing.T) {
	bus := newBus(t)
	otherCalled := false
	_, _ = SubscribeTyped[int](bus, "t", func(ctx context.Context, topic string, payload int) error {
		return nil
	})
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		otherCalled = true
		return nil
	})
	err := bus.Publish(context.Background(), "t", "字符串")
	assertErrCode(t, err, CodeTypeMismatch)
	if !otherCalled {
		t.Fatal("类型不匹配不应影响其他订阅者")
	}
}

func TestPublishCancelledAtEntry(t *testing.T) {
	bus := newBus(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := bus.Publish(ctx, "t", nil)
	assertErrCode(t, err, CodeCancelled)
	err = bus.PublishAsync(ctx, "t", nil)
	assertErrCode(t, err, CodeCancelled)
}

func TestDeliverCancelledMidway(t *testing.T) {
	bus := newBus(t)
	ctx, cancel := context.WithCancel(context.Background())
	secondCalled := false
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		cancel()
		return nil
	})
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		secondCalled = true
		return nil
	})
	err := bus.Publish(ctx, "t", nil)
	assertErrCode(t, err, CodeCancelled)
	if secondCalled {
		t.Fatal("取消后不应继续分发")
	}
}

func TestPublishAsyncCancelledMidDelivery(t *testing.T) {
	errCh := make(chan error, 1)
	bus := newBus(t, WithErrorHandler(func(err error) {
		errCh <- err
	}))
	ctx, cancel := context.WithCancel(context.Background())
	firstStarted := make(chan struct{}, 1)
	release := make(chan struct{})
	secondCalled := false
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		firstStarted <- struct{}{}
		<-release
		return nil
	})
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error {
		secondCalled = true
		return nil
	})
	if err := bus.PublishAsync(ctx, "t", nil); err != nil {
		t.Fatalf("PublishAsync 失败：%v", err)
	}
	<-firstStarted
	cancel()
	close(release)
	select {
	case err := <-errCh:
		assertErrCode(t, err, CodeCancelled)
	case <-time.After(2 * time.Second):
		t.Fatal("取消错误未上报")
	}
	if secondCalled {
		t.Fatal("投递后取消不应继续分发")
	}
}
