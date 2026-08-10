package eventx

import (
	"context"
	"sort"
	"strings"
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

// Handler 是订阅处理函数；返回错误会被聚合。
type Handler func(ctx context.Context, e Event) error

// Filter 是订阅过滤器；返回 false 时跳过该订阅者。
type Filter func(ctx context.Context, e Event) bool

// Subscription 是订阅句柄，可查询与取消。
type Subscription interface {
	// ID 返回订阅唯一标识。
	ID() uint64
	// Topic 返回订阅主题（可能含通配符）。
	Topic() string
	// Unsubscribe 取消订阅；重复取消幂等返回 nil。
	Unsubscribe() error
}

// SubscribeOption 修改订阅配置。
type SubscribeOption func(*subscribeConfig)

// subscribeConfig 是订阅配置。
type subscribeConfig struct {
	filter   Filter
	priority int
}

// WithPriority 设置订阅优先级：数值越小越先执行（默认 0，稳定排序）。
func WithPriority(p int) SubscribeOption {
	return func(c *subscribeConfig) {
		c.priority = p
	}
}

// Bus 是并发安全的事件总线，无全局单例。
type Bus struct {
	mu     sync.RWMutex
	subs   map[string][]*subscription
	wild   []*subscription
	seq    uint64
	closed bool
	async  *asyncDispatcher
}

// subscription 是 Subscription 的内部实现。
type subscription struct {
	id       uint64
	topic    string
	pattern  *pattern
	handler  Handler
	filter   Filter
	priority int
	bus      *Bus
	closed   atomic.Bool
}

// Option 修改总线构造配置。
type Option func(*config)

// config 是总线构造配置。
type config struct {
	workers   int
	queueSize int
	onError   func(error)
}

// WithWorkers 设置异步 worker 数（默认 1）。
func WithWorkers(n int) Option {
	return func(c *config) {
		c.workers = n
	}
}

// WithQueueSize 设置异步队列容量（默认 1024）。
func WithQueueSize(n int) Option {
	return func(c *config) {
		c.queueSize = n
	}
}

// WithErrorHandler 设置异步发布错误回调；
// 未设置时异步 handler 错误被忽略（建议生产环境设置）。
func WithErrorHandler(fn func(error)) Option {
	return func(c *config) {
		c.onError = fn
	}
}

// New 创建事件总线。选项非法时返回 errx 错误。
func New(opts ...Option) (*Bus, error) {
	cfg := config{workers: 1, queueSize: 1024}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.workers < 1 {
		return nil, errx.NewCodef(CodeInvalidOption, "worker 数必须 >= 1，当前 %d", cfg.workers)
	}
	if cfg.queueSize < 1 {
		return nil, errx.NewCodef(CodeInvalidOption, "队列容量必须 >= 1，当前 %d", cfg.queueSize)
	}
	b := &Bus{subs: make(map[string][]*subscription)}
	b.async = newAsyncDispatcher(cfg)
	b.async.start(b)
	return b, nil
}

// Subscribe 按主题注册订阅（支持 `*` / `**` 通配符）。
func (b *Bus) Subscribe(topic string, handler Handler) (Subscription, error) {
	return b.subscribe(topic, handler, subscribeConfig{})
}

// SubscribeFiltered 按主题注册带过滤器的订阅；filter 为 nil 时不过滤。
func (b *Bus) SubscribeFiltered(topic string, filter Filter, handler Handler) (Subscription, error) {
	return b.subscribe(topic, handler, subscribeConfig{filter: filter})
}

// SubscribeWithOptions 按主题注册订阅并应用订阅选项（优先级等）。
func (b *Bus) SubscribeWithOptions(topic string, handler Handler, opts ...SubscribeOption) (Subscription, error) {
	cfg := subscribeConfig{}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return b.subscribe(topic, handler, cfg)
}

// subscribe 统一注册订阅。
func (b *Bus) subscribe(topic string, handler Handler, cfg subscribeConfig) (Subscription, error) {
	if err := validateSubscribeTopic(topic); err != nil {
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
	sub := &subscription{
		id:       b.seq,
		topic:    topic,
		pattern:  compilePattern(topic),
		handler:  handler,
		filter:   cfg.filter,
		priority: cfg.priority,
		bus:      b,
	}
	if strings.Contains(topic, "*") {
		b.wild = append(b.wild, sub)
	} else {
		b.subs[topic] = append(b.subs[topic], sub)
	}
	return sub, nil
}

// Publish 同步发布事件：按优先级与注册顺序调用全部匹配 handler，
// 聚合所有返回的错误（errx.Join）后返回。
func (b *Bus) Publish(ctx context.Context, topic string, payload any) error {
	if err := validatePublishTopic(topic); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return errx.WrapCode(err, CodeCancelled, "发布上下文已取消")
	}
	if b.isClosed() {
		return errx.NewCode(CodeBusClosed, "总线已关闭")
	}
	return deliver(b, ctx, topic, payload, b.snapshot(topic))
}

// Close 关闭总线：拒绝新发布/订阅，排空异步在途任务后清空订阅。
// 重复关闭幂等。
func (b *Bus) Close() error {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return nil
	}
	b.closed = true
	b.mu.Unlock()
	// 先排空异步在途任务（订阅表仍完整），再清空订阅。
	b.async.stopAndWait()
	b.mu.Lock()
	b.subs = make(map[string][]*subscription)
	b.wild = nil
	b.mu.Unlock()
	return nil
}

