package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetRecommendRouter(router *gin.Engine) *gin.Engine {
	assetLocalMint := router.Group("/asset_recommend")
	assetLocalMint.Use(middleware.AuthMiddleware())
	{
		assetLocalMint.GET("/get/user", handlers.GetAssetRecommendByUserId)
		assetLocalMint.GET("/get/asset_id/:asset_id", handlers.GetAssetRecommendByAssetId)
		assetLocalMint.GET("/get/user/asset_id/:asset_id", handlers.GetUserAssetRecommendByAssetId)
		assetLocalMint.POST("/get/user_id_and_asset_id", handlers.GetAssetRecommendByUserIdAndAssetId)
		assetLocalMint.POST("/set", handlers.SetAssetRecommend)
	}
	authorized := router.Group("/asset_recommend", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get/all/simplified", handlers.GetAllAssetRecommendSimplified)
	return router
}
