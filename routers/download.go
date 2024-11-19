package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
)

func SetupDownloadRouter(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/download", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/csv", handlers.CsvDownloadCaptcha)
	return router
}
