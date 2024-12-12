package pool

import (
	"errors"
	"gorm.io/gorm"
	"math/big"
	"strconv"
	"trade/middleware"
	"trade/utils"
)

// PoolShare

type PoolShare struct {
	gorm.Model
	PairId      uint   `json:"pair_id" gorm:"uniqueIndex"`
	TotalSupply string `json:"total_supply" gorm:"type:varchar(255);index"`
}

func getShare(pairId uint) (*PoolShare, error) {
	var share PoolShare
	err := middleware.DB.Where("pair_id = ?", pairId).First(&share).Error
	if err != nil {
		return new(PoolShare), err
	}
	return &share, nil
}

func createShare(pairId uint, totalSupply string) (err error) {
	if pairId <= 0 {
		return errors.New("invalid pairId(" + strconv.Itoa(int(pairId)) + ")")
	}
	_totalSupply, success := new(big.Int).SetString(totalSupply, 10)
	if !success {
		return errors.New("totalSupply SetString(" + totalSupply + ") " + strconv.FormatBool(success))
	}
	if _totalSupply.Sign() < 0 {
		return errors.New("invalid totalSupply(" + _totalSupply.String() + ")")
	}
	share := PoolShare{
		PairId:      pairId,
		TotalSupply: totalSupply,
	}
	return middleware.DB.Create(&share).Error
}

func updateShare(pairId uint, totalSupply string) (err error) {
	_totalSupply, success := new(big.Int).SetString(totalSupply, 10)
	if !success {
		return errors.New("totalSupply SetString(" + totalSupply + ") " + strconv.FormatBool(success))
	}
	if _totalSupply.Sign() <= 0 {
		return errors.New("invalid totalSupply(" + _totalSupply.String() + ")")
	}
	return middleware.DB.
		Model(&PoolShare{}).
		Where("pair_id = ?", pairId).
		Update("total_supply", totalSupply).
		Error
}

func _newShare(pairId uint, totalSupply string) (share *PoolShare, err error) {
	if pairId <= 0 {
		return new(PoolShare), errors.New("invalid pairId(" + strconv.Itoa(int(pairId)) + ")")
	}
	_totalSupply, success := new(big.Int).SetString(totalSupply, 10)
	if !success {
		return new(PoolShare), errors.New("totalSupply SetString(" + totalSupply + ") " + strconv.FormatBool(success))
	}
	if _totalSupply.Sign() < 0 {
		return new(PoolShare), errors.New("invalid totalSupply(" + _totalSupply.String() + ")")
	}
	_share := PoolShare{
		PairId:      pairId,
		TotalSupply: totalSupply,
	}
	return &_share, nil
}

func _mintBig(_amount0 *big.Int, _amount1 *big.Int, _reserve0 *big.Int, _reserve1 *big.Int, _totalSupply *big.Int, isTokenZeroSat bool) (_liquidity *big.Int, err error) {

	if isTokenZeroSat {
		// cmp with minimum liquidity sat
		_minLiquiditySat := new(big.Int).SetUint64(uint64(MinAddLiquiditySat))
		if _amount0.Cmp(_minLiquiditySat) < 0 {
			err = errors.New("insufficient amount0 Sat(" + _reserve0.String() + "), need " + _minLiquiditySat.String())
			return new(big.Int), err
		}
	}

	if _totalSupply.Sign() == 0 {
		// @dev: Make sure that liquidity is not completely removed
		_minLiquidity := new(big.Int).SetUint64(uint64(MinLiquidity))
		_liquidity = new(big.Int).Sub(new(big.Int).Sqrt(new(big.Int).Mul(_amount0, _amount1)), _minLiquidity)
		//fmt.Printf("[_mintBig] _liquidity: %v;(sqrt(_amount0:%v * _amount1:%v) - _minLiquidity:%v)\n", _liquidity, _amount0, _amount1, _minLiquidity)
	} else {
		_liquidity0 := new(big.Int).Div(new(big.Int).Mul(_amount0, _totalSupply), _reserve0)
		//fmt.Printf("[_mintBig] _liquidity0: %v;(_amount0:%v * _totalSupply:%v / _reserve0:%v)\n", _liquidity0, _amount0, _totalSupply, _reserve0)
		_liquidity1 := new(big.Int).Div(new(big.Int).Mul(_amount1, _totalSupply), _reserve1)
		//fmt.Printf("[_mintBig] _liquidity1: %v;(_amount1:%v * _totalSupply:%v / _reserve1:%v)\n", _liquidity1, _amount1, _totalSupply, _reserve1)
		_liquidity = minBigInt(_liquidity0, _liquidity1)
	}
	if _liquidity.Sign() <= 0 {
		return new(big.Int), errors.New("insufficientLiquidityMinted(" + _liquidity.String() + ")")
	}
	return _liquidity, nil
}

