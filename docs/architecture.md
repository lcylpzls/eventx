# 架构

## 1. 包内模块

```text
eventx（根包）
├── errors.go      错误码定义与注册
├── topic.go       主题校验与通配符匹配（v0.1.0/v0.3.0）
├── bus.go         总线：订阅注册、同步发布（v0.1.0）
├── async.go       异步投递、worker 池、优雅关闭（v0.2.0）
├── match.go       过滤器与优先级排序（v0.3.0）
├── typed.go       上下文贯穿与类型化订阅（v0.4.0）
├── audit.go       logx 审计字段与 Metrics（v0.5.0）
```

依赖方向：

```text
audit.go ──→ typed.go ──→ bus.go ──→ topic.go ──→ errors.go
```

## 2. 关键设计

- **注册表并发安全**：订阅表由 RWMutex 保护，发布时持读锁快照执行；
- **同步语义**：Publish 按订阅顺序同步调用，错误聚合后返回；
- **异步语义**：PublishAsync 入队后立即返回，worker 池消费，
  Close 等待在途投递完成；
- **匹配编译**：通配符订阅编译为段模式，避免每次发布做正则；
- **零魔法**：无全局单例；Bus 实例即全部状态。

## 3. 后续演进扩展点

- 中间件链：发布前/后钩子（按需）；
- 分布式形态：独立 mqx 评估，不在 eventx 内实现。
