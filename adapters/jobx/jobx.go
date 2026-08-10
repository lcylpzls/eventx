// Package jobx 提供 jobx.EventHook 到 eventx 总线的适配。
package jobx

import (
	"context"

	"github.com/lcylpzls/eventx"
	"github.com/lcylpzls/jobx"
)

// Hook 返回接入 eventx 的 jobx.EventHook，主题为 jobx.task.<action>。
func Hook(bus *eventx.Bus) jobx.EventHook {
	return jobxHook{bus: bus}
}

type jobxHook struct {
	bus *eventx.Bus
}

func (h jobxHook) OnTaskEvent(ctx context.Context, e jobx.TaskEvent) {
	_ = h.bus.Publish(ctx, "jobx.task."+e.Action, e)
}
