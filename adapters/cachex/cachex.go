// Package cachex 提供 cachex.EventHook 到 eventx 总线的适配。
package cachex

import (
	"context"

	"github.com/lcylpzls/cachex"
	"github.com/lcylpzls/eventx"
)

// Hook 返回接入 eventx 的 cachex.EventHook，主题为 cachex.cache.<action>。
func Hook(bus *eventx.Bus) cachex.EventHook {
	return cachexHook{bus: bus}
}

type cachexHook struct {
	bus *eventx.Bus
}

func (h cachexHook) OnCacheEvent(ctx context.Context, e cachex.CacheEvent) {
	_ = h.bus.Publish(ctx, "cachex.cache."+e.Action, e)
}
