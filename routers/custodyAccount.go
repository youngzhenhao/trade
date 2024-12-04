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
		custody.POST("/getAssetBalanceList", handlers.GetAssetBalanceList)
		Invoice := custody.Group("/invoice")
		{
			Invoice.POST("/apply", handlers.ApplyInvoice)
			Invoice.POST("/pay", handlers.PayInvoice)
			Invoice.POST("/payUserBtc", handlers.PayUserBtc)
			Invoice.POST("/querypayment", handlers.QueryPayment)
			Invoice.POST("/queryinvoice", handlers.QueryInvoice)
			Invoice.POST("/decodeinvoice", handlers.DecodeInvoice)
			//deprecated
			Invoice.POST("/querybalance", handlers.QueryBalance)
		}
		Asset := custody.Group("/Asset")
		{
			Asset.POST("/apply", handlers.ApplyAddress)
			Asset.POST("/send", handlers.SendAsset)
			Asset.POST("/sendToUserAsset", handlers.SendToUserAsset)
			Asset.POST("/queryassets", handlers.QueryAssets)
			Asset.POST("/querypayment", handlers.QueryAssetPayment)
			Asset.POST("/queryaddress", handlers.QueryAddress)
			Asset.POST("/decodeaddr", handlers.DecodeAddress)
			//back
			Asset.POST("/querypayments", handlers.QueryAssetPayments)
			Asset.POST("/queryaddresses", handlers.QueryAddresses)
			//deprecated
			Asset.POST("/queryasset", handlers.QueryAsset)

		}
		locked := custody.Group("/locked")
		{
			locked.POST("/querypayments", handlers.QueryLockedPayments)
		}
	}
	return router
}
