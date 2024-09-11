package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
	"trade/middleware"
)

func SetupAccountAssetRouter(router *gin.Engine) *gin.Engine {
	accountAsset := router.Group("/account_asset")
	accountAsset.Use(middleware.AuthMiddleware())
	{
		balance := accountAsset.Group("/balance")
		{
			balance.GET("/get/asset_id/:asset_id", handlers.GetAccountAssetBalanceByAssetId)
		}
		transfer := accountAsset.Group("/transfer")
		{
			transfer.GET("/get/asset_id/:asset_id", handlers.GetAllAccountAssetTransferByAssetId)
		}

	}
	authorized := router.Group("/account_asset", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get/all/simplified", handlers.GetAllAccountAssetBalanceSimplified)
	return router
}
