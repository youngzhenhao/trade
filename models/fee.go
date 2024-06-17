package models

import "gorm.io/gorm"

type (
	FeePaymentMethod int
	FeeRateType      int
)

const (
	FeePaymentMethodCustodyAccount FeePaymentMethod = iota
)

const (
	FeeRateTypeBtcPerKb FeeRateType = iota
	FeeRateTypeSatPerB
	FeeRateTypeSatPerKw
)

type FeeRateInfo struct {
	gorm.Model
	Name    string      `json:"name" gorm:"type:varchar(255);not null"`
	Unit    FeeRateType `json:"unit"`
	FeeRate float64     `json:"fee_rate"`
}
