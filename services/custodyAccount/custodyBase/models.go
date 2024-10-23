package custodyBase

import (
	"errors"
	"time"
	"trade/models"
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

type PaymentResponse struct {
	Timestamp   int64               `json:"timestamp"`
	BillType    models.BalanceType  `json:"bill_type"`
	Away        models.BalanceAway  `json:"away"`
	Target      *string             `json:"target"`
	PaymentHash *string             `json:"payment_hash"`
	Amount      float64             `json:"amount"`
	AssetId     *string             `json:"asset_id"`
	State       models.BalanceState `json:"state"`
	Fee         uint64              `json:"fee"`
	//deprecated
	Invoice *string `json:"invoice"`
	Address *string `json:"addr"`
}

type PaymentList struct {
	PaymentList []PaymentResponse `json:"payments"`
}

func (r *PaymentList) GetTxString() string {
	return ""
}
