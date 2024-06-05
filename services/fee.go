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
	GasFeeRateOfNumber0                     float64         = 1
	GasFeeRateOfNumber1                                     = 1.2
	GasFeeRateOfNumber2                     float64         = 2
	GasFeeRateOfNumber3                     float64         = 3
	GasFeeRateOfNumber4                     float64         = 4
	GasFeeRateOfNumber5                                     = 4.5
	GasFeeRateOfNumber6                     float64         = 5
	GasFeeRateOfNumber7                                     = 5.4
	GasFeeRateOfNumber8                                     = 5.7
	GasFeeRateOfNumber9                     float64         = 6
	GasFeeRateOfNumber10                                    = 6.5
	GasFeeRateNameMempoolMainnetRecommended FeeRateInfoName = "mempool_mainnet_recommended"
	GasFeeRateNameBitcoind                  FeeRateInfoName = "bitcoind"
	GasFeeRateNameDefault                                   = GasFeeRateNameBitcoind
)

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
		//FEE.Info("Estimate SmartFee And GetResult %v", err)
		return 0, err
	}
	if feeResult.Errors != nil || feeResult.Blocks != int64(blocks) || *feeResult.FeeRate == 0 {
		err = errors.New("fee result got error or blocks is not same or fee rate is zero")
		//FEE.Info("Invalid fee rate result %v", err)
		return 0, err
	}
	return *feeResult.FeeRate, nil
}

func UpdateFeeRate() {
	err := CheckIfUpdateFeeRateInfo()
	if err != nil {
		//FEE.Info("Check If Update FeeRateInfo %v", err)
		return
	}
}

// EstimateSmartFeeRateSatPerKw
// @Note: sat/kw
// @Description: Need UpdateFeeRate first
// @return estimatedFeeSatPerKw
// @return err
func EstimateSmartFeeRateSatPerKw() (estimatedFeeSatPerKw int, err error) {
	estimatedFee, err := GetEstimateSmartFeeRate()
	if err != nil {
		//FEE.Info("Estimate Smart FeeRate %v", err)
		return 0, err
	}
	estimatedFeeSatPerKw = FeeRateBtcPerKbToSatPerKw(estimatedFee)
	return estimatedFeeSatPerKw, nil
}

