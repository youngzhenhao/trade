package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
	"trade/middleware"
)

func SetupFairLaunchFollowRouter(router *gin.Engine) *gin.Engine {
	fairLaunchFollow := router.Group("/fair_launch_follow")
	fairLaunchFollow.Use(middleware.AuthMiddleware())
	{
		fairLaunchFollow.POST("/follow", handlers.SetFollowFairLaunchInfo)
		fairLaunchFollow.POST("/unfollow/asset_id/:asset_id", handlers.SetUnfollowFairLaunchInfo)
		query := fairLaunchFollow.Group("/query")
		query.GET("/user/followed", handlers.GetFollowedFairLaunchInfo)
	}
	authorized := router.Group("/fair_launch_follow", gin.BasicAuth(gin.Accounts{
		config.GetLoadConfig().AdminUser.Username: config.GetLoadConfig().AdminUser.Password,
	}))
	authorized.GET("/get/all/simplified", handlers.GetAllFairLaunchFollowSimplified)
	return router
}
