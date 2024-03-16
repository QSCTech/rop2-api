package utils

import "unicode/utf8"

//检查字符串长度是否在指定范围，按code point计数
func LenBetween(target string, min, max int) bool {
	len := utf8.RuneCountInString(target)
	return len >= min && len <= max
}
