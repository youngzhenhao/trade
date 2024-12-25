package models

import "gorm.io/gorm"

type UnspentUtxo struct {
	AddressType   string `json:"address_type" gorm:"type:varchar(255);index"`
	Address       string `json:"address" gorm:"type:varchar(255);index"`
	AmountSat     int64  `json:"amount_sat"`
	PkScript      string `json:"pk_script" gorm:"type:varchar(255);index"`
	Outpoint      string `json:"outpoint" gorm:"type:varchar(255);uniqueIndex"`
	Confirmations int64  `json:"confirmations"`
}

type BtcUtxo struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(255);index"`
	UnspentUtxo
}

type BtcUtxoHistory struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(255);index"`
	UnspentUtxo
}
