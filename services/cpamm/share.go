package cpamm

import (
	"errors"
	"gorm.io/gorm"
	"math/big"
	"strconv"
	"trade/middleware"
	"trade/utils"
)

type PoolShare struct {
	gorm.Model
	PairId      uint   `json:"pair_id" gorm:"index"`
	Name        string `json:"name" gorm:"varchar(255);index"`
	Symbol      string `json:"symbol" gorm:"varchar(255);index"`
	TotalSupply string `json:"total_supply" gorm:"varchar(255);index"`
}

func (s *PoolShare) BalanceOf(userId uint) (*PoolShareBalance, error) {
	shareBalance, err := balanceOf(userId, s.ID)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "balanceOf")
	}
	return shareBalance, nil
}

// TODO: mutex
func createShareBalance(userId uint, shareId uint) (*PoolShareBalance, error) {
	shareBalance := PoolShareBalance{
		UserId:  userId,
		ShareId: shareId,
		Balance: "0",
	}
	err := middleware.DB.Model(PoolShareBalance{}).Create(&shareBalance).Error
	return &shareBalance, err
}

// TODO: need to test
func balanceOf(userId uint, shareId uint) (*PoolShareBalance, error) {
	var shareBalance *PoolShareBalance
	err := middleware.DB.Model(&PoolShareBalance{}).Where("user_id = ? AND share_id = ?", userId, shareId).First(shareBalance).Error
	if err != nil {
		shareBalance, err = createShareBalance(userId, shareId)
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "createShareBalance")
		}
		return shareBalance, nil
	}
	return shareBalance, nil
}

type PoolShareBalance struct {
	gorm.Model
	UserId  uint   `json:"user_id" gorm:"index"`
	ShareId uint   `json:"share_id" gorm:"index"`
	Balance string `json:"balance" gorm:"varchar(255);index"`
}

type ShareRecordType int64

const (
	ShareMint ShareRecordType = iota
	ShareBurn
	// TODO: Not sure if it's needed
	ShareTransfer
)

type PoolShareRecord struct {
	gorm.Model
	UserId     uint            `json:"user_id" gorm:"index"`
	ShareId    uint            `json:"share_id" gorm:"index"`
	Amount     string          `json:"amount" gorm:"varchar(255);index"`
	RecordType ShareRecordType `json:"record_type" gorm:"index"`
}

func updateShareAndBalance(share *PoolShare, balance *PoolShareBalance) error {
	if share == nil {
		return errors.New("share is nil")
	}
	if balance == nil {
		return errors.New("balance is nil")
	}
	tx := middleware.DB.Begin()
	err := tx.Model(PoolShare{}).Save(share).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(PoolShareBalance{}).Save(balance).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// TODO: need to test
// TODO: consider return value instead of writing to db
func (s *PoolShare) _mint(toUserId uint, value string) error {
	_totalSupply, success := new(big.Int).SetString(s.TotalSupply, 10)
	var err error
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+s.TotalSupply+") "+strconv.FormatBool(success))
	}
	_value, success := new(big.Int).SetString(value, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+value+") "+strconv.FormatBool(success))
	}
	_totalSupply = _totalSupply.Add(_totalSupply, _value)
	//TODO: update share TotalSupply
	s.TotalSupply = _totalSupply.String()
	// TODO: here possibly created shareBalance
	shareBalance, err := s.BalanceOf(toUserId)
	if err != nil {
		return utils.AppendErrorInfo(err, "("+strconv.FormatUint(uint64(s.ID), 10)+")BalanceOf("+strconv.FormatUint(uint64(toUserId), 10)+")")
	}
	_balance, success := new(big.Int).SetString(shareBalance.Balance, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+shareBalance.Balance+") "+strconv.FormatBool(success))
	}
	_balance = _balance.Add(_balance, _value)
	//TODO: update shareBalance Balance
	shareBalance.Balance = _balance.String()
	return updateShareAndBalance(s, shareBalance)
}

func (s *PoolShare) _burn(fromUserId uint, value string) error {
	_totalSupply, success := new(big.Int).SetString(s.TotalSupply, 10)
	var err error
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+s.TotalSupply+") "+strconv.FormatBool(success))
	}
	_value, success := new(big.Int).SetString(value, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+value+") "+strconv.FormatBool(success))
	}
	_totalSupply = _totalSupply.Sub(_totalSupply, _value)
	// @dev: Validate total supply
	if _totalSupply.Sign() < 0 {
		return errors.New("total_supply(" + _totalSupply.String() + ") is negative")
	}
	//TODO: update share TotalSupply
	s.TotalSupply = _totalSupply.String()
	// TODO: here possibly created shareBalance
	shareBalance, err := s.BalanceOf(fromUserId)
	if err != nil {
		return utils.AppendErrorInfo(err, "("+strconv.FormatUint(uint64(s.ID), 10)+")BalanceOf("+strconv.FormatUint(uint64(fromUserId), 10)+")")
	}
	_balance, success := new(big.Int).SetString(shareBalance.Balance, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+shareBalance.Balance+") "+strconv.FormatBool(success))
	}
	_balance = _balance.Sub(_balance, _value)
	// @dev: Validate balance
	if _balance.Sign() < 0 {
		return errors.New("balance(" + _balance.String() + ") is negative")
	}
	//TODO: update shareBalance Balance
	shareBalance.Balance = _balance.String()
	return updateShareAndBalance(s, shareBalance)
}

// TODO: mutex
// TODO: create record of type mint
func (s *PoolShare) mint(userId uint) (liquidity string) {

	return ""
}

func (s *PoolShare) burn(to string) (amount0 string, amount1 string) {
	return "", ""
}

func (s *PoolShare) swap(amount0Out string, amount1Out string, to string, data string) {

}

func (s *PoolShare) skim(to string) {

}

func (s *PoolShare) _sync() {

}
