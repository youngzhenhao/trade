package models

import (
	"gorm.io/gorm"
)

type Balance struct {
	gorm.Model
	AccountId   uint         `gorm:"column:account_id;type:bigint unsigned" json:"accountId"` // 正确地将unique和column选项放在同一个gorm标签内
	BillType    BalanceType  `gorm:"column:bill_type;type:smallint" json:"billType"`
	Away        BalanceAway  `gorm:"column:away;type:smallint" json:"away"`
	Amount      float64      `gorm:"column:amount;type:decimal(10,2)" json:"amount"`
	Unit        BalanceUnit  `gorm:"column:Unit;type:smallint" json:"unit"`
	Invoice     *string      `gorm:"column:invoice;type:varchar(512)" json:"invoice"`
	PaymentHash *string      `gorm:"column:payment_hash;type:varchar(100)" json:"paymentHash"`
	State       BalanceState `gorm:"column:State;type:smallint" json:"State"`
}

func (Balance) TableName() string {
	return "bill_balance"
}

type BalanceType int16

const (
	BILL_TYPE_RECHARGE          BalanceType = 0
	BILL_TYPE_PAYMENT           BalanceType = 1
	BILL_TYPE_ASSET_TRANSFER    BalanceType = 2
	BILL_TYPE_ASSET_MINTED_SEND             = 3
)

type BalanceAway int16

const (
	AWAY_IN  BalanceAway = 0
	AWAY_OUT BalanceAway = 1
)

type BalanceUnit int16

const (
	UNIT_SATOSHIS          BalanceUnit = 0
	UNIT_ASSET_NORMAL      BalanceUnit = 1
	UNIT_ASSET_COLLECTIBLE BalanceUnit = 2
)

type BalanceState int16

const (
	STATE_UNKNOW  BalanceState = 0
	STATE_SUCCESS BalanceState = 1
	STATE_FAILED  BalanceState = 2
)
