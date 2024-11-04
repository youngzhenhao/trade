package custodyModels

import (
	"gorm.io/gorm"
	"time"
)

type LimitBill struct {
	gorm.Model
	UserId        uint      `gorm:"column:user_id;type:bigint;not null;index:idx_user_id_limit_type;not null"`
	LimitType     uint      `gorm:"column:limit_type;type:bigint unsigned;index:idx_user_id_limit_type;not null"`
	TotalAmount   float64   `gorm:"column:total_amount;type:decimal(15,2)" json:"totalAmount"`
	UseAbleAmount float64   `gorm:"column:use_able_amount;type:decimal(15,2)" json:"useAbleAmount"`
	TotalCount    uint      `gorm:"column:total_count;type:bigint unsigned" json:"totalCount"`
	UseAbleCount  uint      `gorm:"column:use_able_count;type:bigint unsigned" json:"useAbleCount"`
	LocalTime     time.Time `gorm:"column:local_time;type:datetime;not null" json:"localTime"`
}

func (LimitBill) TableName() string {
	return "user_limit_bills"
}
