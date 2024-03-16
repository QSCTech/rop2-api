package utils

import "github.com/gin-gonic/gin"

// 每一个参数都应为0~999的整数
func getCode(code ...int) int {
	var result = 0
	for _, v := range code {
		result = result*1000 + v
	}
	return result
}

type CodeMessageObj struct {
	// 为0表示请求成功，否则为HTTP状态码+子错误码(3位)，如404001
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Message(message string, baseStatusCode int, code ...int) (int, *CodeMessageObj) {
	return baseStatusCode, &CodeMessageObj{
		Message: message,
		Code:    getCode(append([]int{baseStatusCode}, code...)...),
	}
}

func Success() (int, gin.H) {
	return 200, gin.H{
		"code": 0,
	}
}
