package pool

import (
	"errors"
	"gorm.io/gorm"
	"math/big"
	"strconv"
	"sync"
	"trade/middleware"
	"trade/utils"
)

// TODO: Pool
type Pool struct {
	gorm.Model
	PairId uint `json:"pair_id" gorm:"uniqueIndex"`
	// TODO: Account
	AccountId uint `json:"account_id" gorm:"index"`
}

var Lock map[string]map[string]*sync.Mutex

// AddLiquidity
// @Description: Add Liquidity
// @param tx transaction, begin first, must be rollback or commit
func AddLiquidity(tokenA string, tokenB string, amountADesired string, amountBDesired string, amountAMin string, amountBMin string, username string) (amountA string, amountB string, liquidity string, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "sortTokens")
	}

	var amount0Desired, amount1Desired, amount0Min, amount1Min string
	if token0 == tokenB {
		amount0Desired, amount1Desired = amountBDesired, amountADesired
		amount0Min, amount1Min = amountBMin, amountAMin
	} else {
		amount0Desired, amount1Desired = amountADesired, amountBDesired
		amount0Min, amount1Min = amountAMin, amountBMin
	}

	// amount0Desired
	_amount0Desired, success := new(big.Int).SetString(amount0Desired, 10)
	if !success {
		return ZeroValue, ZeroValue, ZeroValue, errors.New("amount0Desired SetString(" + amount0Desired + ") " + strconv.FormatBool(success))
	}

	// amount1Desired
	_amount1Desired, success := new(big.Int).SetString(amount1Desired, 10)
	if !success {
		return ZeroValue, ZeroValue, ZeroValue, errors.New("amount1Desired SetString(" + amount1Desired + ") " + strconv.FormatBool(success))
	}

	// amount0Min
	_amount0Min, success := new(big.Int).SetString(amount0Min, 10)
	if !success {
		return ZeroValue, ZeroValue, ZeroValue, errors.New("amount0Min SetString(" + amount0Min + ") " + strconv.FormatBool(success))
	}
	if _amount0Min.Sign() < 0 {
		return ZeroValue, ZeroValue, ZeroValue, errors.New("amount0Min(" + _amount0Min.String() + ") is negative")
	}
	if _amount0Min.Cmp(_amount0Desired) > 0 {
		return ZeroValue, ZeroValue, ZeroValue, errors.New("amount0Min(" + _amount0Min.String() + ") is greater than amount0Desired(" + _amount0Desired.String() + ")")
	}

	// amount1Min
	_amount1Min, success := new(big.Int).SetString(amount1Min, 10)
	if !success {
		return ZeroValue, ZeroValue, ZeroValue, errors.New("amount1Min SetString(" + amount1Min + ") " + strconv.FormatBool(success))
	}
	if _amount1Min.Sign() < 0 {
		return ZeroValue, ZeroValue, ZeroValue, errors.New("amount1Min(" + _amount1Min.String() + ") is negative")
	}
	if _amount1Min.Cmp(_amount1Desired) > 0 {
		return ZeroValue, ZeroValue, ZeroValue, errors.New("amount1Min(" + _amount1Min.String() + ") is greater than amount1Desired(" + _amount1Desired.String() + ")")
	}

	// @dev: lock
	if Lock == nil {
		Lock = make(map[string]map[string]*sync.Mutex)
	}
	if Lock[token0] == nil {
		Lock[token0] = make(map[string]*sync.Mutex)
	}
	if Lock[token0][token1] == nil {
		Lock[token0][token1] = new(sync.Mutex)
	}
	Lock[token0][token1].Lock()
	// @dev: defer finally unlock
	defer Lock[token0][token1].Unlock()

	tx := middleware.DB.Begin()
	// @dev: defer firstly commit
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	var _amount0, _amount1 = new(big.Int), new(big.Int)
	var _reserve0, _reserve1 *big.Int
	// @dev: get pair
	var _pair Pair
	var pairId uint
	err = tx.Model(&Pair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	// @dev: pair does not exist
	if err != nil {
		*_amount0, *_amount1 = *_amount0Desired, *_amount1Desired
		// @dev: create pair
		newPair, err := NewPairBig(token0, token1, _amount0, _amount1)
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "newPair")
		}
		err = tx.Model(&Pair{}).Create(&newPair).Error
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "create pair")
		}
		// @dev: set pairId
		pairId = newPair.ID
		// reserve0, reserve1
		_reserve0, _reserve1 = big.NewInt(0), big.NewInt(0)
	} else {
		// @dev: pair exists
		// reserve0
		_reserve0, success = new(big.Int).SetString(_pair.Reserve0, 10)
		if !success {
			return ZeroValue, ZeroValue, ZeroValue, errors.New("Reserve0 SetString(" + _pair.Reserve0 + ") " + strconv.FormatBool(success))
		}
		// reserve1
		_reserve1, success = new(big.Int).SetString(_pair.Reserve1, 10)
		if !success {
			return ZeroValue, ZeroValue, ZeroValue, errors.New("Reserve1 SetString(" + _pair.Reserve1 + ") " + strconv.FormatBool(success))
		}
		// No fee for adding liquidity
		_amount0, _amount1, err = _addLiquidity(_amount0Desired, _amount1Desired, _amount0Min, _amount1Min, _reserve0, _reserve1)
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "_addLiquidity")
		}
		// update pair
		_newReserve0 := new(big.Int).Add(_reserve0, _amount0)
		if _newReserve0.Cmp(_reserve0) < 0 {
			return ZeroValue, ZeroValue, ZeroValue, errors.New("invalid _newReserve0(" + _newReserve0.String() + ")")
		}
		_newReserve1 := new(big.Int).Add(_reserve1, _amount1)
		if _newReserve1.Cmp(_reserve1) < 0 {
			return ZeroValue, ZeroValue, ZeroValue, errors.New("invalid _newReserve1(" + _newReserve1.String() + ")")
		}
		err = tx.Model(&Pair{}).Where("token0 = ? AND token1 = ?", token0, token1).
			Updates(map[string]any{
				"reserve0": _newReserve0.String(),
				"reserve1": _newReserve1.String(),
			}).Error
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "update pair")
		}
		// @dev: update pairId
		pairId = _pair.ID
	}

	// TODO: check balance then transfer, record transfer Id
	//TransferHelper.safeTransferFrom(tokenA, msg.sender, pair, amountA);
	// TODO: Transfer _amount0 of token0 from user to pool
	//TransferHelper.safeTransferFrom(tokenB, msg.sender, pair, amountB);
	// TODO: Transfer _amount1 of token1 from user to pool

	// get share
	var share Share
	var shareId uint
	var _liquidity *big.Int
	err = tx.Model(&Share{}).Where("pair_id = ?", pairId).First(&share).Error
	if err != nil {
		// @dev: no share
		_liquidity, err = _mintBig(_amount0, _amount1, big.NewInt(0), big.NewInt(0), big.NewInt(0))
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "_mintBig")
		}
		var newShare *Share
		newShare, err = NewShare(pairId, _liquidity.String())
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "NewShare")
		}
		// @dev: update share, shareBalance, shareRecord
		err = tx.Model(&Share{}).Create(&newShare).Error
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "create share")
		}
		shareId = newShare.ID

		shareSupply := big.NewInt(0).String()
		err = UpdateShareBalanceAndRecordMint(tx, shareId, username, _liquidity, _reserve0.String(), _reserve1.String(), _amount0.String(), _amount1.String(), shareSupply, true)
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "UpdateShareBalanceAndRecordMint")
		}
		// record liquidity
		liquidity = _liquidity.String()
	} else {
		shareId = share.ID

		_totalSupply, success := new(big.Int).SetString(share.TotalSupply, 10)
		if !success {
			return ZeroValue, ZeroValue, ZeroValue, errors.New("TotalSupply SetString(" + share.TotalSupply + ") " + strconv.FormatBool(success))
		}
		_liquidity, err = _mintBig(_amount0, _amount1, _reserve0, _reserve1, _totalSupply)
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "_mintBig")
		}
		// @dev: update share, shareBalance, shareRecord
		newSupply := new(big.Int).Add(_totalSupply, _liquidity)
		err = tx.Model(&Share{}).Where("pair_id = ?", pairId).
			Update("total_supply", newSupply.String()).Error
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "update share")
		}
		err = UpdateShareBalanceAndRecordMint(tx, shareId, username, _liquidity, _reserve0.String(), _reserve1.String(), _amount0.String(), _amount1.String(), _totalSupply.String(), false)
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "UpdateShareBalanceAndRecordMint")
		}
		// record liquidity
		liquidity = newSupply.String()
	}
	if token0 == tokenB {
		amountA, amountB = _amount1.String(), _amount0.String()
	} else {
		amountA, amountB = _amount0.String(), _amount1.String()
	}
	err = nil
	return amountA, amountB, liquidity, err
}

