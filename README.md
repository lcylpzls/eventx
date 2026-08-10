# eventx

进程内事件总线基座：主题发布-订阅（同步/异步）、通配符匹配、
过滤器与优先级，与 errx / logx / tracex 家族打通。

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.26-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![CI](https://github.com/lcylpzls/eventx/actions/workflows/ci.yml/badge.svg)](https://github.com/lcylpzls/eventx/actions/workflows/ci.yml)

## 快速开始

```go
package main

import (
	"context"
	"fmt"

	"github.com/lcylpzls/eventx"
)

func main() {
	bus, err := eventx.New()
	if err != nil {
		panic(err)
	}
	sub, _ := bus.Subscribe("orders.created", func(ctx context.Context, e eventx.Event) error {
		fmt.Printf("收到订单：%v\n", e.Payload)
		return nil
	})
	defer sub.Unsubscribe()

	_ = bus.Publish(context.Background(), "orders.created", "10086")
}
```

## 核心特性

- 同步发布：按订阅顺序执行，错误 errx.Join 聚合返回；
- 异步发布：非阻塞入队 + worker 池，优雅关闭排空在途投递；
- 通配符主题：`*` 单段、`**` 多段；
- 过滤器与订阅优先级；
- context.Context 贯穿，tracex 可直接注入；
- 类型化订阅：`SubscribeTyped[T]`；
- 可选 logx 审计与 Metrics 接口；
- 并发安全、无全局单例。

> 当前状态：**v0.6.x（1.0 候选）**；v1.0.0 是否发布由维护者决定。

## 文档

- [docs/research.md](docs/research.md) — 调研与取舍
- [docs/design.md](docs/design.md) — 设计
- [docs/architecture.md](docs/architecture.md) — 架构
- [docs/api.md](docs/api.md) — API 快照
- [docs/errors.md](docs/errors.md) — 错误码手册
- [docs/adapters.md](docs/adapters.md) — 家族接入规范
- [docs/final-review.md](docs/final-review.md) — 1.0 候选终审
- [docs/roadmap.md](docs/roadmap.md) — 路线图

## License

MIT © [lcylpzls](https://github.com/lcylpzls)
