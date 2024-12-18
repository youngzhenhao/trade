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
	todayStart := time.Now().Truncate(24 * time.Hour).Add(-8 * time.Hour)

	q = q.Table("user_limit").
		Joins("left join user on user.id = user_limit.user_id").
		Joins("left join (select * from user_limit_bills where created_at >= ?) as bill on bill.user_id = user_limit.user_id"+
			" and bill.limit_type = user_limit.limit_type", todayStart).
		Joins("left join user_limit_type on user_limit_type.id = user_limit.limit_type ")
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
	if amount < 0 || count < 0 {
		return fmt.Errorf("金额或次数不能为负数")
	}

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
	todayStart := time.Now().Truncate(24 * time.Hour).Add(-8 * time.Hour)

	var bill custodyModels.LimitBill
	// 增加一个条件，筛选今日0点开始的记录
	err = db.Where("created_at >= ?", todayStart).Where("user_id =? and limit_type =?", usr.ID, TypeID).First(&bill).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			bill.LocalTime = todayStart
			bill.LimitType = TypeID
			bill.UserId = usr.ID
			bill.TotalAmount = float64(amount)
			bill.UseAbleAmount = float64(amount)
			bill.TotalCount = uint(count)
			bill.UseAbleCount = uint(count)
			return db.Create(&bill).Error
		}
		return err
	}
	bill.UseAbleAmount = float64(amount)
	bill.UseAbleCount = uint(count)
	return db.Save(&bill).Error
}

type LimitTypes struct {
	AssetId      string `json:"assetId"`
	TransferType int    `json:"transferType"`
	LimitName    string `json:"limitName"`
}

func GetLimitTypes(page, pageSize int) (*[]LimitTypes, int64, error) {
	db := middleware.DB

	var total int64
	err := db.Table("user_limit_type").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	var limitTypes []custodyModels.LimitType
	err = db.Table("user_limit_type").Offset((page - 1) * pageSize).Limit(pageSize).Find(&limitTypes).Error
	if err != nil {
		return nil, 0, err
	}
	var limitTypesArr []LimitTypes
	for _, limitType := range limitTypes {
		limitTypesArr = append(limitTypesArr, LimitTypes{
			AssetId:      limitType.AssetId,
			TransferType: int(limitType.TransferType),
			LimitName:    limitType.Memo,
		})
	}
	return &limitTypesArr, total, nil
}

func CreateOrUpdateLimitType(assetId string, transferType int, limitName string) error {
	if limitName == "" {
		return fmt.Errorf("限额名称不能为空")
	}
	db := middleware.DB
	limitType := custodyModels.LimitType{}
	err := db.FirstOrCreate(&limitType, custodyModels.LimitType{
		AssetId:      assetId,
		TransferType: custodyModels.LimitTransferType(transferType),
	}).Error
	if err != nil {
		return err
	}
	if limitType.Memo != limitName {
		limitType.Memo = limitName
		return db.Save(&limitType).Error
	}
	return nil
}

type LimitLevel struct {
	Level  uint    `json:"level"`
	Amount float64 `json:"amount"`
	Count  uint    `json:"count"`
}

func GetLimitTypeLevels(limitName string, page, pageSize int) (*[]LimitLevel, int64, error) {
	db := middleware.DB
	var limitTypes custodyModels.LimitType
	err := db.Table("user_limit_type").Where("memo =?", limitName).First(&limitTypes).Error
	if err != nil {
		return nil, 0, fmt.Errorf("限额类型不存在: %s", limitName)
	}
	var total int64
	err = db.Table("user_limit_type_level").Where("limit_type_id =?", limitTypes.ID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	var limitLevels []custodyModels.LimitLevel
	err = db.Table("user_limit_type_level").Where("limit_type_id =?", limitTypes.ID).Offset((page - 1) * pageSize).Limit(pageSize).Find(&limitLevels).Error
	if err != nil {
		return nil, 0, err
	}
	var limitLevelArr []LimitLevel
	for _, limitLevel := range limitLevels {
		limitLevelArr = append(limitLevelArr, LimitLevel{
			Level:  limitLevel.Level,
			Amount: limitLevel.Amount,
			Count:  limitLevel.Count,
		})
	}
	return &limitLevelArr, total, nil
}

func CreateOrUpdateLimitTypeLevel(limitName string, level int, amount int, count int) error {
	if level <= 0 || amount < 0 || count < 0 {
		return fmt.Errorf("金额或次数不能为负数")
	}
	db := middleware.DB
	var limitTypes custodyModels.LimitType
	err := db.Table("user_limit_type").Where("memo =?", limitName).First(&limitTypes).Error
	if err != nil {
		return fmt.Errorf("限额类型不存在: %s", limitName)
	}

	limitLevel := custodyModels.LimitLevel{}
	err = db.FirstOrCreate(&limitLevel, custodyModels.LimitLevel{
		LimitTypeId: limitTypes.ID,
		Level:       uint(level),
	}).Error
	if err != nil {
		return err
	}
	if limitLevel.Amount != float64(amount) || limitLevel.Count != uint(count) {
		limitLevel.Amount = float64(amount)
		limitLevel.Count = uint(count)
		db.Save(&limitLevel)
	}
	return err

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
