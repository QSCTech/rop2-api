package handler

import (
	"rop2-api/model"
	"rop2-api/utils"

	"github.com/gin-gonic/gin"
)

func adminInit(routerGroup *gin.RouterGroup) {
	orgGroup := routerGroup.Group("/admin", RequireAdminWithRefresh(true))
	orgGroup.GET("", listAdmins)
	orgGroup.POST("/edit", RequireLevel(model.Maintainer), editAdmin)
}

func listAdmins(ctx *gin.Context) {
	id := ctx.MustGet("identity").(AdminIdentity)
	//只要能以组织的身份登录就可查询
	//考虑加限制？

	type Arg struct {
		Offset int    `form:"offset"`
		Limit  int    `form:"limit"`
		Filter string `form:"filter"`
	}
	arg := &Arg{}
	if ctx.ShouldBindQuery(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	ctx.PureJSON(200, model.GetAdminsInOrg(id.At, arg.Offset, arg.Limit, arg.Filter))
}

func editAdmin(ctx *gin.Context) {
	iden := ctx.MustGet("identity").(AdminIdentity)
	at := iden.At

	type Arg struct {
		ZjuId    string          `json:"zjuId"`
		Nickname string          `json:"nickname"`
		Level    model.PermLevel `json:"level"`
	}
	arg := &Arg{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}
	model.SetAdmin(at, arg.ZjuId, arg.Nickname, arg.Level)
	//NOTE 理论上某个组织可以通过反复修改权限来阻止其他组织的管理员登录
	ForceLogoutAll(arg.ZjuId)
	ctx.PureJSON(utils.Success())
}
