package models

import (
	"gorm.io/gorm"
)

// DeviceManager represents the device_manager table structure
type DeviceManager struct {
	gorm.Model
	NpubKey         string `gorm:"unique;not null" json:"npub_key"`
	DeviceID        string `gorm:"size:40;not null" json:"device_id"`
	EncryptDeviceID string `gorm:"size:1024" json:"encrypt_device_id"`
	Status          int    `gorm:"int" json:"status"`
}

func (DeviceManager) TableName() string {
	return "device_manager"
}
