package models

import (
	"gorm.io/gorm"
	"time"
)

type TradeHistory struct {
	gorm.Model
	TradePair string    `gorm:"size:20;not null"`
	TradeTime time.Time `gorm:"not null"`
	UnitPrice float64   `gorm:"type:decimal(15,2);not null"`
	Buyer     string    `gorm:"size:100;not null"`
	Seller    string    `gorm:"size:100;not null"`
}

func (TradeHistory) TableName() string {
	return "trade_history"
}

func InsertTradeHistory(db *gorm.DB, trade *TradeHistory) error {
	return db.Create(trade).Error
}
