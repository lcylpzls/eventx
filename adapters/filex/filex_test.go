package filex

import (
	"context"
	"testing"

	"github.com/lcylpzls/eventx"
	"github.com/lcylpzls/filex"
	"github.com/lcylpzls/testx"
)

func TestHookForwards(t *testing.T) {
	bus, err := eventx.New()
	testx.RequireNoError(t, err)
	got := make(chan filex.ObjectEvent, 1)
	_, _ = bus.Subscribe("filex.object.put", func(ctx context.Context, e eventx.Event) error {
		got <- e.Payload.(filex.ObjectEvent)
		return nil
	})
	h := Hook(bus)
	h.OnObjectEvent(context.Background(), filex.ObjectEvent{
		Bucket: "b", Key: "k", Action: "put",
	})
	select {
	case e := <-got:
		if e.Action != "put" || e.Key != "k" {
			t.Fatalf("事件不匹配：%+v", e)
		}
	default:
		t.Fatal("适配器未转发事件")
	}
	_ = bus.Close()
}
