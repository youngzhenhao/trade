package pool

import (
	"errors"
	"gorm.io/gorm"
	"math/big"
	"strconv"
	"sync"
	"trade/utils"
)

var LockLpA map[string]*sync.Mutex

type PoolLpAwardBalance struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(255);uniqueIndex"`
	Balance  string `json:"balance" gorm:"type:varchar(255);index"`
}

type PoolLpAwardCumulative struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(255);uniqueIndex"`
	Amount   string `json:"amount" gorm:"type:varchar(255);index"`
}

func newLpAwardCumulative(username string, amount string) (lpAwardCumulative *PoolLpAwardCumulative, err error) {
	_amountFloat, success := new(big.Float).SetString(amount)
	if !success {
		return new(PoolLpAwardCumulative), errors.New("amount SetString(" + amount + ") " + strconv.FormatBool(success))
	}
	if _amountFloat.Sign() < 0 {
		return new(PoolLpAwardCumulative), errors.New("invalid balance(" + _amountFloat.String() + ")")
	}
	_lpAwardCumulative := PoolLpAwardCumulative{
		Username: username,
		Amount:   amount,
	}
	return &_lpAwardCumulative, nil
}

type AwardType int64

const (
	SwapAward AwardType = iota
)

type PoolLpAwardRecord struct {
	gorm.Model
	ShareId      uint      `json:"share_id" gorm:"index"`
	Username     string    `json:"username" gorm:"type:varchar(255);index"`
	Amount       string    `json:"amount" gorm:"type:varchar(255);index"`
	Fee          string    `json:"fee" gorm:"type:varchar(255);index"`
	AwardBalance string    `json:"award_balance" gorm:"type:varchar(255);index"`
	ShareBalance string    `json:"share_balance" gorm:"type:varchar(255);index"`
	TotalSupply  string    `json:"total_supply" gorm:"type:varchar(255);index"`
	SwapRecordId uint      `json:"swap_record_id" gorm:"index"`
	AwardType    AwardType `json:"award_type" gorm:"index"`
}

func newLpAwardBalance(username string, balance string) (lpAwardBalance *PoolLpAwardBalance, err error) {
	_balanceFloat, success := new(big.Float).SetString(balance)
	if !success {
		return new(PoolLpAwardBalance), errors.New("balance SetString(" + balance + ") " + strconv.FormatBool(success))
	}
	if _balanceFloat.Sign() < 0 {
		return new(PoolLpAwardBalance), errors.New("invalid balance(" + _balanceFloat.String() + ")")
	}
	_lpAwardBalance := PoolLpAwardBalance{
		Username: username,
		Balance:  balance,
	}
	return &_lpAwardBalance, nil
}

func createOrUpdateLpAwardBalance(tx *gorm.DB, username string, amount *big.Float) (previousAwardBalance string, err error) {
	var lpAwardBalance *PoolLpAwardBalance
	err = tx.Model(&PoolLpAwardBalance{}).Where("username = ?", username).First(&lpAwardBalance).Error
	if err != nil {
		// @dev: no lpAwardBalance
		lpAwardBalance, err = newLpAwardBalance(username, amount.String())
		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "newLpAwardBalance")
		}
		err = tx.Model(&PoolLpAwardBalance{}).Create(&lpAwardBalance).Error
		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "create lpAwardBalance")
		}
		previousAwardBalance = big.NewFloat(0).String()
	} else {
		oldBalance, success := new(big.Float).SetString(lpAwardBalance.Balance)
		if !success {
			return ZeroValue, errors.New("lpAwardBalance SetString(" + lpAwardBalance.Balance + ") " + strconv.FormatBool(success))
		}
		newBalance := new(big.Float).Add(oldBalance, amount)
		err = tx.Model(&PoolLpAwardBalance{}).Where("username = ?", username).
			Update("balance", newBalance.String()).Error
		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "update lpAwardBalance")
		}
		previousAwardBalance = oldBalance.String()
	}
	return previousAwardBalance, nil
}

func createOrUpdatePoolLpAwardCumulative(tx *gorm.DB, username string, amount *big.Float) (previousAwardCumulative string, err error) {
	var lpAwardCumulative *PoolLpAwardCumulative
	err = tx.Model(&PoolLpAwardCumulative{}).Where("username = ?", username).First(&lpAwardCumulative).Error
	if err != nil {
		// @dev: no lpAwardCumulative
		lpAwardCumulative, err = newLpAwardCumulative(username, amount.String())
		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "newLpAwardCumulative")
		}
		err = tx.Model(&PoolLpAwardCumulative{}).Create(&lpAwardCumulative).Error
		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "create lpAwardCumulative")
		}
		previousAwardCumulative = big.NewFloat(0).String()
	} else {
		oldAmount, success := new(big.Float).SetString(lpAwardCumulative.Amount)
		if !success {
			return ZeroValue, errors.New("lpAwardCumulative SetString(" + lpAwardCumulative.Amount + ") " + strconv.FormatBool(success))
		}
		newBalance := new(big.Float).Add(oldAmount, amount)
		err = tx.Model(&PoolLpAwardCumulative{}).Where("username = ?", username).
			Update("amount", newBalance.String()).Error
		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "update lpAwardCumulative")
		}
		previousAwardCumulative = oldAmount.String()
	}
	return previousAwardCumulative, nil
}

