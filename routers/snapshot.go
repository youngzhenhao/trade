package routers

import (
	"trade/handlers"
	"trade/middleware"

	"github.com/gin-gonic/gin"
)

// download snapshot
func SetupSnapshotRouter(router *gin.Engine) *gin.Engine {
	snapshot := router.Group("/snapshot")
	snapshot.Use(middleware.AuthMiddleware())
	{
		snapshot.GET("/download", handlers.DownloadSnapshot)
	}
	return router
}
