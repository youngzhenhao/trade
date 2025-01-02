package satBackQueue

import (
	"errors"
	"strconv"
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
