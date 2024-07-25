package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetBalanceRouter(router *gin.Engine) *gin.Engine {
	assetBalance := router.Group("/asset_balance")
	assetBalance.Use(middleware.AuthMiddleware())
	{
		assetBalance.GET("/get", handlers.GetAssetBalance)
		assetBalance.GET("/get/holder/number/:asset_id", handlers.GetAssetHolderNumber)
		// TODO: This router should be deprecated
		assetBalance.GET("/get/holder/balance/limit50/:asset_id", handlers.GetAssetHolderBalance)
		assetBalance.GET("/get/holder/balance/records/:asset_id", handlers.GetAssetHolderBalanceRecordsNumber)
		assetBalance.POST("/get/holder/balance/limit_offset", handlers.GetAssetHolderBalanceLimitAndOffset)
		assetBalance.POST("/set", handlers.SetAssetBalance)
		assetBalance.POST("/set_slice", handlers.SetAssetBalances)
	}
	authorized := router.Group("/asset_balance", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get/holder/username/balance/all", handlers.GetAssetHolderUsernameBalanceAll)
	authorized.GET("/get/holder/username/balance/all/simplified", handlers.GetAssetHolderUsernameBalanceAllSimplified)
	return router
}
