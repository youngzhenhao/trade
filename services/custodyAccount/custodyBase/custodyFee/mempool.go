package custodyFee

import (
	"errors"
	"math"
	"sync"
	"time"
	"trade/api"
	"trade/btlLog"
	"trade/config"
	"trade/middleware"
	"trade/models"
	"trade/services/btldb"
	"trade/utils"
)

var (
	lastExecutedTime time.Time
	feeInitOnceLock  sync.Once
)

func EstimateFee() *responseTransformed {
	feeInitOnceLock.Do(func() {
		lastExecutedTime = time.Now()
		updateFeeRateByMempool()
	})
	currentTime := time.Now()
	// 检查距离上次执行的时间
	if currentTime.Sub(lastExecutedTime) >= 10*time.Minute {
		lastExecutedTime = currentTime
		updateFeeRateByMempool()
	}
	FeeList, err := getFeeRateResponseTransformed()
	if err != nil {
		return nil
	}
	return FeeList
}

func updateFeeRateByMempool() {
	err := checkIfUpdateFeeRateInfoByMempool()
	if err != nil {
		return
	}
}

func getFeeRateResponseTransformed() (*responseTransformed, error) {
	var feeRateInfos []models.FeeRateInfo
	units := []models.FeeRateType{models.FeeRateTypeSatPerB, models.FeeRateTypeSatPerKw}
	names := []string{"fastest_fee", "half_hour_fee", "hour_fee", "economy_fee", "minimum_fee"}
	for _, unit := range units {
		for _, name := range names {
			feeRateInfo, err := getFeeRateInfoByNameAndUnit(name, unit)
			if err != nil {
				// @dev: Do not return
				return nil, err
			}
			feeRateInfos = append(feeRateInfos, *feeRateInfo)
		}
	}
	transformed, err := processFeeRateInfosToResponseTransformed(feeRateInfos)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "processFeeRateInfosToResponseTransformed")
	}
	return transformed, nil
}

func processFeeRateInfosToResponseTransformed(feeRateInfos []models.FeeRateInfo) (*responseTransformed, error) {
	var feeRateResponseTransformed responseTransformed
	units := []models.FeeRateType{models.FeeRateTypeSatPerB, models.FeeRateTypeSatPerKw}
	names := []string{"fastest_fee", "half_hour_fee", "hour_fee", "economy_fee", "minimum_fee"}
	for _, feeRateInfo := range feeRateInfos {
		if feeRateInfo.Unit == units[0] {
			if feeRateInfo.Name == names[0] {
				feeRateResponseTransformed.SatPerB.FastestFee = int(feeRateInfo.FeeRate)
			} else if feeRateInfo.Name == names[1] {
				feeRateResponseTransformed.SatPerB.HalfHourFee = int(feeRateInfo.FeeRate)
			} else if feeRateInfo.Name == names[2] {
				feeRateResponseTransformed.SatPerB.HourFee = int(feeRateInfo.FeeRate)
			} else if feeRateInfo.Name == names[3] {
				feeRateResponseTransformed.SatPerB.EconomyFee = int(feeRateInfo.FeeRate)
			} else if feeRateInfo.Name == names[4] {
				feeRateResponseTransformed.SatPerB.MinimumFee = int(feeRateInfo.FeeRate)
			}
		} else if feeRateInfo.Unit == units[1] {
			if feeRateInfo.Name == names[0] {
				feeRateResponseTransformed.SatPerKw.FastestFee = int(feeRateInfo.FeeRate)
			} else if feeRateInfo.Name == names[1] {
				feeRateResponseTransformed.SatPerKw.HalfHourFee = int(feeRateInfo.FeeRate)
			} else if feeRateInfo.Name == names[2] {
				feeRateResponseTransformed.SatPerKw.HourFee = int(feeRateInfo.FeeRate)
			} else if feeRateInfo.Name == names[3] {
				feeRateResponseTransformed.SatPerKw.EconomyFee = int(feeRateInfo.FeeRate)
			} else if feeRateInfo.Name == names[4] {
				feeRateResponseTransformed.SatPerKw.MinimumFee = int(feeRateInfo.FeeRate)
			}
		}
	}
	return &feeRateResponseTransformed, nil
}

func checkIfUpdateFeeRateInfoByMempool() (err error) {
	if config.GetLoadConfig().FairLaunchConfig.IsAutoUpdateFeeRate {
		err = updateFeeRateInfoByMempool()
		if err != nil {
			return utils.AppendErrorInfo(err, "updateFeeRateInfoByMempool")
		}
	}
	return nil
}

