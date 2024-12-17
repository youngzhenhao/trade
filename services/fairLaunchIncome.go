package services

import (
	"gorm.io/gorm"
	"trade/api"
	"trade/models"
	"trade/services/btldb"
	"trade/utils"
)

func CreateFairLaunchIncome(tx *gorm.DB, fairLaunchIncome *models.FairLaunchIncome) error {
	return btldb.CreateFairLaunchIncome(tx, fairLaunchIncome)
}

// @dev: Create FairLaunch Income Records By these functions

func CreateFairLaunchIncomeOfUserPayIssuanceFee(tx *gorm.DB, fairLaunchInfoId int, feePaidId int, satAmount int, userId int, username string) error {
	return CreateFairLaunchIncome(tx, &models.FairLaunchIncome{
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

func CreateFairLaunchIncomeOfServerPayIssuanceFinalizeFee(tx *gorm.DB, fairLaunchInfoId int, txid string) error {
	return CreateFairLaunchIncome(tx, &models.FairLaunchIncome{
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

func CreateFairLaunchIncomeOfServerPaySendReservedFee(tx *gorm.DB, assetId string, fairLaunchInfoId int, txid string) error {
	return CreateFairLaunchIncome(tx, &models.FairLaunchIncome{
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

func CreateFairLaunchIncomeOfUserPayMintedFee(tx *gorm.DB, assetId string, fairLaunchInfoId int, fairLaunchMintedInfoId int, feePaidId int, satAmount int, userId int, username string) error {
	return CreateFairLaunchIncome(tx, &models.FairLaunchIncome{
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

func CreateFairLaunchIncomeOfServerPaySendAssetFee(tx *gorm.DB, assetId string, fairLaunchInfoId int, txid string, addrs string) error {
	return CreateFairLaunchIncome(tx, &models.FairLaunchIncome{
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

func GetFairLaunchIncomesWhoseTxidIsNotNullAndSatAmountIsZero() (*[]models.FairLaunchIncome, error) {
	return btldb.ReadFairLaunchIncomesWhoseTxidIsNotNullAndSatAmountIsZero()
}

func UpdateFairLaunchIncomes(fairLaunchIncomes *[]models.FairLaunchIncome) error {
	if fairLaunchIncomes == nil {
		return nil
	}
	return btldb.UpdateFairLaunchIncomes(fairLaunchIncomes)
}

// @dev: Scheduled task
func UpdateFairLaunchIncomesSatAmountByTxids(network models.Network) error {
	// @dev: Get fair launch incomes whose txid is not null
	fairLaunchIncomes, err := GetFairLaunchIncomesWhoseTxidIsNotNullAndSatAmountIsZero()
	if err != nil {
		return err
	}
	if fairLaunchIncomes == nil || len(*fairLaunchIncomes) == 0 {
		return nil
	}
	var txids []string
	txidMapId := make(map[string]int)
	txidMapFee := make(map[string]int)
	// @dev: Get txid slice and txid map income id
	for _, fairLaunchIncome := range *fairLaunchIncomes {
		txids = append(txids, fairLaunchIncome.Txid)
		txidMapId[fairLaunchIncome.Txid] = int(fairLaunchIncome.ID)
	}
	// @dev: Get transactions by txids
	rawTransactionResponse, err := api.GetRawTransactionsByTxids(network, txids)
	if err != nil {
		return err
	}
	// @dev: Get txid map fee
	for _, rawTransaction := range *rawTransactionResponse {
		if rawTransaction.Error != nil {
			continue
		}
		fee := rawTransaction.Result.Fee
		txidMapFee[rawTransaction.Result.Txid] = utils.ToSat(fee)
	}
	var fairLaunchIncomesUpdated []models.FairLaunchIncome
	// @dev: Set fee by txid
	for _, fairLaunchIncome := range *fairLaunchIncomes {
		fee, ok := txidMapFee[fairLaunchIncome.Txid]
		if ok {
			fairLaunchIncome.SatAmount = fee
			fairLaunchIncomesUpdated = append(fairLaunchIncomesUpdated, fairLaunchIncome)
		}
	}
	if len(fairLaunchIncomesUpdated) == 0 {
		return nil
	}
	// @dev: Update fair launch incomes
	return UpdateFairLaunchIncomes(&fairLaunchIncomesUpdated)
}

func GetFairLaunchIncomesByAssetId(assetId string) (*[]models.FairLaunchIncome, error) {
	return btldb.ReadFairLaunchIncomesByAssetId(assetId)
}

// TODO: Query total incomes and spent by fair launch id
