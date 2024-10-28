package account

import (
	"fmt"
	"time"
	"trade/models"
	"trade/services/btldb"
)

func GetUserInfo(userName string) (*UserInfo, error) {
	var user *UserInfo
	user, exists := pool.GetUser(userName)
	if exists {
		user.LastActiveMux.Lock()
		defer user.LastActiveMux.Unlock()

		if time.Since(user.LastActive) > time.Minute {
			pool.UpdateUserActivity(userName)
		}
		return user, nil
	}
	user, err := pool.CreateUser(userName)
	if err != nil {
		//todo :分析usr是否存在于数据库
		return nil, fmt.Errorf("create user %s failed: %w", userName, err)
	}
	return user, nil
}

func GetUserInfoById(userId uint) (*UserInfo, error) {
	dbUser, err := btldb.ReadUser(userId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", models.ReadDbErr, err)
	}
	return GetUserInfo(dbUser.Username)
}
