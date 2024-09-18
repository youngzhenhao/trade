package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetManagedUtxoRouter(router *gin.Engine) *gin.Engine {
	assetManagedUtxo := router.Group("/asset_managed_utxo")
	assetManagedUtxo.Use(middleware.AuthMiddleware())
	{
		assetManagedUtxo.GET("/get/user", handlers.GetAssetManagedUtxoByUserId)
		// @dev: This two routers may be useless
		{
			assetManagedUtxo.GET("/get/user/ids", handlers.GetAssetManagedUtxoIdsByUserId)
			assetManagedUtxo.GET("/get/user/asset_ids", handlers.GetAssetManagedUtxoAssetIdsByUserId)
		}
		assetManagedUtxo.GET("/get/asset_id/:asset_id", handlers.GetAssetManagedUtxoAssetId)
		assetManagedUtxo.POST("/set", handlers.SetAssetManagedUtxos)
		assetManagedUtxo.POST("/remove", handlers.RemoveAssetManagedUtxos)
		// @dev: Split page
		assetManagedUtxo.POST("/get/limit_offset", handlers.GetAssetManagedUtxoLimitAndOffset)
		assetManagedUtxo.POST("/get/page_number", handlers.GetAssetManagedUtxoPageNumberByPageSize)
	}
	authorized := router.Group("/asset_managed_utxo", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get/all/simplified", handlers.GetAllAssetManagedUtxoSimplified)
	return router
}
