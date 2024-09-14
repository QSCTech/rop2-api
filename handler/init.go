package handler

import (
	"bytes"
	"io"
	"log"
	"rop2-api/model"

	"github.com/gin-gonic/gin"
)

func Init(routerGroup *gin.RouterGroup) {
	routerGroup.Use(logContext)

	authInit(routerGroup)

	orgInit(routerGroup)
	formInit(routerGroup)
	resultInit(routerGroup)
	adminInit(routerGroup)
	interviewInit(routerGroup)
	applicantInit(routerGroup)
}

// 记录有登录态的请求日志
func logContext(ctx *gin.Context) {
	//因为body只能读取一次，所以先读取body并存储
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("GetRawData failed", err.Error())
		return
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body)) //确保后续中间件能正常读取body

	//继续处理请求（登录/实际终结点）
	ctx.Next()

	iden, exists := ctx.Get("identity")
	if !exists { //没有登录态，不记录日志
		return
	}
	zjuId := iden.(userIdentity).getId()
	path := ctx.Request.URL.Path
	if ctx.Request.URL.RawQuery != "" {
		path += "?" + ctx.Request.URL.RawQuery
	}

	status := ctx.Writer.Status()
	model.CreateLog(zjuId, ctx.Request.Method, path, &body, int16(status))
}
