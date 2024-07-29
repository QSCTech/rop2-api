package main

import (
	"fmt"
	"rop2-api/handler"
	"rop2-api/model"
	"rop2-api/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var TestOrgId uint32

func main() {
	utils.Init() //读取配置

	model.Init()

	if utils.DoResetDb || model.CountOrg() <= 0 {
		fmt.Printf("Starting ResetDb\r\n")
		model.ResetDb()
		fmt.Printf("ResetDb done\r\n")
	}

	server := gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	server.SetTrustedProxies(nil)
	server.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"https://localhost:5173", "http://localhost:5173", "https://www.qsc.zju.edu.cn", "http://www.qsc.zju.edu.cn"},
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Content-Type", "Rop-Token"},
		ExposeHeaders: []string{"Rop-Refresh-Token"},
		MaxAge:        10 * time.Minute, //根据规范，预检请求缓存时间不超过10min
	}))

	rootRouter := &server.RouterGroup

	handler.Init(rootRouter)

	server.Run(utils.BindAddr)
}
