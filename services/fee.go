package services

import (
	"errors"
	"strconv"
	"trade/api"
	"trade/config"
	"trade/middleware"
	"trade/models"
	"trade/utils"
)

type (
	FeeRateInfoName string
)

var (
	GasFeeRateOfNumber0          float64         = 1
	GasFeeRateOfNumber1                          = 1.2
	GasFeeRateOfNumber2          float64         = 2
	GasFeeRateOfNumber3          float64         = 3
	GasFeeRateOfNumber4          float64         = 4
	GasFeeRateOfNumber5                          = 4.5
	GasFeeRateOfNumber6          float64         = 5
	GasFeeRateOfNumber7                          = 5.4
	GasFeeRateOfNumber8                          = 5.7
	GasFeeRateOfNumber9          float64         = 6
	GasFeeRateOfNumber10                         = 6.5
	GasFeeRateNameMempoolMainnet FeeRateInfoName = "mempool_mainnet"
	GasFeeRateNameBitcoind       FeeRateInfoName = "bitcoind"
	GasFeeRateNameDefault                        = GasFeeRateNameBitcoind
)

func UpdateAndGetFeeRateResponseTransformed() (*FeeRateResponseTransformed, error) {
	UpdateFeeRateByMempool()
	return GetFeeRateResponseTransformed()
}

func UpdateAndCalculateGasFeeRateByMempool(number int) (*FeeRateResponseTransformed, error) {
	UpdateFeeRateByMempool()
	return CalculateGasFeeRateByMempool(number)
}

func CalculateGasFeeRateByMempool(number int) (*FeeRateResponseTransformed, error) {
	feeRate, err := GetFeeRateResponseTransformed()
	rate, err := NumberToGasFeeRate(number)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "NumberToGasFeeRate")
	}
	return &FeeRateResponseTransformed{
		SatPerB: MempoolFeeRate{
			FastestFee:  int(float64(feeRate.SatPerB.FastestFee) * rate),
			HalfHourFee: int(float64(feeRate.SatPerB.HalfHourFee) * rate),
			HourFee:     int(float64(feeRate.SatPerB.HourFee) * rate),
			EconomyFee:  int(float64(feeRate.SatPerB.EconomyFee) * rate),
			MinimumFee:  int(float64(feeRate.SatPerB.MinimumFee) * rate),
		},
		SatPerKw: MempoolFeeRate{
			FastestFee:  int(float64(feeRate.SatPerKw.FastestFee) * rate),
			HalfHourFee: int(float64(feeRate.SatPerKw.HalfHourFee) * rate),
			HourFee:     int(float64(feeRate.SatPerKw.HourFee) * rate),
			EconomyFee:  int(float64(feeRate.SatPerKw.EconomyFee) * rate),
			MinimumFee:  int(float64(feeRate.SatPerKw.MinimumFee) * rate),
		},
	}, nil
}

func UpdateAndEstimateSmartFeeRateSatPerKw() (estimatedFeeSatPerKw int, err error) {
	UpdateFeeRate()
	return EstimateSmartFeeRateSatPerKw()
}

func UpdateAndEstimateSmartFeeRateSatPerB() (estimatedFeeSatPerB int, err error) {
	UpdateFeeRate()
	return EstimateSmartFeeRateSatPerB()
}

func UpdateAndEstimateSmartFeeRateBtcPerKb() (estimatedFeeBtcPerKb float64, err error) {
	UpdateFeeRate()
	return EstimateSmartFeeRateBtcPerKb()
}

func UpdateAndCalculateGasFeeRateSatPerKw(number int) (int, error) {
	UpdateFeeRate()
	return CalculateGasFeeRateSatPerKw(number)
}

func UpdateAndCalculateGasFeeRateSatPerB(number int) (int, error) {
	UpdateFeeRate()
	return CalculateGasFeeRateSatPerB(number)
}

func UpdateAndCalculateGasFeeRateBtcPerKb(number int) (float64, error) {
	UpdateFeeRate()
	return CalculateGasFeeRateBtcPerKb(number)
}

