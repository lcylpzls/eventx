# 家族接入规范（EventHook 模式）

eventx 采用"库定义零依赖钩子、适配器可插拔、项目层组装"的家族模式。
以下规范供 filex / dbx / authx 等生产者库接入时参考。

## 1. 库侧：定义零依赖 EventHook

```go
// 在库自己的包内定义，不 import eventx。
type EventHook interface {
    OnEvent(ctx context.Context, event Event) error
}
```

规则：

- 钩子必须是**可选字段**（默认 no-op），不设置时库行为与现在完全一致；
- 事件模型由库定义（主题 + 结构化载荷），不依赖 eventx 类型；
- 钩子调用失败不得影响库主流程（记录并继续，或按库语义处理）。

## 2. eventx 侧：提供适配器

```go
// eventx/adapters/filex 子包（独立子模块发布）。
func Hook(bus *eventx.Bus) filex.EventHook {
    return filexHook{bus: bus}
}
```

适配器把库事件映射为 eventx 主题（如 `filex.object.put`）并发布。

## 3. 项目层组装

```go
store, _ := filex.New(filex.WithEventHook(filexadapter.Hook(bus)))
_, _ = eventx.SubscribeTyped[filex.ObjectEvent](bus, "filex.object.put", handler)
```

不传钩子即完全解耦，零成本。

## 4. 接入检查单

- [ ] 库 API 冻结前完成 EventHook 设计（0.x 库在 1.0 定版时一并设计）；
- [ ] 钩子零依赖、可选、并发安全；
- [ ] 主题命名遵循 `库名.领域.动作`（如 `filex.object.put`）；
- [ ] 适配器独立子模块发布（如 `eventx/adapters/filex`）；
- [ ] 接入后同步更新 tracex/metricsx 链路与审计。
