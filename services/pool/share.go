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
	if pairId <= 0 {
		return errors.New("invalid pairId(" + strconv.Itoa(int(pairId)) + ")")
	}
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
	if pairId <= 0 {
		return new(Share), errors.New("invalid pairId(" + strconv.Itoa(int(pairId)) + ")")
	}
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

func _mintBig(_amount0 *big.Int, _amount1 *big.Int, _reserve0 *big.Int, _reserve1 *big.Int, _totalSupply *big.Int) (_liquidity *big.Int, err error) {
	_minLiquidity := new(big.Int).SetUint64(uint64(MinLiquidity))
	if _totalSupply.Sign() == 0 {
		_liquidity = new(big.Int).Sub(new(big.Int).Sqrt(new(big.Int).Mul(_amount0, _amount1)), _minLiquidity)
	} else {
		_liquidity0 := new(big.Int).Div(new(big.Int).Mul(_amount0, _totalSupply), _reserve0)
		_liquidity1 := new(big.Int).Div(new(big.Int).Mul(_amount1, _totalSupply), _reserve1)
		_liquidity = minBigInt(_liquidity0, _liquidity1)
	}
	if _liquidity.Sign() <= 0 {
		return new(big.Int), errors.New("insufficientLiquidityMinted(" + _liquidity.String() + ")")
	}
	return _liquidity, nil
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
	if shareId <= 0 {
		return errors.New("invalid shareId(" + strconv.Itoa(int(shareId)) + ")")
	}
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
	if shareId <= 0 {
		return new(ShareBalance), errors.New("invalid shareId(" + strconv.Itoa(int(shareId)) + ")")
	}
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
	if shareId <= 0 {
		return errors.New("invalid shareId(" + strconv.Itoa(int(shareId)) + ")")
	}
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
	if shareId <= 0 {
		return new(ShareBalance), errors.New("invalid shareId(" + strconv.Itoa(int(shareId)) + ")")
	}
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

func CreateOrUpdateShareBalance(tx *gorm.DB, share *Share, username string, _liquidity *big.Int) (previousShare string, err error) {
	if share.ID <= 0 {
		return ZeroValue, errors.New("invalid shareId(" + strconv.Itoa(int(share.ID)) + ")")
	}
	var shareBalance *ShareBalance
	err = tx.Model(&ShareBalance{}).Where("share_id = ? AND username = ?", share.ID, username).First(&shareBalance).Error
	if err != nil {
		// @dev: no shareBalance
		shareBalance, err = NewShareBalance(share.ID, username, _liquidity.String())
		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "NewShareBalance")
		}
		err = tx.Model(&ShareBalance{}).Create(&shareBalance).Error
		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "create shareBalance")
		}
		previousShare = big.NewInt(0).String()
	} else {
		oldBalance, success := new(big.Int).SetString(shareBalance.Balance, 10)
		if !success {
			return ZeroValue, utils.AppendErrorInfo(err, "SetString("+shareBalance.Balance+") "+strconv.FormatBool(success))
		}
		newBalance := new(big.Int).Add(oldBalance, _liquidity)
		err = tx.Model(&ShareBalance{}).Where("share_id = ? AND username = ?", share.ID, username).
			Update("balance", newBalance.String()).Error
		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "update shareBalance")
		}
		previousShare = oldBalance.String()
	}
	return previousShare, nil
}

// ShareRecord

type ShareRecordType int64

const (
	AddLiquidityShareMint ShareRecordType = iota
	RemoveLiquidityShareBurn
	ShareTransfer
)

