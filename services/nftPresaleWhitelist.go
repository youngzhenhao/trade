package services

import (
	"trade/models"
	"trade/services/btldb"
)

func CreateNftPresaleWhitelist(nftPresaleWhitelist *models.NftPresaleWhitelist) error {
	return btldb.CreateNftPresaleWhitelist(nftPresaleWhitelist)
}

func CreateNftPresaleWhitelists(nftPresaleWhitelists *[]models.NftPresaleWhitelist) error {
	return btldb.CreateNftPresaleWhitelists(nftPresaleWhitelists)
}

func ReadNftPresaleWhitelist(id uint) (*models.NftPresaleWhitelist, error) {
	return btldb.ReadNftPresaleWhitelist(id)
}

func ReadNftPresaleWhitelistByAssetId(assetId string) (*models.NftPresaleWhitelist, error) {
	return btldb.ReadNftPresaleWhitelistByAssetId(assetId)
}

func ReadNftPresaleWhitelistByBatchGroupId(batchGroupId int) (*models.NftPresaleWhitelist, error) {
	return btldb.ReadNftPresaleWhitelistByBatchGroupId(batchGroupId)
}

func ReadAllNftPresaleWhitelists() (*[]models.NftPresaleWhitelist, error) {
	return btldb.ReadAllNftPresaleWhitelists()
}

func UpdateNftPresaleWhitelist(nftPresaleWhitelist *models.NftPresaleWhitelist) error {
	return btldb.UpdateNftPresaleWhitelist(nftPresaleWhitelist)
}

func UpdateNftPresaleWhitelists(nftPresaleWhitelists *[]models.NftPresaleWhitelist) error {
	return btldb.UpdateNftPresaleWhitelists(nftPresaleWhitelists)
}

func DeleteNftPresaleWhitelist(id uint) error {
	return btldb.DeleteNftPresaleWhitelist(id)
}
