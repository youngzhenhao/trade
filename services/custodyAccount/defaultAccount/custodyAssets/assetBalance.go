package custodyAssets

import (
	"errors"
	"gorm.io/gorm"
	"sync"
	"trade/btlLog"
	"trade/models/custodyModels"
	"trade/services/custodyAccount/account"
)

var assetBalanceMutex = &sync.RWMutex{}

var NotEnoughAssetBalance = errors.New("not enough asset balance")

func GetAssetsBalances(Db *gorm.DB, id uint) *[]custodyModels.AccountBalance {
	assetBalanceMutex.RLock()
	defer assetBalanceMutex.RUnlock()

	var accountBalances []custodyModels.AccountBalance
	Db.Where("account_Id =?", id).Find(&accountBalances)
	return &accountBalances
}

func GetAssetBalance(Db *gorm.DB, id uint, assetId string) float64 {
	assetBalanceMutex.RLock()
	defer assetBalanceMutex.RUnlock()

	var accountBalance custodyModels.AccountBalance
	err := Db.Where("account_Id =? and asset_Id =?", id, assetId).First(&accountBalance).Error
	if err != nil {
		return 0
	}

	return accountBalance.Amount
}

func CheckAssetBalance(Db *gorm.DB, usr *account.UserInfo, assetId string, amount float64) bool {
	assetBalanceMutex.RLock()
	defer assetBalanceMutex.RUnlock()
	balance := GetAssetBalance(Db, usr.Account.ID, assetId)
	if balance < amount {
		return false
	}
	return true
}

func AddAssetBalance(Db *gorm.DB, usr *account.UserInfo, amount float64, balanceId uint, assetId string, ChangeType custodyModels.ChangeType) (uint, error) {
	assetBalanceMutex.Lock()
	defer assetBalanceMutex.Unlock()

	balance := custodyModels.AccountBalance{}
	err := Db.FirstOrCreate(&balance, custodyModels.AccountBalance{AccountID: usr.Account.ID, AssetId: assetId}).Error
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
		AssetId:      assetId,
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
func LessAssetBalance(Db *gorm.DB, usr *account.UserInfo, amount float64, balanceId uint, assetId string, ChangeType custodyModels.ChangeType) (uint, error) {
	assetBalanceMutex.Lock()
	defer assetBalanceMutex.Unlock()

	balance := custodyModels.AccountBalance{}
	err := Db.FirstOrCreate(&balance, custodyModels.AccountBalance{AccountID: usr.Account.ID, AssetId: assetId}).Error
	if err != nil {
		btlLog.CUST.Error("AddBtcBalance error: ", err)
		return 0, err
	}
	if balance.Amount < amount {
		btlLog.CUST.Error("LessBtcBalance error: ", "not enough asset balance")
		return 0, NotEnoughAssetBalance
	}
	balance.Amount -= amount
	err = Db.Save(&balance).Error
	if err != nil {
		btlLog.CUST.Error("AddBtcBalance error: ", err)
		return 0, err
	}
	change := custodyModels.AccountBalanceChange{
		AccountId:    usr.Account.ID,
		AssetId:      assetId,
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
