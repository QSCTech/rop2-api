package main

import (
	"rop2-api/handler"
	"rop2-api/model"
	"rop2-api/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.Init() //读取配置

	model.Init()
	model.ResetDb()
	testOrgId, _ := model.InitNewOrg("测试组织", "N/A", "测试管理员")
	model.CreateDepart(testOrgId, "部门1")
	model.InitNewOrg("测试组织2", "N/A", "管理员2")

	server := gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	server.SetTrustedProxies(nil)
	server.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Content-Type", "Rop-Token"},
		ExposeHeaders: []string{"Rop-Refresh-Token"},
		MaxAge:        12 * time.Hour,
	}))

	rootRouter := &server.RouterGroup

	handler.Init(rootRouter)

	server.Run("127.0.0.1:8080")
}
