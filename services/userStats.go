package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"math"
	"os"
	"path/filepath"
	"strconv"
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
		status = "已冻结"
	} else if user.Status == 0 {
		status = "正常"
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

func GetDailyActiveUser(year int, month int, day int) (int64, error) {
	var dailyActiveUserNum int64
	now := time.Now()
	if year > now.Year() {
		return 0, errors.New("invalid year")
	}
	if month > 12 || month < 1 {
		return 0, errors.New("invalid month")
	}
	if day > 31 || day < 1 {
		return 0, errors.New("invalid day")
	}
	date := fmt.Sprintf("%4d-%02d-%02d", year, month, day)
	err := middleware.DB.Model(&models.LoginRecord{}).
		Where("DATE(created_at) = ?", date).
		Distinct("user_id").
		Count(&dailyActiveUserNum).Error
	if err != nil {
		return 0, err
	}
	return dailyActiveUserNum, nil
}

func GetMonthlyActiveUser(year int, month int) (int64, error) {
	var monthlyActiveUser int64
	err := middleware.DB.Model(&models.LoginRecord{}).
		Where("YEAR(created_at) = ? AND MONTH(created_at) = ?", year, month).
		Distinct("user_id").
		Count(&monthlyActiveUser).Error
	if err != nil {
		return 0, err
	}
	return monthlyActiveUser, nil
}

type LoginRecordInfo struct {
	ID                uint   `json:"id"`
	UserId            uint   `json:"user_id"`
	RecentIpAddresses string `json:"recent_ip_addresses"`
	Path              string `json:"path"`
	LoginTime         int    `json:"login_time"`
}

func GetActiveUserBetween(start string, end string) (*[]LoginRecordInfo, error) {
	if len(start) != len(time.DateOnly) {
		return nil, errors.New("invalid start time length(" + strconv.Itoa(len(start)) + "), should be like " + time.DateOnly)
	}
	if len(end) != len(time.DateOnly) {
		return nil, errors.New("invalid end time length(" + strconv.Itoa(len(end)) + "), should be like " + time.DateOnly)
	}
	_start, err := time.Parse(time.DateOnly, start)
	if err != nil {
		return nil, err
	}
	_end, err := time.Parse(time.DateOnly, end)
	if err != nil {
		return nil, err
	}
	var loginRecordInfos []LoginRecordInfo
	// @dev: Do not select times
	err = middleware.DB.Model(&models.LoginRecord{}).
		Select("id, user_id, recent_ip_addresses, path, login_time").
		Where("path = ? and login_time between ? and ?", "/login", _start.Unix(), _end.Unix()).
		Order("login_time desc").
		Scan(&loginRecordInfos).Error
	if err != nil {
		return nil, err
	}
	return &loginRecordInfos, nil
}

func GetSpecifiedDateUserStats(day string) (*models.UserStats, error) {
	specifiedDay, err := utils.DateTimeStringToTimeWithFormat(day, "20060102")
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "DateTimeStringToTimeWithFormat")
	}
	specifiedDate := specifiedDay.Format("2006年01月02日")
	var users []models.User
	var newUserToday []models.User
	errorInfos := new([]string)
	now := time.Now()
	err = middleware.DB.Find(&users).Error
	if err != nil {
		*errorInfos = append(*errorInfos, err.Error())
	} else {
		for _, user := range users {
			if user.CreatedAt.Year() == specifiedDay.Year() && user.CreatedAt.Month() == specifiedDay.Month() && user.CreatedAt.Day() == specifiedDay.Day() {
				newUserToday = append(newUserToday, user)
			}
		}
	}
	var dailyActiveUserNum int64
	var monthlyActiveUserNum int64
	dailyActiveUserNum, err = GetDailyActiveUser(specifiedDay.Year(), int(specifiedDay.Month()), specifiedDay.Day())
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetDailyActiveUser")
	}
	monthlyActiveUserNum, err = GetMonthlyActiveUser(specifiedDay.Year(), int(specifiedDay.Month()))
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetMonthlyActiveUser")
	}
	return &models.UserStats{
		QueryTime:            utils.TimeFormatCN(now),
		TotalUser:            uint64(len(users)),
		SpecifiedDate:        specifiedDate,
		NewUserTodayNum:      uint64(len(newUserToday)),
		DailyActiveUserNum:   uint64(dailyActiveUserNum),
		MonthlyActiveUserNum: uint64(monthlyActiveUserNum),
		ErrorInfos:           errorInfos,
	}, nil
}

