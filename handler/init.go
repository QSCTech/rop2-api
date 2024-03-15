package handler

import (
	"github.com/gin-gonic/gin"
)

func Init(routerGroup gin.RouterGroup) {
	authInit()

	orgRoute := routerGroup.Group("/")
	orgInit(*orgRoute)
}
