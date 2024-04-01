package main

import (
	"rop2-api/handler"
	"rop2-api/model"
	"rop2-api/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.Init() //读取配置

	model.Init()
	model.ResetDb()
	model.InitNewOrg("测试组织", "N/A", "测试管理员")
	model.InitNewOrg("测试组织2", "N/A", "管理员2")

	server := gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	server.SetTrustedProxies(nil)
	server.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*") //TODO 设成具体的域
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Rop-Token")
		ctx.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET")
		ctx.Header("Access-Control-Expose-Headers", "Rop-Refresh-Token")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	})

	rootRouter := &server.RouterGroup

	handler.Init(rootRouter)

	server.Run("127.0.0.1:8080") // listen and serve on 127.0.0.1:8080
}
