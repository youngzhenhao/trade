package services

import (
	"errors"
	"strconv"
	"trade/api"
	"trade/middleware"
	"trade/models"
)

func CreateBatchTransfer(batchTransfer *models.BatchTransfer) error {
	return middleware.DB.Create(batchTransfer).Error
}

func CreateBatchTransfers(batchTransfers *[]models.BatchTransfer) error {
	return middleware.DB.Create(batchTransfers).Error
}

func ReadAllBatchTransfers() (*[]models.BatchTransfer, error) {
	var batchTransfers []models.BatchTransfer
	err := middleware.DB.Find(&batchTransfers).Error
	return &batchTransfers, err
}

func ReadBatchTransfer(id uint) (*models.BatchTransfer, error) {
	var batchTransfer models.BatchTransfer
	err := middleware.DB.First(&batchTransfer, id).Error
	return &batchTransfer, err
}

func ReadBatchTransfersByUserId(userId int) (*[]models.BatchTransfer, error) {
	var batchTransfers []models.BatchTransfer
	err := middleware.DB.Where("user_id = ? AND status = ?", userId, 1).Find(&batchTransfers).Error
	return &batchTransfers, err
}

func ReadBatchTransferByAddrEncoded(encoded string) (*models.BatchTransfer, error) {
	var batchTransfer models.BatchTransfer
	err := middleware.DB.Where("encoded = ? AND status = ?", encoded, 1).First(&batchTransfer).Error
	return &batchTransfer, err
}

func ReadBatchTransferByTxid(txid string) (*models.BatchTransfer, error) {
	var batchTransfer models.BatchTransfer
	err := middleware.DB.Where("txid = ? AND status = ?", txid, 1).First(&batchTransfer).Error
	return &batchTransfer, err
}

func ReadBatchTransferByAddrEncodedAndIndex(encoded string, index int) (*models.BatchTransfer, error) {
	var batchTransfer models.BatchTransfer
	err := middleware.DB.Where("encoded = ? AND index = ? AND status = ?", encoded, index, 1).First(&batchTransfer).Error
	return &batchTransfer, err
}

func ReadBatchTransferByTxidAndIndex(txid string, index int) (*models.BatchTransfer, error) {
	var batchTransfer models.BatchTransfer
	err := middleware.DB.Where("txid = ? AND index = ? AND status = ?", txid, index, 1).First(&batchTransfer).Error
	return &batchTransfer, err
}

func ReadBatchTransferByOutpoint(outpoint string) (*models.BatchTransfer, error) {
	txid, indexStr := api.OutpointToTransactionAndIndex(outpoint)
	if txid == "" || indexStr == "" {
		return nil, errors.New("invalid outpoint or index")
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return nil, err
	}
	return ReadBatchTransferByTxidAndIndex(txid, index)
}

func UpdateBatchTransfer(batchTransfer *models.BatchTransfer) error {
	return middleware.DB.Save(batchTransfer).Error
}

func UpdateBatchTransfers(batchTransfers *[]models.BatchTransfer) error {
	return middleware.DB.Save(batchTransfers).Error
}

func DeleteBatchTransfer(id uint) error {
	var batchTransfer models.BatchTransfer
	return middleware.DB.Delete(&batchTransfer, id).Error
}
