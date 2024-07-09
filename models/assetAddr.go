package models

import "gorm.io/gorm"

type AssetAddr struct {
	gorm.Model
	Encoded          string `json:"encoded"`
	AssetId          string `json:"asset_id" gorm:"type:varchar(255)"`
	AssetType        int    `json:"asset_type"`
	Amount           int    `json:"amount"`
	GroupKey         string `json:"group_key" gorm:"type:varchar(255)"`
	ScriptKey        string `json:"script_key" gorm:"type:varchar(255)"`
	InternalKey      string `json:"internal_key" gorm:"type:varchar(255)"`
	TapscriptSibling string `json:"tapscript_sibling" gorm:"type:varchar(255)"`
	TaprootOutputKey string `json:"taproot_output_key" gorm:"type:varchar(255)"`
	ProofCourierAddr string `json:"proof_courier_addr" gorm:"type:varchar(255)"`
	AssetVersion     int    `json:"asset_version"`
	DeviceID         string `json:"device_id" gorm:"type:varchar(255)"`
	UserId           int    `json:"user_id"`
	Status           int    `json:"status" gorm:"default:1"`
}

type AssetAddrSetRequest struct {
	Encoded          string `json:"encoded"`
	AssetId          string `json:"asset_id"`
	AssetType        int    `json:"asset_type"`
	Amount           int    `json:"amount"`
	GroupKey         string `json:"group_key"`
	ScriptKey        string `json:"script_key"`
	InternalKey      string `json:"internal_key"`
	TapscriptSibling string `json:"tapscript_sibling"`
	TaprootOutputKey string `json:"taproot_output_key"`
	ProofCourierAddr string `json:"proof_courier_addr"`
	AssetVersion     int    `json:"asset_version"`
	DeviceID         string `json:"device_id"`
}
