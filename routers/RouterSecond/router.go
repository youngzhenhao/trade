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

		//todo: add more api
		locked.POST("/getLockedBalanceList", SecondHander.GetLockedBalanceList)
	}
	//Query
	{
		query := r.Group("/LocalQuery")
		query.POST("/QueryBills", SecondHander.QueryBills)
		query.POST("/QueryBalances", SecondHander.QueryBalance)
		query.POST("/QueryBalanceList", SecondHander.GetBalanceList)
	}
	//AssetList
	assetList := r.Group("/asset_list")
	assetList.GET("/is_exist", handlers.IsAssetListRecordExist)
	return r
}
