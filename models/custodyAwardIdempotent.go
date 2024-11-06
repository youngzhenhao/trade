package models

import "gorm.io/gorm"

type AccountAwardIdempotent struct {
	gorm.Model
	AwardId    uint   `gorm:"column:award_id;type:bigint;default:0;" json:"award_id"`
	Idempotent string `gorm:"column:idempotent;type:varchar(128);not null;unique" json:"idempotent"`
}

func (AccountAwardIdempotent) TableName() string {
	return "user_account_award_idempotent"
}
