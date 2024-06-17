package models

import (
	"github.com/lightninglabs/taproot-assets/taprpc"
	"gorm.io/gorm"
)

type (
	AssetIssuanceState int
)

const (
	AssetIssuanceStatePending AssetIssuanceState = iota
	AssetIssuanceStateIssued
)

type AssetIssuance struct {
	gorm.Model
	AssetName      string             `json:"asset_name" gorm:"type:varchar(255)"`
	AssetId        string             `json:"asset_id" gorm:"type:varchar(255)"`
	AssetType      taprpc.AssetType   `json:"asset_type"`
	IssuanceUserId int                `json:"issuance_user_id"`
	IssuanceTime   int                `json:"issuance_time"`
	IsFairLaunch   bool               `json:"is_fair_launch"`
	FairLaunchID   int                `json:"fair_launch_id"`
	Status         int                `json:"status" gorm:"default:1"`
	State          AssetIssuanceState `json:"state"`
}
