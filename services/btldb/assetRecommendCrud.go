package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateAssetRecommend(assetRecommend *models.AssetRecommend) error {
	return middleware.DB.Create(assetRecommend).Error
}

func CreateAssetRecommends(assetRecommends *[]models.AssetRecommend) error {
	return middleware.DB.Create(assetRecommends).Error
}

func ReadAllAssetRecommends() (*[]models.AssetRecommend, error) {
	var assetRecommends []models.AssetRecommend
	err := middleware.DB.Find(&assetRecommends).Error
	return &assetRecommends, err
}

func ReadAllAssetRecommendsUpdatedAtDesc() (*[]models.AssetRecommend, error) {
	var assetRecommends []models.AssetRecommend
	err := middleware.DB.Order("updated_at desc").Find(&assetRecommends).Error
	return &assetRecommends, err
}

func ReadAssetRecommend(id uint) (*models.AssetRecommend, error) {
	var assetRecommend models.AssetRecommend
	err := middleware.DB.First(&assetRecommend, id).Error
	return &assetRecommend, err
}

func ReadAssetRecommendsByUserId(userId int) (*[]models.AssetRecommend, error) {
	var assetRecommends []models.AssetRecommend
	err := middleware.DB.Where("user_id = ? AND status = ?", userId, 1).Find(&assetRecommends).Error
	return &assetRecommends, err
}

func ReadAssetRecommendsByAssetId(assetId string) (*[]models.AssetRecommend, error) {
	var assetRecommends []models.AssetRecommend
	err := middleware.DB.Where("asset_id = ? AND status = ?", assetId, 1).Find(&assetRecommends).Error
	return &assetRecommends, err
}

func ReadAssetRecommendByUserIdAndAssetId(userId int, assetId string) (*models.AssetRecommend, error) {
	var assetRecommend models.AssetRecommend
	err := middleware.DB.Where("user_id = ? AND AND asset_id = ? status = ?", userId, assetId, 1).First(&assetRecommend).Error
	return &assetRecommend, err
}

func UpdateAssetRecommend(assetRecommend *models.AssetRecommend) error {
	return middleware.DB.Save(assetRecommend).Error
}

func UpdateAssetRecommends(assetRecommends *[]models.AssetRecommend) error {
	return middleware.DB.Save(assetRecommends).Error
}

func DeleteAssetRecommend(id uint) error {
	var assetRecommend models.AssetRecommend
	return middleware.DB.Delete(&assetRecommend, id).Error
}
