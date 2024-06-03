package models

import (
	"github.com/lightninglabs/taproot-assets/taprpc"
	"gorm.io/gorm"
)

type (
	FairLaunchState          int
	FairLaunchMintedState    int
	FairLaunchInventoryState int
)

var (
	MintMaxNumber                                             = 10
	StatusDeprecated                                          = 0
	StatusNormal                                              = 1
	StatusPending                                             = 2
	StatusUnknown                                             = 4
	FairLaunchStateNoPay             FairLaunchState          = 0
	FairLaunchStatePaidPending       FairLaunchState          = 1
	FairLaunchStatePaidNoIssue       FairLaunchState          = 2
	FairLaunchStateIssuedPending     FairLaunchState          = 3
	FairLaunchStateIssued            FairLaunchState          = 4
	FairLaunchMintedStateNoPay       FairLaunchMintedState    = 0
	FairLaunchMintedStatePaidPending FairLaunchMintedState    = 1
	FairLaunchMintedStatePaidNoSend  FairLaunchMintedState    = 2
	FairLaunchMintedStateSentPending FairLaunchMintedState    = 3
	FairLaunchMintedStateSent        FairLaunchMintedState    = 4
	FairLaunchInventoryStateOpen     FairLaunchInventoryState = 0
	FairLaunchInventoryStateLocked   FairLaunchInventoryState = 1
	FairLaunchInventoryStateMinted   FairLaunchInventoryState = 2
)

type FairLaunchInfo struct {
	gorm.Model
	ImageData              string           `json:"image_data"`
	Name                   string           `json:"name" gorm:"type:varchar(255);not null"`
	AssetType              taprpc.AssetType `json:"asset_type"`
	Amount                 int              `json:"amount"`
	Reserved               int              `json:"reserved"`
	MintQuantity           int              `json:"mint_quantity"`
	StartTime              int              `json:"start_time"`
	EndTime                int              `json:"end_time"`
	Description            string           `json:"description"`
	FeeRate                int              `json:"fee_rate"`
	SetTime                int              `json:"set_time"`
	ActualReserved         float64          `json:"actual_reserved"`
	ReserveTotal           int              `json:"reserve_total"`
	MintNumber             int              `json:"mint_number"`
	IsFinalEnough          bool             `json:"is_final_enough"`
	FinalQuantity          int              `json:"final_quantity"`
	MintTotal              int              `json:"mint_total"`
	ActualMintTotalPercent float64          `json:"actual_mint_total_percent"`
	CalculationExpression  string           `json:"calculation_expression" gorm:"type:varchar(255)"`
	BatchKey               string           `json:"batch_key" gorm:"type:varchar(255)"`
	BatchState             string           `json:"batch_state" gorm:"type:varchar(255)"`
	BatchTxidAnchor        string           `json:"batch_txid_anchor" gorm:"type:varchar(255)"`
	AssetID                string           `json:"asset_id" gorm:"type:varchar(255)"`
	UserID                 int              `json:"user_id"`
	PayMethod              FeePaymentMethod `json:"pay_method"`
	PaidSuccessTime        int              `json:"paid_success_time"`
	IssuanceFeePaidID      int              `json:"issuance_fee_paid_id"`
	IssuanceTime           int              `json:"issuance_time"`
	ReservedCouldMint      bool             `json:"reserved_could_mint"`
	IsReservedSent         bool             `json:"is_reserved_sent"`
	MintedNumber           int              `json:"minted_number"`
	IsMintAll              bool             `json:"is_mint_all"`
	Status                 int              `json:"status" default:"1" gorm:"default:1"`
	State                  FairLaunchState  `json:"state"`
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
	FairLaunchInfoID      int                   `json:"fair_launch_info_id" gorm:"not null"`
	MintedNumber          int                   `json:"minted_number"`
	MintedFeeRateSatPerKw int                   `json:"minted_fee_rate_sat_per_kw"`
	MintedGasFee          int                   `json:"minted_gas_fee"`
	EncodedAddr           string                `json:"encoded_addr" gorm:"type:varchar(512)"`
	MintFeePaidID         int                   `json:"mint_fee_paid_id"`
	PayMethod             FeePaymentMethod      `json:"pay_method"`
	PaidSuccessTime       int                   `json:"paid_success_time"`
	UserID                int                   `json:"user_id"`
	AssetID               string                `json:"asset_id" gorm:"type:varchar(255)"`
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
	Status                int                   `json:"status" gorm:"default:1"`
	State                 FairLaunchMintedState `json:"state"`
}

type MintFairLaunchRequest struct {
	FairLaunchInfoID      int    `json:"fair_launch_info_id"`
	MintedNumber          int    `json:"minted_number"`
	EncodedAddr           string `json:"encoded_addr" gorm:"type:varchar(255)"`
	MintedFeeRateSatPerKw int    `json:"minted_fee_rate_sat_per_kw"`
}

type FairLaunchMintedUserInfo struct {
	gorm.Model
	UserID                 int `json:"user_id" gorm:"not null"`
	FairLaunchMintedInfoID int `json:"fair_launch_minted_info_id"`
	FairLaunchInfoID       int `json:"fair_launch_info_id"`
	MintedNumber           int `json:"minted_number"`
	Status                 int `json:"status" default:"1" gorm:"default:1"`
}

type FairLaunchInventoryInfo struct {
	gorm.Model
	FairLaunchInfoID       int                      `json:"fair_launch_info_id" gorm:"not null"`
	Quantity               int                      `json:"quantity"`
	IsMinted               bool                     `json:"is_minted"`
	FairLaunchMintedInfoID int                      `json:"fair_launch_minted_info_id"`
	Status                 int                      `json:"status" gorm:"default:1"`
	State                  FairLaunchInventoryState `json:"state"`
}
