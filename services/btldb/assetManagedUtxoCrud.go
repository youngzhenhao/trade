package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateAssetManagedUtxo(assetManagedUtxo *models.AssetManagedUtxo) error {
	return middleware.DB.Create(assetManagedUtxo).Error
}

func CreateAssetManagedUtxos(assetManagedUtxos *[]models.AssetManagedUtxo) error {
	return middleware.DB.Create(assetManagedUtxos).Error
}

func ReadAllAssetManagedUtxos() (*[]models.AssetManagedUtxo, error) {
	var assetManagedUtxos []models.AssetManagedUtxo
	err := middleware.DB.Find(&assetManagedUtxos).Error
	return &assetManagedUtxos, err
}

func ReadAllAssetManagedUtxosUpdatedAtDesc() (*[]models.AssetManagedUtxo, error) {
	var assetManagedUtxos []models.AssetManagedUtxo
	err := middleware.DB.Order("updated_at desc").Find(&assetManagedUtxos).Error
	return &assetManagedUtxos, err
}

func ReadAssetManagedUtxo(id uint) (*models.AssetManagedUtxo, error) {
	var assetManagedUtxo models.AssetManagedUtxo
	err := middleware.DB.First(&assetManagedUtxo, id).Error
	return &assetManagedUtxo, err
}

func ReadAssetManagedUtxosByIds(assetManagedUtxoIds *[]int) (*[]models.AssetManagedUtxo, error) {
	var assetManagedUtxos []models.AssetManagedUtxo
	err := middleware.DB.Where(assetManagedUtxoIds).Find(&assetManagedUtxos).Error
	return &assetManagedUtxos, err
}

func ReadAssetManagedUtxosByUserId(userId int) (*[]models.AssetManagedUtxo, error) {
	var assetManagedUtxos []models.AssetManagedUtxo
	err := middleware.DB.Where("user_id = ? AND status = ?", userId, 1).Find(&assetManagedUtxos).Error
	return &assetManagedUtxos, err
}

func ReadAssetManagedUtxosByAssetId(assetId string) (*[]models.AssetManagedUtxo, error) {
	var assetManagedUtxos []models.AssetManagedUtxo
	err := middleware.DB.Where("asset_id = ? AND status = ?", assetId, 1).Find(&assetManagedUtxos).Error
	return &assetManagedUtxos, err
}

func ReadAssetManagedUtxoByUserIdAndAssetId(userId int, assetId string) (*models.AssetManagedUtxo, error) {
	var assetManagedUtxo models.AssetManagedUtxo
	err := middleware.DB.Where("user_id = ? AND asset_id = ? AND status = ?", userId, assetId, 1).First(&assetManagedUtxo).Error
	return &assetManagedUtxo, err
}

func UpdateAssetManagedUtxo(assetManagedUtxo *models.AssetManagedUtxo) error {
	return middleware.DB.Save(assetManagedUtxo).Error
}

func UpdateAssetManagedUtxos(assetManagedUtxos *[]models.AssetManagedUtxo) error {
	return middleware.DB.Save(assetManagedUtxos).Error
}

func DeleteAssetManagedUtxo(id uint) error {
	var assetManagedUtxo models.AssetManagedUtxo
	return middleware.DB.Delete(&assetManagedUtxo, id).Error
}

func DeleteAssetManagedUtxoByIds(assetManagedUtxoIds *[]int) error {
	var assetManagedUtxos []models.AssetManagedUtxo
	return middleware.DB.Delete(&assetManagedUtxos, &assetManagedUtxoIds).Error
}
