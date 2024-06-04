package models

import (
	"gorm.io/gorm"
	"time"
)

type Invoice struct {
	gorm.Model
	UserID     uint       `gorm:"not null;column:user_id;type:bigint unsigned" json:"userId"`
	AccountID  *uint      `gorm:"column:account_id;type:bigint unsigned" json:"accountId"`
	AssetId    string     `gorm:"column:asset_id;default:00;varchar(100)" json:"assetId"`
	Invoice    string     `gorm:"column:invoice;type:varchar(512)" json:"invoice"`
	Amount     float64    `gorm:"type:decimal(10,2);column:amount" json:"amount"`
	CreateDate *time.Time `gorm:"column:create_date" json:"createDate"`
	Expiry     *int       `gorm:"column:expiry;type:bigint" json:"expiry"`
	Status     int16      `gorm:"column:status;type:smallint" json:"status"`
}

func (Invoice) TableName() string {
	return "user_invoice"
}
