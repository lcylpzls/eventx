package eventx

import "github.com/lcylpzls/errx"

// eventx 错误码全集：所有失败场景统一使用 errx 结构化错误。
const (
	// CodeBusClosed 总线已关闭。
	CodeBusClosed errx.Code = "EVENTX_BUS_CLOSED"
	// CodeInvalidTopic 主题非法。
	CodeInvalidTopic errx.Code = "EVENTX_INVALID_TOPIC"
	// CodeInvalidHandler 订阅处理函数非法。
	CodeInvalidHandler errx.Code = "EVENTX_INVALID_HANDLER"
	// CodeSubscriptionNotFound 订阅不存在。
	CodeSubscriptionNotFound errx.Code = "EVENTX_SUBSCRIPTION_NOT_FOUND"
	// CodeHandlerPanic 订阅处理函数发生未捕获异常。
	CodeHandlerPanic errx.Code = "EVENTX_HANDLER_PANIC"
	// CodeQueueFull 异步队列已满。
	CodeQueueFull errx.Code = "EVENTX_QUEUE_FULL"
	// CodeInvalidOption 选项参数非法。
	CodeInvalidOption errx.Code = "EVENTX_INVALID_OPTION"
	// CodeCancelled 发布或分发被上下文取消。
	CodeCancelled errx.Code = "EVENTX_CANCELLED"
	// CodeTypeMismatch 类型化订阅载荷类型不匹配。
	CodeTypeMismatch errx.Code = "EVENTX_TYPE_MISMATCH"
)

func init() {
	errx.RegisterCode(CodeBusClosed, "总线已关闭")
	errx.RegisterCodeKind(CodeBusClosed, errx.KindUnavailable)
	errx.RegisterCode(CodeInvalidTopic, "主题非法")
	errx.RegisterCodeKind(CodeInvalidTopic, errx.KindInvalid)
	errx.RegisterCode(CodeInvalidHandler, "订阅处理函数非法")
	errx.RegisterCodeKind(CodeInvalidHandler, errx.KindInvalid)
	errx.RegisterCode(CodeSubscriptionNotFound, "订阅不存在")
	errx.RegisterCodeKind(CodeSubscriptionNotFound, errx.KindInvalid)
	errx.RegisterCode(CodeHandlerPanic, "订阅处理函数发生未捕获异常")
	errx.RegisterCodeKind(CodeHandlerPanic, errx.KindInternal)
	errx.RegisterCode(CodeQueueFull, "异步队列已满")
	errx.RegisterCodeKind(CodeQueueFull, errx.KindQuotaExceeded)
	errx.RegisterCode(CodeInvalidOption, "选项参数非法")
	errx.RegisterCodeKind(CodeInvalidOption, errx.KindInvalid)
	errx.RegisterCode(CodeCancelled, "发布或分发被上下文取消")
	errx.RegisterCodeKind(CodeCancelled, errx.KindCancelled)
	errx.RegisterCode(CodeTypeMismatch, "类型化订阅载荷类型不匹配")
	errx.RegisterCodeKind(CodeTypeMismatch, errx.KindInvalid)
}
