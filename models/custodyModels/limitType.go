package custodyModels

import "gorm.io/gorm"

type LimitType struct {
	gorm.Model
	AssetId      string            `gorm:"column:asset_id;type:varchar(128);uniqueIndex:idx_asset_id_transfer_type" json:"assetId"`
	TransferType LimitTransferType `gorm:"column:transfer_type;type:bigint unsigned;uniqueIndex:idx_asset_id_transfer_type" json:"transferType"`
	Memo         string            `gorm:"column:memo;type:varchar(128)" json:"memo"`
}

func (LimitType) TableName() string {
	return "user_limit_type"
}

type LimitTransferType uint

const (
	LimitTransferTypeLocal   LimitTransferType = 0
	LimitTransferTypeOutside LimitTransferType = 1
)
