package services

import (
	"trade/middleware"
	"trade/models"
)

func CreateAssetAddr(assetAddr *models.AssetAddr) error {
	return middleware.DB.Create(assetAddr).Error
}

func CreateAssetAddrs(assetAddrs *[]models.AssetAddr) error {
	return middleware.DB.Create(assetAddrs).Error
}

func ReadAllAssetAddrs() (*[]models.AssetAddr, error) {
	var assetAddrs []models.AssetAddr
	err := middleware.DB.Find(&assetAddrs).Error
	return &assetAddrs, err
}

func ReadAssetAddr(id uint) (*models.AssetAddr, error) {
	var assetAddr models.AssetAddr
	err := middleware.DB.First(&assetAddr, id).Error
	return &assetAddr, err
}

func ReadAssetAddrsByUserId(userId int) (*[]models.AssetAddr, error) {
	var assetAddrs []models.AssetAddr
	err := middleware.DB.Where("user_id = ? AND status = ?", userId, 1).Find(&assetAddrs).Error
	return &assetAddrs, err
}

func ReadAssetAddrByAddrEncoded(addrEncoded string) (*models.AssetAddr, error) {
	var assetAddr models.AssetAddr
	err := middleware.DB.Where("encoded = ? AND status = ?", addrEncoded, 1).First(&assetAddr).Error
	return &assetAddr, err
}

func UpdateAssetAddr(assetAddr *models.AssetAddr) error {
	return middleware.DB.Save(assetAddr).Error
}

func UpdateAssetAddrs(assetAddrs *[]models.AssetAddr) error {
	return middleware.DB.Save(assetAddrs).Error
}

func DeleteAssetAddr(id uint) error {
	var assetAddr models.AssetAddr
	return middleware.DB.Delete(&assetAddr, id).Error
}
