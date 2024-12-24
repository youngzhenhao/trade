package models

import "gorm.io/gorm"

type UnspentUtxo struct {
	AddressType   string `json:"address_type"`
	Address       string `json:"address"`
	AmountSat     int64  `json:"amount_sat"`
	PkScript      string `json:"pk_script"`
	Outpoint      string `json:"outpoint"`
	Confirmations int64  `json:"confirmations"`
}

type BtcUtxo struct {
	gorm.Model
	Username string `json:"username"`
	UnspentUtxo
}
