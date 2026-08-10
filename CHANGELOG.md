# 更新日志

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
