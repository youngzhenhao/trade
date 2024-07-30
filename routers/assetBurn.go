package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetBurnRouter(router *gin.Engine) *gin.Engine {
	assetBurn := router.Group("/asset_burn")
	assetBurn.Use(middleware.AuthMiddleware())
	{
		assetBurn.GET("/get/user", handlers.GetAssetBurnByUserId)
		assetBurn.GET("/get/asset_id/:asset_id", handlers.GetAssetBurnByAssetId)
		assetBurn.POST("/set", handlers.SetAssetBurn)
	}
	authorized := router.Group("/asset_burn", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get/all/simplified", handlers.GetAllAssetBurnSimplified)
	return router
}
