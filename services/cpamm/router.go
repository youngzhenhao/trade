package cpamm

import (
	"errors"
	"time"
)

const (
	MinFeeSat = 20
)

var (
	factory string
	WETH    string
)

func ensure(deadline *time.Time) (err error) {
	now := time.Now()
	if deadline.Before(now) {
		err = errors.New("Router: EXPIRED(" + deadline.Format(time.DateTime) + ";now: " + now.Format(time.DateTime) + ")")
		return err
	}
	return nil
}

// **** ADD LIQUIDITY ****
// @dev: No fee for adding liquidity

// TODO
func _addLiquidity(tokenA string, tokenB string, amountADesired string, amountBDesired string, amountAMin string, amountBMin string) (amountA string, amountB string, err error) {
	// create the pair if it doesn't exist yet
	//if (IUniswapV2Factory(factory).getPair(tokenA, tokenB) == address(0)) {
	//IUniswapV2Factory(factory).createPair(tokenA, tokenB);
	//}
	//(uint reserveA, uint reserveB) = UniswapV2Library.getReserves(factory, tokenA, tokenB);
	//if (reserveA == 0 && reserveB == 0) {
	//(amountA, amountB) = (amountADesired, amountBDesired);
	//} else {
	//uint amountBOptimal = UniswapV2Library.quote(amountADesired, reserveA, reserveB);
	//if (amountBOptimal <= amountBDesired) {
	//require(amountBOptimal >= amountBMin, 'UniswapV2Router: INSUFFICIENT_B_AMOUNT');
	//(amountA, amountB) = (amountADesired, amountBOptimal);
	//} else {
	//uint amountAOptimal = UniswapV2Library.quote(amountBDesired, reserveB, reserveA);
	//assert(amountAOptimal <= amountADesired);
	//require(amountAOptimal >= amountAMin, 'UniswapV2Router: INSUFFICIENT_A_AMOUNT');
	//(amountA, amountB) = (amountAOptimal, amountBDesired);
	//}
	//}
	return "", "", nil
}

// **** REMOVE LIQUIDITY ****

// **** REMOVE LIQUIDITY (supporting fee-on-transfer tokens) ****

// **** SWAP ****

// **** SWAP (supporting fee-on-transfer tokens) ****
