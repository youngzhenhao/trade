package custodyModels

import "gorm.io/gorm"

type AccountBtcBalance struct {
	gorm.Model
	AccountId uint    `gorm:"column:account_id;type:bigint unsigned;uniqueIndex:idx_account_id" json:"accountId"`
	Amount    float64 `gorm:"type:decimal(15,2);column:amount" json:"amount"`
}

func (AccountBtcBalance) TableName() string {
	return "user_account_balance_btc"
}
