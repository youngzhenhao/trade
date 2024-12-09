package pAccount

import "gorm.io/gorm"

type PoolAccount struct {
	gorm.Model
	PairId uint `gorm:"column:pair_id;unique;index;not null"`
	Status uint `gorm:"column:status"`
}

func (PoolAccount) TableName() string {
	return "custody_pool_accounts"
}
