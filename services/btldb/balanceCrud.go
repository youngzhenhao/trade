package btldb

import (
	"sync"
	"trade/middleware"
	"trade/models"
)

var balanceMutex sync.Mutex

// CreateBalance creates a new balance record
func CreateBalance(balance *models.Balance) error {
	balanceMutex.Lock()
	defer balanceMutex.Unlock()
	return middleware.DB.Create(balance).Error
}

// ReadBalance retrieves a balance by Id
func ReadBalance(id uint) (*models.Balance, error) {
	var balance models.Balance
	err := middleware.DB.First(&balance, id).Error
	return &balance, err
}

// UpdateBalance updates an existing balance
func UpdateBalance(balance *models.Balance) error {
	balanceMutex.Lock()
	defer balanceMutex.Unlock()
	return middleware.DB.Save(balance).Error
}

// DeleteBalance soft deletes a balance by Id
func DeleteBalance(id uint) error {
	var balance models.Balance
	return middleware.DB.Delete(&balance, id).Error
}
