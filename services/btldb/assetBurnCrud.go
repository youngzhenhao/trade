package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateAssetBurn(assetBurn *models.AssetBurn) error {
	return middleware.DB.Create(assetBurn).Error
}

func CreateAssetBurns(assetBurns *[]models.AssetBurn) error {
	return middleware.DB.Create(assetBurns).Error
}

func ReadAllAssetBurns() (*[]models.AssetBurn, error) {
	var assetBurns []models.AssetBurn
	err := middleware.DB.Find(&assetBurns).Error
	return &assetBurns, err
}

func ReadAllAssetBurnsUpdatedAt() (*[]models.AssetBurn, error) {
	var assetBurns []models.AssetBurn
	err := middleware.DB.Order("updated_at desc").Find(&assetBurns).Error
	return &assetBurns, err
}

func ReadAssetBurn(id uint) (*models.AssetBurn, error) {
	var assetBurn models.AssetBurn
	err := middleware.DB.First(&assetBurn, id).Error
	return &assetBurn, err
}

func ReadAssetBurnsByUserId(userId int) (*[]models.AssetBurn, error) {
	var assetBurns []models.AssetBurn
	err := middleware.DB.Where("user_id = ?", userId).Find(&assetBurns).Error
	return &assetBurns, err
}

func ReadAssetBurnsByAssetId(assetId string) (*[]models.AssetBurn, error) {
	var assetBurns []models.AssetBurn
	err := middleware.DB.Where("asset_id = ?", assetId).Find(&assetBurns).Error
	return &assetBurns, err
}

func UpdateAssetBurn(assetBurn *models.AssetBurn) error {
	return middleware.DB.Save(assetBurn).Error
}

func UpdateAssetBurns(assetBurns *[]models.AssetBurn) error {
	return middleware.DB.Save(assetBurns).Error
}

func DeleteAssetBurn(id uint) error {
	var assetBurn models.AssetBurn
	return middleware.DB.Delete(&assetBurn, id).Error
}
