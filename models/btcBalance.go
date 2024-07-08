package models

import "gorm.io/gorm"

type BtcBalance struct {
	gorm.Model
	Username           string `json:"username" gorm:"type:varchar(255)"`
	TotalBalance       int    `json:"total_balance"`
	ConfirmedBalance   int    `json:"confirmed_balance"`
	UnconfirmedBalance int    `json:"unconfirmed_balance"`
	LockedBalance      int    `json:"locked_balance"`
	Status             int    `json:"status" gorm:"default:1"`
	DeviceID           string `json:"device_id" gorm:"type:varchar(255)"`
}

type BtcBalanceSetRequest struct {
	TotalBalance       int    `json:"total_balance"`
	ConfirmedBalance   int    `json:"confirmed_balance"`
	UnconfirmedBalance int    `json:"unconfirmed_balance"`
	LockedBalance      int    `json:"locked_balance"`
	DeviceID           string `json:"device_id"`
}
