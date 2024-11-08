package custodyModels

import "gorm.io/gorm"

type BlockedRecord struct {
	gorm.Model
	UserId      uint        `gorm:"column:user_id;type:bigint;not null;index:idx_user_id;not null"`
	BlockedType BlockedType `gorm:"column:blocked_type;type:varchar(128);not null;index:idx_blocked_type;not null"`
	Memo        string      `gorm:"column:memo;type:varchar(255)" json:"memo"`
}

func (BlockedRecord) TableName() string {
	return "user_blocked_record"
}

type BlockedType string

const (
	BlockedUser   BlockedType = "blocked_user"
	UnblockedUser BlockedType = "unblocked_user"
)
