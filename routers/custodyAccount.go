package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupCustodyAccountRouter(router *gin.Engine) *gin.Engine {

	// A routing group that requires authentication
	custody := router.Group("/custodyAccount")

	custody.Use(middleware.AuthMiddleware())
	{
		custody.POST("/balance", handlers.GetBalance)
		Invoice := custody.Group("/invoice")
		{
			Invoice.POST("/apply", handlers.ApplyInvoice)
			Invoice.POST("/pay", handlers.PayInvoice)
			Invoice.POST("/querybalance", handlers.QueryBalance)
			Invoice.POST("/queryinvoice", handlers.QueryInvoice)
			Invoice.POST("/querypayment", handlers.QueryPayment)
			Invoice.POST("/decodeinvoice", handlers.DecodeInvoice)

		}
		Asset := custody.Group("/Asset")
		{
			Asset.POST("/apply", handlers.ApplyAddress)
			Asset.POST("/send", handlers.SendAsset)
			Asset.POST("/queryasset", handlers.QueryAsset)
			Asset.POST("/queryassets", handlers.QueryAssets)
			Asset.POST("/queryaddress", handlers.QueryAddress)
			Asset.POST("/queryaddresses", handlers.QueryAddresses)
			Asset.POST("/querypayment", handlers.QueryAssetPayment)
			Asset.POST("/querypayments", handlers.QueryAssetPayments)
			Asset.POST("/decodeaddr", handlers.DecodeAddress)
		}
	}
	return router
}
