package pool

import (
	"errors"
	"gorm.io/gorm"
	"math/big"
	"strconv"
	"sync"
	"time"
	"trade/middleware"
	"trade/utils"
)

var LockP map[string]map[string]*sync.Mutex

// addLiquidity
// @Description: Add Liquidity, create pair and share if not exist
func addLiquidity(tokenA string, tokenB string, amountADesired string, amountBDesired string, amountAMin string, amountBMin string, username string) (amountA string, amountB string, liquidity string, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "sortTokens")
	}

	isTokenZeroSat := token0 == TokenSatTag

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
	if LockP == nil {
		LockP = make(map[string]map[string]*sync.Mutex)
	}
	if LockP[token0] == nil {
		LockP[token0] = make(map[string]*sync.Mutex)
	}
	if LockP[token0][token1] == nil {
		LockP[token0][token1] = new(sync.Mutex)
	}
	LockP[token0][token1].Lock()
	// @dev: defer finally unlock
	defer LockP[token0][token1].Unlock()

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
	var _pair PoolPair
	var pairId uint
	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	// @dev: pair does not exist
	if err != nil {
		*_amount0, *_amount1 = *_amount0Desired, *_amount1Desired
		// @dev: create pair
		newPair, err := newPairBig(token0, token1, _amount0, _amount1)
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "newPairBig")
		}
		err = tx.Model(&PoolPair{}).Create(&newPair).Error
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "create pair")
		}
		// @dev: set pairId
		pairId = newPair.ID
		// reserve0, reserve1
		_reserve0, _reserve1 = big.NewInt(0), big.NewInt(0)

		// @dev: create account
		err = CreatePoolAccount(tx, pairId, []string{token0, token1})
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "CreatePoolAccount("+strconv.FormatUint(uint64(pairId), 10)+")")
		}

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
		err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).
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

	var token0TransferRecordId, token1TransferRecordId uint

	token0TransferRecordId, err = TransferToPoolAccount(tx, username, pairId, token0, _amount0, "addLiquidity")
	if err != nil {
		return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "TransferToPoolAccount")
	}

	token1TransferRecordId, err = TransferToPoolAccount(tx, username, pairId, token1, _amount1, "addLiquidity")
	if err != nil {
		return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "TransferToPoolAccount")
	}

	// get share
	var share PoolShare
	var shareId uint
	var _liquidity *big.Int
	err = tx.Model(&PoolShare{}).Where("pair_id = ?", pairId).First(&share).Error
	if err != nil {
		// @dev: no share
		_liquidity, err = _mintBig(_amount0, _amount1, big.NewInt(0), big.NewInt(0), big.NewInt(0), isTokenZeroSat)
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "_mintBig")
		}
		var newShare *PoolShare
		newShare, err = _newShare(pairId, _liquidity.String())
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "_newShare")
		}
		// @dev: update share, shareBalance, shareRecord
		err = tx.Model(&PoolShare{}).Create(&newShare).Error
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "create share")
		}
		shareId = newShare.ID

		shareSupply := big.NewInt(0).String()
		// @dev: save recordIds
		err = updateShareBalanceAndRecordMint(tx, shareId, username, _liquidity, token0TransferRecordId, token1TransferRecordId, _reserve0.String(), _reserve1.String(), _amount0.String(), _amount1.String(), shareSupply, true)
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "updateShareBalanceAndRecordMint")
		}
		// record liquidity
		liquidity = _liquidity.String()
	} else {
		shareId = share.ID

		_totalSupply, success := new(big.Int).SetString(share.TotalSupply, 10)
		if !success {
			return ZeroValue, ZeroValue, ZeroValue, errors.New("TotalSupply SetString(" + share.TotalSupply + ") " + strconv.FormatBool(success))
		}
		_liquidity, err = _mintBig(_amount0, _amount1, _reserve0, _reserve1, _totalSupply, isTokenZeroSat)
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "_mintBig")
		}
		// @dev: update share, shareBalance, shareRecord
		newSupply := new(big.Int).Add(_totalSupply, _liquidity)
		//fmt.Printf("newSupply: %v;(_totalSupply: %v + _liquidity: %v)\n", newSupply, _totalSupply, _liquidity)
		err = tx.Model(&PoolShare{}).Where("pair_id = ?", pairId).
			Update("total_supply", newSupply.String()).Error
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "update share")
		}

		err = updateShareBalanceAndRecordMint(tx, shareId, username, _liquidity, token0TransferRecordId, token1TransferRecordId, _reserve0.String(), _reserve1.String(), _amount0.String(), _amount1.String(), _totalSupply.String(), false)
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "updateShareBalanceAndRecordMint")
		}
		// record liquidity
		liquidity = _liquidity.String()
	}
	if token0 == tokenB {
		amountA, amountB = _amount1.String(), _amount0.String()
	} else {
		amountA, amountB = _amount0.String(), _amount1.String()
	}
	err = nil
	return amountA, amountB, liquidity, err
}

// removeLiquidity
// @Description: Remove Liquidity, pair and share must exist
func removeLiquidity(tokenA string, tokenB string, liquidity string, amountAMin string, amountBMin string, username string, feeK uint16) (amountA string, amountB string, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "sortTokens")
	}

	isTokenZeroSat := token0 == TokenSatTag

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

	if isTokenZeroSat {
		_minRemoveLiquiditySat := new(big.Int).SetUint64(uint64(MinRemoveLiquiditySat))
		if _amount0Min.Cmp(_minRemoveLiquiditySat) < 0 {
			return ZeroValue, ZeroValue, errors.New("insufficient _amount0Min(" + _amount0Min.String() + "), need " + _minRemoveLiquiditySat.String())
		}
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
	if LockP == nil {
		LockP = make(map[string]map[string]*sync.Mutex)
	}
	if LockP[token0] == nil {
		LockP[token0] = make(map[string]*sync.Mutex)
	}
	if LockP[token0][token1] == nil {
		LockP[token0][token1] = new(sync.Mutex)
	}
	LockP[token0][token1].Lock()
	// @dev: defer finally unlock
	defer LockP[token0][token1].Unlock()

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
	var _pair PoolPair
	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
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
	var share PoolShare
	err = tx.Model(&PoolShare{}).Where("pair_id = ?", pairId).First(&share).Error
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "share does not exist")
	}
	shareId := share.ID
	if shareId <= 0 {
		return ZeroValue, ZeroValue, errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}

	// @dev: check liquidity
	var shareBalance string
	err = tx.Model(&PoolShareBalance{}).Select("balance").
		Where("share_id = ? AND username = ?", shareId, username).
		Scan(&shareBalance).Error
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "get shareBalance")
	}
	_shareBalance, success := new(big.Int).SetString(shareBalance, 10)
	if !success {
		return ZeroValue, ZeroValue, errors.New("shareBalance SetString(" + shareBalance + ") " + strconv.FormatBool(success))
	}
	if _shareBalance.Cmp(_liquidity) < 0 {
		return ZeroValue, ZeroValue, errors.New("insufficientShareBalance(" + _shareBalance.String() + ")")
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

	var token0TransferRecordId, token1TransferRecordId uint

	token0TransferRecordId, err = PoolAccountTransfer(tx, pairId, username, token0, _amount0, "removeLiquidity")
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "PoolAccountTransfer")
	}

	token1TransferRecordId, err = PoolAccountTransfer(tx, pairId, username, token1, _amount1, "removeLiquidity")
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "PoolAccountTransfer")
	}

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

	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).
		Updates(map[string]any{
			"reserve0": _newReserve0.String(),
			"reserve1": _newReserve1.String(),
		}).Error
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "update pair")
	}

	// update share
	newSupply := new(big.Int).Sub(_totalSupply, _liquidity)
	err = tx.Model(&PoolShare{}).Where("pair_id = ?", pairId).
		Update("total_supply", newSupply.String()).Error
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "update share")
	}

	// update shareBalance and shareRecord
	err = updateShareBalanceAndRecordBurn(tx, shareId, username, _liquidity, token0TransferRecordId, token1TransferRecordId, _reserve0.String(), _reserve1.String(), _amount0.String(), _amount1.String(), _totalSupply.String())
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "updateShareBalanceAndRecordBurn")
	}

	if token0 == tokenB {
		amountA, amountB = _amount1.String(), _amount0.String()
	} else {
		amountA, amountB = _amount0.String(), _amount1.String()
	}

	err = nil
	return amountA, amountB, err
}

func swapExactTokenForTokenNoPath(tokenIn string, tokenOut string, amountIn string, amountOutMin string, username string, projectPartyFeeK uint16, lpAwardFeeK uint16) (amountOut string, err error) {
	feeK := projectPartyFeeK + lpAwardFeeK

	token0, token1, err := sortTokens(tokenIn, tokenOut)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "sortTokens")
	}

	isTokenZeroSat := token0 == TokenSatTag

	// amountIn
	_amountIn, success := new(big.Int).SetString(amountIn, 10)
	if !success {
		return ZeroValue, errors.New("amountIn SetString(" + amountIn + ") " + strconv.FormatBool(success))
	}
	// amountOutMin
	_amountOutMin, success := new(big.Int).SetString(amountOutMin, 10)
	if !success {
		return ZeroValue, errors.New("amountOutMin SetString(" + amountOutMin + ") " + strconv.FormatBool(success))
	}

	if isTokenZeroSat {
		_minSwapSat := new(big.Int).SetUint64(uint64(MinSwapSatFee))
		if tokenIn == TokenSatTag {
			if _amountIn.Cmp(_minSwapSat) <= 0 {
				return ZeroValue, errors.New("insufficient _amountIn(" + _amountIn.String() + "), need " + _minSwapSat.String())
			}
		} else if tokenOut == TokenSatTag {
			if _amountOutMin.Cmp(_minSwapSat) <= 0 {
				return ZeroValue, errors.New("insufficient _amountOutMin(" + _amountOutMin.String() + "), need gt " + _minSwapSat.String())
			}
		}
	}

	// @dev: lock
	if LockP == nil {
		LockP = make(map[string]map[string]*sync.Mutex)
	}
	if LockP[token0] == nil {
		LockP[token0] = make(map[string]*sync.Mutex)
	}
	if LockP[token0][token1] == nil {
		LockP[token0][token1] = new(sync.Mutex)
	}
	LockP[token0][token1].Lock()
	// @dev: defer finally unlock
	defer LockP[token0][token1].Unlock()

	tx := middleware.DB.Begin()
	// @dev: defer firstly commit
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	// get pair
	var _pair PoolPair
	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	if err != nil {
		//@dev: pair does not exist
		return ZeroValue, utils.AppendErrorInfo(err, "pair does not exist")
	}
	pairId := _pair.ID

	var _reserve0, _reserve1 *big.Int
	// reserve0
	_reserve0, success = new(big.Int).SetString(_pair.Reserve0, 10)
	if !success {
		return ZeroValue, errors.New("Reserve0 SetString(" + _pair.Reserve0 + ") " + strconv.FormatBool(success))
	}
	// reserve1
	_reserve1, success = new(big.Int).SetString(_pair.Reserve1, 10)
	if !success {
		return ZeroValue, errors.New("Reserve1 SetString(" + _pair.Reserve1 + ") " + strconv.FormatBool(success))
	}

	var _reserveIn, _reserveOut = new(big.Int), new(big.Int)

	var _amountOut = new(big.Int)

	var swapFeeType SwapFeeType

	var _swapFee = big.NewInt(0)

	var _swapFeeFloat = big.NewFloat(0)

	if token0 == tokenOut {

		*_reserveIn, *_reserveOut = *_reserve1, *_reserve0
		if isTokenZeroSat {
			// @dev: token => ?sat

			var _amountOutWithFee, _amountOutWithoutFee *big.Int

			_amountOutWithFee, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBig")
			}
			if _amountOutWithFee.Cmp(_amountOutMin) < 0 {
				return ZeroValue, errors.New("insufficientAmountOutWithFee(" + _amountOutWithFee.String() + "), need amountOutMin(" + _amountOutMin.String() + ")")
			}

			_amountOutWithoutFee, err = getAmountOutBigWithoutFee(_amountIn, _reserveIn, _reserveOut)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBigWithoutFee")
			}
			if _amountOutWithoutFee.Cmp(_amountOutMin) < 0 {
				return ZeroValue, errors.New("insufficientAmountOutWithoutFee(" + _amountOutWithoutFee.String() + "), need amountOutMin(" + _amountOutMin.String() + ")")
			}

			_minSwapSat := new(big.Int).SetUint64(uint64(MinSwapSatFee))
			if new(big.Int).Sub(_amountOutWithoutFee, _amountOutWithFee).Cmp(_minSwapSat) < 0 {
				// @dev: fee 20 sat

				_amountOut = new(big.Int).Sub(_amountOutWithoutFee, _minSwapSat)
				swapFeeType = SwapFee20Sat
				_swapFee = _minSwapSat
				_swapFeeFloat.SetInt(_swapFee)
				//fmt.Printf("_swapFee: %v\n", _swapFee)
			} else {
				// @dev: fee _amountOutWithoutFee - _amountOutWithFee

				_amountOut = _amountOutWithFee
				swapFeeType = SwapFee6Thousands
				_swapFee = new(big.Int).Sub(_amountOutWithoutFee, _amountOutWithFee)
				_swapFeeFloat.SetInt(_swapFee)
				//fmt.Printf("_swapFee: %v\n", _swapFee)
			}
		} else {
			_amountOut, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBig")
			}

			swapFeeType = SwapFee6ThousandsNotSat
			//fmt.Printf("_swapFee: %v\n", _swapFee)

			// TODO: Swap Fee Calculation

		}
	} else {
		// @dev: sat => ?token

		*_reserveIn, *_reserveOut = *_reserve0, *_reserve1

		if isTokenZeroSat {
			_minSwapSatFee := new(big.Int).SetUint64(uint64(MinSwapSatFee))

			var _amountOutWithFee, _amountOutWithoutFee *big.Int

			_amountOutWithFee, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBig")
			}
			if _amountOutWithFee.Cmp(_amountOutMin) < 0 {
				return ZeroValue, errors.New("insufficientAmountOutWithFee(" + _amountOutWithFee.String() + "), need amountOutMin(" + _amountOutMin.String() + ")")
			}

			_amountOutWithoutFee, err = getAmountOutBigWithoutFee(_amountIn, _reserveIn, _reserveOut)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBigWithoutFee")
			}
			if _amountOutWithoutFee.Cmp(_amountOutMin) < 0 {
				return ZeroValue, errors.New("insufficientAmountOutWithoutFee(" + _amountOutWithoutFee.String() + "), need amountOutMin(" + _amountOutMin.String() + ")")
			}
			_swapFeeToken1 := new(big.Int).Sub(_amountOutWithoutFee, _amountOutWithFee)
			_swapFeeToken1Float := new(big.Float).SetInt(_swapFeeToken1)

			var _price *big.Float
			_price, err = getToken1PriceBig(_reserve0, _reserve1)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getToken1PriceBig")
			}

			_swapFeeToken1ValueFloat := new(big.Float).Mul(_swapFeeToken1Float, _price)
			_minSwapSatFeeFloat := new(big.Float).SetUint64(uint64(MinSwapSatFee))

			//fmt.Printf("_swapFeeToken1ValueFloat: %v; _minSwapSatFeeFloat:%v\n", _swapFeeToken1ValueFloat, _minSwapSatFeeFloat)

			if _swapFeeToken1ValueFloat.Cmp(_minSwapSatFeeFloat) < 0 {
				// @dev: fee 20 sat

				_amountOut, err = getAmountOutBigWithoutFee(new(big.Int).Sub(_amountIn, _minSwapSatFee), _reserveIn, _reserveOut)
				if err != nil {
					return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBigWithoutFee")
				}

				swapFeeType = SwapFee20Sat
				_swapFee = _minSwapSatFee
				_swapFeeFloat.SetInt(_swapFee)
				//fmt.Printf("_swapFee: %v\n", _swapFee)
			} else {
				// @dev: fee _swapFeeToken1ValueFloat

				_amountOut, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
				if err != nil {
					return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBig")
				}

				swapFeeType = SwapFee6Thousands

				_swapFeeFloat = _swapFeeToken1ValueFloat

				// @dev: Set _swapFee
				_swapFeeToken1ValueFloat.Int(_swapFee)

				//fmt.Printf("_swapFee: %v\n", _swapFee)
			}
		} else {
			_amountOut, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBig")
			}
			swapFeeType = SwapFee6ThousandsNotSat
			//fmt.Printf("_swapFee: %v\n", _swapFee)

			// TODO: Swap Fee Calculation

		}
	}

	if _amountOut.Cmp(_amountOutMin) < 0 {
		return ZeroValue, errors.New("insufficientAmountOut(" + _amountOut.String() + "), need amountOutMin(" + _amountOutMin.String() + ")")
	}

	var tokenInTransferRecordId, tokenOutTransferRecordId uint

	tokenInTransferRecordId, err = TransferToPoolAccount(tx, username, pairId, tokenIn, _amountIn, "swapExactTokenForTokenNoPath")
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "TransferToPoolAccount")
	}

	tokenOutTransferRecordId, err = PoolAccountTransfer(tx, pairId, username, tokenOut, _amountOut, "swapExactTokenForTokenNoPath")
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "PoolAccountTransfer")
	}

	// Update pair, swapRecord

	var _newReserve0, _newReserve1 *big.Int
	if token0 == tokenOut {
		_newReserve0 = new(big.Int).Sub(_reserve0, _amountOut)
		_newReserve1 = new(big.Int).Add(_reserve1, _amountIn)
	} else {
		_newReserve0 = new(big.Int).Add(_reserve0, _amountIn)
		_newReserve1 = new(big.Int).Sub(_reserve1, _amountOut)
	}
	if _newReserve0.Sign() <= 0 {
		return ZeroValue, errors.New("invalid _newReserve0(" + _newReserve0.String() + ")")
	}
	if _newReserve1.Sign() <= 0 {
		return ZeroValue, errors.New("invalid _newReserve1(" + _newReserve1.String() + ")")
	}

	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).
		Updates(map[string]any{
			"reserve0": _newReserve0.String(),
			"reserve1": _newReserve1.String(),
		}).Error
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "update pair")
	}

	// @dev: update swapRecord
	var recordId uint

	// @dev : Save recordIds
	recordId, err = createSwapRecord(tx, pairId, username, tokenIn, tokenOut, amountIn, _amountOut.String(), _reserveIn.String(), _reserveOut.String(), tokenInTransferRecordId, tokenOutTransferRecordId, _swapFeeFloat.String(), swapFeeType, SwapExactTokenNoPath)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "createSwapRecord")
	}

	// == COPY ==
	// get share
	var share PoolShare
	var shareId uint
	err = tx.Model(&PoolShare{}).Where("pair_id = ?", pairId).First(&share).Error
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "share does not exist")
	}
	shareId = share.ID
	totalSupply := share.TotalSupply

	_totalSupplyFloat, success := new(big.Float).SetString(totalSupply)
	if !success {
		return ZeroValue, errors.New("TotalSupply SetString(" + totalSupply + ") " + strconv.FormatBool(success))
	}

	type userAndShare struct {
		Username string `json:"username"`
		Balance  string `json:"balance"`
	}

	var userAndShares []userAndShare

	err = tx.Model(&PoolShareBalance{}).
		Select("username, balance").
		Where("share_id = ?", shareId).
		Scan(&userAndShares).Error

	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "get userAndShares")
	}

	for _, _userAndShare := range userAndShares {

		// @dev: lock user's lp award balance
		if LockLpA == nil {
			LockLpA = make(map[string]*sync.Mutex)
		}
		if LockLpA[_userAndShare.Username] == nil {
			LockLpA[_userAndShare.Username] = new(sync.Mutex)
		}
		LockLpA[_userAndShare.Username].Lock()
		// @dev: defer unlock
		defer LockLpA[_userAndShare.Username].Unlock()

		_balanceFloat, success := new(big.Float).SetString(_userAndShare.Balance)
		if !success {
			return ZeroValue, errors.New(_userAndShare.Username + " balance SetString(" + _userAndShare.Balance + ") " + strconv.FormatBool(success))
		}

		var _awardFloat = big.NewFloat(0)

		LpAwardFeeKFloat := new(big.Float).SetUint64(uint64(lpAwardFeeK))
		SwapFeeKFloat := new(big.Float).SetUint64(uint64(feeK))
		_swapFeeForAwardFloat := new(big.Float).Quo(new(big.Float).Mul(_swapFeeFloat, LpAwardFeeKFloat), SwapFeeKFloat)

		_awardFloat = new(big.Float).Quo(new(big.Float).Mul(_swapFeeForAwardFloat, _balanceFloat), _totalSupplyFloat)
		err = updateLpAwardBalanceAndRecordSwap(tx, shareId, _userAndShare.Username, _awardFloat, _swapFeeFloat.String(), _userAndShare.Balance, _totalSupplyFloat.String(), recordId)

		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "updateLpAwardBalanceAndRecordSwap")
		}
	}
	// == COPY END ==

	amountOut = _amountOut.String()
	err = nil
	return amountOut, err
}

