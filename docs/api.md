# API 快照

> 随版本更新。v0.1.0 快照如下；新版本发布后同步替换。

## v0.1.0

### 类型

```go
type Event struct {
    Topic   string
    Payload any
}
type Handler func(ctx context.Context, e Event) error
type Subscription interface {
    ID() uint64
    Topic() string
    Unsubscribe() error
}
```

### 总线

```go
func New() *Bus
func (b *Bus) Subscribe(topic string, handler Handler) (Subscription, error)
func (b *Bus) Publish(ctx context.Context, topic string, payload any) error
func (b *Bus) Close() error
```

- 同一主题按注册顺序同步执行；
- handler 错误被 errx.Join 聚合返回；
- handler panic 恢复为 `EVENTX_HANDLER_PANIC`；
- 关闭后 Subscribe / Publish 返回 `EVENTX_BUS_CLOSED`；
- 取消订阅幂等。

### 主题

- 点分字符串（如 `orders.created`），最长 256 字节；
- 禁止空、控制字符与通配符（通配符订阅 v0.3.0 起）。

### 错误码

| 错误码 | 分类 | 含义 |
| --- | --- | --- |
| `EVENTX_BUS_CLOSED` | unavailable | 总线已关闭 |
| `EVENTX_INVALID_TOPIC` | invalid | 主题非法 |
| `EVENTX_INVALID_HANDLER` | invalid | 订阅处理函数非法 |
| `EVENTX_SUBSCRIPTION_NOT_FOUND` | invalid | 订阅不存在 |
| `EVENTX_HANDLER_PANIC` | internal | 订阅处理函数发生未捕获异常 |