func NumberToGasFeeRate(number int) (gasFeeRate float64, err error) {
	if number < 0 || number > 10 {
		err = errors.New("number out of range")
		// Max rate
		return GasFeeRateOfNumber10, err
	} else if number == 1 {
		return GasFeeRateOfNumber1, nil
	} else if number == 2 {
		return GasFeeRateOfNumber2, nil
	} else if number == 3 {
		return GasFeeRateOfNumber3, nil
	} else if number == 4 {
		return GasFeeRateOfNumber4, nil
	} else if number == 5 {
		return GasFeeRateOfNumber5, nil
	} else if number == 6 {
		return GasFeeRateOfNumber6, nil
	} else if number == 7 {
		return GasFeeRateOfNumber7, nil
	} else if number == 8 {
		return GasFeeRateOfNumber8, nil
	} else if number == 9 {
		return GasFeeRateOfNumber9, nil
	} else if number == 10 {
		return GasFeeRateOfNumber10, nil
	} else {
		return GasFeeRateOfNumber0, nil
	}
}

// BTC/kB
func EstimateSmartFeeRate(blocks int) (gasFeeRate float64, err error) {
	feeResult, err := api.EstimateSmartFeeAndGetResult(blocks)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "EstimateSmartFeeAndGetResult")
	}
	if feeResult.Errors != nil || feeResult.Blocks != int64(blocks) || *feeResult.FeeRate == 0 {
		err = errors.New("fee result got error or blocks is not same or fee rate is zero")
		return 0, err
	}
	return *feeResult.FeeRate, nil
}

func UpdateFeeRate() {
	err := CheckIfUpdateFeeRateInfo()
	if err != nil {
		return
	}
}

// EstimateSmartFeeRateSatPerKw
// @Note: sat/kw
// @Description: Need UpdateFeeRate first
func EstimateSmartFeeRateSatPerKw() (estimatedFeeSatPerKw int, err error) {
	estimatedFee, err := GetEstimateSmartFeeRate()
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "GetEstimateSmartFeeRate")
	}
	estimatedFeeSatPerKw = FeeRateBtcPerKbToSatPerKw(estimatedFee)
	return estimatedFeeSatPerKw, nil
}

// Need UpdateFeeRate first
func EstimateSmartFeeRateSatPerB() (estimatedFeeSatPerB int, err error) {
	var estimatedFeeBtcPerKb float64
	estimatedFeeBtcPerKb, err = GetEstimateSmartFeeRate()
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "GetEstimateSmartFeeRate")
	}
	estimatedFeeSatPerB = int(estimatedFeeBtcPerKb * 1e5)
	return estimatedFeeSatPerB, nil
}

// Need UpdateFeeRate first
func EstimateSmartFeeRateBtcPerKb() (estimatedFeeBtcPerKb float64, err error) {
	return GetEstimateSmartFeeRate()
}

// FeeRateBtcPerKbToSatPerKw
// @Description: BTC/Kb to sat/kw
// 1 sat/vB = 0.25 sat/wu
// https://bitcoin.stackexchange.com/questions/106333/different-fee-rate-units-sat-vb-sat-perkw-sat-perkb
func FeeRateBtcPerKbToSatPerKw(btcPerKb float64) (satPerKw int) {
	// @dev: 1 BTC/kB = 1e8 sat/kB 1e5 sat/B = 0.25e5 sat/w = 0.25e8 sat/kw
	return int(0.25e8 * btcPerKb)
}

// FeeRateBtcPerKbToSatPerB
// @Description: BTC/Kb to sat/b
func FeeRateBtcPerKbToSatPerB(btcPerKb float64) (satPerB int) {
	return int(1e5 * btcPerKb)
}

// FeeRateSatPerKwToBtcPerKb
// @Description: sat/kw to BTC/Kb
func FeeRateSatPerKwToBtcPerKb(feeRateSatPerKw int) (feeRateBtcPerKb float64) {
	return utils.RoundToDecimalPlace(float64(feeRateSatPerKw)/0.25e8, 8)
}

