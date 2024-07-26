package handler

import (
	"errors"
	"rop2-api/model"
	"rop2-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func formInit(routerGroup *gin.RouterGroup) {
	formGroup := routerGroup.Group("/form", RequireAdminWithRefresh(true))

	formGroup.GET("/list", getFormList)
	formGroup.GET("/detail", getFormDetail)
	formGroup.POST("/edit", RequireLevel(model.Maintainer), editForm)
	formGroup.POST("/create", RequireLevel(model.Maintainer), createForm)
	formGroup.POST("/delete", RequireLevel(model.Maintainer), deleteForm)
}

// 获取表单列表，只有简略信息：id,name,start/endAt,create/updateAt
func getFormList(ctx *gin.Context) {
	iden := ctx.MustGet("identity").(AdminIdentity)
	orgId := iden.At

	//不考虑分批查询，一次查询并返回
	ctx.PureJSON(200, model.GetForms(orgId))
}

// 获取单个表单详情，返回全部信息
func getFormDetail(ctx *gin.Context) {
	iden := ctx.MustGet("identity").(AdminIdentity)
	orgId := iden.At

	type Arg struct {
		Id uint32 `form:"id"`
	}
	arg := &Arg{}
	if ctx.ShouldBindQuery(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	formId := arg.Id
	form := model.GetFormDetail(orgId, formId)
	if form != nil {
		ctx.PureJSON(200, form)
	} else {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}
}

func editForm(ctx *gin.Context) {
	iden := ctx.MustGet("identity").(AdminIdentity)
	var formUpdate model.FormUpdate
	if ctx.ShouldBindJSON(&formUpdate) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}
	formId := formUpdate.Id
	orgId := iden.At
	if !model.CheckFormOwner(orgId, formId) {
		ctx.AbortWithStatusJSON(utils.MessageForbidden())
		return
	}

	if err := model.SaveForm(formUpdate); err != nil {
		ctx.PureJSON(utils.Message(err.Error(), 400))
		return
	}
	ctx.PureJSON(utils.Success())
}

// 新建表单
func createForm(ctx *gin.Context) {
	iden := ctx.MustGet("identity").(AdminIdentity)

	type Arg struct {
		Name string `json:"name"`
	}
	arg := &Arg{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	lenDiff := utils.LenBetween(arg.Name, 1, 25)
	if lenDiff != 0 {
		ctx.AbortWithStatusJSON(utils.MessageInvalidLength(lenDiff < 0))
		return
	}

	_, err := model.CreateForm(iden.At, arg.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			ctx.AbortWithStatusJSON(utils.MessageDuplicate())
			return
		}
		ctx.AbortWithStatusJSON(utils.Message("创建失败", 400, 10))
		return
	}
	ctx.PureJSON(utils.Success())
}

func deleteForm(ctx *gin.Context) {
	iden := ctx.MustGet("identity").(AdminIdentity)

	type Arg struct {
		FormId uint32 `json:"formId"`
	}
	arg := &Arg{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}

	if model.DeleteForm(iden.At, arg.FormId) {
		ctx.PureJSON(utils.Success()) //成功删除
		return
	} else {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
	}
}
