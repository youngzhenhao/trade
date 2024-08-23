package models

import "gorm.io/gorm"

type TradeOrder struct {
	gorm.Model
	OrderID       string  `gorm:"column:order_id" json:"order_id"`
	TradingPair   string  `gorm:"column:trading_pair" json:"trading_pair"`
	BitcoinAmount float64 `gorm:"column:bitcoin_amount;type:decimal(15,2)" json:"bitcoin_amount"`
	TapRootAsset  string  `gorm:"column:taproot_asset" json:"taproot_asset"`
	TapRootAmount float64 `gorm:"column:taproot_amount;type:decimal(15,2)" json:"taproot_amount"`
	TotalPrice    float64 `gorm:"column:total_price;type:decimal(15,2)" json:"total_price"`
	UnitPrice     float64 `gorm:"column:unit_price;type:decimal(15,2)" json:"unit_price"`
	OrderType     string  `gorm:"column:order_type"  json:"order_type"`
	Seller        string  `gorm:"column:seller" json:"seller"`
	Buyer         string  `gorm:"column:buyer"  json:"buyer,omitempty"`
	Online        bool    `gorm:"column:online"  json:"online"` // true for online orders, false for offline orders
	PSBTSeller    string  `gorm:"column:psbt_seller" json:"psbt_seller,omitempty"`
	PSBTBuyer     string  `gorm:"column:psbt_buyer" json:"psbt_buyer,omitempty"`
	Type          string  `gorm:"column:Type" json:"type"`
	Status        int16   `gorm:"column:status;type:smallint" json:"status"`
}

func (TradeOrder) TableName() string {
	return "tradeOrder"
}
