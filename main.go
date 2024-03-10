package main

import (
	"errors"
	"os"
	"rop2-api/handler"
	"rop2-api/model"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func dbInit() error {
	if dsn, ok := os.LookupEnv("ROP2_DSN"); ok {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		db.AutoMigrate(&model.Org{})
		return nil
	} else {
		return errors.New("dns not found")
	}
}

func main() {
	err := dbInit()
	if err != nil {
		println(err)
		return
	}

	server := gin.Default()

	handler.Init(server.RouterGroup)

	//仅供测试连通性
	server.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	server.Run("127.0.0.1:8080") // listen and serve on 127.0.0.1:8080
}
