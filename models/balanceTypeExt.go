package models

type BalanceTypeExt struct {
	BalanceID uint               `json:"balance" gorm:"not null;"`
	Type      BalanceTypeExtList `json:"type" gorm:"not null"`
}

func (BalanceTypeExt) TableName() string {
	return "bill_balance_type_ext"
}

type BalanceTypeExtList uint

const (
	BTExtUnknown        BalanceTypeExtList = 0
	BTExtFirLaunch      BalanceTypeExtList = 6
	BTExtLocal          BalanceTypeExtList = 100
	BTExtBackFee        BalanceTypeExtList = 104
	BTExtOnChannel      BalanceTypeExtList = 200
	BTExtAward          BalanceTypeExtList = 300
	BTExtLocked         BalanceTypeExtList = 400
	BTExtLockedTransfer BalanceTypeExtList = 500
)
