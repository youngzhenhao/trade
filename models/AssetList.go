package models

type AssetList struct {
	Version string `json:"version"`
	//@dev: AssetGenesis
	GenesisPoint string `json:"genesis_point"`
	Name         string `json:"name"`
	MetaHash     string `json:"meta_hash"`
	AssetID      string `json:"asset_id"`
	AssetType    string `json:"asset_type"`
	OutputIndex  int    `json:"output_index"`
	//GenesisVersion int    `json:"genesis_version"`

	Amount           int   `json:"amount"`
	LockTime         int32 `json:"lock_time"`
	RelativeLockTime int32 `json:"relative_lock_time"`
	//ScriptVersion    int32  `json:"script_version"`
	ScriptKey string `json:"script_key"`
	//ScriptKeyIsLocal bool   `json:"script_key_is_local"`

	//@dev: ChainAnchor
	//AnchorTx         string `json:"anchor_tx"`
	//AnchorBlockHash string `json:"anchor_block_hash"`
	AnchorOutpoint string `json:"anchor_outpoint"`
	//InternalKey     string `json:"internal_key"`
	//MerkleRoot      string `json:"merkle_root"`
	//TapscriptSibling string `json:"tapscript_sibling"`
	//BlockHeight int `json:"block_height"`

	//@dev: _AssetGroup
	//RawGroupKey     string `json:"raw_group_key"`
	TweakedGroupKey string `json:"tweaked_group_key"`
	//AssetWitness    string `json:"asset_witness"`

	//IsSpent     bool   `json:"is_spent"`
	//LeaseOwner  string `json:"lease_owner"`
	//LeaseExpiry int    `json:"lease_expiry"`
	//IsBurn      bool   `json:"is_burn"`

	DeviceId string `json:"device_id" gorm:"type:varchar(255)"`
	UserId   int    `json:"user_id"`
	Username string `json:"username" gorm:"type:varchar(255)"`
}

type AssetListSetRequest struct {
	Version          string `json:"version" gorm:"type:varchar(255);index"`
	GenesisPoint     string `json:"genesis_point" gorm:"type:varchar(255)"`
	Name             string `json:"name" gorm:"type:varchar(255);index"`
	MetaHash         string `json:"meta_hash" gorm:"type:varchar(255);index"`
	AssetID          string `json:"asset_id" gorm:"type:varchar(255);index"`
	AssetType        string `json:"asset_type" gorm:"type:varchar(255);index"`
	OutputIndex      int    `json:"output_index"`
	Amount           int    `json:"amount"`
	LockTime         int32  `json:"lock_time"`
	RelativeLockTime int32  `json:"relative_lock_time"`
	ScriptKey        string `json:"script_key" gorm:"type:varchar(255);index"`
	AnchorOutpoint   string `json:"anchor_outpoint" gorm:"type:varchar(255);index"`
	TweakedGroupKey  string `json:"tweaked_group_key" gorm:"type:varchar(255);index"`
	DeviceId         string `json:"device_id" gorm:"type:varchar(255);index"`
}
