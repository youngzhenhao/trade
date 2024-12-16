package pool

import (
	"gorm.io/gorm"
	"math/big"
	"trade/models/custodyModels/pAccount"
	"trade/services/custodyAccount/poolAccount"
)

func CreatePoolAccount(tx *gorm.DB, pairId uint, allowTokens []string) (err error) {
	transTokens := make([]string, 0)
	for _, token := range allowTokens {
		if token == TokenSatTag {
			transTokens = append(transTokens, "00")
		} else {
			transTokens = append(transTokens, token)
		}
	}

	return poolAccount.CreatePoolAccount(tx, pairId, transTokens)
}

// @Description: token is the asset_id or the "sat"
func PoolAccountTransfer(tx *gorm.DB, pairId uint, username string, token string, _amount *big.Int, transferDescription string) (recordId uint, err error) {
	return poolAccount.PAccountToUserPay(tx, username, pairId, token, _amount, transferDescription)
}

func TransferToPoolAccount(tx *gorm.DB, username string, pairId uint, token string, _amount *big.Int, transferDescription string) (recordId uint, err error) {
	return poolAccount.UserPayToPAccount(tx, pairId, username, token, _amount, transferDescription)
}

func GetPoolAccountRecords(pairId uint, limit int, offset int) (records *[]pAccount.PAccountBill, err error) {
	return poolAccount.GetAccountRecords(pairId, limit, offset)
}
func GetPoolAccountRecordsCount(pairId uint) (count int64, err error) {
	return poolAccount.GetAccountRecordCount(pairId)
}

func GetPoolAccountInfo(pairId uint) (info *poolAccount.PAccountInfo, err error) {
	return poolAccount.GetPoolAccountInfo(pairId)
}

// TODO 6.LockPoolAccount
func LockPoolAccount(tx *gorm.DB, pairId uint) (err error) {
	return poolAccount.LockPoolAccount(tx, pairId)
}

// TODO 7.UnLockPoolAccount
func UnLockPoolAccount(tx *gorm.DB, pairId uint) (err error) {
	return poolAccount.UnlockPoolAccount(tx, pairId)
}

// @Note: Transfer Sats only
func TransferWithdrawReward(username string, _amount *big.Int, transferDescription string) (recordId uint, err error) {
	return poolAccount.AwardSat(username, _amount, transferDescription)
}
