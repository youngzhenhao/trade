package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	UserId          uint          `gorm:"column:user_id;type:bigint unsigned;uniqueIndex:idx_user_id" json:"userId"` // column选项放在同一个gorm标签内
	UserName        string        `gorm:"column:user_name;type:varchar(100);index:idx_user_name" json:"userName"`
	UserAccountCode string        `gorm:"column:user_account_code;type:varchar(100)" json:"userAccountCode"`
	Status          AccountStatus `gorm:"column:status;type:smallint" json:"status"`
	Label           *string       `gorm:"column:label;type:varchar(100)" json:"label"`
}

func (Account) TableName() string {
	return "user_account"
}

type AccountStatus int16

const (
	AccountStatusDisable AccountStatus = 0
	AccountStatusEnable  AccountStatus = 1
)
