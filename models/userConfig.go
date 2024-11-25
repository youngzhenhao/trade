package models

import "gorm.io/gorm"

type UserConfig struct {
	gorm.Model
	UserID uint   `gorm:"column:user_id;uniqueIndex:user_config_user_id_key" json:"userId"`
	Config string `gorm:"column:config;type:text" json:"config"`
	User   User   `gorm:"foreignKey:UserID" json:"user"`
}

func (UserConfig) TableName() string {
	return "user_config"
}
