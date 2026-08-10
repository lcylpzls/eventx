# 1.0 候选终审

> 本清单用于确认 eventx 达到 1.0 候选标准；**v1.0.0 是否发布由维护者决定**。

## 1. API 冻结

- [ ] 公开 API 签名稳定：Bus/Event/Handler/Subscription/选项语义明确；
- [ ] 主题匹配规则（`*` / `**`）与错误码全集作为稳定约定；
- [ ] pre-1.0 兼容承诺：v0.1.0 起的核心行为无意外破坏。

## 2. 质量门禁

- [ ] 根包语句覆盖率 100%；
- [ ] 测试乱序（`-shuffle=on`）、race 全平台通过；
- [ ] vet / staticcheck / govulncheck 通过；
- [ ] fuzz 目标（主题校验）短跑 5s 通过；
- [ ] 示例模块（同步/异步/通配/过滤/类型化）全绿。

## 3. 设计确认

- [ ] 并发安全：订阅表 RWMutex、快照执行、异步优雅关闭；
- [ ] 发布入口与分发循环响应上下文取消；
- [ ] handler/filter panic 恢复为结构化错误；
- [ ] 事件载荷默认不进入日志；
- [ ] 无全局单例，Bus 实例即全部状态。

## 4. 性能

- [ ] `BENCHMARKS.md` 记录同步/异步/匹配基准。

## 5. 文档与安全

- [ ] README / docs/api.md / docs/errors.md / docs/roadmap.md 一致；
- [ ] SECURITY.md / CONTRIBUTING.md / CODEOWNERS / Issue 模板齐全；
- [ ] 家族接入规范（docs/adapters.md）定稿；
- [ ] 发布流程：tag 触发 Release，CI 全绿后发布。

## 结论

填表完成后更新此处：是否达到 1.0 候选，以及待维护者确认的发布项。