func _burnBig(_reserve0 *big.Int, _reserve1 *big.Int, _totalSupply *big.Int, _liquidity *big.Int, feeK uint16) (_amount0 *big.Int, _amount1 *big.Int, err error) {
	// TODO: consider if allow user to burn all liquidity, now it's not allowed
	if _liquidity.Cmp(_totalSupply) >= 0 {
		return new(big.Int), new(big.Int), errors.New("insufficientLiquidityBurned _liquidity(" + _liquidity.String() + ") _totalSupply(" + _totalSupply.String() + ")")
	}

	k := new(big.Int).SetUint64(uint64(feeK))
	oneThousand := new(big.Int).SetUint64(1000)

	// x_0 * S * (1000 - k)
	_amount0Numerator := new(big.Int).Mul(new(big.Int).Mul(_liquidity, _reserve0), new(big.Int).Sub(oneThousand, k))
	// T * 1000
	_amount0Denominator := new(big.Int).Mul(_totalSupply, oneThousand)
	_amount0 = new(big.Int).Div(_amount0Numerator, _amount0Denominator)
	if !(_amount0.Sign() > 0) {
		return new(big.Int), new(big.Int), errors.New("insufficientAmount0Burned _amount0(" + _amount0.String() + ")")
	}

	// y_0 * S * (1000 - k)
	_amount1Numerator := new(big.Int).Mul(new(big.Int).Mul(_liquidity, _reserve1), new(big.Int).Sub(oneThousand, k))
	// T * 1000
	_amount1Denominator := new(big.Int).Mul(_totalSupply, oneThousand)
	_amount1 = new(big.Int).Div(_amount1Numerator, _amount1Denominator)
	if !(_amount1.Sign() > 0) {
		return new(big.Int), new(big.Int), errors.New("insufficientAmount1Burned _amount1(" + _amount1.String() + ")")
	}

	return _amount0, _amount1, nil
}

// PoolShareBalance

type PoolShareBalance struct {
	gorm.Model
	ShareId  uint   `json:"share_id" gorm:"uniqueIndex:idx_share_id_username"`
	Username string `json:"username" gorm:"type:varchar(255);uniqueIndex:idx_share_id_username"`
	Balance  string `json:"balance" gorm:"type:varchar(255);index"`
}

func getShareBalance(shareId uint, username string) (*PoolShareBalance, error) {
	var shareBalance PoolShareBalance
	err := middleware.DB.Where("share_id = ? AND username = ?", shareId, username).First(&shareBalance).Error
	if err != nil {
		return new(PoolShareBalance), err
	}
	return &shareBalance, nil
}