func swapTokenForExactTokenNoPath(tokenIn string, tokenOut string, amountOut string, amountInMax string, username string, projectPartyFeeK uint16, lpAwardFeeK uint16) (amountIn string, err error) {
	feeK := projectPartyFeeK + lpAwardFeeK

	token0, token1, err := sortTokens(tokenIn, tokenOut)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "sortTokens")
	}

	isTokenZeroSat := token0 == TokenSatTag

	// amountOut
	_amountOut, success := new(big.Int).SetString(amountOut, 10)
	if !success {
		return ZeroValue, errors.New("amountOut SetString(" + amountOut + ") " + strconv.FormatBool(success))
	}
	// amountInMax
	_amountInMax, success := new(big.Int).SetString(amountInMax, 10)
	if !success {
		return ZeroValue, errors.New("amountInMax SetString(" + amountInMax + ") " + strconv.FormatBool(success))
	}

	if isTokenZeroSat {
		_minSwapSat := new(big.Int).SetUint64(uint64(MinSwapSatFee))
		if tokenOut == TokenSatTag {
			if _amountOut.Cmp(_minSwapSat) <= 0 {
				return ZeroValue, errors.New("insufficient _amountOut(" + _amountOut.String() + "), need " + _minSwapSat.String())
			}
		} else if tokenIn == TokenSatTag {
			if _amountInMax.Cmp(_minSwapSat) <= 0 {
				return ZeroValue, errors.New("insufficient _amountInMax(" + _amountInMax.String() + "), need gt " + _minSwapSat.String())
			}
		}
	}

	// @dev: lock
	if LockP == nil {
		LockP = make(map[string]map[string]*sync.Mutex)
	}
	if LockP[token0] == nil {
		LockP[token0] = make(map[string]*sync.Mutex)
	}
	if LockP[token0][token1] == nil {
		LockP[token0][token1] = new(sync.Mutex)
	}
	LockP[token0][token1].Lock()
	// @dev: defer finally unlock
	defer LockP[token0][token1].Unlock()

	tx := middleware.DB.Begin()
	// @dev: defer firstly commit
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	// get pair
	var _pair PoolPair
	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	if err != nil {
		//@dev: pair does not exist
		return ZeroValue, utils.AppendErrorInfo(err, "pair does not exist")
	}
	pairId := _pair.ID

	var _reserve0, _reserve1 *big.Int
	// reserve0
	_reserve0, success = new(big.Int).SetString(_pair.Reserve0, 10)
	if !success {
		return ZeroValue, errors.New("Reserve0 SetString(" + _pair.Reserve0 + ") " + strconv.FormatBool(success))
	}
	// reserve1
	_reserve1, success = new(big.Int).SetString(_pair.Reserve1, 10)
	if !success {
		return ZeroValue, errors.New("Reserve1 SetString(" + _pair.Reserve1 + ") " + strconv.FormatBool(success))
	}

	var _reserveIn, _reserveOut = new(big.Int), new(big.Int)

	var _amountIn = new(big.Int)

	var swapFeeType SwapFeeType

	var _swapFee = big.NewInt(0)

	var _swapFeeFloat = big.NewFloat(0)

	if token0 == tokenOut {
		*_reserveIn, *_reserveOut = *_reserve1, *_reserve0

		if _amountOut.Cmp(_reserveOut) >= 0 {
			return ZeroValue, errors.New("excessive _amountOut(" + _amountOut.String() + "), need lt reserveOut(" + _reserveOut.String() + ")")
		}

		if isTokenZeroSat {
			// @dev: ?token => sat
			var _amountInWithFee, _amountInWithoutFee *big.Int

			_minSwapSatFee := new(big.Int).SetUint64(uint64(MinSwapSatFee))
			_amountInWithFee, err = getAmountInBig(_amountOut, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountInBig")
			}
			if _amountInWithFee.Cmp(_amountInMax) > 0 {
				return ZeroValue, errors.New("excessive _amountInWithFee(" + _amountInWithFee.String() + "), need le amountInMax(" + _amountInMax.String() + ")")
			}

			_amountInWithoutFee, err = getAmountInBigWithoutFee(_amountOut, _reserveIn, _reserveOut)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountInBigWithoutFee")
			}
			if _amountInWithoutFee.Cmp(_amountInMax) > 0 {
				return ZeroValue, errors.New("excessive _amountInWithoutFee(" + _amountInWithoutFee.String() + "), need le amountInMax(" + _amountInMax.String() + ")")
			}

			_amountInFee := new(big.Int).Sub(_amountInWithFee, _amountInWithoutFee)
			_amountInFeeFloat := new(big.Float).SetInt(_amountInFee)

			var _price *big.Float

			_price, err = getToken1PriceBig(_reserve0, _reserve1)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getToken1PriceBig")
			}
			_minSwapSatFeeFloat := new(big.Float).SetUint64(uint64(MinSwapSatFee))
			_feeValueFloat := new(big.Float).Mul(_amountInFeeFloat, _price)

			if _feeValueFloat.Cmp(_minSwapSatFeeFloat) < 0 {
				// @dev: fee 20 sat
				_amountIn, err = getAmountInBigWithoutFee(new(big.Int).Add(_amountOut, _minSwapSatFee), _reserveIn, _reserveOut)
				if err != nil {
					return ZeroValue, utils.AppendErrorInfo(err, "getAmountInBigWithoutFee")
				}
				swapFeeType = SwapFee20Sat

				_swapFee = _minSwapSatFee
				_swapFeeFloat.SetInt(_swapFee)

				//fmt.Printf("_swapFee: %v\n", _swapFee)

			} else {
				// @dev: fee _feeValueFloat
				_amountIn = _amountInWithFee
				swapFeeType = SwapFee6Thousands

				_swapFeeFloat = _feeValueFloat
				// @dev: Set _swapFee
				_feeValueFloat.Int(_swapFee)

				//fmt.Printf("_swapFee: %v\n", _swapFee)
			}

		} else {
			_amountIn, err = getAmountInBig(_amountOut, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountInBig")
			}
			swapFeeType = SwapFee6ThousandsNotSat
			//fmt.Printf("_swapFee: %v\n", _swapFee)

			// TODO: Swap Fee Calculation
		}

	} else {
		*_reserveIn, *_reserveOut = *_reserve0, *_reserve1

		if _amountOut.Cmp(_reserveOut) >= 0 {
			return ZeroValue, errors.New("excessive _amountOut(" + _amountOut.String() + "), need lt reserveOut(" + _reserveOut.String() + ")")
		}

		if isTokenZeroSat {
			// @dev: ?sat => token

			_minSwapSatFee := new(big.Int).SetUint64(uint64(MinSwapSatFee))

			var _amountInWithFee, _amountInWithoutFee *big.Int

			_amountInWithFee, err = getAmountInBig(_amountOut, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountInBig")
			}
			if _amountInWithFee.Cmp(_amountInMax) > 0 {
				return ZeroValue, errors.New("excessive _amountInWithFee(" + _amountInWithFee.String() + "), need le amountInMax(" + _amountInMax.String() + ")")
			}

			_amountInWithoutFee, err = getAmountInBigWithoutFee(_amountOut, _reserveIn, _reserveOut)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountInBigWithoutFee")
			}
			if _amountInWithoutFee.Cmp(_amountInMax) > 0 {
				return ZeroValue, errors.New("excessive _amountInWithoutFee(" + _amountInWithoutFee.String() + "), need le amountInMax(" + _amountInMax.String() + ")")
			}

			if new(big.Int).Sub(_amountInWithFee, _amountInWithoutFee).Cmp(_minSwapSatFee) < 0 {
				// @dev: fee 20 sat

				_amountIn = new(big.Int).Add(_amountInWithoutFee, _minSwapSatFee)
				swapFeeType = SwapFee20Sat

				_swapFee = _minSwapSatFee
				_swapFeeFloat.SetInt(_swapFee)
				//fmt.Printf("_swapFee: %v\n", _swapFee)

			} else {
				// @dev: fee _amountInWithFee - _amountInWithoutFee

				_amountIn = _amountInWithFee
				swapFeeType = SwapFee6Thousands

				_swapFee = new(big.Int).Sub(_amountInWithFee, _amountInWithoutFee)
				_swapFeeFloat.SetInt(_swapFee)

				//fmt.Printf("_swapFee: %v\n", _swapFee)
			}

		} else {
			_amountIn, err = getAmountInBig(_amountOut, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountInBig")
			}
			swapFeeType = SwapFee6ThousandsNotSat

			// TODO: Swap Fee Calculation
			//fmt.Printf("_swapFee: %v\n", _swapFee)

		}
	}

	if _amountIn.Cmp(_amountInMax) > 0 {
		return ZeroValue, errors.New("excessiveAmountIn(" + _amountIn.String() + "), need amountInMax(" + _amountInMax.String() + ")")
	}

	if _amountIn.Sign() <= 0 {
		return ZeroValue, errors.New("invalid _amountIn(" + _amountIn.String() + ")")
	}

	var tokenInTransferRecordId, tokenOutTransferRecordId uint

	tokenInTransferRecordId, err = TransferToPoolAccount(tx, username, pairId, tokenIn, _amountIn, "swapTokenForExactTokenNoPath")
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "TransferToPoolAccount")
	}

	tokenOutTransferRecordId, err = PoolAccountTransfer(tx, pairId, username, tokenOut, _amountOut, "swapTokenForExactTokenNoPath")
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "PoolAccountTransfer")
	}

	// Update pair, swapRecord

	var _newReserve0, _newReserve1 *big.Int
	if token0 == tokenOut {
		_newReserve0 = new(big.Int).Sub(_reserve0, _amountOut)
		_newReserve1 = new(big.Int).Add(_reserve1, _amountIn)
	} else {
		_newReserve0 = new(big.Int).Add(_reserve0, _amountIn)
		_newReserve1 = new(big.Int).Sub(_reserve1, _amountOut)
	}
	if _newReserve0.Sign() <= 0 {
		return ZeroValue, errors.New("invalid _newReserve0(" + _newReserve0.String() + ")")
	}
	if _newReserve1.Sign() <= 0 {
		return ZeroValue, errors.New("invalid _newReserve1(" + _newReserve1.String() + ")")
	}

	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).
		Updates(map[string]any{
			"reserve0": _newReserve0.String(),
			"reserve1": _newReserve1.String(),
		}).Error
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "update pair")
	}

	// @dev: update swapRecord
	var recordId uint
	recordId, err = createSwapRecord(tx, pairId, username, tokenIn, tokenOut, _amountIn.String(), amountOut, _reserveIn.String(), _reserveOut.String(), tokenInTransferRecordId, tokenOutTransferRecordId, _swapFeeFloat.String(), swapFeeType, SwapForExactTokenNoPath)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "createSwapRecord")
	}

	// get share
	var share PoolShare
	var shareId uint
	err = tx.Model(&PoolShare{}).Where("pair_id = ?", pairId).First(&share).Error
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "share does not exist")
	}
	shareId = share.ID
	totalSupply := share.TotalSupply

	_totalSupplyFloat, success := new(big.Float).SetString(totalSupply)
	if !success {
		return ZeroValue, errors.New("TotalSupply SetString(" + totalSupply + ") " + strconv.FormatBool(success))
	}

	type userAndShare struct {
		Username string `json:"username"`
		Balance  string `json:"balance"`
	}

	var userAndShares []userAndShare

	err = tx.Model(&PoolShareBalance{}).
		Select("username, balance").
		Where("share_id = ?", shareId).
		Scan(&userAndShares).Error

	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "get userAndShares")
	}

	for _, _userAndShare := range userAndShares {

		// @dev: lock user's lp award balance
		if LockLpA == nil {
			LockLpA = make(map[string]*sync.Mutex)
		}
		if LockLpA[_userAndShare.Username] == nil {
			LockLpA[_userAndShare.Username] = new(sync.Mutex)
		}
		LockLpA[_userAndShare.Username].Lock()
		// @dev: defer unlock
		defer LockLpA[_userAndShare.Username].Unlock()

		_balanceFloat, success := new(big.Float).SetString(_userAndShare.Balance)
		if !success {
			return ZeroValue, errors.New(_userAndShare.Username + " balance SetString(" + _userAndShare.Balance + ") " + strconv.FormatBool(success))
		}

		var _awardFloat = big.NewFloat(0)

		LpAwardFeeKFloat := new(big.Float).SetUint64(uint64(lpAwardFeeK))
		SwapFeeKFloat := new(big.Float).SetUint64(uint64(feeK))
		_swapFeeForAwardFloat := new(big.Float).Quo(new(big.Float).Mul(_swapFeeFloat, LpAwardFeeKFloat), SwapFeeKFloat)

		_awardFloat = new(big.Float).Quo(new(big.Float).Mul(_swapFeeForAwardFloat, _balanceFloat), _totalSupplyFloat)
		err = updateLpAwardBalanceAndRecordSwap(tx, shareId, _userAndShare.Username, _awardFloat, _swapFeeFloat.String(), _userAndShare.Balance, _totalSupplyFloat.String(), recordId)

		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "updateLpAwardBalanceAndRecordSwap")
		}
	}

	amountIn = _amountIn.String()
	err = nil
	return amountIn, err
}

