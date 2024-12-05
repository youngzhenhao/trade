package lockPayment

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"trade/btlLog"
	"trade/middleware"
	cModels "trade/models/custodyModels"
	caccount "trade/services/custodyAccount/account"
	"trade/services/custodyAccount/custodyBase/custodyMutex"
)

const (
	FeeNpubkey = "bitlong"
	btcId      = "00"
)

var (
	NoEnoughBalance    = errors.New("NoEnoughBalance")
	GetAccountError    = errors.New("GetAccountError")
	RevNpubKeyNotFound = errors.New("RevNpubKeyNotFound")
	AssetIdNotFound    = errors.New("AssetIdNotFound")
	ServiceError       = errors.New("ServiceError")
	BadRequest         = errors.New("BadRequest")
	RepeatedLockId     = errors.New("RepeatedLockId")
)

const (
	InvoiceLocked         = "locked"
	InvoiceUnlocked       = "unlocked"
	InvoicePendingOderPay = "pendingOderPay"
	//InvoicePendingOderPayByLock   = "pendingOderPayByLock"
	InvoicePendingOderReceive = "pendingOderPayReceive"
	InvoicePendingOderAward   = "PENDING_ORDER_AWARD"
)

const LockMutexKey = "lockPayment"

func GetErrorCode(err error) int {
	switch {
	case errors.Is(err, NoEnoughBalance):
		return 10001
	case errors.Is(err, GetAccountError):
		return 10002
	case errors.Is(err, RevNpubKeyNotFound):
		return 10003
	case errors.Is(err, AssetIdNotFound):
		return 10004
	case errors.Is(err, ServiceError):
		return 10005
	case errors.Is(err, BadRequest):
		return 10006
	case errors.Is(err, RepeatedLockId):
		return 10007
	default:
		return 10005
	}
}

func GetBalance(npubkey, assetId string) (err error, unlockedBalance float64, lockedBalance float64) {
	if npubkey == FeeNpubkey {
		npubkey = "admin"
	}
	usr, err := caccount.GetUserInfo(npubkey)
	if err != nil {
		return GetAccountError, 0, 0
	}
	if assetId != btcId {
		err, unlockedBalance, lockedBalance = GetAssetBalance(usr, assetId)
		if err != nil {
			return err, 0, 0
		}
	} else {
		err, unlockedBalance, lockedBalance = GetBtcBalance(usr)
		if err != nil {
			return err, 0, 0
		}
	}
	return
}

func GetBalances(npubkey string) (*[]cModels.LockBalance, error) {
	if npubkey == FeeNpubkey {
		npubkey = "admin"
	}
	usr, err := caccount.GetUserInfo(npubkey)
	if err != nil {
		return nil, GetAccountError
	}
	//获取所有资产余额
	var balances []cModels.LockBalance
	if err = middleware.DB.Where("account_id = ? and asset_id != '00'", usr.LockAccount.ID).
		Find(&balances).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, ServiceError
	}
	return &balances, nil
}

func Lock(npubkey, lockedId, assetId string, amount float64) error {
	usr, err := caccount.GetUserInfo(npubkey)
	if err != nil {
		return fmt.Errorf("%w: %s", GetAccountError, err.Error())
	}
	mutex := GetLockPaymentMutex(usr.User.ID)
	mutex.Lock()
	defer mutex.Unlock()

	if amount <= 0 {
		btlLog.CUST.Error("amount <= 0,lockedId:%s,assetId:%s,amount:%f", lockedId, assetId, amount)
		return BadRequest
	}

	if assetId != btcId {
		err := LockAsset(usr, lockedId, assetId, amount)
		if err != nil {
			return err
		}
	} else {
		err := LockBTC(usr, lockedId, amount)
		if err != nil {
			return err
		}
	}
	return nil
}

func Unlock(npubkey, lockedId, assetId string, amount float64, version int) error {
	usr, err := caccount.GetUserInfo(npubkey)
	if err != nil {
		return GetAccountError
	}

	mutex := GetLockPaymentMutex(usr.User.ID)
	mutex.Lock()
	defer mutex.Unlock()

	if amount <= 0 {
		btlLog.CUST.Error("amount <= 0,lockedId:%s,assetId:%s,amount:%f", lockedId, assetId, amount)
		return BadRequest
	}

	if assetId != btcId {
		err := UnlockAsset(usr, lockedId, assetId, amount, version)
		if err != nil {
			return err
		}
	} else {
		err := UnlockBTC(usr, lockedId, amount, version)
		if err != nil {
			return err
		}
	}
	return nil
}

