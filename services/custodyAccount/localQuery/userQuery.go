package localQuery

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
)

var (
	DBError      = errors.New("database error")
	NotFoundUser = errors.New("not found User")
)

type BlockUserReq struct {
	Username []string `json:"username"`
	Memo     string   `json:"memo"`
}

func BlockUser(username, memo string) error {
	tx, back := middleware.GetTx()
	defer back()
	var err error
	var user models.User
	if err = tx.Where("user_name =?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NotFoundUser
		}
		return fmt.Errorf("%w: %v", DBError, err)
	}
	if user.Status != 0 {
		// user is already blocked
		return nil
	}
	user.Status = 1
	if err = tx.Save(&user).Error; err != nil {
		return fmt.Errorf("%w: %v", DBError, err)
	}
	//build blocked record
	record := custodyModels.BlockedRecord{
		UserId:      user.ID,
		BlockedType: custodyModels.BlockedUser,
		Memo:        memo,
	}
	if err = tx.Save(&record).Error; err != nil {
		return fmt.Errorf("%w: %v", DBError, err)
	}
	tx.Commit()
	return nil
}

type UnblockUserReq struct {
	Username []string `json:"username"`
	Memo     string   `json:"memo"`
}

func UnblockUser(username, memo string) error {
	tx, back := middleware.GetTx()
	defer back()
	var err error
	user := models.User{Username: username}
	if err = tx.Where("user_name =?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NotFoundUser
		}
		return fmt.Errorf("%w: %v", DBError, err)
	}
	if user.Status == 0 {
		// user is already blocked
		return nil
	}
	user.Status = 0
	if err = tx.Save(&user).Error; err != nil {
		return fmt.Errorf("%w: %v", DBError, err)
	}
	//build blocked record
	record := custodyModels.BlockedRecord{
		UserId:      user.ID,
		BlockedType: custodyModels.UnblockedUser,
		Memo:        memo,
	}
	if err = tx.Save(&record).Error; err != nil {
		return fmt.Errorf("%w: %v", DBError, err)
	}
	tx.Commit()
	return nil
}

type UserInfoRep struct {
	Username string `json:"username"`
}

type UserInfo struct {
	Npubkey         string `json:"npubkey"`
	Status          string `json:"status"`
	RecentIp        string `json:"recent_ip"`
	RecentLoginTime string `json:"recent_login_time"`
}

func GetUserInfo(username string) (*UserInfo, error) {
	tx, back := middleware.GetTx()
	defer back()
	var err error
	user := models.User{Username: username}
	if err = tx.Where(&user).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NotFoundUser
		}
		return nil, fmt.Errorf("%w: %v", DBError, err)
	}
	info := UserInfo{
		Npubkey:         user.Username,
		RecentIp:        user.RecentIpAddresses,
		RecentLoginTime: time.Unix(int64(user.RecentLoginTime), 0).Format("2006-01-02 15:04:05"),
	}
	if user.Status != 0 {
		info.Status = "blocked"
	} else {
		info.Status = "active"
	}
	return &info, nil
}
