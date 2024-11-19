package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
)

func SetupBackRewardRouter(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/back_reward", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get", handlers.GetBackRewards)
	return router
}
