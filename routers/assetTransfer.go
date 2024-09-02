package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
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
	query := assetTransfer.Group("/query/:txid")
	{
		query.GET("/txid", handlers.GetAssetTransferByTxid)
	}
	authorized := router.Group("/asset_transfer", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get/combined/all/simplified", handlers.GetAllAssetTransferSimplified)
	authorized.GET("/get/combined/asset_id/all/simplified", handlers.GetAllAssetIdAndAssetTransferSimplified)
	return router
}
