package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
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
	authorized := router.Group("/addr_receive", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get/all/simplified", handlers.GetAllAddrReceiveSimplified)

	return router
}
