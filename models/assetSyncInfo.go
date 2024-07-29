package models

import (
	"gorm.io/gorm"
	"time"
)

type AssetSyncInfo struct {
	gorm.Model
	AssetId      string     `gorm:"column:asset_id;type:varchar(512);not null" json:"asset_Id"`
	Name         string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Point        string     `gorm:"column:point;type:varchar(255);not null" json:"point"`
	AssetType    AssetType  `gorm:"column:asset_type;type:smallint;not null" json:"assetType"`
	GroupName    *string    `gorm:"column:group_name;type:varchar(255)" json:"group_name"`
	GroupKey     *string    `gorm:"column:asset_is_group;type:varchar(255)" json:"group_key"`
	Amount       uint64     `gorm:"column:amount;type:bigint;not null" json:"amount"`
	Meta         *string    `gorm:"column:meta;type:text" json:"meta"`
	CreateHeight int64      `gorm:"column:create_height;type:bigint;not null" json:"create_height"`
	CreateTime   *time.Time `gorm:"column:create_time;not null" json:"create_time"`
	Universe     string     `gorm:"column:universe;type:varchar(100)" json:"universe"`
}

func (AssetSyncInfo) TableName() string {
	return "asset_sync_info"
}

type AssetType uint

const (
	AssettypeNormal AssetType = iota
	AssetTypeNFT
)

var (
	AssetType_name = map[AssetType]string{
		AssettypeNormal: "NORMAL",
		AssetTypeNFT:    "COLLECTIBLE",
	}
	AssetType_value = map[string]AssetType{
		"NORMAL":      AssettypeNormal,
		"COLLECTIBLE": AssetTypeNFT,
	}
)