// TODO: record token transfer Id
type ShareRecord struct {
	gorm.Model
	ShareId     uint            `json:"share_id" gorm:"index"`
	Username    string          `json:"username" gorm:"type:varchar(255);index"`
	Liquidity   string          `json:"liquidity" gorm:"type:varchar(255);index"`
	Reserve0    string          `json:"reserve0" gorm:"type:varchar(255);index"`
	Reserve1    string          `json:"reserve1" gorm:"type:varchar(255);index"`
	Amount0     string          `json:"amount0" gorm:"type:varchar(255);index"`
	Amount1     string          `json:"amount1" gorm:"type:varchar(255);index"`
	ShareSupply string          `json:"share_supply" gorm:"type:varchar(255);index"`
	ShareAmt    string          `json:"share_amt" gorm:"type:varchar(255);index"`
	IsFirstMint bool            `json:"is_first_mint" gorm:"index"`
	RecordType  ShareRecordType `json:"record_type" gorm:"index"`
}

func createShareRecord(shareId uint, username string, liquidity string, reserve0 string, reserve1 string, amount0 string, amount1 string, shareSupply string, shareAmt string, isFirstMint bool, recordType ShareRecordType) (err error) {
	if shareId <= 0 {
		return errors.New("invalid shareId(" + strconv.Itoa(int(shareId)) + ")")
	}
	return middleware.DB.Create(&ShareRecord{
		ShareId:     shareId,
		Username:    username,
		Liquidity:   liquidity,
		Reserve0:    reserve0,
		Reserve1:    reserve1,
		Amount0:     amount0,
		Amount1:     amount1,
		ShareSupply: shareSupply,
		ShareAmt:    shareAmt,
		IsFirstMint: isFirstMint,
		RecordType:  recordType,
	}).Error
}

func NewShareRecord(shareId uint, username string, liquidity string, reserve0 string, reserve1 string, amount0 string, amount1 string, shareSupply string, shareAmt string, isFirstMint bool, recordType ShareRecordType) (shareRecord *ShareRecord, err error) {
	if shareId <= 0 {
		return new(ShareRecord), errors.New("invalid shareId(" + strconv.Itoa(int(shareId)) + ")")
	}
	return &ShareRecord{
		ShareId:     shareId,
		Username:    username,
		Liquidity:   liquidity,
		Reserve0:    reserve0,
		Reserve1:    reserve1,
		Amount0:     amount0,
		Amount1:     amount1,
		ShareSupply: shareSupply,
		ShareAmt:    shareAmt,
		IsFirstMint: isFirstMint,
		RecordType:  recordType,
	}, nil
}

func CreateShareRecord(tx *gorm.DB, share *Share, username string, _liquidity *big.Int, reserve0 string, reserve1 string, amount0 string, amount1 string, shareSupply string, shareAmt string, isFirstMint bool) (err error) {
	if share.ID <= 0 {
		return errors.New("invalid shareId(" + strconv.Itoa(int(share.ID)) + ")")
	}
	var shareRecord *ShareRecord
	shareRecord, err = NewShareRecord(share.ID, username, _liquidity.String(), reserve0, reserve1, amount0, amount1, shareSupply, shareAmt, isFirstMint, AddLiquidityShareMint)
	if err != nil {
		return utils.AppendErrorInfo(err, "NewShareRecord")
	}
	err = tx.Model(&ShareRecord{}).Create(&shareRecord).Error
	if err != nil {
		return utils.AppendErrorInfo(err, "create shareRecord")
	}
	return nil
}

func UpdateShareBalanceAndRecord(tx *gorm.DB, share *Share, username string, _liquidity *big.Int, reserve0 string, reserve1 string, amount0 string, amount1 string, shareSupply string, isFirstMint bool) (err error) {
	if share.ID <= 0 {
		return errors.New("invalid shareId(" + strconv.Itoa(int(share.ID)) + ")")
	}
	var previousShare string
	previousShare, err = CreateOrUpdateShareBalance(tx, share, username, _liquidity)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateOrUpdateShareBalance")
	}
	err = CreateShareRecord(tx, share, username, _liquidity, reserve0, reserve1, amount0, amount1, shareSupply, previousShare, isFirstMint)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateShareRecord")
	}
	return nil
}
