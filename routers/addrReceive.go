package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupAddrReceiveRouter(router *gin.Engine) *gin.Engine {
	addrReceive := router.Group("/addr_receive")
	addrReceive.Use(middleware.AuthMiddleware())
	{
		addrReceive.GET("/get", handlers.GetAddrReceive)
		addrReceive.GET("/get/origin", handlers.GetAddrReceiveOrigin)
		addrReceive.POST("/set", handlers.SetAddrReceive)
	}
	return router
}
