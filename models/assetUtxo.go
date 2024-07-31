package models

import "gorm.io/gorm"

type AssetUtxo struct {
	gorm.Model
	AssetId string `json:"asset_id" gorm:"type:varchar(255)"`
	// TODO:
	DeviceId string `json:"device_id" gorm:"type:varchar(255)"`
	UserId   int    `json:"user_id"`
	Username string `json:"username" gorm:"type:varchar(255)"`
	Status   int    `json:"status" gorm:"default:1"`
}
