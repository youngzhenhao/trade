package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetAddrRouter(router *gin.Engine) *gin.Engine {
	assetAddr := router.Group("/asset_addr")
	assetAddr.Use(middleware.AuthMiddleware())
	{
		assetAddr.GET("/get", handlers.GetAssetAddr)
		assetAddr.GET("/get/script_key/:script_key", handlers.GetAssetAddrByScriptKey)
		assetAddr.POST("/set", handlers.SetAssetAddr)
	}
	return router
}
