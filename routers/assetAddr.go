package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetAddrRouter(router *gin.Engine) *gin.Engine {
	assetAddr := router.Group("/asset_addr")
	assetAddr.Use(middleware.AuthMiddleware())
	{
		assetAddr.GET("/get", handlers.GetAssetAddr)
		assetAddr.GET("/get/script_key/:script_key", handlers.GetAssetAddrByScriptKey)
		assetAddr.GET("/get/encoded/:encoded", handlers.GetAssetAddrByEncoded)
		assetAddr.GET("/migrate/update", handlers.UpdateUsernameByUserId)
		assetAddr.POST("/set", handlers.SetAssetAddr)
	}
	authorized := router.Group("/asset_addr", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get/all", handlers.GetAllAssetAddrs)
	authorized.GET("/get/all/simplified", handlers.GetAllAssetAddrSimplified)

	return router
}
