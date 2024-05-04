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

// 计算SHA256哈希
func Sha256(data []byte, repeatTimes uint) []byte {
	if repeatTimes <= 0 {
		return data
	}
	hasher := sha256.New()
	hasher.Write(data)
	return Sha256(hasher.Sum(nil), repeatTimes-1)
}
