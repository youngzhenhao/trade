package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
)

func SetupDownloadRouter(router *gin.Engine) *gin.Engine {
	download := router.Group("/download")
	download.GET("/csv", handlers.CsvDownloadCaptcha)
	return router
}
