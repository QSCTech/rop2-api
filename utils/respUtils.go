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

//参数绑定失败的公用错误消息&错误码。code为400001
func MessageBindFail() (int, *CodeMessageObj) {
	return Message("参数绑定失败", 400, 1)
}

//部门、表单重名的公用错误消息&错误码。code为409001
func MessageDuplicate() (int, *CodeMessageObj) {
	return Message("存在同名对象", 409, 1)
}

//拒绝访问的公用错误消息&错误码。code为403001，只适用于有权限且权限不足（试图跨组织操作为其它错误）
func MessageForbidden() (int, *CodeMessageObj) {
	return Message("权限不足", 403, 1)
}

//对象不存在的公用错误消息&错误码。code为404001
func MessageNotFound() (int, *CodeMessageObj) {
	return Message("对象不存在", 404, 1)
}

func MessageInvalidLength(isTooShort bool) (int, *CodeMessageObj) {
	var subMessage string
	var subCode int
	if isTooShort {
		subMessage = "过短"
		subCode = 12
	} else {
		subMessage = "过长"
		subCode = 11
	}
	return Message("文本长度"+subMessage, 422, subCode)
}
