package handler

import (
	"time"

	"github.com/gin-gonic/gin"
)

func Init(routerGroup gin.RouterGroup) {
	orgRoute := routerGroup.Group("/")
	orgInit(*orgRoute)

	loginCleanupTicker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			time := <-loginCleanupTicker.C
			for k, v := range loginMap {
				if time.Compare(v.notAfter) > 0 {
					delete(loginMap, k)
				}
			}
		}
	}()
}
