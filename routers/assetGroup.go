package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetGroupRouter(router *gin.Engine) *gin.Engine {
	assetGroup := router.Group("/asset_group")
	assetGroup.Use(middleware.AuthMiddleware())
	{
		assetGroup.GET("/get/first_meta/group_key/:group_key", handlers.GetGroupFirstAssetMeta)
		assetGroup.POST("/set/first_meta/", handlers.SetGroupFirstAssetMeta)
	}
	return router
}
