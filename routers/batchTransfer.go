package routers

import (
	"github.com/gin-gonic/gin"
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
	return router
}
