package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

//TODO: need to test

func SetupIdoRouter(router *gin.Engine) *gin.Engine {
	version := router.Group("/v1")
	ido := version.Group("/ido")
	ido.Use(middleware.AuthMiddleware())
	{
		query := ido.Group("/query")
		{
			query.GET("/all_publish", handlers.GetAllIdoPublishInfo)
			query.GET("/published", handlers.GetIdoPublishedInfo)
			query.GET("/own_publish", handlers.GetOwnIdoPublishInfo)
			query.GET("/own_participate", handlers.GetOwnIdoParticipateInfo)
			query.GET("/publish/id/:id", handlers.GetIdoPublishInfo)
			query.GET("/publish/asset_id/:asset_id", handlers.GetIdoPublishInfoByAssetId)
			query.GET("/participate/:id", handlers.GetIdoParticipateInfo)
			query.POST("/participate", handlers.QueryIdoParticipateIsAvailable)
		}
		ido.POST("/publish", handlers.SetIdoPublishInfo)
		ido.POST("/participate", handlers.SetIdoParticipateInfo)
	}
	return router
}
