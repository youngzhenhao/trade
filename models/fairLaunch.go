package models

import (
	"github.com/lightninglabs/taproot-assets/taprpc"
	"gorm.io/gorm"
)

type (
	FairLaunchState          int
	FairLaunchMintedState    int
	FairLaunchInventoryState int
	FairLaunchStatus         int
)

const (
	FairLaunchStateNoPay FairLaunchState = iota
	FairLaunchStatePaidPending
	FairLaunchStatePaidNoIssue
	FairLaunchStateIssuedPending
	FairLaunchStateIssued
	FairLaunchStateReservedSentPending
	FairLaunchStateReservedSent
)

const (
	FairLaunchMintedStateNoPay FairLaunchMintedState = iota
	FairLaunchMintedStatePaidPending
	FairLaunchMintedStatePaidNoSend
	FairLaunchMintedStateSentPending
	FairLaunchMintedStateSent
)

const (
	FairLaunchStateFail       FairLaunchState       = -1
	FairLaunchMintedStateFail FairLaunchMintedState = -1
)

const (
	FairLaunchInventoryStateOpen FairLaunchInventoryState = iota
	FairLaunchInventoryStateLocked
	FairLaunchInventoryStateMinted
)

const (
	MintMaxNumber = 10
)

const (
	StatusDeprecated FairLaunchStatus = iota
	StatusNormal
	StatusPending
	StatusUnknown
)

// FairLaunchInfo
// TODO: param FeeRate maybe need to rename
type FairLaunchInfo struct {
	gorm.Model
	ImageData                      string           `json:"image_data"`
	Name                           string           `json:"name" gorm:"type:varchar(255);not null"`
	AssetType                      taprpc.AssetType `json:"asset_type"`
	Amount                         int              `json:"amount"`
	Reserved                       int              `json:"reserved"`
	MintQuantity                   int              `json:"mint_quantity"`
	StartTime                      int              `json:"start_time"`
	EndTime                        int              `json:"end_time"`
	Description                    string           `json:"description"`
	FeeRate                        int              `json:"fee_rate"`
	SetGasFee                      int              `json:"set_gas_fee"`
	SetTime                        int              `json:"set_time"`
	ActualReserved                 float64          `json:"actual_reserved"`
	ReserveTotal                   int              `json:"reserve_total"`
	MintNumber                     int              `json:"mint_number"`
	IsFinalEnough                  bool             `json:"is_final_enough"`
	FinalQuantity                  int              `json:"final_quantity"`
	MintTotal                      int              `json:"mint_total"`
	ActualMintTotalPercent         float64          `json:"actual_mint_total_percent"`
	CalculationExpression          string           `json:"calculation_expression" gorm:"type:varchar(255)"`
	BatchKey                       string           `json:"batch_key" gorm:"type:varchar(255)"`
	BatchState                     string           `json:"batch_state" gorm:"type:varchar(255)"`
	BatchTxidAnchor                string           `json:"batch_txid_anchor" gorm:"type:varchar(255)"`
	AssetID                        string           `json:"asset_id" gorm:"type:varchar(255);index"`
	UserID                         int              `json:"user_id" gorm:"index"`
	Username                       string           `json:"username" gorm:"type:varchar(255)"`
	PayMethod                      FeePaymentMethod `json:"pay_method"`
	PaidSuccessTime                int              `json:"paid_success_time"`
	IssuanceFeePaidID              int              `json:"issuance_fee_paid_id"`
	IssuanceTime                   int              `json:"issuance_time"`
	ReservedCouldMint              bool             `json:"reserved_could_mint"`
	IsReservedSent                 bool             `json:"is_reserved_sent"`
	ReservedSentAnchorOutpointTxid string           `json:"reserved_sent_anchor_outpoint_txid" gorm:"type:varchar(255)"`
	ReservedSentAnchorOutpoint     string           `json:"reserved_sent_anchor_outpoint" gorm:"type:varchar(255)"`
	MintedNumber                   int              `json:"minted_number"`
	IsMintAll                      bool             `json:"is_mint_all"`
	Status                         FairLaunchStatus `json:"status" default:"1" gorm:"default:1;index"`
	State                          FairLaunchState  `json:"state"`
	ProcessNumber                  int              `json:"process_number"`
}

type SetFairLaunchInfoRequest struct {
	ImageData    string `json:"image_data"`
	Name         string `json:"name"`
	AssetType    int    `json:"asset_type"`
	Amount       int    `json:"amount"`
	Reserved     int    `json:"reserved"`
	MintQuantity int    `json:"mint_quantity"`
	StartTime    int    `json:"start_time"`
	EndTime      int    `json:"end_time"`
	Description  string `json:"description"`
	FeeRate      int    `json:"fee_rate"`
}

