# 错误码手册

> eventx 全部错误使用 errx 结构化错误码；日志字段不含事件载荷。

| 错误码 | 分类 | 含义 |
| --- | --- | --- |
| `EVENTX_BUS_CLOSED` | unavailable | 总线已关闭 |
| `EVENTX_INVALID_TOPIC` | invalid | 主题非法（空/超长/控制字符/空段/非法通配符） |
| `EVENTX_INVALID_HANDLER` | invalid | 订阅处理函数非法 |
| `EVENTX_SUBSCRIPTION_NOT_FOUND` | invalid | 订阅不存在（通常发生在关闭后取消） |
| `EVENTX_HANDLER_PANIC` | internal | 订阅处理函数/过滤器发生未捕获异常 |
| `EVENTX_QUEUE_FULL` | quota_exceeded | 异步队列已满 |
| `EVENTX_INVALID_OPTION` | invalid | 选项参数非法 |
| `EVENTX_CANCELLED` | cancelled | 发布或分发被上下文取消 |
| `EVENTX_TYPE_MISMATCH` | invalid | 类型化订阅载荷类型不匹配 |

## 安全约定

- 默认日志不记录事件载荷（可能含敏感数据）；
- 审计字段只含主题、订阅者数、耗时与错误码。
