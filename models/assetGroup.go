package models

import "gorm.io/gorm"

type AssetGroup struct {
	gorm.Model
	TweakedGroupKey string `json:"tweaked_group_key" gorm:"type:varchar(255)"`
	FirstAssetMeta  string `json:"first_asset_meta"`
	FirstAssetId    string `json:"first_asset_id" gorm:"type:varchar(255)"`
	DeviceId        string `json:"device_id" gorm:"type:varchar(255)"`
	UserId          int    `json:"user_id"`
	Username        string `json:"username" gorm:"type:varchar(255)"`
}

type AssetGroupSetRequest struct {
	TweakedGroupKey string `json:"tweaked_group_key"`
	FirstAssetMeta  string `json:"first_asset_meta"`
	FirstAssetId    string `json:"first_asset_id" gorm:"type:varchar(255)"`
	DeviceId        string `json:"device_id" gorm:"type:varchar(255)"`
}
