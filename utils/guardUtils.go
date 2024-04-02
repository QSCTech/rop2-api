package utils

import "unicode/utf8"

//检查字符串长度是否在指定范围，按code point计数。
//返回0表示在范围内，负数表示小于min，正数表示大于max
func LenBetween(target string, min, max int) int {
	len := utf8.RuneCountInString(target)
	diff := len - min
	if diff < 0 {
		return diff //如果len < min，返回len-min (负数)
	}
	diff = max - len
	if diff < 0 {
		return -diff //如果max < len，返回max - len (正数)
	}
	return 0
}
