package services

import (
	"trade/middleware"
	"trade/models"
	"trade/services/btldb"
	"trade/utils"
)

func CreateOrUpdateBtcBalance(btcBalance *models.BtcBalance) (err error) {
	return CreateBtcBalanceIfNotExistOrUpdate(btcBalance)
}

func CreateBtcBalanceIfNotExistOrUpdate(btcBalance *models.BtcBalance) (err error) {
	var readBtcBalance *models.BtcBalance
	readBtcBalance, err = btldb.ReadBtcBalanceByUsername(btcBalance.Username)
	if err != nil {
		err = btldb.CreateBtcBalance(btcBalance)
		if err != nil {
			return err
		}
		return nil
	}
	readBtcBalance.TotalBalance = btcBalance.TotalBalance
	readBtcBalance.ConfirmedBalance = btcBalance.ConfirmedBalance
	readBtcBalance.UnconfirmedBalance = btcBalance.UnconfirmedBalance
	readBtcBalance.LockedBalance = btcBalance.LockedBalance
	readBtcBalance.DeviceID = btcBalance.DeviceID
	return btldb.UpdateBtcBalance(readBtcBalance)
}

func GetBtcBalanceByUsername(username string) (btcBalance *models.BtcBalance, err error) {
	btcBalance, err = btldb.ReadBtcBalanceByUsername(username)
	if err != nil {
		err = btldb.CreateBtcBalance(&models.BtcBalance{Username: username})
		if err != nil {
			return nil, err
		}
		btcBalance, err = btldb.ReadBtcBalanceByUsername(username)
		if err != nil {
			return nil, err
		}
		return btcBalance, nil
	}
	return btcBalance, nil
}

type BtcBalanceInfo struct {
	ID                 uint   `json:"id"`
	Username           string `json:"username" gorm:"type:varchar(255)"`
	TotalBalance       int    `json:"total_balance"`
	ConfirmedBalance   int    `json:"confirmed_balance"`
	UnconfirmedBalance int    `json:"unconfirmed_balance"`
	LockedBalance      int    `json:"locked_balance"`
}

func GetBtcBalanceCount() (count int64, err error) {

	tx := middleware.DB.Begin()

	err = tx.Table("btc_balances").
		Count(&count).
		Error
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "btc_balances count")
	}

	tx.Rollback()
	return count, nil
}

func GetBtcBalanceOrderLimitOffset(limit int, offset int) (btcBalanceInfos *[]BtcBalanceInfo, err error) {

	tx := middleware.DB.Begin()

	err = tx.Table("btc_balances").
		Select("id, username, total_balance, confirmed_balance, unconfirmed_balance, locked_balance").
		Order("total_balance desc").
		Limit(limit).
		Offset(offset).
		Scan(&btcBalanceInfos).
		Error
	if err != nil {
		return new([]BtcBalanceInfo), utils.AppendErrorInfo(err, "select btc_balances")
	}

	tx.Rollback()
	return btcBalanceInfos, nil
}
