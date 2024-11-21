package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupAssetBalanceBackupRouter(router *gin.Engine) *gin.Engine {
	assetBalanceBackup := router.Group("/asset_balance_backup")
	assetBalanceBackup.Use(middleware.AuthMiddleware())
	{
		assetBalanceBackup.GET("/get", handlers.GetAssetBalanceBackup)
		assetBalanceBackup.POST("/update", handlers.UpdateAssetBalanceBackup)
	}
	return router
}
