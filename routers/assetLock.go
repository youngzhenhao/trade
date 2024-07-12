package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetLockRouter(router *gin.Engine) *gin.Engine {
	assetLock := router.Group("/asset_lock")
	assetLock.Use(middleware.AuthMiddleware())
	{
		assetLock.GET("/get", handlers.GetAssetLock)
		assetLock.POST("/set", handlers.SetAssetLock)
	}
	return router
}
