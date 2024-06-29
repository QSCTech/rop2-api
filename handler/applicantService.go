package handler

import (
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
		now := time.Now()
		if form.StartAt != nil && form.StartAt.After(now) {
			form.Desc = ""
			//正常的children应为数组
			form.Children = `{"message":"表单未开放"}`
		} else if form.EndAt != nil && form.EndAt.Before(now) {
			form.Desc = ""
			form.Children = `{"message":"表单已结束"}`
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
		FormId  uint32 `json:"formId"`
		Phone   string `json:"phone"`
		Content string `json:"content"`
	}
	arg := &Arg{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	formId := arg.FormId
	form := model.ApplicantGetFormDetail(formId)
	if form != nil {
		now := time.Now()
		if form.StartAt != nil && form.StartAt.After(now) {
			ctx.AbortWithStatusJSON(utils.Message("表单未开放", 400, 21))
			return
		} else if form.EndAt != nil && form.EndAt.Before(now) {
			ctx.AbortWithStatusJSON(utils.Message("表单已结束", 400, 22))
			return
		}
		if err := model.SaveResult(formId, id, arg.Content); err != nil {
			ctx.AbortWithStatusJSON(utils.Message("问卷提交失败", 500, 11))
			return
		}
		if err := model.SaveProfile(id, arg.Phone); err != nil {
			ctx.AbortWithStatusJSON(utils.Message("个人信息保存失败", 500, 12))
			return
		}
		ctx.PureJSON(utils.Success())
	} else {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}
}
