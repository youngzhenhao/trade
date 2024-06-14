package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupProofRouter(router *gin.Engine) *gin.Engine {
	proof := router.Group("/proof")
	proof.Use(middleware.AuthMiddleware())
	{
		proof.POST("/download", handlers.DownloadProof)
	}
	return router
}
