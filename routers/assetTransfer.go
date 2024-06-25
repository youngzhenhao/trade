package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

//TODO: need to test

func SetupAssetTransferRouter(router *gin.Engine) *gin.Engine {
	assetTransfer := router.Group("/asset_transfer")
	assetTransfer.Use(middleware.AuthMiddleware())
	{
		assetTransfer.GET("/get", handlers.GetAssetTransfer)
		//TODO: need to call after send asset to asset's address
		assetTransfer.POST("/set", handlers.SetAssetTransfer)
	}
	return router
}
