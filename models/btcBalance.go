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
}
