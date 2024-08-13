package models

import "gorm.io/gorm"

type FairLaunchIncomeType int

type FairLaunchIncome struct {
	gorm.Model
	AssetId                string               `json:"asset_id" gorm:"type:varchar(255);index"`
	FairLaunchInfoId       int                  `json:"fair_launch_info_id" gorm:"index"`
	FairLaunchMintedInfoId int                  `json:"fair_launch_minted_info_id" gorm:"index"`
	FeePaidId              int                  `json:"fee_paid_id"`
	IncomeType             FairLaunchIncomeType `json:"income_type" gorm:"index"`
	IsIncome               bool                 `json:"is_income" gorm:"index"`
	SatAmount              int                  `json:"sat_amount"`
	Txid                   string               `json:"txid"`
	Addrs                  string               `json:"addrs"`
	UserId                 int                  `json:"user_id" gorm:"index"`
	Username               string               `json:"username" gorm:"type:varchar(255)"`
}

const (
	UserPayIssuanceFee FairLaunchIncomeType = iota
	ServerPayIssuanceFinalizeFee
	ServerPaySendReservedFee
	UserPayMintedFee
	ServerPaySendAssetFee
)
