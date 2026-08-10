package cachex

import (
	"context"
	"testing"

	"github.com/lcylpzls/cachex"
	"github.com/lcylpzls/eventx"
)

func TestHookForwards(t *testing.T) {
	bus, err := eventx.New()
	if err != nil {
		t.Fatalf("New 失败：%v", err)
	}
	got := make(chan cachex.CacheEvent, 1)
	_, _ = bus.Subscribe("cachex.cache.set", func(ctx context.Context, e eventx.Event) error {
		got <- e.Payload.(cachex.CacheEvent)
		return nil
	})
	h := Hook(bus)
	h.OnCacheEvent(context.Background(), cachex.CacheEvent{Action: "set", Key: "k"})
	select {
	case e := <-got:
		if e.Action != "set" || e.Key != "k" {
			t.Fatalf("事件不匹配：%+v", e)
		}
	default:
		t.Fatal("适配器未转发事件")
	}
	_ = bus.Close()
}
