package custodyFee

import "math"

const (
	anchorFee = 1000
	assetBase = 170
)

func GetCustodyAssetFee() int {
	feeList := EstimateFee()
	if feeList == nil {
		return 0
	}
	return anchorFee + int(math.Ceil(float64(feeList.SatPerB.HalfHourFee)*float64(assetBase)*0.5)) + 500
}
