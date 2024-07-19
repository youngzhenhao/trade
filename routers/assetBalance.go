package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetBalanceRouter(router *gin.Engine) *gin.Engine {
	assetBalance := router.Group("/asset_balance")
	assetBalance.Use(middleware.AuthMiddleware())
	{
		assetBalance.GET("/get", handlers.GetAssetBalance)
		assetBalance.GET("/get/holder/number/:asset_id", handlers.GetAssetHolderNumber)
		assetBalance.GET("/get/holder/balance/:asset_id", handlers.GetAssetHolderBalance)
		assetBalance.POST("/get/holder/balance/limit_offset", handlers.GetAssetHolderBalanceLimitAndOffset)
		assetBalance.GET("/get/holder/balance/page/:asset_id", handlers.GetAssetHolderBalancePage)
		assetBalance.POST("/set", handlers.SetAssetBalance)
		assetBalance.POST("/set_slice", handlers.SetAssetBalances)
	}
	return router
}
