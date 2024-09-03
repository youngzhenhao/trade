package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers/test"
	"trade/middleware"
)

func SetupCustodyAccountRouter(router *gin.Engine) *gin.Engine {

	// A routing group that requires authentication
	custody := router.Group("/custodyAccount")

	custody.Use(middleware.AuthMiddleware())
	{
		custody.POST("/create", test.CreateCustodyAccount)
		Invoice := custody.Group("/invoice")
		{
			Invoice.POST("/apply", test.ApplyInvoice)
			Invoice.POST("/pay", test.PayInvoice)
			Invoice.POST("/querybalance", test.QueryBalance)
			Invoice.POST("/queryinvoice", test.QueryInvoice)
			Invoice.POST("/querypayment", test.QueryPayment)
			Invoice.POST("/lookupinvoice", test.LookupInvoice)
		}
	}
	return router
}
