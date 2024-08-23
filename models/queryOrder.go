package models

type QueryOrder struct {
	TradingPair  string  `json:"trading_pair"`
	MinUnitPrice float64 `json:"minUnit_price"`
	MaxUnitPrice float64 `json:"maxUnit_price"`
	Online       bool    `json:"online"`
	OrderType    string  `json:"order_type"`
}
