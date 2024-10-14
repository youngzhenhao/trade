package dao

import (
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
)

func Migrate() error {
	var err error
	err = custodyMigrate(err)
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
	if err = middleware.DB.AutoMigrate(&models.PayInside{}); err != nil {
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
	if err = middleware.DB.AutoMigrate(&models.AccountBalance{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AccountAward{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.PayOutside{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.PayOutsideTx{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AwardInventory{}); err != nil {
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

	return err
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
