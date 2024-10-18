package routers

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
	"trade/middleware"
)

func SetupNftPresaleRouter(router *gin.Engine) *gin.Engine {
	nftPresale := router.Group("/nft_presale")
	nftPresale.Use(middleware.AuthMiddleware())
	{
		get := nftPresale.Group("/get")
		{
			get.GET("/asset_id", handlers.GetNftPresaleByAssetId)
			get.GET("/batch_group_id", handlers.GetNftPresaleByBatchGroupId)

			// @dev: Deprecated
			{
				//get.GET("/launched", handlers.GetLaunchedNftPresale)
				//get.GET("/user_bought", handlers.GetUserBoughtNftPresale)
				//get.GET("/group_key", handlers.GetNftPresaleByGroupKeyPurchasable)
				//get.GET("/no_group_key", handlers.GetNftPresaleNoGroupKeyPurchasable)
			}
		}
		{
			query := nftPresale.Group("/query")
			query.GET("/batch_group", handlers.QueryNftPresaleBatchGroup)

			// @dev: Deprecated
			//{
			//	query.GET("/group_key", handlers.QueryNftPresaleGroupKeyPurchasable)
			//}
		}
		nftPresale.POST("/buy", handlers.BuyNftPresale)
	}
	username := base64.StdEncoding.EncodeToString([]byte(config.GetLoadConfig().AdminUser.Username))
	password := base64.StdEncoding.EncodeToString([]byte(config.GetLoadConfig().AdminUser.Password))
	authorized := router.Group("/nft_presale/auth_op", gin.BasicAuth(gin.Accounts{
		username: password,
	}))
	{
		authorized.POST("/launch", handlers.LaunchNftPresaleBatchGroup)
		authorized.POST("/add_whitelists", handlers.AddNftPresaleWhitelists)

		// @dev: Deprecated temporarily
		{
			//authorized.POST("/set", handlers.SetNftPresale)
			//authorized.POST("/set/batch", handlers.SetNftPresales)
			//authorized.POST("/reset", handlers.ReSetFailOrCanceledNftPresale)}
		}
	}
	return router
}
