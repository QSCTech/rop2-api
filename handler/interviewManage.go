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

	interviewGroup.GET("", listInterviews)
	interviewGroup.GET("/detail", getInterviewDetail)
	interviewGroup.POST("/add", RequireLevel(model.Maintainer), addInterview)
	interviewGroup.POST("/delete", RequireLevel(model.Maintainer), deleteInterview)
	interviewGroup.POST("/freeze", RequireLevel(model.Maintainer), freezeInterview)

	interviewGroup.GET("/schedule", getInterviewScheduledIds)
	//管理员删除/添加面试报名信息。不受冻结/满人限制
	interviewGroup.POST("/schedule/delete", RequireLevel(model.Maintainer), deleteInterviewSchedule)
	interviewGroup.POST("/schedule/add", RequireLevel(model.Maintainer), addInterviewSchedule)
}

func listInterviews(ctx *gin.Context) {
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

// 根据id获取面试详情，仅限管理员
func getInterviewDetail(ctx *gin.Context) {
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
	qualified, interviewInst := checkInterviewOwner(id, interviewId)
	if !qualified {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}
	ctx.PureJSON(200, interviewInst)
}

// 检查指定面试是否由管理员所在组织创建
func checkInterviewOwner(id AdminIdentity, interviewId uint32) (bool, *model.Interview) {
	interviewInst := model.GetInterviewById(interviewId)
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
		Comment  *string   `json:"comment"` //可选，也可能是空字符串（也存数据库，由前端再判断）
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

	model.AddInterview(arg.FormId, arg.Depart, arg.Step, arg.Capacity, arg.Location, arg.StartAt, arg.EndAt, arg.Comment)
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

// 管理员删除面试报名信息
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

	model.DeleteInterviewSchedule(interviewId, arg.ZjuId)
	ctx.PureJSON(utils.Success())
}

// 管理员添加面试安排。无视时间、冻结、容量；但同一志愿只能有一个面试安排
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
	zjuId := arg.ZjuId
	qualified, interviewInst := checkInterviewOwner(id, interviewId)
	if !qualified {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}

	intents := model.QueryIntentsOfPerson(interviewInst.Form, zjuId)
	find := false
	for _, intent := range intents {
		if intent.Depart == interviewInst.Depart && intent.Step == interviewInst.Step {
			find = true
			break
		}
	}
	if !find {
		ctx.AbortWithStatusJSON(utils.Message("候选人无此部门和阶段的志愿", 400, 41))
		return
	}

	// 如果此部门&阶段已有面试安排，删除之
	existantSchedule := model.GetInterviewByIntent(interviewInst.Form, zjuId, interviewInst.Depart, interviewInst.Step)
	if existantSchedule != nil {
		model.DeleteInterviewSchedule(existantSchedule.Id, zjuId)
	}

	// 无视面试时间&冻结&容量，添加面试安排
	model.AddInterviewSchedule(model.DefaultDb(), interviewId, zjuId)
	ctx.PureJSON(utils.Success())
}
