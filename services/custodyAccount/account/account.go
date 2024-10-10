package account

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"trade/models"
	cModels "trade/models/custodyModels"
	"trade/services/btldb"
)

type UserInfo struct {
	User        *models.User
	Account     *models.Account
	LockAccount *cModels.LockAccount
}

// GetUserInfo 获取用户信息
func GetUserInfo(username string) (*UserInfo, error) {
	// 获取用户信息
	user, err := btldb.ReadUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", models.ReadDbErr, err)
	}
	// 获取Lit账户信息
	account := &models.Account{}
	account, err = GetAccountByUserName(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", models.ReadDbErr, err)
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果账户不存在，则创建账户
		account, err = CreateAccount(user)
		if err != nil {
			return nil, CustodyAccountCreateErr
		}
	}
	// 获取冻结账户信息
	lockAccount := &cModels.LockAccount{}
	lockAccount, err = GetLockAccountByUserName(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", models.ReadDbErr, err)
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果账户不存在，则创建账户
		lockAccount, err = CreateLockAccount(user)
		if err != nil {
			return nil, CustodyAccountCreateErr
		}
	}

	return &UserInfo{
		User:        user,
		Account:     account,
		LockAccount: lockAccount,
	}, nil
}

// GetUserInfoById 获取用户信息
func GetUserInfoById(userId uint) (*UserInfo, error) {
	user, err := btldb.ReadUser(userId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", models.ReadDbErr, err)
	}
	// 获取Lit账户信息
	account := &models.Account{}
	account, err = GetAccountByUserName(user.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", models.ReadDbErr, err)
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果账户不存在，则创建账户
		account, err = CreateAccount(user)
		if err != nil {
			return nil, CustodyAccountCreateErr
		}
	}
	// 获取冻结账户信息
	lockAccount := &cModels.LockAccount{}
	lockAccount, err = GetLockAccountByUserName(user.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", models.ReadDbErr, err)
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果账户不存在，则创建账户
		lockAccount, err = CreateLockAccount(user)
		if err != nil {
			return nil, CustodyAccountCreateErr
		}
	}
	return &UserInfo{
		User:        user,
		Account:     account,
		LockAccount: lockAccount,
	}, nil
}
