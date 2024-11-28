package custodyModels

import "gorm.io/gorm"

type AccountInsideMission struct {
	gorm.Model
	AccountId uint    `gorm:"column:account_id;type:bigint unsigned;" json:"accountId"`
	AssetId   string  `gorm:"column:asset_id;type:varchar(128);" json:"assetId"`
	Type      AIMType `gorm:"type:enum('btc','asset');column:type;index:idx_type" json:"type"`

	ReceiverId uint    `gorm:"column:receiver_id;type:bigint unsigned" json:"receiverId"`
	InvoiceId  uint    `gorm:"column:invoice_id;type:bigint unsigned" json:"invoiceId"`
	Amount     float64 `gorm:"type:decimal(15,2);column:amount" json:"amount"`
	Fee        float64 `gorm:"type:decimal(15,2);column:fee" json:"fee"`
	FeeType    string  `gorm:"type:varchar(128);column:fee_type" json:"feeType"`

	PayerBalanceId    uint `gorm:"column:balance_id;type:bigint unsigned" json:"balanceId"`
	ReceiverBalanceId uint `gorm:"column:receiver_balance_id;type:bigint unsigned" json:"receiverBalanceId"`

	Retries int      `gorm:"type:int;column:retries" json:"retries"`
	Error   string   `gorm:"type:text;column:error" json:"error"`
	State   AIMState `gorm:"type:tinyint unsigned;column:state" json:"state"`
}

func (AccountInsideMission) TableName() string {
	return "user_account_inside_mission"
}

type AIMType string

const (
	AIMTypeBtc   AIMType = "btc"
	AIMTypeAsset AIMType = "asset"
)

type AIMState int8

const (
	AIMStateDone    AIMState = -1
	AIMStatePending AIMState = 0
	AIMStatePaid    AIMState = 3
	AIMStateSuccess AIMState = 5
)
