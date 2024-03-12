package utils

import "github.com/gin-gonic/gin"

func BadRequest(ctx *gin.Context, message string) {
	ctx.PureJSON(400, gin.H{
		"message": message,
	})
}

type J map[string]any

func JSON(ctx *gin.Context, data J) {
	ctx.PureJSON(200, data)
}

func NoData(ctx *gin.Context) {
	ctx.PureJSON(204, J{})
}
