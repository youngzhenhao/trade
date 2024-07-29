package models

import (
	"gorm.io/gorm"
)

type AssetLocalMint struct {
	gorm.Model
	AssetVersion    string `json:"asset_version" gorm:"type:varchar(255)"`
	AssetType       string `json:"asset_type" gorm:"type:varchar(255)"`
	Name            string `json:"name" gorm:"type:varchar(255)"`
	AssetMetaData   string `json:"asset_meta_data" gorm:"type:varchar(255)"`
	AssetMetaType   string `json:"asset_meta_type" gorm:"type:varchar(255)"`
	AssetMetaHash   string `json:"asset_meta_hash" gorm:"type:varchar(255)"`
	Amount          int    `json:"amount"`
	NewGroupedAsset bool   `json:"new_grouped_asset"`
	GroupKey        string `json:"group_key" gorm:"type:varchar(255)"`
	GroupAnchor     string `json:"group_anchor" gorm:"type:varchar(255)"`
	GroupedAsset    bool   `json:"grouped_asset"`
	BatchKey        string `json:"batch_key" gorm:"type:varchar(255)"`
	BatchTxid       string `json:"batch_txid" gorm:"type:varchar(255)"`
	AssetId         string `json:"asset_id" gorm:"type:varchar(255)"`
	DeviceId        string `json:"device_id" gorm:"type:varchar(255)"`
	UserId          int    `json:"user_id"`
	Username        string `json:"username" gorm:"type:varchar(255)"`
	Status          int    `json:"status" gorm:"default:1"`
}

type AssetLocalMintSetRequest struct {
	AssetVersion    string `json:"asset_version" gorm:"type:varchar(255)"`
	AssetType       string `json:"asset_type" gorm:"type:varchar(255)"`
	Name            string `json:"name" gorm:"type:varchar(255)"`
	AssetMetaData   string `json:"asset_meta_data" gorm:"type:varchar(255)"`
	AssetMetaType   string `json:"asset_meta_type" gorm:"type:varchar(255)"`
	AssetMetaHash   string `json:"asset_meta_hash" gorm:"type:varchar(255)"`
	Amount          int    `json:"amount"`
	NewGroupedAsset bool   `json:"new_grouped_asset"`
	GroupKey        string `json:"group_key" gorm:"type:varchar(255)"`
	GroupAnchor     string `json:"group_anchor" gorm:"type:varchar(255)"`
	GroupedAsset    bool   `json:"grouped_asset"`
	BatchKey        string `json:"batch_key" gorm:"type:varchar(255)"`
	BatchTxid       string `json:"batch_txid" gorm:"type:varchar(255)"`
	AssetId         string `json:"asset_id" gorm:"type:varchar(255)"`
	DeviceId        string `json:"device_id" gorm:"type:varchar(255)"`
}
