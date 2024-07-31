package handler

import (
	"encoding/json"
	"rop2-api/model"
	"rop2-api/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func applicantInit(routerGroup *gin.RouterGroup) {
	applicantGroup := routerGroup.Group("/applicant", RequireLoginWithRefresh(true))

	applicantGroup.GET("/org", applicantGetOrgDeparts)
	applicantGroup.GET("/form", applicantGetFormDetail)
	applicantGroup.POST("/form", saveForm)
	applicantGroup.GET("/profile", applicantGetProfile)

	applicantGroup.GET("/status", applicantGetStatus)
	applicantGroup.GET("/interview/list", applicantGetInterviewList)
	applicantGroup.POST("/interview/schedule", applicantScheduleInterview)
}

func isFormOpen(form *model.Form) string {
	now := time.Now()
	if form.StartAt != nil && form.StartAt.After(now) {
		return "表单未开放"
	} else if form.EndAt != nil && form.EndAt.Before(now) {
		return "表单已结束"
	}
	return ""
}

// 候选人获取组织部门列表（选择志愿时使用）
func applicantGetOrgDeparts(ctx *gin.Context) {
	type Arg struct {
		Id uint32 `form:"id"`
	}
	arg := &Arg{}
	if ctx.ShouldBindQuery(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	result := model.GetOrgDeparts(arg.Id)
	ctx.PureJSON(200, result)
}

func applicantGetFormDetail(ctx *gin.Context) {
	type Arg struct {
		Id uint32 `form:"id"`
	}
	arg := &Arg{}
	if ctx.ShouldBindQuery(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	formId := arg.Id
	form := model.ApplicantGetFormDetail(formId)
	if form != nil {
		openError := isFormOpen(form)
		if openError != "" {
			form.Desc = ""
			newChildrenBytes, _ := json.Marshal(map[string]string{"message": openError})
			form.Children = utils.RawString(newChildrenBytes)
			//children改成错误信息，无法获取题目，但是表单标题等还是可以获取的
		}
		ctx.PureJSON(200, form)
	} else {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}
}

func saveForm(ctx *gin.Context) {
	id := ctx.MustGet("identity").(userIdentity).getId()
	type Arg struct {
		FormId        uint32   `json:"formId"`
		Phone         string   `json:"phone"`
		IntentDeparts []uint32 `json:"intentDeparts"`
		Content       string   `json:"content"`
	}
	arg := &Arg{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	formId := arg.FormId
	form := model.ApplicantGetFormDetail(formId)
	if form == nil {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}
	openError := isFormOpen(form)
	if openError != "" {
		ctx.AbortWithStatusJSON(utils.Message(openError, 400, 10))
		return
	}
	orgId := form.Owner
	deps := model.GetOrgDeparts(orgId)
	defaultDepartId := model.GetOrg(orgId).DefaultDepart
	intentDeparts := make([]uint32, 0, len(arg.IntentDeparts))
	for _, v := range arg.IntentDeparts {
		for _, dep := range deps {
			if dep.Id == v {
				intentDeparts = append(intentDeparts, v)
				break
			}
		}
	}
	if len(intentDeparts) == 0 {
		intentDeparts = []uint32{defaultDepartId}
	}
	if err := model.SaveResult(formId, id, arg.Content); err != nil {
		ctx.AbortWithStatusJSON(utils.Message("问卷提交失败(答案保存失败)", 500, 11))
		return
	}
	if err := model.SaveIntents(formId, id, intentDeparts); err != nil {
		ctx.AbortWithStatusJSON(utils.Message("问卷提交失败(志愿生成失败)", 500, 11))
		return
	}
	if err := model.SaveProfile(id, arg.Phone); err != nil {
		ctx.AbortWithStatusJSON(utils.Message("个人信息保存失败", 400, 12))
		return
	}
	ctx.PureJSON(utils.Success())
}

func applicantGetProfile(ctx *gin.Context) {
	id := ctx.MustGet("identity").(userIdentity).getId()
	person := model.FindPerson(id)
	if person == nil {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}
	type Profile struct {
		ZjuId string `json:"zjuId"`
		Phone string `json:"phone"`
	}
	ctx.PureJSON(200, Profile{ZjuId: person.ZjuId, Phone: *person.Phone})
}

func applicantGetStatus(ctx *gin.Context) {
	type Arg struct {
		FormId uint32 `form:"formId" binding:"required"`
	}
	arg := &Arg{}
	if ctx.ShouldBindQuery(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}
	zjuId := ctx.MustGet("identity").(userIdentity).getId()

	ctx.PureJSON(200, model.QueryIntentsOfPerson(arg.FormId, zjuId))
}

// 候选人提供formId，departId，获取可用的面试列表（返回数组，length可能为0）
// 如果已经安排了面试，返回已安排的面试（直接返回单个对象json）
func applicantGetInterviewList(ctx *gin.Context) {
	type Arg struct {
		FormId   uint32 `form:"formId" binding:"required"`
		DepartId uint32 `form:"departId" binding:"required"`
	}
	arg := &Arg{}
	if ctx.ShouldBindQuery(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}
	zjuId := ctx.MustGet("identity").(userIdentity).getId()
	formId := arg.FormId
	departId := arg.DepartId

	intents := model.QueryIntentsOfPerson(formId, zjuId)
	for _, v := range intents {
		if v.Depart == departId {
			step := v.Step
			scheduledIv := model.GetInterviewByIntent(formId, zjuId, departId, step)
			if scheduledIv != nil {
				ctx.PureJSON(200, scheduledIv)
				return
			} else {
				ctx.PureJSON(200, model.GetInterviews(formId, []uint32{departId}, step))
				return
			}
		}
	}
	ctx.PureJSON(utils.MessageNotFound())
}

func applicantScheduleInterview(ctx *gin.Context) {
	type Arg struct {
		FormId      uint32 `form:"formId" binding:"required"`
		InterviewId uint32 `form:"interviewId" binding:"required"`
	}
	arg := &Arg{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}
	zjuId := ctx.MustGet("identity").(userIdentity).getId()

	formId := arg.FormId
	intents := model.QueryIntentsOfPerson(formId, zjuId)
	interviewInst := model.GetInterviewById(arg.InterviewId)
	if interviewInst == nil { //面试不存在
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}

	if model.GetInterviewByIntent(formId, zjuId, interviewInst.Depart, interviewInst.Step) != nil { //已经安排了面试
		ctx.AbortWithStatusJSON(utils.Message("已安排面试", 400, 31))
		return
	}
	for _, v := range intents { //对所有志愿遍历，看有没有符合阶段和部门的
		if v.Depart == interviewInst.Depart && v.Step == interviewInst.Step {
			model.AddScheduledId(arg.InterviewId, zjuId)
			ctx.PureJSON(utils.Success())
			return
		}
	}
	ctx.PureJSON(utils.MessageNotFound())
}
