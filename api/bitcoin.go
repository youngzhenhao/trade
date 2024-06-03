package api

import "github.com/btcsuite/btcd/btcjson"

func EstimateSmartFeeAndGetResult(blocks int) (feeResult *btcjson.EstimateSmartFeeResult, err error) {
	return estimateSmartFee(int64(blocks), &btcjson.EstimateModeUnset)
}
