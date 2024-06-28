package models

type Network int

const (
	Mainnet Network = iota
	Testnet
	Regtest
)
