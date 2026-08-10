# API 快照

> 随版本更新。v0.6.2 快照如下；新版本发布后同步替换。

## v0.6.2

### 类型

```go
type Event struct {
    Topic   string
    Payload any
}
type Handler func(ctx context.Context, e Event) error
type Filter func(ctx context.Context, e Event) bool
type Subscription interface {
    ID() uint64
    Topic() string
    Unsubscribe() error
}
```

### 总线

```go
func New(opts ...Option) (*Bus, error)
func (b *Bus) Subscribe(topic string, handler Handler) (Subscription, error)
func (b *Bus) SubscribeFiltered(topic string, filter Filter, handler Handler) (Subscription, error)
func (b *Bus) SubscribeWithOptions(topic string, handler Handler, opts ...SubscribeOption) (Subscription, error)
func (b *Bus) Publish(ctx context.Context, topic string, payload any) error
func (b *Bus) PublishAsync(ctx context.Context, topic string, payload any) error
func (b *Bus) Close() error
```

- 同一主题按注册顺序同步执行；
- handler 错误被 errx.Join 聚合返回；
- handler panic 恢复为 `EVENTX_HANDLER_PANIC`；
- `PublishAsync` 非阻塞入队，队列满返回 `EVENTX_QUEUE_FULL`；
- `Close` 优雅关闭：拒绝新投递，排空在途任务后清空订阅；
- 关闭后 Subscribe / Publish 返回 `EVENTX_BUS_CLOSED`；
- 取消订阅幂等。

### 通配符、过滤器与优先级

```go
func WithPriority(p int) SubscribeOption
```

- 订阅主题支持 `*`（匹配单段）与 `**`（匹配零或多段），
  通配符必须作为独立段出现；
- 发布主题禁止通配符与空段；
- 过滤器返回 false 时跳过该订阅者；过滤器 panic 恢复为
  `EVENTX_HANDLER_PANIC`；
- 执行顺序：优先级升序，同优先级按注册顺序（稳定）。

### 上下文与类型化订阅

```go
type TypedHandler[T any] func(ctx context.Context, topic string, payload T) error
func SubscribeTyped[T any](b *Bus, topic string, handler TypedHandler[T]) (Subscription, error)
```

- 发布入口与分发循环都响应 context 取消：取消时返回
  `EVENTX_CANCELLED`，剩余订阅者不再执行；
- nil context 归一为 context.Background()；
- 未初始化（零值）总线的操作返回 `EVENTX_INVALID_OPTION`；
- `SubscribeTyped` 自动断言载荷类型，不匹配返回
  `EVENTX_TYPE_MISMATCH`（不影响其他订阅者）；
- 类型化订阅处理函数为 nil 时返回 `EVENTX_INVALID_HANDLER`。

### 选项

```go
func WithWorkers(n int) Option
func WithQueueSize(n int) Option
func WithErrorHandler(fn func(error)) Option
func WithLogger(logger logx.Logger) Option
```

- worker 数默认 1，队列容量默认 1024；
- 异步 handler 错误通过 `WithErrorHandler` 回调上报；
- worker/队列容量非法时 `New` 返回 `EVENTX_INVALID_OPTION`。

### 可观测

```go
type Metrics struct {
    Publishes     uint64
    Deliveries    uint64
    Failures      uint64
    Subscriptions uint64
}
func (b *Bus) Metrics() Metrics
```

- `WithLogger` 注入 logx.Logger：分发完成记录
  topic / subscribers / duration_ms / 错误码，载荷不记录；
- `Metrics` 返回发布/投递/失败/订阅数快照，供 metricsx 等外部适配。

### 主题

- 点分字符串（如 `orders.created`），最长 256 字节；
- 禁止空、控制字符与空段；订阅可含独立段通配符 `*` / `**`。

### 错误码

| 错误码 | 分类 | 含义 |
| --- | --- | --- |
| `EVENTX_BUS_CLOSED` | unavailable | 总线已关闭 |
| `EVENTX_INVALID_TOPIC` | invalid | 主题非法 |
| `EVENTX_INVALID_HANDLER` | invalid | 订阅处理函数非法 |
| `EVENTX_SUBSCRIPTION_NOT_FOUND` | invalid | 订阅不存在 |
| `EVENTX_HANDLER_PANIC` | internal | 订阅处理函数发生未捕获异常 |
| `EVENTX_QUEUE_FULL` | quota_exceeded | 异步队列已满 |
| `EVENTX_INVALID_OPTION` | invalid | 选项参数非法 |
| `EVENTX_CANCELLED` | cancelled | 发布或分发被上下文取消 |
| `EVENTX_TYPE_MISMATCH` | invalid | 类型化订阅载荷类型不匹配 |
