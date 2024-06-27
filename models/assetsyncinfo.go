package models

import "gorm.io/gorm"

type AssetSyncInfo struct {
	gorm.Model
	//TODO: Create property
}

func (AssetSyncInfo) TableName() string {
	return "Asset_Sync_Info"
}
