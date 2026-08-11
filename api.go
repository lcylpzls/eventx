package eventx

import (
	"github.com/lcylpzls/eventx/internal/core"
	"github.com/lcylpzls/logx"
)

type (
	TypedHandler[T any] = core.TypedHandler[T]
	Event               = core.Event
	Handler             = core.Handler
	Filter              = core.Filter
	Subscription        = core.Subscription
	SubscribeOption     = core.SubscribeOption
	Bus                 = core.Bus
	Metrics             = core.Metrics
	MetricProvider      = core.MetricProvider
	Option              = core.Option
)

const (
	CodeBusClosed            = core.CodeBusClosed
	CodeInvalidTopic         = core.CodeInvalidTopic
	CodeInvalidHandler       = core.CodeInvalidHandler
	CodeSubscriptionNotFound = core.CodeSubscriptionNotFound
	CodeHandlerPanic         = core.CodeHandlerPanic
	CodeQueueFull            = core.CodeQueueFull
	CodeInvalidOption        = core.CodeInvalidOption
	CodeCancelled            = core.CodeCancelled
	CodeTypeMismatch         = core.CodeTypeMismatch
)

func New(opts ...Option) (*Bus, error)   { return core.New(opts...) }
func WithPriority(p int) SubscribeOption { return core.WithPriority(p) }
func WithWorkers(n int) Option           { return core.WithWorkers(n) }
func WithQueueSize(n int) Option         { return core.WithQueueSize(n) }
func WithErrorHandler(fn func(error)) Option {
	return core.WithErrorHandler(fn)
}
func WithLogger(logger logx.Logger) Option { return core.WithLogger(logger) }

func SubscribeTyped[T any](b *Bus, topic string, handler TypedHandler[T]) (Subscription, error) {
	return core.SubscribeTyped(b, topic, handler)
}
