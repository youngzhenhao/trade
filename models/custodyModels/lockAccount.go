package custodyModels

import "gorm.io/gorm"

type LockAccount struct {
	gorm.Model
	UserId   uint          `gorm:"column:user_id;type:bigint unsigned;uniqueIndex:idx_user_id" json:"userId"` // column选项放在同一个gorm标签内
	UserName string        `gorm:"column:user_name;type:varchar(100);index:idx_user_name" json:"userName"`
	Status   AccountStatus `gorm:"column:status;default:1;type:smallint" json:"status"`
}

func (LockAccount) TableName() string {
	return "user_lock_account"
}

type AccountStatus int8

const (
	AccountStatusDisable AccountStatus = 0
	AccountStatusEnable  AccountStatus = 1
)
