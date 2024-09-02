package btldb

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
	err := middleware.DB.Order("transfer_timestamp desc").Find(&assetTransfers).Error
	return &assetTransfers, err
}

func ReadAssetTransfer(id uint) (*models.AssetTransfer, error) {
	var assetTransfer models.AssetTransfer
	err := middleware.DB.First(&assetTransfer, id).Error
	return &assetTransfer, err
}

func ReadAssetTransfersByUserId(userId int) (*[]models.AssetTransfer, error) {
	var assetTransfers []models.AssetTransfer
	err := middleware.DB.Where("user_id = ?", userId).Order("transfer_timestamp desc").Find(&assetTransfers).Error
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
	err := middleware.DB.Order("transfer_timestamp desc").Find(&assetTransferProcessedSlice).Error
	return &assetTransferProcessedSlice, err
}

func ReadAssetTransferProcessed(id uint) (*models.AssetTransferProcessedDb, error) {
	var assetTransferProcessed models.AssetTransferProcessedDb
	err := middleware.DB.First(&assetTransferProcessed, id).Error
	return &assetTransferProcessed, err
}

func ReadAssetTransferProcessedSliceByUserId(userId int) (*[]models.AssetTransferProcessedDb, error) {
	var assetTransferProcessedSlice []models.AssetTransferProcessedDb
	err := middleware.DB.Where("user_id = ?", userId).Order("transfer_timestamp desc").Find(&assetTransferProcessedSlice).Error
	return &assetTransferProcessedSlice, err
}

func ReadAssetTransferProcessedSliceByAssetId(assetId string) (*[]models.AssetTransferProcessedDb, error) {
	var assetTransferProcessedSlice []models.AssetTransferProcessedDb
	err := middleware.DB.Where("asset_id = ?", assetId).Order("transfer_timestamp desc").Find(&assetTransferProcessedSlice).Error
	return &assetTransferProcessedSlice, err
}

func ReadAssetTransferProcessedSliceByAssetIdLimit(assetId string, limit int) (*[]models.AssetTransferProcessedDb, error) {
	var assetTransferProcessedSlice []models.AssetTransferProcessedDb
	err := middleware.DB.Where("asset_id = ?", assetId).Limit(limit).Order("transfer_timestamp desc").Find(&assetTransferProcessedSlice).Error
	return &assetTransferProcessedSlice, err
}

func ReadAssetTransferProcessedByTxid(txid string) (*models.AssetTransferProcessedDb, error) {
	var assetTransferProcessed models.AssetTransferProcessedDb
	err := middleware.DB.Where("txid = ?", txid).First(&assetTransferProcessed).Error
	return &assetTransferProcessed, err
}

func ReadAssetTransferProcessedSliceByTxid(txid string) (*[]models.AssetTransferProcessedDb, error) {
	var assetTransferProcessedSlice []models.AssetTransferProcessedDb
	err := middleware.DB.Where("txid = ?", txid).First(&assetTransferProcessedSlice).Error
	return &assetTransferProcessedSlice, err
}

