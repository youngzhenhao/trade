package poolAccount

import "gorm.io/gorm"

type PoolAccount struct {
	gorm.Model
	PairId uint
}
