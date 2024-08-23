package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
)

func SetupWsRouter(router *gin.Engine) *gin.Engine {
	router.GET("/ws", func(c *gin.Context) {
		handlers.WsHandler(c.Writer, c.Request)
	})
	return router
}
