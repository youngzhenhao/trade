package models

import "gorm.io/gorm"

type PayOutside struct {
	gorm.Model
}

func (PayOutside) TableName() string {
	return "user_out_inside"
}
