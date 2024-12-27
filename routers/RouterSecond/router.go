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

	// pool
	pool := r.Group("/pool")
	query := pool.Group("/query")
	{
		query.GET("/pool_info", handlers.QueryPoolInfo)

		query.GET("/share_records/count", handlers.QueryShareRecordsCount)
		query.GET("/share_records", handlers.QueryShareRecords)

		query.GET("/all_share_records/count", handlers.QueryUserAllShareRecordsCount)
		query.GET("/all_share_records", handlers.QueryUserAllShareRecords)

		query.GET("share_balance", handlers.QueryUserShareBalance)

		query.GET("/swap_records/count", handlers.QuerySwapRecordsCount)
		query.GET("/swap_records", handlers.QuerySwapRecords)

		query.GET("/all_swap_records/count", handlers.QueryUserAllSwapRecordsCount)
		query.GET("/all_swap_records", handlers.QueryUserAllSwapRecords)

		query.GET("lp_award_balance", handlers.QueryUserLpAwardBalance)

		query.GET("/withdraw_award_records/count", handlers.QueryWithdrawAwardRecordsCount)
		query.GET("/withdraw_award_records", handlers.QueryWithdrawAwardRecords)

		query.GET("/liquidity_and_award_records/count", handlers.QueryLiquidityAndAwardRecordsCount)
		query.GET("/liquidity_and_award_records", handlers.QueryLiquidityAndAwardRecords)

		query.GET("/lp_award_records/count", handlers.QueryLpAwardRecordsCount)
		query.GET("/lp_award_records", handlers.QueryLpAwardRecords)
	}
	calc := pool.Group("/calc")
	{
		calc.GET("/quote", handlers.CalcQuote)
		calc.GET("/burn_liquidity", handlers.CalcBurnLiquidity)
		calc.POST("/add_liquidity", handlers.CalcAddLiquidity)
		calc.POST("/remove_liquidity", handlers.CalcRemoveLiquidity)
		calc.GET("/amount_out", handlers.CalcAmountOut)
		calc.GET("/amount_in", handlers.CalcAmountIn)
		calc.POST("/swap_exact_token_for_token_no_path", handlers.CalcSwapExactTokenForTokenNoPath)
		calc.POST("/swap_token_for_exact_token_no_path", handlers.CalcSwapTokenForExactTokenNoPath)
	}

	return r
}
