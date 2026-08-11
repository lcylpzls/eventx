# eventx

进程内事件总线基座：主题发布-订阅（同步/异步）、通配符匹配、
过滤器与优先级，与 errx / logx / tracex 家族打通。

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.26.5-00ADD8?style=flat&logo=go)](https://go.dev/)
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

## 家族联动

各家族库通过事件钩子（如 `WithEventHook`）把事件转发到 eventx 总线；
钩子实现由使用方在项目层内联完成（事件桥通常只有几行），
完整示例见 `examples/ecosystem`。

> 当前状态：**v1.2.1**。

## 文档

- [docs/architecture.md](docs/architecture.md) — 架构
- [docs/errors.md](docs/errors.md) — 错误码手册

## License

MIT © [lcylpzls](https://github.com/lcylpzls)
