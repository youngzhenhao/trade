package cpamm

import (
	"fmt"
	"trade/utils"
)

func newPair(tokenA string, tokenB string) error {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return utils.AppendErrorInfo(err, "sortTokens")
	}

	pair := new(PoolPair)
	pair.initialize(token0, token1)

	fmt.Println(utils.ValueJsonString(pair))
	// TODO: store to db
	//查询数据库是否已存在

	return nil
}
