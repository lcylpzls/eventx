// Package authx 提供 authx.EventHook 到 eventx 总线的适配。
package authx

import (
	"context"

	"github.com/lcylpzls/authx"
	"github.com/lcylpzls/eventx"
)

// Hook 返回接入 eventx 的 authx.EventHook，主题为 authx.token.<action>。
func Hook(bus *eventx.Bus) authx.EventHook {
	return authxHook{bus: bus}
}

type authxHook struct {
	bus *eventx.Bus
}

func (h authxHook) OnAuthEvent(ctx context.Context, e authx.AuthEvent) {
	_ = h.bus.Publish(ctx, "authx.token."+e.Action, e)
}
