package models

type Order struct {
	OrderID       string  `json:"order_id"`
	TradingPair   string  `json:"trading_pair"`
	BitcoinAmount float64 `json:"bitcoin_amount"`
	RootAsset     string  `json:"root_asset"`
	RootAmount    float64 `json:"root_amount"`
	TotalPrice    float64 `json:"total_price"`
	UnitPrice     float64 `json:"unit_price"`
	OrderType     string  `json:"order_type"`
	Status        int     `json:"status"`
	Seller        string  `json:"seller"`
	Buyer         string  `json:"buyer,omitempty"`
	Online        bool    `json:"online"` // true for online orders, false for offline orders
	PSBTSeller    string  `json:"psbt_seller,omitempty"`
	PSBTBuyer     string  `json:"psbt_buyer,omitempty"`
	Type          string  `json:"type"`
}
