package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username          string `gorm:"unique;column:user_name;type:varchar(255)" json:"userName"` // 正确地将unique和column选项放在同一个gorm标签内
	Password          string `gorm:"column:password" json:"password"`
	Status            int16  `gorm:"column:status;type:smallint" json:"status"`
	RecentIpAddresses string `json:"recent_ip_addresses" gorm:"type:varchar(255)"`
	RecentLoginTime   int    `json:"recent_login_time"`
}

func (User) TableName() string {
	return "user"
}
