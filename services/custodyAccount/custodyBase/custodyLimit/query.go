package custodyLimit

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"time"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
)

var limitSync = sync.Mutex{}

type UserLimit struct {
	UserName          string  `json:"userName" gorm:"column:user_name"`
	LimitType         string  `json:"limit_type" gorm:"column:memo"`
	Level             int     `json:"level" gorm:"column:level"`
	TodayAmount       float64 `json:"todayAmount" gorm:"column:total_amount"`
	TodayCount        int     `json:"todayCount" gorm:"column:total_count"`
	TodayUsefulAmount float64 `json:"todayUsefulAmount" gorm:"column:use_able_amount"`
	TodayUsefulCount  int     `json:"todayUsefulCount" gorm:"column:use_able_count"`
}

func GetUserLimit(userName, LimitType string, page, pageSize int) (int64, *[]UserLimit, error) {
	if !limitSync.TryLock() {
		time.Sleep(5 * time.Second)
		return 0, nil, fmt.Errorf("当前有其他请求正在处理，请稍后再试")
	}
	defer limitSync.Unlock()

	db := middleware.DB
	var userLimits []UserLimit
	GetUserTypeLimitMap()
	q := db.Table("user_limit")
	if userName != "" {
		q = q.Where("user.user_name =?", userName)
	}
	if LimitType != "" {
		if value, exists := typeMapInt[LimitType]; exists {
			q = q.Where("user_limit.limit_type =?", value)
		}
	}
	// 获取今日0点的时间
	todayStart := time.Now().Truncate(24 * time.Hour)

	// 增加一个条件，筛选今日0点开始的记录
	q = q.Where("user_limit.created_at >= ?", todayStart)

	q = q.Debug().Table("user_limit").
		Joins("left join user on user.id = user_limit.user_id").
		Joins("left join (select * from user_limit_bills where created_at >= ?) as bill on bill.user_id = user_limit.user_id", todayStart).
		Joins("left join user_limit_type on user_limit_type.id = user_limit.limit_type")
	var total int64
	err := q.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}
	err = q.Select("user.user_name, " +
		"user_limit_type.memo," +
		"user_limit.level," +
		"bill.total_amount,bill.use_able_amount,bill.total_count,bill.use_able_count").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&userLimits).Error
	if err != nil {
		return 0, nil, err
	}
	return total, &userLimits, nil
}

func SetUserLimitLevel(userName, limitType string, level int) error {
	limitSync.Lock()
	defer limitSync.Unlock()

	limitMux.Lock()
	defer limitMux.Unlock()

	db := middleware.DB
	GetUserTypeLimitMap()

	TypeID, exists := typeMapInt[limitType]
	if !exists {
		return fmt.Errorf("限制类型不存在: %s", limitType)
	}
	usr := models.User{}
	err := db.Where("user_name =?", userName).First(&usr).Error
	if err != nil {
		return err
	}
	if level <= 0 {
		level = 0
	}
	userLimit := custodyModels.Limit{}
	err = db.Where("user_id =? and limit_type =?", usr.ID, TypeID).First(&userLimit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userLimit.Level = uint(level)
			return db.Create(&userLimit).Error
		}
		return err
	}
	userLimit.Level = uint(level)
	return db.Save(&userLimit).Error
}

func SetUserTodayLimit(userName, limitType string, amount int, count int) error {
	limitSync.Lock()
	defer limitSync.Unlock()

	limitMux.Lock()
	defer limitMux.Unlock()

	db := middleware.DB
	GetUserTypeLimitMap()
	TypeID, exists := typeMapInt[limitType]
	if !exists {
		return fmt.Errorf("限制类型不存在: %s", limitType)
	}
	usr := models.User{}
	err := db.Where("user_name =?", userName).First(&usr).Error
	if err != nil {
		return err
	}

	// 获取今日0点的时间
	todayStart := time.Now().Truncate(24 * time.Hour)

	var bill custodyModels.LimitBill
	// 增加一个条件，筛选今日0点开始的记录
	q := db.Where("user_limit.created_at >= ?", todayStart)
	err = q.Where("user_limit.user_id =? and user_limit.limit_type =?", usr.ID, TypeID).First(&bill).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			bill.UserId = usr.ID
			bill.TotalAmount = float64(amount)
			bill.UseAbleAmount = float64(amount)
			bill.TotalCount = uint(count)
			bill.TotalAmount = float64(count)
		}
		return err
	}
	if amount >= 0 {
		bill.UseAbleAmount = float64(amount)
	}
	if count >= 0 {
		bill.UseAbleCount = uint(count)
	}
	return db.Save(&bill).Error
}

// func GetLimitLevel(userName, limitType string) (int, error) {
//
// }
//
// func SetLimitLevel(limitType int, level int, amount int, count int) error {
//
// }
//
// func GetUsersLimit(limitType int, level int) (map[string]int, error) {
//
// }

var typeMapInt = make(map[string]uint)
var typeMapStr = make(map[uint]string)
var initTime = time.Time{}

func GetUserTypeLimitMap() {
	t := time.Now()
	if t.Sub(initTime) < time.Hour*5 {
		return
	}
	initTime = t
	db := middleware.DB
	var limitTypes []custodyModels.LimitType
	err := db.Table("user_limit_type").Find(&limitTypes).Error
	if err != nil {
		return
	}
	for _, limitType := range limitTypes {
		typeMapInt[limitType.Memo] = limitType.ID
		typeMapStr[limitType.ID] = limitType.Memo
	}
}