func withdrawAward(tokenA string, tokenB string, username string, amount string) (newBalance string, err error) {

	// @dev: lock user's lp award balance
	if LockLpA == nil {
		LockLpA = make(map[string]*sync.Mutex)
	}
	if LockLpA[username] == nil {
		LockLpA[username] = new(sync.Mutex)
	}
	LockLpA[username].Lock()
	// @dev: defer unlock
	defer LockLpA[username].Unlock()

	tx := middleware.DB.Begin()
	// @dev: defer firstly commit
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	shareId, err := QueryShareId(tx, tokenA, tokenB)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "QueryShareId")
	}

	_amount, success := new(big.Int).SetString(amount, 10)
	if !success {
		return ZeroValue, errors.New("amount SetString(" + amount + ") " + strconv.FormatBool(success))
	}

	if _amount.Sign() <= 0 {
		return ZeroValue, errors.New("invalid _amount(" + _amount.String() + ")")
	}

	var oldBalance string

	oldBalance, newBalance, err = _withdrawAward2(tx, shareId, username, _amount)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "_withdrawAward2")
	}

	// @dev: previous
	_, _, err = _withdrawAward(tx, username, _amount)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "_withdrawAward")
	}

	var withdrawTransferRecordId uint

	// @dev: Transfer _amount of tokenSat from pool to user
	withdrawTransferRecordId, err = TransferWithdrawReward(username, _amount, "withdrawAward")
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "TransferWithdrawReward")
	}

	err = createWithdrawAwardRecord(tx, username, _amount, withdrawTransferRecordId, oldBalance)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "createWithdrawAwardRecord")
	}

	err = nil
	return newBalance, err
}

// pool public sync

type AddLiquidityResult struct {
	AmountA   string `json:"amountA"`
	AmountB   string `json:"amountB"`
	Liquidity string `json:"liquidity"`
}

func AddLiquidity(request *PoolAddLiquidityRequest) (result *AddLiquidityResult, err error) {
	if request == nil {
		return new(AddLiquidityResult), errors.New("request is nil")
	}

	_, _, err = sortTokens(request.TokenA, request.TokenB)
	if err != nil {
		return new(AddLiquidityResult), utils.AppendErrorInfo(err, "sortTokens")
	}
	if request.AmountADesired == "" {
		err = errors.New("amount_a_desired is empty")
		return new(AddLiquidityResult), err
	}
	if request.AmountBDesired == "" {
		err = errors.New("amount_b_desired is empty")
		return new(AddLiquidityResult), err
	}
	if request.AmountAMin == "" {
		err = errors.New("amount_a_min is empty")
		return new(AddLiquidityResult), err
	}
	if request.AmountBMin == "" {
		err = errors.New("amount_b_min is empty")
		return new(AddLiquidityResult), err
	}
	if request.Username == "" {
		err = errors.New("username is empty")
		return new(AddLiquidityResult), err
	}

	tokenA := request.TokenA
	tokenB := request.TokenB
	amountADesired := request.AmountADesired
	amountBDesired := request.AmountBDesired
	amountAMin := request.AmountAMin
	amountBMin := request.AmountBMin
	username := request.Username

	amountA, amountB, liquidity, err := addLiquidity(tokenA, tokenB, amountADesired, amountBDesired, amountAMin, amountBMin, username)
	return &AddLiquidityResult{
		AmountA:   amountA,
		AmountB:   amountB,
		Liquidity: liquidity,
	}, err
}

type RemoveLiquidityResult struct {
	AmountA string `json:"amountA"`
	AmountB string `json:"amountB"`
}

func RemoveLiquidity(request *PoolRemoveLiquidityRequest) (result *RemoveLiquidityResult, err error) {
	if request == nil {
		return new(RemoveLiquidityResult), errors.New("request is nil")
	}

	_, _, err = sortTokens(request.TokenA, request.TokenB)
	if err != nil {
		return new(RemoveLiquidityResult), utils.AppendErrorInfo(err, "sortTokens")
	}
	if request.Liquidity == "" {
		err = errors.New("liquidity is empty")
		return new(RemoveLiquidityResult), err
	}
	if request.AmountAMin == "" {
		err = errors.New("amount_a_min is empty")
		return new(RemoveLiquidityResult), err
	}
	if request.AmountBMin == "" {
		err = errors.New("amount_b_min is empty")
		return new(RemoveLiquidityResult), err
	}
	if request.Username == "" {
		err = errors.New("username is empty")
		return new(RemoveLiquidityResult), err
	}
	if request.FeeK != RemoveLiquidityFeeK {
		err = errors.New("invalid fee_k(" + strconv.FormatUint(uint64(request.FeeK), 10) + ")")
		return new(RemoveLiquidityResult), err
	}

	tokenA := request.TokenA
	tokenB := request.TokenB
	liquidity := request.Liquidity
	amountAMin := request.AmountAMin
	amountBMin := request.AmountBMin
	username := request.Username
	feeK := request.FeeK

	amountA, amountB, err := removeLiquidity(tokenA, tokenB, liquidity, amountAMin, amountBMin, username, feeK)
	return &RemoveLiquidityResult{
		AmountA: amountA,
		AmountB: amountB,
	}, err
}

type SwapExactTokenForTokenNoPathResult struct {
	AmountOut string `json:"amountOut"`
}

func SwapExactTokenForTokenNoPath(request *PoolSwapExactTokenForTokenNoPathRequest) (result *SwapExactTokenForTokenNoPathResult, err error) {
	if request == nil {
		return new(SwapExactTokenForTokenNoPathResult), errors.New("request is nil")
	}

	_, _, err = sortTokens(request.TokenIn, request.TokenOut)
	if err != nil {
		return new(SwapExactTokenForTokenNoPathResult), utils.AppendErrorInfo(err, "sortTokens")
	}
	if request.AmountIn == "" {
		err = errors.New("amount_in is empty")
		return new(SwapExactTokenForTokenNoPathResult), err
	}
	if request.AmountOutMin == "" {
		err = errors.New("amount_out_min is empty")
		return new(SwapExactTokenForTokenNoPathResult), err
	}
	if request.Username == "" {
		err = errors.New("username is empty")
		return new(SwapExactTokenForTokenNoPathResult), err
	}
	if request.ProjectPartyFeeK != ProjectPartyFeeK {
		err = errors.New("invalid project_party_fee_k(" + strconv.FormatUint(uint64(request.ProjectPartyFeeK), 10) + ")")
		return new(SwapExactTokenForTokenNoPathResult), err
	}
	if request.LpAwardFeeK != LpAwardFeeK {
		err = errors.New("invalid lp_award_fee_k(" + strconv.FormatUint(uint64(request.LpAwardFeeK), 10) + ")")
		return new(SwapExactTokenForTokenNoPathResult), err
	}

	tokenIn := request.TokenIn
	tokenOut := request.TokenOut
	amountIn := request.AmountIn
	amountOutMin := request.AmountOutMin
	username := request.Username
	projectPartyFeeK := request.ProjectPartyFeeK
	lpAwardFeeK := request.LpAwardFeeK

	amountOut, err := swapExactTokenForTokenNoPath(tokenIn, tokenOut, amountIn, amountOutMin, username, projectPartyFeeK, lpAwardFeeK)
	return &SwapExactTokenForTokenNoPathResult{
		AmountOut: amountOut,
	}, err
}

type SwapTokenForExactTokenNoPathResult struct {
	AmountIn string `json:"amountIn"`
}

func SwapTokenForExactTokenNoPath(request *PoolSwapTokenForExactTokenNoPathRequest) (result *SwapTokenForExactTokenNoPathResult, err error) {
	if request == nil {
		return new(SwapTokenForExactTokenNoPathResult), errors.New("request is nil")
	}

	_, _, err = sortTokens(request.TokenIn, request.TokenOut)
	if err != nil {
		return new(SwapTokenForExactTokenNoPathResult), utils.AppendErrorInfo(err, "sortTokens")
	}
	if request.AmountOut == "" {
		err = errors.New("amount_out is empty")
		return new(SwapTokenForExactTokenNoPathResult), err
	}
	if request.AmountInMax == "" {
		err = errors.New("amount_in_max is empty")
		return new(SwapTokenForExactTokenNoPathResult), err
	}
	if request.Username == "" {
		err = errors.New("username is empty")
		return new(SwapTokenForExactTokenNoPathResult), err
	}
	if request.ProjectPartyFeeK != ProjectPartyFeeK {
		err = errors.New("invalid project_party_fee_k(" + strconv.FormatUint(uint64(request.ProjectPartyFeeK), 10) + ")")
		return new(SwapTokenForExactTokenNoPathResult), err
	}
	if request.LpAwardFeeK != LpAwardFeeK {
		err = errors.New("invalid lp_award_fee_k(" + strconv.FormatUint(uint64(request.LpAwardFeeK), 10) + ")")
		return new(SwapTokenForExactTokenNoPathResult), err
	}

	tokenIn := request.TokenIn
	tokenOut := request.TokenOut
	amountOut := request.AmountOut
	amountInMax := request.AmountInMax
	username := request.Username
	projectPartyFeeK := request.ProjectPartyFeeK
	lpAwardFeeK := request.LpAwardFeeK

	amountIn, err := swapTokenForExactTokenNoPath(tokenIn, tokenOut, amountOut, amountInMax, username, projectPartyFeeK, lpAwardFeeK)
	return &SwapTokenForExactTokenNoPathResult{
		AmountIn: amountIn,
	}, err
}

type WithdrawAwardResult struct {
	NewBalance string `json:"newBalance"`
}

func WithdrawAward(request *PoolWithdrawAwardRequest) (result *WithdrawAwardResult, err error) {
	if request == nil {
		return new(WithdrawAwardResult), errors.New("request is nil")
	}

	if request.Username == "" {
		err = errors.New("username is empty")
		return new(WithdrawAwardResult), err
	}
	if request.Amount == "" {
		err = errors.New("amount is empty")
		return new(WithdrawAwardResult), err
	}

	tokenA := request.TokenA
	tokenB := request.TokenB
	username := request.Username
	amount := request.Amount

	newBalance, err := withdrawAward(tokenA, tokenB, username, amount)
	return &WithdrawAwardResult{
		NewBalance: newBalance,
	}, err
}

// calc

func calcAddLiquidity(tokenA string, tokenB string, amountADesired string, amountBDesired string, amountAMin string, amountBMin string, username string) (amountA string, amountB string, liquidity string, shareRecord *PoolShareRecord, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), utils.AppendErrorInfo(err, "sortTokens")
	}

	isTokenZeroSat := token0 == TokenSatTag

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
		return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("amount0Desired SetString(" + amount0Desired + ") " + strconv.FormatBool(success))
	}

	// amount1Desired
	_amount1Desired, success := new(big.Int).SetString(amount1Desired, 10)
	if !success {
		return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("amount1Desired SetString(" + amount1Desired + ") " + strconv.FormatBool(success))
	}

	// amount0Min
	_amount0Min, success := new(big.Int).SetString(amount0Min, 10)
	if !success {
		return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("amount0Min SetString(" + amount0Min + ") " + strconv.FormatBool(success))
	}
	if _amount0Min.Sign() < 0 {
		return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("amount0Min(" + _amount0Min.String() + ") is negative")
	}
	if _amount0Min.Cmp(_amount0Desired) > 0 {
		return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("amount0Min(" + _amount0Min.String() + ") is greater than amount0Desired(" + _amount0Desired.String() + ")")
	}

	// amount1Min
	_amount1Min, success := new(big.Int).SetString(amount1Min, 10)
	if !success {
		return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("amount1Min SetString(" + amount1Min + ") " + strconv.FormatBool(success))
	}
	if _amount1Min.Sign() < 0 {
		return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("amount1Min(" + _amount1Min.String() + ") is negative")
	}
	if _amount1Min.Cmp(_amount1Desired) > 0 {
		return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("amount1Min(" + _amount1Min.String() + ") is greater than amount1Desired(" + _amount1Desired.String() + ")")
	}

	tx := middleware.DB.Begin()
	defer func() {
		tx.Rollback()
	}()

	var _amount0, _amount1 = new(big.Int), new(big.Int)
	var _reserve0, _reserve1 *big.Int
	// @dev: get pair
	var _pair PoolPair
	var pairId uint
	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	// @dev: pair does not exist
	if err != nil {
		*_amount0, *_amount1 = *_amount0Desired, *_amount1Desired

		// reserve0, reserve1
		_reserve0, _reserve1 = big.NewInt(0), big.NewInt(0)

		var _liquidity *big.Int
		_liquidity, err = _mintBig(_amount0, _amount1, big.NewInt(0), big.NewInt(0), big.NewInt(0), isTokenZeroSat)
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), utils.AppendErrorInfo(err, "_mintBig")
		}

		shareRecord, err = calcShareRecord(0, username, _liquidity, _reserve0.String(), _reserve1.String(), _amount0.String(), _amount1.String(), ZeroValue, ZeroValue, true, AddLiquidityShareMint)

		// record liquidity
		liquidity = _liquidity.String()

		return amountADesired, amountBDesired, liquidity, shareRecord, err

	} else {
		// @dev: pair exists
		// reserve0
		_reserve0, success = new(big.Int).SetString(_pair.Reserve0, 10)
		if !success {
			return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("Reserve0 SetString(" + _pair.Reserve0 + ") " + strconv.FormatBool(success))
		}
		// reserve1
		_reserve1, success = new(big.Int).SetString(_pair.Reserve1, 10)
		if !success {
			return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("Reserve1 SetString(" + _pair.Reserve1 + ") " + strconv.FormatBool(success))
		}
		// No fee for adding liquidity
		_amount0, _amount1, err = _addLiquidity(_amount0Desired, _amount1Desired, _amount0Min, _amount1Min, _reserve0, _reserve1)
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), utils.AppendErrorInfo(err, "_addLiquidity")
		}
		// update pair
		_newReserve0 := new(big.Int).Add(_reserve0, _amount0)
		if _newReserve0.Cmp(_reserve0) < 0 {
			return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("invalid _newReserve0(" + _newReserve0.String() + ")")
		}
		_newReserve1 := new(big.Int).Add(_reserve1, _amount1)
		if _newReserve1.Cmp(_reserve1) < 0 {
			return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("invalid _newReserve1(" + _newReserve1.String() + ")")
		}

		// @dev: update pairId
		pairId = _pair.ID
	}

	// get share
	var share PoolShare
	var shareId uint
	var _liquidity *big.Int
	err = tx.Model(&PoolShare{}).Where("pair_id = ?", pairId).First(&share).Error
	if err != nil {
		// @dev: no share
		return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), utils.AppendErrorInfo(err, "share does not exist")

	} else {
		shareId = share.ID

		_totalSupply, success := new(big.Int).SetString(share.TotalSupply, 10)
		if !success {
			return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("TotalSupply SetString(" + share.TotalSupply + ") " + strconv.FormatBool(success))
		}
		_liquidity, err = _mintBig(_amount0, _amount1, _reserve0, _reserve1, _totalSupply, isTokenZeroSat)
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), utils.AppendErrorInfo(err, "_mintBig")
		}
		// @dev: update share, shareBalance, shareRecord
		newSupply := new(big.Int).Add(_totalSupply, _liquidity)
		//fmt.Printf("newSupply: %v;(_totalSupply: %v + _liquidity: %v)\n", newSupply, _totalSupply, _liquidity)

		_ = newSupply

		shareRecord, err = calcUpdateShareBalanceAndRecordMint(tx, shareId, username, _liquidity, _reserve0.String(), _reserve1.String(), _amount0.String(), _amount1.String(), _totalSupply.String(), false)
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, new(PoolShareRecord), utils.AppendErrorInfo(err, "updateShareBalanceAndRecordMint")
		}
		// record liquidity
		liquidity = _liquidity.String()
	}
	if token0 == tokenB {
		amountA, amountB = _amount1.String(), _amount0.String()
	} else {
		amountA, amountB = _amount0.String(), _amount1.String()
	}
	err = nil
	return amountA, amountB, liquidity, shareRecord, err
}

