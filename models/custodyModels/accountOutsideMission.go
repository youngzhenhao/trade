package custodyModels

import (
	"gorm.io/gorm"
)

type AccountOutsideMission struct {
	gorm.Model
	AccountId uint     `gorm:"column:account_id;type:bigint unsigned;" json:"accountId"`
	AssetId   string   `gorm:"column:asset_id;type:varchar(128);" json:"assetId"`
	Type      AOMType  `gorm:"type:enum('btc','asset');column:type;index:idx_type" json:"type"`
	Target    string   `gorm:"type:text;column:target" json:"target"`
	Hash      string   `gorm:"type:varchar(128);column:hash" json:"hash"`
	Amount    float64  `gorm:"type:decimal(15,2);column:amount" json:"amount"`
	FeeLimit  float64  `gorm:"type:decimal(15,2);column:fee_limit" json:"feeLimit"`
	Fee       float64  `gorm:"type:decimal(15,2);column:fee" json:"fee"`
	FeeType   string   `gorm:"type:varchar(128);column:fee_type" json:"feeType"`
	BalanceId uint     `gorm:"column:balance_id;type:bigint unsigned" json:"balanceId"`
	TxId      string   `gorm:"type:varchar(128);column:tx_id" json:"txId"`
	Retries   int      `gorm:"type:int;column:retries" json:"retries"`
	Error     string   `gorm:"type:text;column:error" json:"error"`
	State     AOMState `gorm:"type:int8;column:state" json:"state"`
}

func (AccountOutsideMission) TableName() string {
	return "user_account_outside_mission"
}

type AOMType string

const (
	AOMTypeBtc   AOMType = "btc"
	AOMTypeAsset AOMType = "asset"
)

type AOMState int8

const (
	AOMStateDone      AOMState = -1
	AOMStatePending   AOMState = 0
	AOMStateNotPayFee AOMState = 3
	AOMStateSuccess   AOMState = 5
)
