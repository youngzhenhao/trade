package lockPayment

import caccount "trade/services/custodyAccount/account"

const (
	FeeNpubkey = "btlong"
	btcId      = "00"
)

func GetBalance(npubkey, assetId string) (err error, unlockedBalance float64, lockedBalance float64) {
	usr, err := caccount.GetUserInfo(npubkey)
	if err != nil {
		return err, 0, 0
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
		return err
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
		return err
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
		return err
	}
	toUsr, err := caccount.GetUserInfo(toNpubkey)
	if err != nil {
		return err
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
		return err
	}
	toUsr, err := caccount.GetUserInfo(toNpubkey)
	if err != nil {
		return err
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
