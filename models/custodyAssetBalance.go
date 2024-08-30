package models

import "gorm.io/gorm"

type AccountBalance struct {
	gorm.Model
	AccountID uint    `gorm:"column:account_id;type:bigint unsigned;uniqueIndex:idx_account_id_asset_id" json:"accountId"`
	AssetId   string  `gorm:"column:asset_id;varchar(100);uniqueIndex:idx_account_id_asset_id" json:"assetId"`
	Amount    float64 `gorm:"type:decimal(15,2);column:amount" json:"amount"`
}

func (AccountBalance) TableName() string {
	return "user_account_balance"
}
