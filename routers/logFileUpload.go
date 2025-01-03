package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
)

func SetupLogFileUploadRouter(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/log_file_upload", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	{
		authorized.POST("/upload", handlers.UploadLogFile)
		authorized.POST("/upload_big", handlers.UploadBigFile)
		authorized.GET("/get/all/device_id", handlers.GetAllLogFilesByDeviceId)
		authorized.GET("/get/all", handlers.GetAllLogFiles)
		authorized.GET("/get/download/id/:id", handlers.DownloadLogFileById)
	}
	return router
}
