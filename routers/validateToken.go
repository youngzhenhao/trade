package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupValidateTokenRouter(router *gin.Engine) *gin.Engine {
	assetLock := router.Group("/validate_token")
	assetLock.Use(middleware.AuthMiddleware())
	{
		assetLock.GET("/ping", handlers.ValidateTokenPing)
	}
	return router
}
