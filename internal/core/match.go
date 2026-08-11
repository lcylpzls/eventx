package core

import "strings"

// pattern 是编译后的订阅主题模式。
type pattern struct {
	segments []string
}

// compilePattern 编译订阅主题为段模式。
func compilePattern(topic string) *pattern {
	return &pattern{segments: strings.Split(topic, ".")}
}

// matches 判断发布主题是否命中模式（支持 `*` 与 `**`）。
func (p *pattern) matches(topic string) bool {
	actual := strings.Split(topic, ".")
	return matchSegments(p.segments, actual)
}

// matchSegments 递归匹配模式段与实际段。
func matchSegments(pattern, actual []string) bool {
	var match func(pi, ai int) bool
	match = func(pi, ai int) bool {
		if pi == len(pattern) {
			return ai == len(actual)
		}
		seg := pattern[pi]
		if seg == "**" {
			for k := ai; k <= len(actual); k++ {
				if match(pi+1, k) {
					return true
				}
			}
			return false
		}
		if ai >= len(actual) {
			return false
		}
		if seg != "*" && seg != actual[ai] {
			return false
		}
		return match(pi+1, ai+1)
	}
	return match(0, 0)
}
