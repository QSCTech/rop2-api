package main

import (
	"fmt"
	"rop2-api/handler"
	"rop2-api/model"
	"rop2-api/utils"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var TestOrgId uint32

func main() {
	utils.Init()               //读取配置
	model.Init()               //连接数据库
	if model.CountOrg() <= 0 { //当且仅当没有组织时重建数据库
		fmt.Printf("Starting ResetDb\r\n")
		model.ResetDb()
		fmt.Printf("ResetDb done\r\n")
	}

	server := gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	server.SetTrustedProxies(nil)
	if len(utils.Cfg.CORSAllowOrigins) > 0 {
		println("CORS enabled, allow origins:", strings.Join(utils.Cfg.CORSAllowOrigins, ", "))
		server.Use(cors.New(cors.Config{
			AllowOrigins:     utils.Cfg.CORSAllowOrigins,
			AllowMethods:     []string{"*"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"*"},
			AllowCredentials: false,
			MaxAge:           10 * time.Minute, //根据规范，预检请求缓存时间不超过10min
		}))
	} else {
		println("CORS disabled")
	}

	rootRouter := &server.RouterGroup

	handler.Init(rootRouter)

	server.Run(utils.Cfg.Addr)
}
