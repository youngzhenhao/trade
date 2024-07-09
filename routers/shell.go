package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
)

func SetupShellRouter(router *gin.Engine) *gin.Engine {
	shell := router.Group("/shell", gin.BasicAuth(basicAuthAccounts))
	{
		shell.GET("/generate/1", handlers.GenerateBlockOne)
		shell.GET("/faucet/0.1/:address", handlers.FaucetTransferOneTenthBtc)
		shell.GET("/faucet/0.01/:address", handlers.FaucetTransferOneHundredthBtc)
	}
	return router
}
