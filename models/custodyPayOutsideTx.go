package models

import "gorm.io/gorm"

type PayOutsideTx struct {
	gorm.Model
}

func (PayOutsideTx) TableName() string {
	return "user_out_inside_tx"
}
