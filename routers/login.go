package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupLoginRouter(router *gin.Engine) *gin.Engine {
	router.POST("/getNonce", handlers.GetNonceHandler)
	router.POST("/getDeviceId", handlers.GetDeviceIdHandler)
	// Login routing
	router.POST("/login", handlers.LoginHandler)
	// Refresh the route for the token
	router.POST("/refresh", handlers.RefreshTokenHandler)
	// A routing group that requires authentication
	auth := router.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/getConfig", handlers.GetConfigHandler)
		auth.POST("/setConfig", handlers.SetConfigHandler)
		auth.GET("/userinfo", handlers.UserInfoHandler)
	}
	return router
}