func calcRemoveLiquidity(tokenA string, tokenB string, liquidity string, amountAMin string, amountBMin string, username string, feeK uint16) (amountA string, amountB string, shareRecord *PoolShareRecord, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return ZeroValue, ZeroValue, new(PoolShareRecord), utils.AppendErrorInfo(err, "sortTokens")
	}

	isTokenZeroSat := token0 == TokenSatTag

	var amount0Min, amount1Min string
	if token0 == tokenB {
		amount0Min, amount1Min = amountBMin, amountAMin
	} else {
		amount0Min, amount1Min = amountAMin, amountBMin
	}

	// amount0Min
	_amount0Min, success := new(big.Int).SetString(amount0Min, 10)
	if !success {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("amount0Min SetString(" + amount0Min + ") " + strconv.FormatBool(success))
	}
	if _amount0Min.Sign() < 0 {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("amount0Min(" + _amount0Min.String() + ") is negative")
	}

	if isTokenZeroSat {
		_minRemoveLiquiditySat := new(big.Int).SetUint64(uint64(MinRemoveLiquiditySat))
		if _amount0Min.Cmp(_minRemoveLiquiditySat) < 0 {
			return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("insufficient _amount0Min(" + _amount0Min.String() + "), need " + _minRemoveLiquiditySat.String())
		}
	}

	// amount1Min
	_amount1Min, success := new(big.Int).SetString(amount1Min, 10)
	if !success {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("amount1Min SetString(" + amount1Min + ") " + strconv.FormatBool(success))
	}
	if _amount1Min.Sign() < 0 {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("amount1Min(" + _amount1Min.String() + ") is negative")
	}

	tx := middleware.DB.Begin()
	defer func() {
		tx.Rollback()
	}()

	// @dev: get pair
	var _pair PoolPair
	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	if err != nil {
		//@dev: pair does not exist
		return ZeroValue, ZeroValue, new(PoolShareRecord), utils.AppendErrorInfo(err, "pair does not exist")
	}
	pairId := _pair.ID
	if pairId <= 0 {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("invalid pairId(" + strconv.FormatUint(uint64(pairId), 10) + ")")
	}

	var _reserve0, _reserve1 *big.Int
	// reserve0
	_reserve0, success = new(big.Int).SetString(_pair.Reserve0, 10)
	if !success {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("Reserve0 SetString(" + _pair.Reserve0 + ") " + strconv.FormatBool(success))
	}
	// reserve1
	_reserve1, success = new(big.Int).SetString(_pair.Reserve1, 10)
	if !success {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("Reserve1 SetString(" + _pair.Reserve1 + ") " + strconv.FormatBool(success))
	}
	// liquidity
	_liquidity, success := new(big.Int).SetString(liquidity, 10)
	if !success {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("liquidity SetString(" + liquidity + ") " + strconv.FormatBool(success))
	}

	// get share
	var share PoolShare
	err = tx.Model(&PoolShare{}).Where("pair_id = ?", pairId).First(&share).Error
	if err != nil {
		return ZeroValue, ZeroValue, new(PoolShareRecord), utils.AppendErrorInfo(err, "share does not exist")
	}
	shareId := share.ID
	if shareId <= 0 {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}

	// @dev: check liquidity
	var shareBalance string
	err = tx.Model(&PoolShareBalance{}).Select("balance").
		Where("share_id = ? AND username = ?", shareId, username).
		Scan(&shareBalance).Error
	if err != nil {
		return ZeroValue, ZeroValue, new(PoolShareRecord), utils.AppendErrorInfo(err, "get shareBalance")
	}
	_shareBalance, success := new(big.Int).SetString(shareBalance, 10)
	if !success {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("shareBalance SetString(" + shareBalance + ") " + strconv.FormatBool(success))
	}
	if _shareBalance.Cmp(_liquidity) < 0 {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("insufficientShareBalance(" + _shareBalance.String() + ")")
	}

	_totalSupply, success := new(big.Int).SetString(share.TotalSupply, 10)
	if !success {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("TotalSupply SetString(" + share.TotalSupply + ") " + strconv.FormatBool(success))
	}

	_amount0, _amount1, err := _burnBig(_reserve0, _reserve1, _totalSupply, _liquidity, feeK)
	if err != nil {
		return ZeroValue, ZeroValue, new(PoolShareRecord), utils.AppendErrorInfo(err, "_burnBig")
	}

	if !(_amount0.Cmp(_amount0Min) >= 0) {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("insufficientAmount0(" + _amount0.String() + "), need amount0Min(" + _amount0Min.String() + ")")
	}

	if !(_amount1.Cmp(_amount1Min) >= 0) {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("insufficientAmount1(" + _amount1.String() + "), need amount1Min(" + _amount1Min.String() + ")")
	}

	// @dev: update pair, share, shareBalance, shareRecord

	// update pair
	_newReserve0 := new(big.Int).Sub(_reserve0, _amount0)
	if _newReserve0.Sign() <= 0 {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("invalid _newReserve0(" + _newReserve0.String() + ")")
	}
	_newReserve1 := new(big.Int).Sub(_reserve1, _amount1)
	if _newReserve1.Sign() <= 0 {
		return ZeroValue, ZeroValue, new(PoolShareRecord), errors.New("invalid _newReserve1(" + _newReserve1.String() + ")")
	}

	// update share
	newSupply := new(big.Int).Sub(_totalSupply, _liquidity)

	_ = newSupply

	// update shareBalance and shareRecord
	shareRecord, err = calcUpdateShareBalanceAndRecordBurn(tx, shareId, username, _liquidity, _reserve0.String(), _reserve1.String(), _amount0.String(), _amount1.String(), _totalSupply.String())
	if err != nil {
		return ZeroValue, ZeroValue, new(PoolShareRecord), utils.AppendErrorInfo(err, "updateShareBalanceAndRecordBurn")
	}

	if token0 == tokenB {
		amountA, amountB = _amount1.String(), _amount0.String()
	} else {
		amountA, amountB = _amount0.String(), _amount1.String()
	}

	err = nil
	return amountA, amountB, shareRecord, err
}

func calcSwapExactTokenForTokenNoPath(tokenIn string, tokenOut string, amountIn string, amountOutMin string, username string, projectPartyFeeK uint16, lpAwardFeeK uint16) (amountOut string, swapRecord *PoolSwapRecord, err error) {
	feeK := projectPartyFeeK + lpAwardFeeK

	token0, token1, err := sortTokens(tokenIn, tokenOut)
	if err != nil {
		return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "sortTokens")
	}

	isTokenZeroSat := token0 == TokenSatTag

	// amountIn
	_amountIn, success := new(big.Int).SetString(amountIn, 10)
	if !success {
		return ZeroValue, new(PoolSwapRecord), errors.New("amountIn SetString(" + amountIn + ") " + strconv.FormatBool(success))
	}
	// amountOutMin
	_amountOutMin, success := new(big.Int).SetString(amountOutMin, 10)
	if !success {
		return ZeroValue, new(PoolSwapRecord), errors.New("amountOutMin SetString(" + amountOutMin + ") " + strconv.FormatBool(success))
	}

	if isTokenZeroSat {
		_minSwapSat := new(big.Int).SetUint64(uint64(MinSwapSatFee))
		if tokenIn == TokenSatTag {
			if _amountIn.Cmp(_minSwapSat) <= 0 {
				return ZeroValue, new(PoolSwapRecord), errors.New("insufficient _amountIn(" + _amountIn.String() + "), need " + _minSwapSat.String())
			}
		} else if tokenOut == TokenSatTag {
			if _amountOutMin.Cmp(_minSwapSat) <= 0 {
				return ZeroValue, new(PoolSwapRecord), errors.New("insufficient _amountOutMin(" + _amountOutMin.String() + "), need gt " + _minSwapSat.String())
			}
		}
	}

	tx := middleware.DB.Begin()
	defer func() {
		tx.Rollback()
	}()

	// get pair
	var _pair PoolPair
	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	if err != nil {
		//@dev: pair does not exist
		return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "pair does not exist")
	}
	pairId := _pair.ID

	var _reserve0, _reserve1 *big.Int
	// reserve0
	_reserve0, success = new(big.Int).SetString(_pair.Reserve0, 10)
	if !success {
		return ZeroValue, new(PoolSwapRecord), errors.New("Reserve0 SetString(" + _pair.Reserve0 + ") " + strconv.FormatBool(success))
	}
	// reserve1
	_reserve1, success = new(big.Int).SetString(_pair.Reserve1, 10)
	if !success {
		return ZeroValue, new(PoolSwapRecord), errors.New("Reserve1 SetString(" + _pair.Reserve1 + ") " + strconv.FormatBool(success))
	}

	var _reserveIn, _reserveOut = new(big.Int), new(big.Int)

	var _amountOut = new(big.Int)

	var swapFeeType SwapFeeType

	var _swapFee = big.NewInt(0)

	var _swapFeeFloat = big.NewFloat(0)

	if token0 == tokenOut {

		*_reserveIn, *_reserveOut = *_reserve1, *_reserve0
		if isTokenZeroSat {
			// @dev: token => ?sat

			var _amountOutWithFee, _amountOutWithoutFee *big.Int

			_amountOutWithFee, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getAmountOutBig")
			}
			if _amountOutWithFee.Cmp(_amountOutMin) < 0 {
				return ZeroValue, new(PoolSwapRecord), errors.New("insufficientAmountOutWithFee(" + _amountOutWithFee.String() + "), need amountOutMin(" + _amountOutMin.String() + ")")
			}

			_amountOutWithoutFee, err = getAmountOutBigWithoutFee(_amountIn, _reserveIn, _reserveOut)
			if err != nil {
				return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getAmountOutBigWithoutFee")
			}
			if _amountOutWithoutFee.Cmp(_amountOutMin) < 0 {
				return ZeroValue, new(PoolSwapRecord), errors.New("insufficientAmountOutWithoutFee(" + _amountOutWithoutFee.String() + "), need amountOutMin(" + _amountOutMin.String() + ")")
			}

			_minSwapSat := new(big.Int).SetUint64(uint64(MinSwapSatFee))
			if new(big.Int).Sub(_amountOutWithoutFee, _amountOutWithFee).Cmp(_minSwapSat) < 0 {
				// @dev: fee 20 sat

				_amountOut = new(big.Int).Sub(_amountOutWithoutFee, _minSwapSat)
				swapFeeType = SwapFee20Sat
				_swapFee = _minSwapSat
				_swapFeeFloat.SetInt(_swapFee)
				//fmt.Printf("_swapFee: %v\n", _swapFee)
			} else {
				// @dev: fee _amountOutWithoutFee - _amountOutWithFee

				_amountOut = _amountOutWithFee
				swapFeeType = SwapFee6Thousands
				_swapFee = new(big.Int).Sub(_amountOutWithoutFee, _amountOutWithFee)
				_swapFeeFloat.SetInt(_swapFee)
				//fmt.Printf("_swapFee: %v\n", _swapFee)
			}
		} else {
			_amountOut, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getAmountOutBig")
			}

			swapFeeType = SwapFee6ThousandsNotSat
			//fmt.Printf("_swapFee: %v\n", _swapFee)

		}
	} else {
		// @dev: sat => ?token

		*_reserveIn, *_reserveOut = *_reserve0, *_reserve1

		if isTokenZeroSat {
			_minSwapSatFee := new(big.Int).SetUint64(uint64(MinSwapSatFee))

			var _amountOutWithFee, _amountOutWithoutFee *big.Int

			_amountOutWithFee, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getAmountOutBig")
			}
			if _amountOutWithFee.Cmp(_amountOutMin) < 0 {
				return ZeroValue, new(PoolSwapRecord), errors.New("insufficientAmountOutWithFee(" + _amountOutWithFee.String() + "), need amountOutMin(" + _amountOutMin.String() + ")")
			}

			_amountOutWithoutFee, err = getAmountOutBigWithoutFee(_amountIn, _reserveIn, _reserveOut)
			if err != nil {
				return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getAmountOutBigWithoutFee")
			}
			if _amountOutWithoutFee.Cmp(_amountOutMin) < 0 {
				return ZeroValue, new(PoolSwapRecord), errors.New("insufficientAmountOutWithoutFee(" + _amountOutWithoutFee.String() + "), need amountOutMin(" + _amountOutMin.String() + ")")
			}
			_swapFeeToken1 := new(big.Int).Sub(_amountOutWithoutFee, _amountOutWithFee)
			_swapFeeToken1Float := new(big.Float).SetInt(_swapFeeToken1)

			var _price *big.Float
			_price, err = getToken1PriceBig(_reserve0, _reserve1)
			if err != nil {
				return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getToken1PriceBig")
			}

			_swapFeeToken1ValueFloat := new(big.Float).Mul(_swapFeeToken1Float, _price)
			_minSwapSatFeeFloat := new(big.Float).SetUint64(uint64(MinSwapSatFee))

			//fmt.Printf("_swapFeeToken1ValueFloat: %v; _minSwapSatFeeFloat:%v\n", _swapFeeToken1ValueFloat, _minSwapSatFeeFloat)

			if _swapFeeToken1ValueFloat.Cmp(_minSwapSatFeeFloat) < 0 {
				// @dev: fee 20 sat

				_amountOut, err = getAmountOutBigWithoutFee(new(big.Int).Sub(_amountIn, _minSwapSatFee), _reserveIn, _reserveOut)
				if err != nil {
					return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getAmountOutBigWithoutFee")
				}

				swapFeeType = SwapFee20Sat
				_swapFee = _minSwapSatFee
				_swapFeeFloat.SetInt(_swapFee)
				//fmt.Printf("_swapFee: %v\n", _swapFee)
			} else {
				// @dev: fee _swapFeeToken1ValueFloat

				_amountOut, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
				if err != nil {
					return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getAmountOutBig")
				}

				swapFeeType = SwapFee6Thousands

				_swapFeeFloat = _swapFeeToken1ValueFloat

				// @dev: Set _swapFee
				_swapFeeToken1ValueFloat.Int(_swapFee)

				//fmt.Printf("_swapFee: %v\n", _swapFee)
			}
		} else {
			_amountOut, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getAmountOutBig")
			}
			swapFeeType = SwapFee6ThousandsNotSat
			//fmt.Printf("_swapFee: %v\n", _swapFee)

		}
	}

	if _amountOut.Cmp(_amountOutMin) < 0 {
		return ZeroValue, new(PoolSwapRecord), errors.New("insufficientAmountOut(" + _amountOut.String() + "), need amountOutMin(" + _amountOutMin.String() + ")")
	}

	// @dev: swapRecord
	swapRecord, err = calcSwapRecord(pairId, username, tokenIn, tokenOut, amountIn, _amountOut.String(), _reserveIn.String(), _reserveOut.String(), _swapFeeFloat.String(), swapFeeType, SwapExactTokenNoPath)
	if err != nil {
		return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "calcSwapRecord")
	}

	amountOut = _amountOut.String()
	err = nil
	return amountOut, swapRecord, err
}

