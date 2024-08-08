package models

import "gorm.io/gorm"

type LoginRecord struct {
	gorm.Model
	UserId            uint   `gorm:"column:user_id;type:bigint unsigned" json:"userId"` // column选项放在同一个gorm标签内
	RecentIpAddresses string `json:"recent_ip_addresses" gorm:"type:varchar(255)"`
	Path              string `json:"path" gorm:"type:varchar(128)"`
	LoginTime         int    `json:"login_time" gorm:"type:bigint"`
}

func (LoginRecord) TableName() string {
	return "login_record"
}
