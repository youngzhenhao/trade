package services

import (
	"trade/models"
	"trade/services/btldb"
)

func CreateFairLaunchIncome(fairLaunchIncome *models.FairLaunchIncome) error {
	return btldb.CreateFairLaunchIncome(fairLaunchIncome)
}

// @dev: Create FairLaunch Income Records By these functions

func CreateFairLaunchIncomeOfUserPayIssuanceFee(fairLaunchInfoId int, feePaidId int, satAmount int, userId int, username string) error {
	return CreateFairLaunchIncome(&models.FairLaunchIncome{
		AssetId:                "",
		FairLaunchInfoId:       fairLaunchInfoId,
		FairLaunchMintedInfoId: 0,
		FeePaidId:              feePaidId,
		IncomeType:             models.UserPayIssuanceFee,
		IsIncome:               true,
		SatAmount:              satAmount,
		Txid:                   "",
		Addrs:                  "",
		UserId:                 userId,
		Username:               username,
	})
}

func CreateFairLaunchIncomeOfServerPayIssuanceFinalizeFee(fairLaunchInfoId int, txid string) error {
	return CreateFairLaunchIncome(&models.FairLaunchIncome{
		AssetId:                "",
		FairLaunchInfoId:       fairLaunchInfoId,
		FairLaunchMintedInfoId: 0,
		FeePaidId:              0,
		IncomeType:             models.ServerPayIssuanceFinalizeFee,
		IsIncome:               false,
		SatAmount:              0,
		Txid:                   txid,
		Addrs:                  "",
		UserId:                 0,
		Username:               "",
	})
}

func CreateFairLaunchIncomeOfServerPaySendReservedFee(assetId string, fairLaunchInfoId int, txid string) error {
	return CreateFairLaunchIncome(&models.FairLaunchIncome{
		AssetId:                assetId,
		FairLaunchInfoId:       fairLaunchInfoId,
		FairLaunchMintedInfoId: 0,
		FeePaidId:              0,
		IncomeType:             models.ServerPaySendReservedFee,
		IsIncome:               false,
		SatAmount:              0,
		Txid:                   txid,
		Addrs:                  "",
		UserId:                 0,
		Username:               "",
	})
}

func CreateFairLaunchIncomeOfUserPayMintedFee(assetId string, fairLaunchInfoId int, fairLaunchMintedInfoId int, feePaidId int, satAmount int, userId int, username string) error {
	return CreateFairLaunchIncome(&models.FairLaunchIncome{
		AssetId:                assetId,
		FairLaunchInfoId:       fairLaunchInfoId,
		FairLaunchMintedInfoId: fairLaunchMintedInfoId,
		FeePaidId:              feePaidId,
		IncomeType:             models.UserPayMintedFee,
		IsIncome:               true,
		SatAmount:              satAmount,
		Txid:                   "",
		Addrs:                  "",
		UserId:                 userId,
		Username:               username,
	})
}

func CreateFairLaunchIncomeOfServerPaySendAssetFee(assetId string, fairLaunchInfoId int, txid string, addrs string) error {
	return CreateFairLaunchIncome(&models.FairLaunchIncome{
		AssetId:                assetId,
		FairLaunchInfoId:       fairLaunchInfoId,
		FairLaunchMintedInfoId: 0,
		FeePaidId:              0,
		IncomeType:             models.ServerPaySendAssetFee,
		IsIncome:               false,
		SatAmount:              0,
		Txid:                   txid,
		Addrs:                  addrs,
		UserId:                 0,
		Username:               "",
	})
}

// TODO: Update SatAmount
func UpdateFairLaunchIncomesSatAmountByTxids() error {
	// TODO: Need to complete
	return nil
}

// TODO: Query total incomes and spent by fair launch id
