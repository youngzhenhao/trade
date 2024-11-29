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

	network := assetMeta.Group("/network")
	{
		mainnet := network.Group("/mainnet")
		{
			mainnet.GET("/image/group_first", handlers.GetGroupFirstImageDataInMainnet)
		}
		testnet := network.Group("/testnet")
		{
			testnet.GET("/image/group_first", handlers.GetGroupFirstImageDataInTestnet)
		}
		regtest := network.Group("/regtest")
		{
			regtest.GET("/image/group_first", handlers.GetGroupFirstImageDataInRegtest)
		}
	}
	return router
}
