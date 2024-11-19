package cpamm

import (
	"trade/middleware"
	"trade/utils"
)

func (p *PoolPair) dbCreate() error {
	return middleware.DB.Create(p).Error
}

func (p *PoolPair) dbSave() error {
	return middleware.DB.Save(p).Error
}

func (p *PoolPair) dbRead(token0 string, token1 string) error {
	_token0, _token1, err := sortTokens(token0, token1)
	if err != nil {
		return utils.AppendErrorInfo(err, "sortTokens")
	}
	return middleware.DB.Where("token0 = ? and token1 = ?", _token0, _token1).First(&p).Error
}
