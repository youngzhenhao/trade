package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
)

func SetupLogFileUploadRouter(router *gin.Engine) *gin.Engine {
	logFileUpload := router.Group("/log_file_upload")
	{
		logFileUpload.POST("/upload", handlers.UploadLogFile)
	}
	return router
}
