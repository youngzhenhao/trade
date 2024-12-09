package poolAccount

import "errors"

var (
	ErrorNotEnoughBalance = errors.New("not enough balance")
	ErrorDbError          = errors.New("db error")
)
