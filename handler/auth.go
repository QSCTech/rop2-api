package handler

import (
	"time"

	"github.com/gin-gonic/gin"
)

var loginMap map[string]*struct {
	zjuId    string
	notAfter time.Time //理论过期时间，由于定时批量清除可能略有偏差
}

const authStringLength = 64

// 中间件，要求用户必须登录才能访问API
// 用户学号(string)存至ctx.Keys["zjuId"]
func RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("rop-auth")
		authLen := len(auth)
		if authLen != authStringLength {
			ctx.AbortWithStatus(401)
			return
		}
		str, ok := loginMap[auth]
		if !ok {
			ctx.AbortWithStatus(401)
			return
		}
		ctx.Keys["zjuId"] = str.zjuId
		ctx.Next()
	}
}
