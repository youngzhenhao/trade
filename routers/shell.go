package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
)

func SetupShellRouter(router *gin.Engine) *gin.Engine {
	shell := router.Group("/shell", gin.BasicAuth(gin.Accounts{
		"foo":   "bar",
		"admin": "123456",
	}))
	{
		shell.GET("/generate/1", handlers.GenerateBlockOne)
	}
	return router
}
