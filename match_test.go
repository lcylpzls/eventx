package eventx

import (
	"context"
	"strings"
	"sync"
	"testing"
)

func TestPatternMatch(t *testing.T) {
	tests := []struct {
		pattern string
		topic   string
		want    bool
	}{
		{"orders.created", "orders.created", true},
		{"orders.created", "orders.updated", false},
		{"orders.*", "orders.created", true},
		{"orders.*", "orders.a.b", false},
		{"orders.**", "orders.created", true},
		{"orders.**", "orders.a.b", true},
		{"orders.**", "orders", true},
		{"**", "a.b.c", true},
		{"**", "a", true},
		{"a.*.c", "a.b.c", true},
		{"a.*.c", "a.b.d", false},
		{"a.*.c", "a.b", false},
		{"a.**.b", "a.b", true},
		{"a.**.b", "a.x.b", true},
		{"a.**.b", "a.x.y", false},
	}
	for _, tt := range tests {
		got := compilePattern(tt.pattern).matches(tt.topic)
		if got != tt.want {
			t.Fatalf("pattern %q match %q = %v，期望 %v", tt.pattern, tt.topic, got, tt.want)
		}
	}
}

func TestWildcardSubscribe(t *testing.T) {
	bus := newBus(t)
	mu := sync.Mutex{}
	var topics []string
	_, _ = bus.Subscribe("orders.*", func(ctx context.Context, e Event) error {
		mu.Lock()
		topics = append(topics, e.Topic)
		mu.Unlock()
		return nil
	})
	_ = bus.Publish(context.Background(), "orders.created", nil)
	_ = bus.Publish(context.Background(), "orders.updated", nil)
	_ = bus.Publish(context.Background(), "orders.a.b", nil)
	mu.Lock()
	got := strings.Join(topics, ",")
	mu.Unlock()
	if got != "orders.created,orders.updated" {
		t.Fatalf("通配符订阅命中不匹配：%q", got)
	}
}

func TestWildcardZeroSegment(t *testing.T) {
	bus := newBus(t)
	var got []string
	_, _ = bus.Subscribe("a.**", func(ctx context.Context, e Event) error {
		got = append(got, e.Topic)
		return nil
	})
	_ = bus.Publish(context.Background(), "a", nil)
	_ = bus.Publish(context.Background(), "a.b.c", nil)
	if strings.Join(got, ",") != "a,a.b.c" {
		t.Fatalf("零段通配符匹配不匹配：%v", got)
	}
}

func TestPriorityAndOrder(t *testing.T) {
	bus := newBus(t)
	var order []string
	record := func(name string) Handler {
		return func(ctx context.Context, e Event) error {
			order = append(order, name)
			return nil
		}
	}
	_, _ = bus.Subscribe("a.b", record("精确-后"))
	_, _ = bus.SubscribeWithOptions("a.b", record("精确-前"), WithPriority(-1))
	_, _ = bus.Subscribe("a.*", record("通配-默认"))
	_, _ = bus.SubscribeWithOptions("a.*", record("通配-优先"), WithPriority(-2))
	_ = bus.Publish(context.Background(), "a.b", nil)
	got := strings.Join(order, ",")
	want := "通配-优先,精确-前,精确-后,通配-默认"
	if got != want {
		t.Fatalf("优先级/顺序不匹配：%q，期望 %q", got, want)
	}
}

func TestFilter(t *testing.T) {
	bus := newBus(t)
	var got []any
	_, _ = bus.SubscribeFiltered("t", func(ctx context.Context, e Event) bool {
		return e.Payload.(int)%2 == 0
	}, func(ctx context.Context, e Event) error {
		got = append(got, e.Payload)
		return nil
	})
	_, _ = bus.SubscribeFiltered("t", nil, func(ctx context.Context, e Event) error {
		got = append(got, "all")
		return nil
	})
	for i := 1; i <= 4; i++ {
		_ = bus.Publish(context.Background(), "t", i)
	}
	if len(got) != 6 {
		t.Fatalf("过滤器命中不匹配：%v", got)
	}
}

func TestFilterPanic(t *testing.T) {
	bus := newBus(t)
	called := false
	_, _ = bus.SubscribeFiltered("t", func(ctx context.Context, e Event) bool {
		panic("过滤器崩溃")
	}, func(ctx context.Context, e Event) error {
		called = true
		return nil
	})
	err := bus.Publish(context.Background(), "t", nil)
	assertErrCode(t, err, CodeHandlerPanic)
	if called {
		t.Fatal("过滤器 panic 后不应调用 handler")
	}
}

func TestWildcardAsync(t *testing.T) {
	bus := newBus(t)
	done := make(chan string, 1)
	_, _ = bus.Subscribe("orders.**", func(ctx context.Context, e Event) error {
		done <- e.Topic
		return nil
	})
	if err := bus.PublishAsync(context.Background(), "orders.a.b", nil); err != nil {
		t.Fatalf("PublishAsync 失败：%v", err)
	}
	_ = bus.Close()
	select {
	case got := <-done:
		if got != "orders.a.b" {
			t.Fatalf("异步通配符命中不匹配：%q", got)
		}
	default:
		t.Fatal("异步通配符订阅未命中")
	}
}

func TestWildcardUnsubscribe(t *testing.T) {
	bus := newBus(t)
	count := 0
	sub, _ := bus.Subscribe("orders.*", func(ctx context.Context, e Event) error {
		count++
		return nil
	})
	if err := sub.Unsubscribe(); err != nil {
		t.Fatalf("Unsubscribe 失败：%v", err)
	}
	if err := sub.Unsubscribe(); err != nil {
		t.Fatalf("重复 Unsubscribe 应幂等：%v", err)
	}
	_ = bus.Publish(context.Background(), "orders.created", nil)
	if count != 0 {
		t.Fatalf("取消后通配订阅不应命中：%d", count)
	}
}

func TestWildcardUnsubscribeAfterClose(t *testing.T) {
	bus := newBus(t)
	sub, _ := bus.Subscribe("orders.*", func(ctx context.Context, e Event) error { return nil })
	_ = bus.Close()
	err := sub.Unsubscribe()
	assertErrCode(t, err, CodeSubscriptionNotFound)
}