// FeeRateSatPerKwToSatPerB
// @Description: sat/kw to sat/b
func FeeRateSatPerKwToSatPerB(feeRateSatPerKw int) (feeRateSatPerB int) {
	return feeRateSatPerKw * 4 / 1000
}

// FeeRateSatPerBToBtcPerKb
// @Description: sat/b to BTC/Kb
func FeeRateSatPerBToBtcPerKb(feeRateSatPerB int) (feeRateBtcPerKb float64) {
	return utils.RoundToDecimalPlace(float64(feeRateSatPerB)/100000, 8)
}

// FeeRateSatPerBToSatPerKw
// @Description: sat/b to sat/kw
func FeeRateSatPerBToSatPerKw(feeRateSatPerB int) (feeRateSatPerKw int) {
	return feeRateSatPerB * 1000 / 4
}

// CalculateGasFeeRateSatPerKw
// sat/kw
// Need UpdateFeeRate first
func CalculateGasFeeRateSatPerKw(number int) (int, error) {
	if number <= 0 {
		return 0, errors.New("number to calculate gas fee rate is less equal than zero")
	}
	feeRateSatPerKw, err := EstimateSmartFeeRateSatPerKw()
	rate, err := NumberToGasFeeRate(number)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "NumberToGasFeeRate")
	}
	return int(rate * float64(feeRateSatPerKw)), nil
}

// CalculateGasFeeRateSatPerB
// sat/B
// Need UpdateFeeRate first
func CalculateGasFeeRateSatPerB(number int) (int, error) {
	feeRateSatPerB, err := EstimateSmartFeeRateSatPerB()
	rate, err := NumberToGasFeeRate(number)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "NumberToGasFeeRate")
	}
	return int(rate * float64(feeRateSatPerB)), nil
}

// CalculateGasFeeRateBtcPerKb
// BTC/kB
// Need UpdateFeeRate first
func CalculateGasFeeRateBtcPerKb(number int) (float64, error) {
	feeRateBtcPerKb, err := EstimateSmartFeeRateBtcPerKb()
	rate, err := NumberToGasFeeRate(number)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "NumberToGasFeeRate")
	}
	return rate * feeRateBtcPerKb, nil
}

// GetIssuanceTransactionByteSize
// @dev: Not actual value
func GetIssuanceTransactionByteSize() int {
	// TODO: need to complete
	return GetTapdMintAssetAndFinalizeTransactionByteSize() + GetTapdSendReservedAssetTransactionByteSize()
}

func GetTapdMintAssetAndFinalizeTransactionByteSize() int {
	// TODO: need to complete
	return int(float64(170 * (1 + 0.25)))
}

func GetTapdSendReservedAssetTransactionByteSize() int {
	// TODO: need to complete
	return int(float64(170 * (1 + 0.25)))
}

func GetIssuanceTransactionGasFee(feeRateSatPerKw int) int {
	return FeeRateSatPerKwToSatPerB(feeRateSatPerKw) * GetIssuanceTransactionByteSize()
}

func GetMintTransactionByteSize() int {
	// TODO: need to complete
	return 170
}

func GetMintedTransactionGasFee(feeRateSatPerKw int) int {
	return FeeRateSatPerKwToSatPerB(feeRateSatPerKw) * GetMintTransactionByteSize()
}

func CalculateGasFee(number int, byteSize int) (int, error) {
	calculatedGasFeeRateSatPerB, err := CalculateGasFeeRateSatPerB(number)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "CalculateGasFeeRateSatPerB")
	}
	gasFee := byteSize * calculatedGasFeeRateSatPerB
	return gasFee, nil
}

func IsMintFeePaid(paidId int) bool {
	state, err := CheckPayInsideStatus(uint(paidId))
	if err != nil {
		return false
	}
	return state
}

func IsIssuanceFeePaid(paidId int) bool {
	state, err := CheckPayInsideStatus(uint(paidId))
	if err != nil {
		return false
	}
	return state
}

