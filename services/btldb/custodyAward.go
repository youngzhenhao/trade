package btldb

import (
	"sync"
	"trade/middleware"
	"trade/models"
)

var awardMutex sync.Mutex

// CreateAward creates a new Award
func CreateAward(award *models.AccountAward) error {
	awardMutex.Lock()
	defer awardMutex.Unlock()
	return middleware.DB.Create(award).Error
}

// ReadAward retrieves an Award by id
func ReadAward(id uint) (*models.AccountAward, error) {
	var award models.AccountAward
	err := middleware.DB.First(&award, id).Error
	return &award, err
}

// UpdateAward updates an existing Award
func UpdateAward(award *models.AccountAward) error {
	payOutsideMutex.Lock()
	defer payOutsideMutex.Unlock()
	return middleware.DB.Save(award).Error
}

func DeleteAward(id uint) error {
	var award models.AccountAward
	return middleware.DB.Delete(&award, id).Error
}
