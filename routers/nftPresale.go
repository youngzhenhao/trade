package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

// TODO: ONLY FOR TEST
func SetupNftPresaleRouter(router *gin.Engine) *gin.Engine {
	nftPresale := router.Group("/nft_presale")
	nftPresale.Use(middleware.AuthMiddleware())
	{
		nftPresale.GET("/get/asset_id/:asset_id", handlers.GetNftPresaleByAssetId)
		nftPresale.POST("/set", handlers.SetNftPresale)
	}
	return router
}
