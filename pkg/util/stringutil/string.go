package stringutil

import (
	"strings"
)

func BothEmptyOrNone(str1, str2 string) bool {
	return (str1 == "" && str2 == "") || (str1 != "" && str2 != "")
}

func HasAnyPrefix(str string, prefixes ...string) bool {
	if str != "" {
		for _, p := range prefixes {
			if p != "" {
				if strings.HasPrefix(str, p) {
					return true
				}
			}
		}
	}
	return false
}

func PointerToString(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}
