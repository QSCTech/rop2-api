package handler

import (
	"rop2-api/model"

	"github.com/gin-gonic/gin"
)

func orgInit(routerGroup *gin.RouterGroup) {
	routerGroup.GET("org", AuthWithRefresh(true), getOrgInfo)
}

// 返回登录所在组织的信息
func getOrgInfo(ctx *gin.Context) {
	id := ctx.MustGet("identity").(*UserIdentity)
	org := model.GetOrg(id.At)
	if org == nil {
		ctx.PureJSON(204, gin.H{})
	} else {
		ctx.PureJSON(200, org)
	}
}
