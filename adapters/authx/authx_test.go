package authx

import (
	"context"
	"errors"
	"testing"

	"github.com/lcylpzls/authx"
	"github.com/lcylpzls/eventx"
)

func TestHookForwards(t *testing.T) {
	bus, err := eventx.New()
	if err != nil {
		t.Fatalf("New 失败：%v", err)
	}
	got := make(chan authx.AuthEvent, 1)
	_, _ = bus.Subscribe("authx.token.issue", func(ctx context.Context, e eventx.Event) error {
		got <- e.Payload.(authx.AuthEvent)
		return nil
	})
	h := Hook(bus)
	h.OnAuthEvent(context.Background(), authx.AuthEvent{Action: "issue"})
	h.OnAuthEvent(context.Background(), authx.AuthEvent{Action: "validate", Err: errors.New("失败")})
	select {
	case e := <-got:
		if e.Action != "issue" {
			t.Fatalf("事件不匹配：%+v", e)
		}
	default:
		t.Fatal("适配器未转发事件")
	}
	_ = bus.Close()
}