type FairLaunchMintedInfo struct {
	gorm.Model
	FairLaunchInfoID      int                   `json:"fair_launch_info_id" gorm:"not null;index"`
	MintedNumber          int                   `json:"minted_number"`
	MintedFeeRateSatPerKw int                   `json:"minted_fee_rate_sat_per_kw"`
	MintedGasFee          int                   `json:"minted_gas_fee"`
	EncodedAddr           string                `json:"encoded_addr" gorm:"type:varchar(512)"`
	MintFeePaidID         int                   `json:"mint_fee_paid_id"`
	PayMethod             FeePaymentMethod      `json:"pay_method"`
	PaidSuccessTime       int                   `json:"paid_success_time"`
	UserID                int                   `json:"user_id" gorm:"index"`
	Username              string                `json:"username" gorm:"type:varchar(255)"`
	AssetID               string                `json:"asset_id" gorm:"type:varchar(255)" gorm:"index"`
	AssetName             string                `json:"asset_name" gorm:"type:varchar(255)"`
	AssetType             int                   `json:"asset_type"`
	AddrAmount            int                   `json:"amount_addr"`
	ScriptKey             string                `json:"script_key" gorm:"type:varchar(255)"`
	InternalKey           string                `json:"internal_key" gorm:"type:varchar(255)"`
	TaprootOutputKey      string                `json:"taproot_output_key" gorm:"type:varchar(255)"`
	ProofCourierAddr      string                `json:"proof_courier_addr" gorm:"type:varchar(512)"`
	MintedSetTime         int                   `json:"minted_set_time"`
	SendAssetTime         int                   `json:"send_asset_time"`
	IsAddrSent            bool                  `json:"is_addr_sent"`
	OutpointTxHash        string                `json:"outpoint_tx_hash" gorm:"type:varchar(255)"`
	Outpoint              string                `json:"outpoint" gorm:"type:varchar(255)"`
	Address               string                `json:"address" gorm:"type:varchar(255)"`
	Status                FairLaunchStatus      `json:"status" gorm:"default:1"`
	State                 FairLaunchMintedState `json:"state"`
	ProcessNumber         int                   `json:"process_number"`
}

type MintFairLaunchRequest struct {
	FairLaunchInfoID      int    `json:"fair_launch_info_id"`
	MintedNumber          int    `json:"minted_number"`
	EncodedAddr           string `json:"encoded_addr" gorm:"type:varchar(255)"`
	MintedFeeRateSatPerKw int    `json:"minted_fee_rate_sat_per_kw"`
}

type MintFairLaunchReservedRequest struct {
	AssetID     string `json:"asset_id"`
	EncodedAddr string `json:"encoded_addr"`
}

type FairLaunchMintedUserInfo struct {
	gorm.Model
	UserID                 int              `json:"user_id" gorm:"not null"`
	FairLaunchMintedInfoID int              `json:"fair_launch_minted_info_id"`
	FairLaunchInfoID       int              `json:"fair_launch_info_id"`
	MintedNumber           int              `json:"minted_number"`
	Status                 FairLaunchStatus `json:"status" default:"1" gorm:"default:1"`
}

// Deprecated: Use FairLaunchMintedAndAvailableInfo instead
type FairLaunchInventoryInfo struct {
	gorm.Model
	FairLaunchInfoID       int                      `json:"fair_launch_info_id" gorm:"not null;index"`
	Quantity               int                      `json:"quantity" gorm:"index"`
	IsMinted               bool                     `json:"is_minted" gorm:"index"`
	FairLaunchMintedInfoID int                      `json:"fair_launch_minted_info_id" gorm:"index"`
	Status                 FairLaunchStatus         `json:"status" gorm:"default:1;index"`
	State                  FairLaunchInventoryState `json:"state" gorm:"index"`
}

type FairLaunchMintedAndAvailableInfo struct {
	gorm.Model
	FairLaunchInfoID      int    `json:"fair_launch_info_id" gorm:"not null;index"`
	MintedNumber          int    `json:"minted_number"`
	MintedAmount          int    `json:"minted_amount"`
	AvailableNumber       int    `json:"available_number"`
	AvailableAmount       int    `json:"available_amount"`
	ReserveTotal          int    `json:"reserve_total"`
	MintTotal             int    `json:"mint_total"`
	MintNumber            int    `json:"mint_number"`
	MintQuantity          int    `json:"mint_quantity"`
	FinalQuantity         int    `json:"final_quantity"`
	CalculationExpression string `json:"calculation_expression" gorm:"type:varchar(255)"`
}
