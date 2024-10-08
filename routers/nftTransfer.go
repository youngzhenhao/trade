package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupNftTransferRouter(router *gin.Engine) *gin.Engine {
	assetGroup := router.Group("/nft_transfer")
	assetGroup.Use(middleware.AuthMiddleware())
	{
		assetGroup.GET("/get/asset_id/:asset_id", handlers.GetNftTransferByAssetId)
		assetGroup.POST("/set", handlers.SetNftTransfer)
	}
	return router
}
