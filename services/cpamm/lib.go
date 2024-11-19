package cpamm

import (
	"errors"
	"math/big"
	"strconv"
	"trade/utils"
)

/**
 * Reference
 * https://github.com/Uniswap/v2-periphery/blob/master/contracts/libraries/UniswapV2Library.sol
 */

const (
	ProjectParty = 3
	// TODO: All lps award 0.3% fee
	LpAward = 0
	// FeeK is 1000 * f, i.e. f = k / 1000
	// fee rate should be 3/1000 (0.3%), and f = k / 1000
	FeeK uint16 = ProjectParty + LpAward
)

// TODO: TOKEN0 sat
// @dev: Tested
func sortTokens(tokenA string, tokenB string) (token0 string, token1 string, err error) {
	if tokenA == tokenB {
		err = errors.New("lib:IDENTICAL_TOKENS(" + tokenA + ")")
		return "", "", err
	}
	if !(len(tokenA) == 3 || len(tokenA) == 64) {
		err = errors.New("invalid tokenA length(" + strconv.Itoa(len(tokenA)) + ")")
		return "", "", err
	}
	if !(len(tokenA) == 3 || len(tokenA) == 64) {
		err = errors.New("invalid tokenB length(" + strconv.Itoa(len(tokenB)) + ")")
		return "", "", err
	}
	// @dev: sat is always token0
	if tokenA == "sat" {
		token0, token1 = tokenA, tokenB
	} else if tokenB == "sat" {
		token0, token1 = tokenB, tokenA
	} else if tokenA < tokenB {
		token0, token1 = tokenA, tokenB
	} else {
		token0, token1 = tokenB, tokenA
	}
	if token0 == "" {
		err = errors.New("lib:ZERO_TOKENS(" + token0 + ")")
		return "", "", err
	}
	return token0, token1, nil
}

// calculates the CREATE2 address for a pair without making any external calls
// TODO: need to complete and test
func pairFor(tokenA string, tokenB string) (pair *PoolPair, err error) {
	//TODO:
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return nil, err
	}
	_ = token0
	_ = token1
	pair = &PoolPair{}

	return pair, nil
}

// fetches and sorts the reserves for a pair
// TODO: need to test
func getReserves(tokenA string, tokenB string) (reserveA string, reserveB string, err error) {
	// @dev: Get smaller token
	token0, _, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return "", "", utils.AppendErrorInfo(err, "sortTokens("+tokenA+","+tokenB+")")
	}
	// @dev: Get pair
	pair, err := pairFor(tokenA, tokenB)
	if err != nil {
		return "", "", utils.AppendErrorInfo(err, "pairFor("+tokenA+","+tokenB+")")
	}
	// @dev: Get pair's reserves
	reserve0, reserve1 := pair.getReserves()
	if tokenA == token0 {
		reserveA, reserveB = reserve0, reserve1
	} else {
		reserveA, reserveB = reserve1, reserve0
	}
	return reserveA, reserveB, nil
}

// @dev: Tested
// given some amount of an asset and pair reserves, returns an equivalent amount of the other asset
func quote(amountA string, reserveA string, reserveB string) (amountB string, err error) {
	_amountA, success := new(big.Int).SetString(amountA, 10)
	if !success {
		return "", utils.AppendErrorInfo(err, "SetString("+amountA+") "+strconv.FormatBool(success))
	}
	if _amountA.Sign() < 0 {
		err = errors.New("lib:INSUFFICIENT_AMOUNT(" + _amountA.String() + ")")
		return "", err
	}
	_reserveA, success := new(big.Int).SetString(reserveA, 10)
	if !success {
		return "", utils.AppendErrorInfo(err, "SetString("+reserveA+") "+strconv.FormatBool(success))
	}
	_reserveB, success := new(big.Int).SetString(reserveB, 10)
	if !success {
		return "", utils.AppendErrorInfo(err, "SetString("+reserveB+") "+strconv.FormatBool(success))
	}
	if !((_reserveA.Sign() > 0) && (_reserveB.Sign() > 0)) {
		err = errors.New("lib:INSUFFICIENT_LIQUIDITY(" + _reserveA.String() + "," + _reserveB.String() + ")")
		return "", err
	}
	_amountB := _amountA.Mul(_amountA, _reserveB).Div(_amountA, _reserveA)
	return _amountB.String(), nil
}

