package pool

import (
	"errors"
	"gorm.io/gorm"
	"math/big"
	"strconv"
	"trade/middleware"
	"trade/utils"
)

type PoolPair struct {
	gorm.Model
	IsTokenZeroSat bool   `json:"is_token_zero_sat" gorm:"index"`
	Token0         string `json:"token0" gorm:"type:varchar(255);uniqueIndex:idx_token_0_token_1"`
	Token1         string `json:"token1" gorm:"type:varchar(255);uniqueIndex:idx_token_0_token_1"`
	Reserve0       string `json:"reserve0" gorm:"type:varchar(255)"`
	Reserve1       string `json:"reserve1" gorm:"type:varchar(255)"`
}

func getPair(token0 string, token1 string) (pair *PoolPair, err error) {
	_token0, _token1, err := sortTokens(token0, token1)
	if err != nil {
		return new(PoolPair), utils.AppendErrorInfo(err, "sortTokens")
	}
	var _pair PoolPair
	err = middleware.DB.Where("token0 = ? AND token1 = ?", _token0, _token1).First(&_pair).Error
	return &_pair, err
}

func createPair(token0 string, token1 string, reserve0 string, reserve1 string) (err error) {
	// sort token
	_token0, _token1, err := sortTokens(token0, token1)
	if err != nil {
		return utils.AppendErrorInfo(err, "sortTokens")
	}
	// check reserves
	_reserve0, success := new(big.Int).SetString(reserve0, 10)
	if !success {
		return errors.New("reserve0 SetString(" + reserve0 + ") " + strconv.FormatBool(success))
	}
	_reserve1, success := new(big.Int).SetString(reserve1, 10)
	if !success {
		return errors.New("reserve1 SetString(" + reserve1 + ") " + strconv.FormatBool(success))
	}
	if !((_reserve0.Sign() > 0) && (_reserve1.Sign() > 0)) {
		err = errors.New("insufficientLiquidity(" + _reserve0.String() + "," + _reserve1.String() + ")")
		return err
	}
	// cmp with minimum liquidity sat
	if _token0 == TokenSatTag {
		_minLiquiditySat := new(big.Int).SetUint64(uint64(MinAddLiquiditySat))
		if _reserve0.Cmp(_minLiquiditySat) < 0 {
			err = errors.New("insufficientLiquidity Sat(" + _reserve0.String() + "), need " + _minLiquiditySat.String())
			return err
		}
	} else {
		// _token0 != TokenSatTag
		// TODO: Add MinLiquidity check for pair whose token0 is not sat
		_liquidity := new(big.Int).Sqrt(new(big.Int).Mul(_reserve0, _reserve1))
		_minLiquidity := new(big.Int).SetUint64(uint64(MinLiquidity))
		if _liquidity.Cmp(_minLiquidity) < 0 {
			err = errors.New("insufficientLiquidity k_sqrt(" + _liquidity.String() + "), need " + _minLiquidity.String())
			return err
		}
	}
	tokenMapReserve := make(map[string]string)
	tokenMapReserve[token0] = reserve0
	tokenMapReserve[token1] = reserve1
	pair := PoolPair{
		Model:          gorm.Model{},
		IsTokenZeroSat: _token0 == TokenSatTag,
		Token0:         _token0,
		Token1:         _token1,
		Reserve0:       tokenMapReserve[_token0],
		Reserve1:       tokenMapReserve[_token1],
	}
	// save in db
	return middleware.DB.Create(&pair).Error
}

func updatePair(token0 string, token1 string, reserve0 string, reserve1 string) (err error) {
	_token0, _token1, err := sortTokens(token0, token1)
	if err != nil {
		return utils.AppendErrorInfo(err, "sortTokens")
	}
	_reserve0, success := new(big.Int).SetString(reserve0, 10)
	if !success {
		return errors.New("reserve0 SetString(" + reserve0 + ") " + strconv.FormatBool(success))
	}
	_reserve1, success := new(big.Int).SetString(reserve1, 10)
	if !success {
		return errors.New("reserve1 SetString(" + reserve1 + ") " + strconv.FormatBool(success))
	}
	if !((_reserve0.Sign() > 0) && (_reserve1.Sign() > 0)) {
		err = errors.New("insufficientLiquidity(" + _reserve0.String() + "," + _reserve1.String() + ")")
		return err
	}
	// cmp with minimum liquidity sat
	if _token0 == TokenSatTag {
		_minLiquiditySat := new(big.Int).SetUint64(uint64(MinAddLiquiditySat))
		if _reserve0.Cmp(_minLiquiditySat) < 0 {
			err = errors.New("insufficientLiquidity Sat(" + _reserve0.String() + "), need " + _minLiquiditySat.String())
			return err
		}
	} else {
		// _token0 != TokenSatTag
		// TODO: Add MinLiquidity check for pair whose token0 is not sat
		_liquidity := new(big.Int).Sqrt(new(big.Int).Mul(_reserve0, _reserve1))
		_minLiquidity := new(big.Int).SetUint64(uint64(MinLiquidity))
		if _liquidity.Cmp(_minLiquidity) < 0 {
			err = errors.New("insufficientLiquidity k_sqrt(" + _liquidity.String() + "), need " + _minLiquidity.String())
			return err
		}
	}
	tokenMapReserve := make(map[string]string)
	tokenMapReserve[token0] = reserve0
	tokenMapReserve[token1] = reserve1
	pair, err := getPair(token0, token1)
	if err != nil {
		return utils.AppendErrorInfo(err, "getPair")
	}
	pair.Reserve0 = tokenMapReserve[_token0]
	pair.Reserve1 = tokenMapReserve[_token1]

	return middleware.DB.
		Model(&PoolPair{}).
		Where("token0 = ? AND token1 = ?", _token0, _token1).
		Updates(map[string]any{
			"reserve0": tokenMapReserve[_token0],
			"reserve1": tokenMapReserve[_token1],
		}).
		Error
}

