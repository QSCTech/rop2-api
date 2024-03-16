package utils

type ErrCodeObj struct {
	// 为0表示请求成功，否则为HTTP状态码+子错误码(3位)，如404001
	Code int32 `json:"code"`
}

// 每一个参数都应为0~999的整数
func CodeObj(code ...int32) *ErrCodeObj {
	var result int32 = 0
	for _, v := range code {
		result = result*1000 + v
	}
	return &ErrCodeObj{Code: result}
}