func createShareBalance(shareId uint, username string, balance string) (err error) {
	if shareId <= 0 {
		return errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	_balance, success := new(big.Int).SetString(balance, 10)
	if !success {
		return errors.New("balance SetString(" + balance + ") " + strconv.FormatBool(success))
	}
	if _balance.Sign() < 0 {
		return errors.New("invalid balance(" + _balance.String() + ")")
	}
	shareBalance := PoolShareBalance{
		ShareId:  shareId,
		Username: username,
		Balance:  balance,
	}
	return middleware.DB.Create(&shareBalance).Error
}

func updateShareBalance(shareId uint, username string, balance string) (err error) {
	_balance, success := new(big.Int).SetString(balance, 10)
	if !success {
		return errors.New("balance SetString(" + balance + ") " + strconv.FormatBool(success))
	}
	if _balance.Sign() < 0 {
		return errors.New("invalid balance(" + _balance.String() + ")")
	}
	return middleware.DB.
		Model(&PoolShareBalance{}).
		Where("share_id = ? AND username = ?", shareId, username).
		Update("balance", balance).
		Error
}

func getShareBalanceIfNotExistCreate(shareId uint, username string) (*PoolShareBalance, error) {
	if shareId <= 0 {
		return new(PoolShareBalance), errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	shareBalance, err := getShareBalance(shareId, username)
	if err != nil {
		shareBalance = &PoolShareBalance{
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
		return errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	_balance, success := new(big.Int).SetString(balance, 10)
	if !success {
		return errors.New("balance SetString(" + balance + ") " + strconv.FormatBool(success))
	}
	if _balance.Sign() < 0 {
		return errors.New("invalid balance(" + _balance.String() + ")")
	}
	shareBalance, err := getShareBalance(shareId, username)
	if err != nil {
		shareBalance = &PoolShareBalance{
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
		Model(&PoolShareBalance{}).
		Where("share_id = ? AND username = ?", shareId, username).
		Update("balance", balance).
		Error
}

func newShareBalance(shareId uint, username string, balance string) (shareBalance *PoolShareBalance, err error) {
	if shareId <= 0 {
		return new(PoolShareBalance), errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	_balance, success := new(big.Int).SetString(balance, 10)
	if !success {
		return new(PoolShareBalance), errors.New("balance SetString(" + balance + ") " + strconv.FormatBool(success))
	}
	if _balance.Sign() < 0 {
		return new(PoolShareBalance), errors.New("invalid balance(" + _balance.String() + ")")
	}
	_shareBalance := PoolShareBalance{
		ShareId:  shareId,
		Username: username,
		Balance:  balance,
	}
	return &_shareBalance, nil
}

func createOrUpdateShareBalance(tx *gorm.DB, shareId uint, username string, _liquidity *big.Int) (previousShare string, err error) {
	if shareId <= 0 {
		return ZeroValue, errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	var shareBalance *PoolShareBalance
	err = tx.Model(&PoolShareBalance{}).Where("share_id = ? AND username = ?", shareId, username).First(&shareBalance).Error
	if err != nil {
		// @dev: no shareBalance
		shareBalance, err = newShareBalance(shareId, username, _liquidity.String())
		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "newShareBalance")
		}
		err = tx.Model(&PoolShareBalance{}).Create(&shareBalance).Error
		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "create shareBalance")
		}
		previousShare = big.NewInt(0).String()
	} else {
		oldBalance, success := new(big.Int).SetString(shareBalance.Balance, 10)
		if !success {
			return ZeroValue, errors.New("shareBalance SetString(" + shareBalance.Balance + ") " + strconv.FormatBool(success))
		}
		newBalance := new(big.Int).Add(oldBalance, _liquidity)
		err = tx.Model(&PoolShareBalance{}).Where("share_id = ? AND username = ?", shareId, username).
			Update("balance", newBalance.String()).Error
		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "update shareBalance")
		}
		previousShare = oldBalance.String()
	}
	return previousShare, nil
}

func updateShareBalanceBurn(tx *gorm.DB, shareId uint, username string, _liquidity *big.Int) (previousShare string, err error) {
	if shareId <= 0 {
		return ZeroValue, errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	var shareBalance *PoolShareBalance
	err = tx.Model(&PoolShareBalance{}).Where("share_id = ? AND username = ?", shareId, username).First(&shareBalance).Error
	if err != nil {
		// @dev: no shareBalance
		return ZeroValue, errors.New("no shareBalance")
	}

	oldBalance, success := new(big.Int).SetString(shareBalance.Balance, 10)
	if !success {
		return ZeroValue, errors.New("shareBalance SetString(" + shareBalance.Balance + ") " + strconv.FormatBool(success))
	}
	newBalance := new(big.Int).Sub(oldBalance, _liquidity)
	if newBalance.Sign() < 0 {
		return ZeroValue, errors.New("insufficient newBalance(" + newBalance.String() + ")")
	}

	err = tx.Model(&PoolShareBalance{}).Where("share_id = ? AND username = ?", shareId, username).
		Update("balance", newBalance.String()).Error
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "update shareBalance")
	}
	previousShare = oldBalance.String()

	return previousShare, nil
}

// PoolShareRecord

type ShareRecordType int64

const (
	AddLiquidityShareMint ShareRecordType = iota
	RemoveLiquidityShareBurn
	ShareTransfer
)

// TODO: record token transfer Id
type PoolShareRecord struct {
	gorm.Model
	ShareId                uint            `json:"share_id" gorm:"index"`
	Username               string          `json:"username" gorm:"type:varchar(255);index"`
	Liquidity              string          `json:"liquidity" gorm:"type:varchar(255);index"`
	Token0TransferRecordId uint            `json:"token0_transfer_record_id" gorm:"index"`
	Token1TransferRecordId uint            `json:"token1_transfer_record_id" gorm:"index"`
	Reserve0               string          `json:"reserve0" gorm:"type:varchar(255);index"`
	Reserve1               string          `json:"reserve1" gorm:"type:varchar(255);index"`
	Amount0                string          `json:"amount0" gorm:"type:varchar(255);index"`
	Amount1                string          `json:"amount1" gorm:"type:varchar(255);index"`
	ShareSupply            string          `json:"share_supply" gorm:"type:varchar(255);index"`
	ShareAmt               string          `json:"share_amt" gorm:"type:varchar(255);index"`
	IsFirstMint            bool            `json:"is_first_mint" gorm:"index"`
	RecordType             ShareRecordType `json:"record_type" gorm:"index"`
}

func _createShareRecord(shareId uint, username string, liquidity string, reserve0 string, reserve1 string, amount0 string, amount1 string, shareSupply string, shareAmt string, isFirstMint bool, recordType ShareRecordType) (err error) {
	if shareId <= 0 {
		return errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	return middleware.DB.Create(&PoolShareRecord{
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

func newShareRecord(shareId uint, username string, liquidity string, token0TransferRecordId uint, token1TransferRecordId uint, reserve0 string, reserve1 string, amount0 string, amount1 string, shareSupply string, shareAmt string, isFirstMint bool, recordType ShareRecordType) (shareRecord *PoolShareRecord, err error) {
	if shareId <= 0 {
		return new(PoolShareRecord), errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	return &PoolShareRecord{
		ShareId:                shareId,
		Username:               username,
		Liquidity:              liquidity,
		Token0TransferRecordId: token0TransferRecordId,
		Token1TransferRecordId: token1TransferRecordId,
		Reserve0:               reserve0,
		Reserve1:               reserve1,
		Amount0:                amount0,
		Amount1:                amount1,
		ShareSupply:            shareSupply,
		ShareAmt:               shareAmt,
		IsFirstMint:            isFirstMint,
		RecordType:             recordType,
	}, nil
}

func createShareRecord(tx *gorm.DB, shareId uint, username string, _liquidity *big.Int, token0TransferRecordId uint, token1TransferRecordId uint, reserve0 string, reserve1 string, amount0 string, amount1 string, shareSupply string, shareAmt string, isFirstMint bool, recordType ShareRecordType) (err error) {
	if shareId <= 0 {
		return errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	var shareRecord *PoolShareRecord
	shareRecord, err = newShareRecord(shareId, username, _liquidity.String(), token0TransferRecordId, token1TransferRecordId, reserve0, reserve1, amount0, amount1, shareSupply, shareAmt, isFirstMint, recordType)
	if err != nil {
		return utils.AppendErrorInfo(err, "newShareRecord")
	}
	err = tx.Model(&PoolShareRecord{}).Create(&shareRecord).Error
	if err != nil {
		return utils.AppendErrorInfo(err, "create shareRecord")
	}
	return nil
}

func updateShareBalanceAndRecordMint(tx *gorm.DB, shareId uint, username string, _liquidity *big.Int, token0TransferRecordId uint, token1TransferRecordId uint, reserve0 string, reserve1 string, amount0 string, amount1 string, shareSupply string, isFirstMint bool) (err error) {
	if shareId <= 0 {
		return errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	var previousShare string
	previousShare, err = createOrUpdateShareBalance(tx, shareId, username, _liquidity)
	if err != nil {
		return utils.AppendErrorInfo(err, "createOrUpdateShareBalance")
	}
	err = createShareRecord(tx, shareId, username, _liquidity, token0TransferRecordId, token1TransferRecordId, reserve0, reserve1, amount0, amount1, shareSupply, previousShare, isFirstMint, AddLiquidityShareMint)
	if err != nil {
		return utils.AppendErrorInfo(err, "createShareRecord")
	}
	return nil
}

func updateShareBalanceAndRecordBurn(tx *gorm.DB, shareId uint, username string, _liquidity *big.Int, token0TransferRecordId uint, token1TransferRecordId uint, reserve0 string, reserve1 string, amount0 string, amount1 string, shareSupply string) (err error) {
	if shareId <= 0 {
		return errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	var previousShare string
	previousShare, err = updateShareBalanceBurn(tx, shareId, username, _liquidity)
	if err != nil {
		return utils.AppendErrorInfo(err, "updateShareBalanceBurn")
	}
	err = createShareRecord(tx, shareId, username, _liquidity, token0TransferRecordId, token1TransferRecordId, reserve0, reserve1, amount0, amount1, shareSupply, previousShare, false, RemoveLiquidityShareBurn)
	if err != nil {
		return utils.AppendErrorInfo(err, "createShareRecord")
	}
	return nil
}

// calc

func calcNewShareRecord(shareId uint, username string, liquidity string, reserve0 string, reserve1 string, amount0 string, amount1 string, shareSupply string, shareAmt string, isFirstMint bool, recordType ShareRecordType) (shareRecord *PoolShareRecord, err error) {
	return &PoolShareRecord{
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

func calcShareRecord(shareId uint, username string, _liquidity *big.Int, reserve0 string, reserve1 string, amount0 string, amount1 string, shareSupply string, shareAmt string, isFirstMint bool, recordType ShareRecordType) (shareRecord *PoolShareRecord, err error) {
	shareRecord, err = calcNewShareRecord(shareId, username, _liquidity.String(), reserve0, reserve1, amount0, amount1, shareSupply, shareAmt, isFirstMint, recordType)
	if err != nil {
		return new(PoolShareRecord), utils.AppendErrorInfo(err, "calcNewShareRecord")
	}
	return shareRecord, nil
}

func calcCreateOrUpdateShareBalance(tx *gorm.DB, shareId uint, username string, _liquidity *big.Int) (previousShare string, err error) {
	if shareId <= 0 {
		return ZeroValue, errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	var shareBalance *PoolShareBalance
	err = tx.Model(&PoolShareBalance{}).Where("share_id = ? AND username = ?", shareId, username).First(&shareBalance).Error
	if err != nil {
		previousShare = big.NewInt(0).String()
	} else {
		oldBalance, success := new(big.Int).SetString(shareBalance.Balance, 10)
		if !success {
			return ZeroValue, errors.New("shareBalance SetString(" + shareBalance.Balance + ") " + strconv.FormatBool(success))
		}
		previousShare = oldBalance.String()
	}
	return previousShare, nil
}

func calcUpdateShareBalanceBurn(tx *gorm.DB, shareId uint, username string, _liquidity *big.Int) (previousShare string, err error) {
	if shareId <= 0 {
		return ZeroValue, errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	var shareBalance *PoolShareBalance
	err = tx.Model(&PoolShareBalance{}).Where("share_id = ? AND username = ?", shareId, username).First(&shareBalance).Error
	if err != nil {
		// @dev: no shareBalance
		return ZeroValue, errors.New("no shareBalance")
	}

	oldBalance, success := new(big.Int).SetString(shareBalance.Balance, 10)
	if !success {
		return ZeroValue, errors.New("shareBalance SetString(" + shareBalance.Balance + ") " + strconv.FormatBool(success))
	}

	previousShare = oldBalance.String()
	return previousShare, nil
}

func calcUpdateShareBalanceAndRecordMint(tx *gorm.DB, shareId uint, username string, _liquidity *big.Int, reserve0 string, reserve1 string, amount0 string, amount1 string, shareSupply string, isFirstMint bool) (shareRecord *PoolShareRecord, err error) {
	var previousShare string
	previousShare, err = calcCreateOrUpdateShareBalance(tx, shareId, username, _liquidity)
	if err != nil {
		return new(PoolShareRecord), utils.AppendErrorInfo(err, "createOrUpdateShareBalance")
	}
	shareRecord, err = calcShareRecord(shareId, username, _liquidity, reserve0, reserve1, amount0, amount1, shareSupply, previousShare, isFirstMint, AddLiquidityShareMint)
	if err != nil {
		return new(PoolShareRecord), utils.AppendErrorInfo(err, "createShareRecord")
	}
	return shareRecord, nil
}

func calcUpdateShareBalanceAndRecordBurn(tx *gorm.DB, shareId uint, username string, _liquidity *big.Int, reserve0 string, reserve1 string, amount0 string, amount1 string, shareSupply string) (shareRecord *PoolShareRecord, err error) {
	if shareId <= 0 {
		return new(PoolShareRecord), errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}
	var previousShare string
	previousShare, err = calcUpdateShareBalanceBurn(tx, shareId, username, _liquidity)
	if err != nil {
		return new(PoolShareRecord), utils.AppendErrorInfo(err, "updateShareBalanceBurn")
	}
	shareRecord, err = calcShareRecord(shareId, username, _liquidity, reserve0, reserve1, amount0, amount1, shareSupply, previousShare, false, RemoveLiquidityShareBurn)
	if err != nil {
		return new(PoolShareRecord), utils.AppendErrorInfo(err, "createShareRecord")
	}
	return shareRecord, nil
}
