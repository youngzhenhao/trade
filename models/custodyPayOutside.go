package models

import "gorm.io/gorm"

type PayOutside struct {
	gorm.Model
	AccountID uint             `gorm:"column:account_id;type:bigint unsigned;index:idx_account_id" json:"accountId"`
	AssetId   string           `gorm:"column:asset_id;default:00;varchar(100)" json:"assetId"`
	Address   string           `gorm:"column:address;type:varchar(512)" json:"address"`
	Amount    float64          `gorm:"type:decimal(15,2);column:amount" json:"amount"`
	TxHash    string           `gorm:"column:tx_hash;type:varchar(100)" json:"txHash"`
	BalanceId uint             `gorm:"column:balance_id;type:bigint;default:0" json:"balance_id"`
	Status    PayOutsideStatus `gorm:"column:status;type:smallint" json:"status"`
}

func (PayOutside) TableName() string {
	return "user_outside"
}

type PayOutsideStatus int16

const (
	PayOutsideStatusPending PayOutsideStatus = 0
	PayOutsideStatusPaid    PayOutsideStatus = 1
	PayOutsideStatusSuccess PayOutsideStatus = 2
)
