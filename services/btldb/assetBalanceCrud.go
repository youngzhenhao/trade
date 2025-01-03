package btldb

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

func ReadAllAssetBalancesNonZeroUpdatedAtDesc() (*[]models.AssetBalance, error) {
	var assetBalances []models.AssetBalance
	err := middleware.DB.Where("balance <> ?", 0).Order("updated_at desc").Find(&assetBalances).Error
	return &assetBalances, err
}

func ReadAllAssetBalancesNonZero() (*[]models.AssetBalance, error) {
	var assetBalances []models.AssetBalance
	err := middleware.DB.Where("balance <> ?", 0).Order("balance desc").Find(&assetBalances).Error
	return &assetBalances, err
}

func ReadAllAssetBalancesNonZeroByAssetId(assetId string) (*[]models.AssetBalance, error) {
	var assetBalances []models.AssetBalance
	err := middleware.DB.Where("asset_id = ? AND balance <> ?", assetId, 0).Order("balance desc").Find(&assetBalances).Error
	return &assetBalances, err
}

func ReadAllAssetBalancesNonZeroLimit(limit int) (*[]models.AssetBalance, error) {
	var assetBalances []models.AssetBalance
	err := middleware.DB.Where("balance <> ?", 0).Limit(limit).Order("balance desc").Find(&assetBalances).Error
	return &assetBalances, err
}

func ReadAllAssetBalancesNonZeroLimitAndOffset(limit int, offset int) (*[]models.AssetBalance, error) {
	var assetBalances []models.AssetBalance
	err := middleware.DB.Where("balance <> ?", 0).Order("updated_at desc").Limit(limit).Offset(offset).Find(&assetBalances).Error
	return &assetBalances, err
}

func ReadAssetBalance(id uint) (*models.AssetBalance, error) {
	var assetBalance models.AssetBalance
	err := middleware.DB.First(&assetBalance, id).Error
	return &assetBalance, err
}

func ReadAssetBalancesByUserId(userId int) (*[]models.AssetBalance, error) {
	var assetBalances []models.AssetBalance
	err := middleware.DB.Where("user_id = ?", userId).Find(&assetBalances).Error
	return &assetBalances, err
}

func ReadAssetBalancesByUserIdNonZero(userId int) (*[]models.AssetBalance, error) {
	var assetBalances []models.AssetBalance
	err := middleware.DB.Where("user_id = ? AND balance <> ?", userId, 0).Find(&assetBalances).Error
	return &assetBalances, err
}

func ReadAssetBalanceByAssetId(assetId string) (*[]models.AssetBalance, error) {
	var assetBalances []models.AssetBalance
	err := middleware.DB.Where("asset_id = ?", assetId).Find(&assetBalances).Error
	return &assetBalances, err
}

func ReadAssetBalanceByAssetIdNonZero(assetId string) (*[]models.AssetBalance, error) {
	var assetBalances []models.AssetBalance
	err := middleware.DB.Where("asset_id = ? AND balance <> ?", assetId, 0).Order("updated_at desc").Find(&assetBalances).Error
	return &assetBalances, err
}

// ReadAssetBalanceByAssetIdNonZeroLimitAndOffset
// @Description: read asset balance by asset id non-zero limit and offset
func ReadAssetBalanceByAssetIdNonZeroLimitAndOffset(assetId string, limit int, offset int) (*[]models.AssetBalance, error) {
	var assetBalances []models.AssetBalance
	err := middleware.DB.Where("asset_id = ? AND balance <> ?", assetId, 0).Order("balance desc").Limit(limit).Offset(offset).Find(&assetBalances).Error
	return &assetBalances, err
}

func ReadAssetBalanceByAssetIdAndUserId(assetId string, userId int) (*models.AssetBalance, error) {
	var assetBalance models.AssetBalance
	err := middleware.DB.Where("asset_id = ? AND user_id = ?", assetId, userId).First(&assetBalance).Error
	return &assetBalance, err
}

func ReadAssetBalanceByUsername(username string) (*[]models.AssetBalance, error) {
	var assetBalances []models.AssetBalance
	err := middleware.DB.Where("username = ?", username).Find(&assetBalances).Error
	return &assetBalances, err
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
