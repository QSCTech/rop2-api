package utils

import (
	"strconv"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

func ToStr(value uint32) string {
	return strconv.FormatUint(uint64(value), 10)
}

func Stringify(literal interface{}) string {
	str, _ := jsoniter.MarshalToString(literal)
	return str
}

// 获得相对时间戳
func ToRelTimestamp(time time.Time) uint32 {
	return uint32(time.Unix() - TimeOffset)
}

// 将相对时间戳转换回golang的time.Time
func ToTime(relTimestamp uint32) time.Time {
	return time.Unix(int64(relTimestamp)+TimeOffset, 0)
}

// 获取字符串的只读[]byte，修改slice会出错
func ToBytes(from string) []byte {
	return unsafe.Slice(unsafe.StringData(from), len(from))
}
