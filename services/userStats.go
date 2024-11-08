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

func userToUserInfo(user models.User) models.StatsUserInfo {
	loginTime := utils.TimestampToTime(user.RecentLoginTime)
	var status string
	if user.Status == 1 {
		status = "正常"
	} else if user.Status == 0 {
		status = "已冻结"
	}
	return models.StatsUserInfo{
		ID:                user.ID,
		CreatedAt:         utils.TimeFormatCN(user.CreatedAt),
		UpdatedAt:         utils.TimeFormatCN(user.UpdatedAt),
		Username:          user.Username,
		Status:            status,
		RecentIpAddresses: user.RecentIpAddresses,
		RecentLoginTime:   utils.TimeFormatCN(loginTime),
	}
}

func usersToUserInfos(users *[]models.User) *[]models.StatsUserInfo {
	var userInfos []models.StatsUserInfo
	for _, user := range *users {
		userInfos = append(userInfos, userToUserInfo(user))
	}
	return &userInfos
}

//TODO: Mysql data to csv

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
	var _newUserToday *[]models.StatsUserInfo = nil
	var _dailyActiveUser *[]models.StatsUserInfo = nil
	var _monthlyActiveUser *[]models.StatsUserInfo = nil
	if _new {
		_newUserToday = usersToUserInfos(&newUserToday)
	}
	if _day {
		_dailyActiveUser = usersToUserInfos(&dailyActiveUser)
	}
	if _month {
		_monthlyActiveUser = usersToUserInfos(&monthlyActiveUser)
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
