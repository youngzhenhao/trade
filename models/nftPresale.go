package models

import (
	"gorm.io/gorm"
	"time"
)

type NftPresale struct {
	gorm.Model
	BatchGroupId    int              `json:"batch_group_id" gorm:"index"`
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
	AddrScriptKey   string           `json:"addr_script_key" gorm:"type:varchar(255)"`
	AddrInternalKey string           `json:"addr_internal_key" gorm:"type:varchar(255)"`
	PayMethod       FeePaymentMethod `json:"pay_method" gorm:"index"`
	LaunchTime      int              `json:"launch_time"`
	StartTime       int              `json:"start_time"`
	EndTime         int              `json:"end_time"`
	BoughtTime      int              `json:"bought_time"`
	PaidId          int              `json:"paid_id" gorm:"index"`
	PaidSuccessTime int              `json:"paid_success_time"`
	SentTime        int              `json:"sent_time"`
	SentTxid        string           `json:"sent_txid" gorm:"type:varchar(255)"`
	SentOutpoint    string           `json:"sent_outpoint" gorm:"type:varchar(255)"`
	SentAddress     string           `json:"sent_address" gorm:"type:varchar(255)"`
	State           NftPresaleState  `json:"state" gorm:"index"`
	ProcessNumber   int              `json:"process_number"`
	IsReLaunched    bool             `json:"is_re_launched"`
}

type (
	NftPresaleState int
)

const (
	NftPresaleStateLaunched NftPresaleState = iota
	NftPresaleStateBoughtNotPay
	NftPresaleStatePaidPending
	NftPresaleStatePaidNotSend
	NftPresaleStateSentPending
	NftPresaleStateSent
	NftPresaleStateFailOrCanceled = -1
)

func (n NftPresaleState) String() string {
	nftPresaleStateMapString := map[NftPresaleState]string{
		NftPresaleStateLaunched:       "NftPresaleStateLaunched",
		NftPresaleStateBoughtNotPay:   "NftPresaleStateBoughtNotPay",
		NftPresaleStatePaidPending:    "NftPresaleStatePaidPending",
		NftPresaleStatePaidNotSend:    "NftPresaleStatePaidNotSend",
		NftPresaleStateSentPending:    "NftPresaleStateSentPending",
		NftPresaleStateSent:           "NftPresaleStateSent",
		NftPresaleStateFailOrCanceled: "NftPresaleStateFailOrCanceled",
	}
	return nftPresaleStateMapString[n]
}

type NftPresaleSetRequest struct {
	BatchGroupId int    `json:"batch_group_id" gorm:"index"`
	AssetId      string `json:"asset_id"`
	Price        int    `json:"price"`
}

type BuyNftPresaleRequest struct {
	AssetId     string `json:"asset_id"`
	ReceiveAddr string `json:"receive_addr"`
	DeviceId    string `json:"device_id"`
}

type NftPresaleSimplified struct {
	ID              uint `gorm:"primarykey"`
	UpdatedAt       time.Time
	BatchGroupId    int              `json:"batch_group_id" gorm:"index"`
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
	AddrScriptKey   string           `json:"addr_script_key" gorm:"type:varchar(255)"`
	AddrInternalKey string           `json:"addr_internal_key" gorm:"type:varchar(255)"`
	PayMethod       FeePaymentMethod `json:"pay_method" gorm:"index"`
	LaunchTime      int              `json:"launch_time"`
	StartTime       int              `json:"start_time"`
	EndTime         int              `json:"end_time"`
	BoughtTime      int              `json:"bought_time"`
	PaidId          int              `json:"paid_id" gorm:"index"`
	PaidSuccessTime int              `json:"paid_success_time"`
	SentTime        int              `json:"sent_time"`
	SentTxid        string           `json:"sent_txid" gorm:"type:varchar(255)"`
	SentOutpoint    string           `json:"sent_outpoint" gorm:"type:varchar(255)"`
	SentAddress     string           `json:"sent_address" gorm:"type:varchar(255)"`
	State           NftPresaleState  `json:"state" gorm:"index"`
	ProcessNumber   int              `json:"process_number"`
	IsReLaunched    bool             `json:"is_re_launched"`
	MetaStr         string           `json:"meta_str"`
}