func updateFeeRateInfoByMempool() error {
	//var feeRateInfos []models.FeeRateInfo
	feeRateResponse, err := getFeeRateResponseTransformedByMempool()
	if err != nil {
		return utils.AppendErrorInfo(err, "getFeeRateResponseTransformedByMempool")
	}
	units := []models.FeeRateType{models.FeeRateTypeSatPerB, models.FeeRateTypeSatPerKw}
	names := []string{"fastest_fee", "half_hour_fee", "hour_fee", "economy_fee", "minimum_fee"}
	for _, unit := range units {
		for _, name := range names {
			var feeRate int
			feeRate, err = feeRateResponse.GetFeeRateByNameAndUnit(unit, name)
			err = updateFeeRateInfoByNameAndUnitIfNotExistThenCreate(name, unit, feeRate)
			if err != nil {
				// @dev: Do not return
			}
		}
	}
	return nil
}

func getFeeRateResponseTransformedByMempool() (*responseTransformed, error) {
	fees, err := api.MempoolGetRecommendedFees()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "MempoolGetRecommendedFees")
	}
	return &responseTransformed{
		SatPerB: mempoolFeeRate{
			FastestFee:  fees.FastestFee,
			HalfHourFee: fees.HalfHourFee,
			HourFee:     fees.HourFee,
			EconomyFee:  fees.EconomyFee,
			MinimumFee:  fees.MinimumFee,
		},
		SatPerKw: mempoolFeeRate{
			FastestFee:  satPerBToSatPerKw(fees.FastestFee),
			HalfHourFee: satPerBToSatPerKw(fees.HalfHourFee),
			HourFee:     satPerBToSatPerKw(fees.HourFee),
			EconomyFee:  satPerBToSatPerKw(fees.EconomyFee),
			MinimumFee:  satPerBToSatPerKw(fees.MinimumFee),
		},
	}, nil
}

type responseTransformed struct {
	SatPerB  mempoolFeeRate
	SatPerKw mempoolFeeRate
}

type mempoolFeeRate struct {
	FastestFee  int `json:"fastest_fee"`
	HalfHourFee int `json:"half_hour_fee"`
	HourFee     int `json:"hour_fee"`
	EconomyFee  int `json:"economy_fee"`
	MinimumFee  int `json:"minimum_fee"`
}

// satPerBToSatPerKw
// @Description: sat/b to sat/kw
func satPerBToSatPerKw(feeRateSatPerB int) (feeRateSatPerKw int) {
	return int(math.Ceil(float64(feeRateSatPerB) * 1000 / 4))
}

func (feeRate *responseTransformed) GetFeeRateByNameAndUnit(unit models.FeeRateType, name string) (rate int, err error) {
	units := []models.FeeRateType{models.FeeRateTypeSatPerB, models.FeeRateTypeSatPerKw}
	names := []string{"fastest_fee", "half_hour_fee", "hour_fee", "economy_fee", "minimum_fee"}
	if unit == units[0] {
		if name == names[0] {
			return feeRate.SatPerB.FastestFee, nil
		} else if name == names[1] {
			return feeRate.SatPerB.HalfHourFee, nil
		} else if name == names[2] {
			return feeRate.SatPerB.HourFee, nil
		} else if name == names[3] {
			return feeRate.SatPerB.EconomyFee, nil
		} else if name == names[4] {
			return feeRate.SatPerB.MinimumFee, nil
		}
	} else if unit == units[1] {
		if name == names[0] {
			return feeRate.SatPerKw.FastestFee, nil
		} else if name == names[1] {
			return feeRate.SatPerKw.HalfHourFee, nil
		} else if name == names[2] {
			return feeRate.SatPerKw.HourFee, nil
		} else if name == names[3] {
			return feeRate.SatPerKw.EconomyFee, nil
		} else if name == names[4] {
			return feeRate.SatPerKw.MinimumFee, nil
		}
	}
	err = errors.New("can't get fee rate info by match name and unit")
	return 0, err
}

func updateFeeRateInfoByNameAndUnitIfNotExistThenCreate(name string, unit models.FeeRateType, feeRate int) error {
	f := btldb.FeeRateInfoStore{DB: middleware.DB}
	feeRateInfo, err := getFeeRateInfoByNameAndUnit(name, unit)
	if err != nil {
		feeRateInfo = &models.FeeRateInfo{
			Name:    name,
			Unit:    unit,
			FeeRate: float64(feeRate),
		}
		err = f.CreateFeeRateInfo(feeRateInfo)
		if err != nil {
			return utils.AppendErrorInfo(err, "CreateFeeRateInfo")
		}
		btlLog.FEE.Info("%v %v %v", name, unit, "FeeRateInfo record created.")
	} else {
		feeRateInfo.FeeRate = float64(feeRate)
		err = f.UpdateFeeRateInfo(feeRateInfo)
		if err != nil {
			return utils.AppendErrorInfo(err, "UpdateFeeRateInfo")
		}
	}
	return nil
}

func getFeeRateInfoByNameAndUnit(name string, unit models.FeeRateType) (*models.FeeRateInfo, error) {
	var feeRateInfo models.FeeRateInfo
	err := middleware.DB.Where("name = ? AND unit = ?", name, unit).First(&feeRateInfo).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "First feeRateInfo")
	}
	return &feeRateInfo, nil
}
