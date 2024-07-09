package models

import "gorm.io/gorm"

type AssetLock struct {
	gorm.Model
	AssetId          string `json:"asset_id" gorm:"type:varchar(255)"`
	AssetName        string `json:"asset_name" gorm:"type:varchar(255)"`
	AssetType        string `json:"asset_type" gorm:"type:varchar(255)"`
	LockAmount       int    `json:"lock_amount"`
	LockTime         int    `json:"lock_time"`
	RelativeLockTime int    `json:"relative_lock_time"`
	HashLock         string `json:"hash_lock" gorm:"type:varchar(255)"`
	Invoice          string `json:"invoice" gorm:"type:varchar(255)"`
	DeviceId         string `json:"device_id" gorm:"type:varchar(255)"`
	UserId           int    `json:"user_id"`
	Status           int    `json:"status" gorm:"default:1"`
}

type AssetLockSetRequest struct {
	AssetId          string `json:"asset_id" gorm:"type:varchar(255)"`
	AssetName        string `json:"asset_name" gorm:"type:varchar(255)"`
	AssetType        string `json:"asset_type" gorm:"type:varchar(255)"`
	LockAmount       int    `json:"lock_amount"`
	LockTime         int    `json:"lock_time"`
	RelativeLockTime int    `json:"relative_lock_time"`
	HashLock         string `json:"hash_lock" gorm:"type:varchar(255)"`
	Invoice          string `json:"invoice" gorm:"type:varchar(255)"`
	DeviceId         string `json:"device_id" gorm:"type:varchar(255)"`
}
