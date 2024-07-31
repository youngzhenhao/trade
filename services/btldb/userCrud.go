package btldb

import (
	"trade/middleware"
	"trade/models"
)

// CreateUser creates a new user record
func CreateUser(user *models.User) error {
	return middleware.DB.Create(user).Error
}

// ReadUser retrieves a user by Id
func ReadUser(id uint) (*models.User, error) {
	var user models.User
	err := middleware.DB.First(&user, id).Error
	return &user, err
}

func ReadAllUser() (*[]models.User, error) {
	var users []models.User
	err := middleware.DB.Find(&users).Error
	return &users, err
}

// ReadUserByUsername retrieves a user by username
func ReadUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := middleware.DB.Where("user_name = ?", username).First(&user).Error
	return &user, err
}

// UpdateUser updates an existing user
func UpdateUser(user *models.User) error {
	return middleware.DB.Save(user).Error
}

// DeleteUser soft deletes a user by Id
func DeleteUser(id uint) error {
	var user models.User
	return middleware.DB.Delete(&user, id).Error
}
