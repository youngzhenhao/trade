package btldb

import (
	"trade/middleware"
	"trade/models"
)

// CreateLoginRecord creates a new user login record
func CreateLoginRecord(loginRecord *models.LoginRecord) error {
	return middleware.DB.Create(loginRecord).Error
}