func GetSpecifiedDateUserStatsYaml(day string) (string, error) {
	userStats, err := GetSpecifiedDateUserStats(day)
	if err != nil {
		return "", utils.AppendErrorInfo(err, "GetUserStats")
	}
	userStatsBytes, _ := yaml.Marshal(userStats)
	return string(userStatsBytes), nil

}

func GetUserStatsYaml(_new bool, _day bool, _month bool) (string, error) {
	userStats, err := GetUserStats(_new, _day, _month)
	if err != nil {
		return "", utils.AppendErrorInfo(err, "GetUserStats")
	}
	userStatsBytes, _ := yaml.Marshal(userStats)
	return string(userStatsBytes), nil
}

type StatsUserInfo struct {
	ID                uint   `json:"用户ID" yaml:"用户ID"`
	CreatedAt         string `json:"用户创建时间" yaml:"用户创建时间"`
	UpdatedAt         string `json:"更新时间" yaml:"更新时间"`
	Username          string `json:"用户名;NpubKey;Nostr地址" yaml:"用户名;NpubKey;Nostr地址"`
	Status            string `json:"用户状态" yaml:"用户状态"`
	RecentIpAddresses string `json:"最近IP地址" yaml:"最近IP地址"`
	RecentLoginTime   string `json:"最近登录时间" yaml:"最近登录时间"`
}

func StatsUserInfoToCsv(filename string, statsUserInfos *[]models.StatsUserInfo) (string, error) {
	if statsUserInfos == nil {
		return "", errors.New("statsUserInfos is nil")
	}
	var path = filepath.Join(".", filepath.Join("csv", filename+"-user.csv"))
	utils.CreateFile(path, "")
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.Write([]string{"用户ID", "用户创建时间", "更新时间", "用户名;NpubKey;Nostr地址", "用户状态", "最近IP地址", "最近登录时间"})
	if err != nil {
		return "", err
	}
	for _, statsUserInfo := range *statsUserInfos {
		err = writer.Write([]string{
			strconv.Itoa(int(statsUserInfo.ID)),
			statsUserInfo.CreatedAt,
			statsUserInfo.UpdatedAt,
			statsUserInfo.Username,
			statsUserInfo.Status,
			statsUserInfo.RecentIpAddresses,
			statsUserInfo.RecentLoginTime,
		})
		if err != nil {
			return "", err
		}
	}
	return path, nil
}

func GetActiveUserCountBetween(start string, end string) (int64, error) {
	if len(start) != len(time.DateOnly) {
		return 0, errors.New("invalid start time length(" + strconv.Itoa(len(start)) + "), should be like " + time.DateOnly)
	}
	if len(end) != len(time.DateOnly) {
		return 0, errors.New("invalid end time length(" + strconv.Itoa(len(end)) + "), should be like " + time.DateOnly)
	}
	var activeUser int64
	err := middleware.DB.Model(&models.LoginRecord{}).
		Where("created_at between ? and ?", start, end).
		Distinct("user_id").
		Count(&activeUser).Error
	if err != nil {
		return 0, err
	}
	return activeUser, nil
}

