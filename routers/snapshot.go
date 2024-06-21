package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupSnapshotRouter(router *gin.Engine) *gin.Engine {
	snapshot := router.Group("/snapshot")
	snapshot.Use(middleware.AuthMiddleware())
	{
		snapshot.GET("/download", handlers.DownloadSnapshot)
	}
	return router
}
