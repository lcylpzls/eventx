// Package filex 提供 filex.EventHook 到 eventx 总线的适配。
package filex

import (
	"context"

	"github.com/lcylpzls/eventx"
	"github.com/lcylpzls/filex"
)

// Hook 返回接入 eventx 的 filex.EventHook，主题为 filex.object.<action>。
func Hook(bus *eventx.Bus) filex.EventHook {
	return filexHook{bus: bus}
}

type filexHook struct {
	bus *eventx.Bus
}

func (h filexHook) OnObjectEvent(ctx context.Context, e filex.ObjectEvent) {
	_ = h.bus.Publish(ctx, "filex.object."+e.Action, e)
}