func TransferByUnlock(lockedId, npubkey, toNpubkey, assetId string, amount float64) error {
	if npubkey == FeeNpubkey {
		npubkey = "admin"
	}
	usr, err := caccount.GetUserInfo(npubkey)
	if err != nil {
		return GetAccountError
	}
	mutex := GetLockPaymentMutex(usr.User.ID)
	mutex.Lock()
	defer mutex.Unlock()

	if toNpubkey == FeeNpubkey {
		toNpubkey = "admin"
	}
	toUsr, err := caccount.GetUserInfo(toNpubkey)
	if err != nil {
		return RevNpubKeyNotFound
	}
	mutexTo := GetLockPaymentMutex(toUsr.User.ID)
	mutexTo.Lock()
	defer mutexTo.Unlock()

	if amount <= 0 {
		btlLog.CUST.Error("amount <= 0,lockedId:%s,assetId:%s,amount:%f", lockedId, assetId, amount)
		return BadRequest
	}

	if assetId != btcId {
		err := transferAsset(usr, lockedId, assetId, amount, toUsr)
		if err != nil {
			return err
		}
	} else {
		err := transferBTC(usr, lockedId, amount, toUsr)
		if err != nil {
			return err
		}
	}
	return nil
}

func TransferByLock(lockedId, npubkey, toNpubkey, assetId string, amount float64) error {
	if npubkey == FeeNpubkey {
		npubkey = "admin"
	}
	usr, err := caccount.GetUserInfo(npubkey)
	if err != nil {
		return GetAccountError
	}
	mutex := GetLockPaymentMutex(usr.User.ID)
	mutex.Lock()
	defer mutex.Unlock()

	if toNpubkey == FeeNpubkey {
		toNpubkey = "admin"
	}
	toUsr, err := caccount.GetUserInfo(toNpubkey)
	if err != nil {
		return RevNpubKeyNotFound
	}
	mutexTo := GetLockPaymentMutex(toUsr.User.ID)
	mutexTo.Lock()
	defer mutexTo.Unlock()
	if amount <= 0 {
		btlLog.CUST.Error("amount <= 0,lockedId:%s,assetId:%s,amount:%f", lockedId, assetId, amount)
		return BadRequest
	}
	if assetId != btcId {
		err := transferLockedAsset(usr, lockedId, assetId, amount, toUsr)
		if err != nil {
			return err
		}
	} else {
		err := transferLockedBTC(usr, lockedId, amount, toUsr)
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckLockId(lockedId string) error {
	bill := cModels.LockBill{}
	if err := middleware.DB.Where("locked_id = ?", lockedId).First(&bill).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return RepeatedLockId
}

func GetLockPaymentMutex(userId uint) *sync.Mutex {
	mutexKey := fmt.Sprintf("%s_%d", LockMutexKey, userId)
	return custodyMutex.GetCustodyMutex(mutexKey)
}

// ListTransferBTC 列出转账记录
func ListTransferBTC(usr *caccount.UserInfo, assetId string, page, pageSize, away int) ([]cModels.LockBill, error) {
	var err error
	var bills []cModels.LockBill
	offset := (page - 1) * pageSize
	limit := pageSize
	tx := middleware.DB
	q := tx.Where("account_id =? AND asset_id = ?  ", usr.LockAccount.ID, assetId)
	switch away {
	case 0:
		q.Where("bill_type = ? OR bill_type = ? ",
			cModels.LockBillTypeLock,
			cModels.LockBillTypeAward)
	case 1:
		q.Where("bill_type = ? OR bill_type = ? OR bill_type = ?",
			cModels.LockBillTypeTransferByLockAsset,
			cModels.LockBillTypeUnlock,
			cModels.LockBillTypeTransferByUnlockAsset)
	default:
		q.Where("bill_type != ? ", cModels.LockErr)
	}
	if err = q.Order("id desc").Offset(offset).
		Limit(limit).Find(&bills).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, ServiceError
	}
	return bills, nil
}
