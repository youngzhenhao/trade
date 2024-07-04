package services

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
	err := middleware.DB.Find(&addrReceiveEvents).Error
	return &addrReceiveEvents, err
}

func ReadAddrReceiveEvent(id uint) (*models.AddrReceiveEvent, error) {
	var addrReceiveEvent models.AddrReceiveEvent
	err := middleware.DB.First(&addrReceiveEvent, id).Error
	return &addrReceiveEvent, err
}

func ReadAddrReceiveEventsByUserId(userId int) (*[]models.AddrReceiveEvent, error) {
	var addrReceiveEvents []models.AddrReceiveEvent
	err := middleware.DB.Where("user_id = ? AND status = ?", userId, 1).Find(&addrReceiveEvents).Error
	return &addrReceiveEvents, err
}

func ReadAddrReceiveEventByAddrEncoded(addrEncoded string) (*models.AddrReceiveEvent, error) {
	var addrReceiveEvent models.AddrReceiveEvent
	err := middleware.DB.Where("addr_encoded = ? AND status = ?", addrEncoded, 1).First(&addrReceiveEvent).Error
	return &addrReceiveEvent, err
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
