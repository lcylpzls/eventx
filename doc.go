// Package eventx 是进程内事件总线基座：主题发布-订阅（同步/异步）、
// 通配符匹配、过滤器与优先级，与 errx / logx / tracex 家族打通。
//
// 典型用法：
//
//	bus := eventx.New()
//	sub, _ := bus.Subscribe("orders.created", func(ctx context.Context, e eventx.Event) error {
//	    return nil
//	})
//	defer sub.Unsubscribe()
//	_ = bus.Publish(context.Background(), "orders.created", "10086")
package eventx
