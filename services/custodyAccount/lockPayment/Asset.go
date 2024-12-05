package lockPayment

import (
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"trade/middleware"
	"trade/models"
	cModels "trade/models/custodyModels"
	caccount "trade/services/custodyAccount/account"
	"trade/services/custodyAccount/defaultAccount/custodyAssets"
)

// GetAssetBalance 获取用户资产余额
func GetAssetBalance(usr *caccount.UserInfo, assetId string) (err error, unlock float64, locked float64) {
	tx := middleware.DB.Begin()
	defer tx.Rollback()
	lockedBalance := cModels.LockBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.LockAccount.ID, assetId).First(&lockedBalance).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return ServiceError, 0, 0
		}
		locked = 0
		err = nil
	}
	locked = lockedBalance.Amount

	assetBalance := cModels.AccountBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.Account.ID, assetId).First(&assetBalance).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return ServiceError, 0, 0
		}
		unlock = 0
		err = nil
	}
	unlock = assetBalance.Amount

	tx.Commit()
	return
}

// LockAsset 冻结Asset
func LockAsset(usr *caccount.UserInfo, lockedId string, assetId string, amount float64) error {

	tx := middleware.DB.Begin()
	defer tx.Rollback()
	var err error

	// lock btc
	lockedBalance := cModels.LockBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.LockAccount.ID, assetId).First(&lockedBalance).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return ServiceError
		}
		// Init Balance record
		lockedBalance.AssetId = assetId
		lockedBalance.AccountID = usr.LockAccount.ID
		lockedBalance.Amount = 0
	}
	lockedBalance.Amount += amount
	if err = tx.Save(&lockedBalance).Error; err != nil {
		return ServiceError
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
		var mySQLErr *mysql.MySQLError
		if errors.As(err, &mySQLErr) {
			if mySQLErr.Number == 1062 {
				return RepeatedLockId
			}
		}
		return ServiceError
	}
	Invoice := InvoiceLocked
	// update user account record
	balanceBill := models.Balance{
		AccountId:   usr.Account.ID,
		BillType:    models.BiLLTypeLock,
		Away:        models.AWAY_OUT,
		Amount:      amount,
		Unit:        models.UNIT_ASSET_NORMAL,
		ServerFee:   0,
		AssetId:     &assetId,
		Invoice:     &Invoice,
		PaymentHash: &lockedId,
		State:       models.STATE_SUCCESS,
		TypeExt: &models.BalanceTypeExt{
			Type: models.BTExtLocked,
		},
	}
	if err = tx.Create(&balanceBill).Error; err != nil {
		return ServiceError
	}
	_, err = custodyAssets.LessAssetBalance(tx, usr, balanceBill.Amount, balanceBill.ID, *balanceBill.AssetId, cModels.ChangeTypeLock)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func UnlockAsset(usr *caccount.UserInfo, lockedId string, assetId string, amount float64, version int) error {
	tx := middleware.DB.Begin()
	defer tx.Rollback()
	var err error

	// check locked balance
	lockedBalance := cModels.LockBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.LockAccount.ID, assetId).First(&lockedBalance).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return ServiceError
		}
		lockedBalance.Amount = 0
	}
	// version 1.0 check awardAmount
	if version == 0 {
		if lockedBalance.Amount < amount {
			return NoEnoughBalance
		}
		if (lockedBalance.Amount - lockedBalance.AwardAmount) < amount {
			return fmt.Errorf("%w,have  %f is awardAmount", NoEnoughBalance, lockedBalance.AwardAmount)
		}
		// update locked balance
		lockedBalance.Amount -= amount
	} else if version == 1 {
		if lockedBalance.AwardAmount < amount {
			return fmt.Errorf("%w,have  %f is awardAmount", NoEnoughBalance, lockedBalance.AwardAmount)
		}
		// update locked balance
		lockedBalance.Amount -= amount
		lockedBalance.AwardAmount -= amount
	}

	if err = tx.Save(&lockedBalance).Error; err != nil {
		return ServiceError
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
		var mySQLErr *mysql.MySQLError
		if errors.As(err, &mySQLErr) {
			if mySQLErr.Number == 1062 {
				return RepeatedLockId
			}
		}
		return ServiceError
	}

	Invoice := InvoiceUnlocked
	// update user account record
	balanceBill := models.Balance{
		AccountId:   usr.Account.ID,
		BillType:    models.BiLLTypeLock,
		Away:        models.AWAY_IN,
		Amount:      amount,
		Unit:        models.UNIT_ASSET_NORMAL,
		ServerFee:   0,
		AssetId:     &assetId,
		Invoice:     &Invoice,
		PaymentHash: &lockedId,
		State:       models.STATE_SUCCESS,
		TypeExt: &models.BalanceTypeExt{
			Type: models.BTExtLocked,
		},
	}
	if err = tx.Create(&balanceBill).Error; err != nil {
		return ServiceError
	}
	// update user account
	_, err = custodyAssets.AddAssetBalance(tx, usr, balanceBill.Amount, balanceBill.ID, *balanceBill.AssetId, cModels.ChangeTypeUnlock)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func transferLockedAsset(usr *caccount.UserInfo, lockedId string, assetId string, amount float64, toUser *caccount.UserInfo) error {
	tx := middleware.DB.Begin()
	defer tx.Rollback()

	var err error
	// check locked balance
	lockedBalance := cModels.LockBalance{}
	if err = tx.Where("account_id =? AND asset_id =?", usr.LockAccount.ID, assetId).First(&lockedBalance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ServiceError
		}
		lockedBalance.Amount = 0
	}
	if lockedBalance.Amount < amount {
		return NoEnoughBalance
	}

	// update locked balance
	lockedBalance.Amount -= amount
	if err = tx.Save(&lockedBalance).Error; err != nil {
		return ServiceError
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
		var mySQLErr *mysql.MySQLError
		if errors.As(err, &mySQLErr) {
			if mySQLErr.Number == 1062 {
				return RepeatedLockId
			}
		}
		return ServiceError
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
		return ServiceError
	}

	invoice := InvoicePendingOderReceive
	if usr.User.Username == FeeNpubkey {
		invoice = InvoicePendingOderAward
	}

	// update user account record
	balanceBill := models.Balance{
		AccountId:   toUser.Account.ID,
		BillType:    models.BillTypePendingOder,
		Away:        models.AWAY_IN,
		Amount:      amount,
		Unit:        models.UNIT_ASSET_NORMAL,
		ServerFee:   0,
		AssetId:     &assetId,
		Invoice:     &invoice,
		PaymentHash: &lockedId,
		State:       models.STATE_SUCCESS,
		TypeExt: &models.BalanceTypeExt{
			Type: models.BTExtLockedTransfer,
		},
	}
	if err = tx.Create(&balanceBill).Error; err != nil {
		return ServiceError
	}

	// update user account
	_, err = custodyAssets.AddAssetBalance(tx, toUser, balanceBill.Amount, balanceBill.ID, *balanceBill.AssetId, cModels.ChangeTypeLockedTransfer)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func transferAsset(usr *caccount.UserInfo, lockedId string, assetId string, amount float64, toUser *caccount.UserInfo) error {
	tx := middleware.DB.Begin()
	defer tx.Rollback()

	var err error
	// Create transferBTC Bill
	transferBill := cModels.LockBill{
		AccountID: usr.LockAccount.ID,
		LockId:    lockedId,
		AssetId:   assetId,
		Amount:    amount,
		BillType:  cModels.LockBillTypeTransferByUnlockAsset,
	}
	if err = tx.Create(&transferBill).Error; err != nil {
		var mySQLErr *mysql.MySQLError
		if errors.As(err, &mySQLErr) {
			if mySQLErr.Number == 1062 {
				return RepeatedLockId
			}
		}
		return ServiceError
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
		var mySQLErr *mysql.MySQLError
		if errors.As(err, &mySQLErr) {
			if mySQLErr.Number == 1062 {
				return RepeatedLockId
			}
		}
		return ServiceError
	}

	payInvoice := InvoicePendingOderPay
	if usr.User.Username == FeeNpubkey {
		payInvoice = InvoicePendingOderAward
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
		Invoice:     &payInvoice,
		PaymentHash: &lockedId,
		State:       models.STATE_SUCCESS,
		TypeExt: &models.BalanceTypeExt{
			Type: models.BTExtLockedTransfer,
		},
	}
	if err = tx.Create(&balanceBill).Error; err != nil {
		return ServiceError
	}
	// update user account
	_, err = custodyAssets.LessAssetBalance(tx, usr, balanceBill.Amount, balanceBill.ID, *balanceBill.AssetId, cModels.ChangeTypeLockedTransfer)
	if err != nil {
		return err
	}
	tx.Commit()

	//Second tx
	txRev := middleware.DB.Begin()
	defer txRev.Rollback()

	recInvoice := InvoicePendingOderReceive
	if usr.User.Username == FeeNpubkey {
		recInvoice = InvoicePendingOderAward
	}

	// update user account record
	balanceBillRev := models.Balance{
		AccountId:   toUser.Account.ID,
		BillType:    models.BillTypePendingOder,
		Away:        models.AWAY_IN,
		Amount:      amount,
		Unit:        models.UNIT_ASSET_NORMAL,
		ServerFee:   0,
		AssetId:     &assetId,
		Invoice:     &recInvoice,
		PaymentHash: &lockedId,
		State:       models.STATE_SUCCESS,
		TypeExt: &models.BalanceTypeExt{
			Type: models.BTExtLockedTransfer,
		},
	}
	if err = txRev.Create(&balanceBillRev).Error; err != nil {
		return ServiceError
	}
	// update billExt record
	BillExt.Status = cModels.LockBillExtStatusSuccess
	if err = txRev.Save(&BillExt).Error; err != nil {
		return ServiceError
	}

	// update user account
	_, err = custodyAssets.AddAssetBalance(txRev, toUser, balanceBillRev.Amount, balanceBillRev.ID, *balanceBillRev.AssetId, cModels.ChangeTypeLockedTransfer)
	if err != nil {
		return err
	}
	txRev.Commit()
	return nil
}
