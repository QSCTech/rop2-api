package main

import (
	"os"
	"rop2-api/handler"
	"rop2-api/model"

	"github.com/gin-gonic/gin"
)

func main() {
	//本地测试用
	os.Setenv("ROP2_DSN", "root:root@tcp(localhost:3306)/rop2?charset=utf8mb4&parseTime=true")

	model.Init()
	model.ResetDb()

	server := gin.Default()
	server.SetTrustedProxies(nil)

	rootRouter := server.RouterGroup
	handler.Init(rootRouter)

	//仅供测试连通性
	rootRouter.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	server.Run("127.0.0.1:8080") // listen and serve on 127.0.0.1:8080
}
