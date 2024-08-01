package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateAssetLocalMintHistory(assetLocalMintHistory *models.AssetLocalMintHistory) error {
	return middleware.DB.Create(assetLocalMintHistory).Error
}

func CreateAssetLocalMintHistories(assetLocalMintHistories *[]models.AssetLocalMintHistory) error {
	return middleware.DB.Create(assetLocalMintHistories).Error
}

func ReadAllAssetLocalMintHistories() (*[]models.AssetLocalMintHistory, error) {
	var assetLocalMintHistories []models.AssetLocalMintHistory
	err := middleware.DB.Find(&assetLocalMintHistories).Error
	return &assetLocalMintHistories, err
}

func ReadAllAssetLocalMintHistoriesUpdatedAtDesc() (*[]models.AssetLocalMintHistory, error) {
	var assetLocalMintHistories []models.AssetLocalMintHistory
	err := middleware.DB.Order("updated_at desc").Find(&assetLocalMintHistories).Error
	return &assetLocalMintHistories, err
}

func ReadAssetLocalMintHistory(id uint) (*models.AssetLocalMintHistory, error) {
	var assetLocalMintHistory models.AssetLocalMintHistory
	err := middleware.DB.First(&assetLocalMintHistory, id).Error
	return &assetLocalMintHistory, err
}

func ReadAssetLocalMintHistoriesByUserId(userId int) (*[]models.AssetLocalMintHistory, error) {
	var assetLocalMintHistories []models.AssetLocalMintHistory
	err := middleware.DB.Where("user_id = ? AND status = ?", userId, 1).Find(&assetLocalMintHistories).Error
	return &assetLocalMintHistories, err
}

func ReadAssetLocalMintHistoryByAssetId(assetId string) (*models.AssetLocalMintHistory, error) {
	var assetLocalMintHistory models.AssetLocalMintHistory
	err := middleware.DB.Where("asset_id = ? AND status = ?", assetId, 1).First(&assetLocalMintHistory).Error
	return &assetLocalMintHistory, err
}

func UpdateAssetLocalMintHistory(assetLocalMintHistory *models.AssetLocalMintHistory) error {
	return middleware.DB.Save(assetLocalMintHistory).Error
}

func UpdateAssetLocalMintHistories(assetLocalMintHistories *[]models.AssetLocalMintHistory) error {
	return middleware.DB.Save(assetLocalMintHistories).Error
}

func DeleteAssetLocalMintHistory(id uint) error {
	var assetLocalMintHistory models.AssetLocalMintHistory
	return middleware.DB.Delete(&assetLocalMintHistory, id).Error
}
