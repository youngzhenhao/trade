package models

import "gorm.io/gorm"

type FairLaunchFollow struct {
	gorm.Model
	FairLaunchInfoId int    `json:"fair_launch_info_id"`
	AssetId          string `json:"asset_id" gorm:"type:varchar(255)"`
	DeviceId         string `json:"device_id" gorm:"type:varchar(255)"`
	UserId           int    `json:"user_id"`
	Username         string `json:"username" gorm:"type:varchar(255)"`
	Status           int    `json:"status" gorm:"default:1"`
}

type FairLaunchFollowSetRequest struct {
	FairLaunchInfoId int    `json:"fair_launch_info_id"`
	AssetId          string `json:"asset_id" gorm:"type:varchar(255)"`
	DeviceId         string `json:"device_id" gorm:"type:varchar(255)"`
}
