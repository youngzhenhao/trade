package routers

import (
	"github.com/gin-gonic/gin"
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
	return router
}
