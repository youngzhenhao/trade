package pool

import (
	"errors"
	"math/big"
	"strconv"
	"trade/utils"
)

func sortTokens(tokenA string, tokenB string) (token0 string, token1 string, err error) {
	if tokenA == tokenB {
		err = errors.New("identicalTokens(" + tokenA + ")")
		return "", "", err
	}
	if !(len(tokenA) == len(TokenSatTag) || len(tokenA) == AssetIdLength) {
		err = errors.New("invalid tokenA length(" + strconv.Itoa(len(tokenA)) + ")")
		return "", "", err
	}
	if !(len(tokenA) == len(TokenSatTag) || len(tokenA) == AssetIdLength) {
		err = errors.New("invalid tokenB length(" + strconv.Itoa(len(tokenB)) + ")")
		return "", "", err
	}
	// @dev: sat is always token0
	if tokenA == TokenSatTag {
		token0, token1 = tokenA, tokenB
	} else if tokenB == TokenSatTag {
		token0, token1 = tokenB, tokenA
	} else if tokenA < tokenB {
		token0, token1 = tokenA, tokenB
	} else {
		token0, token1 = tokenB, tokenA
	}
	if token0 == "" {
		err = errors.New("zeroTokens(" + token0 + ")")
		return "", "", err
	}
	return token0, token1, nil
}

// quote
// given some amount of an asset and pair reserves, returns an equivalent amount of the other asset
func quote(amountA string, reserveA string, reserveB string) (amountB string, err error) {
	_amountA, success := new(big.Int).SetString(amountA, 10)
	if !success {
		return "", errors.New("amountA SetString(" + amountA + ") " + strconv.FormatBool(success))
	}
	if _amountA.Sign() < 0 {
		err = errors.New("insufficientAmount(" + _amountA.String() + ")")
		return "", err
	}
	_reserveA, success := new(big.Int).SetString(reserveA, 10)
	if !success {
		return "", errors.New("reserveA SetString(" + reserveA + ") " + strconv.FormatBool(success))
	}
	_reserveB, success := new(big.Int).SetString(reserveB, 10)
	if !success {
		return "", errors.New("reserveB SetString(" + reserveB + ") " + strconv.FormatBool(success))
	}
	if !((_reserveA.Sign() > 0) && (_reserveB.Sign() > 0)) {
		err = errors.New("insufficientLiquidity(" + _reserveA.String() + "," + _reserveB.String() + ")")
		return "", err
	}
	_amountB := new(big.Int).Div(new(big.Int).Mul(_amountA, _reserveB), _reserveA)
	return _amountB.String(), nil
}

// getAmountOut
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
// @dev: e.g. fee rate is 3/1000 (0.3%)
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
		return "", errors.New("amountIn SetString(" + amountIn + ") " + strconv.FormatBool(success))
	}
	if !(_amountIn.Sign() > 0) {
		err = errors.New("insufficientInputAmount(" + _amountIn.String() + ")")
		return "", err
	}
	_reserveIn, success := new(big.Int).SetString(reserveIn, 10)
	if !success {
		return "", errors.New("reserveIn SetString(" + reserveIn + ") " + strconv.FormatBool(success))
	}
	_reserveOut, success := new(big.Int).SetString(reserveOut, 10)
	if !success {
		return "", errors.New("reserveOut SetString(" + reserveOut + ") " + strconv.FormatBool(success))
	}
	if !((_reserveIn.Sign() > 0) && (_reserveOut.Sign() > 0)) {
		err = errors.New("insufficientLiquidity(" + _reserveIn.String() + "," + _reserveOut.String() + ")")
		return "", err
	}
	k := new(big.Int).SetUint64(uint64(feeK))
	oneThousand := new(big.Int).SetUint64(1000)
	// @dev: dx(1000 - k)
	amountInWithFee := new(big.Int).Mul(_amountIn, new(big.Int).Sub(oneThousand, k))
	// @dev: numerator dx(1000 - k)y_0
	numerator := new(big.Int).Mul(_reserveOut, amountInWithFee)
	// @dev: denominator 1000x_0 + dx(1000 - k)
	denominator := new(big.Int).Add(new(big.Int).Mul(_reserveIn, oneThousand), amountInWithFee)
	// @dev: dy = numerator / denominator
	_amountOut := new(big.Int).Div(numerator, denominator)
	amountOut = _amountOut.String()
	return amountOut, nil
}

