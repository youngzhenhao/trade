package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateNftPresale(nftPresale *models.NftPresale) error {
	return middleware.DB.Create(nftPresale).Error
}

func CreateNftPresales(nftPresales *[]models.NftPresale) error {
	return middleware.DB.Create(nftPresales).Error
}

func ReadNftPresale(id uint) (*models.NftPresale, error) {
	var nftPresale models.NftPresale
	err := middleware.DB.First(&nftPresale, id).Error
	return &nftPresale, err
}

func ReadNftPresaleByAssetId(assetId string) (*models.NftPresale, error) {
	var nftPresale models.NftPresale
	err := middleware.DB.Where("asset_id = ?", assetId).First(&nftPresale).Error
	return &nftPresale, err
}

func ReadAllNftPresales() (*[]models.NftPresale, error) {
	var nftPresales []models.NftPresale
	err := middleware.DB.Order("launch_time desc").Find(&nftPresales).Error
	return &nftPresales, err
}

func ReadNftPresalesByNftPresaleState(nftPresaleState models.NftPresaleState) (*[]models.NftPresale, error) {
	var nftPresales []models.NftPresale
	err := middleware.DB.Where("state = ?", nftPresaleState).Order("launch_time desc").Find(&nftPresales).Error
	return &nftPresales, err
}

func ReadNftPresalesBetweenNftPresaleState(stateStart models.NftPresaleState, stateEnd models.NftPresaleState) (*[]models.NftPresale, error) {
	var nftPresales []models.NftPresale
	err := middleware.DB.Where("state BETWEEN ? AND ?", stateStart, stateEnd).Order("launch_time desc").Find(&nftPresales).Error
	return &nftPresales, err
}

func ReadNftPresalesByBuyerUserId(userId int) (*[]models.NftPresale, error) {
	var nftPresales []models.NftPresale
	err := middleware.DB.Where("buyer_user_id = ?", userId).Order("launch_time desc").Find(&nftPresales).Error
	return &nftPresales, err
}

func UpdateNftPresale(nftPresale *models.NftPresale) error {
	return middleware.DB.Save(nftPresale).Error
}

func UpdateNftPresales(nftPresales *[]models.NftPresale) error {
	return middleware.DB.Save(nftPresales).Error
}

func DeleteNftPresale(id uint) error {
	var nftPresale models.NftPresale
	return middleware.DB.Delete(&nftPresale, id).Error
}
