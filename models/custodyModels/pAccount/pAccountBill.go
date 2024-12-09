package pAccount

import "gorm.io/gorm"

type PAccountBill struct {
	gorm.Model
	PoolAccountId uint             `gorm:"index;column:pool_account_id;not null"`
	Away          PAccountBillAway `gorm:"column:away;type:smallint" json:"away"`
	Target        string           `gorm:"column:target;type:varchar(100)" json:"target"`
	Amount        float64          `gorm:"column:amount;type:decimal(15,2)" json:"amount"`
	AssetId       string           `gorm:"column:asset_id;varchar(128);default:'00'" json:"assetId"`
	PaymentHash   string           `gorm:"column:payment_hash;type:varchar(100)" json:"paymentHash"`
	State         PAccountState    `gorm:"column:State;type:smallint" json:"State"`

	// 外键关联
	PoolAccount *PoolAccount `gorm:"foreignkey:PoolAccountId"`
}

func (PAccountBill) TableName() string {
	return "custody_pool_account_bills"
}

type PAccountBillAway uint8

const (
	PAccountBillAwayIn PAccountBillAway = iota
	PAccountBillAwayOut
)

type PAccountState uint8
