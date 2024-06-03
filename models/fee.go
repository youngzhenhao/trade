package models

import "gorm.io/gorm"

type (
	FeePaymentMethod int
)

var (
	FeePaymentMethodCustodyAccount FeePaymentMethod = 0
)

type FeeRateInfo struct {
	gorm.Model
	Name                 string  `json:"name" gorm:"type:varchar(255);not null"`
	EstimateSmartFeeRate float64 `json:"estimate_smart_fee_rate"`
}