// snapshot 返回主题全部匹配订阅（精确 + 通配符），
// 按优先级升序、同优先级按注册顺序稳定排序。
func (b *Bus) snapshot(topic string) []*subscription {
	b.mu.RLock()
	defer b.mu.RUnlock()
	subs := append([]*subscription(nil), b.subs[topic]...)
	for _, sub := range b.wild {
		if sub.pattern.matches(topic) {
			subs = append(subs, sub)
		}
	}
	if len(subs) > 1 {
		sort.SliceStable(subs, func(i, j int) bool {
			if subs[i].priority != subs[j].priority {
				return subs[i].priority < subs[j].priority
			}
			return subs[i].id < subs[j].id
		})
	}
	return subs
}

// isClosed 返回总线是否已关闭。
func (b *Bus) isClosed() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.closed
}

// deliver 按快照顺序执行订阅（过滤器 + handler），聚合错误。
func deliver(b *Bus, ctx context.Context, topic string, payload any, subs []*subscription) error {
	var errs []error
	for _, sub := range subs {
		if err := ctx.Err(); err != nil {
			return errx.WrapCode(err, CodeCancelled, "发布上下文已取消")
		}
		if sub.closed.Load() {
			continue
		}
		e := Event{Topic: topic, Payload: payload}
		if sub.filter != nil {
			ok, ferr := applyFilter(sub.filter, ctx, e)
			if ferr != nil {
				errs = append(errs, ferr)
				continue
			}
			if !ok {
				continue
			}
		}
		if err := invoke(sub.handler, ctx, e); err != nil {
			errs = append(errs, err)
		}
	}
	return errx.Join(errs...)
}

// applyFilter 调用过滤器并恢复 panic。
func applyFilter(f Filter, ctx context.Context, e Event) (ok bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
			err = errx.NewCodef(CodeHandlerPanic, "订阅过滤器发生未捕获异常：%v", r)
		}
	}()
	return f(ctx, e), nil
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
	if list, ok := removeSub(s.bus.subs[s.topic], s.id); ok {
		if len(list) == 0 {
			delete(s.bus.subs, s.topic)
		} else {
			s.bus.subs[s.topic] = list
		}
		return nil
	}
	var ok bool
	s.bus.wild, ok = removeSub(s.bus.wild, s.id)
	if ok {
		return nil
	}
	return errx.NewCode(CodeSubscriptionNotFound, "订阅不存在")
}

// removeSub 从列表移除指定 id 的订阅。
func removeSub(list []*subscription, id uint64) ([]*subscription, bool) {
	for i, sub := range list {
		if sub.id == id {
			return append(list[:i], list[i+1:]...), true
		}
	}
	return list, false
}
