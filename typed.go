package eventx

import (
	"context"

	"github.com/lcylpzls/errx"
)

// TypedHandler 是类型化订阅处理函数：载荷已断言为 T。
type TypedHandler[T any] func(ctx context.Context, topic string, payload T) error

// SubscribeTyped 订阅类型化事件：发布载荷类型不匹配时返回
// EVENTX_TYPE_MISMATCH（该订阅者失败，不影响其他订阅者）。
func SubscribeTyped[T any](b *Bus, topic string, handler TypedHandler[T]) (Subscription, error) {
	if b == nil {
		return nil, errx.NewCode(CodeInvalidOption, "总线不能为空")
	}
	if handler == nil {
		return nil, errx.NewCode(CodeInvalidHandler, "类型化订阅处理函数不能为空")
	}
	return b.Subscribe(topic, func(ctx context.Context, e Event) error {
		p, ok := e.Payload.(T)
		if !ok {
			return errx.NewCodef(CodeTypeMismatch,
				"载荷类型不匹配：期望 %T，实际 %T", *new(T), e.Payload)
		}
		return handler(ctx, e.Topic, p)
	})
}
