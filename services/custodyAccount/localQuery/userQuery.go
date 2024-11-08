package localQuery

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"trade/middleware"
	"trade/models"
)

var (
	DBError      = errors.New("database error")
	NotFoundUser = errors.New("not found User")
)

func BlockUser(username string) error {
	tx, back := middleware.GetTx()
	defer back()
	var err error
	user := models.User{Username: username}
	if err = tx.Where(&user).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NotFoundUser
		}
		return fmt.Errorf("%w: %v", DBError, err)
	}
	if user.Status != 0 {
		// user is already blocked
		return fmt.Errorf("user is already blocked")
	}
	user.Status = 1
	if err = tx.Save(&user).Error; err != nil {
		return fmt.Errorf("%w: %v", DBError, err)
	}
	tx.Commit()
	return nil
}

func UnblockUser(username, memo string) error {
	tx, back := middleware.GetTx()
	defer back()
	var err error
	user := models.User{Username: username}
	if err = tx.Where(&user).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NotFoundUser
		}
		return fmt.Errorf("%w: %v", DBError, err)
	}
	if user.Status != 0 {
		// user is already blocked
		return fmt.Errorf("user is already blocked")
	}
	user.Status = 1
	if err = tx.Save(&user).Error; err != nil {
		return fmt.Errorf("%w: %v", DBError, err)
	}
	tx.Commit()
	return nil
}

func GetUserInfo(userID string) {}
