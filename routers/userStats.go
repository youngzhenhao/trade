package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
)

func SetupUserStatsRouter(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/user_stats", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get", handlers.GetUserStats)
	authorized.GET("/csv", handlers.DownloadCsv)
	return router
}
