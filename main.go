package main

import (
	"rop2-api/handler"
	"rop2-api/model"

	"github.com/gin-gonic/gin"
)

func main() {
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