// getAmountIn
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
// @dev: e.g. fee rate is 3/1000 (0.3%)
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
		return "", errors.New("amountOut SetString(" + amountOut + ") " + strconv.FormatBool(success))
	}
	if !(_amountOut.Sign() > 0) {
		err = errors.New("insufficientOutputAmount(" + _amountOut.String() + ")")
		return "", err
	}
	_reserveIn, success := new(big.Int).SetString(reserveIn, 10)
	if !success {
		return "", errors.New("reserveIn SetString(" + reserveIn + ") " + strconv.FormatBool(success))
	}
	_reserveOut, success := new(big.Int).SetString(reserveOut, 10)
	if !success {
		return "", errors.New("reserveOut SetString(" + reserveOut + ") " + strconv.FormatBool(success))
	}
	if !((_reserveIn.Sign() > 0) && (_reserveOut.Sign() > 0)) {
		err = errors.New("insufficientLiquidity(" + _reserveIn.String() + "," + _reserveOut.String() + ")")
		return "", err
	}
	k := new(big.Int).SetUint64(uint64(feeK))
	oneThousand := new(big.Int).SetUint64(1000)
	// @dev: numerator 1000x_0dy
	numerator := new(big.Int).Mul(new(big.Int).Mul(_reserveIn, _amountOut), oneThousand)
	// @dev: denominator (y_0 - dy)(1000 - k)
	denominator := new(big.Int).Mul(new(big.Int).Sub(_reserveOut, _amountOut), new(big.Int).Sub(oneThousand, k))
	// @dev: Addition of 1 is to compensate for the loss of precision that may result from integer division
	one := new(big.Int).SetUint64(1)
	m := new(big.Int)
	// @dev: dy, mod = numerator div mod denominator
	_amountIn, m := new(big.Int).DivMod(numerator, denominator, m)
	if m.Sign() != 0 {
		// @dev: dy = (numerator / denominator) + 1
		_amountIn = new(big.Int).Add(_amountIn, one)
	}
	amountIn = _amountIn.String()
	return amountIn, nil
}

func getLiquidityK(reserve0 string, reserve1 string) (k string, err error) {
	_reserve0, success := new(big.Int).SetString(reserve0, 10)
	if !success {
		return "", errors.New("reserve0 SetString(" + reserve0 + ") " + strconv.FormatBool(success))
	}
	_reserve1, success := new(big.Int).SetString(reserve1, 10)
	if !success {
		return "", errors.New("reserve0 SetString(" + reserve1 + ") " + strconv.FormatBool(success))
	}
	if !((_reserve0.Sign() > 0) && (_reserve1.Sign() > 0)) {
		err = errors.New("insufficientLiquidity(" + _reserve0.String() + "," + _reserve1.String() + ")")
		return "", err
	}
	_k := new(big.Int).Mul(_reserve0, _reserve1)
	return _k.String(), nil
}

func quoteBig(_amountA *big.Int, _reserveA *big.Int, _reserveB *big.Int) (_amountB *big.Int, err error) {
	if _amountA.Sign() < 0 {
		err = errors.New("insufficientAmount(" + _amountA.String() + ")")
		return new(big.Int), err
	}
	if !((_reserveA.Sign() > 0) && (_reserveB.Sign() > 0)) {
		err = errors.New("insufficientLiquidity(" + _reserveA.String() + "," + _reserveB.String() + ")")
		return new(big.Int), err
	}
	_amountB = new(big.Int).Div(new(big.Int).Mul(_amountA, _reserveB), _reserveA)
	return _amountB, nil
}

