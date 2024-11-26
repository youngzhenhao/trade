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
	PairId    uint `json:"pair_id" gorm:"uniqueIndex"`
	AccountId uint `json:"account_id" gorm:"index"`
}

var Lock map[string]map[string]*sync.Mutex

func _addLiquidity(_amount0Desired *big.Int, _amount1Desired *big.Int, _amount0Min *big.Int, _amount1Min *big.Int, _reserve0 *big.Int, _reserve1 *big.Int) (_amount0 *big.Int, _amount1 *big.Int, err error) {
	_amount1Optimal, err := quoteBig(_amount0Desired, _reserve0, _reserve1)
	if err != nil {
		return new(big.Int), new(big.Int), utils.AppendErrorInfo(err, "quoteBig(_amount0Desired, _reserve0, _reserve1)")
	}
	if _amount1Optimal.Cmp(_amount1Desired) <= 0 {
		if !(_amount1Optimal.Cmp(_amount1Min) >= 0) {
			return new(big.Int), new(big.Int), errors.New("insufficientAmount1(" + _amount1Optimal.String() + ")")
		}
		_amount0, _amount1 = _amount0Desired, _amount1Optimal
	} else {
		_amount0Optimal, err := quoteBig(_amount1Desired, _reserve1, _reserve0)
		if err != nil {
			return new(big.Int), new(big.Int), utils.AppendErrorInfo(err, "quoteBig(_amount1Desired, _reserve1, _reserve0)")
		}
		if !(_amount0Optimal.Cmp(_amount0Desired) <= 0) {
			return new(big.Int), new(big.Int), errors.New("amount0Optimal(" + _amount0Optimal.String() + ") is greater than amount0Desired(" + _amount0Desired.String() + ")")
		}
		if !(_amount0Optimal.Cmp(_amount0Min) >= 0) {
			return new(big.Int), new(big.Int), errors.New("insufficientAmount0(" + _amount0Optimal.String() + ")")
		}
		_amount0, _amount1 = _amount0Optimal, _amount1Desired
	}
	return _amount0, _amount1, nil
}

// TODO: Add Liquidity
func AddLiquidity(tokenA string, tokenB string, amountADesired string, amountBDesired string, amountAMin string, amountBMin string) error {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return utils.AppendErrorInfo(err, "sortTokens")
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
		return utils.AppendErrorInfo(err, "SetString("+amount0Desired+") "+strconv.FormatBool(success))
	}
	// amount1Desired
	_amount1Desired, success := new(big.Int).SetString(amount1Desired, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+amount1Desired+") "+strconv.FormatBool(success))
	}
	// amount0Min
	_amount0Min, success := new(big.Int).SetString(amount0Min, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+amount0Min+") "+strconv.FormatBool(success))
	}
	if _amount0Min.Sign() < 0 {
		return errors.New("amount0Min(" + _amount0Min.String() + ") is negative")
	}
	if _amount0Min.Cmp(_amount0Desired) > 0 {
		return errors.New("amount0Min(" + _amount0Min.String() + ") is greater than amount0Desired(" + _amount0Desired.String() + ")")
	}
	// amount1Min
	_amount1Min, success := new(big.Int).SetString(amount1Min, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+amount1Min+") "+strconv.FormatBool(success))
	}
	if _amount1Min.Sign() < 0 {
		return errors.New("amount1Min(" + _amount1Min.String() + ") is negative")
	}
	if _amount1Min.Cmp(_amount1Desired) > 0 {
		return errors.New("amount1Min(" + _amount1Min.String() + ") is greater than amount1Desired(" + _amount1Desired.String() + ")")
	}
	// @dev: lock
	Lock[token0][token1].Lock()
	defer Lock[token0][token1].Unlock()
	// @dev: transaction
	tx := middleware.DB.Begin()
	// @dev: get pair
	var _pair Pair
	err = tx.Where("token0 = ? AND token1 = ?", token0, token1).First(&_pair).Error
	// @dev: pair does not exist
	if err != nil {
		// @dev: create pair
		newPair, err := NewPair(token0, token1, amount0Desired, amount1Desired)
		if err != nil {
			tx.Rollback()
			return utils.AppendErrorInfo(err, "newPair")
		}
		err = tx.Create(&newPair).Error
		if err != nil {
			tx.Rollback()
			return utils.AppendErrorInfo(err, "create pair")
		}
	} else {
		// @dev: pair exists
		// reserve0
		_reserve0, success := new(big.Int).SetString(_pair.Reserve0, 10)
		if !success {
			return utils.AppendErrorInfo(err, "SetString("+_pair.Reserve0+") "+strconv.FormatBool(success))
		}
		// reserve1
		_reserve1, success := new(big.Int).SetString(_pair.Reserve1, 10)
		if !success {
			return utils.AppendErrorInfo(err, "SetString("+_pair.Reserve1+") "+strconv.FormatBool(success))
		}
		// No fee for adding liquidity
		_amount0, _amount1, err := _addLiquidity(_amount0Desired, _amount1Desired, _amount0Min, _amount1Min, _reserve0, _reserve1)
		if err != nil {
			tx.Rollback()
			return utils.AppendErrorInfo(err, "_addLiquidity")
		}
		_, _ = _amount0, _amount1

		//	TODO: mint liquidity
		//address pair = UniswapV2Library.pairFor(factory, tokenA, tokenB);

		//TransferHelper.safeTransferFrom(tokenA, msg.sender, pair, amountA);
		// TODO: Transfer _amount0 of token0 from user to pool

		//TransferHelper.safeTransferFrom(tokenB, msg.sender, pair, amountB);
		// TODO: Transfer _amount1 of token1 from user to pool

		//liquidity = IUniswapV2Pair(pair).mint(to);
		//	TODO: mint share
		// @dev: calculate liquidity
		
		//	TODO: update pair, share, shareBalance, shareRecord

	}

	return tx.Commit().Error
}

// TODO: Remove Liquidity
func RemoveLiquidity() {

}

// TODO: Swap In
func SwapIn() {

}

// TODO: Swap Out
func SwapOut() {

}
