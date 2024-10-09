package models

import (
	"gorm.io/gorm"
)

type NftPresale struct {
	gorm.Model
	AssetId         string           `json:"asset_id" gorm:"type:varchar(255);index"`
	Name            string           `json:"name"`
	AssetType       string           `json:"asset_type" gorm:"type:varchar(255);index"`
	Meta            string           `json:"meta"`
	GroupKey        string           `json:"group_key" gorm:"type:varchar(255);index"`
	Amount          int              `json:"amount" gorm:"index"`
	Price           int              `json:"price"`
	Info            string           `json:"info"`
	BuyerUserId     int              `json:"buyer_user_id" gorm:"index"`
	BuyerUsername   string           `json:"buyer_username" gorm:"type:varchar(255);index"`
	BuyerDeviceId   string           `json:"buyer_device_id" gorm:"type:varchar(255);index"`
	PayMethod       FeePaymentMethod `json:"pay_method" gorm:"index"`
	LaunchTime      int              `json:"launch_time"`
	BoughtTime      int              `json:"bought_time"`
	PaidSuccessTime int              `json:"paid_success_time"`
	SentTime        int              `json:"sent_time"`
	State           NftPresaleState  `json:"state" gorm:"index"`
	ProcessNumber   int              `json:"process_number"`
}

type (
	NftPresaleState int
)

const (
	NftPresaleStateLaunched NftPresaleState = iota
	NftPresaleStateBoughtNotPay
	NftPresaleStatePaidNotSend
	NftPresaleStateSent
)

type NftPresaleSetRequest struct {
	AssetId string `json:"asset_id"`
	Price   int    `json:"price"`
	Info    string `json:"info"`
}
