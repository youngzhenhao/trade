package pool

import (
	"errors"
	"gorm.io/gorm"
	"math/big"
	"strconv"
	"trade/middleware"
	"trade/utils"
)

// Share

type Share struct {
	gorm.Model
	PairId      uint   `json:"pair_id" gorm:"uniqueIndex"`
	TotalSupply string `json:"total_supply" gorm:"type:varchar(255);index"`
}

func getShare(pairId uint) (*Share, error) {
	var share Share
	err := middleware.DB.Where("pair_id = ?", pairId).First(&share).Error
	if err != nil {
		return new(Share), err
	}
	return &share, nil
}

func createShare(pairId uint, totalSupply string) (err error) {
	_totalSupply, success := new(big.Int).SetString(totalSupply, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+totalSupply+") "+strconv.FormatBool(success))
	}
	if _totalSupply.Sign() < 0 {
		return errors.New("invalid totalSupply(" + _totalSupply.String() + ")")
	}
	share := Share{
		PairId:      pairId,
		TotalSupply: totalSupply,
	}
	return middleware.DB.Create(&share).Error
}

func updateShare(pairId uint, totalSupply string) (err error) {
	_totalSupply, success := new(big.Int).SetString(totalSupply, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+totalSupply+") "+strconv.FormatBool(success))
	}
	if _totalSupply.Sign() <= 0 {
		return errors.New("invalid totalSupply(" + _totalSupply.String() + ")")
	}
	return middleware.DB.
		Model(&Share{}).
		Where("pair_id = ?", pairId).
		Update("total_supply", totalSupply).
		Error
}

func NewShare(pairId uint, totalSupply string) (share *Share, err error) {
	_totalSupply, success := new(big.Int).SetString(totalSupply, 10)
	if !success {
		return new(Share), utils.AppendErrorInfo(err, "SetString("+totalSupply+") "+strconv.FormatBool(success))
	}
	if _totalSupply.Sign() < 0 {
		return new(Share), errors.New("invalid totalSupply(" + _totalSupply.String() + ")")
	}
	_share := Share{
		PairId:      pairId,
		TotalSupply: totalSupply,
	}
	return &_share, nil
}

// ShareBalance

type ShareBalance struct {
	gorm.Model
	ShareId  uint   `json:"share_id" gorm:"uniqueIndex:idx_share_id_username"`
	Username string `json:"username" gorm:"type:varchar(255);uniqueIndex:idx_share_id_username"`
	Balance  string `json:"balance" gorm:"type:varchar(255);index"`
}

func getShareBalance(shareId uint, username string) (*ShareBalance, error) {
	var shareBalance ShareBalance
	err := middleware.DB.Where("share_id = ? AND username = ?", shareId, username).First(&shareBalance).Error
	if err != nil {
		return new(ShareBalance), err
	}
	return &shareBalance, nil
}

func createShareBalance(shareId uint, username string, balance string) (err error) {
	_balance, success := new(big.Int).SetString(balance, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+balance+") "+strconv.FormatBool(success))
	}
	if _balance.Sign() < 0 {
		return errors.New("invalid balance(" + _balance.String() + ")")
	}
	shareBalance := ShareBalance{
		ShareId:  shareId,
		Username: username,
		Balance:  balance,
	}
	return middleware.DB.Create(&shareBalance).Error
}

func updateShareBalance(shareId uint, username string, balance string) (err error) {
	_balance, success := new(big.Int).SetString(balance, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+balance+") "+strconv.FormatBool(success))
	}
	if _balance.Sign() < 0 {
		return errors.New("invalid balance(" + _balance.String() + ")")
	}
	return middleware.DB.
		Model(&ShareBalance{}).
		Where("share_id = ? AND username = ?", shareId, username).
		Update("balance", balance).
		Error
}

func getShareBalanceIfNotExistCreate(shareId uint, username string) (*ShareBalance, error) {
	shareBalance, err := getShareBalance(shareId, username)
	if err != nil {
		shareBalance = &ShareBalance{
			ShareId:  shareId,
			Username: username,
			Balance:  ZeroValue,
		}
		err = middleware.DB.Create(shareBalance).Error
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "Create")
		}
	}
	return shareBalance, nil
}

func updateShareBalanceIfNotExistCreate(shareId uint, username string, balance string) (err error) {
	_balance, success := new(big.Int).SetString(balance, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+balance+") "+strconv.FormatBool(success))
	}
	if _balance.Sign() < 0 {
		return errors.New("invalid balance(" + _balance.String() + ")")
	}
	shareBalance, err := getShareBalance(shareId, username)
	if err != nil {
		shareBalance = &ShareBalance{
			ShareId:  shareId,
			Username: username,
			Balance:  balance,
		}
		err = middleware.DB.Create(shareBalance).Error
		if err != nil {
			return utils.AppendErrorInfo(err, "Create")
		}
	}
	return middleware.DB.
		Model(&ShareBalance{}).
		Where("share_id = ? AND username = ?", shareId, username).
		Update("balance", balance).
		Error
}

func NewShareBalance(shareId uint, username string, balance string) (shareBalance *ShareBalance, err error) {
	_balance, success := new(big.Int).SetString(balance, 10)
	if !success {
		return new(ShareBalance), utils.AppendErrorInfo(err, "SetString("+balance+") "+strconv.FormatBool(success))
	}
	if _balance.Sign() < 0 {
		return new(ShareBalance), errors.New("invalid balance(" + _balance.String() + ")")
	}
	_shareBalance := ShareBalance{
		ShareId:  shareId,
		Username: username,
		Balance:  balance,
	}
	return &_shareBalance, nil
}

// ShareRecord

type ShareRecordType int64

const (
	ShareMint ShareRecordType = iota
	ShareBurn
	ShareTransfer
)

type ShareRecord struct {
	gorm.Model
	ShareId    uint            `json:"share_id" gorm:"index"`
	Username   string          `json:"username" gorm:"type:varchar(255);index"`
	Amount     string          `json:"amount" gorm:"type:varchar(255);index"`
	RecordType ShareRecordType `json:"record_type" gorm:"index"`
}

func createShareRecord(shareId uint, username string, amount string, recordType ShareRecordType) error {
	return middleware.DB.Create(&ShareRecord{
		ShareId:    shareId,
		Username:   username,
		Amount:     amount,
		RecordType: recordType,
	}).Error
}

func NewShareRecord(shareId uint, username string, amount string, recordType ShareRecordType) (shareRecord *ShareRecord, err error) {
	return &ShareRecord{
		ShareId:    shareId,
		Username:   username,
		Amount:     amount,
		RecordType: recordType,
	}, nil
}
