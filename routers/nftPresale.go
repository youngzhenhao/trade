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
		{
			nftPresale.GET("/get/asset_id/:asset_id", handlers.GetNftPresaleByAssetId)
			nftPresale.GET("/get/launched", handlers.GetLaunchedNftPresale)
			nftPresale.GET("/get/user_bought", handlers.GetUserBoughtNftPresale)
		}
		nftPresale.POST("/buy", handlers.BuyNftPresale)

	}
	username := base64.StdEncoding.EncodeToString([]byte(config.GetLoadConfig().AdminUser.Username))
	password := base64.StdEncoding.EncodeToString([]byte(config.GetLoadConfig().AdminUser.Password))
	authorized := router.Group("/nft_presale/auth_op", gin.BasicAuth(gin.Accounts{
		username: password,
	}))
	{
		authorized.POST("/set", handlers.SetNftPresale)
		authorized.POST("/set/batch", handlers.SetNftPresales)
	}
	return router
}
