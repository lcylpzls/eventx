package eventx

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/lcylpzls/logx"
	"github.com/lcylpzls/testx"
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
	testx.RequireEqual(t, int(m.Publishes), 5)
	testx.RequireEqual(t, int(m.Deliveries), 10)
	testx.RequireEqual(t, int(m.Failures), 5)
	testx.RequireEqual(t, int(m.Subscriptions), 2)
}

func TestMetricsAsync(t *testing.T) {
	bus := newBus(t)
	_, _ = bus.Subscribe("t", func(ctx context.Context, e Event) error { return nil })
	for i := 0; i < 3; i++ {
		testx.RequireNoError(t, bus.PublishAsync(context.Background(), "t", nil))
	}
	_ = bus.Close()
	m := bus.Metrics()
	testx.RequireEqual(t, int(m.Publishes), 3)
	testx.RequireEqual(t, int(m.Deliveries), 3)
	testx.RequireEqual(t, int(m.Failures), 0)
	testx.RequireEqual(t, int(m.Subscriptions), 0)
}

func TestWithLogger(t *testing.T) {
	logger := &fakeLogger{}
	bus := newBus(t, WithLogger(logger))
	_, _ = bus.Subscribe("orders.created", func(ctx context.Context, e Event) error { return nil })
	testx.RequireNoError(t, bus.Publish(context.Background(), "orders.created", "10086"))
	logger.mu.Lock()
	defer logger.mu.Unlock()
	testx.RequireLen(t, logger.debugs, 1)
	testx.RequireEqual(t, logger.debugs[0], "事件分发完成")
}

func TestWithLoggerNil(t *testing.T) {
	bus := newBus(t, WithLogger(nil))
	testx.RequireNoError(t, bus.Publish(context.Background(), "t", nil))
}

func TestAuditFields(t *testing.T) {
	g := auditFields("t", 3, time.Now(), nil)
	keys := map[string]bool{}
	for i := 0; i < g.Len(); i++ {
		keys[g.At(i).Key] = true
	}
	for _, want := range []string{"eventx.topic", "eventx.subscribers", "eventx.duration_ms"} {
		testx.RequireTrue(t, keys[want])
	}
	testx.RequireFalse(t, keys["err.code"])
}

func TestAuditFieldsWithError(t *testing.T) {
	g := auditFields("t", 1, time.Now(), errors.New("失败"))
	var hasErrCode bool
	for i := 0; i < g.Len(); i++ {
		if g.At(i).Key == "err.code" {
			hasErrCode = true
		}
	}
	testx.RequireTrue(t, hasErrCode)
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
