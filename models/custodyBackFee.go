package models

import "gorm.io/gorm"

type BackFee struct {
	gorm.Model
	PayInsideId   uint         `gorm:"column:pay_inside_id;type:bigint;not null;index:idx_pay_inside_id;unique" json:"pay_inside_id"`
	BackBalanceId uint         `gorm:"column:back_balance_id;type:bigint" json:"back_balance_id"`
	Status        BackFeeState `gorm:"column:status;type:smallint;not null" json:"status"`
}

func (BackFee) TableName() string {
	return "user_back_fees"
}

type BackFeeState uint8

const (
	BackFeeStatePending BackFeeState = iota
	BackFeeStatePaid
)
