package lockPayment

import (
	"gorm.io/gorm"
	"strconv"
	"trade/btlLog"
	"trade/middleware"
)

type DealUser struct {
	ID      int    `gorm:"column:id;primary_key"`
	NpubKey string `gorm:"column:npubkey"`
	Type    string `gorm:"column:type"`
	IsDeal  int    `gorm:"column:is_deal"`
}

func LoadDealUser(db *gorm.DB) (*[]DealUser, error) {
	var dealUsers []DealUser
	err := db.Table("z_user").Where("is_deal =?", 0).Find(&dealUsers).Error
	if err != nil {
		return nil, err
	}
	return &dealUsers, nil
}
func RunDealUser() error {
	db := middleware.DB
	usr, err := LoadDealUser(db)
	if err != nil {
		return err
	}
	for _, user := range *usr {
		lockedId := "/Unlock/alllocked/cancelOder/N2/" + strconv.Itoa(user.ID)
		if user.Type == "btc" {
			err, _, f2, f3 := GetBalance(user.NpubKey, btcId)
			if err != nil {
				btlLog.CUST.Error("GetBalance error,%s,%s", err, user.NpubKey)
				continue
			}
			amount := f2 - f3
			if amount <= 0 {
				continue
			}
			err = Unlock(user.NpubKey, lockedId, btcId, amount, 0)
			if err != nil {
				btlLog.CUST.Error("Unlock error,%s,%s", err, user.NpubKey)
				continue
			}
		} else if user.Type == "asset" {
			assetId := "47ed120d4b173eb79ba46cd1959bb9c881cb69332cf8a21336110bda05402308"
			err, _, f2, f3 := GetBalance(user.NpubKey, assetId)
			if err != nil {
				btlLog.CUST.Error("GetBalance error,%s,%s", err, user.NpubKey)
				continue
			}
			amount := f2 - f3
			if amount <= 0 {
				continue
			}
			err = Unlock(user.NpubKey, lockedId, assetId, amount, 0)
			if err != nil {
				btlLog.CUST.Error("Unlock error,%s,%s", err, user.NpubKey)
				continue
			}

		}
		err = db.Table("z_user").Where("id = ?", user.ID).Update("is_deal", 1).Error
		if err != nil {
			btlLog.CUST.Error("Save error,%s,%s", err, user.NpubKey)
			continue
		}
	}
	return nil
}
