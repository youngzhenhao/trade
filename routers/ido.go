package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

//TODO: need to test

func SetupIdoRouter(router *gin.Engine) *gin.Engine {
	ido := router.Group("/ido")
	version := ido.Group("/v1")
	version.Use(middleware.AuthMiddleware())
	{
		query := version.Group("/query")
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
		version.POST("/publish", handlers.SetIdoPublishInfo)
		version.POST("/participate", handlers.SetIdoParticipateInfo)
	}
	return router
}