func _newPair(token0 string, token1 string, reserve0 string, reserve1 string) (pair *PoolPair, err error) {
	// sort token
	_token0, _token1, err := sortTokens(token0, token1)
	if err != nil {
		return new(PoolPair), utils.AppendErrorInfo(err, "sortTokens")
	}
	// check reserves
	_reserve0, success := new(big.Int).SetString(reserve0, 10)
	if !success {
		return new(PoolPair), errors.New("reserve0 SetString(" + reserve0 + ") " + strconv.FormatBool(success))
	}
	_reserve1, success := new(big.Int).SetString(reserve1, 10)
	if !success {
		return new(PoolPair), errors.New("reserve1 SetString(" + reserve1 + ") " + strconv.FormatBool(success))
	}
	if !((_reserve0.Sign() > 0) && (_reserve1.Sign() > 0)) {
		err = errors.New("insufficientLiquidity(" + _reserve0.String() + "," + _reserve1.String() + ")")
		return new(PoolPair), err
	}
	// cmp with minimum liquidity sat
	if _token0 == TokenSatTag {
		_minLiquiditySat := new(big.Int).SetUint64(uint64(MinAddLiquiditySat))
		if _reserve0.Cmp(_minLiquiditySat) < 0 {
			err = errors.New("insufficientLiquidity Sat(" + _reserve0.String() + "), need " + _minLiquiditySat.String())
			return new(PoolPair), err
		}
	} else {
		// _token0 != TokenSatTag
		// TODO: Add MinLiquidity check for pair whose token0 is not sat
		_liquidity := new(big.Int).Sqrt(new(big.Int).Mul(_reserve0, _reserve1))
		_minLiquidity := new(big.Int).SetUint64(uint64(MinLiquidity))
		if _liquidity.Cmp(_minLiquidity) < 0 {
			err = errors.New("insufficientLiquidity k_sqrt(" + _liquidity.String() + "), need " + _minLiquidity.String())
			return new(PoolPair), err
		}
	}
	tokenMapReserve := make(map[string]string)
	tokenMapReserve[token0] = reserve0
	tokenMapReserve[token1] = reserve1
	pair = &PoolPair{
		IsTokenZeroSat: _token0 == TokenSatTag,
		Token0:         _token0,
		Token1:         _token1,
		Reserve0:       tokenMapReserve[_token0],
		Reserve1:       tokenMapReserve[_token1],
	}
	return pair, nil
}

func newPairBig(token0 string, token1 string, _reserve0 *big.Int, _reserve1 *big.Int) (pair *PoolPair, err error) {
	// sort token
	_token0, _token1, err := sortTokens(token0, token1)
	if err != nil {
		return new(PoolPair), utils.AppendErrorInfo(err, "sortTokens")
	}
	// check reserves
	if !((_reserve0.Sign() > 0) && (_reserve1.Sign() > 0)) {
		err = errors.New("insufficientLiquidity(" + _reserve0.String() + "," + _reserve1.String() + ")")
		return new(PoolPair), err
	}
	// cmp with minimum liquidity sat
	if _token0 == TokenSatTag {
		_minLiquiditySat := new(big.Int).SetUint64(uint64(MinAddLiquiditySat))
		if _reserve0.Cmp(_minLiquiditySat) < 0 {
			err = errors.New("insufficientLiquidity Sat(" + _reserve0.String() + "), need " + _minLiquiditySat.String())
			return new(PoolPair), err
		}
	} else {
		// _token0 != TokenSatTag
		// TODO: Add MinLiquidity check for pair whose token0 is not sat
		_liquidity := new(big.Int).Sqrt(new(big.Int).Mul(_reserve0, _reserve1))
		_minLiquidity := new(big.Int).SetUint64(uint64(MinLiquidity))
		if _liquidity.Cmp(_minLiquidity) < 0 {
			err = errors.New("insufficientLiquidity k_sqrt(" + _liquidity.String() + "), need " + _minLiquidity.String())
			return new(PoolPair), err
		}
	}
	tokenMapReserve := make(map[string]string)
	tokenMapReserve[token0] = _reserve0.String()
	tokenMapReserve[token1] = _reserve1.String()
	pair = &PoolPair{
		IsTokenZeroSat: _token0 == TokenSatTag,
		Token0:         _token0,
		Token1:         _token1,
		Reserve0:       tokenMapReserve[_token0],
		Reserve1:       tokenMapReserve[_token1],
	}
	return pair, nil
}

func getToken1PriceBig(_reserve0 *big.Int, _reserve1 *big.Int) (price *big.Float, err error) {
	_reserve0Float := new(big.Float).SetInt(_reserve0)
	_reserve1Float := new(big.Float).SetInt(_reserve1)
	_price := new(big.Float).Quo(_reserve0Float, _reserve1Float)
	if _price.IsInf() {
		return new(big.Float), errors.New("price is Inf")
	}
	if _price.Sign() <= 0 {
		return new(big.Float), errors.New("price(" + _price.String() + ") is less equal than zero")
	}
	return _price, nil
}
