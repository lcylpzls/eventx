package basic_test

import (
	"context"
	"testing"

	"github.com/lcylpzls/eventx"
)

func TestExamplePubSub(t *testing.T) {
	bus := eventx.New()
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
}
