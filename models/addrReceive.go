package models

import "gorm.io/gorm"

type AddrReceiveEvent struct {
	gorm.Model
	CreationTimeUnixSeconds int    `json:"creation_time_unix_seconds"`
	AddrEncoded             string `json:"addr_encoded"`
	AddrAssetID             string `json:"addr_asset_id" gorm:"type:varchar(255)"`
	AddrAmount              int    `json:"addr_amount"`
	AddrScriptKey           string `json:"addr_script_key" gorm:"type:varchar(255)"`
	AddrInternalKey         string `json:"addr_internal_key" gorm:"type:varchar(255)"`
	AddrTaprootOutputKey    string `json:"addr_taproot_output_key" gorm:"type:varchar(255)"`
	AddrProofCourierAddr    string `json:"addr_proof_courier_addr"`
	EventStatus             string `json:"event_status" gorm:"type:varchar(255)"`
	Outpoint                string `json:"outpoint" gorm:"type:varchar(255)"`
	UtxoAmtSat              int    `json:"utxo_amt_sat"`
	ConfirmationHeight      int    `json:"confirmation_height"`
	HasProof                bool   `json:"has_proof,omitempty"`
	DeviceID                string `json:"device_id" gorm:"type:varchar(255)"`
	UserID                  int    `json:"user_id"`
	Username                string `json:"username" gorm:"type:varchar(255)"`
	Status                  int    `json:"status" gorm:"default:1"`
}

type AddrReceiveEventSetRequest struct {
	CreationTimeUnixSeconds int                            `json:"creation_time_unix_seconds"`
	Addr                    AddrReceiveEventSetRequestAddr `json:"addr"`
	Status                  string                         `json:"status"`
	Outpoint                string                         `json:"outpoint"`
	UtxoAmtSat              int                            `json:"utxo_amt_sat"`
	ConfirmationHeight      int                            `json:"confirmation_height"`
	HasProof                bool                           `json:"has_proof,omitempty"`
	DeviceID                string                         `json:"device_id"`
}

type AddrReceiveEventSetRequestAddr struct {
	Encoded          string `json:"encoded"`
	AssetID          string `json:"asset_id"`
	Amount           int    `json:"amount"`
	ScriptKey        string `json:"script_key"`
	InternalKey      string `json:"internal_key"`
	TaprootOutputKey string `json:"taproot_output_key"`
	ProofCourierAddr string `json:"proof_courier_addr"`
}
