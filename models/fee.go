package models

import "gorm.io/gorm"

type (
	FeePaymentMethod int
	FeeRateType      int
)

var (
	FeePaymentMethodCustodyAccount FeePaymentMethod = 0
	FeeRateTypeBtcPerKb            FeeRateType      = 0
	FeeRateTypeSatPerB             FeeRateType      = 1
	FeeRateTypeSatPerKw            FeeRateType      = 2
)

type FeeRateInfo struct {
	gorm.Model
	Name    string      `json:"name" gorm:"type:varchar(255);not null"`
	Unit    FeeRateType `json:"unit"`
	FeeRate float64     `json:"fee_rate"`
}

type MempoolFeeRateInfo struct {
	gorm.Model
	Name                 string  `json:"name" gorm:"type:varchar(255);not null"`
	EstimateSmartFeeRate float64 `json:"estimate_smart_fee_rate"`
}
