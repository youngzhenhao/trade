package btldb

import (
	"trade/middleware"
	"trade/models"
)

// CreateAssetSyncInfo creates a new AssetSyncInfo record
func CreateAssetSyncInfo(assetSyncInfo *models.AssetSyncInfo) error {
	return middleware.DB.Create(assetSyncInfo).Error
}

// ReadAssetSyncInfo retrieves a AssetSyncInfo by Id
func ReadAssetSyncInfo(id uint) (*models.AssetSyncInfo, error) {
	var assetSyncInfo models.AssetSyncInfo
	err := middleware.DB.First(&assetSyncInfo, id).Error
	return &assetSyncInfo, err
}

func ReadAssetSyncInfoByAssetID(assetID string) (*models.AssetSyncInfo, error) {
	var assetSyncInfo models.AssetSyncInfo
	err := middleware.DB.Where("asset_id =?", assetID).First(&assetSyncInfo).Error
	return &assetSyncInfo, err
}

// UpdateAssetSyncInfo updates an existing AssetSyncInfo
func UpdateAssetSyncInfo(assetSyncInfo *models.AssetSyncInfo) error {
	return middleware.DB.Save(assetSyncInfo).Error
}

// DeleteAssetSyncInfo soft deletes a AssetSyncInfo by Id
func DeleteAssetSyncInfo(id uint) error {
	var assetSyncInfo models.AssetSyncInfo
	return middleware.DB.Delete(&assetSyncInfo, id).Error
}
