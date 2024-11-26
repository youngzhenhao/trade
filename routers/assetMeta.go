package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetMetaRouter(router *gin.Engine) *gin.Engine {
	assetMeta := router.Group("/asset_meta")
	assetMeta.Use(middleware.AuthMiddleware())
	assetMeta.POST("/image/query", handlers.GetAssetMetaImage)
	return router
}
