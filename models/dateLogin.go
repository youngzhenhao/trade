package models

import (
	"gorm.io/gorm"
)

type DateLogin struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(255);index;uniqueIndex:idx_username_date"`
	Date     string `json:"date" gorm:"type:varchar(255);index;uniqueIndex:idx_username_date"`
}
