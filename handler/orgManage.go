package handler

import (
	"rop2-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func orgInit(routerGroup gin.RouterGroup) {
	routerGroup.GET("", func(ctx *gin.Context) {
		if id, exist := ctx.GetQuery("id"); exist {
			if idNum, err := strconv.ParseUint(id, 10, 32); err == nil {
				//TODO
			}
		} else {
			utils.BadRequest(ctx, "缺少id")
		}
	})
}
