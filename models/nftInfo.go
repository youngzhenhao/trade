package models

import "gorm.io/gorm"

type NftInfo struct {
	gorm.Model
	AssetId         string `json:"asset_id" gorm:"type:varchar(255);index"`
	Name            string `json:"name"`
	Version         string `json:"version" gorm:"type:varchar(255);index"`
	AssetType       string `json:"asset_type" gorm:"type:varchar(255);index"`
	GenesisPoint    string `json:"genesis_point" gorm:"type:varchar(255)"`
	MetaHash        string `json:"meta_hash" gorm:"type:varchar(255)"`
	TweakedGroupKey string `json:"tweaked_group_key" gorm:"type:varchar(255)"`
	DeviceId        string `json:"device_id" gorm:"type:varchar(255);index"`
	UserId          int    `json:"user_id" gorm:"index"`
	Username        string `json:"username" gorm:"type:varchar(255);index"`
}

type NftInfoSetRequest struct {
	AssetId         string `json:"asset_id" gorm:"type:varchar(255);index"`
	Name            string `json:"name"`
	Version         string `json:"version" gorm:"type:varchar(255);index"`
	AssetType       string `json:"asset_type" gorm:"type:varchar(255);index"`
	GenesisPoint    string `json:"genesis_point" gorm:"type:varchar(255)"`
	MetaHash        string `json:"meta_hash" gorm:"type:varchar(255)"`
	TweakedGroupKey string `json:"tweaked_group_key" gorm:"type:varchar(255)"`
	DeviceId        string `json:"device_id" gorm:"type:varchar(255);index"`
}
