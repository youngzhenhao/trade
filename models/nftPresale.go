package models

import "gorm.io/gorm"

type NftPresale struct {
	gorm.Model
	AssetId string `json:"asset_id" gorm:"type:varchar(255);index"`
	//TODO
}
