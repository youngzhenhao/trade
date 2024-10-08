package btldb

import (
	"trade/middleware"
	"trade/models"
)

// CreateAccount creates a new account
func CreateAccount(account *models.Account) error {
	return middleware.DB.Create(account).Error
}

// ReadAccount retrieves an account by accountId
func ReadAccount(id uint) (*models.Account, error) {
	var account models.Account
	err := middleware.DB.First(&account, id).Error
	return &account, err
}

// ReadAccountByName  retrieves an account by Username
func ReadAccountByName(name string) (*models.Account, error) {
	var account models.Account
	err := middleware.DB.Where("user_name =?", name).First(&account).Error
	return &account, err
}

// ReadAccountByUserId retrieves an account by userId
func ReadAccountByUserId(userId uint) (*models.Account, error) {
	var account models.Account
	err := middleware.DB.Where("user_id =?", userId).First(&account).Error
	return &account, err
}

// UpdateAccount updates an existing account
func UpdateAccount(account *models.Account) error {
	return middleware.DB.Save(account).Error
}

func DeleteAccount(id uint) error {
	var account models.Account
	return middleware.DB.Delete(&account, id).Error
}
