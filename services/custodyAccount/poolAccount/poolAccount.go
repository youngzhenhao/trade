package poolAccount

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math/big"
	"sync"
	"trade/btlLog"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
	"trade/models/custodyModels/pAccount"
	"trade/services/custodyAccount/account"
	"trade/services/custodyAccount/defaultAccount/custodyAssets"
	"trade/services/custodyAccount/defaultAccount/custodyBtc"
)

const btcId = "00"

var poolAccountMutex = sync.Mutex{}

func CreatePoolAccount(tx *gorm.DB, pairId uint, allowTokens []string) error {
	if tx == nil {
		return fmt.Errorf("tx is nil")
	}
	if allowTokens == nil || len(allowTokens) == 0 {
		return errors.New("allow tokens is empty")
	}
	var Account pAccount.PoolAccount
	err := tx.Where(pAccount.PoolAccount{PairId: pairId}).Attrs(pAccount.PoolAccount{Status: 1}).FirstOrCreate(&Account).Error
	if err != nil {
		return err
	}
	var pAccountAssets []pAccount.PAccountAssetId
	err = tx.Where("pool_account_id =?", Account.ID).Find(&pAccountAssets).Error
	if err != nil {
		return err
	}
	if len(pAccountAssets) > 0 {
		return fmt.Errorf("pool account is exist %d", Account.PairId)
	}
	for _, token := range allowTokens {
		pAccountAssets = append(pAccountAssets, pAccount.PAccountAssetId{
			PoolAccountId: Account.ID,
			AssetId:       token,
		})
	}
	err = tx.Create(&pAccountAssets).Error
	if err != nil {
		return err
	}
	return err
}

func GetPoolAccount(tx *gorm.DB, pairId uint) (*pAccount.PoolAccount, error) {
	if tx == nil {
		return nil, fmt.Errorf("tx is nil")
	}
	var Account pAccount.PoolAccount
	err := tx.Where(pAccount.PoolAccount{PairId: pairId}).First(&Account).Error
	if err != nil {
		return nil, err
	}
	if Account.Status != 1 {
		return nil, fmt.Errorf("pool account is locked")
	}
	return &Account, nil
}

func UserPayToPAccount(tx *gorm.DB, pairId uint, username string, token string, _amount *big.Int, transferDesc string) (uint, error) {
	if tx == nil {
		return 0, fmt.Errorf("tx is nil")
	}

	amount := bigIntToFloat64(_amount)

	poolAccount, err := GetPoolAccount(tx, pairId)
	if err != nil {
		return 0, err
	}

	usr, err := account.GetUserInfo(username)
	if err != nil {
		return 0, err
	}
	poolAccountMutex.Lock()
	defer poolAccountMutex.Unlock()

	b := getBillBalanceModel(usr, amount, token, models.AWAY_OUT, transferDesc)
	if err = tx.Create(b).Error; err != nil {
		return 0, ErrorDbError
	}
	//创建扣款记录
	switch token {
	case btcId:
		_, err = custodyBtc.LessBtcBalance(tx, usr, amount, b.ID, custodyModels.ChangeTypePayToPoolAccount)
		if err != nil {
			if errors.Is(err, custodyBtc.NotEnoughBalance) {
				return 0, ErrorNotEnoughBalance
			}
			btlLog.CUST.Error("LessBtcBalance error:%s", err)
			return 0, ErrorDbError
		}
	default:
		_, err = custodyAssets.LessAssetBalance(tx, usr, amount, b.ID, token, custodyModels.ChangeTypePayToPoolAccount)
		if err != nil {
			if errors.Is(err, custodyAssets.NotEnoughAssetBalance) {
				btlLog.CUST.Error("NotEnoughAssetBalance:%s, amount:%f", err, amount)
				return 0, ErrorNotEnoughBalance
			}
			btlLog.CUST.Error("LessAssetBalance error:%s", err)
			//todo 处理余额不足的情况
			return 0, ErrorDbError
		}
	}
	return addBalance(tx, poolAccount.ID, token, amount, username, transferDesc)
}

