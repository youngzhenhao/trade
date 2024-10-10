package lockPayment

import (
	"errors"
	"gorm.io/gorm"
	"trade/middleware"
	"trade/models"
	cModels "trade/models/custodyModels"
	caccount "trade/services/custodyAccount/account"

	rpc "trade/services/servicesrpc"
)

func GetBtcBalance(usr *caccount.UserInfo) (err error, unlock float64, locked float64) {
	tx := middleware.DB.Begin()
	lockedBalance := cModels.LockBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.LockAccount.ID, btcId).First(&lockedBalance).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err, 0, 0
		}
		locked = 0
	}
	locked = lockedBalance.Amount

	acc, err := rpc.AccountInfo(usr.Account.UserAccountCode)
	if err != nil {
		return err, 0, 0
	}
	unlock = float64(acc.CurrentBalance)
	tx.Commit()
	return
}

// LockBTC 冻结BTC
func LockBTC(usr *caccount.UserInfo, lockedId string, amount float64) error {
	tx := middleware.DB.Begin()
	var err error
	// check balance
	acc, err := rpc.AccountInfo(usr.Account.UserAccountCode)
	if err != nil {
		return err
	}
	if float64(acc.CurrentBalance) < amount {
		return nil
	}
	// lock btc
	lockedBalance := cModels.LockBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.LockAccount.ID, btcId).First(&lockedBalance).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}
		// Init Balance record
		lockedBalance.AssetId = btcId
		lockedBalance.AccountID = usr.LockAccount.ID
		lockedBalance.Amount = 0
	}
	lockedBalance.Amount += amount
	if err = tx.Save(&lockedBalance).Error; err != nil {
		tx.Rollback()
		return err
	}

	// lockBill record
	lockBill := cModels.LockBill{
		AccountID: usr.LockAccount.ID,
		AssetId:   btcId,
		Amount:    amount,
		LockId:    lockedId,
		BillType:  cModels.LockBillTypeLock,
	}
	if err = tx.Create(&lockBill).Error; err != nil {
		tx.Rollback()
		return err
	}

	BtcId := btcId
	// update user account record
	balanceBill := models.Balance{
		AccountId:   usr.Account.ID,
		BillType:    models.BiLLTypeLock,
		Away:        models.AWAY_OUT,
		Amount:      amount,
		Unit:        models.UNIT_SATOSHIS,
		ServerFee:   0,
		AssetId:     &BtcId,
		Invoice:     &lockedId,
		PaymentHash: nil,
		State:       models.STATE_SUCCESS,
	}
	if err = tx.Create(&balanceBill).Error; err != nil {
		tx.Rollback()
		return err
	}
	// update user account
	resultAmount := float64(acc.CurrentBalance) - amount
	_, err = rpc.AccountUpdate(usr.Account.UserAccountCode, int64(resultAmount), -1)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// UnlockBTC 解冻BTC
func UnlockBTC(usr *caccount.UserInfo, lockedId string, amount float64) error {
	tx := middleware.DB.Begin()
	var err error

	// check locked balance
	lockedBalance := cModels.LockBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.LockAccount.ID, btcId).First(&lockedBalance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}
		tx.Rollback()
		return err
	}
	if lockedBalance.Amount < amount {
		tx.Rollback()
		return err
	}

	// update locked balance
	lockedBalance.Amount -= amount
	if err = tx.Save(&lockedBalance).Error; err != nil {
		tx.Rollback()
		return err
	}

	// unlockBill record
	unlockBill := cModels.LockBill{
		AccountID: usr.LockAccount.ID,
		AssetId:   btcId,
		Amount:    amount,
		LockId:    lockedId,
		BillType:  cModels.LockBillTypeUnlock,
	}
	if err = tx.Create(&unlockBill).Error; err != nil {
		tx.Rollback()
		return err
	}

	BtcId := btcId

	// update user account record
	balanceBill := models.Balance{
		AccountId:   usr.Account.ID,
		BillType:    models.BiLLTypeLock,
		Away:        models.AWAY_IN,
		Amount:      amount,
		Unit:        models.UNIT_SATOSHIS,
		ServerFee:   0,
		AssetId:     &BtcId,
		Invoice:     &lockedId,
		PaymentHash: nil,
		State:       models.STATE_SUCCESS,
	}
	if err = tx.Create(&balanceBill).Error; err != nil {
		tx.Rollback()
		return err
	}

	// update user account
	acc, err := rpc.AccountInfo(usr.Account.UserAccountCode)
	if err != nil {
		tx.Rollback()
		return err
	}
	resultAmount := float64(acc.CurrentBalance) + amount
	_, err = rpc.AccountUpdate(usr.Account.UserAccountCode, int64(resultAmount), -1)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// transferLockedBTC 转账冻结的BTC
