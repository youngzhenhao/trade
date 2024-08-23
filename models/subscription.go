package models

import "time"

type Subscription struct {
	Channel     string
	Interval    time.Duration
	TradingPair string
	Online      string
	OrderType   string
	Status      string
}
