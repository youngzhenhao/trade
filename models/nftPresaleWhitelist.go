package models

import (
	"gorm.io/gorm"
)

// TODO: Only check white list when processing purchase request
type NftPresaleWhitelist struct {
	gorm.Model
	WhitelistType WhitelistType `json:"whitelist_type" gorm:"index"`
	AssetId       string        `json:"asset_id" gorm:"type:varchar(255);index"`
	BatchGroupId  int           `json:"batch_group_id" gorm:"index"`
	UserId        int           `json:"user_id" gorm:"index"`
	Username      string        `json:"username" gorm:"type:varchar(255);index"`
}

type WhitelistType int

const (
	WhitelistTypeAsset WhitelistType = iota
	WhitelistTypeGroupBatch
)

type NftPresaleWhitelistSetRequest struct {
	WhitelistType WhitelistType `json:"whitelist_type" gorm:"index"`
	AssetId       string        `json:"asset_id" gorm:"type:varchar(255);index"`
	BatchGroupId  int           `json:"batch_group_id" gorm:"index"`
	UserId        int           `json:"user_id" gorm:"index"`
	Username      string        `json:"username" gorm:"type:varchar(255);index"`
}
