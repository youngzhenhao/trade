package models

import (
	"time"
)

type UserData struct {
	QueryTime               time.Time                      `json:"查询时间" yaml:"查询时间"`
	UserInfo                *UserInfoData                  `json:"用户基本信息" yaml:"用户基本信息"`
	UserBtcBalance          *UserBtcBalanceData            `json:"上传的比特币余额记录" yaml:"上传的比特币余额记录"`
	UserAssetBalance        *[]UserAssetBalanceData        `json:"上传的资产余额记录" yaml:"上传的资产余额记录"`
	UserAddrReceive         *[]UserAddrReceiveData         `json:"上传的资产接收记录" yaml:"上传的资产接收记录"`
	UserAssetTransfer       *[]UserAssetTransferData       `json:"上传的资产转出记录" yaml:"上传的资产转出记录"`
	UserAccountBtcBalance   *UserAccountBtcBalanceData     `json:"用户托管账户比特币余额" yaml:"用户托管账户比特币余额"`
	UserAccountAssetBalance *[]UserAccountAssetBalanceData `json:"用户托管账户资产余额" yaml:"用户托管账户资产余额"`
	ErrorInfos              *[]string                      `json:"错误信息" yaml:"错误信息"`
}

type UserInfoData struct {
	Username     string    `json:"用户名;NpubKey;Nostr地址"  yaml:"用户名;NpubKey;Nostr地址"`
	UserId       int       `json:"用户ID"  yaml:"用户ID"`
	CreatedAt    time.Time `json:"用户创建时间"  yaml:"用户创建时间"`
	Account      string    `json:"托管账户码"  yaml:"托管账户码"`
	AccountId    int       `json:"托管账户ID"  yaml:"托管账户ID"`
	UserRecentIp string    `json:"最近登录IP"  yaml:"最近登录IP"`
}

type UserBtcBalanceData struct {
	CreatedAt          time.Time `json:"本条记录创建时间"  yaml:"本条记录创建时间"`
	UpdatedAt          time.Time `json:"本条记录最近修改时间"  yaml:"本条记录最近修改时间"`
	TotalBalance       int       `json:"总余额"  yaml:"总余额"`
	ConfirmedBalance   int       `json:"已确认"  yaml:"已确认"`
	UnconfirmedBalance int       `json:"未确认"  yaml:"未确认"`
	LockedBalance      int       `json:"锁定"  yaml:"锁定"`
}

type UserAssetBalanceData struct {
	CreatedAt time.Time `json:"本条记录创建时间"  yaml:"本条记录创建时间"`
	UpdatedAt time.Time `json:"本条记录最近修改时间"  yaml:"本条记录最近修改时间"`
	AssetId   string    `json:"资产ID"  yaml:"资产ID"`
	Name      string    `json:"资产名称"  yaml:"资产名称"`
	AssetType string    `json:"资产类型"  yaml:"资产类型"`
	Balance   int       `json:"余额"  yaml:"余额"`
}

type UserAddrReceiveData struct {
	CreatedAt    time.Time `json:"本条记录创建时间"  yaml:"本条记录创建时间"`
	UpdatedAt    time.Time `json:"本条记录最近修改时间"  yaml:"本条记录最近修改时间"`
	CreationTime time.Time `json:"用户本地创建时间"  yaml:"用户本地创建时间"`
	AssetId      string    `json:"资产ID"  yaml:"资产ID"`
	AddrEncoded  string    `json:"编码的资产地址"  yaml:"编码的资产地址"`
	Amount       int       `json:"地址接收资产数量"  yaml:"地址接收资产数量"`
	Outpoint     string    `json:"输出点"  yaml:"输出点"`
}

type UserAssetTransferData struct {
	CreatedAt    time.Time `json:"本条记录创建时间"  yaml:"本条记录创建时间"`
	UpdatedAt    time.Time `json:"本条记录最近修改时间"  yaml:"本条记录最近修改时间"`
	Txid         string    `json:"交易ID"  yaml:"交易ID"`
	AssetId      string    `json:"资产ID"  yaml:"资产ID"`
	TransferTime time.Time `json:"用户本地转账时间"  yaml:"用户本地转账时间"`
}

type UserAccountBtcBalanceData struct {
	Amount int `json:"数量"  yaml:"数量"`
}

type UserAccountAssetBalanceData struct {
	CreatedAt time.Time `json:"本条记录创建时间"  yaml:"本条记录创建时间"`
	UpdatedAt time.Time `json:"本条记录最近修改时间"  yaml:"本条记录最近修改时间"`
	AssetId   string    `json:"资产ID"  yaml:"资产ID"`
	Amount    float64   `json:"数量"  yaml:"数量"`
}
