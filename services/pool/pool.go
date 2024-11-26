package pool

import "trade/middleware"

// TODO: Add Liquidity
func AddLiquidity(tokenA string, tokenB string, amountADesired string, amountBDesired string, amountAMin string, amountBMin string) error {
	token0, token1, err := sortTokens(tokenA, tokenB)
	tx := middleware.DB.Begin()

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
