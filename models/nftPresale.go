package models

import (
	"gorm.io/gorm"
	"time"
)

type NftPresale struct {
	gorm.Model
	AssetId       string `json:"asset_id" gorm:"type:varchar(255);index"`
	Name          string `json:"name"`
	AssetType     string `json:"asset_type" gorm:"type:varchar(255);index"`
	Meta          string `json:"meta"`
	GroupKey      string `json:"group_key" gorm:"type:varchar(255);index"`
	Amount        int    `json:"amount" gorm:"index"`
	Price         int    `json:"price"`
	Info          string `json:"info"`
	BuyerUserId   int    `json:"buyer_user_id" gorm:"index"`
	BuyerUsername string `json:"buyer_username" gorm:"type:varchar(255);index"`
	BuyerDeviceId string `json:"buyer_device_id" gorm:"type:varchar(255);index"`
	ReceiveAddr   string `json:"receive_addr"`
	// TODO
	PayMethod  FeePaymentMethod `json:"pay_method" gorm:"index"`
	LaunchTime int              `json:"launch_time"`
	BoughtTime int              `json:"bought_time"`
	// TODO
	PaidId int `json:"paid_id" gorm:"index"`
	// TODO
	PaidSuccessTime int `json:"paid_success_time"`
	// TODO
	SentTime int `json:"sent_time"`
	// TODO
	State NftPresaleState `json:"state" gorm:"index"`
	// TODO
	ProcessNumber int `json:"process_number"`
}

type (
	NftPresaleState int
)

const (
	NftPresaleStateLaunched NftPresaleState = iota
	NftPresaleStateBoughtNotPay
	NftPresaleStatePaidNotSend
	NftPresaleStateSent
	NftPresaleStateCanceled = -1
)

func (n NftPresaleState) String() string {
	nftPresaleStateMapString := map[NftPresaleState]string{
		NftPresaleStateLaunched:     "NftPresaleStateLaunched",
		NftPresaleStateBoughtNotPay: "NftPresaleStateBoughtNotPay",
		NftPresaleStatePaidNotSend:  "NftPresaleStatePaidNotSend",
		NftPresaleStateSent:         "NftPresaleStateSent",
		NftPresaleStateCanceled:     "NftPresaleStateCanceled",
	}
	return nftPresaleStateMapString[n]
}

type NftPresaleSetRequest struct {
	AssetId string `json:"asset_id"`
	Price   int    `json:"price"`
	Info    string `json:"info"`
}

type BuyNftPresaleRequest struct {
	AssetId     string `json:"asset_id"`
	ReceiveAddr string `json:"receive_addr"`
	DeviceId    string `json:"device_id"`
}

type NftPresaleSimplified struct {
	ID              uint `gorm:"primarykey"`
	UpdatedAt       time.Time
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
	ReceiveAddr     string           `json:"receive_addr"`
	PayMethod       FeePaymentMethod `json:"pay_method" gorm:"index"`
	LaunchTime      int              `json:"launch_time"`
	BoughtTime      int              `json:"bought_time"`
	PaidId          int              `json:"paid_id" gorm:"index"`
	PaidSuccessTime int              `json:"paid_success_time"`
	SentTime        int              `json:"sent_time"`
	State           NftPresaleState  `json:"state" gorm:"index"`
	ProcessNumber   int              `json:"process_number"`
}
