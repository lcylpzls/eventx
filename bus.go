package eventx

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/lcylpzls/errx"
)

// Event 是总线传递的事件。
type Event struct {
	// Topic 发布主题。
	Topic string
	// Payload 事件载荷。
	Payload any
}

// Handler 是订阅处理函数；返回错误会被 Publish 聚合。
type Handler func(ctx context.Context, e Event) error

// Subscription 是订阅句柄，可查询与取消。
type Subscription interface {
	// ID 返回订阅唯一标识。
	ID() uint64
	// Topic 返回订阅主题。
	Topic() string
	// Unsubscribe 取消订阅；重复取消幂等返回 nil。
	Unsubscribe() error
}

// Bus 是并发安全的事件总线，无全局单例。
type Bus struct {
	mu     sync.RWMutex
	subs   map[string][]*subscription
	seq    uint64
	closed bool
}

// subscription 是 Subscription 的内部实现。
type subscription struct {
	id      uint64
	topic   string
	handler Handler
	bus     *Bus
	closed  atomic.Bool
}

// New 创建事件总线。
func New() *Bus {
	return &Bus{subs: make(map[string][]*subscription)}
}

// Subscribe 按主题注册订阅，返回订阅句柄。
// 同一主题的订阅按注册顺序执行。
func (b *Bus) Subscribe(topic string, handler Handler) (Subscription, error) {
	if err := validateTopic(topic); err != nil {
		return nil, err
	}
	if handler == nil {
		return nil, errx.NewCode(CodeInvalidHandler, "订阅处理函数不能为空")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return nil, errx.NewCode(CodeBusClosed, "总线已关闭")
	}
	b.seq++
	sub := &subscription{id: b.seq, topic: topic, handler: handler, bus: b}
	b.subs[topic] = append(b.subs[topic], sub)
	return sub, nil
}

// Publish 同步发布事件：按订阅顺序调用全部匹配 handler，
// 聚合所有返回的错误（errx.Join）后返回。
// handler 的 panic 会被恢复并转为 EVENTX_HANDLER_PANIC。
func (b *Bus) Publish(ctx context.Context, topic string, payload any) error {
	if err := validateTopic(topic); err != nil {
		return err
	}
	if b.isClosed() {
		return errx.NewCode(CodeBusClosed, "总线已关闭")
	}
	subs := b.snapshot(topic)
	if len(subs) == 0 {
		return nil
	}
	var errs []error
	for _, sub := range subs {
		if sub.closed.Load() {
			continue
		}
		if err := invoke(sub.handler, ctx, Event{Topic: topic, Payload: payload}); err != nil {
			errs = append(errs, err)
		}
	}
	return errx.Join(errs...)
}

// Close 关闭总线：关闭后 Subscribe / Publish 返回 EVENTX_BUS_CLOSED。
// 重复关闭幂等。异步在途投递的等待由 v0.2.0 的队列语义补充。
func (b *Bus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return nil
	}
	b.closed = true
	b.subs = make(map[string][]*subscription)
	return nil
}

// snapshot 返回主题订阅的快照副本，避免执行 handler 时持锁。
func (b *Bus) snapshot(topic string) []*subscription {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return append([]*subscription(nil), b.subs[topic]...)
}

// isClosed 返回总线是否已关闭。
func (b *Bus) isClosed() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.closed
}

// invoke 调用 handler 并恢复 panic。
func invoke(h Handler, ctx context.Context, e Event) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errx.NewCodef(CodeHandlerPanic, "订阅处理函数发生未捕获异常：%v", r)
		}
	}()
	return h(ctx, e)
}

// ID 返回订阅唯一标识。
func (s *subscription) ID() uint64 {
	return s.id
}

// Topic 返回订阅主题。
func (s *subscription) Topic() string {
	return s.topic
}

// Unsubscribe 取消订阅；重复取消幂等返回 nil。
func (s *subscription) Unsubscribe() error {
	if s.closed.Swap(true) {
		return nil
	}
	s.bus.mu.Lock()
	defer s.bus.mu.Unlock()
	list := s.bus.subs[s.topic]
	for i, sub := range list {
		if sub.id == s.id {
			s.bus.subs[s.topic] = append(list[:i], list[i+1:]...)
			if len(s.bus.subs[s.topic]) == 0 {
				delete(s.bus.subs, s.topic)
			}
			return nil
		}
	}
	return errx.NewCode(CodeSubscriptionNotFound, "订阅不存在")
}
