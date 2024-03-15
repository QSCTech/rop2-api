package utils

import (
	"crypto/hmac"
	"crypto/sha256"
)

// 计算HMAC-SHA256签名
func HmacSha256(data, key []byte) []byte {
	hasher := hmac.New(sha256.New, key)
	_, _ = hasher.Write(data)
	return hasher.Sum(nil)
}