func RemoveLiquidity(tokenA string, tokenB string, liquidity string, amountAMin string, amountBMin string, username string, feeK uint16) (amountA string, amountB string, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "sortTokens")
	}

	var amount0Min, amount1Min string
	if token0 == tokenB {
		amount0Min, amount1Min = amountBMin, amountAMin
	} else {
		amount0Min, amount1Min = amountAMin, amountBMin
	}

	// amount0Min
	_amount0Min, success := new(big.Int).SetString(amount0Min, 10)
	if !success {
		return ZeroValue, ZeroValue, errors.New("amount0Min SetString(" + amount0Min + ") " + strconv.FormatBool(success))
	}
	if _amount0Min.Sign() < 0 {
		return ZeroValue, ZeroValue, errors.New("amount0Min(" + _amount0Min.String() + ") is negative")
	}

	// amount1Min
	_amount1Min, success := new(big.Int).SetString(amount1Min, 10)
	if !success {
		return ZeroValue, ZeroValue, errors.New("amount1Min SetString(" + amount1Min + ") " + strconv.FormatBool(success))
	}
	if _amount1Min.Sign() < 0 {
		return ZeroValue, ZeroValue, errors.New("amount1Min(" + _amount1Min.String() + ") is negative")
	}

	// @dev: lock
	if Lock == nil {
		Lock = make(map[string]map[string]*sync.Mutex)
	}
	if Lock[token0] == nil {
		Lock[token0] = make(map[string]*sync.Mutex)
	}
	if Lock[token0][token1] == nil {
		Lock[token0][token1] = new(sync.Mutex)
	}
	Lock[token0][token1].Lock()
	// @dev: defer finally unlock
	defer Lock[token0][token1].Unlock()

	tx := middleware.DB.Begin()
	// @dev: defer firstly commit
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	// @dev: get pair
	var _pair Pair
	err = tx.Model(&Pair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	if err != nil {
		//@dev: pair does not exist
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "pair does not exist")
	}
	pairId := _pair.ID
	if pairId <= 0 {
		return ZeroValue, ZeroValue, errors.New("invalid pairId(" + strconv.FormatUint(uint64(pairId), 10) + ")")
	}

	var _reserve0, _reserve1 *big.Int
	// reserve0
	_reserve0, success = new(big.Int).SetString(_pair.Reserve0, 10)
	if !success {
		return ZeroValue, ZeroValue, errors.New("Reserve0 SetString(" + _pair.Reserve0 + ") " + strconv.FormatBool(success))
	}
	// reserve1
	_reserve1, success = new(big.Int).SetString(_pair.Reserve1, 10)
	if !success {
		return ZeroValue, ZeroValue, errors.New("Reserve1 SetString(" + _pair.Reserve1 + ") " + strconv.FormatBool(success))
	}
	// liquidity
	_liquidity, success := new(big.Int).SetString(liquidity, 10)
	if !success {
		return ZeroValue, ZeroValue, errors.New("liquidity SetString(" + liquidity + ") " + strconv.FormatBool(success))
	}

	// get share
	var share Share
	err = tx.Model(&Share{}).Where("pair_id = ?", pairId).First(&share).Error
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "share does not exist")
	}
	shareId := share.ID
	if shareId <= 0 {
		return ZeroValue, ZeroValue, errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}

	_totalSupply, success := new(big.Int).SetString(share.TotalSupply, 10)
	if !success {
		return ZeroValue, ZeroValue, errors.New("TotalSupply SetString(" + share.TotalSupply + ") " + strconv.FormatBool(success))
	}

	_amount0, _amount1, err := _burnBig(_reserve0, _reserve1, _totalSupply, _liquidity, feeK)
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "_burnBig")
	}

	if !(_amount0.Cmp(_amount0Min) >= 0) {
		return ZeroValue, ZeroValue, errors.New("insufficientAmount0(" + _amount0.String() + "), need amount0Min(" + _amount0Min.String() + ")")
	}

	if !(_amount1.Cmp(_amount1Min) >= 0) {
		return ZeroValue, ZeroValue, errors.New("insufficientAmount1(" + _amount1.String() + "), need amount1Min(" + _amount1Min.String() + ")")
	}

	//        _safeTransfer(_token0, to, amount0);
	// TODO: transfer _amount0 of token0 from pool to user
	//        _safeTransfer(_token1, to, amount1);
	// TODO: transfer _amount1 of token1 from pool to user

	// @dev: update pair, share, shareBalance, shareRecord

	// update pair
	_newReserve0 := new(big.Int).Sub(_reserve0, _amount0)
	if _newReserve0.Sign() <= 0 {
		return ZeroValue, ZeroValue, errors.New("invalid _newReserve0(" + _newReserve0.String() + ")")
	}
	_newReserve1 := new(big.Int).Sub(_reserve1, _amount1)
	if _newReserve1.Sign() <= 0 {
		return ZeroValue, ZeroValue, errors.New("invalid _newReserve1(" + _newReserve1.String() + ")")
	}

	err = tx.Model(&Pair{}).Where("token0 = ? AND token1 = ?", token0, token1).
		Updates(map[string]any{
			"reserve0": _newReserve0.String(),
			"reserve1": _newReserve1.String(),
		}).Error
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "update pair")
	}

	// update share
	newSupply := new(big.Int).Sub(_totalSupply, _liquidity)
	err = tx.Model(&Share{}).Where("pair_id = ?", pairId).
		Update("total_supply", newSupply.String()).Error
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "update share")
	}

	// update shareBalance and shareRecord
	err = UpdateShareBalanceAndRecordBurn(tx, shareId, username, _liquidity, _reserve0.String(), _reserve1.String(), _amount0.String(), _amount1.String(), _totalSupply.String())
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "UpdateShareBalanceAndRecordBurn")
	}

	if token0 == tokenB {
		amountA, amountB = _amount1.String(), _amount0.String()
	} else {
		amountA, amountB = _amount0.String(), _amount1.String()
	}

	err = nil
	return amountA, amountB, err
}

// TODO: Swap Exact Tokens For Tokens
func swapExactTokensForTokens() {

}

// TODO: Swap Tokens For Exact Tokens
func swapTokensForExactTokens() {

}