// getAmountOut
// @dev: Tested
// @Description:
//
//	x_0y_0=(x_0 + dx)(y_0 - dy)
//
//	x_0:	reserveIn
//	y_0:	reserveOut
//	dx:		amountIn
//	dy:		amountOut
//
//	dx with fee: dx(1 - f)
//
//			dx(1 - f)y_0
//	dy = —————————————————————
//			x_0 + dx(1 - f)
//
// ========================================
//
//	define: 	f = k / 1000
//
//					 k
//			dx(1 - ——————)y_0
//					1000
//	dy = ——————————————————————————
//						   k
//			x_0 + dx(1 - ——————)
//						  1000
//
// ========================================
//
//			dx(1000 - k)y_0
//	dy = ——————————————————————————
//			1000x_0 + dx(1000 - k)
//
// ========================================
//
// @dev: fee rate should be 3/1000 (0.3%)
//
//			dx(997)y_0
//	dy = —————————————————————— (k = 3)
//			1000x_0 + dx(997)
//
// given an input amount of an asset and pair reserves, returns the maximum output amount of the other asset
func getAmountOut(amountIn string, reserveIn string, reserveOut string, feeK uint16) (amountOut string, err error) {
	if !(feeK <= 1000) {
		err = errors.New("invalid fee rate k(" + strconv.FormatUint(uint64(feeK), 10) + "), must less equal than 1000")
		return "", err
	}
	_amountIn, success := new(big.Int).SetString(amountIn, 10)
	if !success {
		return "", utils.AppendErrorInfo(err, "SetString("+amountIn+") "+strconv.FormatBool(success))
	}
	if !(_amountIn.Sign() > 0) {
		err = errors.New("lib:INSUFFICIENT_INPUT_AMOUNT(" + _amountIn.String() + ")")
		return "", err
	}
	_reserveIn, success := new(big.Int).SetString(reserveIn, 10)
	if !success {
		return "", utils.AppendErrorInfo(err, "SetString("+reserveIn+") "+strconv.FormatBool(success))
	}
	_reserveOut, success := new(big.Int).SetString(reserveOut, 10)
	if !success {
		return "", utils.AppendErrorInfo(err, "SetString("+reserveOut+") "+strconv.FormatBool(success))
	}
	if !((_reserveIn.Sign() > 0) && (_reserveOut.Sign() > 0)) {
		err = errors.New("lib:INSUFFICIENT_LIQUIDITY(" + _reserveIn.String() + "," + _reserveOut.String() + ")")
		return "", err
	}
	k := new(big.Int).SetUint64(uint64(feeK))
	oneThousand := new(big.Int).SetUint64(1000)
	// @dev: dx(1000 - k)
	amountInWithFee := _amountIn.Mul(_amountIn, new(big.Int).Sub(oneThousand, k))
	// @dev: numerator dx(1000 - k)y_0
	numerator := _reserveOut.Mul(_reserveOut, amountInWithFee)
	// @dev: denominator 1000x_0 + dx(1000 - k)
	denominator := _reserveIn.Mul(_reserveIn, oneThousand).Add(_reserveIn, amountInWithFee)
	// @dev: dy = numerator / denominator
	_amountOut := numerator.Div(numerator, denominator)
	amountOut = _amountOut.String()
	return amountOut, nil
}

