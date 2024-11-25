package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetBalanceHistoryRouter(router *gin.Engine) *gin.Engine {
	assetBalanceHistory := router.Group("/asset_balance_history")
	assetBalanceHistory.Use(middleware.AuthMiddleware())
	{
		get := assetBalanceHistory.Group("/get")
		{
			get.GET("/latest", handlers.GetLatestAssetBalanceHistories)
		}
		assetBalanceHistory.POST("/create", handlers.CreateAssetBalanceHistories)
	}
	return router
}
