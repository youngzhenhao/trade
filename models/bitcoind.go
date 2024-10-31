package models

type Network int

const (
	Mainnet Network = iota
	Testnet
	Regtest
)

func (n Network) String() string {
	networkMap := map[Network]string{
		Mainnet: "mainnet",
		Testnet: "testnet",
		Regtest: "regtest",
	}
	return networkMap[n]
}
