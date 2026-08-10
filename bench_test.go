package eventx

import (
	"context"
	"testing"

	"github.com/lcylpzls/errx"
)

// BenchmarkPublishSync 基准：10 个订阅者的同步发布。
func BenchmarkPublishSync(b *testing.B) {
	bus, err := New()
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		_, _ = bus.Subscribe("orders.created", func(ctx context.Context, e Event) error { return nil })
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := bus.Publish(context.Background(), "orders.created", i); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPublishAsync 基准：4 worker、1024 队列的异步发布。
func BenchmarkPublishAsync(b *testing.B) {
	bus, err := New(WithWorkers(4), WithQueueSize(1<<20))
	if err != nil {
		b.Fatal(err)
	}
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error { return nil })
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := bus.PublishAsync(context.Background(), "t", i); err != nil {
			if !errx.Is(err, CodeQueueFull) {
				b.Fatal(err)
			}
		}
	}
	_ = bus.Close()
}

// BenchmarkWildcardMatch 基准：通配符模式匹配。
func BenchmarkWildcardMatch(b *testing.B) {
	p := compilePattern("orders.**")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if !p.matches("orders.a.b.c") {
			b.Fatal("匹配失败")
		}
	}
}
