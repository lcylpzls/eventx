package eventx

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/lcylpzls/logx"
)

func TestMetricsSync(t *testing.T) {
	bus := newBus(t)
	_, _ = bus.Subscribe("t.x", func(ctx context.Context, e Event) error {
		return errors.New("失败")
	})
	_, _ = bus.Subscribe("t.*", func(ctx context.Context, e Event) error {
		return nil
	})
	for i := 0; i < 5; i++ {
		_ = bus.Publish(context.Background(), "t.x", nil)
	}
	m := bus.Metrics()
	if m.Publishes != 5 || m.Deliveries != 10 || m.Failures != 5 || m.Subscriptions != 2 {
		t.Fatalf("指标不匹配：%+v", m)
	}
}

func TestMetricsAsync(t *testing.T) {
	bus := newBus(t)
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error { return nil })
	for i := 0; i < 3; i++ {
		if err := bus.PublishAsync(context.Background(), "t", nil); err != nil {
			t.Fatalf("PublishAsync 失败：%v", err)
		}
	}
	_ = bus.Close()
	m := bus.Metrics()
	if m.Publishes != 3 || m.Deliveries != 3 || m.Failures != 0 || m.Subscriptions != 0 {
		t.Fatalf("异步指标不匹配：%+v", m)
	}
}

func TestWithLogger(t *testing.T) {
	logger := &fakeLogger{}
	bus := newBus(t, WithLogger(logger))
	_, _ = bus.Subscribe("orders.created", func(ctx context.Context, e Event) error { return nil })
	if err := bus.Publish(context.Background(), "orders.created", "10086"); err != nil {
		t.Fatalf("Publish 失败：%v", err)
	}
	logger.mu.Lock()
	defer logger.mu.Unlock()
	if len(logger.debugs) != 1 {
		t.Fatalf("期望 1 条 Debug 审计日志，得到 %d", len(logger.debugs))
	}
	if logger.debugs[0] != "事件分发完成" {
		t.Fatalf("审计消息不匹配：%q", logger.debugs[0])
	}
}

func TestWithLoggerNil(t *testing.T) {
	bus := newBus(t, WithLogger(nil))
	if err := bus.Publish(context.Background(), "t", nil); err != nil {
		t.Fatalf("nil logger 发布失败：%v", err)
	}
}

func TestAuditFields(t *testing.T) {
	g := auditFields("t", 3, time.Now(), nil)
	keys := map[string]bool{}
	for i := 0; i < g.Len(); i++ {
		keys[g.At(i).Key] = true
	}
	for _, want := range []string{"eventx.topic", "eventx.subscribers", "eventx.duration_ms"} {
		if !keys[want] {
			t.Fatalf("审计字段缺少 %q：%v", want, keys)
		}
	}
	if keys["err.code"] {
		t.Fatal("无错误时不应包含 err.code")
	}
}

func TestAuditFieldsWithError(t *testing.T) {
	g := auditFields("t", 1, time.Now(), errors.New("失败"))
	var hasErrCode bool
	for i := 0; i < g.Len(); i++ {
		if g.At(i).Key == "err.code" {
			hasErrCode = true
		}
	}
	if !hasErrCode {
		t.Fatal("错误审计字段应包含 err.code")
	}
}

// fakeLogger 是 logx.Logger 的最小实现。
type fakeLogger struct {
	mu     sync.Mutex
	debugs []string
}

func (l *fakeLogger) IsDebugEnabled() bool { return true }

func (l *fakeLogger) Debug(msg string, _ logx.FieldGroup) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.debugs = append(l.debugs, msg)
}

func (l *fakeLogger) Info(msg string, _ logx.FieldGroup)  {}
func (l *fakeLogger) Warn(msg string, _ logx.FieldGroup)  {}
func (l *fakeLogger) Error(msg string, _ logx.FieldGroup) {}
func (l *fakeLogger) Panic(msg string, _ logx.FieldGroup) {}
func (l *fakeLogger) Fatal(msg string, _ logx.FieldGroup) {}
func (l *fakeLogger) Debugf(format string, args ...any)   {}
func (l *fakeLogger) Infof(format string, args ...any)    {}
func (l *fakeLogger) Warnf(format string, args ...any)    {}
func (l *fakeLogger) Errorf(format string, args ...any)   {}
func (l *fakeLogger) Panicf(format string, args ...any)   {}
func (l *fakeLogger) Fatalf(format string, args ...any)   {}
func (l *fakeLogger) WithContext(ctx context.Context) logx.Logger {
	return l
}
func (l *fakeLogger) WithField(key string, val any) logx.Logger {
	return l
}
func (l *fakeLogger) Sync() error  { return nil }
func (l *fakeLogger) Close() error { return nil }
func (l *fakeLogger) SafeExit(exitFunc func()) {
	if exitFunc != nil {
		exitFunc()
	}
}
