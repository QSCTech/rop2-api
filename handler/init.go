package handler

import (
	"github.com/gin-gonic/gin"
)

func Init(routerGroup *gin.RouterGroup) {
	authInit(routerGroup)

	orgInit(routerGroup)
	formInit(routerGroup)
	resultInit(routerGroup)
	adminInit(routerGroup)
	interviewInit(routerGroup)
	applicantInit(routerGroup)
}