func transferLockedBTC(usr *caccount.UserInfo, lockedId string, amount float64, toUser *caccount.UserInfo) error {
	tx := middleware.DB.Begin()
	defer tx.Commit()
	BtcId := btcId

	var err error

	// check locked balance
	lockedBalance := cModels.LockBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.LockAccount.ID, btcId).First(&lockedBalance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}
		tx.Rollback()
		return err
	}
	if lockedBalance.Amount < amount {
		tx.Rollback()
		return err
	}

	// update locked balance
	lockedBalance.Amount -= amount
	if err = tx.Save(&lockedBalance).Error; err != nil {
		tx.Rollback()
		return err
	}

	// unlockBill record
	transferBill := cModels.LockBill{
		AccountID: usr.LockAccount.ID,
		AssetId:   btcId,
		Amount:    amount,
		LockId:    lockedId,
		BillType:  cModels.LockBillTypeTransferByLockAsset,
	}
	if err = tx.Create(&transferBill).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create transferBTC BillExt
	BillExt := cModels.LockBillExt{
		BillId:     transferBill.ID,
		LockId:     lockedId,
		PayAccType: cModels.LockBillExtPayAccTypeLock,
		PayAccId:   usr.LockAccount.ID,
		RevAccId:   toUser.Account.ID,
		Amount:     amount,
		AssetId:    btcId,
		Status:     cModels.LockBillExtStatusSuccess,
	}
	if err = tx.Create(&BillExt).Error; err != nil {
		tx.Rollback()
		return err
	}

	// update user account record
	balanceBill := models.Balance{
		AccountId:   toUser.Account.ID,
		BillType:    models.BillTypePendingOder,
		Away:        models.AWAY_IN,
		Amount:      amount,
		Unit:        models.UNIT_SATOSHIS,
		ServerFee:   0,
		AssetId:     &BtcId,
		Invoice:     &lockedId,
		PaymentHash: nil,
		State:       models.STATE_SUCCESS,
	}
	if err = tx.Create(&balanceBill).Error; err != nil {
		tx.Rollback()
		return err
	}

	// update user account
	acc, err := rpc.AccountInfo(toUser.Account.UserAccountCode)
	if err != nil {
		tx.Rollback()
		return err
	}
	resultAmount := float64(acc.CurrentBalance) + amount
	_, err = rpc.AccountUpdate(toUser.Account.UserAccountCode, int64(resultAmount), -1)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// transferBTC 转账非冻结的BTC
func transferBTC(usr *caccount.UserInfo, lockedId string, amount float64, toUser *caccount.UserInfo) error {
	BtcId := btcId
	tx := middleware.DB.Begin()

	var err error

	// check balance
	acc, err := rpc.AccountInfo(usr.Account.UserAccountCode)
	if err != nil {
		return err
	}
	if float64(acc.CurrentBalance) < amount {
		return nil
	}

	// Create transferBTC Bill
	transferBill := cModels.LockBill{
		AccountID: usr.LockAccount.ID,
		LockId:    lockedId,
		AssetId:   btcId,
		Amount:    amount,
		BillType:  cModels.LockBillTypeTransferByUnlockAsset,
	}
	if err = tx.Create(&transferBill).Error; err != nil {
		tx.Rollback()
		return err
	}
	// Create transferBTC BillExt
	BillExt := cModels.LockBillExt{
		BillId:     transferBill.ID,
		LockId:     lockedId,
		PayAccType: cModels.LockBillExtPayAccTypeUnlock,
		PayAccId:   usr.Account.ID,
		RevAccId:   toUser.Account.ID,
		Amount:     amount,
		AssetId:    btcId,
		Status:     cModels.LockBillExtStatusInit,
	}
	if err = tx.Create(&BillExt).Error; err != nil {
		tx.Rollback()
		return err
	}

	// transfer balance record
	balanceBill := models.Balance{
		AccountId:   usr.Account.ID,
		BillType:    models.BillTypePendingOder,
		Away:        models.AWAY_OUT,
		Amount:      amount,
		Unit:        models.UNIT_SATOSHIS,
		ServerFee:   0,
		AssetId:     &BtcId,
		Invoice:     &lockedId,
		PaymentHash: nil,
		State:       models.STATE_SUCCESS,
	}
	if err = tx.Create(&balanceBill).Error; err != nil {
		tx.Rollback()
		return err
	}
	// update user account
	resultAmount := float64(acc.CurrentBalance) - amount
	_, err = rpc.AccountUpdate(usr.Account.UserAccountCode, int64(resultAmount), -1)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	//Second tx
	txRev := middleware.DB.Begin()
	// update user account record
	balanceBillRev := models.Balance{
		AccountId:   toUser.Account.ID,
		BillType:    models.BillTypePendingOder,
		Away:        models.AWAY_IN,
		Amount:      amount,
		Unit:        models.UNIT_SATOSHIS,
		ServerFee:   0,
		AssetId:     &BtcId,
		Invoice:     &lockedId,
		PaymentHash: nil,
		State:       models.STATE_SUCCESS,
	}
	if err = txRev.Create(&balanceBillRev).Error; err != nil {
		txRev.Rollback()
		return err
	}
	// update billExt record
	BillExt.Status = cModels.LockBillExtStatusSuccess
	if err = txRev.Save(&BillExt).Error; err != nil {
		txRev.Rollback()
		return err
	}

	// update user account
	accRev, err := rpc.AccountInfo(toUser.Account.UserAccountCode)
	if err != nil {
		txRev.Rollback()
		return err
	}
	resultAmountRev := float64(accRev.CurrentBalance) + amount
	_, err = rpc.AccountUpdate(toUser.Account.UserAccountCode, int64(resultAmountRev), -1)
	if err != nil {
		txRev.Rollback()
		return err
	}
	txRev.Commit()
	return nil
}
