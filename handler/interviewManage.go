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
	interviewGroup.POST("/freeze", RequireLevel(model.Maintainer), freezeInterview)
	interviewGroup.GET("/schedule", getInterviewScheduledIds)
	interviewGroup.POST("/schedule/delete", RequireLevel(model.Maintainer), deleteInterviewSchedule)
	interviewGroup.POST("/schedule/add", RequireLevel(model.Maintainer), addInterviewSchedule)
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

func checkInterviewOwner(id AdminIdentity, interviewId uint32) (bool, *model.Interview) {
	interviewInst := model.GetInterview(interviewId)
	if interviewInst == nil {
		return false, nil
	} else {
		return model.CheckFormOwner(id.At, interviewInst.Form), interviewInst
	}
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
	if qualified, _ := checkInterviewOwner(id, interviewId); !qualified {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}

	model.DeleteInterview(interviewId)
	ctx.PureJSON(utils.Success())
}

func freezeInterview(ctx *gin.Context) {
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
	if qualified, _ := checkInterviewOwner(id, interviewId); !qualified {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}

	model.FreezeInterview(interviewId)
	ctx.PureJSON(utils.Success())
}

// 获取报名某场面试的所有学生
func getInterviewScheduledIds(ctx *gin.Context) {
	id := ctx.MustGet("identity").(AdminIdentity)
	type Arg struct {
		Id uint32 `form:"id" binding:"required"`
	}
	arg := &Arg{}
	if ctx.ShouldBindQuery(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	interviewId := arg.Id
	if qualified, _ := checkInterviewOwner(id, interviewId); !qualified {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}

	ctx.PureJSON(200, model.GetScheduledIds(interviewId))
}

func deleteInterviewSchedule(ctx *gin.Context) {
	id := ctx.MustGet("identity").(AdminIdentity)
	type Arg struct {
		Id    uint32         `json:"id" binding:"required"`
		ZjuId model.PersonId `json:"zjuId" binding:"required"`
	}
	arg := &Arg{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	interviewId := arg.Id
	if qualified, _ := checkInterviewOwner(id, interviewId); !qualified {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}

	model.RemoveScheduledId(interviewId, arg.ZjuId)
	ctx.PureJSON(utils.Success())
}

// 添加面试安排，不检查唯一性。管理员可以给一个学生添加多个面试安排
func addInterviewSchedule(ctx *gin.Context) {
	id := ctx.MustGet("identity").(AdminIdentity)
	type Arg struct {
		Id    uint32         `json:"id" binding:"required"`
		ZjuId model.PersonId `json:"zjuId" binding:"required"`
	}
	arg := &Arg{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	interviewId := arg.Id
	if qualified, _ := checkInterviewOwner(id, interviewId); !qualified {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}

	model.AddScheduledId(interviewId, arg.ZjuId)
	ctx.PureJSON(utils.Success())
}
