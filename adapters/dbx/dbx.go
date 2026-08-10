// Package dbx 提供 dbx.EventHook 到 eventx 总线的适配。
package dbx

import (
	"context"

	"github.com/lcylpzls/dbx"
	"github.com/lcylpzls/eventx"
)

// Hook 返回接入 eventx 的 dbx.EventHook，主题为 dbx.query.<operation>。
func Hook(bus *eventx.Bus) dbx.EventHook {
	return dbxHook{bus: bus}
}

type dbxHook struct {
	bus *eventx.Bus
}

func (h dbxHook) OnQueryEvent(ctx context.Context, e dbx.QueryEvent) {
	_ = h.bus.Publish(ctx, "dbx.query."+e.Operation, e)
}
