package custodyBtc

import (
	"errors"
	"gorm.io/gorm"
	"sync"
	"trade/btlLog"
	"trade/models/custodyModels"
	"trade/services/custodyAccount/account"
)

var btcMutex = &sync.RWMutex{}

const BtcId = "00"

var NotEnoughBalance = errors.New("not enough balance")

func getBtcBalance(Db *gorm.DB, id uint) float64 {
	btcMutex.RLock()
	defer btcMutex.RUnlock()
	balance := custodyModels.AccountBtcBalance{}
	err := Db.Where("account_id =?", id).First(&balance).Error
	if err != nil {
		return 0
	}
	return balance.Amount
}

func CheckBtcBalance(Db *gorm.DB, usr *account.UserInfo, amount float64) bool {
	btcMutex.RLock()
	defer btcMutex.RUnlock()
	balance := getBtcBalance(Db, usr.Account.ID)
	if balance < amount {
		return false
	}
	return true
}

func AddBtcBalance(Db *gorm.DB, usr *account.UserInfo, amount float64, balanceId uint, ChangeType custodyModels.ChangeType) (uint, error) {
	btcMutex.Lock()
	defer btcMutex.Unlock()

	balance := custodyModels.AccountBtcBalance{}
	err := Db.FirstOrCreate(&balance, custodyModels.AccountBtcBalance{AccountId: usr.Account.ID}).Error
	if err != nil {
		btlLog.CUST.Error("AddBtcBalance error: ", err)
		return 0, err
	}
	balance.Amount += amount
	err = Db.Save(&balance).Error
	if err != nil {
		btlLog.CUST.Error("AddBtcBalance error: ", err)
		return 0, err
	}
	change := custodyModels.AccountBalanceChange{
		AccountId:    usr.Account.ID,
		AssetId:      BtcId,
		ChangeAmount: amount,
		Away:         custodyModels.ChangeAwayAdd,
		FinalBalance: balance.Amount,
		BalanceId:    balanceId,
		ChangeType:   ChangeType,
	}
	err = Db.Create(&change).Error
	if err != nil {
		btlLog.CUST.Error("AddBtcBalance error: ", err)
		return 0, err
	}
	return change.ID, nil
}
func LessBtcBalance(Db *gorm.DB, usr *account.UserInfo, amount float64, balanceId uint, ChangeType custodyModels.ChangeType) (uint, error) {
	btcMutex.Lock()
	defer btcMutex.Unlock()

	balance := custodyModels.AccountBtcBalance{}
	err := Db.FirstOrCreate(&balance, custodyModels.AccountBtcBalance{AccountId: usr.Account.ID}).Error
	if err != nil {
		btlLog.CUST.Error("AddBtcBalance error: ", err)
		return 0, err
	}
	if balance.Amount < amount {
		btlLog.CUST.Error("LessBtcBalance error: ", "not enough balance")
		return 0, NotEnoughBalance
	}
	balance.Amount -= amount
	err = Db.Save(&balance).Error
	if err != nil {
		btlLog.CUST.Error("AddBtcBalance error: ", err)
		return 0, err
	}
	change := custodyModels.AccountBalanceChange{
		AccountId:    usr.Account.ID,
		AssetId:      BtcId,
		ChangeAmount: amount,
		Away:         custodyModels.ChangeAwayLess,
		FinalBalance: balance.Amount,
		BalanceId:    balanceId,
		ChangeType:   ChangeType,
	}
	err = Db.Create(&change).Error
	if err != nil {
		btlLog.CUST.Error("AddBtcBalance error: ", err)
		return 0, err
	}
	return change.ID, nil
}
