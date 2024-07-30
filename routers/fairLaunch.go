package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupFairLaunchRouter(router *gin.Engine) *gin.Engine {
	version := router.Group("/v1")
	fairLaunch := version.Group("/fair_launch")
	fairLaunch.Use(middleware.AuthMiddleware())
	{
		fairLaunch.POST("/set", handlers.SetFairLaunchInfo)
		fairLaunch.POST("/mint", handlers.SetFairLaunchMintedInfo)
		fairLaunch.POST("/mint_reserved", handlers.MintFairLaunchReserved)
		query := fairLaunch.Group("/query")
		{
			query.GET("/all", handlers.GetAllFairLaunchInfo)
			query.GET("/issued", handlers.GetIssuedFairLaunchInfo)
			query.GET("/closed", handlers.GetClosedFairLaunchInfo)
			query.GET("/not_started", handlers.GetNotStartedFairLaunchInfo)
			query.GET("/own_set", handlers.GetOwnFairLaunchInfo)
			query.GET("/own_set/issued/simplified", handlers.GetOwnFairLaunchInfoIssuedSimplified)
			query.GET("/own_mint", handlers.GetOwnFairLaunchMintedInfo)
			query.GET("/info/:id", handlers.GetFairLaunchInfo)
			query.GET("/minted/:id", handlers.GetMintedInfo)
			query.GET("/inventory/:id", handlers.QueryInventory)
			query.POST("/mint", handlers.QueryMintIsAvailable)
			query.GET("/asset/:asset_id", handlers.GetFairLaunchInfoByAssetId)
			query.GET("/inventory_mint_number/:asset_id", handlers.GetFairLaunchInventoryMintNumberAssetId)
		}
	}
	return router
}
