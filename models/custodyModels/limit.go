package custodyModels

import "gorm.io/gorm"

type Limit struct {
	gorm.Model
	UserId    uint `gorm:"column:user_id;type:bigint;not null;uniqueIndex:idx_user_id_limit_type;not null"`
	LimitType uint `gorm:"column:limit_type;type:bigint unsigned;not null;uniqueIndex:idx_user_id_limit_type;not null" json:"limitType"`
	Level     uint `gorm:"column:level;type:bigint unsigned;default:1" json:"level"`
}

func (Limit) TableName() string {
	return "user_limit"
}
