package handler

import (
	"rop2-api/model"
	"rop2-api/utils"

	"github.com/gin-gonic/gin"
)

func orgInit(routerGroup *gin.RouterGroup) {
	orgGroup := routerGroup.Group("/org", AuthWithRefresh(true))
	orgGroup.GET("", getOrgInfo) //对应路径：/org，末尾没有/
	orgGroup.POST("/addDepart", addDepart)
	orgGroup.POST("/deleteDepart", deleteDepart)
	orgGroup.POST("/renameDepart", renameDepart)
}

// 获取登录所在组织（含所有部门）的信息
func getOrgInfo(ctx *gin.Context) {
	id := ctx.MustGet("identity").(*AdminIdentity)
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
	iden := ctx.MustGet("identity").(*AdminIdentity)
	orgId := iden.At

	type AddDepartBody struct {
		Name string
	}
	arg := &AddDepartBody{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}
	name := arg.Name
	lenDiff := utils.LenBetween(name, 1, 20)
	if lenDiff != 0 {
		ctx.AbortWithStatusJSON(utils.MessageInvalidLength(lenDiff < 0))
		return
	}

	if iden.Level < model.Maintainer {
		ctx.AbortWithStatusJSON(utils.MessageForbidden())
		return
	}

	if ok, _ := model.CreateDepart(orgId, arg.Name); ok {
		ctx.PureJSON(utils.Success())
	} else {
		ctx.AbortWithStatusJSON(utils.MessageDuplicate())
	}
}

func deleteDepart(ctx *gin.Context) {
	iden := ctx.MustGet("identity").(*AdminIdentity)
	orgId := iden.At

	type DeleteDepartBody struct {
		Id uint32
	}
	arg := &DeleteDepartBody{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}
	depIdToDelete := arg.Id

	if iden.Level < model.Maintainer {
		ctx.AbortWithStatusJSON(utils.MessageForbidden())
		return
	}

	org := model.GetOrg(orgId)
	if org.DefaultDepart == depIdToDelete {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}
	depToDelete := model.GetDepart(depIdToDelete)
	if depToDelete == nil || depToDelete.Owner != orgId {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}

	model.DeleteDepart(depIdToDelete)
	ctx.PureJSON(utils.Success())
}

func renameDepart(ctx *gin.Context) {
	iden := ctx.MustGet("identity").(*AdminIdentity)
	orgId := iden.At

	type RenameDepartBody struct {
		Id      uint32
		NewName string
	}
	arg := &RenameDepartBody{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}
	depIdToRename := arg.Id
	newName := arg.NewName

	lenDiff := utils.LenBetween(newName, 1, 20)
	if lenDiff != 0 {
		ctx.AbortWithStatusJSON(utils.MessageInvalidLength(lenDiff < 0))
		return
	}

	if iden.Level < model.Maintainer {
		ctx.AbortWithStatusJSON(utils.MessageForbidden())
		return
	}

	//可以重命名默认部门，没有什么作用

	depToRename := model.GetDepart(depIdToRename)
	if depToRename == nil || depToRename.Owner != orgId {
		ctx.AbortWithStatusJSON(utils.MessageNotFound())
		return
	}

	if model.RenameDepart(depIdToRename, newName) {
		ctx.PureJSON(utils.Success())
	} else {
		ctx.AbortWithStatusJSON(utils.MessageDuplicate())
	}
}
