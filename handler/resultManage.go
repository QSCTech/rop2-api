package handler

import (
	"rop2-api/utils"

	"github.com/gin-gonic/gin"
)

func resultInit(routerGroup *gin.RouterGroup) {
	formGroup := routerGroup.Group("/result", RequireAdminWithRefresh(true))

	formGroup.GET("", listResults)
}

func listResults(ctx *gin.Context) {
	id := ctx.MustGet("identity").(AdminIdentity)
	_ = id
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

	// ctx.PureJSON(200, model.GetAdminsInOrg(id.At, arg.Offset, arg.Limit, arg.Filter))
}
