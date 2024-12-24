package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupBtcUtxoRouter(router *gin.Engine) *gin.Engine {
	btcUtxo := router.Group("/btc_utxo")
	btcUtxo.Use(middleware.AuthMiddleware())

	btcUtxo.POST("set", handlers.SetBtcUtxo)
	return router
}
