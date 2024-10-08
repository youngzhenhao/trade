package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateNftTransfer(nftTransfer *models.NftTransfer) error {
	return middleware.DB.Create(nftTransfer).Error
}

func CreateNftTransfers(nftTransfers *[]models.NftTransfer) error {
	return middleware.DB.Create(nftTransfers).Error
}

func ReadNftTransfer(id uint) (*models.NftTransfer, error) {
	var nftTransfer models.NftTransfer
	err := middleware.DB.First(&nftTransfer, id).Error
	return &nftTransfer, err
}

func ReadAllNftTransfers() (*[]models.NftTransfer, error) {
	var nftTransfers []models.NftTransfer
	err := middleware.DB.Order("updated_at desc").Find(&nftTransfers).Error
	return &nftTransfers, err
}

func ReadNftTransfersByAssetId(assetId string) (*[]models.NftTransfer, error) {
	var nftTransfers []models.NftTransfer
	err := middleware.DB.Where("asset_id = ?", assetId).Find(&nftTransfers).Error
	return &nftTransfers, err
}

func UpdateNftTransfer(nftTransfer *models.NftTransfer) error {
	return middleware.DB.Save(nftTransfer).Error
}

func UpdateNftTransfers(nftTransfers *[]models.NftTransfer) error {
	return middleware.DB.Save(nftTransfers).Error
}

func DeleteNftTransfer(id uint) error {
	var nftTransfer models.NftTransfer
	return middleware.DB.Delete(&nftTransfer, id).Error
}