func GetDateLoginCount(start string, end string) (int64, error) {
	if len(start) != len(time.DateOnly) {
		return 0, errors.New("invalid start time length(" + strconv.Itoa(len(start)) + "), should be like " + time.DateOnly)
	}
	if len(end) != len(time.DateOnly) {
		return 0, errors.New("invalid end time length(" + strconv.Itoa(len(end)) + "), should be like " + time.DateOnly)
	}
	var count int64
	err := middleware.DB.Model(&models.DateLogin{}).
		Where("date between ? and ?", start, end).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

type UserActiveResult struct {
	UserName          string
	RecentIpAddresses string
}

func GetUserActiveRecord(start string, end string, limit int, offset int) (*[]UserActiveResult, error) {
	if len(start) != len(time.DateOnly) {
		return nil, errors.New("invalid start time length(" + strconv.Itoa(len(start)) + "), should be like " + time.DateOnly)
	}
	if len(end) != len(time.DateOnly) {
		return nil, errors.New("invalid end time length(" + strconv.Itoa(len(end)) + "), should be like " + time.DateOnly)
	}
	var userActiveResults []UserActiveResult
	err := middleware.DB.Table("login_record").
		Select("user.user_name, login_record.recent_ip_addresses").
		Joins("JOIN user ON user.id = login_record.user_id").
		Where("login_record.created_at BETWEEN ? AND ?", start, end).
		Group("user.user_name, login_record.recent_ip_addresses").
		Limit(limit).
		Offset(offset).
		Scan(&userActiveResults).Error
	if err != nil {
		return nil, err
	}
	return &userActiveResults, nil
}

func GetUserActiveRecordNum(start string, end string) (int64, error) {
	if len(start) != len(time.DateOnly) {
		return 0, errors.New("invalid start time length(" + strconv.Itoa(len(start)) + "), should be like " + time.DateOnly)
	}
	if len(end) != len(time.DateOnly) {
		return 0, errors.New("invalid end time length(" + strconv.Itoa(len(end)) + "), should be like " + time.DateOnly)
	}
	var recordNum int64
	err := middleware.DB.Table("login_record").
		Select("user.user_name, login_record.recent_ip_addresses").
		Joins("JOIN user ON user.id = login_record.user_id").
		Where("login_record.created_at BETWEEN ? AND ?", start, end).
		Group("user.user_name, login_record.recent_ip_addresses").
		Count(&recordNum).Error
	if err != nil {
		return 0, err
	}
	return recordNum, nil
}

type DateIpLoginRecord struct {
	Username string `json:"username"`
	Date     string `json:"date"`
	Ip       string `json:"ip"`
}

func GetDateIpLoginRecord(start string, end string, limit int, offset int) (*[]DateIpLoginRecord, error) {
	var dateIpLoginRecords []DateIpLoginRecord
	if len(start) != len(time.DateOnly) {
		return &dateIpLoginRecords, errors.New("invalid start time length(" + strconv.Itoa(len(start)) + "), should be like " + time.DateOnly)
	}
	if len(end) != len(time.DateOnly) {
		return &dateIpLoginRecords, errors.New("invalid end time length(" + strconv.Itoa(len(end)) + "), should be like " + time.DateOnly)
	}
	err := middleware.DB.Model(models.DateIpLogin{}).
		Where("date between ? and ?", start, end).
		Limit(limit).
		Offset(offset).
		Order("id desc").
		Scan(&dateIpLoginRecords).Error
	if err != nil {
		return &dateIpLoginRecords, err
	}
	return &dateIpLoginRecords, nil
}

func GetDateIpLoginRecordAll(start string, end string) (*[]DateIpLoginRecord, error) {
	var dateIpLoginRecords []DateIpLoginRecord
	if len(start) != len(time.DateOnly) {
		return &dateIpLoginRecords, errors.New("invalid start time length(" + strconv.Itoa(len(start)) + "), should be like " + time.DateOnly)
	}
	if len(end) != len(time.DateOnly) {
		return &dateIpLoginRecords, errors.New("invalid end time length(" + strconv.Itoa(len(end)) + "), should be like " + time.DateOnly)
	}
	err := middleware.DB.Model(models.DateIpLogin{}).
		Where("date between ? and ?", start, end).
		Order("id desc").
		Scan(&dateIpLoginRecords).Error
	if err != nil {
		return &dateIpLoginRecords, err
	}
	return &dateIpLoginRecords, nil
}

func GetDateIpLoginCount(start string, end string) (int64, error) {
	if len(start) != len(time.DateOnly) {
		return 0, errors.New("invalid start time length(" + strconv.Itoa(len(start)) + "), should be like " + time.DateOnly)
	}
	if len(end) != len(time.DateOnly) {
		return 0, errors.New("invalid end time length(" + strconv.Itoa(len(end)) + "), should be like " + time.DateOnly)
	}
	var count int64
	err := middleware.DB.Model(&models.DateIpLogin{}).
		Where("date between ? and ?", start, end).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetDateIpLoginPageNumber(start string, end string, size int) (int, error) {
	_count, err := GetDateIpLoginCount(start, end)
	if err != nil {
		return 0, err
	}
	if _count > math.MaxInt {
		return 0, errors.New("count overflows")
	}
	var pageNumber int
	count := int(_count)
	pageNumber = count / size
	if count%size != 0 {
		pageNumber++
	}
	return pageNumber, nil
}

func PageAndSizeToLimitAndOffset(page uint, size uint) (limit uint, offset uint) {
	return size, (page - 1) * size
}

func GetNewUserCount(start string, end string) (int64, error) {
	if len(start) != len(time.DateOnly) {
		return 0, errors.New("invalid start time length(" + strconv.Itoa(len(start)) + "), should be like " + time.DateOnly)
	}
	if len(end) != len(time.DateOnly) {
		return 0, errors.New("invalid end time length(" + strconv.Itoa(len(end)) + "), should be like " + time.DateOnly)
	}
	var count int64
	err := middleware.DB.Model(&models.User{}).
		Where("created_at between ? and ?", start, end).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetNewUserPageNumber(start string, end string, size int) (int, error) {
	_count, err := GetNewUserCount(start, end)
	if err != nil {
		return 0, err
	}
	if _count > math.MaxInt {
		return 0, errors.New("count overflows")
	}
	var pageNumber int
	count := int(_count)
	pageNumber = count / size
	if count%size != 0 {
		pageNumber++
	}
	return pageNumber, nil
}

type userInfo struct {
	CreatedAt         time.Time
	Username          string `gorm:"unique;column:user_name;type:varchar(255)" json:"userName"` // 正确地将unique和column选项放在同一个gorm标签内
	RecentIpAddresses string `json:"recent_ip_addresses" gorm:"type:varchar(255)"`
}

func UserToDateIpLoginRecord(user userInfo) DateIpLoginRecord {
	return DateIpLoginRecord{
		Username: user.Username,
		Date:     user.CreatedAt.Format(time.DateOnly),
		Ip:       user.RecentIpAddresses,
	}
}

func UsersToDateIpLoginRecords(users *[]userInfo) *[]DateIpLoginRecord {
	if users == nil {
		return new([]DateIpLoginRecord)
	}
	var dateIpLoginRecords []DateIpLoginRecord
	for _, user := range *users {
		dateIpLoginRecords = append(dateIpLoginRecords, UserToDateIpLoginRecord(user))
	}
	return &dateIpLoginRecords
}

func GetNewUserRecord(start string, end string, limit int, offset int) (*[]DateIpLoginRecord, error) {
	var dateIpLoginRecords = new([]DateIpLoginRecord)
	if len(start) != len(time.DateOnly) {
		return dateIpLoginRecords, errors.New("invalid start time length(" + strconv.Itoa(len(start)) + "), should be like " + time.DateOnly)
	}
	if len(end) != len(time.DateOnly) {
		return dateIpLoginRecords, errors.New("invalid end time length(" + strconv.Itoa(len(end)) + "), should be like " + time.DateOnly)
	}
	var users []userInfo
	err := middleware.DB.Model(models.User{}).
		Where("created_at between ? and ?", start, end).
		Limit(limit).
		Offset(offset).
		Order("id desc").
		Scan(&users).Error
	dateIpLoginRecords = UsersToDateIpLoginRecords(&users)
	if err != nil {
		return dateIpLoginRecords, err
	}
	return dateIpLoginRecords, nil
}

func GetNewUserRecordAll(start string, end string) (*[]DateIpLoginRecord, error) {
	var dateIpLoginRecords = new([]DateIpLoginRecord)
	if len(start) != len(time.DateOnly) {
		return dateIpLoginRecords, errors.New("invalid start time length(" + strconv.Itoa(len(start)) + "), should be like " + time.DateOnly)
	}
	if len(end) != len(time.DateOnly) {
		return dateIpLoginRecords, errors.New("invalid end time length(" + strconv.Itoa(len(end)) + "), should be like " + time.DateOnly)
	}
	var users []userInfo
	err := middleware.DB.Model(models.User{}).
		Where("created_at between ? and ?", start, end).
		Order("id desc").
		Scan(&users).Error
	dateIpLoginRecords = UsersToDateIpLoginRecords(&users)
	if err != nil {
		return dateIpLoginRecords, err
	}
	return dateIpLoginRecords, nil
}

func DateIpLoginRecordToCsv(filename string, dateIpLoginRecords *[]DateIpLoginRecord) (string, error) {
	if dateIpLoginRecords == nil {
		return "", errors.New("dateIpLoginRecords is nil")
	}
	var path = filepath.Join(".", filepath.Join("csv", filename+"-user.csv"))
	utils.CreateFile(path, "")
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.Write([]string{"用户名", "日期", "IP"})
	if err != nil {
		return "", err
	}
	for _, statsUserInfo := range *dateIpLoginRecords {
		err = writer.Write([]string{
			statsUserInfo.Username,
			statsUserInfo.Date,
			statsUserInfo.Ip,
		})
		if err != nil {
			return "", err
		}
	}
	return path, nil
}
