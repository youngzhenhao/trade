package custodyModels

import "gorm.io/gorm"

type LimitLevel struct {
	gorm.Model
	LimitTypeId uint    `gorm:"column:limit_type_id;type:bigint unsigned;uniqueIndex:idx_limit_type_id_level";not null" json:"limitTypeId"`
	Level       uint    `gorm:"column:level;type:bigint unsigned;uniqueIndex:idx_limit_type_id_level;not null" json:"level"`
	Amount      float64 `gorm:"column:amount;type:decimal(15,2)" json:"amount"`
	Count       uint    `gorm:"column:count;type:bigint unsigned" json:"count"`
}

func (LimitLevel) TableName() string {
	return "user_limit_type_level"
}
