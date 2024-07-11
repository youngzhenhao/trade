package models

import "gorm.io/gorm"

type AssetBalance struct {
	gorm.Model
	GenesisPoint string `json:"genesis_point"`
	Name         string `json:"name"`
	MetaHash     string `json:"meta_hash"`
	AssetID      string `json:"asset_id"`
	AssetType    string `json:"asset_type"`
	OutputIndex  int    `json:"output_index"`
	Version      int    `json:"version"`
	Balance      int    `json:"balance"`
	DeviceId     string `json:"device_id" gorm:"type:varchar(255)"`
	UserId       int    `json:"user_id"`
	Status       int    `json:"status" gorm:"default:1"`
}

type AssetBalanceSetRequest struct {
	GenesisPoint string `json:"genesis_point"`
	Name         string `json:"name"`
	MetaHash     string `json:"meta_hash"`
	AssetID      string `json:"asset_id"`
	AssetType    string `json:"asset_type"`
	OutputIndex  int    `json:"output_index"`
	Version      int    `json:"version"`
	Balance      int    `json:"balance"`
	DeviceId     string `json:"device_id" gorm:"type:varchar(255)"`
}
