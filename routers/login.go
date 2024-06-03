package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupLoginRouter(router *gin.Engine) *gin.Engine {
	// Login routing
	router.POST("/login", handlers.LoginHandler)
	// Refresh the route for the token
	router.POST("/refresh", handlers.RefreshTokenHandler)
	// A routing group that requires authentication
	auth := router.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/userinfo", handlers.UserInfoHandler)
	}
	return router
}
