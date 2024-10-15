package models

import "gorm.io/gorm"

type AccountAwardExt struct {
	gorm.Model
	BalanceId uint `gorm:"column:balance_id;type:bigint;default:0;uniqueIndex:uix_balance_award_id" json:"balance_id"`
	AwardId   uint `gorm:"column:award_id;type:bigint;default:0;uniqueIndex:uix_balance_award_id" json:"award_id"`
}

func (AccountAwardExt) TableName() string {
	return "user_account_award_ext"
}
