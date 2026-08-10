package basic_test

import (
	"context"
	"testing"

	"github.com/lcylpzls/eventx"
)

func TestExamplePubSub(t *testing.T) {
	bus, err := eventx.New()
	if err != nil {
		t.Fatalf("New 失败：%v", err)
	}
	var got any
	sub, err := bus.Subscribe("orders.created", func(ctx context.Context, e eventx.Event) error {
		got = e.Payload
		return nil
	})
	if err != nil {
		t.Fatalf("Subscribe 失败：%v", err)
	}
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
	if err := bus.Close(); err != nil {
		t.Fatalf("Close 失败：%v", err)
	}
}
