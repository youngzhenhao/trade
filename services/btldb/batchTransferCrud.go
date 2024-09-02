package btldb

import (
	"errors"
	"strconv"
	"trade/middleware"
	"trade/models"
	"trade/utils"
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

func ReadAllBatchTransfersUpdatedAtDesc() (*[]models.BatchTransfer, error) {
	var batchTransfers []models.BatchTransfer
	err := middleware.DB.Order("updated_at desc").Find(&batchTransfers).Error
	return &batchTransfers, err
}

func ReadBatchTransfer(id uint) (*models.BatchTransfer, error) {
	var batchTransfer models.BatchTransfer
	err := middleware.DB.First(&batchTransfer, id).Error
	return &batchTransfer, err
}

func ReadBatchTransfersByUserId(userId int) (*[]models.BatchTransfer, error) {
	var batchTransfers []models.BatchTransfer
	err := middleware.DB.Where("user_id = ?", userId).Find(&batchTransfers).Error
	return &batchTransfers, err
}

func ReadBatchTransferByAddrEncoded(encoded string) (*models.BatchTransfer, error) {
	var batchTransfer models.BatchTransfer
	err := middleware.DB.Where("encoded = ?", encoded).First(&batchTransfer).Error
	return &batchTransfer, err
}

func ReadBatchTransferByTxid(txid string) (*models.BatchTransfer, error) {
	var batchTransfer models.BatchTransfer
	err := middleware.DB.Where("txid = ?", txid).First(&batchTransfer).Error
	return &batchTransfer, err
}

// ReadBatchTransferByAddrEncodedAndIndex
// @dev: `index`
func ReadBatchTransferByAddrEncodedAndIndex(encoded string, index int) (*models.BatchTransfer, error) {
	var batchTransfer models.BatchTransfer
	err := middleware.DB.Where("encoded = ? AND `index` = ?", encoded, index).First(&batchTransfer).Error
	return &batchTransfer, err
}

// ReadBatchTransferByTxidAndIndex
// @dev: `index`
func ReadBatchTransferByTxidAndIndex(txid string, index int) (*models.BatchTransfer, error) {
	var batchTransfer models.BatchTransfer
	err := middleware.DB.Where("txid = ? AND `index` = ?", txid, index).First(&batchTransfer).Error
	return &batchTransfer, err
}

func ReadBatchTransferByOutpoint(outpoint string) (*models.BatchTransfer, error) {
	txid, indexStr := utils.OutpointToTransactionAndIndex(outpoint)
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