func PayMintFee(userId int, feeRateSatPerKw int) (mintFeePaidId int, err error) {
	fee := GetMintedTransactionGasFee(feeRateSatPerKw)
	return PayGasFee(userId, fee)
}

func PayIssuanceFee(userId int, feeRateSatPerKw int) (IssuanceFeePaidId int, err error) {
	fee := GetIssuanceTransactionGasFee(feeRateSatPerKw)
	return PayGasFee(userId, fee)
}

func PayGasFee(payUserId int, gasFee int) (int, error) {
	id, err := PayAmountToAdmin(uint(payUserId), uint64(gasFee), 0)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "PayAmountToAdmin")
	}
	return int(id), nil
}

func GetFeeRateInfoByName(name string) (feeRateInfo *models.FeeRateInfo, err error) {
	err = middleware.DB.Where("name = ?", name).First(&feeRateInfo).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "First feeRateInfo")
	}
	return feeRateInfo, nil
}

func GetFeeRateInfoEstimateSmartFeeRateByName(name string) (estimateSmartFeeRate float64, err error) {
	var feeRateInfo *models.FeeRateInfo
	feeRateInfo, err = GetFeeRateInfoByName(name)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "GetFeeRateInfoByName")
	}
	return feeRateInfo.FeeRate, nil
}

func UpdateFeeRateInfoByBitcoind() (err error) {
	var feeRateInfo *models.FeeRateInfo
	var f = FeeRateInfoStore{DB: middleware.DB}
	feeRateInfo, err = GetFeeRateInfoByName(string(GasFeeRateNameBitcoind))
	if err != nil {
		// @dev: Create FeeRateInfo
		feeRateInfo = &models.FeeRateInfo{
			Name: string(GasFeeRateNameBitcoind),
		}
		err = f.CreateFeeRateInfo(feeRateInfo)
		if err != nil {
			return utils.AppendErrorInfo(err, "CreateFeeRateInfo")
		}
		// @dev: Create new record
		FEE.Info("Bitcoind FeeRateInfo record created. %v", err)
	}
	feeRateInfo.FeeRate, err = EstimateSmartFeeRate(config.GetLoadConfig().FairLaunchConfig.EstimateSmartFeeRateBlocks)
	if err != nil {
		return utils.AppendErrorInfo(err, "EstimateSmartFeeRate")
	}
	err = f.UpdateFeeRateInfo(feeRateInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateFeeRateInfo")
	}
	return nil
}

func UpdateFeeRateInfoByBlock(block int) (err error) {
	var feeRateInfo *models.FeeRateInfo
	var f = FeeRateInfoStore{DB: middleware.DB}
	name := strconv.Itoa(block)
	feeRateInfo, err = GetFeeRateInfoByName(name)
	if err != nil {
		// @dev: Create FeeRateInfo
		feeRateInfo = &models.FeeRateInfo{
			Name: name,
		}
		err = f.CreateFeeRateInfo(feeRateInfo)
		if err != nil {
			return utils.AppendErrorInfo(err, "CreateFeeRateInfo")
		}
		FEE.Info("%s %v", name, "FeeRateInfo record created.")
	}
	feeRateInfo.FeeRate, err = EstimateSmartFeeRate(block)
	if err != nil {
		return utils.AppendErrorInfo(err, "EstimateSmartFeeRate")
	}
	err = f.UpdateFeeRateInfo(feeRateInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateFeeRateInfo")
	}
	return nil
}

// CheckIfUpdateFeeRateInfo
// @dev: 1.Update fee rate or not
func CheckIfUpdateFeeRateInfo() (err error) {
	if config.GetLoadConfig().FairLaunchConfig.IsAutoUpdateFeeRate {
		err = UpdateFeeRateInfoByBitcoind()
		if err != nil {
			return utils.AppendErrorInfo(err, "UpdateFeeRateInfoByBitcoind")
		}
	}
	return nil
}

