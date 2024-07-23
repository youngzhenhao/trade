package routers

import (
	"github.com/gin-gonic/gin"
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
	return router
}
