package btldb

import (
	"sync"
	"trade/middleware"
	"trade/models"
)

var payOutsideTxMutex sync.Mutex

// CreatePayOutsideTx creates a new payOutside
func CreatePayOutsideTx(pay *models.PayOutsideTx) error {
	payOutsideTxMutex.Lock()
	defer payOutsideTxMutex.Unlock()
	return middleware.DB.Create(pay).Error
}

// ReadPayOutsideTx retrieves an payOutside by  Id
func ReadPayOutsideTx(id uint) (*models.PayOutsideTx, error) {
	var pay models.PayOutsideTx
	err := middleware.DB.First(&pay, id).Error
	return &pay, err
}

// UpdatePayOutsideTx updates an existing payOutside
func UpdatePayOutsideTx(pay *models.PayOutsideTx) error {
	payOutsideTxMutex.Lock()
	defer payOutsideTxMutex.Unlock()
	return middleware.DB.Save(pay).Error
}

func DeletePayOutsideTx(id uint) error {
	var pay models.PayOutsideTx
	return middleware.DB.Delete(&pay, id).Error
}
