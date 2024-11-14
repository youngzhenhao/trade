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
	authorized.GET("/specified", handlers.GetSpecifiedDateUserStats)
	authorized.GET("/csv", handlers.DownloadCsv)

	// TODO: ONLY FOR TEST
	authorized.GET("/count", handlers.GetDateLoginCount)
	authorized.GET("/record", handlers.GetDateIpLoginRecord)
	authorized.GET("/record_count", handlers.GetDateIpLoginRecordCount)
	authorized.GET("/new_count", handlers.GetNewUserCount)

	return router
}
