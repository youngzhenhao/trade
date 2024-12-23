package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateAddrReceiveEvent(addrReceiveEvent *models.AddrReceiveEvent) error {
	return middleware.DB.Create(addrReceiveEvent).Error
}

func CreateAddrReceiveEvents(addrReceiveEvents *[]models.AddrReceiveEvent) error {
	return middleware.DB.Create(addrReceiveEvents).Error
}

func ReadAllAddrReceiveEvents() (*[]models.AddrReceiveEvent, error) {
	var addrReceiveEvents []models.AddrReceiveEvent
	err := middleware.DB.Order("creation_time_unix_seconds desc").Find(&addrReceiveEvents).Error
	return &addrReceiveEvents, err
}

func ReadAddrReceiveEvent(id uint) (*models.AddrReceiveEvent, error) {
	var addrReceiveEvent models.AddrReceiveEvent
	err := middleware.DB.First(&addrReceiveEvent, id).Error
	return &addrReceiveEvent, err
}

func ReadAddrReceiveEventsByUserId(userId int) (*[]models.AddrReceiveEvent, error) {
	var addrReceiveEvents []models.AddrReceiveEvent
	err := middleware.DB.Where("user_id = ?", userId).Order("creation_time_unix_seconds desc").Find(&addrReceiveEvents).Error
	return &addrReceiveEvents, err
}

func ReadAddrReceiveEventByAddrEncoded(addrEncoded string) (*models.AddrReceiveEvent, error) {
	var addrReceiveEvent models.AddrReceiveEvent
	err := middleware.DB.Where("addr_encoded = ?", addrEncoded).First(&addrReceiveEvent).Error
	return &addrReceiveEvent, err
}

func ReadAddrReceiveEventsByAssetId(assetId string) (*[]models.AddrReceiveEvent, error) {
	var addrReceiveEvents []models.AddrReceiveEvent
	err := middleware.DB.Where("addr_asset_id = ?", assetId).Order("creation_time_unix_seconds desc").Find(&addrReceiveEvents).Error
	return &addrReceiveEvents, err
}

func ReadAddrReceiveEventsByUsername(username string) (*[]models.AddrReceiveEvent, error) {
	var addrReceiveEvents []models.AddrReceiveEvent
	err := middleware.DB.Where("username = ?", username).Order("creation_time_unix_seconds desc").Find(&addrReceiveEvents).Error
	return &addrReceiveEvents, err
}

func UpdateAddrReceiveEvent(addrReceiveEvent *models.AddrReceiveEvent) error {
	return middleware.DB.Save(addrReceiveEvent).Error
}

func UpdateAddrReceiveEvents(addrReceiveEvents *[]models.AddrReceiveEvent) error {
	return middleware.DB.Save(addrReceiveEvents).Error
}

func DeleteAddrReceiveEvent(id uint) error {
	var addrReceiveEvent models.AddrReceiveEvent
	return middleware.DB.Delete(&addrReceiveEvent, id).Error
}