func calcSwapTokenForExactTokenNoPath(tokenIn string, tokenOut string, amountOut string, amountInMax string, username string, projectPartyFeeK uint16, lpAwardFeeK uint16) (amountIn string, swapRecord *PoolSwapRecord, err error) {
	feeK := projectPartyFeeK + lpAwardFeeK

	token0, token1, err := sortTokens(tokenIn, tokenOut)
	if err != nil {
		return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "sortTokens")
	}

	isTokenZeroSat := token0 == TokenSatTag

	// amountOut
	_amountOut, success := new(big.Int).SetString(amountOut, 10)
	if !success {
		return ZeroValue, new(PoolSwapRecord), errors.New("amountOut SetString(" + amountOut + ") " + strconv.FormatBool(success))
	}
	// amountInMax
	_amountInMax, success := new(big.Int).SetString(amountInMax, 10)
	if !success {
		return ZeroValue, new(PoolSwapRecord), errors.New("amountInMax SetString(" + amountInMax + ") " + strconv.FormatBool(success))
	}

	if isTokenZeroSat {
		_minSwapSat := new(big.Int).SetUint64(uint64(MinSwapSatFee))
		if tokenOut == TokenSatTag {
			if _amountOut.Cmp(_minSwapSat) <= 0 {
				return ZeroValue, new(PoolSwapRecord), errors.New("insufficient _amountOut(" + _amountOut.String() + "), need " + _minSwapSat.String())
			}
		} else if tokenIn == TokenSatTag {
			if _amountInMax.Cmp(_minSwapSat) <= 0 {
				return ZeroValue, new(PoolSwapRecord), errors.New("insufficient _amountInMax(" + _amountInMax.String() + "), need gt " + _minSwapSat.String())
			}
		}
	}

	tx := middleware.DB.Begin()
	defer func() {
		tx.Rollback()
	}()

	// get pair
	var _pair PoolPair
	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	if err != nil {
		//@dev: pair does not exist
		return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "pair does not exist")
	}
	pairId := _pair.ID

	var _reserve0, _reserve1 *big.Int
	// reserve0
	_reserve0, success = new(big.Int).SetString(_pair.Reserve0, 10)
	if !success {
		return ZeroValue, new(PoolSwapRecord), errors.New("Reserve0 SetString(" + _pair.Reserve0 + ") " + strconv.FormatBool(success))
	}
	// reserve1
	_reserve1, success = new(big.Int).SetString(_pair.Reserve1, 10)
	if !success {
		return ZeroValue, new(PoolSwapRecord), errors.New("Reserve1 SetString(" + _pair.Reserve1 + ") " + strconv.FormatBool(success))
	}

	var _reserveIn, _reserveOut = new(big.Int), new(big.Int)

	var _amountIn = new(big.Int)

	var swapFeeType SwapFeeType

	var _swapFee = big.NewInt(0)

	var _swapFeeFloat = big.NewFloat(0)

	if token0 == tokenOut {
		*_reserveIn, *_reserveOut = *_reserve1, *_reserve0

		if _amountOut.Cmp(_reserveOut) >= 0 {
			return ZeroValue, new(PoolSwapRecord), errors.New("excessive _amountOut(" + _amountOut.String() + "), need lt reserveOut(" + _reserveOut.String() + ")")
		}

		if isTokenZeroSat {
			// @dev: ?token => sat
			var _amountInWithFee, _amountInWithoutFee *big.Int

			_minSwapSatFee := new(big.Int).SetUint64(uint64(MinSwapSatFee))
			_amountInWithFee, err = getAmountInBig(_amountOut, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getAmountInBig")
			}
			if _amountInWithFee.Cmp(_amountInMax) > 0 {
				return ZeroValue, new(PoolSwapRecord), errors.New("excessive _amountInWithFee(" + _amountInWithFee.String() + "), need le amountInMax(" + _amountInMax.String() + ")")
			}

			_amountInWithoutFee, err = getAmountInBigWithoutFee(_amountOut, _reserveIn, _reserveOut)
			if err != nil {
				return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getAmountInBigWithoutFee")
			}
			if _amountInWithoutFee.Cmp(_amountInMax) > 0 {
				return ZeroValue, new(PoolSwapRecord), errors.New("excessive _amountInWithoutFee(" + _amountInWithoutFee.String() + "), need le amountInMax(" + _amountInMax.String() + ")")
			}

			_amountInFee := new(big.Int).Sub(_amountInWithFee, _amountInWithoutFee)
			_amountInFeeFloat := new(big.Float).SetInt(_amountInFee)

			var _price *big.Float

			_price, err = getToken1PriceBig(_reserve0, _reserve1)
			if err != nil {
				return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getToken1PriceBig")
			}
			_minSwapSatFeeFloat := new(big.Float).SetUint64(uint64(MinSwapSatFee))
			_feeValueFloat := new(big.Float).Mul(_amountInFeeFloat, _price)

			if _feeValueFloat.Cmp(_minSwapSatFeeFloat) < 0 {
				// @dev: fee 20 sat
				_amountIn, err = getAmountInBigWithoutFee(new(big.Int).Add(_amountOut, _minSwapSatFee), _reserveIn, _reserveOut)
				if err != nil {
					return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getAmountInBigWithoutFee")
				}
				swapFeeType = SwapFee20Sat

				_swapFee = _minSwapSatFee
				_swapFeeFloat.SetInt(_swapFee)

				//fmt.Printf("_swapFee: %v\n", _swapFee)

			} else {
				// @dev: fee _feeValueFloat
				_amountIn = _amountInWithFee
				swapFeeType = SwapFee6Thousands

				_swapFeeFloat = _feeValueFloat
				// @dev: Set _swapFee
				_feeValueFloat.Int(_swapFee)

				//fmt.Printf("_swapFee: %v\n", _swapFee)
			}

		} else {
			_amountIn, err = getAmountInBig(_amountOut, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getAmountInBig")
			}
			swapFeeType = SwapFee6ThousandsNotSat
			//fmt.Printf("_swapFee: %v\n", _swapFee)

		}

	} else {
		*_reserveIn, *_reserveOut = *_reserve0, *_reserve1

		if _amountOut.Cmp(_reserveOut) >= 0 {
			return ZeroValue, new(PoolSwapRecord), errors.New("excessive _amountOut(" + _amountOut.String() + "), need lt reserveOut(" + _reserveOut.String() + ")")
		}

		if isTokenZeroSat {
			// @dev: ?sat => token

			_minSwapSatFee := new(big.Int).SetUint64(uint64(MinSwapSatFee))

			var _amountInWithFee, _amountInWithoutFee *big.Int

			_amountInWithFee, err = getAmountInBig(_amountOut, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getAmountInBig")
			}
			if _amountInWithFee.Cmp(_amountInMax) > 0 {
				return ZeroValue, new(PoolSwapRecord), errors.New("excessive _amountInWithFee(" + _amountInWithFee.String() + "), need le amountInMax(" + _amountInMax.String() + ")")
			}

			_amountInWithoutFee, err = getAmountInBigWithoutFee(_amountOut, _reserveIn, _reserveOut)
			if err != nil {
				return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getAmountInBigWithoutFee")
			}
			if _amountInWithoutFee.Cmp(_amountInMax) > 0 {
				return ZeroValue, new(PoolSwapRecord), errors.New("excessive _amountInWithoutFee(" + _amountInWithoutFee.String() + "), need le amountInMax(" + _amountInMax.String() + ")")
			}

			if new(big.Int).Sub(_amountInWithFee, _amountInWithoutFee).Cmp(_minSwapSatFee) < 0 {
				// @dev: fee 20 sat

				_amountIn = new(big.Int).Add(_amountInWithoutFee, _minSwapSatFee)
				swapFeeType = SwapFee20Sat

				_swapFee = _minSwapSatFee
				_swapFeeFloat.SetInt(_swapFee)
				//fmt.Printf("_swapFee: %v\n", _swapFee)

			} else {
				// @dev: fee _amountInWithFee - _amountInWithoutFee

				_amountIn = _amountInWithFee
				swapFeeType = SwapFee6Thousands

				_swapFee = new(big.Int).Sub(_amountInWithFee, _amountInWithoutFee)
				_swapFeeFloat.SetInt(_swapFee)

				//fmt.Printf("_swapFee: %v\n", _swapFee)
			}

		} else {
			_amountIn, err = getAmountInBig(_amountOut, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "getAmountInBig")
			}
			swapFeeType = SwapFee6ThousandsNotSat
		}
	}

	if _amountIn.Cmp(_amountInMax) > 0 {
		return ZeroValue, new(PoolSwapRecord), errors.New("excessiveAmountIn(" + _amountIn.String() + "), need amountInMax(" + _amountInMax.String() + ")")
	}

	if _amountIn.Sign() <= 0 {
		return ZeroValue, new(PoolSwapRecord), errors.New("invalid _amountIn(" + _amountIn.String() + ")")
	}

	// @dev: swapRecord
	swapRecord, err = calcSwapRecord(pairId, username, tokenIn, tokenOut, _amountIn.String(), amountOut, _reserveIn.String(), _reserveOut.String(), _swapFeeFloat.String(), swapFeeType, SwapForExactTokenNoPath)
	if err != nil {
		return ZeroValue, new(PoolSwapRecord), utils.AppendErrorInfo(err, "calcSwapRecord")
	}

	amountIn = _amountIn.String()
	err = nil
	return amountIn, swapRecord, err
}

// public calc

func CalcQuote(tokenA string, tokenB string, amountA string) (amountB string, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "sortTokens")
	}

	// amountA
	_amountA, success := new(big.Int).SetString(amountA, 10)
	if !success {
		return ZeroValue, errors.New("amountA SetString(" + amountA + ") " + strconv.FormatBool(success))
	}
	if _amountA.Sign() < 0 {
		return ZeroValue, errors.New("amountA(" + _amountA.String() + ") is negative")
	}

	tx := middleware.DB.Begin()
	defer func() {
		tx.Rollback()
	}()

	// @dev: get pair
	var _pair PoolPair
	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	if err != nil {
		//@dev: pair does not exist
		return ZeroValue, utils.AppendErrorInfo(err, "pair does not exist")
	}
	pairId := _pair.ID
	if pairId <= 0 {
		return ZeroValue, errors.New("invalid pairId(" + strconv.FormatUint(uint64(pairId), 10) + ")")
	}

	var _reserve0, _reserve1 *big.Int
	// reserve0
	_reserve0, success = new(big.Int).SetString(_pair.Reserve0, 10)
	if !success {
		return ZeroValue, errors.New("Reserve0 SetString(" + _pair.Reserve0 + ") " + strconv.FormatBool(success))
	}
	// reserve1
	_reserve1, success = new(big.Int).SetString(_pair.Reserve1, 10)
	if !success {
		return ZeroValue, errors.New("Reserve1 SetString(" + _pair.Reserve1 + ") " + strconv.FormatBool(success))
	}

	var _reserveA, _reserveB = new(big.Int), new(big.Int)
	if token0 == tokenB {
		*_reserveA, *_reserveB = *_reserve1, *_reserve0
	} else {
		*_reserveA, *_reserveB = *_reserve0, *_reserve1
	}

	var _amountB *big.Int

	_amountB, err = quoteBig(_amountA, _reserveA, _reserveB)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "quoteBig")
	}

	return _amountB.String(), nil
}

type CalcBurnLiquidityResponse struct {
	AmountA string `json:"amount_a"`
	AmountB string `json:"amount_b"`
}

func CalcBurnLiquidity(tokenA string, tokenB string, liquidity string, username string, feeK uint16) (amountA string, amountB string, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "sortTokens")
	}

	tx := middleware.DB.Begin()
	// @dev: defer firstly commit
	defer func() {
		tx.Rollback()
	}()

	// @dev: get pair
	var _pair PoolPair
	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	if err != nil {
		//@dev: pair does not exist
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "pair does not exist")
	}
	pairId := _pair.ID
	if pairId <= 0 {
		return ZeroValue, ZeroValue, errors.New("invalid pairId(" + strconv.FormatUint(uint64(pairId), 10) + ")")
	}

	var success bool
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
	var share PoolShare
	err = tx.Model(&PoolShare{}).Where("pair_id = ?", pairId).First(&share).Error
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "share does not exist")
	}
	shareId := share.ID
	if shareId <= 0 {
		return ZeroValue, ZeroValue, errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}

	// @dev: check liquidity
	var shareBalance string
	err = tx.Model(&PoolShareBalance{}).Select("balance").
		Where("share_id = ? AND username = ?", shareId, username).
		Scan(&shareBalance).Error
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "get shareBalance")
	}
	_shareBalance, success := new(big.Int).SetString(shareBalance, 10)
	if !success {
		return ZeroValue, ZeroValue, errors.New("shareBalance SetString(" + shareBalance + ") " + strconv.FormatBool(success))
	}
	if _shareBalance.Cmp(_liquidity) < 0 {
		return ZeroValue, ZeroValue, errors.New("insufficientShareBalance(" + _shareBalance.String() + ")")
	}

	_totalSupply, success := new(big.Int).SetString(share.TotalSupply, 10)
	if !success {
		return ZeroValue, ZeroValue, errors.New("TotalSupply SetString(" + share.TotalSupply + ") " + strconv.FormatBool(success))
	}

	_amount0, _amount1, err := _burnBig(_reserve0, _reserve1, _totalSupply, _liquidity, feeK)
	if err != nil {
		return ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "_burnBig")
	}

	if token0 == tokenB {
		amountA, amountB = _amount1.String(), _amount0.String()
	} else {
		amountA, amountB = _amount0.String(), _amount1.String()
	}

	err = nil
	return amountA, amountB, err
}

func CalcAmountOut(tokenIn string, tokenOut string, amountIn string) (amountOut string, err error) {
	token0, token1, err := sortTokens(tokenIn, tokenOut)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "sortTokens")
	}

	// amountIn
	_amountIn, success := new(big.Int).SetString(amountIn, 10)
	if !success {
		return ZeroValue, errors.New("amountIn SetString(" + amountIn + ") " + strconv.FormatBool(success))
	}
	if _amountIn.Sign() < 0 {
		return ZeroValue, errors.New("amountIn(" + _amountIn.String() + ") is negative")
	}

	tx := middleware.DB.Begin()
	defer func() {
		tx.Rollback()
	}()

	// @dev: get pair
	var _pair PoolPair
	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	if err != nil {
		//@dev: pair does not exist
		return ZeroValue, utils.AppendErrorInfo(err, "pair does not exist")
	}
	pairId := _pair.ID
	if pairId <= 0 {
		return ZeroValue, errors.New("invalid pairId(" + strconv.FormatUint(uint64(pairId), 10) + ")")
	}

	var _reserve0, _reserve1 *big.Int
	// reserve0
	_reserve0, success = new(big.Int).SetString(_pair.Reserve0, 10)
	if !success {
		return ZeroValue, errors.New("Reserve0 SetString(" + _pair.Reserve0 + ") " + strconv.FormatBool(success))
	}
	// reserve1
	_reserve1, success = new(big.Int).SetString(_pair.Reserve1, 10)
	if !success {
		return ZeroValue, errors.New("Reserve1 SetString(" + _pair.Reserve1 + ") " + strconv.FormatBool(success))
	}

	var _reserveIn, _reserveOut = new(big.Int), new(big.Int)
	if token0 == tokenOut {
		*_reserveIn, *_reserveOut = *_reserve1, *_reserve0
	} else {
		*_reserveIn, *_reserveOut = *_reserve0, *_reserve1
	}

	var _amountOut *big.Int

	feeK := ProjectPartyFeeK + LpAwardFeeK

	_amountOut, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBig")
	}
	return _amountOut.String(), nil
}

func CalcAmountOut2(tokenIn string, tokenOut string, amountIn string) (amountOut string, err error) {
	feeK := ProjectPartyFeeK + LpAwardFeeK

	token0, token1, err := sortTokens(tokenIn, tokenOut)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "sortTokens")
	}

	isTokenZeroSat := token0 == TokenSatTag

	// amountIn
	_amountIn, success := new(big.Int).SetString(amountIn, 10)
	if !success {
		return ZeroValue, errors.New("amountIn SetString(" + amountIn + ") " + strconv.FormatBool(success))
	}

	if isTokenZeroSat {
		_minSwapSat := new(big.Int).SetUint64(uint64(MinSwapSatFee))
		if tokenIn == TokenSatTag {
			if _amountIn.Cmp(_minSwapSat) <= 0 {
				return ZeroValue, errors.New("insufficient _amountIn(" + _amountIn.String() + "), need " + _minSwapSat.String())
			}
		}
	}

	tx := middleware.DB.Begin()
	defer func() {
		tx.Rollback()
	}()

	// get pair
	var _pair PoolPair
	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	if err != nil {
		//@dev: pair does not exist
		return ZeroValue, utils.AppendErrorInfo(err, "pair does not exist")
	}

	var _reserve0, _reserve1 *big.Int
	// reserve0
	_reserve0, success = new(big.Int).SetString(_pair.Reserve0, 10)
	if !success {
		return ZeroValue, errors.New("Reserve0 SetString(" + _pair.Reserve0 + ") " + strconv.FormatBool(success))
	}
	// reserve1
	_reserve1, success = new(big.Int).SetString(_pair.Reserve1, 10)
	if !success {
		return ZeroValue, errors.New("Reserve1 SetString(" + _pair.Reserve1 + ") " + strconv.FormatBool(success))
	}

	var _reserveIn, _reserveOut = new(big.Int), new(big.Int)

	var _amountOut = new(big.Int)

	var _swapFee = big.NewInt(0)

	var _swapFeeFloat = big.NewFloat(0)

	if token0 == tokenOut {

		*_reserveIn, *_reserveOut = *_reserve1, *_reserve0
		if isTokenZeroSat {
			// @dev: token => ?sat

			var _amountOutWithFee, _amountOutWithoutFee *big.Int

			_amountOutWithFee, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBig")
			}

			_amountOutWithoutFee, err = getAmountOutBigWithoutFee(_amountIn, _reserveIn, _reserveOut)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBigWithoutFee")
			}

			_minSwapSat := new(big.Int).SetUint64(uint64(MinSwapSatFee))
			if new(big.Int).Sub(_amountOutWithoutFee, _amountOutWithFee).Cmp(_minSwapSat) < 0 {
				// @dev: fee 20 sat

				_amountOut = new(big.Int).Sub(_amountOutWithoutFee, _minSwapSat)
				//swapFeeType = SwapFee20Sat
				_swapFee = _minSwapSat
				_swapFeeFloat.SetInt(_swapFee)
				//fmt.Printf("_swapFee: %v\n", _swapFee)
			} else {
				// @dev: fee _amountOutWithoutFee - _amountOutWithFee

				_amountOut = _amountOutWithFee
				//swapFeeType = SwapFee6Thousands
				_swapFee = new(big.Int).Sub(_amountOutWithoutFee, _amountOutWithFee)
				_swapFeeFloat.SetInt(_swapFee)
				//fmt.Printf("_swapFee: %v\n", _swapFee)
			}
		} else {
			_amountOut, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBig")
			}

			//swapFeeType = SwapFee6ThousandsNotSat
			//fmt.Printf("_swapFee: %v\n", _swapFee)

		}
	} else {
		// @dev: sat => ?token

		*_reserveIn, *_reserveOut = *_reserve0, *_reserve1

		if isTokenZeroSat {
			_minSwapSatFee := new(big.Int).SetUint64(uint64(MinSwapSatFee))

			var _amountOutWithFee, _amountOutWithoutFee *big.Int

			_amountOutWithFee, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBig")
			}

			_amountOutWithoutFee, err = getAmountOutBigWithoutFee(_amountIn, _reserveIn, _reserveOut)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBigWithoutFee")
			}

			_swapFeeToken1 := new(big.Int).Sub(_amountOutWithoutFee, _amountOutWithFee)
			_swapFeeToken1Float := new(big.Float).SetInt(_swapFeeToken1)

			var _price *big.Float
			_price, err = getToken1PriceBig(_reserve0, _reserve1)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getToken1PriceBig")
			}

			_swapFeeToken1ValueFloat := new(big.Float).Mul(_swapFeeToken1Float, _price)
			_minSwapSatFeeFloat := new(big.Float).SetUint64(uint64(MinSwapSatFee))

			//fmt.Printf("_swapFeeToken1ValueFloat: %v; _minSwapSatFeeFloat:%v\n", _swapFeeToken1ValueFloat, _minSwapSatFeeFloat)

			if _swapFeeToken1ValueFloat.Cmp(_minSwapSatFeeFloat) < 0 {
				// @dev: fee 20 sat

				_amountOut, err = getAmountOutBigWithoutFee(new(big.Int).Sub(_amountIn, _minSwapSatFee), _reserveIn, _reserveOut)
				if err != nil {
					return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBigWithoutFee")
				}

				//swapFeeType = SwapFee20Sat
				_swapFee = _minSwapSatFee
				_swapFeeFloat.SetInt(_swapFee)
				//fmt.Printf("_swapFee: %v\n", _swapFee)
			} else {
				// @dev: fee _swapFeeToken1ValueFloat

				_amountOut, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
				if err != nil {
					return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBig")
				}

				//swapFeeType = SwapFee6Thousands

				_swapFeeFloat = _swapFeeToken1ValueFloat

				// @dev: Set _swapFee
				_swapFeeToken1ValueFloat.Int(_swapFee)

				//fmt.Printf("_swapFee: %v\n", _swapFee)
			}
		} else {
			_amountOut, err = getAmountOutBig(_amountIn, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountOutBig")
			}
			//swapFeeType = SwapFee6ThousandsNotSat
			//fmt.Printf("_swapFee: %v\n", _swapFee)

		}
	}

	amountOut = _amountOut.String()
	err = nil
	return amountOut, err
}

