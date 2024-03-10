package utils

import "github.com/gin-gonic/gin"

func BadRequest(ctx *gin.Context, message string) {
	ctx.PureJSON(400, gin.H{
		"message": message,
	})
}
