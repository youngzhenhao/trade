package models

import "gorm.io/gorm"

type LoginRecord struct {
	gorm.Model
	UserId            uint   `gorm:"column:user_id;type:bigint unsigned;index" json:"user_id"` // column选项放在同一个gorm标签内
	RecentIpAddresses string `json:"recent_ip_addresses" gorm:"type:varchar(255);index"`
	Path              string `json:"path" gorm:"type:varchar(128);index"`
	LoginTime         int    `json:"login_time" gorm:"type:bigint;index"`
}

func (LoginRecord) TableName() string {
	return "login_record"
}
