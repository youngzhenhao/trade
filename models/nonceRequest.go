package models

type NonceRequest struct {
	Username string `json:"userName"`
	Nonce    string `json:"nonce"`
}
