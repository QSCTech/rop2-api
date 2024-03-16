package handler

import (
	"rop2-api/model"
	"rop2-api/utils"

	"github.com/gin-gonic/gin"
)

func orgInit(routerGroup *gin.RouterGroup) {
	orgGroup := routerGroup.Group("/org", AuthWithRefresh(true))
	orgGroup.GET("/", getOrgInfo)
	orgGroup.POST("/addDepart", addDepart)
	orgGroup.POST("/deleteDepart", deleteDepart)
	orgGroup.POST("/renameDepart", renameDepart)
}

// 获取登录所在组织（含所有部门）的信息
func getOrgInfo(ctx *gin.Context) {
	id := ctx.MustGet("identity").(*UserIdentity)
	//只要能以组织的身份登录，就可查询，对具体权限级别无要求
	orgId := id.At
	org := model.GetOrg(orgId)
	if org == nil {
		ctx.PureJSON(204, gin.H{})
	} else {
		departs := model.GetOrgDeparts(orgId)
		ctx.PureJSON(200, gin.H{
			"org":     org,
			"departs": departs,
		})
	}
}

func addDepart(ctx *gin.Context) {
	iden := ctx.MustGet("identity").(*UserIdentity)
	orgId := iden.At

	type AddDepartBody struct {
		Name string
	}
	arg := &AddDepartBody{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.Message("参数绑定失败", 400, 0))
		return
	}
	name := arg.Name
	if !utils.LenBetween(name, 1, 20) {
		ctx.AbortWithStatusJSON(utils.Message("名称长度无效", 400, 1))
		return
	}

	if !model.HasOrgLevel(iden.Perm, orgId, model.Maintainer) {
		ctx.AbortWithStatusJSON(utils.Message("权限不足", 403, 1))
		return
	}

	if model.CreateDepart(orgId, arg.Name) {
		ctx.PureJSON(utils.Success())
	} else {
		ctx.AbortWithStatusJSON(utils.Message("部门命名重复", 422, 11))
	}
}

func deleteDepart(ctx *gin.Context) {
	iden := ctx.MustGet("identity").(*UserIdentity)
	orgId := iden.At

	type DeleteDepartBody struct {
		Id uint32
	}
	arg := &DeleteDepartBody{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.Message("参数绑定失败", 400, 0))
		return
	}
	depIdToDelete := arg.Id

	if !model.HasOrgLevel(iden.Perm, orgId, model.Maintainer) {
		ctx.AbortWithStatusJSON(utils.Message("权限不足", 403, 1))
		return
	}

	depToDelete := model.GetDepart(depIdToDelete)
	if depToDelete == nil || depToDelete.Parent != orgId {
		ctx.AbortWithStatusJSON(utils.Message("部门不存在", 422, 2))
		return
	}

	if model.DeleteDepart(depIdToDelete) {
		ctx.PureJSON(utils.Success())
	} else {
		ctx.AbortWithStatusJSON(utils.Message("无法删除默认部门", 422, 1))
	}
}

func renameDepart(ctx *gin.Context) {
	iden := ctx.MustGet("identity").(*UserIdentity)
	orgId := iden.At

	type RenameDepartBody struct {
		Id      uint32
		NewName string
	}
	arg := &RenameDepartBody{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.Message("参数绑定失败", 400, 0))
		return
	}
	depIdToRename := arg.Id
	newName := arg.NewName

	if !utils.LenBetween(newName, 1, 20) {
		ctx.AbortWithStatusJSON(utils.Message("名称长度无效", 400, 1))
		return
	}

	if !model.HasOrgLevel(iden.Perm, orgId, model.Maintainer) {
		ctx.AbortWithStatusJSON(utils.Message("权限不足", 403, 1))
		return
	}

	//可以重命名默认部门，没有什么作用

	depToRename := model.GetDepart(depIdToRename)
	if depToRename == nil || depToRename.Parent != orgId {
		ctx.AbortWithStatusJSON(utils.Message("部门不存在", 422, 2))
		return
	}

	if model.RenameDepart(depIdToRename, newName) {
		ctx.PureJSON(utils.Success())
	} else {
		ctx.AbortWithStatusJSON(utils.Message("部门命名重复", 422, 11))
	}
}