// getAmountIn
// @dev: Tested
// @Description:
//
//	x_0y_0=(x_0 + dx)(y_0 - dy)
//
//	x_0:	reserveIn
//	y_0:	reserveOut
//	dx:		amountIn
//	dy:		amountOut
//
//	dx with fee: dx(1 - f)
//
//			x_0dy		   1
//	dx = ————————————— —————————
//			y_0 - dy	 1 - f
//
// ========================================
//
//	define: 	f = k / 1000
//
//				x_0dy
//	dx = ——————————————————————————
//						   	 k
//			(y_0 - dy)(1 - ——————)
//						  	1000
//
// ========================================
//
//				1000x_0dy
//	dx = ——————————————————————————
//			(y_0 - dy)(1000 - k)
//
// ========================================
//
// @dev: fee rate should be 3/1000 (0.3%)
//
//				1000x_0dy
//	dx = —————————————————————————— (k = 3)
//			(y_0 - dy)(997)
//
// given an output amount of an asset and pair reserves, returns a required input amount of the other asset
func getAmountIn(amountOut string, reserveIn string, reserveOut string, feeK uint16) (amountIn string, err error) {
	if !(feeK <= 1000) {
		err = errors.New("invalid fee rate k(" + strconv.FormatUint(uint64(feeK), 10) + "), must less equal than 1000")
		return "", err
	}
	_amountOut, success := new(big.Int).SetString(amountOut, 10)
	if !success {
		return "", utils.AppendErrorInfo(err, "SetString("+amountOut+") "+strconv.FormatBool(success))
	}
	if !(_amountOut.Sign() > 0) {
		err = errors.New("lib:INSUFFICIENT_OUTPUT_AMOUNT(" + _amountOut.String() + ")")
		return "", err
	}
	_reserveIn, success := new(big.Int).SetString(reserveIn, 10)
	if !success {
		return "", utils.AppendErrorInfo(err, "SetString("+reserveIn+") "+strconv.FormatBool(success))
	}
	_reserveOut, success := new(big.Int).SetString(reserveOut, 10)
	if !success {
		return "", utils.AppendErrorInfo(err, "SetString("+reserveOut+") "+strconv.FormatBool(success))
	}
	if !((_reserveIn.Sign() > 0) && (_reserveOut.Sign() > 0)) {
		err = errors.New("lib:INSUFFICIENT_LIQUIDITY(" + _reserveIn.String() + "," + _reserveOut.String() + ")")
		return "", err
	}
	k := new(big.Int).SetUint64(uint64(feeK))
	oneThousand := new(big.Int).SetUint64(1000)
	// @dev: numerator 1000x_0dy
	numerator := _reserveIn.Mul(_reserveIn, _amountOut).Mul(_reserveIn, oneThousand)
	// @dev: denominator (y_0 - dy)(1000 - k)
	denominator := _reserveOut.Sub(_reserveOut, _amountOut).Mul(_reserveOut, new(big.Int).Sub(oneThousand, k))
	// @dev: Addition of 1 is to compensate for the loss of precision that may result from integer division
	one := new(big.Int).SetUint64(1)
	m := new(big.Int)
	// @dev: dy, mod = numerator div mod denominator
	_amountIn, m := numerator.DivMod(numerator, denominator, m)
	if m.Sign() != 0 {
		// @dev: dy = (numerator / denominator) + 1
		_amountIn = _amountIn.Add(_amountIn, one)
	}
	amountIn = _amountIn.String()
	return amountIn, nil
}

// TODO: need to complete and test
// performs chained getAmountOut calculations on any number of pairs
func getAmountsOut(amountIn string, path []string) (amounts []string, err error) {
	if !(len(path) >= 2) {
		err = errors.New("lib:INVALID_PATH(" + strconv.Itoa(len(path)) + ")")
	}
	amounts = make([]string, len(path), len(path))
	amounts[0] = amountIn
	for i := 0; i < len(path)-1; i++ {
		// @dev: Get reserves x_0 and y_0
		var reserveIn, reserveOut string
		reserveIn, reserveOut, err = getReserves(path[i], path[i+1])
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "getReserves("+path[i]+","+path[i+1]+")")
		}
		// @dev: Get Amount Out
		amounts[i+1], err = getAmountOut(amounts[i], reserveIn, reserveOut, FeeK)
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "getAmountOut("+amounts[i]+","+reserveIn+","+reserveOut+")")
		}
	}
	return amounts, nil
}

// TODO: need to complete and test
// performs chained getAmountIn calculations on any number of pairs
func getAmountsIn(amountOut string, path []string) (amounts []string, err error) {
	if !(len(path) >= 2) {
		err = errors.New("lib:INVALID_PATH(length:" + strconv.Itoa(len(path)) + ")")
	}
	amounts = make([]string, len(path), len(path))
	amounts[len(amounts)-1] = amountOut
	for i := len(path) - 1; i > 0; i-- {
		// @dev: Get reserves x_0 and y_0
		var reserveIn, reserveOut string
		reserveIn, reserveOut, err = getReserves(path[i-1], path[i])
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "getReserves("+path[i-1]+","+path[i]+")")
		}
		// @dev: Get Amount In
		amounts[i-1], err = getAmountIn(amounts[i], reserveIn, reserveOut, FeeK)
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "getAmountIn("+amounts[i]+","+reserveIn+","+reserveOut+")")
		}
	}
	return amounts, nil
}

// TODO: LiquidityMathLib

// TODO: OracleLib
