package btldb

import (
	"sync"
	"trade/middleware"
	"trade/models/custodyModels"
)

var payOutsideMutex sync.Mutex

// CreatePayOutside creates a new payOutside
func CreatePayOutside(pay *custodyModels.PayOutside) error {
	payOutsideMutex.Lock()
	defer payOutsideMutex.Unlock()
	return middleware.DB.Create(pay).Error
}

// ReadPayOutside retrieves an payOutside by  Id
func ReadPayOutside(id uint) (*custodyModels.PayOutside, error) {
	var pay custodyModels.PayOutside
	err := middleware.DB.First(&pay, id).Error
	return &pay, err
}

func LoadPendingOutsides() (*[]custodyModels.PayOutside, error) {
	var pay []custodyModels.PayOutside
	err := middleware.DB.Where("status =?", custodyModels.PayOutsideStatusPending).Find(&pay).Error
	return &pay, err
}

// UpdatePayOutside updates an existing payOutside
func UpdatePayOutside(pay *custodyModels.PayOutside) error {
	payOutsideMutex.Lock()
	defer payOutsideMutex.Unlock()
	return middleware.DB.Save(pay).Error
}

func DeletePayOutside(id uint) error {
	var pay custodyModels.PayOutside
	return middleware.DB.Delete(&pay, id).Error
}
