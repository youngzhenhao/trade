package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupPingRouter(router *gin.Engine) *gin.Engine {
	ping := router.Group("/ping")
	ping.Use(middleware.AuthMiddleware())
	{
		ping.GET("/ip_test", handlers.PingIpTestToken)
	}
	return router
}
