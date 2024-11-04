package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateAssetList(assetList *models.AssetList) error {
	return middleware.DB.Create(assetList).Error
}

func CreateAssetLists(assetLists *[]models.AssetList) error {
	return middleware.DB.Create(assetLists).Error
}

func ReadAllAssetLists() (*[]models.AssetList, error) {
	var assetLists []models.AssetList
	err := middleware.DB.Order("updated_at desc").Find(&assetLists).Error
	return &assetLists, err
}

func ReadAllAssetListsNonZeroUpdatedAtDesc() (*[]models.AssetList, error) {
	var assetLists []models.AssetList
	err := middleware.DB.Where("amount <> ?", 0).Order("updated_at desc").Find(&assetLists).Error
	return &assetLists, err
}

func ReadAllAssetListsNonZero() (*[]models.AssetList, error) {
	var assetLists []models.AssetList
	err := middleware.DB.Where("amount <> ?", 0).Order("amount desc").Find(&assetLists).Error
	return &assetLists, err
}

func ReadAllAssetListsNonZeroByAssetId(assetId string) (*[]models.AssetList, error) {
	var assetLists []models.AssetList
	err := middleware.DB.Where("asset_id = ? AND amount <> ?", assetId, 0).Order("amount desc").Find(&assetLists).Error
	return &assetLists, err
}

func ReadAllAssetListsNonZeroLimit(limit int) (*[]models.AssetList, error) {
	var assetLists []models.AssetList
	err := middleware.DB.Where("amount <> ?", 0).Limit(limit).Order("amount desc").Find(&assetLists).Error
	return &assetLists, err
}

func ReadAllAssetListsNonZeroLimitAndOffset(limit int, offset int) (*[]models.AssetList, error) {
	var assetLists []models.AssetList
	err := middleware.DB.Where("amount <> ?", 0).Order("updated_at desc").Limit(limit).Offset(offset).Find(&assetLists).Error
	return &assetLists, err
}

func ReadAssetList(id uint) (*models.AssetList, error) {
	var assetList models.AssetList
	err := middleware.DB.First(&assetList, id).Error
	return &assetList, err
}

func ReadAssetListsByUserId(userId int) (*[]models.AssetList, error) {
	var assetLists []models.AssetList
	err := middleware.DB.Where("user_id = ?", userId).Find(&assetLists).Error
	return &assetLists, err
}

func ReadAssetListsByUserIdNonZero(userId int) (*[]models.AssetList, error) {
	var assetLists []models.AssetList
	err := middleware.DB.Where("user_id = ? AND amount <> ?", userId, 0).Find(&assetLists).Error
	return &assetLists, err
}

func ReadAssetListByAssetId(assetId string) (*[]models.AssetList, error) {
	var assetLists []models.AssetList
	err := middleware.DB.Where("asset_id = ?", assetId).Find(&assetLists).Error
	return &assetLists, err
}

func ReadAssetListByAssetIdNonZero(assetId string) (*[]models.AssetList, error) {
	var assetLists []models.AssetList
	err := middleware.DB.Where("asset_id = ? AND amount <> ?", assetId, 0).Order("updated_at desc").Find(&assetLists).Error
	return &assetLists, err
}

func ReadAssetListByAssetIdNonZeroLimitAndOffset(assetId string, limit int, offset int) (*[]models.AssetList, error) {
	var assetLists []models.AssetList
	err := middleware.DB.Where("asset_id = ? AND amount <> ?", assetId, 0).Order("amount desc").Limit(limit).Offset(offset).Find(&assetLists).Error
	return &assetLists, err
}

func ReadAssetListByAssetIdAndUserId(assetId string, userId int) (*models.AssetList, error) {
	var assetList models.AssetList
	err := middleware.DB.Where("asset_id = ? AND user_id = ?", assetId, userId).First(&assetList).Error
	return &assetList, err
}

func ReadAssetListByUsername(username string) (*[]models.AssetList, error) {
	var assetLists []models.AssetList
	err := middleware.DB.Where("username = ?", username).Find(&assetLists).Error
	return &assetLists, err
}

func UpdateAssetList(assetList *models.AssetList) error {
	return middleware.DB.Save(assetList).Error
}

func UpdateAssetLists(assetLists *[]models.AssetList) error {
	return middleware.DB.Save(assetLists).Error
}

func DeleteAssetList(id uint) error {
	var assetList models.AssetList
	return middleware.DB.Delete(&assetList, id).Error
}
