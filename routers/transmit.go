package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupTransmitRouter(router *gin.Engine) *gin.Engine {

	transmit := router.Group("/TradeTransmit")
	{
		transmit.POST("/login", handlers.LoginHandler)
	}
	// A routing group that requires authentication
	custody := transmit.Group("/custodyAccount")
	custody.Use(middleware.AuthMiddleware())
	{
		custody.POST("/create", handlers.CreateCustodyAccount)
		Invoice := custody.Group("/invoice")
		{
			Invoice.POST("/apply", handlers.ApplyInvoice)
			Invoice.POST("/pay", handlers.PayInvoice)
			Invoice.POST("/querybalance", handlers.QueryBalance)
			Invoice.POST("/queryinvoice", handlers.QueryInvoice)
			Invoice.POST("/querypayment", handlers.QueryPayment)
			Invoice.POST("/lookupinvoice", handlers.LookupInvoice)
		}
		Asset := custody.Group("/Asset")
		{
			Asset.POST("/queryassets", handlers.QueryAssets)
		}
	}
	return router
}
