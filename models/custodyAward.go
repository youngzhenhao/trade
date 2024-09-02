package models

import "gorm.io/gorm"

type AccountAward struct {
	gorm.Model
	AccountID uint    `gorm:"column:account_id;type:bigint unsigned" json:"accountId"`
	AssetId   string  `gorm:"column:asset_id;type:varchar(128)" json:"assetId"`
	Amount    float64 `gorm:"column:amount;type:decimal(15,2)" json:"amount"`
	Memo      *string `gorm:"column:memo;type:varchar(255)" json:"memo"`
}

func (AccountAward) TableName() string {
	return "user_account_award"
}
