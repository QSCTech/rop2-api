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
		if form.StartAt != nil && form.StartAt.Compare(now) > 0 {
			form.Desc = ""
			//正常的children应为数组
			form.Children = `{"message":"表单未开放"}`
		} else if form.EndAt != nil && form.EndAt.Compare(now) < 0 {
			form.Desc = ""
			form.Children = `{"message":"表单已结束"}`
		}
		ctx.PureJSON(200, form)
	} else {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}
}
