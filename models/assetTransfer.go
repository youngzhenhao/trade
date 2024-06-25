package models

import (
	"github.com/lightninglabs/taproot-assets/taprpc"
	"gorm.io/gorm"
)

type AssetTransferType int

const (
	AssetTransferTypeOut = iota
	AssetTransferTypeIn
)

type AssetTransfer struct {
	gorm.Model
	AssetID           string                `json:"asset_id" gorm:"type:varchar(255)"`
	AssetName         string                `json:"asset_name" gorm:"type:varchar(255)"`
	AssetType         taprpc.AssetType      `json:"asset_type"`
	AssetAddressFrom  string                `json:"address_from" gorm:"type:varchar(255)"`
	AssetAddressTo    string                `json:"address_to" gorm:"type:varchar(255)"`
	Amount            int                   `json:"amount"`
	TransferType      AssetTransferType     `json:"transfer_type"`
	Inputs            []AssetTransferInput  `json:"inputs"`
	Outputs           []AssetTransferOutput `json:"outputs"`
	UserID            int                   `json:"user_id"`
	TransactionID     string                `json:"transaction_id" gorm:"type:varchar(255)"`
	TransferTimestamp int                   `json:"transfer_timestamp"`
	AnchorTxChainFees int                   `json:"anchor_tx_chain_fees"`
	ConfirmedBlocks   int                   `json:"confirmed_blocks"`
	Status            int                   `json:"status" gorm:"default:1"`
}

type AssetTransferInput struct {
	AnchorPoint string `json:"anchor_point"`
	ScriptKey   string `json:"script_key"`
	Amount      string `json:"amount"`
}

type AssetTransferOutputAnchor struct {
	Outpoint string `json:"outpoint"`
	Value    string `json:"value"`
}

type AssetTransferOutput struct {
	Anchor           AssetTransferOutputAnchor
	ScriptKey        string `json:"script_key"`
	ScriptKeyIsLocal bool   `json:"script_key_is_local"`
	Amount           string `json:"amount"`
}

type AssetTransaction struct {
	gorm.Model
	AssetID         string                 `json:"asset_id" gorm:"type:varchar(255)"`
	AssetName       string                 `json:"asset_name" gorm:"type:varchar(255)"`
	AssetType       taprpc.AssetType       `json:"asset_type"`
	Input           []AssetTransactionItem `json:"input"`
	Output          []AssetTransactionItem `json:"output"`
	TransactionID   string                 `json:"transaction_id" gorm:"type:varchar(255)"`
	Time            int                    `json:"time"`
	FeeRate         int                    `json:"fee_rate"`
	Fee             int                    `json:"fee"`
	ConfirmedBlocks int                    `json:"confirmed_blocks"`
	Status          int                    `json:"status" gorm:"default:1"`
}

type AssetTransactionItem struct {
	Address string `json:"address" gorm:"type:varchar(255)"`
	Value   int    `json:"value"`
}

type AssetTransferSetRequest struct {
	AssetID           string                `json:"asset_id" gorm:"type:varchar(255)"`
	AssetAddressFrom  string                `json:"address_from" gorm:"type:varchar(255)"`
	AssetAddressTo    string                `json:"address_to" gorm:"type:varchar(255)"`
	Amount            int                   `json:"amount"`
	TransferType      AssetTransferType     `json:"transfer_type"`
	Inputs            []AssetTransferInput  `json:"inputs"`
	Outputs           []AssetTransferOutput `json:"outputs"`
	TransactionID     string                `json:"transaction_id" gorm:"type:varchar(255)"`
	TransferTimestamp int                   `json:"transfer_timestamp"`
	AnchorTxChainFees int                   `json:"anchor_tx_chain_fees"`
}
