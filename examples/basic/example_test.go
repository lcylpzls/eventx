package basic_test

import (
	"context"
	"testing"

	"github.com/lcylpzls/eventx"
	"github.com/lcylpzls/testx"
)

func TestExamplePubSub(t *testing.T) {
	bus, err := eventx.New()
	testx.RequireNoError(t, err)
	var got any
	sub, err := bus.Subscribe("orders.created", func(ctx context.Context, e eventx.Event) error {
		got = e.Payload
		return nil
	})
	testx.RequireNoError(t, err)
	defer sub.Unsubscribe()
	if err := bus.Publish(context.Background(), "orders.created", "10086"); err != nil {
		t.Fatalf("Publish 失败：%v", err)
	}
	if got != "10086" {
		t.Fatalf("载荷不匹配：%v", got)
	}

	if err := bus.PublishAsync(context.Background(), "orders.created", "10087"); err != nil {
		t.Fatalf("PublishAsync 失败：%v", err)
	}
	_, err = bus.SubscribeFiltered("orders.*", func(ctx context.Context, e eventx.Event) bool {
		return e.Payload != nil
	}, func(ctx context.Context, e eventx.Event) error {
		return nil
	})
	testx.RequireNoError(t, err)
	_, err = eventx.SubscribeTyped[int](bus, "orders.count", func(ctx context.Context, topic string, payload int) error {
		return nil
	})
	testx.RequireNoError(t, err)
	m := bus.Metrics()
	if m.Publishes == 0 {
		t.Fatal("指标应记录发布次数")
	}
	if err := bus.Close(); err != nil {
		t.Fatalf("Close 失败：%v", err)
	}
}
