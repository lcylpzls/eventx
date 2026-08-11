package core

import (
	"context"
	"github.com/lcylpzls/testx"
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
		testx.RequireEqual(t, got, tt.want)
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
	testx.RequireEqual(t, got, "orders.created,orders.updated")
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
	testx.RequireEqual(t, strings.Join(got, ","), "a,a.b.c")
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
	testx.RequireEqual(t, got, want)
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
	testx.RequireLen(t, got, 6)
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
	testx.RequireFalse(t, called)
}

func TestWildcardAsync(t *testing.T) {
	bus := newBus(t)
	done := make(chan string, 1)
	_, _ = bus.Subscribe("orders.**", func(ctx context.Context, e Event) error {
		done <- e.Topic
		return nil
	})
	testx.RequireNoError(t, bus.PublishAsync(context.Background(), "orders.a.b", nil))
	_ = bus.Close()
	select {
	case got := <-done:
		testx.RequireEqual(t, got, "orders.a.b")
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
	testx.RequireNoError(t, sub.Unsubscribe())
	testx.RequireNoError(t, sub.Unsubscribe())
	_ = bus.Publish(context.Background(), "orders.created", nil)
	testx.RequireEqual(t, count, 0)
}

func TestWildcardUnsubscribeAfterClose(t *testing.T) {
	bus := newBus(t)
	sub, _ := bus.Subscribe("orders.*", func(ctx context.Context, e Event) error { return nil })
	_ = bus.Close()
	err := sub.Unsubscribe()
	assertErrCode(t, err, CodeSubscriptionNotFound)
}
