package models

import "gorm.io/gorm"

type AssetRecommend struct {
	gorm.Model
	AssetId           string `json:"asset_id" gorm:"type:varchar(255)"`
	AssetFromAddr     string `json:"asset_from_addr" gorm:"type:varchar(255)"`
	RecommendUserId   int    `json:"recommend_user_id"`
	RecommendUsername string `json:"recommend_username" gorm:"type:varchar(255)"`
	RecommendTime     int    `json:"recommend_time"`
	DeviceId          string `json:"device_id" gorm:"type:varchar(255)"`
	UserId            int    `json:"user_id"`
	Username          string `json:"username" gorm:"type:varchar(255)"`
	Status            int    `json:"status" gorm:"default:1"`
}

type AssetRecommendSetRequest struct {
	AssetId           string `json:"asset_id"`
	AssetFromAddr     string `json:"asset_from_addr"`
	RecommendUserId   int    `json:"recommend_user_id"`
	RecommendUsername string `json:"recommend_username"`
	RecommendTime     int    `json:"recommend_time"`
	DeviceId          string `json:"device_id"`
}
