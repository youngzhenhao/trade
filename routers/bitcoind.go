package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupBitcoindRouter(router *gin.Engine) *gin.Engine {
	bitcoind := router.Group("/bitcoind")
	bitcoind.Use(middleware.AuthMiddleware())
	mainnet := bitcoind.Group("/mainnet")
	{
		address := mainnet.Group("/address")
		{
			address.GET("/outpoint/:op", handlers.GetAddressByOutpointInMainnet)
		}
	}
	testnet := bitcoind.Group("/testnet")
	{
		address := testnet.Group("/address")
		{
			address.GET("/outpoint/:op", handlers.GetAddressByOutpointInTestnet)
		}
	}
	regtest := bitcoind.Group("/regtest")
	{
		address := regtest.Group("/address")
		{
			address.GET("/outpoint/:op", handlers.GetAddressByOutpointInRegtest)
		}
	}
	return router
}
