package pool

import (
	"errors"
	"math/big"
	"time"
	"trade/middleware"
	"trade/utils"
)

type SwapTr struct {
	ID uint `json:"id"`
	// amount of sat / amount of asset; float string
	Price string `json:"price"`
	// amount of asset; int string
	Number string `json:"number"`
	// amount of sat; int string
	TotalPrice string `json:"total_price"`
	NpubKey    string `json:"npub_key"`
	// microsecond
	TrUnixtimeMs int64  `json:"tr_unixtime_ms"`
	AssetsID     string `json:"assets_id"`
	Type         string `json:"type"`
}

const (
	SwapTrTypeBuy  = "buy"
	SwapTrTypeSell = "sell"
)

func SwapRecordInfoToSwapTr(swapRecordInfo SwapRecordInfo) (swapTr SwapTr, err error) {

	var token0, token1 string
	token0, token1, err = sortTokens(swapRecordInfo.TokenIn, swapRecordInfo.TokenOut)
	if err != nil {
		return SwapTr{}, utils.AppendErrorInfo(err, "sortTokens")
	}

	isTokenZeroSat := token0 == TokenSatTag

	if isTokenZeroSat {

		var AmountSat, AmountAsset string

		if token0 == swapRecordInfo.TokenIn {
			// sat In, asset Out
			AmountSat = swapRecordInfo.AmountIn
			AmountAsset = swapRecordInfo.AmountOut

		} else {
			//	sat Out, asset In
			AmountSat = swapRecordInfo.AmountOut
			AmountAsset = swapRecordInfo.AmountIn

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

		_time := utils.TimestampToTime(swapRecordInfo.Time)
		return SwapTr{
			ID:           swapRecordInfo.ID,
			Price:        price.String(),
			Number:       AmountAsset,
			TotalPrice:   AmountSat,
			NpubKey:      swapRecordInfo.Username,
			TrUnixtimeMs: _time.UnixMilli(),
			AssetsID:     token1,
		}, nil

	} else {
		return SwapTr{}, errors.New("not support assets swap now")
	}
}

func SwapRecordInfosToSwapTrs(swapRecordInfos *[]SwapRecordInfo) (swapTrs *[]SwapTr, err error) {
	if swapRecordInfos == nil {
		return nil, errors.New("swapRecordInfos is nil")
	}
	swapTrs = &[]SwapTr{}
	for _, swapRecordInfo := range *swapRecordInfos {
		swapTr, err := SwapRecordInfoToSwapTr(swapRecordInfo)
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "SwapRecordInfoToSwapTr")
		}
		*swapTrs = append(*swapTrs, swapTr)
	}
	return swapTrs, nil
}

func QuerySwapTrsScanCount(tokenA string, tokenB string) (count int64, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "sortTokens")
	}

	tx := middleware.DB.Begin()

	err = tx.Table("pool_pairs").
		Joins("join pool_swap_records on pool_pairs.id = pool_swap_records.pair_id").
		Where("pool_pairs.token0 = ? AND pool_pairs.token1 = ?", token0, token1).
		Count(&count).
		Error
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "select SwapTrsScan count")
	}

	tx.Rollback()

	return count, nil
}

type SwapTrsScan struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	TokenIn   string    `json:"token_in"`
	TokenOut  string    `json:"token_out"`
	AmountIn  string    `json:"amount_in"`
	AmountOut string    `json:"amount_out"`
}

func QuerySwapTrsScan(tokenA string, tokenB string, limit int, offset int) (swapTrsScans *[]SwapTrsScan, err error) {
	token0, token1, err := sortTokens(tokenA, tokenB)
	if err != nil {
		return new([]SwapTrsScan), utils.AppendErrorInfo(err, "sortTokens")
	}

	tx := middleware.DB.Begin()

	swapTrsScans = new([]SwapTrsScan)

	err = tx.Table("pool_pairs").
		Select("pool_swap_records.id,pool_swap_records.created_at,pool_swap_records.username,pool_swap_records.token_in,pool_swap_records.token_out,pool_swap_records.amount_in,pool_swap_records.amount_out").
		Joins("join pool_swap_records on pool_pairs.id = pool_swap_records.pair_id").
		Where("pool_pairs.token0 = ? AND pool_pairs.token1 = ?", token0, token1).
		Order("pool_swap_records.id desc").
		Limit(limit).
		Offset(offset).
		Scan(&swapTrsScans).
		Error
	if err != nil {
		return new([]SwapTrsScan), utils.AppendErrorInfo(err, "select SwapTrsScan")
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

func SwapTrsScansToSwapTrs(swapTrsScans *[]SwapTrsScan) (swapTrs *[]SwapTr, err error) {
	if swapTrsScans == nil {
		return nil, errors.New("swapTrsScans is nil")
	}
	swapTrs = &[]SwapTr{}
	for _, swapTrsScan := range *swapTrsScans {
		swapTr, err := SwapTrsScanToSwapTr(swapTrsScan)
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "SwapTrsScanToSwapTr")
		}
		*swapTrs = append(*swapTrs, swapTr)
	}
	return swapTrs, nil
}

func QuerySwapTrs(tokenA string, tokenB string, limit int, offset int) (swapTrs *[]SwapTr, err error) {
	swapTrsScans, err := QuerySwapTrsScan(tokenA, tokenB, limit, offset)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "QuerySwapRecords")
	}
	swapTrs, err = SwapTrsScansToSwapTrs(swapTrsScans)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "SwapRecordInfosToSwapTrs")
	}
	return swapTrs, nil
}
