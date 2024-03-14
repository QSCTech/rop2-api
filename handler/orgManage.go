package handler

import (
	"rop2-api/model"
	"rop2-api/utils"
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
				utils.JSON(ctx, utils.J{})
			} else {
				//TODO 鉴权 需要在该组织为管理
				utils.JSON(ctx, org)
			}
		} else {
			utils.BadRequest(ctx, "缺少id")
		}
	} else {
		utils.BadRequest(ctx, "缺少id")
	}
}
