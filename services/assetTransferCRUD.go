package services

import (
	"trade/middleware"
	"trade/models"
)

// AssetTransfer

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

// AssetTransferProcessedDb

func CreateAssetTransferProcessed(assetTransferProcessed *models.AssetTransferProcessedDb) error {
	return middleware.DB.Create(assetTransferProcessed).Error
}

func CreateAssetTransferProcessedSlice(assetTransferProcessedSlice *[]models.AssetTransferProcessedDb) error {
	return middleware.DB.Create(assetTransferProcessedSlice).Error
}

func ReadAllAssetTransferProcessedSlice() (*[]models.AssetTransferProcessedDb, error) {
	var assetTransferProcessedSlice []models.AssetTransferProcessedDb
	err := middleware.DB.Find(&assetTransferProcessedSlice).Error
	return &assetTransferProcessedSlice, err
}

func ReadAssetTransferProcessed(id uint) (*models.AssetTransferProcessedDb, error) {
	var assetTransferProcessed models.AssetTransferProcessedDb
	err := middleware.DB.First(&assetTransferProcessed, id).Error
	return &assetTransferProcessed, err
}

func ReadAssetTransferProcessedSliceByUserId(userId int) (*[]models.AssetTransferProcessedDb, error) {
	var assetTransferProcessedSlice []models.AssetTransferProcessedDb
	err := middleware.DB.Where("user_id = ? AND status = ?", userId, 1).Find(&assetTransferProcessedSlice).Error
	return &assetTransferProcessedSlice, err
}

func ReadAssetTransferProcessedByTxid(txid string) (*models.AssetTransferProcessedDb, error) {
	var assetTransferProcessed models.AssetTransferProcessedDb
	err := middleware.DB.Where("txid = ? AND status = ?", txid, 1).First(&assetTransferProcessed).Error
	return &assetTransferProcessed, err
}

func ReadAssetTransferProcessedByAnchorTxHash(anchorTxHash string) (*models.AssetTransferProcessedDb, error) {
	var assetTransferProcessed models.AssetTransferProcessedDb
	err := middleware.DB.Where("anchor_tx_hash = ? AND status = ?", anchorTxHash, 1).First(&assetTransferProcessed).Error
	return &assetTransferProcessed, err
}

func UpdateAssetTransferProcessed(assetTransferProcessed *models.AssetTransferProcessedDb) error {
	return middleware.DB.Save(assetTransferProcessed).Error
}

func UpdateAssetTransferProcessedSlice(assetTransferProcessedSlice *[]models.AssetTransferProcessedDb) error {
	return middleware.DB.Save(assetTransferProcessedSlice).Error
}

func DeleteAssetTransferProcessed(id uint) error {
	var assetTransferProcessed models.AssetTransferProcessedDb
	return middleware.DB.Delete(&assetTransferProcessed, id).Error
}

// TODO: Inputs and outputs