func PAccountToUserPay(tx *gorm.DB, username string, pairId uint, token string, _amount *big.Int, transferDesc string) (uint, error) {
	if tx == nil {
		return 0, fmt.Errorf("tx is nil")
	}
	amount := bigIntToFloat64(_amount)

	poolAccount, err := GetPoolAccount(tx, pairId)
	if err != nil {
		return 0, err
	}
	usr, err := account.GetUserInfo(username)
	if err != nil {
		return 0, err
	}
	poolAccountMutex.Lock()
	defer poolAccountMutex.Unlock()

	b := getBillBalanceModel(usr, amount, token, models.AWAY_IN, transferDesc)
	if err = tx.Create(b).Error; err != nil {
		return 0, ErrorDbError
	}
	//创建收款记录
	switch token {
	case btcId:
		_, err = custodyBtc.AddBtcBalance(tx, usr, amount, b.ID, custodyModels.ChangeTypeReceiveFromPoolAccount)
		if err != nil {
			btlLog.CUST.Error("AddBtcBalance error:%s", err)
			return 0, ErrorDbError
		}
	default:
		_, err = custodyAssets.AddAssetBalance(tx, usr, amount, b.ID, token, custodyModels.ChangeTypeReceiveFromPoolAccount)
		if err != nil {
			btlLog.CUST.Error("AddAssetBalance error:%s", err)
			return 0, ErrorDbError
		}
	}
	return lessBalance(tx, poolAccount.ID, token, amount, username, transferDesc)
}

func GetAccountRecords(pairId uint, limit, offset int) (*[]pAccount.PAccountBill, error) {
	db := middleware.DB
	poolAccount, err := GetPoolAccount(db, pairId)
	if err != nil {
		return nil, err
	}
	var bills []pAccount.PAccountBill
	err = db.Where("pool_account_id = ?", poolAccount.ID).Offset(offset).Limit(limit).Find(&bills).Error
	if err != nil {
		return nil, err
	}
	return &bills, nil
}

func GetAccountRecordCount(pairId uint) (int64, error) {
	db := middleware.DB
	poolAccount, err := GetPoolAccount(db, pairId)
	if err != nil {
		return 0, err
	}
	var total int64
	err = db.Table("custody_pool_account_bills").Where("pool_account_id = ?", poolAccount.ID).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func LockPoolAccount(tx *gorm.DB, pairId uint) error {
	if tx == nil {
		return fmt.Errorf("tx is nil")
	}
	return tx.Table("custody_pool_accounts").Where("pair_id = ?", pairId).Update("status", 0).Error
}

func UnlockPoolAccount(tx *gorm.DB, pairId uint) error {
	if tx == nil {
		return fmt.Errorf("tx is nil")
	}
	return tx.Table("custody_pool_accounts").Where("pair_id = ?", pairId).Update("status", 1).Error
}

func AwardSat(username string, _amount *big.Int, transferDesc string) (uint, error) {
	amount := bigIntToFloat64(_amount)
	usr, err := account.GetUserInfo(username)
	if err != nil {
		return 0, err
	}
	awardType := "swapLP"
	award, err := custodyBtc.PutInAward(usr, "", int(amount), &awardType, transferDesc)
	if err != nil {
		return 0, err
	}
	return award.ID, nil
}

type PAccountInfo struct {
	Status   uint
	Balances *[]pAccount.PAccountBalance
}

func GetPoolAccountInfo(pairId uint) (*PAccountInfo, error) {
	db := middleware.DB
	var info PAccountInfo
	poolAccount, err := GetPoolAccount(db, pairId)
	if err != nil {
		return nil, err
	}
	info.Status = poolAccount.Status
	var balances []pAccount.PAccountBalance
	err = db.Where("pool_account_id = ?", poolAccount.ID).Find(&balances).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		btlLog.CUST.Error("AddBtcBalance error: ", err)
		return nil, err
	}
	info.Balances = &balances
	return &info, nil
}

func bigIntToFloat64(bi *big.Int) float64 {
	// 将*big.Int转换为int64
	intValue := bi.Int64()
	// 将int64转换为float64并返回
	return float64(intValue)
}

func getBillBalanceModel(usr *account.UserInfo, amount float64, assetId string, away models.BalanceAway, transferDesc string) *models.Balance {
	ba := models.Balance{}

	var i string
	var typeExt models.BalanceTypeExtList
	if away == models.AWAY_OUT {
		i = "PayToPoolAccount"
		typeExt = models.BTExtPayToPoolAccount
	} else {
		i = "ReceiveFromPoolAccount"
		typeExt = models.BTExtReceivePoolAccount
	}
	ba.AccountId = usr.Account.ID
	ba.Amount = amount
	ba.AssetId = &assetId
	ba.Unit = models.UNIT_ASSET_NORMAL
	if assetId == btcId {
		ba.Unit = models.UNIT_SATOSHIS
	}
	ba.BillType = models.BillTypePoolAccount
	ba.Away = away
	ba.Invoice = &i
	ba.PaymentHash = &transferDesc
	ba.State = models.STATE_SUCCESS
	ba.TypeExt = &models.BalanceTypeExt{Type: typeExt}
	return &ba
}
