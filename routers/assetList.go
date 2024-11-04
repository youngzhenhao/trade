package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetListRouter(router *gin.Engine) *gin.Engine {
	assetBalance := router.Group("/asset_list")
	assetBalance.Use(middleware.AuthMiddleware())
	{
		assetBalance.POST("/set_slice", handlers.SetAssetLists)
	}
	return router
}
