package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
)

func SetupAssetMetaRouter(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/asset_meta", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.POST("/query", handlers.GetAssetMeta)
	return router
}
