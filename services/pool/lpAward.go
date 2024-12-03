package pool

import (
	"gorm.io/gorm"
)

// TODO
type LpAwardBalance struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(255);uniqueIndex"`
	Balance  string `json:"balance" gorm:"type:varchar(255);index"`
}

// TODO
type LpAwardRecord struct {
	gorm.Model
}

//func t() {
//	fromString, err := decimal.NewFromString("-123.4567")
//	if err != nil {
//		return
//	}
//}
