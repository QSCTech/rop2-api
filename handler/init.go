package handler

import (
	"github.com/gin-gonic/gin"
)

func Init(routerGroup gin.RouterGroup) {
	orgRoute := routerGroup.Group("/")
	orgInit(*orgRoute)
}