func CalcAmountIn(tokenIn string, tokenOut string, amountOut string) (amountIn string, err error) {
	token0, token1, err := sortTokens(tokenIn, tokenOut)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "sortTokens")
	}

	// amountOut
	_amountOut, success := new(big.Int).SetString(amountOut, 10)
	if !success {
		return ZeroValue, errors.New("amountOut SetString(" + amountOut + ") " + strconv.FormatBool(success))
	}
	if _amountOut.Sign() < 0 {
		return ZeroValue, errors.New("amountOut(" + _amountOut.String() + ") is negative")
	}

	tx := middleware.DB.Begin()
	defer func() {
		tx.Rollback()
	}()

	// @dev: get pair
	var _pair PoolPair
	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	if err != nil {
		//@dev: pair does not exist
		return ZeroValue, utils.AppendErrorInfo(err, "pair does not exist")
	}
	pairId := _pair.ID
	if pairId <= 0 {
		return ZeroValue, errors.New("invalid pairId(" + strconv.FormatUint(uint64(pairId), 10) + ")")
	}

	var _reserve0, _reserve1 *big.Int
	// reserve0
	_reserve0, success = new(big.Int).SetString(_pair.Reserve0, 10)
	if !success {
		return ZeroValue, errors.New("Reserve0 SetString(" + _pair.Reserve0 + ") " + strconv.FormatBool(success))
	}
	// reserve1
	_reserve1, success = new(big.Int).SetString(_pair.Reserve1, 10)
	if !success {
		return ZeroValue, errors.New("Reserve1 SetString(" + _pair.Reserve1 + ") " + strconv.FormatBool(success))
	}

	var _reserveIn, _reserveOut = new(big.Int), new(big.Int)
	if token0 == tokenOut {
		*_reserveIn, *_reserveOut = *_reserve1, *_reserve0
	} else {
		*_reserveIn, *_reserveOut = *_reserve0, *_reserve1
	}

	var _amountIn *big.Int

	feeK := ProjectPartyFeeK + LpAwardFeeK

	_amountIn, err = getAmountInBig(_amountOut, _reserveIn, _reserveOut, feeK)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "getAmountInBig")
	}
	return _amountIn.String(), nil
}

func CalcAmountIn2(tokenIn string, tokenOut string, amountOut string) (amountIn string, err error) {
	feeK := ProjectPartyFeeK + LpAwardFeeK

	token0, token1, err := sortTokens(tokenIn, tokenOut)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "sortTokens")
	}

	isTokenZeroSat := token0 == TokenSatTag

	// amountOut
	_amountOut, success := new(big.Int).SetString(amountOut, 10)
	if !success {
		return ZeroValue, errors.New("amountOut SetString(" + amountOut + ") " + strconv.FormatBool(success))
	}

	if isTokenZeroSat {
		_minSwapSat := new(big.Int).SetUint64(uint64(MinSwapSatFee))
		if tokenOut == TokenSatTag {
			if _amountOut.Cmp(_minSwapSat) <= 0 {
				return ZeroValue, errors.New("insufficient _amountOut(" + _amountOut.String() + "), need " + _minSwapSat.String())
			}
		}
	}

	tx := middleware.DB.Begin()
	defer func() {
		tx.Rollback()
	}()

	// get pair
	var _pair PoolPair
	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	if err != nil {
		//@dev: pair does not exist
		return ZeroValue, utils.AppendErrorInfo(err, "pair does not exist")
	}

	var _reserve0, _reserve1 *big.Int
	// reserve0
	_reserve0, success = new(big.Int).SetString(_pair.Reserve0, 10)
	if !success {
		return ZeroValue, errors.New("Reserve0 SetString(" + _pair.Reserve0 + ") " + strconv.FormatBool(success))
	}
	// reserve1
	_reserve1, success = new(big.Int).SetString(_pair.Reserve1, 10)
	if !success {
		return ZeroValue, errors.New("Reserve1 SetString(" + _pair.Reserve1 + ") " + strconv.FormatBool(success))
	}

	var _reserveIn, _reserveOut = new(big.Int), new(big.Int)

	var _amountIn = new(big.Int)

	var _swapFee = big.NewInt(0)

	var _swapFeeFloat = big.NewFloat(0)

	if token0 == tokenOut {
		*_reserveIn, *_reserveOut = *_reserve1, *_reserve0

		if _amountOut.Cmp(_reserveOut) >= 0 {
			return ZeroValue, errors.New("excessive _amountOut(" + _amountOut.String() + "), need lt reserveOut(" + _reserveOut.String() + ")")
		}

		if isTokenZeroSat {
			// @dev: ?token => sat
			var _amountInWithFee, _amountInWithoutFee *big.Int

			_minSwapSatFee := new(big.Int).SetUint64(uint64(MinSwapSatFee))
			_amountInWithFee, err = getAmountInBig(_amountOut, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountInBig")
			}

			_amountInWithoutFee, err = getAmountInBigWithoutFee(_amountOut, _reserveIn, _reserveOut)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountInBigWithoutFee")
			}
			_amountInFee := new(big.Int).Sub(_amountInWithFee, _amountInWithoutFee)
			_amountInFeeFloat := new(big.Float).SetInt(_amountInFee)

			var _price *big.Float

			_price, err = getToken1PriceBig(_reserve0, _reserve1)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getToken1PriceBig")
			}
			_minSwapSatFeeFloat := new(big.Float).SetUint64(uint64(MinSwapSatFee))
			_feeValueFloat := new(big.Float).Mul(_amountInFeeFloat, _price)

			if _feeValueFloat.Cmp(_minSwapSatFeeFloat) < 0 {
				// @dev: fee 20 sat
				_amountIn, err = getAmountInBigWithoutFee(new(big.Int).Add(_amountOut, _minSwapSatFee), _reserveIn, _reserveOut)
				if err != nil {
					return ZeroValue, utils.AppendErrorInfo(err, "getAmountInBigWithoutFee")
				}
				//swapFeeType = SwapFee20Sat

				_swapFee = _minSwapSatFee
				_swapFeeFloat.SetInt(_swapFee)

				//fmt.Printf("_swapFee: %v\n", _swapFee)

			} else {
				// @dev: fee _feeValueFloat
				_amountIn = _amountInWithFee
				//swapFeeType = SwapFee6Thousands

				_swapFeeFloat = _feeValueFloat
				// @dev: Set _swapFee
				_feeValueFloat.Int(_swapFee)

				//fmt.Printf("_swapFee: %v\n", _swapFee)
			}

		} else {
			_amountIn, err = getAmountInBig(_amountOut, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountInBig")
			}
			//swapFeeType = SwapFee6ThousandsNotSat
			//fmt.Printf("_swapFee: %v\n", _swapFee)

		}

	} else {
		*_reserveIn, *_reserveOut = *_reserve0, *_reserve1

		if _amountOut.Cmp(_reserveOut) >= 0 {
			return ZeroValue, errors.New("excessive _amountOut(" + _amountOut.String() + "), need lt reserveOut(" + _reserveOut.String() + ")")
		}

		if isTokenZeroSat {
			// @dev: ?sat => token

			_minSwapSatFee := new(big.Int).SetUint64(uint64(MinSwapSatFee))

			var _amountInWithFee, _amountInWithoutFee *big.Int

			_amountInWithFee, err = getAmountInBig(_amountOut, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountInBig")
			}

			_amountInWithoutFee, err = getAmountInBigWithoutFee(_amountOut, _reserveIn, _reserveOut)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountInBigWithoutFee")
			}

			if new(big.Int).Sub(_amountInWithFee, _amountInWithoutFee).Cmp(_minSwapSatFee) < 0 {
				// @dev: fee 20 sat

				_amountIn = new(big.Int).Add(_amountInWithoutFee, _minSwapSatFee)
				//swapFeeType = SwapFee20Sat

				_swapFee = _minSwapSatFee
				_swapFeeFloat.SetInt(_swapFee)
				//fmt.Printf("_swapFee: %v\n", _swapFee)

			} else {
				// @dev: fee _amountInWithFee - _amountInWithoutFee

				_amountIn = _amountInWithFee
				//swapFeeType = SwapFee6Thousands

				_swapFee = new(big.Int).Sub(_amountInWithFee, _amountInWithoutFee)
				_swapFeeFloat.SetInt(_swapFee)

				//fmt.Printf("_swapFee: %v\n", _swapFee)
			}

		} else {
			_amountIn, err = getAmountInBig(_amountOut, _reserveIn, _reserveOut, feeK)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "getAmountInBig")
			}
			//swapFeeType = SwapFee6ThousandsNotSat
		}
	}

	if _amountIn.Sign() <= 0 {
		return ZeroValue, errors.New("invalid _amountIn(" + _amountIn.String() + ")")
	}

	amountIn = _amountIn.String()
	err = nil
	return amountIn, err
}

func (p *PoolShareRecord) ToShareRecordInfo() *ShareRecordInfo {
	if p == nil {
		return nil
	}
	return &ShareRecordInfo{
		ID:          p.ID,
		ShareId:     p.ShareId,
		Username:    p.Username,
		Liquidity:   p.Liquidity,
		Reserve0:    p.Reserve0,
		Reserve1:    p.Reserve1,
		Amount0:     p.Amount0,
		Amount1:     p.Amount1,
		ShareSupply: p.ShareSupply,
		ShareAmt:    p.ShareAmt,
		IsFirstMint: p.IsFirstMint,
		RecordType:  p.RecordType,
	}
}

func (p *PoolSwapRecord) ToSwapRecordInfo() *SwapRecordInfo {
	if p == nil {
		return nil
	}
	return &SwapRecordInfo{
		ID:             p.ID,
		PairId:         p.PairId,
		Username:       p.Username,
		TokenIn:        p.TokenIn,
		TokenOut:       p.TokenOut,
		AmountIn:       p.AmountIn,
		AmountOut:      p.AmountOut,
		ReserveIn:      p.ReserveIn,
		ReserveOut:     p.ReserveOut,
		SwapFee:        p.SwapFee,
		SwapFeeType:    p.SwapFeeType,
		SwapRecordType: p.SwapRecordType,
	}
}

type CalcAddLiquidityResponse struct {
	AmountA     string           `json:"amount_a"`
	AmountB     string           `json:"amount_b"`
	Liquidity   string           `json:"liquidity"`
	ShareRecord *ShareRecordInfo `json:"share_record"`
}

func CalcAddLiquidity(tokenA string, tokenB string, amountADesired string, amountBDesired string, amountAMin string, amountBMin string, username string) (amountA string, amountB string, liquidity string, shareRecord *PoolShareRecord, err error) {
	return calcAddLiquidity(tokenA, tokenB, amountADesired, amountBDesired, amountAMin, amountBMin, username)
}

type CalcRemoveLiquidityResponse struct {
	AmountA     string           `json:"amount_a"`
	AmountB     string           `json:"amount_b"`
	ShareRecord *ShareRecordInfo `json:"share_record"`
}

func CalcRemoveLiquidity(tokenA string, tokenB string, liquidity string, amountAMin string, amountBMin string, username string, feeK uint16) (amountA string, amountB string, shareRecord *PoolShareRecord, err error) {
	return calcRemoveLiquidity(tokenA, tokenB, liquidity, amountAMin, amountBMin, username, feeK)
}

// get amount

// TODO

type CalcSwapExactTokenForTokenNoPathResponse struct {
	AmountOut  string          `json:"amount_out"`
	SwapRecord *SwapRecordInfo `json:"swap_record"`
}

func CalcSwapExactTokenForTokenNoPath(tokenIn string, tokenOut string, amountIn string, amountOutMin string, username string, projectPartyFeeK uint16, lpAwardFeeK uint16) (amountOut string, swapRecord *PoolSwapRecord, err error) {
	return calcSwapExactTokenForTokenNoPath(tokenIn, tokenOut, amountIn, amountOutMin, username, projectPartyFeeK, lpAwardFeeK)
}

type CalcSwapTokenForExactTokenNoPathResponse struct {
	AmountIn   string          `json:"amount_in"`
	SwapRecord *SwapRecordInfo `json:"swap_record"`
}

func CalcSwapTokenForExactTokenNoPath(tokenIn string, tokenOut string, amountOut string, amountInMax string, username string, projectPartyFeeK uint16, lpAwardFeeK uint16) (amountIn string, swapRecord *PoolSwapRecord, err error) {
	return calcSwapTokenForExactTokenNoPath(tokenIn, tokenOut, amountOut, amountInMax, username, projectPartyFeeK, lpAwardFeeK)
}

// Query

type PoolInfo struct {
	PairId         uint   `json:"pair_id"`
	ShareId        uint   `json:"share_id"`
	IsTokenZeroSat bool   `json:"is_token_zero_sat"`
	Token0         string `json:"token0"`
	Token1         string `json:"token1"`
	Reserve0       string `json:"reserve0"`
	Reserve1       string `json:"reserve1"`
	Liquidity      string `json:"liquidity"`
}

func queryPoolInfo(tokenA string, tokenB string) (poolInfo *PoolInfo, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return new(PoolInfo), utils.AppendErrorInfo(err, "sortTokens")
	}

	tx := middleware.DB.Begin()

	var _poolInfo PoolInfo
	err = tx.Table("pool_pairs").
		Joins("join pool_shares on pool_pairs.id = pool_shares.pair_id").
		Select("pool_pairs.id as pair_id,pool_shares.id as share_id,pool_pairs.is_token_zero_sat,pool_pairs.token0,pool_pairs.token1,pool_pairs.reserve0,pool_pairs.reserve1,pool_shares.total_supply as liquidity").
		Where("pool_pairs.token0 = ? AND pool_pairs.token1 = ?", token0, token1).
		Scan(&_poolInfo).
		Error
	if err != nil {
		return new(PoolInfo), utils.AppendErrorInfo(err, "pair info not found")
	}

	tx.Rollback()
	if _poolInfo.PairId == 0 {
		_poolInfo.Token0 = token0
		_poolInfo.Token1 = token1
		_poolInfo.Reserve0 = ZeroValue
		_poolInfo.Reserve1 = ZeroValue
		_poolInfo.Liquidity = ZeroValue
	}

	poolInfo = &_poolInfo
	return poolInfo, nil
}

var PoolDoesNotExistErr = errors.New("pool does not exist")

func QueryPoolInfo(tokenA string, tokenB string) (poolInfo *PoolInfo, err error) {
	poolInfo, err = queryPoolInfo(tokenA, tokenB)
	if err != nil {
		return new(PoolInfo), utils.AppendErrorInfo(err, "queryPoolInfo")
	} else {
		if poolInfo == nil {
			return new(PoolInfo), errors.New("pool info nil")
		} else if poolInfo.PairId == 0 {
			return poolInfo, PoolDoesNotExistErr
		}
		return poolInfo, nil
	}
}

