package custodyLimit

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"time"
	"trade/models/custodyModels"
	caccount "trade/services/custodyAccount/account"
)

const (
	defaultLimitLevel = 1
)

var (
	ErrLimitEntirely = errors.New("ErrLimitEntirely")
)
var limitMux = new(sync.Mutex)

func GetLimit(db *gorm.DB, user *caccount.UserInfo, limitType *custodyModels.LimitType) (*custodyModels.LimitBill, error) {
	err := db.Where(limitType).First(limitType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 无限制限额
			return nil, nil
		}
		// 返回错误
		return nil, err
	}
	limitMux.Lock()
	defer limitMux.Unlock()

	limitBill := custodyModels.LimitBill{
		UserId:    user.User.ID,
		LimitType: limitType.ID,
	}
	err = db.Where("created_at >= CURDATE() AND created_at < CURDATE() + INTERVAL 1 DAY").Where(limitBill).First(&limitBill).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新的限制额度Bill

			// 查询limit表中user的额度等级
			limit := custodyModels.Limit{
				UserId:    user.User.ID,
				LimitType: limitType.ID,
			}
			err = db.Where(limit).First(&limit).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				// 返回错误
				return nil, err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 设置默认限制额度，并保存
				limit.Level = defaultLimitLevel
				if err := db.Create(&limit).Error; err != nil {
					return nil, err // 错误处理
				}
			}

			if limit.Level == 0 {
				// 0无额度
				return nil, ErrLimitEntirely
			}

			// 根据用户的限制额度等级，计算限制额度
			levelLimit := custodyModels.LimitLevel{
				LimitTypeId: limitType.ID,
				Level:       limit.Level,
			}
			err = db.Where(levelLimit).First(&levelLimit).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				// 返回错误
				return nil, err
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 返回不限额额度
				return nil, nil
			}

			// 计算限制额度
			limitBill.TotalAmount = levelLimit.Amount
			limitBill.UseAbleAmount = levelLimit.Amount
			limitBill.TotalCount = levelLimit.Count
			limitBill.UseAbleCount = levelLimit.Count
			limitBill.LocalTime = time.Now().Add(-time.Minute * 15)

			// 创建并返回今日限制额度Bill
			if err := db.Create(&limitBill).Error; err != nil {
				return nil, err // 错误处理
			}

			return &limitBill, nil
		}
		// 返回错误
		return nil, err
	}

	// 返回今日额度Bill
	return &limitBill, nil
}

func CheckLimit(db *gorm.DB, user *caccount.UserInfo, limitType *custodyModels.LimitType, amount float64) error {
	limitBill, err := GetLimit(db, user, limitType)
	if err != nil {
		return err // 错误处理
	}

	if limitBill == nil {
		// 无限制额度
		return nil
	}

	// 检查可用额度是否足够
	if limitBill.UseAbleAmount < amount {
		return fmt.Errorf("%w,剩余额度：%v", errors.New("今日可用交易额度不足"), limitBill.UseAbleAmount)
	}
	if limitBill.UseAbleCount <= 0 {
		return fmt.Errorf("%w,剩余交易次数：%v", errors.New("今日可用交易次数不足"), limitBill.UseAbleCount)
	}
	if time.Now().Sub(limitBill.LocalTime).Minutes() < 10 {
		return errors.New("交易频繁，请稍后再试")
	}
	return nil
}

func MinusLimit(db *gorm.DB, user *caccount.UserInfo, limitType *custodyModels.LimitType, amount float64) error {

	err := db.Where(limitType).First(limitType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 无限制限额
			return nil
		}
		// 返回错误
		return err
	}

	limitMux.Lock()
	defer limitMux.Unlock()

	limitBill := custodyModels.LimitBill{
		UserId:    user.User.ID,
		LimitType: limitType.ID,
	}

	err = db.Where("created_at >= CURDATE() AND created_at < CURDATE() + INTERVAL 1 DAY").Where(limitBill).First(&limitBill).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 未查询到今日限制额度视为无限制额度
			return nil
		}
		// 返回错误
		return err
	}

	// 检查可用额度是否足够
	if limitBill.UseAbleAmount < amount+20 {
		return errors.New("可用额度不足")
	}
	if limitBill.UseAbleCount <= 0 {
		return errors.New("可用交易次数不足")
	}
	if time.Now().Sub(limitBill.LocalTime).Minutes() < 10 {
		return errors.New("交易频繁，请稍后再试")
	}

	limitBill.UseAbleAmount -= amount
	limitBill.UseAbleCount -= 1
	limitBill.LocalTime = time.Now()

	// 保存更新后的 limitBill
	if err := db.Save(&limitBill).Error; err != nil {
		return err // 错误处理
	}

	return nil
}

func AddLimit(limitType int, userId int, limit int) {
	// 该函数可以实现添加限制额度的逻辑
}
