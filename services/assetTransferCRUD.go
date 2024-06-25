package services

import (
	"trade/middleware"
	"trade/models"
)

func CreateAssetTransfer(assetTransfer *models.AssetTransfer) error {
	return middleware.DB.Create(assetTransfer).Error
}

func CreateAssetTransfers(assetTransfers *[]models.AssetTransfer) error {
	return middleware.DB.Create(assetTransfers).Error
}

func ReadAllAssetTransfers() (*[]models.AssetTransfer, error) {
	var assetTransfers []models.AssetTransfer
	err := middleware.DB.Find(&assetTransfers).Error
	return &assetTransfers, err
}

func ReadAssetTransfer(id uint) (*models.AssetTransfer, error) {
	var assetTransfer models.AssetTransfer
	err := middleware.DB.First(&assetTransfer, id).Error
	return &assetTransfer, err
}

func ReadAssetTransfersByUserId(userId int) (*[]models.AssetTransfer, error) {
	var assetTransfers []models.AssetTransfer
	err := middleware.DB.Where("user_id = ? AND status = ?", userId, 1).Find(&assetTransfers).Error
	return &assetTransfers, err
}

func UpdateAssetTransfer(assetTransfer *models.AssetTransfer) error {
	return middleware.DB.Save(assetTransfer).Error
}

func UpdateAssetTransfers(assetTransfers *[]models.AssetTransfer) error {
	return middleware.DB.Save(assetTransfers).Error
}

func DeleteAssetTransfer(id uint) error {
	var assetTransfer models.AssetTransfer
	return middleware.DB.Delete(&assetTransfer, id).Error
}
