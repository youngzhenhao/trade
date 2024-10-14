package models

import "gorm.io/gorm"

type AssetMeta struct {
	gorm.Model
	AssetID   string `json:"asset_id" gorm:"type:varchar(255);index"`
	AssetMeta string `json:"asset_meta"`
}
