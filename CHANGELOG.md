# 更新日志

## [v0.4.0] - 2026-08-10

### 新增

- 上下文取消贯穿：Publish / PublishAsync 入口与分发循环响应
  context 取消，返回 `EVENTX_CANCELLED` 并停止剩余订阅者；
- 类型化订阅：`SubscribeTyped[T]` 自动断言载荷类型，
  不匹配返回 `EVENTX_TYPE_MISMATCH`；
- 错误码：`EVENTX_CANCELLED` / `EVENTX_TYPE_MISMATCH`。

### 质量

- 根包语句覆盖率 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v0.3.0] - 2026-08-10

### 新增

- 通配符订阅：`*` 单段、`**` 多段/零段，编译为段模式匹配；
- 过滤器：`SubscribeFiltered`，false 跳过订阅者，panic 恢复；
- 订阅优先级：`SubscribeWithOptions` + `WithPriority`，
  优先级升序、同优先级按注册顺序稳定执行；
- 主题校验收紧：发布/订阅均禁止空段；订阅通配符必须独立段；
- 精确/通配符订阅统一快照排序，异步与同步共用同一分发逻辑。

### 质量

- 根包语句覆盖率 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v0.2.0] - 2026-08-10

### 新增

- 异步发布：`PublishAsync` 非阻塞入队，worker 池并发消费；
- 选项：`WithWorkers` / `WithQueueSize` / `WithErrorHandler`；
- `New` 变更为 `New(opts ...Option) (*Bus, error)`（pre-1.0 破坏性变更）；
- 优雅关闭：`Close` 排空在途任务后清空订阅（先排空、后清表）；
- 错误码：`EVENTX_QUEUE_FULL` / `EVENTX_INVALID_OPTION`；
- 修复：异步入队与关闭之间的竞态（读锁内完成检查与入队）。

### 质量

- 根包语句覆盖率 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v0.1.0] - 2026-08-10

### 新增

- 同步发布-订阅核心：`New` / `Subscribe` / `Unsubscribe` / `Publish`；
- 主题校验（非空、长度上限、禁止控制字符）；
- 同步发布聚合全部订阅者错误（errx.Join）；
- errx 错误码：`EVENTX_BUS_CLOSED` / `EVENTX_INVALID_TOPIC` /
  `EVENTX_INVALID_HANDLER` / `EVENTX_SUBSCRIPTION_NOT_FOUND`；
- fuzz 目标（`FuzzTopic`）接入 CI；
- 三平台 CI + Linux 多发行版容器矩阵 + Release 工作流。

### 质量

- 根包语句覆盖率 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。
