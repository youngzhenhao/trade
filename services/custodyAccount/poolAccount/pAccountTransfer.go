package poolAccount

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"trade/btlLog"
	"trade/models/custodyModels/pAccount"
)

func addBalance(tx *gorm.DB, poolAccountId uint, AssetId string, amount float64, target string, transferDesc string) (uint, error) {
	if !checkAssetId(tx, poolAccountId, AssetId) {
		return 0, fmt.Errorf("AssetId not found in pool account")
	}
	var err error
	balance := pAccount.PAccountBalance{}
	err = tx.FirstOrCreate(&balance, pAccount.PAccountBalance{PoolAccountId: poolAccountId, AssetId: AssetId}).Error
	if err != nil {
		btlLog.CUST.Error("AddBtcBalance error: ", err)
		return 0, err
	}
	balance.Balance += amount
	err = tx.Save(&balance).Error
	if err != nil {
		return 0, err
	}
	bill := pAccount.PAccountBill{
		PoolAccountId: poolAccountId,
		Away:          pAccount.PAccountBillAwayIn,
		AssetId:       AssetId,
		Amount:        amount,
		Target:        target,
		PaymentHash:   transferDesc,
		State:         1,
	}
	err = tx.Create(&bill).Error
	if err != nil {
		return 0, err
	}
	change := pAccount.PAccountBalanceChange{
		PoolAccountId: poolAccountId,
		AssetId:       AssetId,
		BillId:        bill.ID,
		Amount:        amount,
		FinalBalance:  balance.Balance,
	}
	err = tx.Create(&change).Error
	if err != nil {
		return 0, err
	}
	return balance.Id, nil
}

func lessBalance(tx *gorm.DB, poolAccountId uint, AssetId string, amount float64, target string, transferDesc string) (uint, error) {
	if !checkAssetId(tx, poolAccountId, AssetId) {
		return 0, fmt.Errorf("AssetId not found in pool account")
	}
	var err error
	balance := pAccount.PAccountBalance{}
	err = tx.Where("pool_account_id = ? AND asset_id = ?", poolAccountId, AssetId).First(&balance).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		btlLog.CUST.Error("AddBtcBalance error: ", err)
		return 0, err
	}
	if balance.Balance < amount {
		return 0, ErrorNotEnoughBalance
	}
	balance.Balance -= amount
	err = tx.Save(&balance).Error
	if err != nil {
		return 0, err
	}
	bill := pAccount.PAccountBill{
		PoolAccountId: poolAccountId,
		Away:          pAccount.PAccountBillAwayOut,
		AssetId:       AssetId,
		Amount:        amount,
		Target:        target,
		PaymentHash:   transferDesc,
		State:         1,
	}
	err = tx.Create(&bill).Error
	if err != nil {
		return 0, err
	}
	change := pAccount.PAccountBalanceChange{
		PoolAccountId: poolAccountId,
		AssetId:       AssetId,
		BillId:        bill.ID,
		Amount:        amount,
		FinalBalance:  balance.Balance,
	}
	err = tx.Create(&change).Error
	if err != nil {
		return 0, err
	}
	return balance.Id, nil
}

func checkAssetId(tx *gorm.DB, poolAccountId uint, AssetId string) bool {
	err := tx.Where("pool_account_id = ? AND asset_id = ?", poolAccountId, AssetId).First(&pAccount.PAccountAssetId{}).Error
	if err != nil {
		return false
	}
	return true
}
