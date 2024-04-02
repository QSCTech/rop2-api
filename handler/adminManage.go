package handler

import (
	"rop2-api/model"
	"rop2-api/utils"

	"github.com/gin-gonic/gin"
)

func adminInit(routerGroup *gin.RouterGroup) {
	orgGroup := routerGroup.Group("/admin", AuthWithRefresh(true))
	orgGroup.GET("", getAdminList)
	orgGroup.POST("/editLevel", RequireLevel(model.Maintainer), editLevel)
}

func getAdminList(ctx *gin.Context) {
	id := ctx.MustGet("identity").(*AdminIdentity)
	//只要能以组织的身份登录就可查询
	//考虑加限制？

	type Arg struct {
		Offset int `form:"offset"`
		Limit  int `form:"limit"`
	}
	arg := &Arg{}
	if ctx.ShouldBindQuery(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	ctx.PureJSON(200, model.GetAdminsInOrg(id.At, arg.Offset, arg.Limit))
}

func editLevel(ctx *gin.Context) {
	//TODO 设为0即为删除
	// id := ctx.MustGet("identity").(*AdminIdentity)

	// type Arg struct {
	// 	ZjuId    string          `json:"zjuId"`
	// 	NewLevel model.PermLevel `json:"newLevel"`
	// }
	// arg := &Arg{}
	// if ctx.ShouldBindQuery(arg) != nil {
	// 	ctx.AbortWithStatusJSON(utils.MessageBindFail())
	// 	return
	// }
}
