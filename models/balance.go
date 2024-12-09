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
	ServerFee   uint64       `gorm:"column:server_fee;type:bigint unsigned" json:"serverFee"`
	AssetId     *string      `gorm:"column:asset_id;type:varchar(512);default:'00'" json:"assetId"`
	Invoice     *string      `gorm:"column:invoice;type:varchar(512)" json:"invoice"`
	PaymentHash *string      `gorm:"column:payment_hash;type:varchar(100)" json:"paymentHash"`
	State       BalanceState `gorm:"column:State;type:smallint" json:"State"`
	TypeExt     *BalanceTypeExt
}

func (Balance) TableName() string {
	return "bill_balance"
}

type BalanceType int16

const (
	BillTypeRecharge        BalanceType = 0
	BillTypePayment         BalanceType = 1
	BillTypeAssetTransfer   BalanceType = 2
	BillTypeAssetMintedSend             = 3
	BILL_TYPE_BACK_FEE                  = 4
	BillTypeAwardSat                    = 5
	BillTypeAwardAsset                  = 6
	BiLLTypeLock                        = 7
	BillTypePendingOder                 = 8
	BillTypePoolAccount                 = 9

	//locked 仅作为查询时的标识，不参与任何业务逻辑
	LockedTransfer BalanceType = 1000
)

func (bt BalanceType) String() string {
	balanceTypeMapString := map[BalanceType]string{
		BillTypeRecharge:        "BillTypeRecharge",
		BillTypePayment:         "BillTypePayment",
		BillTypeAssetTransfer:   "BillTypeAssetTransfer",
		BillTypeAssetMintedSend: "BillTypeAssetMintedSend",
		BILL_TYPE_BACK_FEE:      "BILL_TYPE_BACK_FEE",
		BillTypeAwardSat:        "BillTypeAwardSat",
		BillTypeAwardAsset:      "BillTypeAwardAsset",
	}
	return balanceTypeMapString[bt]
}

type BalanceAway int16

const (
	AWAY_IN  BalanceAway = 0
	AWAY_OUT BalanceAway = 1
)

func (ba BalanceAway) String() string {
	balanceAwayMapString := map[BalanceAway]string{
		AWAY_IN:  "AWAY_IN",
		AWAY_OUT: "AWAY_OUT",
	}
	return balanceAwayMapString[ba]
}

type BalanceUnit int16

const (
	UNIT_SATOSHIS          BalanceUnit = 0
	UNIT_ASSET_NORMAL      BalanceUnit = 1
	UNIT_ASSET_COLLECTIBLE BalanceUnit = 2
)

func (bu BalanceUnit) String() string {
	balanceAwayMapString := map[BalanceUnit]string{
		UNIT_SATOSHIS:          "UNIT_SATOSHIS",
		UNIT_ASSET_NORMAL:      "UNIT_ASSET_NORMAL",
		UNIT_ASSET_COLLECTIBLE: "UNIT_ASSET_COLLECTIBLE",
	}
	return balanceAwayMapString[bu]
}

type BalanceState int16

const (
	STATE_UNKNOW  BalanceState = 0
	STATE_SUCCESS BalanceState = 1
	STATE_FAILED  BalanceState = 2
)

func (bs BalanceState) String() string {
	balanceStateMapString := map[BalanceState]string{
		STATE_UNKNOW:  "STATE_UNKNOW",
		STATE_SUCCESS: "STATE_SUCCESS",
		STATE_FAILED:  "STATE_FAILED",
	}
	return balanceStateMapString[bs]
}
