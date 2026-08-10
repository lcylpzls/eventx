# 设计

## 1. 定位与范围

eventx 是家族生态的**进程内事件总线基座**：

- 提供主题发布-订阅（同步/异步）、通配符匹配、过滤器与优先级；
- 与 errx（错误码）、logx（审计）、tracex（上下文）打通；
- 家族库通过零依赖 EventHook + `eventx/adapters` 可插拔接入；
- 项目可直接使用 eventx 做模块间解耦。

非目标：

- 不做分布式消息中间件抽象（跨进程需求另立 mqx）；
- 不内置持久化、重放与死信（进程内形态）；
- 不要求任何家族库强制依赖 eventx。

## 2. 核心模型

```text
eventx
├── Bus：线程安全的事件总线
├── Event：{ Topic string, Payload any }
├── Handler：func(ctx context.Context, e Event) error
├── Subscribe / Unsubscribe：订阅管理（v0.1.0）
├── Publish：同步发布，聚合订阅者错误（v0.1.0）
├── PublishAsync：异步投递 + worker 池 + 优雅关闭（v0.2.0）
├── 通配符主题 / Filter / Priority（v0.3.0）
├── 上下文贯穿 / SubscribeTyped（v0.4.0）
└── logx 审计 / Metrics / 基准（v0.5.0）
```

## 3. 主题与匹配

- 主题为点分字符串（如 `orders.created`、`filex.object.put`）；
- 发布主题禁止使用通配符；订阅主题支持 `*`（单段）与 `**`（多段）；
- 订阅匹配在注册时编译，发布时 O(段数) 查表；
- 同一主题的订阅按注册顺序执行（v0.1/v0.2），
  优先级排序自 v0.3.0 起。

## 4. 错误与可观测

- 同步发布收集全部订阅者错误并 errx.Join 聚合返回；
- 异步发布通过 `OnError` 回调上报（不阻塞投递）；
- 已关闭总线的发布/订阅返回 `EVENTX_BUS_CLOSED`；
- 日志只含 topic/handler/耗时/错误码，事件载荷默认不记录。

## 5. 版本与兼容

- 语义化版本；pre-1.0 允许按路线图演进；
- 每版完成即发布 tag，CI 全绿后 Release 自动生成；
- v1.0.0 是否发布由维护者决定，eventx 只推进到 1.0 候选即停；
- 家族接入（EventHook + adapters）在 1.0 确认后按需求开展。
