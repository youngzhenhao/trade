package dao

import (
	"trade/middleware"
	"trade/services/cpamm"
)

func cpAmmAutoMigrate() (err error) {
	if err = middleware.DB.AutoMigrate(&cpamm.PoolPair{}); err != nil {
		return err
	}
	return nil
}
