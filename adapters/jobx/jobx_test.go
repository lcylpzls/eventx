package jobx

import (
	"context"
	"testing"

	"github.com/lcylpzls/eventx"
	"github.com/lcylpzls/jobx"
	"github.com/lcylpzls/testx"
)

func TestHookForwards(t *testing.T) {
	bus, err := eventx.New()
	testx.RequireNoError(t, err)
	got := make(chan jobx.TaskEvent, 1)
	_, _ = bus.Subscribe("jobx.task.completed", func(ctx context.Context, e eventx.Event) error {
		got <- e.Payload.(jobx.TaskEvent)
		return nil
	})
	h := Hook(bus)
	h.OnTaskEvent(context.Background(), jobx.TaskEvent{Action: "completed", Name: "task"})
	select {
	case e := <-got:
		if e.Action != "completed" || e.Name != "task" {
			t.Fatalf("事件不匹配：%+v", e)
		}
	default:
		t.Fatal("适配器未转发事件")
	}
	_ = bus.Close()
}
