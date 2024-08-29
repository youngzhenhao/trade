package models

import "gorm.io/gorm"

type AssetBalance struct {
	gorm.Model
	GenesisPoint string `json:"genesis_point" gorm:"type:varchar(255)"`
	Name         string `json:"name" gorm:"type:varchar(255)"`
	MetaHash     string `json:"meta_hash" gorm:"type:varchar(255)"`
	AssetID      string `json:"asset_id" gorm:"type:varchar(255)"`
	AssetType    string `json:"asset_type" gorm:"type:varchar(255)"`
	OutputIndex  int    `json:"output_index"`
	Version      int    `json:"version"`
	Balance      int    `json:"balance"`
	DeviceId     string `json:"device_id" gorm:"type:varchar(255)"`
	UserId       int    `json:"user_id"`
	Username     string `json:"username" gorm:"type:varchar(255)"`
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

type AssetHolderBalanceLimitAndOffsetRequest struct {
	AssetId string `json:"asset_id"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}

type GetAssetBalanceByUserIdAndAssetIdRequest struct {
	UserId  int    `json:"user_id"`
	AssetId string `json:"asset_id"`
}
