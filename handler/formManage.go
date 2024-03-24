package handler

import (
	"rop2-api/model"
	"rop2-api/utils"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

func formInit(routerGroup *gin.RouterGroup) {
	formGroup := routerGroup.Group("/form", AuthWithRefresh(true))

	formGroup.GET("/list", getFormList)
	formGroup.GET("/detail", getFormDetail)
	formGroup.POST("/edit", editForm)
}

// 获取表单列表，只有简略信息：id,name,start/endAt,create/updateAt
func getFormList(ctx *gin.Context) {
	iden := ctx.MustGet("identity").(*AdminIdentity)
	orgId := iden.At

	//不考虑分批查询，一次查询并返回
	ctx.PureJSON(200, model.GetForms(orgId))
}

// 获取单个表单详情，返回全部信息
func getFormDetail(ctx *gin.Context) {
	iden := ctx.MustGet("identity").(*AdminIdentity)
	orgId := iden.At

	type Arg struct {
		Id uint32 `form:"id"`
	}
	arg := &Arg{}
	if ctx.ShouldBindQuery(arg) != nil {
		ctx.AbortWithStatusJSON(utils.Message("参数绑定失败", 400, 0))
		return
	}

	formId := arg.Id
	form := model.GetFormDetail(orgId, formId)
	if form != nil {
		ctx.PureJSON(200, form)
	} else {
		ctx.AbortWithStatusJSON(utils.Message("表单不存在", 422, 1))
		return
	}
}

// 编辑表单，query传入id，body为json包含要编辑的字段和新值
func editForm(ctx *gin.Context) {
	iden := ctx.MustGet("identity").(*AdminIdentity)

	if iden.Level < model.Maintainer {
		ctx.AbortWithStatusJSON(utils.Message("权限不足", 403, 1))
		return
	}

	type Arg struct {
		Id uint32 `form:"id"`
	}
	arg := &Arg{}
	if ctx.ShouldBindQuery(arg) != nil {
		ctx.AbortWithStatusJSON(utils.Message("参数绑定失败", 400, 0))
		return
	}

	formId := arg.Id
	form := model.GetFormDetail(iden.At, formId)
	if form == nil {
		ctx.AbortWithStatusJSON(utils.Message("表单不存在", 422, 1))
		return
	}

	iter := jsoniter.Parse(jsoniter.ConfigFastest, ctx.Request.Body, 256)
	matched := false
	var timePointer **time.Time
	for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
		switch field {
		case "name":
			newValue := iter.ReadString()
			if !utils.LenBetween(newValue, 1, 25) {
				ctx.AbortWithStatusJSON(utils.Message("表单名称长度无效", 400, 11))
				return
			}
			form.Name = newValue
		case "desc":
			newValue := iter.ReadString()
			//不考虑长度限制(longtext)
			form.Desc = newValue
		case "startAt":
			timePointer = &form.StartAt
			matched = true
			fallthrough
		case "endAt":
			if matched {
				matched = false
			} else {
				timePointer = &form.EndAt
			}
			newValue := iter.ReadString()
			newTime, err := time.Parse(time.RFC3339, newValue)
			if err != nil {
				ctx.AbortWithStatusJSON(utils.Message("时间转换失败", 400, 12))
				return
			}
			*timePointer = &newTime
		case "children":
			form.Children = iter.ReadString()
		}
	}
	model.SaveForm(form)
}
