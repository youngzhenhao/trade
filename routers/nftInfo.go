package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupNftInfoRouter(router *gin.Engine) *gin.Engine {
	nftInfo := router.Group("/nft_info")
	nftInfo.Use(middleware.AuthMiddleware())
	{
		nftInfo.GET("/get/asset_id/:asset_id", handlers.GetNftInfoByAssetId)
		nftInfo.POST("/set", handlers.SetNftInfo)
	}
	return router
}