func getAmountOutBig(_amountIn *big.Int, _reserveIn *big.Int, _reserveOut *big.Int, feeK uint16) (_amountOut *big.Int, err error) {
	if !(feeK <= 1000) {
		err = errors.New("invalid fee rate k(" + strconv.FormatUint(uint64(feeK), 10) + "), must less equal than 1000")
		return new(big.Int), err
	}
	if !(_amountIn.Sign() > 0) {
		err = errors.New("insufficientInputAmount(" + _amountIn.String() + ")")
		return new(big.Int), err
	}
	if !((_reserveIn.Sign() > 0) && (_reserveOut.Sign() > 0)) {
		err = errors.New("insufficientLiquidity(" + _reserveIn.String() + "," + _reserveOut.String() + ")")
		return new(big.Int), err
	}
	k := new(big.Int).SetUint64(uint64(feeK))
	oneThousand := new(big.Int).SetUint64(1000)
	// @dev: dx(1000 - k)
	amountInWithFee := new(big.Int).Mul(_amountIn, new(big.Int).Sub(oneThousand, k))
	// @dev: numerator dx(1000 - k)y_0
	numerator := new(big.Int).Mul(_reserveOut, amountInWithFee)
	// @dev: denominator 1000x_0 + dx(1000 - k)
	denominator := new(big.Int).Add(new(big.Int).Mul(_reserveIn, oneThousand), amountInWithFee)
	// @dev: dy = numerator / denominator
	_amountOut = new(big.Int).Div(numerator, denominator)
	return _amountOut, nil
}

func getAmountInBig(_amountOut *big.Int, _reserveIn *big.Int, _reserveOut *big.Int, feeK uint16) (_amountIn *big.Int, err error) {
	if !(feeK <= 1000) {
		err = errors.New("invalid fee rate k(" + strconv.FormatUint(uint64(feeK), 10) + "), must less equal than 1000")
		return new(big.Int), err
	}
	if !(_amountOut.Sign() > 0) {
		err = errors.New("insufficientOutputAmount(" + _amountOut.String() + ")")
		return new(big.Int), err
	}
	if !((_reserveIn.Sign() > 0) && (_reserveOut.Sign() > 0)) {
		err = errors.New("insufficientLiquidity(" + _reserveIn.String() + "," + _reserveOut.String() + ")")
		return new(big.Int), err
	}
	k := new(big.Int).SetUint64(uint64(feeK))
	oneThousand := new(big.Int).SetUint64(1000)
	// @dev: numerator 1000x_0dy
	numerator := new(big.Int).Mul(new(big.Int).Mul(_reserveIn, _amountOut), oneThousand)
	// @dev: denominator (y_0 - dy)(1000 - k)
	denominator := new(big.Int).Mul(new(big.Int).Sub(_reserveOut, _amountOut), new(big.Int).Sub(oneThousand, k))
	// @dev: Addition of 1 is to compensate for the loss of precision that may result from integer division
	one := new(big.Int).SetUint64(1)
	m := new(big.Int)
	// @dev: dy, mod = numerator div mod denominator
	_amountIn, m = new(big.Int).DivMod(numerator, denominator, m)
	if m.Sign() != 0 {
		// @dev: dy = (numerator / denominator) + 1
		_amountIn = new(big.Int).Add(_amountIn, one)
	}
	return _amountIn, nil
}

func getLiquidityKBig(_reserve0 *big.Int, _reserve1 *big.Int) (_k *big.Int, err error) {
	if !((_reserve0.Sign() > 0) && (_reserve1.Sign() > 0)) {
		err = errors.New("insufficientLiquidity(" + _reserve0.String() + "," + _reserve1.String() + ")")
		return new(big.Int), err
	}
	_k = new(big.Int).Mul(_reserve0, _reserve1)
	return _k, nil
}

func minBigInt(a *big.Int, b *big.Int) *big.Int {
	if a.Cmp(b) < 0 {
		return a
	}
	return b
}

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
