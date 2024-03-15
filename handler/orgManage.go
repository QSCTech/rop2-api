package handler

import (
	"rop2-api/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func orgInit(routerGroup gin.RouterGroup) {
	routerGroup.GET("org", getOrgInfo)
}

func getOrgInfo(ctx *gin.Context) {
	if id, exist := ctx.GetQuery("id"); exist {
		if id64, err := strconv.ParseUint(id, 10, 32); err == nil {
			id32 := uint32(id64)
			org := model.GetOrg(id32)
			if org == nil {
				ctx.PureJSON(204, gin.H{})
			} else {
				//TODO 鉴权 需要在该组织为管理
				ctx.PureJSON(200, *org)
			}
		} else {
			ctx.AbortWithStatus(400)
		}
	} else {
		ctx.AbortWithStatus(400)
	}
}
