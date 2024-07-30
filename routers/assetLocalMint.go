package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetLocalMintRouter(router *gin.Engine) *gin.Engine {
	assetLocalMint := router.Group("/asset_local_mint")
	assetLocalMint.Use(middleware.AuthMiddleware())
	{
		assetLocalMint.GET("/get/user", handlers.GetAssetLocalMintByUserId)
		assetLocalMint.GET("/get/asset_id/:asset_id", handlers.GetAssetLocalMintAssetId)
		assetLocalMint.POST("/set", handlers.SetAssetLocalMint)
		assetLocalMint.POST("/set/slice", handlers.SetAssetLocalMints)
	}
	authorized := router.Group("/asset_local_mint", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get/all/simplified", handlers.GetAllAssetLocalMintSimplified)
	return router
}
