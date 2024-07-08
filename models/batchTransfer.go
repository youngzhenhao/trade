package models

import "gorm.io/gorm"

type BatchTransfer struct {
	gorm.Model
	Encoded            string `json:"encoded"`
	AssetID            string `json:"asset_id" gorm:"type:varchar(255)"`
	Amount             int    `json:"amount"`
	ScriptKey          string `json:"script_key" gorm:"type:varchar(255)"`
	InternalKey        string `json:"internal_key" gorm:"type:varchar(255)"`
	TaprootOutputKey   string `json:"taproot_output_key" gorm:"type:varchar(255)"`
	ProofCourierAddr   string `json:"proof_courier_addr" gorm:"type:varchar(255)"`
	Txid               string `json:"txid" gorm:"type:varchar(255)"`
	Index              int    `json:"index"`
	TransferTimestamp  int    `json:"transfer_timestamp"`
	AnchorTxHash       string `json:"anchor_tx_hash" gorm:"type:varchar(255)"`
	AnchorTxHeightHint int    `json:"anchor_tx_height_hint"`
	AnchorTxChainFees  int    `json:"anchor_tx_chain_fees"`
	DeviceID           string `json:"device_id" gorm:"type:varchar(255)"`
	UserID             int    `json:"user_id"`
	Status             int    `json:"status" gorm:"default:1"`
}

type BatchTransferRequest struct {
	Encoded            string `json:"encoded"`
	AssetID            string `json:"asset_id"`
	Amount             int    `json:"amount"`
	ScriptKey          string `json:"script_key"`
	InternalKey        string `json:"internal_key"`
	TaprootOutputKey   string `json:"taproot_output_key"`
	ProofCourierAddr   string `json:"proof_courier_addr"`
	Txid               string `json:"txid"`
	Index              int    `json:"index"`
	TransferTimestamp  int    `json:"transfer_timestamp"`
	AnchorTxHash       string `json:"anchor_tx_hash"`
	AnchorTxHeightHint int    `json:"anchor_tx_height_hint"`
	AnchorTxChainFees  int    `json:"anchor_tx_chain_fees"`
	DeviceID           string `json:"device_id"`
}
