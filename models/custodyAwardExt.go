package models

import "gorm.io/gorm"

type AccountAwardExt struct {
	gorm.Model
	BalanceId   uint         `gorm:"column:balance_id;type:bigint;default:0;uniqueIndex:uix_balance_award_id" json:"balance_id"`
	AwardId     uint         `gorm:"column:award_id;type:bigint;default:0;uniqueIndex:uix_balance_award_id" json:"award_id"`
	AccountType AccountTypes `gorm:"column:account_type;type:tinyint;default:0" json:"account_type"`
}

func (AccountAwardExt) TableName() string {
	return "user_account_award_ext"
}

type AccountTypes uint

const (
	DefaultAccount AccountTypes = iota
	LockedAccount
)
