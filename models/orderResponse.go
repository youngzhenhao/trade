package models

type OrderResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Order   *Order `json:"order,omitempty"`
}
