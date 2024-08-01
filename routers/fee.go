package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupFeeRouter(router *gin.Engine) *gin.Engine {
	version := router.Group("/v1")
	fee := version.Group("/fee")
	fee.Use(middleware.AuthMiddleware())
	{
		query := fee.Group("/query")
		{
			query.GET("/rate", handlers.QueryFeeRate)
			query.GET("/recommended", handlers.QueryRecommendedFeeRate)
			fairLaunch := query.Group("/fair_launch")
			{
				fairLaunch.GET("/issuance", handlers.QueryFairLaunchIssuanceFee)
				fairLaunch.GET("/mint", handlers.QueryFairLaunchMintFee)
			}
			//query.GET("/all", handlers.QueryAllFeeRate)
		}
	}
	return router
}
