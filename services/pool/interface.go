package pool

import (
	"gorm.io/gorm"
	"math/big"
)

func CreatePoolAccount(tx *gorm.DB, pairId uint, allowTokens []string) (err error) {
	// TODO 1.
	return err
}

// @Description: token is the asset_id or the "sat"
func PoolAccountTransfer(tx *gorm.DB, pairId uint, username string, token string, _amount *big.Int) (recordId uint, err error) {
	// TODO 2.
	return recordId, err
}

func TransferToPoolAccount(tx *gorm.DB, username string, pairId uint, token string, _amount *big.Int) (recordId uint, err error) {
	// TODO 3.
	return recordId, err
}

func GetPoolAccountRecords(pairId uint, limit uint64, offset uint64) (records *[]any, err error) {
	// TODO 4.
	return records, err
}

// TODO 5.GetPoolAccountInfo

// TODO 6.LockPoolAccount

// TODO 7.UnLockPoolAccount
