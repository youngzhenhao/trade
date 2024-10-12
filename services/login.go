package services

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
	"trade/middleware"
	"trade/models"
	"trade/services/btldb"
)

const fixedSalt = "bitlongwallet7238baee9c2638664"

func SplitStringAndVerifyChecksum(extstring string) bool {
	originalString, checksum := spilt(extstring)
	if originalString == "" {
		return false
	}
	if checksum == "" {
		return false
	}
	return verifyChecksumWithSalt(originalString, checksum)
}

func spilt(extstring string) (string, string) {
	parts := strings.Split(extstring, "_e_")
	if len(parts) != 2 {
		return "", ""
	}
	originalString := parts[0]
	checksum := parts[1]
	return originalString, checksum
}

func generateMD5WithSalt(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input + fixedSalt))
	return hex.EncodeToString(hasher.Sum(nil))
}

func verifyChecksumWithSalt(originalString, checksum string) bool {
	expectedChecksum := generateMD5WithSalt(originalString)
	return checksum == expectedChecksum
}

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
	fmt.Println("login user:", user)
	if !CheckPassword(user.Password, creds.Password) {
		return "", errors.New("invalid credentials")
	}
	fmt.Println("CheckPassword user success:", user)
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
		originalString, _ := spilt(creds.Password)
		if originalString != "" {
			password, err := hashPassword(originalString)
			if err != nil {
				return "", err
			}
			user.Password = password
			err = btldb.UpdateUser(&user)
			if err != nil {
				return "", err
			}
		}
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
