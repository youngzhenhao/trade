package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
	"trade/middleware"
)

// TODO: need to test
func SetupAssetLocalMintHistoryRouter(router *gin.Engine) *gin.Engine {
	assetLocalMint := router.Group("/asset_local_mint_history")
	assetLocalMint.Use(middleware.AuthMiddleware())
	{
		assetLocalMint.GET("/get/user", handlers.GetAssetLocalMintHistoryByUserId)
		assetLocalMint.GET("/get/asset_id/:asset_id", handlers.GetAssetLocalMintHistoryAssetId)
		assetLocalMint.POST("/set", handlers.SetAssetLocalMintHistories)
	}
	authorized := router.Group("/asset_local_mint_history", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get/all/simplified", handlers.GetAllAssetLocalMintHistorySimplified)
	return router
}
