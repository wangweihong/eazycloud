package skipper

import (
	"github.com/wangweihong/eazycloud/pkg/sets"
)

type SkipperFunc func(string) bool

// 跳过指定前缀的方法.
func AllowPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(method string) bool {
		return sets.NewString(prefixes...).IsPrefixOf(method)
	}
}

// 跳过非指定前缀的方法.
func AllowPathPrefixNoSkipper(prefixes ...string) SkipperFunc {
	return func(method string) bool {
		return !sets.NewString(prefixes...).IsPrefixOf(method)
	}
}

// SkipInterceptor 跳过指定的中间件.
func SkipInterceptor(method string, skippers ...SkipperFunc) bool {
	for _, skipper := range skippers {
		if skipper(method) {
			return true
		}
	}
	return false
}
