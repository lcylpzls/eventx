# 1.0 候选终审

> 本清单用于确认 eventx 达到 1.0 候选标准；**v1.0.0 是否发布由维护者决定**。

## 1. API 冻结

- [x] 公开 API 签名稳定：Bus/Event/Handler/Subscription/选项语义明确；
- [x] 主题匹配规则（`*` / `**`）与错误码全集作为稳定约定；
- [x] pre-1.0 兼容承诺：v0.1.0 起的核心行为无意外破坏。

## 2. 质量门禁

- [x] 根包语句覆盖率 100%；
- [x] 测试乱序（`-shuffle=on`）、race 全平台通过；
- [x] vet / staticcheck / govulncheck 通过；
- [x] fuzz 目标（主题校验）短跑 5s 通过；
- [x] 示例模块（同步/异步/通配/过滤/类型化）全绿。

## 3. 设计确认

- [x] 并发安全：订阅表 RWMutex、快照执行、异步优雅关闭；
- [x] 发布入口与分发循环响应上下文取消；
- [x] handler/filter panic 恢复为结构化错误；
- [x] 事件载荷默认不进入日志；
- [x] 无全局单例，Bus 实例即全部状态。

## 4. 性能

- [x] `BENCHMARKS.md` 记录同步/异步/匹配基准。

## 5. 文档与安全

- [x] README / docs/api.md / docs/errors.md / docs/roadmap.md 一致；
- [x] SECURITY.md / CONTRIBUTING.md / CODEOWNERS / Issue 模板齐全；
- [x] 家族接入规范（docs/adapters.md）定稿；
- [x] 发布流程：tag 触发 Release，CI 全绿后发布。

## 结论

eventx 已通过 1.0 候选终审清单，达到 1.0 候选标准。
**v1.0.0 是否发布由维护者决定**；确认发布前不再自动推进版本。
1.0 确认后，按 docs/adapters.md 规范开展家族 EventHook + adapters 接入。
