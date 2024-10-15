package models

import (
	"time"
)

type UserData struct {
	QueryTime               time.Time
	UserInfo                *UserInfoData
	UserBtcBalance          *UserBtcBalanceData
	UserAssetBalance        *[]UserAssetBalanceData
	UserAddrReceive         *[]UserAddrReceiveData
	UserAssetTransfer       *[]UserAssetTransferData
	UserAccountBtcBalance   *UserAccountBtcBalanceData
	UserAccountAssetBalance *[]UserAccountAssetBalanceData
}

type UserInfoData struct {
	Username     string
	UserId       int
	CreatedAt    time.Time
	Account      string
	AccountId    int
	UserRecentIp string
}

type UserBtcBalanceData struct {
	CreatedAt          time.Time
	UpdatedAt          time.Time
	TotalBalance       int
	ConfirmedBalance   int
	UnconfirmedBalance int
	LockedBalance      int
}

type UserAssetBalanceData struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	AssetId   string
	Name      string
	AssetType string
	Balance   int
}

type UserAddrReceiveData struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CreationTime time.Time
	AssetId      string
	AddrEncoded  string
	Amount       int
	Outpoint     string
}

type UserAssetTransferData struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Txid         string
	AssetId      string
	TransferTime time.Time
}

type UserAccountBtcBalanceData struct {
	Amount int
}

type UserAccountAssetBalanceData struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	AssetId   string
	Amount    float64
}
