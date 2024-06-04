package models

import "gorm.io/gorm"

type BalanceExt struct {
	gorm.Model
	BalanceId   uint    `gorm:"column:balance_id;type:bigint unsigned" json:"balanceId"`
	BillExtDesc *string `gorm:"column:bill_ext_desc;type:longtext" json:"billExtDesc"`
}

func (BalanceExt) TableName() string {
	return "bill_balance_ext"
}
