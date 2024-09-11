package routers

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"trade/config"
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
			query.GET("/followed", handlers.GetFollowedFairLaunchInfo)
			query.GET("/hot", handlers.GetHotFairLaunchInfo)
			query.GET("/issued", handlers.GetIssuedFairLaunchInfo)
			query.GET("/not_started", handlers.GetNotStartedFairLaunchInfo)
			query.GET("/closed", handlers.GetClosedFairLaunchInfo)
			query.GET("/own_set", handlers.GetOwnFairLaunchInfo)
			query.GET("/own_set/issued/simplified", handlers.GetOwnFairLaunchInfoIssuedSimplified)
			query.GET("/own_mint", handlers.GetOwnFairLaunchMintedInfo)
			query.GET("/info/:id", handlers.GetFairLaunchInfo)
			query.GET("/minted/:id", handlers.GetMintedInfo)
			//query.GET("/inventory/:id", handlers.QueryInventory)
			query.POST("/mint", handlers.QueryMintIsAvailable)
			query.GET("/asset/:asset_id", handlers.GetFairLaunchInfoByAssetId)
			query.GET("/asset/plus_info/:asset_id", handlers.GetFairLaunchInfoPlusByAssetId)
			//query.GET("/inventory_mint_number/:asset_id", handlers.GetFairLaunchInventoryMintNumberAssetId)
		}
	}
	username := base64.StdEncoding.EncodeToString([]byte(config.GetLoadConfig().AdminUser.Username))
	password := base64.StdEncoding.EncodeToString([]byte(config.GetLoadConfig().AdminUser.Password))
	authorized := router.Group("/fair_launch/auth_op", gin.BasicAuth(gin.Accounts{
		username: password,
	}))
	authorized.POST("/refund", handlers.RefundUserFirstMintByUsernameAndAssetId)
	return router
}
