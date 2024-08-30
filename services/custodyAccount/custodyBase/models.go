package custodyBase

import (
	"errors"
	"time"
)

const (
	Timeout = 20 * time.Second
)

type BtcChannel error

var (
	GetbalanceErr BtcChannel = errors.New("GetbalanceErr")
)

var TimeoutErr = errors.New("TimeoutErr")

type AssetPacketErr error

var (
	NotEnoughFeeFunds   AssetPacketErr = errors.New("not enough Fee funds")
	NotEnoughAssetFunds AssetPacketErr = errors.New("not enough Asset funds")
	DecodeAddressFail   AssetPacketErr = errors.New("decode Address fail")
	GetBalanceErr       AssetPacketErr = errors.New("get balance fail")
)
var (
	AssetPaymentFee int64 = 2000
	ServerFee             = 0
)

type Balance struct {
	AssetId string `json:"assetId"`
	Amount  int64  `json:"amount"`
}