type ShareRecordInfo struct {
	ID          uint            `json:"id"`
	Time        int64           `json:"time"`
	ShareId     uint            `json:"share_id"`
	Username    string          `json:"username"`
	Liquidity   string          `json:"liquidity"`
	Reserve0    string          `json:"reserve0"`
	Reserve1    string          `json:"reserve1"`
	Amount0     string          `json:"amount0"`
	Amount1     string          `json:"amount1"`
	ShareSupply string          `json:"share_supply"`
	ShareAmt    string          `json:"share_amt"`
	IsFirstMint bool            `json:"is_first_mint"`
	RecordType  ShareRecordType `json:"record_type"`
}

type ShareRecordInfoScan struct {
	ID          uint            `json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	ShareId     uint            `json:"share_id"`
	Username    string          `json:"username"`
	Liquidity   string          `json:"liquidity"`
	Reserve0    string          `json:"reserve0"`
	Reserve1    string          `json:"reserve1"`
	Amount0     string          `json:"amount0"`
	Amount1     string          `json:"amount1"`
	ShareSupply string          `json:"share_supply"`
	ShareAmt    string          `json:"share_amt"`
	IsFirstMint bool            `json:"is_first_mint"`
	RecordType  ShareRecordType `json:"record_type"`
}

func ProcessShareRecordInfoScan(record ShareRecordInfoScan) ShareRecordInfo {
	return ShareRecordInfo{
		ID:          record.ID,
		Time:        record.CreatedAt.Unix(),
		ShareId:     record.ShareId,
		Username:    record.Username,
		Liquidity:   record.Liquidity,
		Reserve0:    record.Reserve0,
		Reserve1:    record.Reserve1,
		Amount0:     record.Amount0,
		Amount1:     record.Amount1,
		ShareSupply: record.ShareSupply,
		ShareAmt:    record.ShareAmt,
		IsFirstMint: record.IsFirstMint,
		RecordType:  record.RecordType,
	}
}

func QueryShareRecords(tokenA string, tokenB string, limit int, offset int) (shareRecordInfos *[]ShareRecordInfo, err error) {

	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return new([]ShareRecordInfo), utils.AppendErrorInfo(err, "sortTokens")
	}

	tx := middleware.DB.Begin()

	var _shareRecordInfos []ShareRecordInfo
	var _shareRecordInfosScan []ShareRecordInfoScan

	err = tx.Table("pool_pairs").
		Select("pool_share_records.id,pool_share_records.created_at,pool_share_records.share_id,pool_share_records.username,pool_share_records.liquidity,pool_share_records.reserve0,pool_share_records.reserve1,pool_share_records.amount0,pool_share_records.amount1,pool_share_records.share_supply,pool_share_records.share_amt,pool_share_records.is_first_mint,pool_share_records.record_type").
		Joins("join pool_shares on pool_pairs.id = pool_shares.pair_id").
		Joins("join pool_share_records on pool_shares.id = pool_share_records.share_id").
		Where("pool_pairs.token0 = ? AND pool_pairs.token1 = ?", token0, token1).
		Order("pool_share_records.id desc").
		Limit(limit).
		Offset(offset).
		Scan(&_shareRecordInfosScan).
		Error
	if err != nil {
		return new([]ShareRecordInfo), utils.AppendErrorInfo(err, "select ShareRecordInfo")
	}

	tx.Rollback()

	if _shareRecordInfosScan == nil {
		_shareRecordInfos = make([]ShareRecordInfo, 0)
	} else {
		for _, record := range _shareRecordInfosScan {
			_shareRecordInfos = append(_shareRecordInfos, ProcessShareRecordInfoScan(record))
		}
	}

	shareRecordInfos = &_shareRecordInfos
	return shareRecordInfos, nil
}

func QueryUserShareRecords(tokenA string, tokenB string, username string, limit int, offset int) (shareRecordInfos *[]ShareRecordInfo, err error) {

	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return new([]ShareRecordInfo), utils.AppendErrorInfo(err, "sortTokens")
	}

	tx := middleware.DB.Begin()

	var _shareRecordInfos []ShareRecordInfo
	var _shareRecordInfosScan []ShareRecordInfoScan

	err = tx.Table("pool_pairs").
		Select("pool_share_records.id,pool_share_records.created_at,pool_share_records.share_id,pool_share_records.username,pool_share_records.liquidity,pool_share_records.reserve0,pool_share_records.reserve1,pool_share_records.amount0,pool_share_records.amount1,pool_share_records.share_supply,pool_share_records.share_amt,pool_share_records.is_first_mint,pool_share_records.record_type").
		Joins("join pool_shares on pool_pairs.id = pool_shares.pair_id").
		Joins("join pool_share_records on pool_shares.id = pool_share_records.share_id").
		Where("pool_pairs.token0 = ? and pool_pairs.token1 = ? and pool_share_records.username = ?", token0, token1, username).
		Order("pool_share_records.id desc").
		Limit(limit).
		Offset(offset).
		Scan(&_shareRecordInfosScan).
		Error
	if err != nil {
		return new([]ShareRecordInfo), utils.AppendErrorInfo(err, "select ShareRecordInfo")
	}

	tx.Rollback()

	if _shareRecordInfosScan == nil {
		_shareRecordInfos = make([]ShareRecordInfo, 0)
	} else {
		for _, record := range _shareRecordInfosScan {
			_shareRecordInfos = append(_shareRecordInfos, ProcessShareRecordInfoScan(record))
		}
	}
	shareRecordInfos = &_shareRecordInfos

	return shareRecordInfos, nil
}

type ShareRecordInfoIncludeToken struct {
	ID          uint            `json:"id"`
	Time        int64           `json:"time"`
	Token0      string          `json:"token0"`
	Token1      string          `json:"token1"`
	ShareId     uint            `json:"share_id"`
	Username    string          `json:"username"`
	Liquidity   string          `json:"liquidity"`
	Reserve0    string          `json:"reserve0"`
	Reserve1    string          `json:"reserve1"`
	Amount0     string          `json:"amount0"`
	Amount1     string          `json:"amount1"`
	ShareSupply string          `json:"share_supply"`
	ShareAmt    string          `json:"share_amt"`
	IsFirstMint bool            `json:"is_first_mint"`
	RecordType  ShareRecordType `json:"record_type"`
}

type ShareRecordInfoIncludeTokenScan struct {
	ID          uint            `json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	Token0      string          `json:"token0"`
	Token1      string          `json:"token1"`
	ShareId     uint            `json:"share_id"`
	Username    string          `json:"username"`
	Liquidity   string          `json:"liquidity"`
	Reserve0    string          `json:"reserve0"`
	Reserve1    string          `json:"reserve1"`
	Amount0     string          `json:"amount0"`
	Amount1     string          `json:"amount1"`
	ShareSupply string          `json:"share_supply"`
	ShareAmt    string          `json:"share_amt"`
	IsFirstMint bool            `json:"is_first_mint"`
	RecordType  ShareRecordType `json:"record_type"`
}

func ProcessShareRecordInfoIncludeTokenScan(record ShareRecordInfoIncludeTokenScan) ShareRecordInfoIncludeToken {
	return ShareRecordInfoIncludeToken{
		ID:          record.ID,
		Time:        record.CreatedAt.Unix(),
		Token0:      record.Token0,
		Token1:      record.Token1,
		ShareId:     record.ShareId,
		Username:    record.Username,
		Liquidity:   record.Liquidity,
		Reserve0:    record.Reserve0,
		Reserve1:    record.Reserve1,
		Amount0:     record.Amount0,
		Amount1:     record.Amount1,
		ShareSupply: record.ShareSupply,
		ShareAmt:    record.ShareAmt,
		IsFirstMint: record.IsFirstMint,
		RecordType:  record.RecordType,
	}
}

func QueryUserAllShareRecordsCount(username string) (count int64, err error) {

	tx := middleware.DB.Begin()

	err = tx.Table("pool_pairs").
		Joins("join pool_shares on pool_pairs.id = pool_shares.pair_id").
		Joins("join pool_share_records on pool_shares.id = pool_share_records.share_id").
		Where("pool_share_records.username = ?", username).
		Count(&count).
		Error
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "select ShareRecordInfo")
	}

	tx.Rollback()

	return count, nil
}

func QueryUserAllShareRecords(username string, limit int, offset int) (shareRecordInfos *[]ShareRecordInfoIncludeToken, err error) {

	tx := middleware.DB.Begin()

	var _shareRecordInfos []ShareRecordInfoIncludeToken
	var _shareRecordInfosScan []ShareRecordInfoIncludeTokenScan

	err = tx.Table("pool_pairs").
		Select("pool_share_records.id,pool_share_records.created_at,pool_pairs.token0,pool_pairs.token1,pool_share_records.share_id,pool_share_records.username,pool_share_records.liquidity,pool_share_records.reserve0,pool_share_records.reserve1,pool_share_records.amount0,pool_share_records.amount1,pool_share_records.share_supply,pool_share_records.share_amt,pool_share_records.is_first_mint,pool_share_records.record_type").
		Joins("join pool_shares on pool_pairs.id = pool_shares.pair_id").
		Joins("join pool_share_records on pool_shares.id = pool_share_records.share_id").
		Where("pool_share_records.username = ?", username).
		Order("pool_share_records.id desc").
		Limit(limit).
		Offset(offset).
		Scan(&_shareRecordInfosScan).
		Error
	if err != nil {
		return new([]ShareRecordInfoIncludeToken), utils.AppendErrorInfo(err, "select ShareRecordInfo")
	}

	tx.Rollback()

	if _shareRecordInfosScan == nil {
		_shareRecordInfos = make([]ShareRecordInfoIncludeToken, 0)
	} else {
		for _, record := range _shareRecordInfosScan {
			_shareRecordInfos = append(_shareRecordInfos, ProcessShareRecordInfoIncludeTokenScan(record))
		}
	}

	shareRecordInfos = &_shareRecordInfos

	return shareRecordInfos, nil
}

type SwapRecordInfo struct {
	ID             uint           `json:"id"`
	Time           int64          `json:"time"`
	PairId         uint           `json:"pair_id"`
	Username       string         `json:"username"`
	TokenIn        string         `json:"token_in"`
	TokenOut       string         `json:"token_out"`
	AmountIn       string         `json:"amount_in"`
	AmountOut      string         `json:"amount_out"`
	ReserveIn      string         `json:"reserve_in"`
	ReserveOut     string         `json:"reserve_out"`
	SwapFee        string         `json:"swap_fee"`
	SwapFeeType    SwapFeeType    `json:"swap_fee_type"`
	SwapRecordType SwapRecordType `json:"swap_record_type"`
}

type SwapRecordInfoScan struct {
	ID             uint           `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	PairId         uint           `json:"pair_id"`
	Username       string         `json:"username"`
	TokenIn        string         `json:"token_in"`
	TokenOut       string         `json:"token_out"`
	AmountIn       string         `json:"amount_in"`
	AmountOut      string         `json:"amount_out"`
	ReserveIn      string         `json:"reserve_in"`
	ReserveOut     string         `json:"reserve_out"`
	SwapFee        string         `json:"swap_fee"`
	SwapFeeType    SwapFeeType    `json:"swap_fee_type"`
	SwapRecordType SwapRecordType `json:"swap_record_type"`
}

func PrecessSwapRecordInfoScan(record SwapRecordInfoScan) SwapRecordInfo {
	return SwapRecordInfo{
		ID:             record.ID,
		Time:           record.CreatedAt.Unix(),
		PairId:         record.PairId,
		Username:       record.Username,
		TokenIn:        record.TokenIn,
		TokenOut:       record.TokenOut,
		AmountIn:       record.AmountIn,
		AmountOut:      record.AmountOut,
		ReserveIn:      record.ReserveIn,
		ReserveOut:     record.ReserveOut,
		SwapFee:        record.SwapFee,
		SwapFeeType:    record.SwapFeeType,
		SwapRecordType: record.SwapRecordType,
	}

}

func QuerySwapRecords(tokenA string, tokenB string, limit int, offset int) (swapRecordInfos *[]SwapRecordInfo, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return new([]SwapRecordInfo), utils.AppendErrorInfo(err, "sortTokens")
	}

	tx := middleware.DB.Begin()

	var _swapRecordInfos []SwapRecordInfo
	var _swapRecordInfosScan []SwapRecordInfoScan

	err = tx.Table("pool_pairs").
		Select("pool_swap_records.id,pool_swap_records.created_at,pool_swap_records.pair_id,pool_swap_records.username,pool_swap_records.token_in,pool_swap_records.token_out,pool_swap_records.amount_in,pool_swap_records.amount_out,pool_swap_records.reserve_in,pool_swap_records.reserve_out,pool_swap_records.swap_fee,pool_swap_records.swap_fee_type,pool_swap_records.swap_record_type").
		Joins("join pool_swap_records on pool_pairs.id = pool_swap_records.pair_id").
		Where("pool_pairs.token0 = ? AND pool_pairs.token1 = ?", token0, token1).
		Order("pool_swap_records.id desc").
		Limit(limit).
		Offset(offset).
		Scan(&_swapRecordInfosScan).
		Error
	if err != nil {
		return new([]SwapRecordInfo), utils.AppendErrorInfo(err, "select SwapRecordInfo")
	}

	tx.Rollback()

	if _swapRecordInfosScan == nil {
		_swapRecordInfos = make([]SwapRecordInfo, 0)
	} else {
		for _, record := range _swapRecordInfosScan {
			_swapRecordInfos = append(_swapRecordInfos, PrecessSwapRecordInfoScan(record))
		}
	}

	swapRecordInfos = &_swapRecordInfos
	return swapRecordInfos, nil
}

func QueryUserSwapRecords(tokenA string, tokenB string, username string, limit int, offset int) (swapRecordInfos *[]SwapRecordInfo, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return new([]SwapRecordInfo), utils.AppendErrorInfo(err, "sortTokens")
	}

	tx := middleware.DB.Begin()

	var _swapRecordInfos []SwapRecordInfo
	var _swapRecordInfosScan []SwapRecordInfoScan

	err = tx.Table("pool_pairs").
		Select("pool_swap_records.id,pool_swap_records.created_at,pool_swap_records.pair_id,pool_swap_records.username,pool_swap_records.token_in,pool_swap_records.token_out,pool_swap_records.amount_in,pool_swap_records.amount_out,pool_swap_records.reserve_in,pool_swap_records.reserve_out,pool_swap_records.swap_fee,pool_swap_records.swap_fee_type,pool_swap_records.swap_record_type").
		Joins("join pool_swap_records on pool_pairs.id = pool_swap_records.pair_id").
		Where("pool_pairs.token0 = ? and pool_pairs.token1 = ? and pool_swap_records.username = ?", token0, token1, username).
		Order("pool_swap_records.id desc").
		Limit(limit).
		Offset(offset).
		Scan(&_swapRecordInfosScan).
		Error
	if err != nil {
		return new([]SwapRecordInfo), utils.AppendErrorInfo(err, "select SwapRecordInfo")
	}

	tx.Rollback()

	if _swapRecordInfosScan == nil {
		_swapRecordInfos = make([]SwapRecordInfo, 0)
	} else {
		for _, record := range _swapRecordInfosScan {
			_swapRecordInfos = append(_swapRecordInfos, PrecessSwapRecordInfoScan(record))
		}
	}
	swapRecordInfos = &_swapRecordInfos
	return swapRecordInfos, nil
}

func QueryUserAllSwapRecordsCount(username string) (count int64, err error) {
	tx := middleware.DB.Begin()

	err = tx.Table("pool_pairs").
		Joins("join pool_swap_records on pool_pairs.id = pool_swap_records.pair_id").
		Where("pool_swap_records.username = ?", username).
		Count(&count).
		Error
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "select SwapRecordInfoIncludeToken count")
	}

	tx.Rollback()
	return count, nil
}

type SwapRecordInfoIncludeToken struct {
	ID             uint           `json:"id"`
	Time           int64          `json:"time"`
	Token0         string         `json:"token0"`
	Token1         string         `json:"token1"`
	PairId         uint           `json:"pair_id"`
	Username       string         `json:"username"`
	TokenIn        string         `json:"token_in"`
	TokenOut       string         `json:"token_out"`
	AmountIn       string         `json:"amount_in"`
	AmountOut      string         `json:"amount_out"`
	ReserveIn      string         `json:"reserve_in"`
	ReserveOut     string         `json:"reserve_out"`
	SwapFee        string         `json:"swap_fee"`
	SwapFeeType    SwapFeeType    `json:"swap_fee_type"`
	SwapRecordType SwapRecordType `json:"swap_record_type"`
}

type SwapRecordInfoIncludeTokenScan struct {
	ID             uint           `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	Token0         string         `json:"token0"`
	Token1         string         `json:"token1"`
	PairId         uint           `json:"pair_id"`
	Username       string         `json:"username"`
	TokenIn        string         `json:"token_in"`
	TokenOut       string         `json:"token_out"`
	AmountIn       string         `json:"amount_in"`
	AmountOut      string         `json:"amount_out"`
	ReserveIn      string         `json:"reserve_in"`
	ReserveOut     string         `json:"reserve_out"`
	SwapFee        string         `json:"swap_fee"`
	SwapFeeType    SwapFeeType    `json:"swap_fee_type"`
	SwapRecordType SwapRecordType `json:"swap_record_type"`
}

