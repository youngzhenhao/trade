package services

import (
	"trade/middleware"
	"trade/models"
)

func CreateAssetBalance(assetBalance *models.AssetBalance) error {
	return middleware.DB.Create(assetBalance).Error
}

func CreateAssetBalances(assetBalances *[]models.AssetBalance) error {
	return middleware.DB.Create(assetBalances).Error
}

func ReadAllAssetBalances() (*[]models.AssetBalance, error) {
	var assetBalances []models.AssetBalance
	err := middleware.DB.Order("updated_at desc").Find(&assetBalances).Error
	return &assetBalances, err
}

func ReadAssetBalance(id uint) (*models.AssetBalance, error) {
	var assetBalance models.AssetBalance
	err := middleware.DB.First(&assetBalance, id).Error
	return &assetBalance, err
}

func ReadAssetBalancesByUserId(userId int) (*[]models.AssetBalance, error) {
	var assetBalances []models.AssetBalance
	err := middleware.DB.Where("user_id = ? AND status = ?", userId, 1).Find(&assetBalances).Error
	return &assetBalances, err
}

func ReadAssetBalanceByAssetID(assetId string) (*models.AssetBalance, error) {
	var assetBalance models.AssetBalance
	err := middleware.DB.Where("asset_id = ? AND status = ?", assetId, 1).First(&assetBalance).Error
	return &assetBalance, err
}

func UpdateAssetBalance(assetBalance *models.AssetBalance) error {
	return middleware.DB.Save(assetBalance).Error
}

func UpdateAssetBalances(assetBalances *[]models.AssetBalance) error {
	return middleware.DB.Save(assetBalances).Error
}

func DeleteAssetBalance(id uint) error {
	var assetBalance models.AssetBalance
	return middleware.DB.Delete(&assetBalance, id).Error
}
