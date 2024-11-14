package models

type BalanceTypeExt struct {
	BalanceID uint `json:"balance" gorm:"not null;unique"`
	Type      uint `json:"type" gorm:"not null"`
}

func (BalanceTypeExt) TableName() string {
	return "bill_balance_type_ext"
}

type BalanceTypeExtList uint

const (
	BTExtUnknown   BalanceTypeExtList = 0
	BTExtLocal     BalanceTypeExtList = 100
	BTExtOutSide   BalanceTypeExtList = 200
	BTExtAward     BalanceTypeExtList = 300
	BTExtLocked    BalanceTypeExtList = 400
	BTExtLockedPay BalanceTypeExtList = 500
)
