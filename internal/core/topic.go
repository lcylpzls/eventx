package core

import (
	"strings"
	"unicode"

	"github.com/lcylpzls/errx"
	"github.com/lcylpzls/validx"
)

// maxTopicLength 是主题最大长度。
const maxTopicLength = 256

// init 注册主题规则到 validx 全局规则表，
// 参数区分发布（publish）与订阅（subscribe）语义。
func init() {
	_ = validx.RegisterRule("eventx_topic", func(value any, param, path string) error {
		// 内部调用保证 value 为 string、param 为 publish/subscribe。
		topic := value.(string)
		if err := checkBaseTopic(topic); err != nil {
			return err
		}
		switch param {
		case "publish":
			if strings.Contains(topic, "*") {
				return errx.NewCode(CodeInvalidTopic, "发布主题不能包含通配符")
			}
		case "subscribe":
			for _, seg := range strings.Split(topic, ".") {
				if strings.Contains(seg, "*") && seg != "*" && seg != "**" {
					return errx.NewCodef(CodeInvalidTopic, "通配符只能作为独立段：%q", seg)
				}
			}
		}
		return nil
	})
}

// validatePublishTopic 校验发布主题：非空、长度受限、无控制字符、
// 段非空、无通配符。
func validatePublishTopic(topic string) error {
	return validx.ValidateField(topic, "eventx_topic=publish")
}

// validateSubscribeTopic 校验订阅主题：允许 `*`（单段）与 `**`（多段），
// 通配符必须作为独立段出现。
func validateSubscribeTopic(topic string) error {
	return validx.ValidateField(topic, "eventx_topic=subscribe")
}

// checkBaseTopic 校验主题公共规则：非空、长度受限、无控制字符、段非空。
func checkBaseTopic(topic string) error {
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
