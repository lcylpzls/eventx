package eventx

import (
	"testing"
)

// FuzzTopic 验证任意主题字符串下校验不 panic。
func FuzzTopic(f *testing.F) {
	f.Add("orders.created")
	f.Add("")
	f.Add("a.*")
	f.Add(string(make([]byte, 300)))
	f.Fuzz(func(t *testing.T, topic string) {
		_ = validateTopic(topic)
	})
}
