package pool

import (
	"errors"
	"gorm.io/gorm"
	"math/big"
	"strconv"
	"trade/middleware"
	"trade/utils"
)

type Pair struct {
	gorm.Model
	IsTokenZeroSat bool   `json:"is_token_zero_sat" gorm:"index"`
	Token0         string `json:"token0" gorm:"type:varchar(255);uniqueIndex:idx_token_0_token_1"`
	Token1         string `json:"token1" gorm:"type:varchar(255);uniqueIndex:idx_token_0_token_1"`
	Reserve0       string `json:"reserve0" gorm:"type:varchar(255)"`
	Reserve1       string `json:"reserve1" gorm:"type:varchar(255)"`
}

func getPair(token0 string, token1 string) (pair *Pair, err error) {
	_token0, _token1, err := sortTokens(token0, token1)
	if err != nil {
		return new(Pair), utils.AppendErrorInfo(err, "sortTokens")
	}
	var _pair Pair
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
		return utils.AppendErrorInfo(err, "SetString("+reserve0+") "+strconv.FormatBool(success))
	}
	_reserve1, success := new(big.Int).SetString(reserve1, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+reserve1+") "+strconv.FormatBool(success))
	}
	if !((_reserve0.Sign() > 0) && (_reserve1.Sign() > 0)) {
		err = errors.New("insufficientLiquidity(" + _reserve0.String() + "," + _reserve1.String() + ")")
		return err
	}
	// cmp with minimum liquidity
	_k := _reserve0.Mul(_reserve0, _reserve1)
	_minLiquidity := new(big.Int).SetUint64(uint64(MinLiquidity))
	if _k.Cmp(_minLiquidity) < 0 {
		err = errors.New("insufficientLiquidity k(" + _k.String() + "), need " + _minLiquidity.String())
		return
	}
	tokenMapReserve := make(map[string]string)
	tokenMapReserve[token0] = reserve0
	tokenMapReserve[token1] = reserve1
	pair := Pair{
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

// TODO: Mutex
func updatePair(token0 string, token1 string, reserve0 string, reserve1 string) (err error) {
	_token0, _token1, err := sortTokens(token0, token1)
	if err != nil {
		return utils.AppendErrorInfo(err, "sortTokens")
	}
	_reserve0, success := new(big.Int).SetString(reserve0, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+reserve0+") "+strconv.FormatBool(success))
	}
	_reserve1, success := new(big.Int).SetString(reserve1, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+reserve1+") "+strconv.FormatBool(success))
	}
	if !((_reserve0.Sign() > 0) && (_reserve1.Sign() > 0)) {
		err = errors.New("insufficientLiquidity(" + _reserve0.String() + "," + _reserve1.String() + ")")
		return err
	}
	_k := _reserve0.Mul(_reserve0, _reserve1)
	_minLiquidity := new(big.Int).SetUint64(uint64(MinLiquidity))
	if _k.Cmp(_minLiquidity) < 0 {
		err = errors.New("insufficientLiquidity k(" + _k.String() + "), need " + _minLiquidity.String())
		return
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
		Model(&Pair{}).
		Where("token0 = ? AND token1 = ?", _token0, _token1).
		Updates(map[string]any{
			"reserve0": tokenMapReserve[_token0],
			"reserve1": tokenMapReserve[_token1],
		}).
		Error
}

func NewPair(token0 string, token1 string, reserve0 string, reserve1 string) (pair *Pair, err error) {
	// sort token
	_token0, _token1, err := sortTokens(token0, token1)
	if err != nil {
		return new(Pair), utils.AppendErrorInfo(err, "sortTokens")
	}
	// check reserves
	_reserve0, success := new(big.Int).SetString(reserve0, 10)
	if !success {
		return new(Pair), utils.AppendErrorInfo(err, "SetString("+reserve0+") "+strconv.FormatBool(success))
	}
	_reserve1, success := new(big.Int).SetString(reserve1, 10)
	if !success {
		return new(Pair), utils.AppendErrorInfo(err, "SetString("+reserve1+") "+strconv.FormatBool(success))
	}
	if !((_reserve0.Sign() > 0) && (_reserve1.Sign() > 0)) {
		err = errors.New("insufficientLiquidity(" + _reserve0.String() + "," + _reserve1.String() + ")")
		return new(Pair), err
	}
	// cmp with minimum liquidity
	_k := _reserve0.Mul(_reserve0, _reserve1)
	_minLiquidity := new(big.Int).SetUint64(uint64(MinLiquidity))
	if _k.Cmp(_minLiquidity) < 0 {
		err = errors.New("insufficientLiquidity k(" + _k.String() + "), need " + _minLiquidity.String())
		return new(Pair), err
	}
	tokenMapReserve := make(map[string]string)
	tokenMapReserve[token0] = reserve0
	tokenMapReserve[token1] = reserve1
	pair = &Pair{
		IsTokenZeroSat: _token0 == TokenSatTag,
		Token0:         _token0,
		Token1:         _token1,
		Reserve0:       tokenMapReserve[_token0],
		Reserve1:       tokenMapReserve[_token1],
	}
	return pair, nil
}
