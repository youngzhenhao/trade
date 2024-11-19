package cpamm

import (
	"gorm.io/gorm"
	"math/big"
	"strconv"
	"sync"
	"time"
	"trade/utils"
)

const (
	MINIMUM_LIQUIDITY uint64 = 1e3
	TokenSatTag       string = "sat"
)

var PairMapMutex map[uint]sync.Mutex

type PoolPair struct {
	gorm.Model
	IsTokenZeroSat       bool   `json:"is_token_zero_sat" gorm:"index"`
	Token0               string `json:"token0" gorm:"type:varchar(255);index;uniqueIndex:idx_token_0_token_1"`
	Token1               string `json:"token1" gorm:"type:varchar(255);index;uniqueIndex:idx_token_0_token_1"`
	Reserve0             string `json:"reserve0" gorm:"type:varchar(255)"`
	Reserve1             string `json:"reserve1" gorm:"type:varchar(255)"`
	Price0CumulativeLast string `json:"price0_cumulative_last" gorm:"type:varchar(255)"`
	Price1CumulativeLast string `json:"price1_cumulative_last" gorm:"type:varchar(255)"`
	KLast                string `json:"k_last" gorm:"type:varchar(255)"`
}

func (p *PoolPair) getReserves() (_reserve0 string, _reserve1 string) {
	_reserve0 = p.Reserve0
	_reserve1 = p.Reserve1
	return _reserve0, _reserve1
}

// TODO:
func _safeTransfer(token string, to string, uint string) {
	//(bool success, bytes memory data) = token.call(abi.encodeWithSelector(SELECTOR, to, value));
	//require(success && (data.length == 0 || abi.decode(data, (bool))), 'UniswapV2: TRANSFER_FAILED');
}

func (p *PoolPair) initialize(_token0 string, _token1 string) {
	if _token0 == TokenSatTag {
		p.IsTokenZeroSat = true
	}
	p.Token0 = _token0
	p.Token1 = _token1
	p.Reserve0 = "0"
	p.Reserve1 = "0"
	p.Price0CumulativeLast = "0"
	p.Price1CumulativeLast = "0"
	p.KLast = "0"
}

// TODO: need to test
// update reserves and, on the first call per block, price accumulators
func (p *PoolPair) _update(balance0 string, balance1 string) (err error) {
	now := time.Now()
	timeElapsed := now.Sub(p.UpdatedAt).Seconds()
	_reserve0Float, success := new(big.Float).SetString(p.Reserve0)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+p.Reserve0+") "+strconv.FormatBool(success))
	}
	_reserve1Float, success := new(big.Float).SetString(p.Reserve1)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+p.Reserve1+") "+strconv.FormatBool(success))
	}
	timeElapsedFloat := new(big.Float).SetFloat64(timeElapsed)
	Price0CumulativeLastFloat, success := new(big.Float).SetString(p.Price0CumulativeLast)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+p.Price0CumulativeLast+") "+strconv.FormatBool(success))
	}
	Price1CumulativeLastFloat, success := new(big.Float).SetString(p.Price1CumulativeLast)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+p.Price1CumulativeLast+") "+strconv.FormatBool(success))
	}
	if timeElapsedFloat.Sign() > 0 && _reserve0Float.Sign() != 0 && _reserve1Float.Sign() != 0 {
		Price0CumulativeLastFloat = Price0CumulativeLastFloat.Add(Price0CumulativeLastFloat, new(big.Float).Quo(new(big.Float).Mul(timeElapsedFloat, _reserve1Float), _reserve0Float))
		Price1CumulativeLastFloat = Price1CumulativeLastFloat.Add(Price1CumulativeLastFloat, new(big.Float).Quo(new(big.Float).Mul(timeElapsedFloat, _reserve0Float), _reserve1Float))
		// @dev: Update prices
		p.Price0CumulativeLast = Price0CumulativeLastFloat.String()
		p.Price1CumulativeLast = Price1CumulativeLastFloat.String()
	}
	// @dev: Validate balances
	_, success = new(big.Int).SetString(balance0, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+balance0+") "+strconv.FormatBool(success))
	}
	_, success = new(big.Int).SetString(balance1, 10)
	if !success {
		return utils.AppendErrorInfo(err, "SetString("+balance1+") "+strconv.FormatBool(success))
	}
	// @dev: Update reserves
	p.Reserve0 = balance0
	p.Reserve1 = balance1
	return nil
}

// if fee is on, mint liquidity equivalent to 1/6th of the growth in sqrt(k)
func (p *PoolPair) _mintFee(_reserve0 string, _reserve1 string) (feeOn bool, err error) {

	//address feeTo = IUniswapV2Factory(factory).feeTo();
	//feeOn = feeTo != address(0);
	//uint _kLast = kLast; // gas savings
	//if (feeOn) {
	//	if (_kLast != 0) {
	//		uint rootK = Math.sqrt(uint(_reserve0).mul(_reserve1));
	//		uint rootKLast = Math.sqrt(_kLast);
	//		if (rootK > rootKLast) {
	//			uint numerator = totalSupply.mul(rootK.sub(rootKLast));
	//			uint denominator = rootK.mul(5).add(rootKLast);
	//			uint liquidity = numerator / denominator;
	//			if (liquidity > 0) _mint(feeTo, liquidity);
	//		}
	//	}
	//} else if (_kLast != 0) {
	//	kLast = 0;
	//}
	return false, nil
}

// TODO: mutex
func mint(to string) (liquidity string) {
	return ""
}

func burn(to string) (amount0 string, amount1 string) {
	return "", ""
}

func swap(amount0Out string, amount1Out string, to string, data string) {

}

func skim(to string) {

}

func _sync() {

}

// TODO: Create share when creating pair
