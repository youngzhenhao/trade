package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"trade/middleware"
	"trade/models"
	"trade/services/btldb"
)

func Login(creds models.User) (string, error) {
	var user models.User
	result := middleware.DB.Where("user_name = ?", creds.Username).First(&user)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// If there are other database errors, an error is returned
			return "", result.Error
		} else {
			user.Username = creds.Username
			password, err := hashPassword(creds.Password)
			if err != nil {
				return "", err
			}
			user.Password = password
			err = btldb.CreateUser(&user)
			if err != nil {
				return "", err
			}
		}
	}
	if !CheckPassword(user.Password, creds.Password) {
		return "", errors.New("invalid credentials")
	}
	token, err := middleware.GenerateToken(creds.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidateUserAndGenerateToken(creds models.User) (string, error) {
	var user models.User
	result := middleware.DB.Where("user_name = ?", creds.Username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", errors.New("invalid credentials")
	}
	if !CheckPassword(user.Password, creds.Password) {
		return "", errors.New("invalid credentials")
	}
	token, err := middleware.GenerateToken(creds.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (cs *CronService) FiveSecondTask() {
	fmt.Println("5 secs runs")
	log.Println("5 secs runs")
}
