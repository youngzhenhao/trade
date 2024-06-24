package models

import (
	"github.com/lightninglabs/taproot-assets/taprpc"
	"gorm.io/gorm"
)

type (
	IdoPublishState     int
	IdoParticipateState int
	IdoStatus           int
)

const (
	IdoPublishStateNoPay IdoPublishState = iota
	IdoPublishStatePaidPending
	IdoPublishStatePaidNoPublish
	IdoPublishStatePublishedPending
	IdoPublishStatePublished
	IdoPublishStateRefundedPending
	IdoPublishStateRefunded
)

const (
	IdoParticipateStateNoPay IdoParticipateState = iota
	IdoParticipateStatePaidPending
	IdoParticipateStatePaidNoSend
	IdoParticipateStateSentPending
	IdoParticipateStateSent
)

const (
	IdoStatusDeprecated IdoStatus = iota
	IdoStatusNormal
	IdoStatusPending
	IdoStatusUnknown
)

type IdoPublishInfo struct {
	gorm.Model
	AssetID           string           `json:"asset_id" gorm:"type:varchar(255)"`
	AssetName         string           `json:"asset_name" gorm:"type:varchar(255)"`
	AssetType         taprpc.AssetType `json:"asset_type"`
	TotalAmount       int              `json:"total_amount"`
	MinimumQuantity   int              `json:"min_quantity"`
	UnitPrice         int              `json:"unit_price"`
	StartTime         int              `json:"start_time"`
	EndTime           int              `json:"end_time"`
	FeeRate           int              `json:"fee_rate"`
	GasFee            int              `json:"gas_fee"`
	SetTime           int              `json:"set_time"`
	UserID            int              `json:"user_id"`
	PayMethod         FeePaymentMethod `json:"pay_method"`
	FeePaidID         int              `json:"fee_paid_id"`
	PaidSuccessTime   int              `json:"paid_success_time"`
	EncodedAddr       string           `json:"encoded_addr" gorm:"type:varchar(512)"`
	ScriptKey         string           `json:"script_key" gorm:"type:varchar(255)"`
	InternalKey       string           `json:"internal_key" gorm:"type:varchar(255)"`
	TaprootOutputKey  string           `json:"taproot_output_key" gorm:"type:varchar(255)"`
	ProofCourierAddr  string           `json:"proof_courier_addr" gorm:"type:varchar(512)"`
	ParticipateAmount int              `json:"participate_amount"`
	IsParticipateAll  bool             `json:"is_participate_all"`
	Status            IdoStatus        `json:"status" default:"1" gorm:"default:1"`
	State             IdoPublishState  `json:"state"`
	ProcessNumber     int              `json:"process_number"`
}

type IdoParticipateInfo struct {
	gorm.Model
	IdoPublishInfoID int                 `json:"ido_publish_info_id"`
	AssetID          string              `json:"asset_id" gorm:"type:varchar(255)"`
	AssetName        string              `json:"asset_name" gorm:"type:varchar(255)"`
	AssetType        taprpc.AssetType    `json:"asset_type"`
	BoughtAmount     int                 `json:"bought_amount"`
	FeeRate          int                 `json:"fee_rate"`
	GasFee           int                 `json:"gas_fee"`
	SetTime          int                 `json:"set_time"`
	UserID           int                 `json:"user_id"`
	PayMethod        FeePaymentMethod    `json:"pay_method"`
	FeePaidID        int                 `json:"fee_paid_id"`
	PaidSuccessTime  int                 `json:"paid_success_time"`
	EncodedAddr      string              `json:"encoded_addr" gorm:"type:varchar(512)"`
	ScriptKey        string              `json:"script_key" gorm:"type:varchar(255)"`
	InternalKey      string              `json:"internal_key" gorm:"type:varchar(255)"`
	TaprootOutputKey string              `json:"taproot_output_key" gorm:"type:varchar(255)"`
	ProofCourierAddr string              `json:"proof_courier_addr" gorm:"type:varchar(512)"`
	SendAssetTime    int                 `json:"send_asset_time"`
	IsAddrSent       bool                `json:"is_addr_sent"`
	OutpointTxHash   string              `json:"outpoint_tx_hash" gorm:"type:varchar(255)"`
	Outpoint         string              `json:"outpoint" gorm:"type:varchar(255)"`
	Address          string              `json:"address" gorm:"type:varchar(255)"`
	Status           IdoStatus           `json:"status" gorm:"default:1"`
	State            IdoParticipateState `json:"state"`
	ProcessNumber    int                 `json:"process_number"`
}

type PublishIdoRequest struct {
	gorm.Model
	AssetID         string `json:"asset_id" gorm:"type:varchar(255)"`
	TotalAmount     int    `json:"total_amount"`
	MinimumQuantity int    `json:"min_quantity"`
	UnitPrice       int    `json:"unit_price"`
	StartTime       int    `json:"start_time"`
	EndTime         int    `json:"end_time"`
	FeeRate         int    `json:"fee_rate"`
}

type ParticipateIdoRequest struct {
	gorm.Model
	IdoPublishInfoID int    `json:"ido_publish_info_id"`
	BoughtAmount     int    `json:"bought_amount"`
	FeeRate          int    `json:"fee_rate"`
	EncodedAddr      string `json:"encoded_addr" gorm:"type:varchar(512)"`
}

type IdoParticipateUserInfo struct {
	gorm.Model
	UserID               int       `json:"user_id" gorm:"not null"`
	IdoParticipateInfoID int       `json:"ido_participate_info_id"`
	IdoPublishInfoID     int       `json:"ido_publish_info_id"`
	AssetID              string    `json:"asset_id"`
	BoughtAmount         int       `json:"bought_amount"`
	Status               IdoStatus `json:"status" gorm:"default:1"`
}