func CheckIfUpdateFeeRateInfoByBlockOfWeek() (err error) {
	if config.GetLoadConfig().FairLaunchConfig.IsAutoUpdateFeeRate {
		for i := 2; i <= 1008; i++ {
			block := i
			err = UpdateFeeRateInfoByBlock(block)
			if err != nil {
				// @dev: Do not return
			}
		}

	}
	return nil
}

func CheckIfUpdateFeeRateInfoByBlockOfDay() (err error) {
	if config.GetLoadConfig().FairLaunchConfig.IsAutoUpdateFeeRate {
		for i := 2; i <= 144; i++ {
			block := i
			err = UpdateFeeRateInfoByBlock(block)
			if err != nil {
				return utils.AppendErrorInfo(err, "UpdateFeeRateInfoByBlock")
			}
		}

	}
	return nil
}

func CheckIfUpdateFeeRateInfoByBlockCustom() (err error) {
	if config.GetLoadConfig().FairLaunchConfig.IsAutoUpdateFeeRate {
		for _, block := range []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 15, 18, 21, 24, 27, 30, 36, 42, 48, 54, 60, 66, 72, 84, 96, 108, 126, 144, 288, 432, 576, 720, 864, 1008} {
			err = UpdateFeeRateInfoByBlock(block)
			if err != nil {
				return utils.AppendErrorInfo(err, "UpdateFeeRateInfoByBlock")
			}
		}

	}
	return nil
}

// GetEstimateSmartFeeRate
// @dev: 2.Get fee rate
func GetEstimateSmartFeeRate() (estimateSmartFeeRate float64, err error) {
	return GetFeeRateInfoEstimateSmartFeeRateByName(string(GasFeeRateNameDefault))
}

type FeeRateResponse struct {
	SatPerKw int     `json:"sat_per_kw"`
	SatPerB  int     `json:"sat_per_b"`
	BtcPerKb float64 `json:"btc_per_kb"`
}

func GetMempoolFeeRate() (*FeeRateResponseTransformed, error) {
	return UpdateAndGetFeeRateResponseTransformed()
}

func GetFeeRate() (*FeeRateResponse, error) {
	UpdateFeeRate()
	var feeRateResponse FeeRateResponse
	var err error
	feeRateResponse.SatPerKw, err = EstimateSmartFeeRateSatPerKw()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "EstimateSmartFeeRateSatPerKw")
	}
	feeRateResponse.SatPerB, err = EstimateSmartFeeRateSatPerB()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "EstimateSmartFeeRateSatPerB")
	}
	feeRateResponse.BtcPerKb, err = EstimateSmartFeeRateBtcPerKb()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "EstimateSmartFeeRateBtcPerKb")
	}
	return &feeRateResponse, nil
}

func GetAllFeeRateInfos() (*[]models.FeeRateInfo, error) {
	var feeRateInfos []models.FeeRateInfo
	err := middleware.DB.Find(&feeRateInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find feeRateInfos")
	}
	return &feeRateInfos, nil
}

type MempoolFeeRate struct {
	FastestFee  int `json:"fastest_fee"`
	HalfHourFee int `json:"half_hour_fee"`
	HourFee     int `json:"hour_fee"`
	EconomyFee  int `json:"economy_fee"`
	MinimumFee  int `json:"minimum_fee"`
}

type FeeRateResponseTransformed struct {
	SatPerB  MempoolFeeRate
	SatPerKw MempoolFeeRate
}

func GetFeeRateResponseTransformedByMempool() (*FeeRateResponseTransformed, error) {
	fees, err := api.MempoolGetRecommendedFees()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "MempoolGetRecommendedFees")
	}
	return &FeeRateResponseTransformed{
		SatPerB: MempoolFeeRate{
			FastestFee:  fees.FastestFee,
			HalfHourFee: fees.HalfHourFee,
			HourFee:     fees.HourFee,
			EconomyFee:  fees.EconomyFee,
			MinimumFee:  fees.MinimumFee,
		},
		SatPerKw: MempoolFeeRate{
			FastestFee:  FeeRateSatPerBToSatPerKw(fees.FastestFee),
			HalfHourFee: FeeRateSatPerBToSatPerKw(fees.HalfHourFee),
			HourFee:     FeeRateSatPerBToSatPerKw(fees.HourFee),
			EconomyFee:  FeeRateSatPerBToSatPerKw(fees.EconomyFee),
			MinimumFee:  FeeRateSatPerBToSatPerKw(fees.MinimumFee),
		},
	}, nil
}

