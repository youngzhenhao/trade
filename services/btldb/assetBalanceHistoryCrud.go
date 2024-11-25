package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateAssetBalanceHistory(assetBalanceHistory *models.AssetBalanceHistory) error {
	return middleware.DB.Create(assetBalanceHistory).Error
}

func CreateAssetBalanceHistorys(assetBalanceHistorys *[]models.AssetBalanceHistory) error {
	return middleware.DB.Create(assetBalanceHistorys).Error
}

func ReadAssetBalanceHistory(assetId string, username string) (*models.AssetBalanceHistory, error) {
	var assetBalanceHistory models.AssetBalanceHistory
	err := middleware.DB.Where("asset_id = ? and username = ?", assetId, username).Order("id desc").First(&assetBalanceHistory).Error
	return &assetBalanceHistory, err
}

func UpdateAssetBalanceHistory(assetBalanceHistory *models.AssetBalanceHistory) error {
	return middleware.DB.Save(assetBalanceHistory).Error
}

func UpdateAssetBalanceHistorys(assetBalanceHistorys *[]models.AssetBalanceHistory) error {
	return middleware.DB.Save(assetBalanceHistorys).Error
}

func DeleteAssetBalanceHistory(id uint) error {
	var assetBalanceHistory models.AssetBalanceHistory
	return middleware.DB.Delete(&assetBalanceHistory, id).Error
}
