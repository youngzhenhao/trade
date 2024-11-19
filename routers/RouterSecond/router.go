package RouterSecond

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
	"trade/handlers/SecondHander"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	//award
	{
		award := r.Group("/award")
		award.POST("/PutInSatoshiAward", SecondHander.PutInSatoshiAward)
		award.POST("/PutAssetAward", SecondHander.PutAssetAward)
	}
	//fair_launch
	{
		username := base64.StdEncoding.EncodeToString([]byte(config.GetLoadConfig().AdminUser.Username))
		password := base64.StdEncoding.EncodeToString([]byte(config.GetLoadConfig().AdminUser.Password))
		authorized := r.Group("/fair_launch/auth_op", gin.BasicAuth(gin.Accounts{
			username: password,
		}))
		authorized.POST("/refund", handlers.RefundUserFirstMintByUsernameAndAssetId)
	}
	//lockAccount
	{
		locked := r.Group("/lockAccount")
		locked.POST("/getBalance", SecondHander.GetBalance)
		locked.POST("/lock", SecondHander.Lock)
		locked.POST("/unlock", SecondHander.Unlock)
		locked.POST("/payAsset", SecondHander.PayAsset)
	}
	//Query
	{
		query := r.Group("/LocalQuery")
		query.POST("/QueryBills", SecondHander.QueryBills)
		query.POST("/QueryBalances", SecondHander.QueryBalance)
		query.POST("/QueryBalanceList", SecondHander.GetBalanceList)
		query.POST("/TotalBillList", SecondHander.TotalBillList)
		user := query.Group("/user")
		{
			user.POST("/userinfo", SecondHander.QueryUserInfo)
			user.POST("/block", SecondHander.BlockUser)
			user.POST("/unblock", SecondHander.UnBlockUser)
		}
		locked := query.Group("/locked")
		{
			locked.POST("/QueryLockedBills", SecondHander.QueryLockedBills)
		}
	}
	//AssetList
	assetList := r.Group("/asset_list")
	assetList.GET("/is_exist", handlers.IsAssetListRecordExist)

	// userStats
	userStats := r.Group("/user_stats")
	userStats.GET("/count", handlers.GetDateLoginCount)
	userStats.GET("/record", handlers.GetDateIpLoginRecord)
	userStats.GET("/new_count", handlers.GetNewUserCount)

	// backReward
	backReward := r.Group("/back_reward")
	backReward.GET("/get", handlers.GetBackRewards)

	return r
}
