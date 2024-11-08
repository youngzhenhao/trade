package services

import (
	"gopkg.in/yaml.v3"
	"time"
	"trade/middleware"
	"trade/models"
	"trade/utils"
)

func GetTotalUserNumber() (uint64, error) {
	type onlyId struct {
		Id uint `json:"id"`
	}
	var ids []onlyId
	err := middleware.DB.Table("user").Select("id").Scan(&ids).Error
	if err != nil {
		return 0, err
	}
	return uint64(len(ids)), nil
}

type userInfo struct {
	ID                uint   `json:"用户ID" yaml:"用户ID"`
	CreatedAt         string `json:"用户创建时间" yaml:"用户创建时间"`
	UpdatedAt         string `json:"更新时间" yaml:"更新时间"`
	Username          string `json:"用户名;NpubKey;Nostr地址" yaml:"用户名;NpubKey;Nostr地址"`
	Status            int16  `json:"用户状态" yaml:"用户状态"`
	RecentIpAddresses string `json:"最近IP地址" yaml:"最近IP地址"`
	RecentLoginTime   string `json:"最近登录时间" yaml:"最近登录时间"`
}

func userToUserInfo(user *models.User) *userInfo {
	if user == nil {
		return nil
	}
	loginTime := utils.TimestampToTime(user.RecentLoginTime)
	return &userInfo{
		ID:                user.ID,
		CreatedAt:         utils.TimeFormatCN(user.CreatedAt),
		UpdatedAt:         utils.TimeFormatCN(user.UpdatedAt),
		Username:          user.Username,
		Status:            user.Status,
		RecentIpAddresses: user.RecentIpAddresses,
		RecentLoginTime:   utils.TimeFormatCN(loginTime),
	}
}

func GetUserStats(_new bool, _day bool, _month bool) (*models.UserStats, error) {
	var users []models.User
	var newUserToday []models.User
	var dailyActiveUser []models.User
	var monthlyActiveUser []models.User
	errorInfos := new([]string)
	now := time.Now()
	err := middleware.DB.Find(&users).Error
	if err != nil {
		*errorInfos = append(*errorInfos, err.Error())
	} else {
		for _, user := range users {
			if user.CreatedAt.Year() == now.Year() && user.CreatedAt.Month() == now.Month() && user.CreatedAt.Day() == now.Day() {
				newUserToday = append(newUserToday, user)
			}
			if user.UpdatedAt.Year() == now.Year() {
				if user.UpdatedAt.Month() == now.Month() {
					monthlyActiveUser = append(monthlyActiveUser, user)
					if user.UpdatedAt.Day() == now.Day() {
						dailyActiveUser = append(dailyActiveUser, user)
					}
				}
			}
		}
	}
	var _newUserToday *[]models.User = nil
	var _dailyActiveUser *[]models.User = nil
	var _monthlyActiveUser *[]models.User = nil
	if _new {
		_newUserToday = &newUserToday
	}
	if _day {
		_dailyActiveUser = &dailyActiveUser
	}
	if _month {
		_monthlyActiveUser = &monthlyActiveUser
	}
	return &models.UserStats{
		QueryTime:            utils.TimeFormatCN(now),
		TotalUser:            uint64(len(users)),
		NewUserTodayNum:      uint64(len(newUserToday)),
		DailyActiveUserNum:   uint64(len(dailyActiveUser)),
		MonthlyActiveUserNum: uint64(len(monthlyActiveUser)),
		ErrorInfos:           errorInfos,
		NewUserToday:         _newUserToday,
		DailyActiveUser:      _dailyActiveUser,
		MonthlyActiveUser:    _monthlyActiveUser,
	}, nil
}

func GetUserStatsYaml(_new bool, _day bool, _month bool) (string, error) {
	userStats, err := GetUserStats(_new, _day, _month)
	if err != nil {
		return "", utils.AppendErrorInfo(err, "GetUserStats")
	}
	userStatsBytes, _ := yaml.Marshal(userStats)
	return string(userStatsBytes), nil
}
