# 基准测试

> 采集环境：Windows / AMD Ryzen 5 7600 / Go 1.26.5
> 采集日期：2026-08-10
> 命令：`go test -bench=. -benchmem -run '^$' .`

## BenchmarkPublishSync

10 个订阅者的同步发布：

| 指标 | 数值 |
| --- | --- |
| 耗时 | 215.8 ns/op |
| 内存 | 176 B/op |
| 分配 | 4 allocs/op |

## BenchmarkPublishAsync

4 worker、大队列异步投递（不含订阅者执行）：

| 指标 | 数值 |
| --- | --- |
| 耗时 | 185.1 ns/op |
| 内存 | 39 B/op |
| 分配 | 2 allocs/op |

## BenchmarkWildcardMatch

`orders.**` 模式匹配 4 段发布主题：

| 指标 | 数值 |
| --- | --- |
| 耗时 | 47.7 ns/op |
| 内存 | 64 B/op |
| 分配 | 1 allocs/op |

## 说明

- 基准仅反映本机相对量级；CI 不设硬性性能门槛；
- 同步发布按订阅者数线性增长，异步投递路径保持轻量。
