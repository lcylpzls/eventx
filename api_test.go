package eventx_test

import (
	"context"
	"testing"

	"github.com/lcylpzls/eventx"
)

// TestPublicAPI 黑盒冒烟测试：覆盖根包全部转发函数、类型别名与常量。
func TestPublicAPI(t *testing.T) {
	bus, err := eventx.New(
		eventx.WithWorkers(2),
		eventx.WithQueueSize(16),
		eventx.WithErrorHandler(func(error) {}),
		eventx.WithLogger(nil),
	)
	if err != nil || bus == nil {
		t.Fatalf("New 失败：%v", err)
	}
	defer bus.Close()

	sub, err := bus.Subscribe("orders.created", func(context.Context, eventx.Event) error {
		return nil
	})
	if err != nil || sub == nil {
		t.Fatalf("Subscribe 失败：%v", err)
	}
	sub2, err := bus.SubscribeFiltered("orders.*", func(context.Context, eventx.Event) bool {
		return true
	}, func(context.Context, eventx.Event) error {
		return nil
	})
	if err != nil || sub2 == nil {
		t.Fatalf("SubscribeFiltered 失败：%v", err)
	}
	sub3, err := bus.SubscribeWithOptions("orders.updated", func(context.Context, eventx.Event) error {
		return nil
	}, eventx.WithPriority(10))
	if err != nil || sub3 == nil {
		t.Fatalf("SubscribeWithOptions 失败：%v", err)
	}
	sub4, err := eventx.SubscribeTyped[int](bus, "orders.typed",
		func(context.Context, string, int) error { return nil })
	if err != nil || sub4 == nil {
		t.Fatalf("SubscribeTyped 失败：%v", err)
	}

	if err := bus.Publish(context.Background(), "orders.created", "x"); err != nil {
		t.Fatalf("Publish 失败：%v", err)
	}
	if err := bus.PublishAsync(context.Background(), "orders.typed", 1); err != nil {
		t.Fatalf("PublishAsync 失败：%v", err)
	}
	_ = bus.Metrics()

	_ = sub.ID()
	_ = sub.Topic()
	_ = sub.Unsubscribe()
	_ = sub2.Unsubscribe()
	_ = sub3.Unsubscribe()
	_ = sub4.Unsubscribe()

	var _ eventx.Event
	var _ eventx.Handler
	var _ eventx.Filter
	var _ eventx.Subscription
	var _ eventx.SubscribeOption
	var _ eventx.Metrics
	var _ eventx.MetricProvider
	var _ eventx.Option
	var _ eventx.TypedHandler[int]
	_ = eventx.CodeBusClosed
	_ = eventx.CodeInvalidTopic
	_ = eventx.CodeInvalidHandler
	_ = eventx.CodeSubscriptionNotFound
	_ = eventx.CodeHandlerPanic
	_ = eventx.CodeQueueFull
	_ = eventx.CodeInvalidOption
	_ = eventx.CodeCancelled
	_ = eventx.CodeTypeMismatch
}
