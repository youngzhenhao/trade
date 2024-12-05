package pool

import (
	"errors"
	"gorm.io/gorm"
	"math/big"
	"strconv"
	"trade/utils"
)

type LpAwardBalance struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(255);uniqueIndex"`
	Balance  string `json:"balance" gorm:"type:varchar(255);index"`
}

type AwardType int64

const (
	SwapAward AwardType = iota
)

type LpAwardRecord struct {
	gorm.Model
	ShareId      uint      `json:"share_id" gorm:"index"`
	Username     string    `json:"username" gorm:"type:varchar(255);index"`
	Amount       string    `json:"amount" gorm:"type:varchar(255);index"`
	AwardBalance string    `json:"award_balance" gorm:"type:varchar(255);index"`
	ShareBalance string    `json:"share_balance" gorm:"type:varchar(255);index"`
	TotalSupply  string    `json:"total_supply" gorm:"type:varchar(255);index"`
	SwapRecordId uint      `json:"swap_record_id" gorm:"index"`
	AwardType    AwardType `json:"award_type" gorm:"index"`
}

func NewLpAwardBalance(username string, balance string) (lpAwardBalance *LpAwardBalance, err error) {
	_balanceFloat, success := new(big.Float).SetString(balance)
	if !success {
		return new(LpAwardBalance), errors.New("balance SetString(" + balance + ") " + strconv.FormatBool(success))
	}
	if _balanceFloat.Sign() < 0 {
		return new(LpAwardBalance), errors.New("invalid balance(" + _balanceFloat.String() + ")")
	}
	_lpAwardBalance := LpAwardBalance{
		Username: username,
		Balance:  balance,
	}
	return &_lpAwardBalance, nil
}

func CreateOrUpdateLpAwardBalance(tx *gorm.DB, username string, amount *big.Float) (previousAwardBalance string, err error) {
	var lpAwardBalance *LpAwardBalance
	err = tx.Model(&LpAwardBalance{}).Where("username = ?", username).First(&lpAwardBalance).Error
	if err != nil {
		// @dev: no lpAwardBalance
		lpAwardBalance, err = NewLpAwardBalance(username, amount.String())
		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "NewLpAwardBalance")
		}
		err = tx.Model(&LpAwardBalance{}).Create(&lpAwardBalance).Error
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
		err = tx.Model(&LpAwardBalance{}).Where("username = ?", username).
			Update("balance", newBalance.String()).Error
		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "update lpAwardBalance")
		}
		previousAwardBalance = oldBalance.String()
	}
	return previousAwardBalance, nil
}

func NewLpAwardRecord(shareId uint, username string, amount string, awardBalance string, shareBalance string, totalSupply string, swapRecordId uint, awardType AwardType) (lpAwardRecord *LpAwardRecord, err error) {
	if shareId <= 0 {
		return new(LpAwardRecord), errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	return &LpAwardRecord{
		ShareId:      shareId,
		Username:     username,
		Amount:       amount,
		AwardBalance: awardBalance,
		ShareBalance: shareBalance,
		TotalSupply:  totalSupply,
		SwapRecordId: swapRecordId,
		AwardType:    awardType,
	}, nil
}

func CreateLpAwardRecord(tx *gorm.DB, shareId uint, username string, amount *big.Float, awardBalance string, shareBalance string, totalSupply string, swapRecordId uint, awardType AwardType) (err error) {
	if shareId <= 0 {
		return errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	var lpAwardRecord *LpAwardRecord
	lpAwardRecord, err = NewLpAwardRecord(shareId, username, amount.String(), awardBalance, shareBalance, totalSupply, swapRecordId, awardType)
	if err != nil {
		return utils.AppendErrorInfo(err, "NewLpAwardRecord")
	}
	err = tx.Model(&LpAwardRecord{}).Create(&lpAwardRecord).Error
	if err != nil {
		return utils.AppendErrorInfo(err, "create lpAwardRecord")
	}
	return nil
}

func UpdateLpAwardBalanceAndRecordSwap(tx *gorm.DB, shareId uint, username string, amount *big.Float, shareBalance string, totalSupply string, swapRecordId uint) (err error) {
	if shareId <= 0 {
		return errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	var previousAwardBalance string
	previousAwardBalance, err = CreateOrUpdateLpAwardBalance(tx, username, amount)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateOrUpdateLpAwardBalance")
	}
	err = CreateLpAwardRecord(tx, shareId, username, amount, previousAwardBalance, shareBalance, totalSupply, swapRecordId, SwapAward)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateLpAwardRecord")
	}
	return nil
}
