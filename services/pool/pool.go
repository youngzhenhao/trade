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
// @Description: Add Liquidity, create pair and share if not exist
func AddLiquidity(tokenA string, tokenB string, amountADesired string, amountBDesired string, amountAMin string, amountBMin string, username string) (amountA string, amountB string, liquidity string, err error) {
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
		_liquidity, err = _mintBig(_amount0, _amount1, big.NewInt(0), big.NewInt(0), big.NewInt(0), isTokenZeroSat)
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
		_liquidity, err = _mintBig(_amount0, _amount1, _reserve0, _reserve1, _totalSupply, isTokenZeroSat)
		if err != nil {
			return ZeroValue, ZeroValue, ZeroValue, utils.AppendErrorInfo(err, "_mintBig")
		}
		// @dev: update share, shareBalance, shareRecord
		newSupply := new(big.Int).Add(_totalSupply, _liquidity)
		//fmt.Printf("newSupply: %v;(_totalSupply: %v + _liquidity: %v)\n", newSupply, _totalSupply, _liquidity)
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

// RemoveLiquidity
// @Description: Remove Liquidity, pair and share must exist
func RemoveLiquidity(tokenA string, tokenB string, liquidity string, amountAMin string, amountBMin string, username string, feeK uint16) (amountA string, amountB string, err error) {
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

	// @dev: check liquidity
	var shareBalance string
	err = tx.Model(&ShareBalance{}).Select("balance").
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

	// TODO: check balance then transfer, record transfer Id
	//_safeTransfer(_token0, to, amount0);
	// TODO: transfer _amount0 of token0 from pool to user
	//_safeTransfer(_token1, to, amount1);
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
// TODO: Test
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

	// get pair
	var _pair Pair
	err = tx.Model(&Pair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
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
			_price, err = GetToken1PriceBig(_reserve0, _reserve1)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "GetToken1PriceBig")
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

	// TODO: check balance then transfer, record transfer Id

	// TODO: Transfer _amountIn of tokenIn from user to pool

	// TODO: Transfer _amountOut of tokenOut from pool to user

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

	err = tx.Model(&Pair{}).Where("token0 = ? AND token1 = ?", token0, token1).
		Updates(map[string]any{
			"reserve0": _newReserve0.String(),
			"reserve1": _newReserve1.String(),
		}).Error
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "update pair")
	}

	// @dev: update swapRecord
	var recordId uint
	recordId, err = CreateSwapRecord(tx, pairId, username, tokenIn, tokenOut, amountIn, _amountOut.String(), _reserveIn.String(), _reserveOut.String(), _swapFeeFloat.String(), swapFeeType, SwapExactTokenNoPath)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "CreateSwapRecord")
	}

	// == COPY ==
	// get share
	var share Share
	var shareId uint
	err = tx.Model(&Share{}).Where("pair_id = ?", pairId).First(&share).Error
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

	err = tx.Model(&ShareBalance{}).
		Select("username, balance").
		Where("share_id = ?", shareId).
		Scan(&userAndShares).Error

	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "get userAndShares")
	}

	for _, _userAndShare := range userAndShares {
		_balanceFloat, success := new(big.Float).SetString(_userAndShare.Balance)
		if !success {
			return ZeroValue, errors.New(_userAndShare.Username + " balance SetString(" + _userAndShare.Balance + ") " + strconv.FormatBool(success))
		}

		var _awardFloat = big.NewFloat(0)

		LpAwardFeeKFloat := new(big.Float).SetUint64(uint64(lpAwardFeeK))
		SwapFeeKFloat := new(big.Float).SetUint64(uint64(feeK))
		_swapFeeForAwardFloat := new(big.Float).Quo(new(big.Float).Mul(_swapFeeFloat, LpAwardFeeKFloat), SwapFeeKFloat)

		_awardFloat = new(big.Float).Quo(new(big.Float).Mul(_swapFeeForAwardFloat, _balanceFloat), _totalSupplyFloat)
		err = UpdateLpAwardBalanceAndRecordSwap(tx, shareId, _userAndShare.Username, _awardFloat, _swapFeeFloat.String(), _userAndShare.Balance, _totalSupplyFloat.String(), recordId)

		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "UpdateLpAwardBalanceAndRecordSwap")
		}
	}
	// == COPY END ==

	amountOut = _amountOut.String()
	err = nil
	return amountOut, err
}

// TODO: Swap Tokens For Exact Tokens
// TODO: Test
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

	// get pair
	var _pair Pair
	err = tx.Model(&Pair{}).Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
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

			_price, err = GetToken1PriceBig(_reserve0, _reserve1)
			if err != nil {
				return ZeroValue, utils.AppendErrorInfo(err, "GetToken1PriceBig")
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

	// TODO: check balance then transfer, record transfer Id

	// TODO: Transfer _amountIn of tokenIn from user to pool

	// TODO: Transfer _amountOut of tokenOut from pool to user

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

	err = tx.Model(&Pair{}).Where("token0 = ? AND token1 = ?", token0, token1).
		Updates(map[string]any{
			"reserve0": _newReserve0.String(),
			"reserve1": _newReserve1.String(),
		}).Error
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "update pair")
	}

	// @dev: update swapRecord
	var recordId uint
	recordId, err = CreateSwapRecord(tx, pairId, username, tokenIn, tokenOut, _amountIn.String(), amountOut, _reserveIn.String(), _reserveOut.String(), _swapFeeFloat.String(), swapFeeType, SwapForExactTokenNoPath)
	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "CreateSwapRecord")
	}

	// get share
	var share Share
	var shareId uint
	err = tx.Model(&Share{}).Where("pair_id = ?", pairId).First(&share).Error
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

	err = tx.Model(&ShareBalance{}).
		Select("username, balance").
		Where("share_id = ?", shareId).
		Scan(&userAndShares).Error

	if err != nil {
		return ZeroValue, utils.AppendErrorInfo(err, "get userAndShares")
	}

	for _, _userAndShare := range userAndShares {
		_balanceFloat, success := new(big.Float).SetString(_userAndShare.Balance)
		if !success {
			return ZeroValue, errors.New(_userAndShare.Username + " balance SetString(" + _userAndShare.Balance + ") " + strconv.FormatBool(success))
		}

		var _awardFloat = big.NewFloat(0)

		LpAwardFeeKFloat := new(big.Float).SetUint64(uint64(lpAwardFeeK))
		SwapFeeKFloat := new(big.Float).SetUint64(uint64(feeK))
		_swapFeeForAwardFloat := new(big.Float).Quo(new(big.Float).Mul(_swapFeeFloat, LpAwardFeeKFloat), SwapFeeKFloat)

		_awardFloat = new(big.Float).Quo(new(big.Float).Mul(_swapFeeForAwardFloat, _balanceFloat), _totalSupplyFloat)
		err = UpdateLpAwardBalanceAndRecordSwap(tx, shareId, _userAndShare.Username, _awardFloat, _swapFeeFloat.String(), _userAndShare.Balance, _totalSupplyFloat.String(), recordId)

		if err != nil {
			return ZeroValue, utils.AppendErrorInfo(err, "UpdateLpAwardBalanceAndRecordSwap")
		}
	}

	amountIn = _amountIn.String()
	err = nil
	return amountIn, err
}

// TODO: Tolerance

// TODO: Encapsulate API

// TODO: 1. Award Query
// TODO: 2. Award Withdraw