func ProcessSwapRecordInfoIncludeTokenScan(record SwapRecordInfoIncludeTokenScan) SwapRecordInfoIncludeToken {
	return SwapRecordInfoIncludeToken{
		ID:             record.ID,
		Time:           record.CreatedAt.Unix(),
		Token0:         record.Token0,
		Token1:         record.Token1,
		PairId:         record.PairId,
		Username:       record.Username,
		TokenIn:        record.TokenIn,
		TokenOut:       record.TokenOut,
		AmountIn:       record.AmountIn,
		AmountOut:      record.AmountOut,
		ReserveIn:      record.ReserveIn,
		ReserveOut:     record.ReserveOut,
		SwapFee:        record.SwapFee,
		SwapFeeType:    record.SwapFeeType,
		SwapRecordType: record.SwapRecordType,
	}
}

func QueryUserAllSwapRecords(username string, limit int, offset int) (swapRecordInfos *[]SwapRecordInfoIncludeToken, err error) {
	tx := middleware.DB.Begin()

	var _swapRecordInfos []SwapRecordInfoIncludeToken

	var _swapRecordInfosScan []SwapRecordInfoIncludeTokenScan

	err = tx.Table("pool_pairs").
		Select("pool_swap_records.id,pool_swap_records.created_at,pool_pairs.token0,pool_pairs.token1,pool_swap_records.pair_id,pool_swap_records.username,pool_swap_records.token_in,pool_swap_records.token_out,pool_swap_records.amount_in,pool_swap_records.amount_out,pool_swap_records.reserve_in,pool_swap_records.reserve_out,pool_swap_records.swap_fee,pool_swap_records.swap_fee_type,pool_swap_records.swap_record_type").
		Joins("join pool_swap_records on pool_pairs.id = pool_swap_records.pair_id").
		Where("pool_swap_records.username = ?", username).
		Order("pool_swap_records.id desc").
		Limit(limit).
		Offset(offset).
		Scan(&_swapRecordInfosScan).
		Error
	if err != nil {
		return new([]SwapRecordInfoIncludeToken), utils.AppendErrorInfo(err, "select SwapRecordInfoIncludeToken")
	}

	tx.Rollback()

	if _swapRecordInfosScan == nil {
		_swapRecordInfos = make([]SwapRecordInfoIncludeToken, 0)
	} else {

		for _, record := range _swapRecordInfosScan {
			_swapRecordInfos = append(_swapRecordInfos, ProcessSwapRecordInfoIncludeTokenScan(record))
		}
	}

	swapRecordInfos = &_swapRecordInfos
	return swapRecordInfos, nil
}

type LpAwardBalanceInfo struct {
	ID      uint   `json:"id"`
	Balance string `json:"balance"`
}

func QueryUserLpAwardBalance(username string) (lpAwardBalanceInfo *LpAwardBalanceInfo, err error) {

	tx := middleware.DB.Begin()

	var _lpAwardBalanceInfo LpAwardBalanceInfo

	err = tx.Table("pool_lp_award_balances").
		Select("id,balance").
		Where("username = ?", username).
		Scan(&_lpAwardBalanceInfo).
		Error
	if err != nil {
		return new(LpAwardBalanceInfo), utils.AppendErrorInfo(err, "select LpAwardBalanceInfo")
	}

	tx.Rollback()

	if _lpAwardBalanceInfo.ID == 0 {
		_lpAwardBalanceInfo.Balance = ZeroValue
	}

	lpAwardBalanceInfo = &_lpAwardBalanceInfo
	return lpAwardBalanceInfo, nil
}

type WithdrawAwardRecordInfo struct {
	ID           uint   `json:"id"`
	Time         int64  `json:"time"`
	Username     string `json:"username"`
	Amount       string `json:"amount"`
	AwardBalance string `json:"award_balance"`
}

type WithdrawAwardRecordInfoScan struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Username     string    `json:"username"`
	Amount       string    `json:"amount"`
	AwardBalance string    `json:"award_balance"`
}

func ProcessWithdrawAwardRecordInfoScan(record WithdrawAwardRecordInfoScan) WithdrawAwardRecordInfo {
	return WithdrawAwardRecordInfo{
		ID:           record.ID,
		Time:         record.CreatedAt.Unix(),
		Username:     record.Username,
		Amount:       record.Amount,
		AwardBalance: record.AwardBalance,
	}
}

// not used
func QueryWithdrawAwardRecords(limit int, offset int) (withdrawAwardRecords *[]PoolWithdrawAwardRecord, err error) {
	tx := middleware.DB.Begin()

	var _withdrawAwardRecords []PoolWithdrawAwardRecord

	err = tx.Table("pool_withdraw_award_records").
		Select("id,username,amount,award_balance").
		Order("id desc").
		Limit(limit).
		Offset(offset).
		Scan(&_withdrawAwardRecords).
		Error
	if err != nil {
		return new([]PoolWithdrawAwardRecord), utils.AppendErrorInfo(err, "select PoolWithdrawAwardRecord")
	}

	tx.Rollback()
	withdrawAwardRecords = &_withdrawAwardRecords
	return withdrawAwardRecords, nil
}

func QueryUserWithdrawAwardRecords(username string, limit int, offset int) (withdrawAwardRecords *[]WithdrawAwardRecordInfo, err error) {
	tx := middleware.DB.Begin()

	var _withdrawAwardRecords []WithdrawAwardRecordInfo
	var _withdrawAwardRecordsScan []WithdrawAwardRecordInfoScan

	err = tx.Table("pool_withdraw_award_records").
		Select("id,created_at,username,amount,award_balance").
		Where("username = ?", username).
		Order("id desc").
		Limit(limit).
		Offset(offset).
		Scan(&_withdrawAwardRecordsScan).
		Error
	if err != nil {
		return new([]WithdrawAwardRecordInfo), utils.AppendErrorInfo(err, "select PoolWithdrawAwardRecord")
	}

	tx.Rollback()

	if _withdrawAwardRecordsScan == nil {
		_withdrawAwardRecords = make([]WithdrawAwardRecordInfo, 0)
	} else {
		for _, record := range _withdrawAwardRecordsScan {
			_withdrawAwardRecords = append(_withdrawAwardRecords, ProcessWithdrawAwardRecordInfoScan(record))
		}
	}

	withdrawAwardRecords = &_withdrawAwardRecords
	return withdrawAwardRecords, nil
}

// count

func QueryShareRecordsCount(tokenA string, tokenB string) (count int64, err error) {

	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "sortTokens")
	}

	tx := middleware.DB.Begin()

	err = tx.Table("pool_pairs").
		Joins("join pool_shares on pool_pairs.id = pool_shares.pair_id").
		Joins("join pool_share_records on pool_shares.id = pool_share_records.share_id").
		Where("pool_pairs.token0 = ? AND pool_pairs.token1 = ?", token0, token1).
		Count(&count).
		Error
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "select ShareRecordInfo count")
	}

	tx.Rollback()
	return count, nil
}

func QueryUserShareRecordsCount(tokenA string, tokenB string, username string) (count int64, err error) {

	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "sortTokens")
	}

	tx := middleware.DB.Begin()

	err = tx.Table("pool_pairs").
		Joins("join pool_shares on pool_pairs.id = pool_shares.pair_id").
		Joins("join pool_share_records on pool_shares.id = pool_share_records.share_id").
		Where("pool_pairs.token0 = ? and pool_pairs.token1 = ? and pool_share_records.username = ?", token0, token1, username).
		Count(&count).
		Error
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "select ShareRecordInfo count")
	}

	tx.Rollback()
	return count, nil
}

func QuerySwapRecordsCount(tokenA string, tokenB string) (count int64, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "sortTokens")
	}

	tx := middleware.DB.Begin()

	err = tx.Table("pool_pairs").
		Joins("join pool_swap_records on pool_pairs.id = pool_swap_records.pair_id").
		Where("pool_pairs.token0 = ? AND pool_pairs.token1 = ?", token0, token1).
		Count(&count).
		Error
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "select SwapRecordInfo count")
	}

	tx.Rollback()
	return count, nil
}

func QueryUserSwapRecordsCount(tokenA string, tokenB string, username string) (count int64, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "sortTokens")
	}

	tx := middleware.DB.Begin()

	err = tx.Table("pool_pairs").
		Joins("join pool_swap_records on pool_pairs.id = pool_swap_records.pair_id").
		Where("pool_pairs.token0 = ? and pool_pairs.token1 = ? and pool_swap_records.username = ?", token0, token1, username).
		Count(&count).
		Error
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "select SwapRecordInfo count")
	}

	tx.Rollback()
	return count, nil
}

// not used
func QueryWithdrawAwardRecordsCount() (count int64, err error) {
	tx := middleware.DB.Begin()

	err = tx.Table("pool_withdraw_award_records").
		Count(&count).
		Error
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "select PoolWithdrawAwardRecord count")
	}

	tx.Rollback()
	return count, nil
}

func QueryUserWithdrawAwardRecordsCount(username string) (count int64, err error) {
	tx := middleware.DB.Begin()

	err = tx.Table("pool_withdraw_award_records").
		Where("username = ?", username).
		Count(&count).
		Error
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "select PoolWithdrawAwardRecord count")
	}

	tx.Rollback()
	return count, nil
}

// @dev: Tokens to share id

func QueryShareId(tx *gorm.DB, tokenA string, tokenB string) (shareId uint, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "sortTokens")
	}

	err = tx.Table("pool_pairs").
		Joins("join pool_shares on pool_pairs.id = pool_shares.pair_id").
		Where("pool_pairs.token0 = ? AND pool_pairs.token1 = ?", token0, token1).
		Select("pool_shares.id").
		Scan(&shareId).
		Error
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "select share id")
	}

	return shareId, nil
}

func QueryPairAndShareId(tx *gorm.DB, tokenA string, tokenB string) (pairId uint, shareId uint, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return 0, 0, utils.AppendErrorInfo(err, "sortTokens")
	}

	// @dev: get pair
	var _pair PoolPair
	err = tx.Model(&PoolPair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	if err != nil {
		//@dev: pair does not exist
		return 0, 0, utils.AppendErrorInfo(err, "pair does not exist")
	}

	pairId = _pair.ID
	if pairId <= 0 {
		return 0, 0, errors.New("invalid pairId(" + strconv.FormatUint(uint64(pairId), 10) + ")")
	}

	// get share
	var share PoolShare
	err = tx.Model(&PoolShare{}).Where("pair_id = ?", pairId).First(&share).Error
	if err != nil {
		return 0, 0, utils.AppendErrorInfo(err, "share does not exist")
	}

	shareId = share.ID
	if shareId <= 0 {
		return 0, 0, errors.New("invalid shareId(" + strconv.FormatUint(uint64(shareId), 10) + ")")
	}

	return pairId, shareId, nil
}

// @dev: Query liquidity, lp_award_balance, lp_award_cumulative

func QueryLiquidityAndAwardRecordsCount(username string) (count int64, err error) {
	tx := middleware.DB.Begin()

	err = tx.Table("pool_pairs").
		Joins("join pool_shares on pool_pairs.id = pool_shares.pair_id").
		Joins("join pool_share_balances on pool_shares.id = pool_share_balances.share_id").
		Joins("join pool_share_lp_award_balances on pool_share_balances.username = pool_share_lp_award_balances.username").
		Joins("join pool_share_lp_award_cumulatives on pool_share_balances.username = pool_share_lp_award_cumulatives.username").
		Where("pool_share_balances.username = ?", username).
		Count(&count).Error
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "select liquidity, lp_award_balance, lp_award_cumulative count")
	}
	tx.Rollback()

	return count, nil
}

type LiquidityAndAwardRecordInfo struct {
	Token0            string `json:"token_0"`
	Token1            string `json:"token_1"`
	Liquidity         string `json:"liquidity"`
	LpAwardBalance    string `json:"lp_award_balance"`
	LpAwardCumulative string `json:"lp_award_cumulative"`
}

func QueryLiquidityAndAwardRecords(username string, limit int, offset int) (records *[]LiquidityAndAwardRecordInfo, err error) {
	tx := middleware.DB.Begin()

	var _liquidityAndAwardRecordInfo []LiquidityAndAwardRecordInfo

	err = tx.Table("pool_pairs").
		Joins("join pool_shares on pool_pairs.id = pool_shares.pair_id").
		Joins("join pool_share_balances on pool_shares.id = pool_share_balances.share_id").
		Joins("join pool_share_lp_award_balances on pool_share_balances.username = pool_share_lp_award_balances.username").
		Joins("join pool_share_lp_award_cumulatives on pool_share_balances.username = pool_share_lp_award_cumulatives.username").
		Select("pool_pairs.token0, pool_pairs.token1, pool_share_balances.balance as liquidity, pool_share_lp_award_balances.balance as lp_award_balance, pool_share_lp_award_cumulatives.amount as lp_award_cumulative").
		Where("pool_share_balances.username = ?", username).
		Order("pool_share_balances.id desc").
		Limit(limit).
		Offset(offset).
		Scan(&_liquidityAndAwardRecordInfo).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "select liquidity, lp_award_balance, lp_award_cumulative")
	}

	if _liquidityAndAwardRecordInfo == nil {
		_liquidityAndAwardRecordInfo = make([]LiquidityAndAwardRecordInfo, 0)
	}
	tx.Rollback()

	return &_liquidityAndAwardRecordInfo, nil
}

func QueryLpAwardRecordsCount(username string) (count int64, err error) {
	tx := middleware.DB.Begin()

	err = tx.Table("pool_lp_award_records").
		Where("username = ?", username).
		Count(&count).
		Error
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "select LpAwardRecordInfo count")
	}

	tx.Rollback()
	return count, nil
}

type LpAwardRecordInfo struct {
	ID           uint      `json:"id"`
	Time         int64     `json:"time"`
	ShareId      uint      `json:"share_id" gorm:"index"`
	Amount       string    `json:"amount" gorm:"type:varchar(255);index"`
	Fee          string    `json:"fee" gorm:"type:varchar(255);index"`
	AwardBalance string    `json:"award_balance" gorm:"type:varchar(255);index"`
	ShareBalance string    `json:"share_balance" gorm:"type:varchar(255);index"`
	TotalSupply  string    `json:"total_supply" gorm:"type:varchar(255);index"`
	SwapRecordId uint      `json:"swap_record_id" gorm:"index"`
	AwardType    AwardType `json:"award_type" gorm:"index"`
}

type LpAwardRecordInfoScan struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	ShareId      uint      `json:"share_id" gorm:"index"`
	Amount       string    `json:"amount" gorm:"type:varchar(255);index"`
	Fee          string    `json:"fee" gorm:"type:varchar(255);index"`
	AwardBalance string    `json:"award_balance" gorm:"type:varchar(255);index"`
	ShareBalance string    `json:"share_balance" gorm:"type:varchar(255);index"`
	TotalSupply  string    `json:"total_supply" gorm:"type:varchar(255);index"`
	SwapRecordId uint      `json:"swap_record_id" gorm:"index"`
	AwardType    AwardType `json:"award_type" gorm:"index"`
}

func ProcessLpAwardRecordInfoScan(record LpAwardRecordInfoScan) LpAwardRecordInfo {
	return LpAwardRecordInfo{
		ID:           record.ID,
		Time:         record.CreatedAt.Unix(),
		ShareId:      record.ShareId,
		Amount:       record.Amount,
		Fee:          record.Fee,
		AwardBalance: record.AwardBalance,
		ShareBalance: record.ShareBalance,
		TotalSupply:  record.TotalSupply,
		SwapRecordId: record.SwapRecordId,
		AwardType:    record.AwardType,
	}
}

func QueryLpAwardRecords(username string, limit int, offset int) (lpAwardRecords *[]LpAwardRecordInfo, err error) {
	tx := middleware.DB.Begin()

	var _lpAwardRecords []LpAwardRecordInfo
	var _lpAwardRecordsScan []LpAwardRecordInfoScan

	err = tx.Table("pool_lp_award_records").
		Select("id,created_at,share_id,amount,fee,award_balance,share_balance,total_supply,swap_record_id,award_type").
		Where("username = ?", username).
		Order("id desc").
		Limit(limit).
		Offset(offset).
		Scan(&_lpAwardRecordsScan).
		Error
	if err != nil {
		return new([]LpAwardRecordInfo), utils.AppendErrorInfo(err, "select LpAwardRecordInfo")
	}

	tx.Rollback()

	if _lpAwardRecordsScan == nil {
		_lpAwardRecords = make([]LpAwardRecordInfo, 0)
	} else {
		for _, record := range _lpAwardRecordsScan {
			_lpAwardRecords = append(_lpAwardRecords, ProcessLpAwardRecordInfoScan(record))
		}
	}

	lpAwardRecords = &_lpAwardRecords
	return lpAwardRecords, nil
}
