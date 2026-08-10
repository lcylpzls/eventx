package eventx

import (
	"strings"
	"unicode"

	"github.com/lcylpzls/errx"
)

// maxTopicLength 是主题最大长度。
const maxTopicLength = 256

// validateTopic 校验发布主题：非空、长度受限、无控制字符、无通配符。
func validateTopic(topic string) error {
	if topic == "" {
		return errx.NewCode(CodeInvalidTopic, "主题不能为空")
	}
	if len(topic) > maxTopicLength {
		return errx.NewCodef(CodeInvalidTopic, "主题长度超过上限 %d", maxTopicLength)
	}
	for _, r := range topic {
		if unicode.IsControl(r) {
			return errx.NewCodef(CodeInvalidTopic, "主题不能包含控制字符：%q", r)
		}
	}
	if strings.Contains(topic, "*") {
		return errx.NewCode(CodeInvalidTopic, "发布主题不能包含通配符")
	}
	return nil
}
