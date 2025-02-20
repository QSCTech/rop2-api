package handler

import (
	"rop2-api/model"
	"rop2-api/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func resultInit(routerGroup *gin.RouterGroup) {
	formGroup := routerGroup.Group("/result", RequireAdminWithRefresh(true))

	formGroup.GET("/intents", listIntents)
	formGroup.GET("", getResultDetail)
	formGroup.POST("/set", RequireLevel(model.Maintainer), setIntents)
}

func listIntents(ctx *gin.Context) {
	id := ctx.MustGet("identity").(AdminIdentity)
	type Arg struct {
		//注：binding:"required"会拒绝0值
		Offset int            `form:"offset"`
		Limit  int            `form:"limit" binding:"required"`
		Filter string         `form:"filter"`
		Depart string         `form:"depart" binding:"required"` //格式: 1,2,3
		Step   model.StepType `form:"step"`
		FormId uint32         `form:"formId" binding:"required"`
	}
	arg := &Arg{}
	if ctx.ShouldBindQuery(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	formId := arg.FormId
	if !model.CheckFormOwner(id.At, formId) {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}

	departs := strings.Split(arg.Depart, ",")
	departIds := make([]uint32, len(departs))
	for i, v := range departs {
		parsedUint, _ := strconv.ParseUint(v, 10, 32)
		departIds[i] = uint32(parsedUint)
	}

	ctx.PureJSON(200, model.ListIntents(formId, departIds, arg.Step, arg.Offset, arg.Limit, arg.Filter))
}

// 获取候选人的详细答卷，一次性获取多个并返回数组
func getResultDetail(ctx *gin.Context) {
	id := ctx.MustGet("identity").(AdminIdentity)
	type Arg struct {
		//注：binding:"required"会拒绝0值
		FormId uint32 `form:"formId" binding:"required"`
		Target string `form:"target" binding:"required"` //格式：1,2,3
	}
	arg := &Arg{}
	if ctx.ShouldBindQuery(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	formId := arg.FormId
	if !model.CheckFormOwner(id.At, formId) {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}

	ctx.PureJSON(200, model.GetResult(formId, strings.Split(arg.Target, ",")))
}

func setIntents(ctx *gin.Context) {
	id := ctx.MustGet("identity").(AdminIdentity)
	type Arg struct {
		IntentIds []uint32 `json:"intentIds" binding:"required"`
		//允许设成0 即已填表
		Step   model.StepType `json:"step"`
		FormId uint32         `json:"formId" binding:"required"`
	}
	arg := &Arg{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	formId := arg.FormId
	if !model.CheckFormOwner(id.At, formId) {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}

	model.SetIntents(arg.FormId, arg.IntentIds, arg.Step)
	ctx.PureJSON(utils.Success())
}
