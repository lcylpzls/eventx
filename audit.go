package eventx

import (
	"time"
	"github.com/lcylpzls/logx"
)

// auditFields 生成事件分发审计字段：主题、订阅者数、耗时与错误码，
// 不包含事件载荷（可能含敏感数据）。
func auditFields(topic string, subscribers int, start time.Time, err error) logx.FieldGroup {
	groups := []logx.FieldGroup{
		logx.Fields(
			logx.String("eventx.topic", topic),
			logx.Int("eventx.subscribers", subscribers),
			logx.Int64("eventx.duration_ms", time.Since(start).Milliseconds()),
		),
	}
	if err != nil {
		groups = append(groups, logx.FieldsFromError(err))
	}
	var fs []logx.Field
	for _, g := range groups {
		for i := 0; i < g.Len(); i++ {
			fs = append(fs, g.At(i))
		}
	}
	return logx.Fields(fs...)
}
