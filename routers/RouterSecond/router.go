package RouterSecond

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
	"trade/handlers/SecondHandler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	//award
	{
		award := r.Group("/award")
		award.POST("/PutInSatoshiAward", SecondHandler.PutInSatoshiAward)
		award.POST("/PutAssetAward", SecondHandler.PutAssetAward)
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
		locked.POST("/CheckUserStatus", SecondHandler.CheckUserStatus)
		locked.POST("/getBalance", SecondHandler.GetBalance)
		locked.POST("/lock", SecondHandler.Lock)
		locked.POST("/unlock", SecondHandler.Unlock)
		locked.POST("/payAsset", SecondHandler.PayAsset)
	}
	//Query
	{
		query := r.Group("/LocalQuery")
		{
			query.POST("/QueryBills", SecondHandler.QueryBills)
			query.POST("/QueryBalances", SecondHandler.QueryBalance)
			query.POST("/QueryBalancesChange", SecondHandler.QueryBalancesChange)
			query.POST("/QueryBalanceList", SecondHandler.GetBalanceList)
			query.POST("/TotalBillList", SecondHandler.TotalBillList)

			custodyInfo := query.Group("/custody_info")
			{
				custodyInfo.POST("/QueryChannelAssetInfo", SecondHandler.QueryChannelAssetInfo)
			}
		}

		limit := query.Group("/limit")
		{
			limit.POST("/GetUserLimit", SecondHandler.GetUserLimitHandler)
			limit.POST("/SetUserLimitLevel", SecondHandler.SetUserLimitLevelHandler)
			limit.POST("/SetUserTodayLimit", SecondHandler.SetUserTodayLimitHandler)

			limit.POST("/GetLimitType", SecondHandler.GetLimitTypeHandler)
			limit.POST("/CreateOrUpdateLimitType", SecondHandler.CreateOrUpdateLimitTypeHandle)
			limit.POST("/GetLimitTypeLevels", SecondHandler.GetLimitTypeLevelsHandle)
			limit.POST("/CreateOrUpdateLimitTypeLevel", SecondHandler.CreateOrUpdateLimitTypeLevelHandle)
		}

		user := query.Group("/user")
		{
			user.POST("/userinfo", SecondHandler.QueryUserInfo)
			user.POST("/block", SecondHandler.BlockUser)
			user.POST("/unblock", SecondHandler.UnBlockUser)
		}

		locked := query.Group("/locked")
		{
			locked.POST("/QueryLockedBills", SecondHandler.QueryLockedBills)
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

	// assetMeta
	assetMeta := r.Group("/asset_meta")
	assetMeta.POST("/image/query", handlers.GetAssetMetaImage)

	// nftPresale
	nftPresale := r.Group("/nft_presale")
	nftPresale.GET("/get/purchased", handlers.GetPurchasedNftPresaleInfo)
	nftPresale.GET("/get/offline/purchased", handlers.GetNftPresaleOfflinePurchaseData)
	nftPresale.POST("/get/offline/update", handlers.UpdateNftPresaleOfflinePurchaseData)

	// btcBalance
	btcBalance := r.Group("/btc_balance")
	btcBalance.GET("/get/rank/count", handlers.GetBtcBalanceCount)
	btcBalance.GET("/get/rank", handlers.GetBtcBalanceOrderLimitOffset)

	return r
}
