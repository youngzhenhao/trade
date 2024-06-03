package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func setupFairLaunchRouter(router *gin.Engine) *gin.Engine {
	version := router.Group("/v1")
	fairLaunch := version.Group("/fair_launch")
	fairLaunch.Use(middleware.AuthMiddleware())
	{
		fairLaunch.POST("/set", handlers.SetFairLaunchInfo)
		fairLaunch.POST("/mint", handlers.SetFairLaunchMintedInfo)
		fairLaunch.POST("/mint_reserved/:id", handlers.MintFairLaunchReserved)
		query := fairLaunch.Group("/query")
		{
			query.GET("/all", handlers.GetAllFairLaunchInfo)
			query.GET("/issued", handlers.GetIssuedFairLaunchInfo)
			query.GET("/own_set", handlers.GetOwnFairLaunchInfo)
			query.GET("/own_mint", handlers.GetOwnFairLaunchMintedInfo)
			query.GET("/info/:id", handlers.GetFairLaunchInfo)
			query.GET("/minted/:id", handlers.GetMintedInfo)
			query.GET("/inventory/:id", handlers.QueryInventory)
			query.POST("/mint", handlers.QueryMintIsAvailable)
			query.GET("/asset/:id", handlers.GetFairLaunchInfoByAssetId)
		}
	}
	return router
}
