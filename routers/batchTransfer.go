package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
	"trade/middleware"
)

func SetupBatchTransferRouter(router *gin.Engine) *gin.Engine {
	addrReceive := router.Group("/batch_transfer")
	addrReceive.Use(middleware.AuthMiddleware())
	{
		addrReceive.GET("/get", handlers.GetBatchTransfer)
		addrReceive.POST("/set", handlers.SetBatchTransfer)
		addrReceive.POST("/set_slice", handlers.SetBatchTransfers)
	}
	authorized := router.Group("/batch_transfer", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get/asset_id/all/simplified", handlers.GetAllAssetIdAndBatchTransferSimplified)
	return router
}
