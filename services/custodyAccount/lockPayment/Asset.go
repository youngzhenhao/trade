package lockPayment

import (
	"errors"
	"gorm.io/gorm"
	"trade/middleware"
	"trade/models"
	cModels "trade/models/custodyModels"
	caccount "trade/services/custodyAccount/account"
)

// 冻结Asset
func GetAssetBalance(usr *caccount.UserInfo, assetId string) (err error, unlock float64, locked float64) {
	tx := middleware.DB.Begin()
	lockedBalance := cModels.LockBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.LockAccount.ID, assetId).First(&lockedBalance).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err, 0, 0
		}
		locked = 0
	}
	locked = lockedBalance.Amount

	assetBalance := models.AccountBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.Account.ID, assetId).First(&assetBalance).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err, 0, 0
		}
		unlock = 0
	}
	unlock = assetBalance.Amount

	tx.Commit()
	return
}

// 冻结Asset
func LockAsset(usr *caccount.UserInfo, lockedId string, assetId string, amount float64) error {

	tx := middleware.DB.Begin()
	var err error
	// check balance
	assetBalance := models.AccountBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.Account.ID, assetId).First(&assetBalance).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}
		assetBalance.Amount = 0
	}
	if assetBalance.Amount < amount {
		return NoEnoughBalance
	}

	// lock btc
	lockedBalance := cModels.LockBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.LockAccount.ID, assetId).First(&lockedBalance).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}
		// Init Balance record
		lockedBalance.AssetId = assetId
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
		AssetId:   assetId,
		Amount:    amount,
		LockId:    lockedId,
		BillType:  cModels.LockBillTypeLock,
	}
	if err = tx.Create(&lockBill).Error; err != nil {
		tx.Rollback()
		return err
	}

	// update user account record
	balanceBill := models.Balance{
		AccountId:   usr.Account.ID,
		BillType:    models.BiLLTypeLock,
		Away:        models.AWAY_OUT,
		Amount:      amount,
		Unit:        models.UNIT_SATOSHIS,
		ServerFee:   0,
		AssetId:     &assetId,
		Invoice:     &lockedId,
		PaymentHash: nil,
		State:       models.STATE_SUCCESS,
	}
	if err = tx.Create(&balanceBill).Error; err != nil {
		tx.Rollback()
		return err
	}
	// update user account
	assetBalance.Amount -= amount
	if err = tx.Save(&assetBalance).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func UnlockAsset(usr *caccount.UserInfo, lockedId string, assetId string, amount float64) error {
	tx := middleware.DB.Begin()
	var err error

	// check locked balance
	lockedBalance := cModels.LockBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.LockAccount.ID, assetId).First(&lockedBalance).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}
		lockedBalance.Amount = 0
	}
	if lockedBalance.Amount < amount {
		tx.Rollback()
		return NoEnoughBalance
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
		AssetId:   assetId,
		Amount:    amount,
		LockId:    lockedId,
		BillType:  cModels.LockBillTypeUnlock,
	}
	if err = tx.Create(&unlockBill).Error; err != nil {
		tx.Rollback()
		return err
	}

	// update user account record
	balanceBill := models.Balance{
		AccountId:   usr.Account.ID,
		BillType:    models.BiLLTypeLock,
		Away:        models.AWAY_IN,
		Amount:      amount,
		Unit:        models.UNIT_SATOSHIS,
		ServerFee:   0,
		AssetId:     &assetId,
		Invoice:     &lockedId,
		PaymentHash: nil,
		State:       models.STATE_SUCCESS,
	}
	if err = tx.Create(&balanceBill).Error; err != nil {
		tx.Rollback()
		return err
	}
	// update user account
	assetBalance := models.AccountBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.Account.ID, assetId).First(&assetBalance).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}
		assetBalance.AssetId = assetId
		assetBalance.Amount = 0
		assetBalance.AccountID = usr.Account.ID
	}
	assetBalance.Amount += amount
	if err = tx.Save(&assetBalance).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func transferLockedAsset(usr *caccount.UserInfo, lockedId string, assetId string, amount float64, toUser *caccount.UserInfo) error {
	tx := middleware.DB.Begin()
	defer tx.Commit()

	var err error
	// check locked balance
	lockedBalance := cModels.LockBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.LockAccount.ID, assetId).First(&lockedBalance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}
		lockedBalance.Amount = 0
	}
	if lockedBalance.Amount < amount {
		tx.Rollback()
		return NoEnoughBalance
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
		AssetId:   assetId,
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
		AssetId:    assetId,
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
		AssetId:     &assetId,
		Invoice:     &lockedId,
		PaymentHash: nil,
		State:       models.STATE_SUCCESS,
	}
	if err = tx.Create(&balanceBill).Error; err != nil {
		tx.Rollback()
		return err
	}

	// update user account
	assetBalance := models.AccountBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", toUser.Account.ID, assetId).First(&assetBalance).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}
		assetBalance.AssetId = assetId
		assetBalance.Amount = 0
		assetBalance.AccountID = toUser.Account.ID
	}
	assetBalance.Amount += amount
	if err = tx.Save(&assetBalance).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func transferAsset(usr *caccount.UserInfo, lockedId string, assetId string, amount float64, toUser *caccount.UserInfo) error {
	tx := middleware.DB.Begin()

	var err error
	// check balance
	assetBalance := models.AccountBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.Account.ID, assetId).First(&assetBalance).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}
		assetBalance.Amount = 0
	}
	if assetBalance.Amount < amount {
		return NoEnoughBalance
	}

	// Create transferBTC Bill
	transferBill := cModels.LockBill{
		AccountID: usr.LockAccount.ID,
		LockId:    lockedId,
		AssetId:   assetId,
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
		AssetId:    assetId,
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
		Unit:        models.UNIT_ASSET_NORMAL,
		ServerFee:   0,
		AssetId:     &assetId,
		Invoice:     &lockedId,
		PaymentHash: nil,
		State:       models.STATE_SUCCESS,
	}
	if err = tx.Create(&balanceBill).Error; err != nil {
		tx.Rollback()
		return err
	}
	// update user account
	assetBalance.Amount -= amount
	if err = tx.Save(&assetBalance).Error; err != nil {
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
		Unit:        models.UNIT_ASSET_NORMAL,
		ServerFee:   0,
		AssetId:     &assetId,
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
	assetBalanceRev := models.AccountBalance{}
	if err = txRev.Where("account_id =? AND asset_id =?", toUser.Account.ID, assetId).First(&assetBalanceRev).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			txRev.Rollback()
			return err
		}
		assetBalanceRev.AssetId = assetId
		assetBalanceRev.Amount = 0
		assetBalanceRev.AccountID = toUser.Account.ID
	}
	assetBalanceRev.Amount += amount
	if err = txRev.Save(&assetBalanceRev).Error; err != nil {
		txRev.Rollback()
		return err
	}
	txRev.Commit()
	return nil
}
