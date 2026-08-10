package clix

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lcylpzls/clix"
	"github.com/lcylpzls/eventx"
	"github.com/lcylpzls/testx"
)

func TestObserverForwards(t *testing.T) {
	bus, err := eventx.New()
	testx.RequireNoError(t, err)
	started := make(chan CommandEvent, 1)
	finished := make(chan CommandEvent, 1)
	_, _ = bus.Subscribe("clix.command.start", func(ctx context.Context, e eventx.Event) error {
		started <- e.Payload.(CommandEvent)
		return nil
	})
	_, _ = bus.Subscribe("clix.command.finish", func(ctx context.Context, e eventx.Event) error {
		finished <- e.Payload.(CommandEvent)
		return nil
	})
	obs := Observer(bus)
	obs.OnCommandStart(context.Background(), "greet hello", []string{"a"})
	obs.OnCommandFinish(context.Background(), "greet hello", []string{"a"}, errors.New("失败"), time.Second)

	select {
	case e := <-started:
		if e.Command != "greet hello" || len(e.Args) != 1 {
			t.Fatalf("start 事件不匹配：%+v", e)
		}
	default:
		t.Fatal("start 事件未转发")
	}
	select {
	case e := <-finished:
		if e.Command != "greet hello" || e.Err == nil || e.Duration != time.Second {
			t.Fatalf("finish 事件不匹配：%+v", e)
		}
	default:
		t.Fatal("finish 事件未转发")
	}
	_ = bus.Close()
}

var _ clix.Observer = clixObserver{}
