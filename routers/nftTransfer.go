package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupNftTransferRouter(router *gin.Engine) *gin.Engine {
	nftTransfer := router.Group("/nft_transfer")
	nftTransfer.Use(middleware.AuthMiddleware())
	{
		nftTransfer.GET("/get/asset_id/:asset_id", handlers.GetNftTransferByAssetId)
		nftTransfer.POST("/set", handlers.SetNftTransfer)
	}
	return router
}
