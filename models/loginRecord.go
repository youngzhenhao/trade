package models

import "gorm.io/gorm"

type LoginRecord struct {
	gorm.Model
	UserId            uint   `gorm:"column:user_id;type:bigint unsigned" json:"userId"` // column选项放在同一个gorm标签内
	RecentIpAddresses string `json:"recent_ip_addresses" gorm:"type:varchar(255)"`
}

func (LoginRecord) TableName() string {
	return "login_record"
}
