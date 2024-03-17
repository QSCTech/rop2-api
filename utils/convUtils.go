package utils

import (
	"encoding/base64"
	"strconv"
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

// 获取字符串的只读[]byte，修改slice会出错
func RawBytes(from string) []byte {
	if from == "" {
		return []byte{}
	}
	return unsafe.Slice(unsafe.StringData(from), len(from))
}

func RawString(from []byte) string {
	return unsafe.String(&from[0], len(from))
}

func Base64Encode(from []byte) string {
	return base64.RawStdEncoding.EncodeToString(from)
}

func Base64Decode(from string) []byte {
	result, err := base64.RawStdEncoding.DecodeString(from)
	if err != nil {
		panic(err)
	}
	return result
}

// 类似js的map工具方法。
func MapArray[TElement any, TResult any](array []TElement, hanlder func(TElement, int) TResult) []TResult {
	result := make([]TResult, len(array))
	for i, v := range array {
		result[i] = hanlder(v, i)
	}
	return result
}
