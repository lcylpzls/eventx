// Package clix 提供 clix.Observer 到 eventx 总线的适配。
package clix

import (
	"context"
	"time"

	"github.com/lcylpzls/clix"
	"github.com/lcylpzls/eventx"
)

// CommandEvent 描述一次命令生命周期事件。
type CommandEvent struct {
	// Command 命令完整路径（根命令为"（根命令）"）。
	Command string
	// Args 命令参数。
	Args []string
	// Err 执行结果错误；nil 表示成功。
	Err error
	// Duration 执行耗时。
	Duration time.Duration
}

// Observer 返回接入 eventx 的 clix.Observer，
// 主题为 clix.command.start / clix.command.finish。
func Observer(bus *eventx.Bus) clix.Observer {
	return clixObserver{bus: bus}
}

type clixObserver struct {
	bus *eventx.Bus
}

func (o clixObserver) OnCommandStart(ctx context.Context, command string, args []string) {
	_ = o.bus.Publish(ctx, "clix.command.start", CommandEvent{
		Command: command,
		Args:    append([]string(nil), args...),
	})
}

func (o clixObserver) OnCommandFinish(ctx context.Context, command string, args []string, err error, duration time.Duration) {
	_ = o.bus.Publish(ctx, "clix.command.finish", CommandEvent{
		Command:  command,
		Args:     append([]string(nil), args...),
		Err:      err,
		Duration: duration,
	})
}
