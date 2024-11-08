package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetListRouter(router *gin.Engine) *gin.Engine {
	assetList := router.Group("/asset_list")
	assetList.Use(middleware.AuthMiddleware())
	{
		assetList.GET("/get", handlers.GetAssetList)
		assetList.POST("/set_slice", handlers.SetAssetLists)
		assetList.GET("/is_exist", handlers.IsAssetListRecordExist)
	}
	return router
}
