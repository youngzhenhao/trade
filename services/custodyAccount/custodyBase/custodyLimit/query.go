package custodyLimit

import (
	"time"
	"trade/middleware"
	"trade/models/custodyModels"
)

type UserLimit struct {
	UserName          string  `json:"user_name" gorm:"column:user_name"`
	LimitType         string  `json:"limit_type" gorm:"column:memo"`
	Level             int     `json:"level" gorm:"column:level"`
	TodayAmount       float64 `json:"today_amount" gorm:"column:total_amount"`
	TodayCount        int     `json:"today_count" gorm:"column:total_count"`
	TodayUsefulAmount float64 `json:"today_useful_amount" gorm:"column:use_able_amount"`
	TodayUsefulCount  int     `json:"today_useful_count" gorm:"column:use_able_count"`
}

func GetUserLimit(userName, LimitType string) (*[]UserLimit, error) {
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

	err := q.Debug().Table("user_limit").
		Joins("left join user on user.id = user_limit.user_id").
		Joins("left join (select * from user_limit_bills where created_at >= ?) as bill on bill.user_id = user_limit.user_id", todayStart).
		Joins("left join user_limit_type on user_limit_type.id = user_limit.limit_type").
		Select("user.user_name, " +
			"user_limit_type.memo," +
			"user_limit.level," +
			"bill.total_amount,bill.use_able_amount,bill.total_count,bill.use_able_count").
		Scan(&userLimits).Error
	if err != nil {
		return nil, err
	}
	return &userLimits, nil
}

// func SetUserLimitLevel(userName, limitType string, level int) error {
//
// }
//
// func SetUserTodayLimit(userName, limitType string, amount int, count int) error {
//
// }
//
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
