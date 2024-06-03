package models

import "gorm.io/gorm"

type BalanceExt struct {
	BalanceId   uint    `gorm:"column:balance_id;type:bigint unsigned" json:"balanceId"`
	BillExtDesc *string `gorm:"column:bill_ext_desc;type:longtext" json:"billExtDesc"`
	gorm.Model
}

func (BalanceExt) TableName() string {
	return "bill_balance_ext"
}
