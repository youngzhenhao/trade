package services

import (
	"trade/models"
)

func CreateOrUpdateBtcBalance(btcBalance *models.BtcBalance) (err error) {
	return CreateBtcBalanceIfNotExistOrUpdate(btcBalance)
}

func CreateBtcBalanceIfNotExistOrUpdate(btcBalance *models.BtcBalance) (err error) {
	var readBtcBalance *models.BtcBalance
	readBtcBalance, err = ReadBtcBalanceByUsername(btcBalance.Username)
	if err != nil {
		err = CreateBtcBalance(btcBalance)
		if err != nil {
			return err
		}
		return nil
	}
	readBtcBalance.TotalBalance = btcBalance.TotalBalance
	readBtcBalance.ConfirmedBalance = btcBalance.ConfirmedBalance
	readBtcBalance.UnconfirmedBalance = btcBalance.UnconfirmedBalance
	readBtcBalance.LockedBalance = btcBalance.LockedBalance
	return UpdateBtcBalance(readBtcBalance)
}

func GetBtcBalanceByUsername(username string) (btcBalance *models.BtcBalance, err error) {
	btcBalance, err = ReadBtcBalanceByUsername(username)
	if err != nil {
		err = CreateBtcBalance(&models.BtcBalance{Username: username})
		if err != nil {
			return nil, err
		}
		btcBalance, err = ReadBtcBalanceByUsername(username)
		if err != nil {
			return nil, err
		}
		return btcBalance, nil
	}
	return btcBalance, nil
}
