package models

import "gorm.io/gorm"

type NftTransfer struct {
	gorm.Model
	Txid     string `json:"txid" gorm:"type:varchar(255)"`
	AssetId  string `json:"asset_id" gorm:"type:varchar(255);index"`
	Time     int    `json:"time"`
	FromAddr string `json:"from" gorm:"type:varchar(255)"`
	ToAddr   string `json:"to" gorm:"type:varchar(255)"`
	FromInfo string `json:"from_info" gorm:"type:varchar(255)"`
	ToInfo   string `json:"to_info" gorm:"type:varchar(255)"`
	DeviceId string `json:"device_id" gorm:"type:varchar(255);index"`
	UserId   int    `json:"user_id" gorm:"index"`
	Username string `json:"username" gorm:"type:varchar(255);index"`
}

type NftTransferSetRequest struct {
	Txid     string `json:"txid" gorm:"type:varchar(255)"`
	AssetId  string `json:"asset_id" gorm:"type:varchar(255);index"`
	Time     int    `json:"time"`
	FromAddr string `json:"from" gorm:"type:varchar(255)"`
	ToAddr   string `json:"to" gorm:"type:varchar(255)"`
	FromInfo string `json:"from_info" gorm:"type:varchar(255)"`
	ToInfo   string `json:"to_info" gorm:"type:varchar(255)"`
	DeviceId string `json:"device_id" gorm:"type:varchar(255);index"`
}
