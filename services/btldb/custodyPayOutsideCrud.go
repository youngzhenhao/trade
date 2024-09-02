package btldb

import (
	"sync"
	"trade/middleware"
	"trade/models"
)

var payOutsideMutex sync.Mutex

// CreatePayOutside creates a new payOutside
func CreatePayOutside(pay *models.PayOutside) error {
	payOutsideMutex.Lock()
	defer payOutsideMutex.Unlock()
	return middleware.DB.Create(pay).Error
}

// ReadPayOutside retrieves an payOutside by  Id
func ReadPayOutside(id uint) (*models.PayOutside, error) {
	var pay models.PayOutside
	err := middleware.DB.First(&pay, id).Error
	return &pay, err
}

// UpdatePayOutside updates an existing payOutside
func UpdatePayOutside(pay *models.PayOutside) error {
	payOutsideMutex.Lock()
	defer payOutsideMutex.Unlock()
	return middleware.DB.Save(pay).Error
}

func DeletePayOutside(id uint) error {
	var pay models.PayOutside
	return middleware.DB.Delete(&pay, id).Error
}
