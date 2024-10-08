package btldb

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
	err := middleware.DB.Order("updated_at desc").Find(&assetAddrs).Error
	return &assetAddrs, err
}

func ReadAssetAddr(id uint) (*models.AssetAddr, error) {
	var assetAddr models.AssetAddr
	err := middleware.DB.First(&assetAddr, id).Error
	return &assetAddr, err
}

func ReadAssetAddrsByUserId(userId int) (*[]models.AssetAddr, error) {
	var assetAddrs []models.AssetAddr
	err := middleware.DB.Where("user_id = ?", userId).Find(&assetAddrs).Error
	return &assetAddrs, err
}

func ReadAssetAddrsByScriptKey(scriptKey string) (*[]models.AssetAddr, error) {
	var assetAddrs []models.AssetAddr
	err := middleware.DB.Where("script_key = ?", scriptKey).Find(&assetAddrs).Error
	return &assetAddrs, err
}

func ReadAssetAddrByAddrEncoded(addrEncoded string) (*models.AssetAddr, error) {
	var assetAddr models.AssetAddr
	err := middleware.DB.Where("encoded = ?", addrEncoded).First(&assetAddr).Error
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
