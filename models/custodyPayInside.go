package models

import "gorm.io/gorm"

type PayInside struct {
	gorm.Model
	PayUserId     uint            `gorm:"column:pay_user_id;type:bigint;not null;index:idx_pay_user_id" json:"pay_user_id"`
	GasFee        uint64          `gorm:"column:gas_fee;type:bigint" json:"gas_fee"`
	ServeFee      uint64          `gorm:"column:serve_fee;type:bigint" json:"serve_fee"`
	ReceiveUserId uint            `gorm:"column:receive_user_id;type:bigint" json:"receive_user_id"`
	PayType       PayInsideType   `gorm:"column:pay_type;type:smallint" json:"pay_type"`
	AssetType     string          `gorm:"column:asset_type;type:varchar(128);default:'00'" json:"asset_type"`
	PayReq        *string         `gorm:"column:pay_req;type:varchar(512)" json:"pay_req"`
	BalanceId     uint            `gorm:"column:balance_id;type:bigint;default:0" json:"balance_id"`
	Status        PayInsideStatus `gorm:"column:status;type:smallint" json:"status"`
}

func (PayInside) TableName() string {
	return "user_pay_inside"
}

type PayInsideType uint16

const (
	PayInsideToAdmin   PayInsideType = 1
	PayInsideByInvoice PayInsideType = 2
	PayInsideByAddress PayInsideType = 3

	FairLunchFee         PayInsideType = 5
	ChannelBTCFee        PayInsideType = 6
	ChannelBTCOutSideFee PayInsideType = 7
	AssetInSideFee       PayInsideType = 8
	AssetOutSideFee      PayInsideType = 9
)

type PayInsideStatus uint16

const (
	PayInsideStatusPending PayInsideStatus = 0
	PayInsideStatusSuccess PayInsideStatus = 1
	PayInsideStatusFailed  PayInsideStatus = 2
)
