package models

import (
	"github.com/lightninglabs/taproot-assets/taprpc"
	"gorm.io/gorm"
)

type AssetTransferType int

const (
	AssetTransferTypeOut AssetTransferType = iota
	AssetTransferTypeIn
)

type AssetTransferProcessed struct {
	gorm.Model
	Txid               string                         `json:"txid" gorm:"type:varchar(255)"`
	AssetID            string                         `json:"asset_id" gorm:"type:varchar(255)"`
	TransferTimestamp  int                            `json:"transfer_timestamp"`
	AnchorTxHash       string                         `json:"anchor_tx_hash" gorm:"type:varchar(255)"`
	AnchorTxHeightHint int                            `json:"anchor_tx_height_hint"`
	AnchorTxChainFees  int                            `json:"anchor_tx_chain_fees"`
	Inputs             []AssetTransferProcessedInput  `json:"inputs"`
	Outputs            []AssetTransferProcessedOutput `json:"outputs"`
	UserID             int                            `json:"user_id"`
	Status             int                            `json:"status" gorm:"default:1"`
}

type AssetTransferProcessedInput struct {
	Address     string `json:"address" gorm:"type:varchar(255)"`
	Amount      int    `json:"amount"`
	AnchorPoint string `json:"anchor_point" gorm:"type:varchar(255)"`
	ScriptKey   string `json:"script_key" gorm:"type:varchar(255)"`
}

type AssetTransferProcessedOutput struct {
	Address                string `json:"address" gorm:"type:varchar(255)"`
	Amount                 int    `json:"amount"`
	AnchorOutpoint         string `json:"anchor_outpoint" gorm:"type:varchar(255)"`
	AnchorValue            int    `json:"anchor_value"`
	AnchorInternalKey      string `json:"anchor_internal_key" gorm:"type:varchar(255)"`
	AnchorTaprootAssetRoot string `json:"anchor_taproot_asset_root" gorm:"type:varchar(255)"`
	AnchorMerkleRoot       string `json:"anchor_merkle_root" gorm:"type:varchar(255)"`
	AnchorTapscriptSibling string `json:"anchor_tapscript_sibling" gorm:"type:varchar(255)"`
	AnchorNumPassiveAssets int    `json:"anchor_num_passive_assets"`
	ScriptKey              string `json:"script_key" gorm:"type:varchar(255)"`
	ScriptKeyIsLocal       bool   `json:"script_key_is_local"`
	NewProofBlob           string `json:"new_proof_blob"`
	SplitCommitRootHash    string `json:"split_commit_root_hash" gorm:"type:varchar(255)"`
	OutputType             string `json:"output_type" gorm:"type:varchar(255)"`
	AssetVersion           string `json:"asset_version" gorm:"type:varchar(255)"`
}

type AssetTransferProcessedSetRequest struct {
	Txid               string                         `json:"txid" gorm:"type:varchar(255)"`
	AssetID            string                         `json:"asset_id" gorm:"type:varchar(255)"`
	TransferTimestamp  int                            `json:"transfer_timestamp"`
	AnchorTxHash       string                         `json:"anchor_tx_hash" gorm:"type:varchar(255)"`
	AnchorTxHeightHint int                            `json:"anchor_tx_height_hint"`
	AnchorTxChainFees  int                            `json:"anchor_tx_chain_fees"`
	Inputs             []AssetTransferProcessedInput  `json:"inputs"`
	Outputs            []AssetTransferProcessedOutput `json:"outputs"`
}

// @dev: These models may be deprecated.

type AssetTransfer struct {
	gorm.Model
	AssetID           string            `json:"asset_id" gorm:"type:varchar(255)"`
	AssetName         string            `json:"asset_name" gorm:"type:varchar(255)"`
	AssetType         taprpc.AssetType  `json:"asset_type"`
	AssetAddressFrom  string            `json:"address_from" gorm:"type:varchar(255)"`
	AssetAddressTo    string            `json:"address_to" gorm:"type:varchar(255)"`
	Amount            int               `json:"amount"`
	TransferType      AssetTransferType `json:"transfer_type"`
	UserID            int               `json:"user_id"`
	TransactionID     string            `json:"transaction_id" gorm:"type:varchar(255)"`
	TransferTimestamp int               `json:"transfer_timestamp"`
	AnchorTxChainFees int               `json:"anchor_tx_chain_fees"`
	ConfirmedBlocks   int               `json:"confirmed_blocks"`
	Status            int               `json:"status" gorm:"default:1"`
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
	AssetID           string            `json:"asset_id" gorm:"type:varchar(255)"`
	AssetAddressFrom  string            `json:"address_from" gorm:"type:varchar(255)"`
	AssetAddressTo    string            `json:"address_to" gorm:"type:varchar(255)"`
	Amount            int               `json:"amount"`
	TransferType      AssetTransferType `json:"transfer_type"`
	TransactionID     string            `json:"transaction_id" gorm:"type:varchar(255)"`
	TransferTimestamp int               `json:"transfer_timestamp"`
	AnchorTxChainFees int               `json:"anchor_tx_chain_fees"`
}