func newLpAwardRecord(shareId uint, username string, amount string, fee string, awardBalance string, shareBalance string, totalSupply string, swapRecordId uint, awardType AwardType) (lpAwardRecord *PoolLpAwardRecord, err error) {
	if shareId <= 0 {
		return new(PoolLpAwardRecord), errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	return &PoolLpAwardRecord{
		ShareId:      shareId,
		Username:     username,
		Amount:       amount,
		Fee:          fee,
		AwardBalance: awardBalance,
		ShareBalance: shareBalance,
		TotalSupply:  totalSupply,
		SwapRecordId: swapRecordId,
		AwardType:    awardType,
	}, nil
}

func createLpAwardRecord(tx *gorm.DB, shareId uint, username string, amount *big.Float, fee string, awardBalance string, shareBalance string, totalSupply string, swapRecordId uint, awardType AwardType) (err error) {
	if shareId <= 0 {
		return errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	var lpAwardRecord *PoolLpAwardRecord
	lpAwardRecord, err = newLpAwardRecord(shareId, username, amount.String(), fee, awardBalance, shareBalance, totalSupply, swapRecordId, awardType)
	if err != nil {
		return utils.AppendErrorInfo(err, "newLpAwardRecord")
	}
	err = tx.Model(&PoolLpAwardRecord{}).Create(&lpAwardRecord).Error
	if err != nil {
		return utils.AppendErrorInfo(err, "create lpAwardRecord")
	}
	return nil
}

func updateLpAwardBalanceAndRecordSwap(tx *gorm.DB, shareId uint, username string, amount *big.Float, fee string, shareBalance string, totalSupply string, swapRecordId uint) (err error) {
	if shareId <= 0 {
		return errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	var previousAwardBalance string
	previousAwardBalance, err = createOrUpdateLpAwardBalance(tx, username, amount)
	if err != nil {
		return utils.AppendErrorInfo(err, "createOrUpdateLpAwardBalance")
	}

	_, err = createOrUpdatePoolLpAwardCumulative(tx, username, amount)
	if err != nil {
		return utils.AppendErrorInfo(err, "createOrUpdatePoolLpAwardCumulative")
	}

	err = createLpAwardRecord(tx, shareId, username, amount, fee, previousAwardBalance, shareBalance, totalSupply, swapRecordId, SwapAward)
	if err != nil {
		return utils.AppendErrorInfo(err, "createLpAwardRecord")
	}
	return nil
}

type PoolWithdrawAwardRecord struct {
	gorm.Model
	Username                 string `json:"username" gorm:"type:varchar(255);index"`
	Amount                   string `json:"amount" gorm:"type:varchar(255);index"`
	WithdrawTransferRecordId uint   `json:"withdraw_transfer_record_id" gorm:"index"`
	AwardBalance             string `json:"award_balance" gorm:"type:varchar(255);index"`
}

func newWithdrawAwardRecord(username string, amount string, withdrawTransferRecordId uint, awardBalance string) (withdrawAwardRecord *PoolWithdrawAwardRecord, err error) {
	return &PoolWithdrawAwardRecord{
		Username:                 username,
		Amount:                   amount,
		WithdrawTransferRecordId: withdrawTransferRecordId,
		AwardBalance:             awardBalance,
	}, nil
}

func createWithdrawAwardRecord(tx *gorm.DB, username string, amount *big.Int, withdrawTransferRecordId uint, awardBalance string) (err error) {
	var withdrawAwardRecord *PoolWithdrawAwardRecord
	withdrawAwardRecord, err = newWithdrawAwardRecord(username, amount.String(), withdrawTransferRecordId, awardBalance)
	if err != nil {
		return utils.AppendErrorInfo(err, "newWithdrawAwardRecord")
	}
	err = tx.Model(&PoolWithdrawAwardRecord{}).Create(&withdrawAwardRecord).Error
	if err != nil {
		return utils.AppendErrorInfo(err, "create withdrawAwardRecord")
	}
	return nil
}

func _withdrawAward(tx *gorm.DB, username string, amount *big.Int) (oldBalance string, newBalance string, err error) {
	var lpAwardBalance *PoolLpAwardBalance
	err = tx.Model(&PoolLpAwardBalance{}).Where("username = ?", username).First(&lpAwardBalance).Error
	if err != nil {
		// @dev: no lpAwardBalance
		return ZeroValue, ZeroValue, errors.New("lpAwardBalance of " + username + " not found")
	} else {
		_oldBalance, success := new(big.Float).SetString(lpAwardBalance.Balance)
		if !success {
			return ZeroValue, ZeroValue, errors.New("lpAwardBalance SetString(" + lpAwardBalance.Balance + ") " + strconv.FormatBool(success))
		}
		_amountFloat := new(big.Float).SetInt(amount)
		_minWithdrawAwardSat := new(big.Float).SetUint64(uint64(MinWithdrawAwardSat))

		if _amountFloat.Cmp(_minWithdrawAwardSat) < 0 {
			return ZeroValue, ZeroValue, errors.New("insufficient _amountFloat(" + _amountFloat.String() + "), need ge " + _minWithdrawAwardSat.String())
		}

		_newBalance := new(big.Float).Sub(_oldBalance, _amountFloat)
		if _newBalance.Sign() < 0 {
			return ZeroValue, ZeroValue, errors.New("insufficient _newBalance(" + _newBalance.String() + "), _oldBalance(" + _oldBalance.String() + ") _amountFloat(" + _amountFloat.String() + ")")
		}

		err = tx.Model(&PoolLpAwardBalance{}).Where("username = ?", username).
			Update("balance", _newBalance.String()).Error
		if err != nil {
			return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "update lpAwardBalance")
		}
		newBalance = _newBalance.String()
		oldBalance = _oldBalance.String()
	}
	return oldBalance, newBalance, nil
}
