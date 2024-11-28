package models

type BalanceTypeExt struct {
	BalanceID uint               `json:"balance" gorm:"not null;unique;index;"`
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

func (b BalanceTypeExtList) ToString() string {
	balanceTypeExtString := map[BalanceTypeExtList]string{
		BTExtUnknown:        "Unknown",
		BTExtFirLaunch:      "FirLaunch",
		BTExtLocal:          "Local",
		BTExtBackFee:        "BackFee",
		BTExtOnChannel:      "OnChannel",
		BTExtAward:          "Award",
		BTExtLocked:         "Locked",
		BTExtLockedTransfer: "LockedTransfer",
	}
	return balanceTypeExtString[b]
}
func ToBalanceTypeExtList(s string) BalanceTypeExtList {
	balanceTypeExtList := map[string]BalanceTypeExtList{
		"Unknown":        BTExtUnknown,
		"FirLaunch":      BTExtFirLaunch,
		"Local":          BTExtLocal,
		"BackFee":        BTExtBackFee,
		"OnChannel":      BTExtOnChannel,
		"Award":          BTExtAward,
		"Locked":         BTExtLocked,
		"LockedTransfer": BTExtLockedTransfer,
	}
	return balanceTypeExtList[s]
}
