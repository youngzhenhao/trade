package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateAssetMeta(assetMeta *models.AssetMeta) error {
	return middleware.DB.Create(assetMeta).Error
}

func CreateAssetMetas(assetMetas *[]models.AssetMeta) error {
	return middleware.DB.Create(assetMetas).Error
}

func ReadAssetMeta(id uint) (*models.AssetMeta, error) {
	var assetMeta models.AssetMeta
	err := middleware.DB.First(&assetMeta, id).Error
	return &assetMeta, err
}

func ReadAssetMetaByAssetId(assetId string) (*models.AssetMeta, error) {
	var assetMeta models.AssetMeta
	err := middleware.DB.Where("asset_id = ?", assetId).First(&assetMeta).Error
	return &assetMeta, err
}

func ReadAllAssetMetas() (*[]models.AssetMeta, error) {
	var assetMetas []models.AssetMeta
	err := middleware.DB.Find(&assetMetas).Error
	return &assetMetas, err
}

func UpdateAssetMeta(assetMeta *models.AssetMeta) error {
	return middleware.DB.Save(assetMeta).Error
}

func UpdateAssetMetas(assetMetas *[]models.AssetMeta) error {
	return middleware.DB.Save(assetMetas).Error
}

func DeleteAssetMeta(id uint) error {
	var assetMeta models.AssetMeta
	return middleware.DB.Delete(&assetMeta, id).Error
}
