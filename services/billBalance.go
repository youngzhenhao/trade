package services

import (
	"trade/middleware"
	"trade/models"
)

func ReadBillBalanceByBalanceType(balanceType models.BalanceType) (*[]models.Balance, error) {
	var billBalances []models.Balance
	err := middleware.DB.Where("amount <> ? AND bill_type = ?", 0, balanceType).Find(&billBalances).Error
	return &billBalances, err
}

func ReadBillBalanceByBalanceTypeAndAssetId(balanceType models.BalanceType, assetId string) (*[]models.Balance, error) {
	var billBalances []models.Balance
	err := middleware.DB.Where("amount <> ? AND bill_type = ? AND asset_id = ?", 0, balanceType, assetId).Find(&billBalances).Error
	return &billBalances, err
}

func GetBillBalanceAssetTransfer() (*[]models.Balance, error) {
	return ReadBillBalanceByBalanceType(models.BillTypeAssetTransfer)
}

func GetBillBalanceAwardAsset() (*[]models.Balance, error) {
	return ReadBillBalanceByBalanceType(models.BillTypeAwardAsset)
}

func ReadBillBalanceAssetTransferAndAwardAssetByAssetId(assetId string) (*[]models.Balance, error) {
	var billBalances []models.Balance
	err := middleware.DB.Where("amount <> ? AND bill_type IN ? AND asset_id = ?", 0, []models.BalanceType{models.BillTypeAssetTransfer, models.BillTypeAwardAsset}, assetId).Order("updated_at desc").Find(&billBalances).Error
	return &billBalances, err
}
