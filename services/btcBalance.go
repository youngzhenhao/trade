package services

import (
	"trade/models"
	"trade/services/btldb"
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
