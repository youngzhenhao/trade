package btldb

import (
	"gorm.io/gorm"
	"sync"
	"trade/middleware"
	"trade/models"
)

var balanceMutex sync.Mutex

// CreateBalance creates a new balance record
func CreateBalance(tx *gorm.DB, balance *models.Balance) error {
	balanceMutex.Lock()
	defer balanceMutex.Unlock()
	return tx.Create(balance).Error
}

// ReadBalance retrieves a balance by Id
func ReadBalance(id uint) (*models.Balance, error) {
	var balance models.Balance
	err := middleware.DB.First(&balance, id).Error
	return &balance, err
}

// UpdateBalance updates an existing balance
func UpdateBalance(tx *gorm.DB, balance *models.Balance) error {
	balanceMutex.Lock()
	defer balanceMutex.Unlock()
	return tx.Save(balance).Error
}

// DeleteBalance soft deletes a balance by Id
func DeleteBalance(id uint) error {
	var balance models.Balance
	return middleware.DB.Delete(&balance, id).Error
}
