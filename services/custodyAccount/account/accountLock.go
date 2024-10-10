package account

import (
	"trade/middleware"
	"trade/models"
	cModels "trade/models/custodyModels"
)

func CreateLockAccount(user *models.User) (*cModels.LockAccount, error) {
	tx := middleware.DB.Begin()
	account := cModels.LockAccount{
		UserId:   user.ID,
		UserName: user.Username,
		Status:   cModels.AccountStatusEnable,
	}
	if err := tx.Create(&account).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &account, nil
}

func GetLockAccountByUserName(username string) (*cModels.LockAccount, error) {
	tx := middleware.DB.Begin()
	account := cModels.LockAccount{}
	if err := tx.Where("user_name =?", username).First(&account).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return &account, nil
}
