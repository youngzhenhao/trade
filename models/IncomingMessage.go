package models

type IncomingMessage struct {
	Type     string  `json:"type"`
	OrderID  string  `json:"order_id,omitempty"`
	Buyer    string  `json:"buyer,omitempty"`
	Quantity int     `json:"quantity,omitempty"`
	Price    float64 `json:"price,omitempty"`
}
