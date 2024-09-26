package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateAssetGroup(assetGroup *models.AssetGroup) error {
	return middleware.DB.Create(assetGroup).Error
}

func CreateAssetGroups(assetGroups *[]models.AssetGroup) error {
	return middleware.DB.Create(assetGroups).Error
}

func ReadAllAssetGroups() (*[]models.AssetGroup, error) {
	var assetGroups []models.AssetGroup
	err := middleware.DB.Order("updated_at desc").Find(&assetGroups).Error
	return &assetGroups, err
}

func ReadAssetGroup(id uint) (*models.AssetGroup, error) {
	var assetGroup models.AssetGroup
	err := middleware.DB.First(&assetGroup, id).Error
	return &assetGroup, err
}

func ReadAssetGroupByTweakedGroupKey(tweakedGroupKey string) (*models.AssetGroup, error) {
	var assetGroup models.AssetGroup
	err := middleware.DB.Where("tweaked_group_key = ?", tweakedGroupKey).First(&assetGroup).Error
	return &assetGroup, err
}

func UpdateAssetGroup(assetGroup *models.AssetGroup) error {
	return middleware.DB.Save(assetGroup).Error
}

func UpdateAssetGroups(assetGroups *[]models.AssetGroup) error {
	return middleware.DB.Save(assetGroups).Error
}

func DeleteAssetGroup(id uint) error {
	var assetGroup models.AssetGroup
	return middleware.DB.Delete(&assetGroup, id).Error
}