// Need UpdateFeeRate first
func EstimateSmartFeeRateSatPerB() (estimatedFeeSatPerB int, err error) {
	var estimatedFeeBtcPerKb float64
	estimatedFeeBtcPerKb, err = GetEstimateSmartFeeRate()
	if err != nil {
		return 0, err
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
// @param btcPerKb
// @return satPerB
func FeeRateBtcPerKbToSatPerB(btcPerKb float64) (satPerB int) {
	return int(1e5 * btcPerKb)
}

// FeeRateSatPerKwToBtcPerKb
// @Description: sat/kw to BTC/Kb
// @param feeRateSatPerKw
// @return feeRateBtcPerKb
func FeeRateSatPerKwToBtcPerKb(feeRateSatPerKw int) (feeRateBtcPerKb float64) {
	return utils.RoundToDecimalPlace(float64(feeRateSatPerKw)/0.25e8, 8)
}

// FeeRateSatPerKwToSatPerB
// @Description: sat/kw to sat/b
// @param feeRateSatPerKw
// @return feeRateSatPerB
func FeeRateSatPerKwToSatPerB(feeRateSatPerKw int) (feeRateSatPerB int) {
	return feeRateSatPerKw * 4 / 1000
}

// FeeRateSatPerBToBtcPerKb
// @Description: sat/b to BTC/Kb
// @param feeRateSatPerB
// @return feeRateBtcPerKb
func FeeRateSatPerBToBtcPerKb(feeRateSatPerB int) (feeRateBtcPerKb float64) {
	return utils.RoundToDecimalPlace(float64(feeRateSatPerB)/100000, 8)
}

// FeeRateSatPerBToSatPerKw
// @Description: sat/b to sat/kw
// @param feeRateSatPerB
// @return feeRateSatPerKw
func FeeRateSatPerBToSatPerKw(feeRateSatPerB int) (feeRateSatPerKw int) {
	return feeRateSatPerB * 1000 / 4
}

// sat/kw
// Need UpdateFeeRate first
func CalculateGasFeeRateSatPerKw(number int) (int, error) {
	if number <= 0 {
		return 0, errors.New("number to calculate gas fee rate is less equal than zero")
	}
	feeRateSatPerKw, err := EstimateSmartFeeRateSatPerKw()
	rate, err := NumberToGasFeeRate(number)
	if err != nil {
		//FEE.Info("Number To Gas FeeRate %v", err)
		return 0, err
	}
	return int(rate * float64(feeRateSatPerKw)), nil
}

// sat/B
// Need UpdateFeeRate first
func CalculateGasFeeRateSatPerB(number int) (int, error) {
	feeRateSatPerB, err := EstimateSmartFeeRateSatPerB()
	rate, err := NumberToGasFeeRate(number)
	if err != nil {
		//FEE.Info("Number To Gas FeeRate %v", err)
		return 0, err
	}
	return int(rate * float64(feeRateSatPerB)), nil
}

// BTC/kB
// Need UpdateFeeRate first
func CalculateGasFeeRateBtcPerKb(number int) (float64, error) {
	feeRateBtcPerKb, err := EstimateSmartFeeRateBtcPerKb()
	rate, err := NumberToGasFeeRate(number)
	if err != nil {
		//FEE.Info("Number To Gas FeeRate %v", err)
		return 0, err
	}
	return rate * feeRateBtcPerKb, nil
}

// @dev: not actual value
func GetTransactionByteSize() int {
	// TODO: need to complete
	return 170
}

func CalculateGasFee(number int, byteSize int) (int, error) {
	calculatedGasFeeRateSatPerB, err := CalculateGasFeeRateSatPerB(number)
	if err != nil {
		//FEE.Info("Calculate GasFeeRate SatPerB %v", err)
		return 0, err
	}
	gasFee := byteSize * calculatedGasFeeRateSatPerB
	return gasFee, nil
}

func IsMintFeePaid(paidId int) bool {
	state, err := CheckPayInsideStatus(uint(paidId))
	if err != nil {
		//FEE.Info("GetBalance %v", err)
		return false
	}
	return state
}

func IsIssuanceFeePaid(paidId int) bool {
	state, err := CheckPayInsideStatus(uint(paidId))
	if err != nil {
		//FEE.Info("GetBalance %v", err)
		return false
	}
	return state
}

func PayMintFee(userId int, feeRateSatPerKw int) (mintFeePaidId int, err error) {
	fee := FeeRateSatPerKwToSatPerB(feeRateSatPerKw) * GetTransactionByteSize()
	return PayGasFee(userId, fee)
}

func PayIssuanceFee(userId int, feeRateSatPerKw int) (IssuanceFeePaidId int, err error) {
	// TODO: User need to pay more fee
	fee := FeeRateSatPerKwToSatPerB(feeRateSatPerKw) * GetTransactionByteSize()
	return PayGasFee(userId, fee)
}

func PayGasFee(payUserId int, gasFee int) (int, error) {
	id, err := PayAmountToAdmin(uint(payUserId), uint64(gasFee), 0)
	return int(id), err
}

func GetFeeRateInfoByName(name string) (feeRateInfo *models.FeeRateInfo, err error) {
	err = middleware.DB.Where("name = ?", name).First(&feeRateInfo).Error
	if err != nil {
		//FEE.Info("Find FeeRateInfo %v", err)
		return nil, err
	}
	return feeRateInfo, nil
}

func GetFeeRateInfoEstimateSmartFeeRateByName(name string) (estimateSmartFeeRate float64, err error) {
	var feeRateInfo *models.FeeRateInfo
	feeRateInfo, err = GetFeeRateInfoByName(name)
	if err != nil {
		//FEE.Info("Get FeeRateInfo By Name %v", err)
		return 0, err
	}
	return feeRateInfo.FeeRate, nil
}

func UpdateFeeRateInfoByBitcoind() (err error) {
	var feeRateInfo *models.FeeRateInfo
	var f = FeeRateInfoStore{DB: middleware.DB}
	feeRateInfo, err = GetFeeRateInfoByName(string(GasFeeRateNameBitcoind))
	if err != nil {
		//FEE.Info("Get FeeRateInfo By Bitcoind, Create now. %v", err)
		//	Create FeeRateInfo
		feeRateInfo = &models.FeeRateInfo{
			Name: string(GasFeeRateNameBitcoind),
		}
		err = f.CreateFeeRateInfo(feeRateInfo)
		if err != nil {
			//FEE.Info("Create FeeRate Info %v", err)
			return err
		}
		//@dev: create new record
		FEE.Info("Bitcoind FeeRateInfo record created. %v", err)
	}
	feeRateInfo.FeeRate, err = EstimateSmartFeeRate(config.GetLoadConfig().FairLaunchConfig.EstimateSmartFeeRateBlocks)
	if err != nil {
		//FEE.Info("Estimate Smart FeeRate %v", err)
		return err
	}
	err = f.UpdateFeeRateInfo(feeRateInfo)
	if err != nil {
		//FEE.Info("Update FeeRateInfo %v", err)
		return err
	}
	return nil
}

func UpdateFeeRateInfoByBlock(block int) (err error) {
	var feeRateInfo *models.FeeRateInfo
	var f = FeeRateInfoStore{DB: middleware.DB}
	name := strconv.Itoa(block)
	feeRateInfo, err = GetFeeRateInfoByName(name)
	if err != nil {
		//FEE.Info("%s %v %s", "Get FeeRateInfo By Block", block, "Create now.")
		//	Create FeeRateInfo
		feeRateInfo = &models.FeeRateInfo{
			Name: name,
		}
		err = f.CreateFeeRateInfo(feeRateInfo)
		if err != nil {
			//FEE.Info("Create FeeRate Info %v", err)
			return err
		}
		FEE.Info("%s %v", name, "FeeRateInfo record created.")
	}
	feeRateInfo.FeeRate, err = EstimateSmartFeeRate(block)
	if err != nil {
		//FEE.Info("Estimate Smart FeeRate %v", err)
		return err
	}
	err = f.UpdateFeeRateInfo(feeRateInfo)
	if err != nil {
		//FEE.Info("Update FeeRateInfo %v", err)
		return err
	}
	return nil
}

// @dev: 1.update fee rate or not
func CheckIfUpdateFeeRateInfo() (err error) {
	if config.GetLoadConfig().FairLaunchConfig.IsAutoUpdateFeeRate {
		err = UpdateFeeRateInfoByBitcoind()
		if err != nil {
			//FEE.Info("Update FeeRateInfo By Bitcoind %v", err)
			return err
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
				//FEE.Info("Update FeeRateInfo By %v %v", block, err)
				// @dev: do not return
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
				//FEE.Info("Update FeeRateInfo By %v %v", block, err)
				return err
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
				//FEE.Info("Update FeeRateInfo By %v %v", block, err)
				return err
			}
		}

	}
	return nil
}

// @dev: 2.get fee rate
func GetEstimateSmartFeeRate() (estimateSmartFeeRate float64, err error) {
	return GetFeeRateInfoEstimateSmartFeeRateByName(string(GasFeeRateNameDefault))
}

type FeeRateResponse struct {
	SatPerKw int     `json:"sat_per_kw"`
	SatPerB  int     `json:"sat_per_b"`
	BtcPerKb float64 `json:"btc_per_kb"`
}

func GetFeeRate() (*FeeRateResponse, error) {
	UpdateFeeRate()
	var feeRateResponse FeeRateResponse
	var err error
	feeRateResponse.SatPerKw, err = EstimateSmartFeeRateSatPerKw()
	if err != nil {
		//FEE.Info("Estimate Smart FeeRate SatPerKw %v", err)
		return nil, err
	}
	feeRateResponse.SatPerB, err = EstimateSmartFeeRateSatPerB()
	if err != nil {
		//FEE.Info("Estimate Smart FeeRate SatPerB %v", err)
		return nil, err
	}
	feeRateResponse.BtcPerKb, err = EstimateSmartFeeRateBtcPerKb()
	if err != nil {
		//FEE.Info("Estimate Smart FeeRate BtcPerKb %v", err)
		return nil, err
	}
	return &feeRateResponse, nil
}

func GetAllFeeRateInfos() (*[]models.FeeRateInfo, error) {
	var feeRateInfos []models.FeeRateInfo
	err := middleware.DB.Find(&feeRateInfos).Error
	return &feeRateInfos, err
}
