package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
	"trade/middleware"
)

func SetupBtcBalanceRouter(router *gin.Engine) *gin.Engine {
	snapshot := router.Group("/btc_balance")
	snapshot.Use(middleware.AuthMiddleware())
	{
		snapshot.GET("/get", handlers.GetBtcBalance)
		snapshot.POST("/set", handlers.SetBtcBalance)
	}
	authorized := router.Group("/btc_balance", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get/btc_balance_rank/count", handlers.GetBtcBalanceCount)
	authorized.GET("/get/btc_balance_rank", handlers.GetBtcBalanceOrderLimitOffset)
	return router
}
