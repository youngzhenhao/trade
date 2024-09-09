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
	}
	return router
}
