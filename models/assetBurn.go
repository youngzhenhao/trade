package models

import "gorm.io/gorm"

type AssetBurn struct {
	gorm.Model
	AssetId  string `json:"asset_id" gorm:"type:varchar(255)"`
	Amount   string `json:"amount" gorm:"type:varchar(255)"`
	DeviceId string `json:"device_id" gorm:"type:varchar(255)"`
	UserId   int    `json:"user_id"`
	Username string `json:"username" gorm:"type:varchar(255)"`
	Status   int    `json:"status" gorm:"default:1"`
}

type AssetBurnSetRequest struct {
	AssetId  string `json:"asset_id"`
	Amount   string `json:"amount"`
	DeviceId string `json:"device_id"`
}
