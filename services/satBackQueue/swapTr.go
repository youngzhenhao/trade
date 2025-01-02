package satBackQueue

import (
	"errors"
	"math/big"
	"time"
	"trade/middleware"
	"trade/utils"
)

type SwapTr struct {
	ID           uint   `json:"id"`
	Price        string `json:"price"`
	Number       string `json:"number"`
	TotalPrice   string `json:"total_price"`
	NpubKey      string `json:"npub_key"`
	TrUnixtimeMs int64  `json:"tr_unixtime_ms"`
	AssetsID     string `json:"assets_id"`
	Type         string `json:"type"`
}

const (
	SwapTrTypeBuy  = "buy"
	SwapTrTypeSell = "sell"
)

type SwapTrsScan struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	TokenIn   string    `json:"token_in"`
	TokenOut  string    `json:"token_out"`
	AmountIn  string    `json:"amount_in"`
	AmountOut string    `json:"amount_out"`
}

func QueryNotPushedSwapTrsScan() (swapTrsScans []SwapTrsScan, err error) {

	tx := middleware.DB.Begin()

	swapTrsScans = []SwapTrsScan{}

	err = tx.Table("pool_swap_records").
		Select("id,created_at,username,token_in,token_out,amount_in,amount_out").
		Order("id desc").
		Where("is_pushed_queue = ?", false).
		Scan(&swapTrsScans).
		Error
	if err != nil {
		return []SwapTrsScan{}, utils.AppendErrorInfo(err, "select SwapTrsScan")
	}

	tx.Rollback()

	return swapTrsScans, nil
}

func SwapTrsScanToSwapTr(swapTrsScan SwapTrsScan) (swapTr SwapTr, err error) {

	var token0, token1 string
	token0, token1, err = sortTokens(swapTrsScan.TokenIn, swapTrsScan.TokenOut)
	if err != nil {
		return SwapTr{}, utils.AppendErrorInfo(err, "sortTokens")
	}

	isTokenZeroSat := token0 == TokenSatTag

	if isTokenZeroSat {

		var AmountSat, AmountAsset, _type string

		if token0 == swapTrsScan.TokenIn {
			// sat In, asset Out
			AmountSat = swapTrsScan.AmountIn
			AmountAsset = swapTrsScan.AmountOut
			_type = SwapTrTypeBuy

		} else {
			//	sat Out, asset In
			AmountSat = swapTrsScan.AmountOut
			AmountAsset = swapTrsScan.AmountIn
			_type = SwapTrTypeSell
		}

		var satFloat, assetFloat, price *big.Float
		var success bool

		assetFloat, success = new(big.Float).SetString(AmountAsset)
		if !success {
			return SwapTr{}, errors.New("AmountAsset big.Float.SetString")
		}
		satFloat, success = new(big.Float).SetString(AmountSat)
		if !success {
			return SwapTr{}, errors.New("AmountSat big.Float.SetString")
		}

		price = new(big.Float).Quo(satFloat, assetFloat)

		return SwapTr{
			ID:           swapTrsScan.ID,
			Price:        price.String(),
			Number:       AmountAsset,
			TotalPrice:   AmountSat,
			NpubKey:      swapTrsScan.Username,
			TrUnixtimeMs: swapTrsScan.CreatedAt.UnixMilli(),
			AssetsID:     token1,
			Type:         _type,
		}, nil

	} else {
		return SwapTr{}, errors.New("not support assets swap now")
	}
}

func SwapTrsScansToSwapTrs(swapTrsScans []SwapTrsScan) (swapTrs []SwapTr, err error) {
	swapTrs = []SwapTr{}
	for _, swapTrsScan := range swapTrsScans {
		swapTr, err := SwapTrsScanToSwapTr(swapTrsScan)
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "SwapTrsScanToSwapTr")
		}
		swapTrs = append(swapTrs, swapTr)
	}
	return swapTrs, nil
}

func QueryNotPushedSwapTrs() (swapTrs []SwapTr, err error) {
	swapTrsScans, err := QueryNotPushedSwapTrsScan()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "QueryNotPushedSwapTrsScan")
	}
	swapTrs, err = SwapTrsScansToSwapTrs(swapTrsScans)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "SwapRecordInfosToSwapTrs")
	}
	return swapTrs, nil
}
