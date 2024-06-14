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
		proof.GET("/download/:asset_id/:proof_name", handlers.DownloadProof)
		proof.POST("/download2/:asset_id/:proof_name", handlers.DownloadProof2)
	}
	return router
}
