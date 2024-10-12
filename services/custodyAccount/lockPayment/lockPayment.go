package lockPayment

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"trade/middleware"
	cModels "trade/models/custodyModels"
	caccount "trade/services/custodyAccount/account"
)

const (
	FeeNpubkey = "btlong"
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

func Lock(npubkey, lockedId, assetId string, amount float64) error {
	usr, err := caccount.GetUserInfo(npubkey)
	if err != nil {
		return fmt.Errorf("%w: %s", GetAccountError, err.Error())
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

func Unlock(npubkey, lockedId, assetId string, amount float64) error {
	usr, err := caccount.GetUserInfo(npubkey)
	if err != nil {
		return GetAccountError
	}
	if assetId != btcId {
		err := UnlockAsset(usr, lockedId, assetId, amount)
		if err != nil {
			return err
		}
	} else {
		err := UnlockBTC(usr, lockedId, amount)
		if err != nil {
			return err
		}
	}
	return nil
}

func TransferByUnlock(lockedId, npubkey, toNpubkey, assetId string, amount float64) error {
	usr, err := caccount.GetUserInfo(npubkey)
	if err != nil {
		return GetAccountError
	}
	if toNpubkey == FeeNpubkey {
		toNpubkey = "admin"
	}
	toUsr, err := caccount.GetUserInfo(toNpubkey)
	if err != nil {
		return RevNpubKeyNotFound
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
	usr, err := caccount.GetUserInfo(npubkey)
	if err != nil {
		return GetAccountError
	}
	if toNpubkey == FeeNpubkey {
		toNpubkey = "admin"
	}
	toUsr, err := caccount.GetUserInfo(toNpubkey)
	if err != nil {
		return RevNpubKeyNotFound
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
