package utils

import (
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func ToStr(value uint32) string {
	return strconv.FormatUint(uint64(value), 10)
}

func Stringify(literal interface{}) string {
	str, _ := jsoniter.MarshalToString(literal)
	return str
}
