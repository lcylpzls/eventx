package eventx

import (
	"strings"
	"unicode"

	"github.com/lcylpzls/errx"
)

// maxTopicLength 是主题最大长度。
const maxTopicLength = 256

// validatePublishTopic 校验发布主题：非空、长度受限、无控制字符、
// 段非空、无通配符。
func validatePublishTopic(topic string) error {
	if err := validateBaseTopic(topic); err != nil {
		return err
	}
	if strings.Contains(topic, "*") {
		return errx.NewCode(CodeInvalidTopic, "发布主题不能包含通配符")
	}
	return nil
}

// validateSubscribeTopic 校验订阅主题：允许 `*`（单段）与 `**`（多段），
// 通配符必须作为独立段出现。
func validateSubscribeTopic(topic string) error {
	if err := validateBaseTopic(topic); err != nil {
		return err
	}
	for _, seg := range strings.Split(topic, ".") {
		if strings.Contains(seg, "*") && seg != "*" && seg != "**" {
			return errx.NewCodef(CodeInvalidTopic, "通配符只能作为独立段：%q", seg)
		}
	}
	return nil
}

// validateBaseTopic 校验主题公共规则：非空、长度受限、无控制字符、段非空。
func validateBaseTopic(topic string) error {
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
	for _, seg := range strings.Split(topic, ".") {
		if seg == "" {
			return errx.NewCode(CodeInvalidTopic, "主题段不能为空")
		}
	}
	return nil
}
