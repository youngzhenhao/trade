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
		assetBalance.POST("/set", handlers.SetAssetBalance)
		assetBalance.POST("/set_slice", handlers.SetAssetBalances)
	}
	return router
}
