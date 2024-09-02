package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateAssetLock(assetLock *models.AssetLock) error {
	return middleware.DB.Create(assetLock).Error
}

func CreateAssetLocks(assetLocks *[]models.AssetLock) error {
	return middleware.DB.Create(assetLocks).Error
}

func ReadAllAssetLocks() (*[]models.AssetLock, error) {
	var assetLocks []models.AssetLock
	err := middleware.DB.Find(&assetLocks).Error
	return &assetLocks, err
}

func ReadAssetLock(id uint) (*models.AssetLock, error) {
	var assetLock models.AssetLock
	err := middleware.DB.First(&assetLock, id).Error
	return &assetLock, err
}

func ReadAssetLocksByUserId(userId int) (*[]models.AssetLock, error) {
	var assetLocks []models.AssetLock
	err := middleware.DB.Where("user_id = ?", userId).Find(&assetLocks).Error
	return &assetLocks, err
}

func ReadAssetLockByInvoice(invoice string) (*models.AssetLock, error) {
	var assetLock models.AssetLock
	err := middleware.DB.Where("invoice = ?", invoice).First(&assetLock).Error
	return &assetLock, err
}

func UpdateAssetLock(assetLock *models.AssetLock) error {
	return middleware.DB.Save(assetLock).Error
}

func UpdateAssetLocks(assetLocks *[]models.AssetLock) error {
	return middleware.DB.Save(assetLocks).Error
}

func DeleteAssetLock(id uint) error {
	var assetLock models.AssetLock
	return middleware.DB.Delete(&assetLock, id).Error
}
