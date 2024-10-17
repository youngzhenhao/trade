package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateNftPresaleWhitelist(nftPresaleWhitelist *models.NftPresaleWhitelist) error {
	return middleware.DB.Create(nftPresaleWhitelist).Error
}

func CreateNftPresaleWhitelists(nftPresaleWhitelists *[]models.NftPresaleWhitelist) error {
	return middleware.DB.Create(nftPresaleWhitelists).Error
}

func ReadNftPresaleWhitelist(id uint) (*models.NftPresaleWhitelist, error) {
	var nftPresaleWhitelist models.NftPresaleWhitelist
	err := middleware.DB.First(&nftPresaleWhitelist, id).Error
	return &nftPresaleWhitelist, err
}

func ReadNftPresaleWhitelistByAssetId(assetId string) (*models.NftPresaleWhitelist, error) {
	var nftPresaleWhitelist models.NftPresaleWhitelist
	err := middleware.DB.Where("asset_id = ?", assetId).First(&nftPresaleWhitelist).Error
	return &nftPresaleWhitelist, err
}

func ReadNftPresaleWhitelistByBatchGroupId(batchGroupId int) (*models.NftPresaleWhitelist, error) {
	var nftPresaleWhitelist models.NftPresaleWhitelist
	err := middleware.DB.Where("batch_group_id = ?", batchGroupId).First(&nftPresaleWhitelist).Error
	return &nftPresaleWhitelist, err
}

func ReadAllNftPresaleWhitelists() (*[]models.NftPresaleWhitelist, error) {
	var nftPresaleWhitelists []models.NftPresaleWhitelist
	err := middleware.DB.Find(&nftPresaleWhitelists).Error
	return &nftPresaleWhitelists, err
}

func UpdateNftPresaleWhitelist(nftPresaleWhitelist *models.NftPresaleWhitelist) error {
	return middleware.DB.Save(nftPresaleWhitelist).Error
}

func UpdateNftPresaleWhitelists(nftPresaleWhitelists *[]models.NftPresaleWhitelist) error {
	return middleware.DB.Save(nftPresaleWhitelists).Error
}

func DeleteNftPresaleWhitelist(id uint) error {
	var nftPresaleWhitelist models.NftPresaleWhitelist
	return middleware.DB.Delete(&nftPresaleWhitelist, id).Error
}
