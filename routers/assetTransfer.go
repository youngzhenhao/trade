package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetTransferRouter(router *gin.Engine) *gin.Engine {
	assetTransfer := router.Group("/asset_transfer")
	assetTransfer.Use(middleware.AuthMiddleware())
	{
		assetTransfer.GET("/get", handlers.GetAssetTransfer)
		assetTransfer.GET("/get/:asset_id", handlers.GetAssetTransferByAssetId)
		assetTransfer.GET("/get/txids", handlers.GetAssetTransferTxids)
		assetTransfer.POST("/set", handlers.SetAssetTransfer)
	}
	return router
}
