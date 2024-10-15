package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
)

func SetupUserDataRouter(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/user_data", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get", handlers.GetUserData)
	authorized.GET("/get/yaml", handlers.GetUserDataYaml)
	return router
}