func UpdateFeeRateByMempool() {
	err := CheckIfUpdateFeeRateInfoByMempool()
	if err != nil {
		return
	}
}

func CheckIfUpdateFeeRateInfoByMempool() (err error) {
	if config.GetLoadConfig().FairLaunchConfig.IsAutoUpdateFeeRate {
		err = UpdateFeeRateInfoByMempool()
		if err != nil {
			return utils.AppendErrorInfo(err, "UpdateFeeRateInfoByMempool")
		}
	}
	return nil
}

func GetFeeRateResponseTransformed() (*FeeRateResponseTransformed, error) {
	var feeRateInfos []models.FeeRateInfo
	units := []models.FeeRateType{models.FeeRateTypeSatPerB, models.FeeRateTypeSatPerKw}
	names := []string{"fastest_fee", "half_hour_fee", "hour_fee", "economy_fee", "minimum_fee"}
	for _, unit := range units {
		for _, name := range names {
			feeRateInfo, err := GetFeeRateInfoByNameAndUnit(name, unit)
			if err != nil {
				// @dev: Do not return
			}
			feeRateInfos = append(feeRateInfos, *feeRateInfo)
		}
	}
	transformed, err := ProcessFeeRateInfosToResponseTransformed(feeRateInfos)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "ProcessFeeRateInfosToResponseTransformed")
	}
	return transformed, nil
}

func ProcessFeeRateInfosToResponseTransformed(feeRateInfos []models.FeeRateInfo) (*FeeRateResponseTransformed, error) {
	var feeRateResponseTransformed FeeRateResponseTransformed
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

func GetFeeRateInfoByNameAndUnit(name string, unit models.FeeRateType) (*models.FeeRateInfo, error) {
	var feeRateInfo models.FeeRateInfo
	err := middleware.DB.Where("name = ? AND unit = ?", name, unit).First(&feeRateInfo).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "First feeRateInfo")
	}
	return &feeRateInfo, nil
}

func UpdateFeeRateInfoByNameAndUnitIfNotExistThenCreate(name string, unit models.FeeRateType, feeRate int) error {
	f := FeeRateInfoStore{DB: middleware.DB}
	feeRateInfo, err := GetFeeRateInfoByNameAndUnit(name, unit)
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
		FEE.Info("%v %v %v", name, unit, "FeeRateInfo record created.")
	} else {
		feeRateInfo.FeeRate = float64(feeRate)
		err = f.UpdateFeeRateInfo(feeRateInfo)
		if err != nil {
			return utils.AppendErrorInfo(err, "UpdateFeeRateInfo")
		}
	}
	return nil
}

func (feeRate *FeeRateResponseTransformed) GetFeeRateByNameAndUnit(unit models.FeeRateType, name string) (rate int, err error) {
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

func UpdateFeeRateInfoByMempool() error {
	//var feeRateInfos []models.FeeRateInfo
	feeRateResponse, err := GetFeeRateResponseTransformedByMempool()
	if err != nil {
		return utils.AppendErrorInfo(err, "GetFeeRateResponseTransformedByMempool")
	}
	units := []models.FeeRateType{models.FeeRateTypeSatPerB, models.FeeRateTypeSatPerKw}
	names := []string{"fastest_fee", "half_hour_fee", "hour_fee", "economy_fee", "minimum_fee"}
	for _, unit := range units {
		for _, name := range names {
			var feeRate int
			feeRate, err = feeRateResponse.GetFeeRateByNameAndUnit(unit, name)
			err = UpdateFeeRateInfoByNameAndUnitIfNotExistThenCreate(name, unit, feeRate)
			if err != nil {
				// @dev: Do not return
			}
		}
	}
	return nil
}
