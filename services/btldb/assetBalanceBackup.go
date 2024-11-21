package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateAssetBalanceBackup(assetBalanceBackup *models.AssetBalanceBackup) error {
	return middleware.DB.Create(assetBalanceBackup).Error
}

func ReadAssetBalanceBackupByUsername(username string) (*models.AssetBalanceBackup, error) {
	var assetBalanceBackup models.AssetBalanceBackup
	err := middleware.DB.Where("username = ?", username).First(&assetBalanceBackup).Error
	return &assetBalanceBackup, err
}

func UpdateAssetBalanceBackup(assetBalanceBackup *models.AssetBalanceBackup) error {
	return middleware.DB.Save(assetBalanceBackup).Error
}
