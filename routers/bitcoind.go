package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupBitcoindRouter(router *gin.Engine) *gin.Engine {
	SetupBitcoindMainnetRouter(router)
	SetupBitcoindTestnetRouter(router)
	SetupBitcoindRegtestRouter(router)
	return router
}

func SetupBitcoindMainnetRouter(router *gin.Engine) *gin.Engine {
	bitcoind := router.Group("/bitcoind")
	bitcoind.Use(middleware.AuthMiddleware())
	mainnet := bitcoind.Group("/mainnet")
	{
		address := mainnet.Group("/address")
		{
			address.GET("/outpoint/:op", handlers.GetAddressByOutpointInMainnet)
			address.POST("/outpoints", handlers.GetAddressesByOutpointSliceInMainnet)
		}
		transaction := mainnet.Group("/transaction")
		{
			transaction.GET("/outpoint/:op", handlers.GetTransactionByOutpointInMainnet)
			transaction.POST("/outpoints", handlers.GetTransactionsByOutpointSliceInMainnet)
		}
		decode := mainnet.Group("/decode")
		{
			decode.GET("/transaction/:tx", handlers.DecodeTransactionInMainnet)
			decode.POST("/transactions", handlers.DecodeTransactionSliceInMainnet)
			decode.POST("/query/transactions", handlers.DecodeAndQueryTransactionSliceInMainnet)
		}
		time := mainnet.Group("/time")
		{
			time.GET("/outpoint/:op", handlers.GetTimeByOutpointInMainnet)
			time.POST("/outpoints", handlers.GetTimesByOutpointSliceInMainnet)
		}
		blockchain := mainnet.Group("/blockchain")
		{
			blockchain.POST("/get_blockchain_info", handlers.GetBlockchainInfoInMainnet)
		}
	}
	return router
}

func SetupBitcoindTestnetRouter(router *gin.Engine) *gin.Engine {
	bitcoind := router.Group("/bitcoind")
	bitcoind.Use(middleware.AuthMiddleware())
	testnet := bitcoind.Group("/testnet")
	{
		address := testnet.Group("/address")
		{
			address.GET("/outpoint/:op", handlers.GetAddressByOutpointInTestnet)
			address.POST("/outpoints", handlers.GetAddressesByOutpointSliceInTestnet)
		}
		transaction := testnet.Group("/transaction")
		{
			transaction.GET("/outpoint/:op", handlers.GetTransactionByOutpointInTestnet)
			transaction.POST("/outpoints", handlers.GetTransactionsByOutpointSliceInTestnet)
		}
		decode := testnet.Group("/decode")
		{
			decode.GET("/transaction/:tx", handlers.DecodeTransactionInTestnet)
			decode.POST("/transactions", handlers.DecodeTransactionSliceInTestnet)
			decode.POST("/query/transactions", handlers.DecodeAndQueryTransactionSliceInTestnet)
		}
		time := testnet.Group("/time")
		{
			time.GET("/outpoint/:op", handlers.GetTimeByOutpointInTestnet)
			time.POST("/outpoints", handlers.GetTimesByOutpointSliceInTestnet)
		}
		blockchain := testnet.Group("/blockchain")
		{
			blockchain.POST("/get_blockchain_info", handlers.GetBlockchainInfoInTestnet)
		}
	}
	return router
}

func SetupBitcoindRegtestRouter(router *gin.Engine) *gin.Engine {
	bitcoind := router.Group("/bitcoind")
	bitcoind.Use(middleware.AuthMiddleware())
	regtest := bitcoind.Group("/regtest")
	{
		address := regtest.Group("/address")
		{
			address.GET("/outpoint/:op", handlers.GetAddressByOutpointInRegtest)
			address.POST("/outpoints", handlers.GetAddressesByOutpointSliceInRegtest)
		}
		transaction := regtest.Group("/transaction")
		{
			transaction.GET("/outpoint/:op", handlers.GetTransactionByOutpointInRegtest)
			transaction.POST("/outpoints", handlers.GetTransactionsByOutpointSliceInRegtest)
		}
		decode := regtest.Group("/decode")
		{
			decode.GET("/transaction/:tx", handlers.DecodeTransactionInRegtest)
			decode.POST("/transactions", handlers.DecodeTransactionSliceInRegtest)
			decode.POST("/query/transactions", handlers.DecodeAndQueryTransactionSliceInRegtest)
		}
		time := regtest.Group("/time")
		{
			time.GET("/outpoint/:op", handlers.GetTimeByOutpointInRegtest)
			time.POST("/outpoints", handlers.GetTimesByOutpointSliceInRegtest)
		}
		blockchain := regtest.Group("/blockchain")
		{
			blockchain.POST("/get_blockchain_info", handlers.GetBlockchainInfoInRegtest)
		}
	}
	return router
}
