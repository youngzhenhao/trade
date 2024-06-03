package models

type AssetIssuanceLeaf struct {
	Version            string `json:"version"`
	GenesisPoint       string `json:"genesis_point"`
	Name               string `json:"name"`
	MetaHash           string `json:"meta_hash"`
	AssetID            string `json:"asset_id"`
	AssetType          string `json:"asset_type"`
	GenesisOutputIndex int    `json:"genesis_output_index"`
	Amount             int    `json:"amount"`
	LockTime           int    `json:"lock_time"`
	RelativeLockTime   int    `json:"relative_lock_time"`
	ScriptVersion      int    `json:"script_version"`
	ScriptKey          string `json:"script_key"`
	ScriptKeyIsLocal   bool   `json:"script_key_is_local"`
	IsSpent            bool   `json:"is_spent"`
	LeaseOwner         string `json:"lease_owner"`
	LeaseExpiry        int    `json:"lease_expiry"`
	IsBurn             bool   `json:"is_burn"`
	Proof              string `json:"proof"`
}
