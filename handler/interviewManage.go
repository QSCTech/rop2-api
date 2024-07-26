package handler

import (
	"rop2-api/model"
	"rop2-api/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func interviewInit(routerGroup *gin.RouterGroup) {
	interviewGroup := routerGroup.Group("/interview", RequireAdminWithRefresh(true))

	interviewGroup.GET("", getInterviews)
	interviewGroup.POST("/add", RequireLevel(model.Maintainer), addInterview)
	interviewGroup.POST("/delete", RequireLevel(model.Maintainer), deleteInterview)
}

func getInterviews(ctx *gin.Context) {
	id := ctx.MustGet("identity").(AdminIdentity)
	type Arg struct {
		FormId uint32         `form:"formId" binding:"required"`
		Step   model.StepType `form:"step"`
		Depart string         `form:"depart" binding:"required"` //格式: 1,2,3
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

	ctx.PureJSON(200, model.GetInterviews(formId, departIds, arg.Step))
}

func addInterview(ctx *gin.Context) {
	id := ctx.MustGet("identity").(AdminIdentity)
	type Arg struct {
		FormId   uint32         `json:"formId" binding:"required"`
		Depart   uint32         `json:"depart" binding:"required"`
		Step     model.StepType `json:"step" binding:"required"`
		Capacity int32          `json:"capacity"`

		Location string    `json:"location" binding:"required"`
		StartAt  time.Time `json:"startAt" binding:"required"`
		EndAt    time.Time `json:"endAt" binding:"required"`
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

	model.AddInterview(arg.FormId, arg.Depart, arg.Step, arg.Capacity, arg.Location, arg.StartAt, arg.EndAt)
	ctx.PureJSON(200, gin.H{
		"code": 0,
	})
}

func deleteInterview(ctx *gin.Context) {
	id := ctx.MustGet("identity").(AdminIdentity)
	type Arg struct {
		Id uint32 `json:"id" binding:"required"`
	}
	arg := &Arg{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	interviewId := arg.Id
	if !model.CheckFormOwner(id.At, interviewId) {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}

	interviewInst := model.GetInterview(interviewId)
	if interviewInst == nil || !model.CheckFormOwner(id.At, interviewInst.Form) {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}

	if err := model.DeleteInterview(interviewId); err != nil {
		ctx.PureJSON(utils.MessageNotFound())
		return
	}
	ctx.PureJSON(utils.Success())
}