func ReadAssetTransferProcessedByAnchorTxHash(anchorTxHash string) (*models.AssetTransferProcessedDb, error) {
	var assetTransferProcessed models.AssetTransferProcessedDb
	err := middleware.DB.Where("anchor_tx_hash = ?", anchorTxHash).First(&assetTransferProcessed).Error
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

func DeleteAssetTransferProcessedSlice(assetTransferProcessedSlice *[]models.AssetTransferProcessedDb) error {
	return middleware.DB.Delete(&assetTransferProcessedSlice).Error
}

// AssetTransferProcessedInputDb

func CreateAssetTransferProcessedInput(assetTransferProcessedInput *models.AssetTransferProcessedInputDb) error {
	return middleware.DB.Create(assetTransferProcessedInput).Error
}

func CreateAssetTransferProcessedInputSlice(assetTransferProcessedInputSlice *[]models.AssetTransferProcessedInputDb) error {
	return middleware.DB.Create(assetTransferProcessedInputSlice).Error
}

func ReadAllAssetTransferProcessedInputSlice() (*[]models.AssetTransferProcessedInputDb, error) {
	var assetTransferProcessedInputSlice []models.AssetTransferProcessedInputDb
	err := middleware.DB.Find(&assetTransferProcessedInputSlice).Error
	return &assetTransferProcessedInputSlice, err
}

func ReadAssetTransferProcessedInput(id uint) (*models.AssetTransferProcessedInputDb, error) {
	var assetTransferProcessedInput models.AssetTransferProcessedInputDb
	err := middleware.DB.First(&assetTransferProcessedInput, id).Error
	return &assetTransferProcessedInput, err
}

func ReadAssetTransferProcessedInputSliceByUserId(userId int) (*[]models.AssetTransferProcessedInputDb, error) {
	var assetTransferProcessedInputSlice []models.AssetTransferProcessedInputDb
	err := middleware.DB.Where("user_id = ?", userId).Find(&assetTransferProcessedInputSlice).Error
	return &assetTransferProcessedInputSlice, err
}

func ReadAssetTransferProcessedInputSliceByAssetId(assetId string) (*[]models.AssetTransferProcessedInputDb, error) {
	var assetTransferProcessedInputSlice []models.AssetTransferProcessedInputDb
	err := middleware.DB.Where("asset_id = ?", assetId).Find(&assetTransferProcessedInputSlice).Error
	return &assetTransferProcessedInputSlice, err
}

func ReadAssetTransferProcessedInputSliceByTxid(txid string) (*[]models.AssetTransferProcessedInputDb, error) {
	var assetTransferProcessedInputSlice []models.AssetTransferProcessedInputDb
	err := middleware.DB.Where("txid = ?", txid).Find(&assetTransferProcessedInputSlice).Error
	return &assetTransferProcessedInputSlice, err
}

// Deprecated
func ReadAssetTransferProcessedInputSliceByAssetIdLimit(assetId string, limit int) (*[]models.AssetTransferProcessedInputDb, error) {
	var assetTransferProcessedInputSlice []models.AssetTransferProcessedInputDb
	err := middleware.DB.Where("asset_id = ?", assetId).Limit(limit).Find(&assetTransferProcessedInputSlice).Error
	return &assetTransferProcessedInputSlice, err
}

func ReadAssetTransferProcessedInputByTxid(txid string) (*models.AssetTransferProcessedInputDb, error) {
	var assetTransferProcessedInput models.AssetTransferProcessedInputDb
	err := middleware.DB.Where("txid = ?", txid).First(&assetTransferProcessedInput).Error
	return &assetTransferProcessedInput, err
}

// ReadAssetTransferProcessedInputByTxidAndIndex
// @dev: `index`
func ReadAssetTransferProcessedInputByTxidAndIndex(txid string, index int) (*models.AssetTransferProcessedInputDb, error) {
	var assetTransferProcessedInput models.AssetTransferProcessedInputDb
	err := middleware.DB.Where("txid = ? AND `index` = ?", txid, index).First(&assetTransferProcessedInput).Error
	return &assetTransferProcessedInput, err
}

func UpdateAssetTransferProcessedInput(assetTransferProcessedInput *models.AssetTransferProcessedInputDb) error {
	return middleware.DB.Save(assetTransferProcessedInput).Error
}

func UpdateAssetTransferProcessedInputSlice(assetTransferProcessedInputSlice *[]models.AssetTransferProcessedInputDb) error {
	return middleware.DB.Save(assetTransferProcessedInputSlice).Error
}

func DeleteAssetTransferProcessedInput(id uint) error {
	var assetTransferProcessedInput models.AssetTransferProcessedInputDb
	return middleware.DB.Delete(&assetTransferProcessedInput, id).Error
}

func DeleteAssetTransferProcessedInputSlice(assetTransferProcessedInputSlice *[]models.AssetTransferProcessedInputDb) error {
	return middleware.DB.Delete(assetTransferProcessedInputSlice).Error
}

// AssetTransferProcessedOutputDb

func CreateAssetTransferProcessedOutput(assetTransferProcessedOutput *models.AssetTransferProcessedOutputDb) error {
	return middleware.DB.Create(assetTransferProcessedOutput).Error
}

func CreateAssetTransferProcessedOutputSlice(assetTransferProcessedOutputSlice *[]models.AssetTransferProcessedOutputDb) error {
	return middleware.DB.Create(assetTransferProcessedOutputSlice).Error
}

func ReadAllAssetTransferProcessedOutputSlice() (*[]models.AssetTransferProcessedOutputDb, error) {
	var assetTransferProcessedOutputSlice []models.AssetTransferProcessedOutputDb
	err := middleware.DB.Find(&assetTransferProcessedOutputSlice).Error
	return &assetTransferProcessedOutputSlice, err
}

func ReadAssetTransferProcessedOutput(id uint) (*models.AssetTransferProcessedOutputDb, error) {
	var assetTransferProcessedOutput models.AssetTransferProcessedOutputDb
	err := middleware.DB.First(&assetTransferProcessedOutput, id).Error
	return &assetTransferProcessedOutput, err
}

func ReadAssetTransferProcessedOutputSliceByUserId(userId int) (*[]models.AssetTransferProcessedOutputDb, error) {
	var assetTransferProcessedOutputSlice []models.AssetTransferProcessedOutputDb
	err := middleware.DB.Where("user_id = ?", userId).Find(&assetTransferProcessedOutputSlice).Error
	return &assetTransferProcessedOutputSlice, err
}

func ReadAssetTransferProcessedOutputSliceByAssetId(assetId string) (*[]models.AssetTransferProcessedOutputDb, error) {
	var assetTransferProcessedOutputSlice []models.AssetTransferProcessedOutputDb
	err := middleware.DB.Where("asset_id = ?", assetId).Find(&assetTransferProcessedOutputSlice).Error
	return &assetTransferProcessedOutputSlice, err
}

// Deprecated
func ReadAssetTransferProcessedOutputSliceByAssetIdLimit(assetId string, limit int) (*[]models.AssetTransferProcessedOutputDb, error) {
	var assetTransferProcessedOutputSlice []models.AssetTransferProcessedOutputDb
	err := middleware.DB.Where("asset_id = ?", assetId).Limit(limit).Find(&assetTransferProcessedOutputSlice).Error
	return &assetTransferProcessedOutputSlice, err
}

func ReadAssetTransferProcessedOutputByTxid(txid string) (*models.AssetTransferProcessedOutputDb, error) {
	var assetTransferProcessedOutput models.AssetTransferProcessedOutputDb
	err := middleware.DB.Where("txid = ?", txid).First(&assetTransferProcessedOutput).Error
	return &assetTransferProcessedOutput, err
}

func ReadAssetTransferProcessedOutputSliceByTxid(txid string) (*[]models.AssetTransferProcessedOutputDb, error) {
	var assetTransferProcessedOutputSlice []models.AssetTransferProcessedOutputDb
	err := middleware.DB.Where("txid = ?", txid).Find(&assetTransferProcessedOutputSlice).Error
	return &assetTransferProcessedOutputSlice, err
}

func ReadAssetTransferProcessedOutputSliceWhoseAddressIsNull() (*[]models.AssetTransferProcessedOutputDb, error) {
	var assetTransferProcessedOutputSlice []models.AssetTransferProcessedOutputDb
	err := middleware.DB.Where("address = ?", "").Find(&assetTransferProcessedOutputSlice).Error
	return &assetTransferProcessedOutputSlice, err
}

// ReadAssetTransferProcessedOutputByTxidAndIndex
// @dev: `index`
func ReadAssetTransferProcessedOutputByTxidAndIndex(txid string, index int) (*models.AssetTransferProcessedOutputDb, error) {
	var assetTransferProcessedOutput models.AssetTransferProcessedOutputDb
	err := middleware.DB.Where("txid = ? AND `index` = ?", txid, index).First(&assetTransferProcessedOutput).Error
	return &assetTransferProcessedOutput, err
}

func UpdateAssetTransferProcessedOutput(assetTransferProcessedOutput *models.AssetTransferProcessedOutputDb) error {
	return middleware.DB.Save(assetTransferProcessedOutput).Error
}

func UpdateAssetTransferProcessedOutputSlice(assetTransferProcessedOutputSlice *[]models.AssetTransferProcessedOutputDb) error {
	return middleware.DB.Save(assetTransferProcessedOutputSlice).Error
}

func DeleteAssetTransferProcessedOutput(id uint) error {
	var assetTransferProcessedOutput models.AssetTransferProcessedOutputDb
	return middleware.DB.Delete(&assetTransferProcessedOutput, id).Error
}

func DeleteAssetTransferProcessedOutputSlice(assetTransferProcessedOutputSlice *[]models.AssetTransferProcessedOutputDb) error {
	return middleware.DB.Delete(&assetTransferProcessedOutputSlice).Error
}
