package models

import (
	"gorm.io/gorm"
)

type DateIpLogin struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(255);index;uniqueIndex:idx_username_ip_date"`
	Date     string `json:"date" gorm:"type:varchar(255);index;uniqueIndex:idx_username_ip_date"`
	Ip       string `json:"ip" gorm:"type:varchar(255);index;uniqueIndex:idx_username_ip_date"`
}
