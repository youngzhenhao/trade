package dao

import (
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
	"trade/models/custodyModels/pAccount"
	"trade/services"
	"trade/services/pool"
	"trade/services/satBackQueue"
)

func Migrate() error {
	var err error
	{
		if err = custodyMigrate(err); err != nil {
			return err
		}
		if err = custodyAwardMigrate(err); err != nil {
			return err
		}
		if err = custodyLimitMigrate(err); err != nil {
			return err
		}
		if err = custodyBTCMigrate(err); err != nil {
			return err
		}
		if err = custodyPAccountMigrate(err); err != nil {
			return err
		}
	}

	if err = middleware.DB.AutoMigrate(&models.Account{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.Balance{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.BalanceExt{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.UserConfig{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.ScheduledTask{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.Invoice{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.FairLaunchInfo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.FairLaunchMintedInfo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.FairLaunchMintedUserInfo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.FeeRateInfo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetIssuance{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.IdoPublishInfo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.IdoParticipateInfo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.IdoParticipateUserInfo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetSyncInfo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.BtcBalance{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetTransferProcessedDb{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetTransferProcessedInputDb{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetTransferProcessedOutputDb{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AddrReceiveEvent{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.BatchTransfer{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetAddr{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetLock{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetBalance{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetBurn{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetLocalMint{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetRecommend{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.LoginRecord{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.FairLaunchFollow{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetLocalMintHistory{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetManagedUtxo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.FairLaunchMintedAndAvailableInfo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.FairLaunchIncome{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.BackFee{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&custodyModels.AccountBalance{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&custodyModels.PayOutside{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&custodyModels.PayOutsideTx{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.LogFileUpload{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AccountAssetReceive{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetGroup{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.NftTransfer{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.NftInfo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.NftPresale{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetMeta{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.NftPresaleBatchGroup{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.NftPresaleWhitelist{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetList{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.DateIpLogin{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.DateLogin{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.BalanceTypeExt{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetBalanceBackup{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetBalanceHistory{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&satBackQueue.PushQueueRecord{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.RestRecord{}); err != nil {
		return err
	}
	if err = poolMigrate(); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&services.NftPresaleOfflinePurchaseData{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.BtcUtxo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pool.PoolLpAwardCumulative{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pool.PoolShareLpAwardBalance{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pool.PoolShareLpAwardCumulative{}); err != nil {
		return err
	}

	return err
}

func poolMigrate() (err error) {
	if err = middleware.DB.AutoMigrate(&pool.PoolPair{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pool.PoolShare{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pool.PoolShareBalance{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pool.PoolShareRecord{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pool.PoolSwapRecord{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pool.PoolLpAwardBalance{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pool.PoolLpAwardRecord{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pool.PoolWithdrawAwardRecord{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pool.PoolAddLiquidityBatch{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pool.PoolRemoveLiquidityBatch{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pool.PoolSwapExactTokenForTokenNoPathBatch{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pool.PoolSwapTokenForExactTokenNoPathBatch{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pool.PoolWithdrawAwardBatch{}); err != nil {
		return err
	}
	return nil
}

func custodyMigrate(err error) error {
	if err = middleware.DB.AutoMigrate(&custodyModels.LockBill{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&custodyModels.LockAccount{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&custodyModels.LockBalance{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&custodyModels.LockBillExt{}); err != nil {
		return err
	}
	return err
}

func custodyAwardMigrate(err error) error {
	if err = middleware.DB.AutoMigrate(&models.AccountAwardExt{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AwardInventory{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AccountAward{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AccountAwardIdempotent{}); err != nil {
		return err
	}
	return err
}

func custodyLimitMigrate(err error) error {
	if err = middleware.DB.AutoMigrate(&custodyModels.Limit{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&custodyModels.LimitBill{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&custodyModels.LimitLevel{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&custodyModels.LimitType{}); err != nil {
		return err
	}
	//blocked
	if err = middleware.DB.AutoMigrate(&custodyModels.BlockedRecord{}); err != nil {
		return err
	}

	return err
}

func custodyBTCMigrate(err error) error {
	if err = middleware.DB.AutoMigrate(&custodyModels.AccountInsideMission{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&custodyModels.AccountOutsideMission{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&custodyModels.AccountBalanceChange{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&custodyModels.AccountBtcBalance{}); err != nil {
		return err
	}
	return err
}
func custodyPAccountMigrate(err error) error {
	if err = middleware.DB.AutoMigrate(&pAccount.PoolAccount{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pAccount.PAccountAssetId{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pAccount.PAccountBalance{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pAccount.PAccountBill{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&pAccount.PAccountBalanceChange{}); err != nil {
		return err
	}
	return err
}
