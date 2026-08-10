# 更新日志

## [v1.0.3] - 2026-08-10

### 变更

- 家族正式基线锁定：依赖统一指向 v1 基线已发布版本（errx v1.5.5 / logx v1.3.2 / testx v1.4.3 / validx v1.2.4 / cryptox v1.0.2 / confx v1.0.2 / webx v1.5.4 等），此后家族依赖不再前进。

### 质量

- 全部库包语句覆盖率保持 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v1.0.2] - 2026-08-10

### 变更

- 家族依赖最终对齐到 v1 正式版基线（errx v1.5.4 / logx v1.3.1 / testx v1.4.2 / validx v1.2.3 / confx v1.0.1 / cryptox v1.0.1 等），无 API 变更。

### 质量

- 全部库包语句覆盖率保持 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v1.0.1] - 2026-08-10

### 变更

- 家族依赖统一对齐到最新基线（errx v1.5.4 / logx v1.3.0 / testx v1.4.1 / validx v1.2.2 等），无 API 变更。

### 质量

- 全部库包语句覆盖率保持 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v1.0.0] - 2026-08-10

### 发布

- 家族正式版 v1.0.0：当前 API 与行为作为 v1 基线；按家族规则，v1.*.* 内允许破坏性修改，不承诺向后兼容。

### 质量

- 覆盖率、race、vet、staticcheck、fuzz、govulncheck 全绿；CI/Release 自动化发布。

## [v0.8.0] - 2026-08-10

### 变更

- 主题校验统一迁移至家族 `validx`：发布/订阅主题注册为
  `eventx_topic` 全局规则（参数区分 publish/subscribe），
  调用点走 `validx.ValidateField`；
- errx 错误码保持 `CodeInvalidTopic` 语义，行为不变。

### 质量

- 根包语句覆盖率保持 100%；race / vet / staticcheck / fuzz /
  govulncheck 全绿。

## [v0.7.2] - 2026-08-10

### 新增

- `eventx/adapters/authx`：authx.EventHook 适配器，
  令牌操作发布为 `authx.token.<action>` 事件。

## [v0.7.1] - 2026-08-10

### 新增

- `eventx/adapters/clix`：clix.Observer 适配器，
  命令生命周期发布为 `clix.command.start` / `clix.command.finish` 事件。

## [v0.7.0] - 2026-08-10

### 新增

- `eventx/adapters` 子模块（独立发布）：filex / cachex / jobx / dbx
  官方事件适配器，把各库 EventHook 桥接到总线；
- 端到端联动示例 examples/ecosystem：filex 上传 → eventx →
  cachex 失效 + jobx 异步处理；
- CI/Release 覆盖 adapters 子模块与子模块 tag。

### 质量

- 根包语句覆盖率 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v0.6.4] - 2026-08-10

### 变更

- go 指令与 CI/Release 工作流统一为 Go 1.26.5；
- README Go 版本徽章同步更新。

## [v0.6.3] - 2026-08-10

### 收口

- 终审清单全部通过，eventx 达到 1.0 候选标准；
- **v1.0.0 是否发布由维护者决定**；
- 1.0 确认后按 `docs/adapters.md` 规范开展家族接入。

### 质量

- 根包语句覆盖率 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v0.6.2] - 2026-08-10

### 修复

- 零值 Bus（未经 New 构造）的 Subscribe / Publish / PublishAsync /
  Close 统一返回 `EVENTX_INVALID_OPTION`，不再 panic；
- Metrics 对零值总线返回全 0 快照。

### 质量

- 根包语句覆盖率 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v0.6.1] - 2026-08-10

### 修复

- `Publish` / `PublishAsync` 对 nil context 归一为 Background，不再 panic；
- `SubscribeTyped` 对 nil 总线返回 `EVENTX_INVALID_OPTION`。

### 质量

- 根包语句覆盖率 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v0.6.0] - 2026-08-10

### 新增

- `docs/errors.md` 错误码手册（9 个错误码）；
- `docs/adapters.md` 家族接入规范：EventHook 零依赖钩子 +
  `eventx/adapters` 适配器 + 项目层组装；
- `docs/final-review.md` 1.0 候选终审清单；
- Issue 模板与 README CI 徽章/状态行。

### 质量

- 根包语句覆盖率 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v0.5.0] - 2026-08-10

### 新增

- logx 审计：`WithLogger` 注入后记录分发完成（主题/订阅数/耗时/错误码，
  载荷不记录）；
- 运行指标：`Metrics` 返回发布/投递/失败/订阅数快照，
  供 metricsx 等外部适配；
- 基准测试与 `BENCHMARKS.md`（同步发布 ~216ns、异步投递 ~185ns、
  通配匹配 ~48ns）。

### 质量

- 根包语句覆盖率 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

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
