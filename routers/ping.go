package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
)

func SetupPingRouter(router *gin.Engine) *gin.Engine {
	version := router.Group("/v1")
	authorized := version.Group("/admin", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/ping", handlers.PingHandler)
	return router
}
