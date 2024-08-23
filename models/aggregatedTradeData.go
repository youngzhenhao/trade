package models

import "time"

type AggregatedTradeData struct {
	Period        time.Time
	AvgPrice      float64
	TotalQuantity float64
	TradeCount    int
}
