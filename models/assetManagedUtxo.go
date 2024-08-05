package models

import "gorm.io/gorm"

type AssetManagedUtxo struct {
	gorm.Model
	Op                          string `json:"op" gorm:"type:varchar(255)"`
	OutPoint                    string `json:"out_point" gorm:"type:varchar(255)"`
	Time                        int    `json:"time"`
	AmtSat                      int    `json:"amt_sat"`
	InternalKey                 string `json:"internal_key" gorm:"type:varchar(255)"`
	TaprootAssetRoot            string `json:"taproot_asset_root" gorm:"type:varchar(255)"`
	MerkleRoot                  string `json:"merkle_root" gorm:"type:varchar(255)"`
	Version                     string `json:"version" gorm:"type:varchar(255)"`
	AssetGenesisPoint           string `json:"asset_genesis_point" gorm:"type:varchar(255)"`
	AssetGenesisName            string `json:"asset_genesis_name" gorm:"type:varchar(255)"`
	AssetGenesisMetaHash        string `json:"asset_genesis_meta_hash" gorm:"type:varchar(255)"`
	AssetGenesisAssetID         string `json:"asset_genesis_asset_id" gorm:"type:varchar(255)"`
	AssetGenesisAssetType       string `json:"asset_genesis_asset_type" gorm:"type:varchar(255)"`
	AssetGenesisOutputIndex     int    `json:"asset_genesis_output_index"`
	AssetGenesisVersion         int    `json:"asset_genesis_version"`
	Amount                      int    `json:"amount"`
	LockTime                    int    `json:"lock_time"`
	RelativeLockTime            int    `json:"relative_lock_time"`
	ScriptVersion               int    `json:"script_version"`
	ScriptKey                   string `json:"script_key" gorm:"type:varchar(255)"`
	ScriptKeyIsLocal            bool   `json:"script_key_is_local"`
	AssetGroupRawGroupKey       string `json:"asset_group_raw_group_key" gorm:"type:varchar(255)"`
	AssetGroupTweakedGroupKey   string `json:"asset_group_tweaked_group_key" gorm:"type:varchar(255)"`
	AssetGroupAssetWitness      string `json:"asset_group_asset_witness"`
	ChainAnchorTx               string `json:"chain_anchor_tx"`
	ChainAnchorBlockHash        string `json:"chain_anchor_block_hash" gorm:"type:varchar(255)"`
	ChainAnchorOutpoint         string `json:"chain_anchor_outpoint" gorm:"type:varchar(255)"`
	ChainAnchorInternalKey      string `json:"chain_anchor_internal_key" gorm:"type:varchar(255)"`
	ChainAnchorMerkleRoot       string `json:"chain_anchor_merkle_root" gorm:"type:varchar(255)"`
	ChainAnchorTapscriptSibling string `json:"chain_anchor_tapscript_sibling"`
	ChainAnchorBlockHeight      int    `json:"chain_anchor_block_height"`
	IsSpent                     bool   `json:"is_spent"`
	LeaseOwner                  string `json:"lease_owner" gorm:"type:varchar(255)"`
	LeaseExpiry                 int    `json:"lease_expiry"`
	IsBurn                      bool   `json:"is_burn"`
	DeviceId                    string `json:"device_id" gorm:"type:varchar(255)"`
	UserId                      int    `json:"user_id"`
	Username                    string `json:"username" gorm:"type:varchar(255)"`
	Status                      int    `json:"status" gorm:"default:1"`
}

type AssetManagedUtxoSetRequest struct {
	Op                          string `json:"op"`
	OutPoint                    string `json:"out_point"`
	Time                        int    `json:"time"`
	AmtSat                      int    `json:"amt_sat"`
	InternalKey                 string `json:"internal_key"`
	TaprootAssetRoot            string `json:"taproot_asset_root"`
	MerkleRoot                  string `json:"merkle_root"`
	Version                     string `json:"version"`
	AssetGenesisPoint           string `json:"asset_genesis_point"`
	AssetGenesisName            string `json:"asset_genesis_name"`
	AssetGenesisMetaHash        string `json:"asset_genesis_meta_hash"`
	AssetGenesisAssetID         string `json:"asset_genesis_asset_id"`
	AssetGenesisAssetType       string `json:"asset_genesis_asset_type"`
	AssetGenesisOutputIndex     int    `json:"asset_genesis_output_index"`
	AssetGenesisVersion         int    `json:"asset_genesis_version"`
	Amount                      int    `json:"amount"`
	LockTime                    int    `json:"lock_time"`
	RelativeLockTime            int    `json:"relative_lock_time"`
	ScriptVersion               int    `json:"script_version"`
	ScriptKey                   string `json:"script_key"`
	ScriptKeyIsLocal            bool   `json:"script_key_is_local"`
	AssetGroupRawGroupKey       string `json:"asset_group_raw_group_key"`
	AssetGroupTweakedGroupKey   string `json:"asset_group_tweaked_group_key"`
	AssetGroupAssetWitness      string `json:"asset_group_asset_witness"`
	ChainAnchorTx               string `json:"chain_anchor_tx"`
	ChainAnchorBlockHash        string `json:"chain_anchor_block_hash"`
	ChainAnchorOutpoint         string `json:"chain_anchor_outpoint"`
	ChainAnchorInternalKey      string `json:"chain_anchor_internal_key"`
	ChainAnchorMerkleRoot       string `json:"chain_anchor_merkle_root"`
	ChainAnchorTapscriptSibling string `json:"chain_anchor_tapscript_sibling"`
	ChainAnchorBlockHeight      int    `json:"chain_anchor_block_height"`
	IsSpent                     bool   `json:"is_spent"`
	LeaseOwner                  string `json:"lease_owner"`
	LeaseExpiry                 int    `json:"lease_expiry"`
	IsBurn                      bool   `json:"is_burn"`
	DeviceId                    string `json:"device_id"`
}
