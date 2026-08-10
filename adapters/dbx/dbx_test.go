package dbx

import (
	"context"
	"testing"

	"github.com/lcylpzls/dbx"
	"github.com/lcylpzls/eventx"
	"github.com/lcylpzls/testx"
)

func TestHookForwards(t *testing.T) {
	bus, err := eventx.New()
	testx.RequireNoError(t, err)
	got := make(chan dbx.QueryEvent, 1)
	_, _ = bus.Subscribe("dbx.query.exec", func(ctx context.Context, e eventx.Event) error {
		got <- e.Payload.(dbx.QueryEvent)
		return nil
	})
	h := Hook(bus)
	h.OnQueryEvent(context.Background(), dbx.QueryEvent{Operation: "exec", Statement: "SELECT 1"})
	select {
	case e := <-got:
		if e.Operation != "exec" || e.Statement != "SELECT 1" {
			t.Fatalf("事件不匹配：%+v", e)
		}
	default:
		t.Fatal("适配器未转发事件")
	}
	_ = bus.Close()
}
