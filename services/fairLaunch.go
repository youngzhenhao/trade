package services

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"gorm.io/gorm"
	"math"
	"reflect"
	"sort"
	"strconv"
	"time"
	"trade/api"
	"trade/btlLog"
	"trade/middleware"
	"trade/models"
	"trade/services/btldb"
	"trade/services/custodyAccount"
	"trade/utils"
)

func PrintProcessionResult(processionResult *[]ProcessionResult) {
	for _, result := range *processionResult {
		if !result.Success {
			btlLog.FairLaunchDebugLogger.Info("%d:%v", result.Id, utils.ValueJsonString(result.Error))
		}
	}
}

// FairLaunchIssuance
// @Description: Scheduled Task
func FairLaunchIssuance(tx *gorm.DB) {
	processionResult, err := ProcessAllFairLaunchInfos(tx)
	if err != nil {
		return
	}
	if processionResult == nil || len(*processionResult) == 0 {
		return
	}
	PrintProcessionResult(processionResult)
}

// @dev: Process by state

// ProcessFairLaunchNoPay
// @Description: Scheduled Task
func ProcessFairLaunchNoPay(tx *gorm.DB) {
	processionResult, err := ProcessAllFairLaunchStateNoPayInfoService(tx)
	if err != nil {
		return
	}
	if processionResult == nil || len(*processionResult) == 0 {
		return
	}
	PrintProcessionResult(processionResult)
}

// ProcessFairLaunchPaidPending
// @Description: Scheduled Task
func ProcessFairLaunchPaidPending(tx *gorm.DB) {
	processionResult, err := ProcessAllFairLaunchStatePaidPendingInfoService(tx)
	if err != nil {
		return
	}
	if processionResult == nil || len(*processionResult) == 0 {
		return
	}
	PrintProcessionResult(processionResult)
}

// ProcessFairLaunchPaidNoIssue
// @Description: Scheduled Task
func ProcessFairLaunchPaidNoIssue(tx *gorm.DB) {
	processionResult, err := ProcessAllFairLaunchStatePaidNoIssueInfoService(tx)
	if err != nil {
		return
	}
	if processionResult == nil || len(*processionResult) == 0 {
		return
	}
	PrintProcessionResult(processionResult)
}

// ProcessFairLaunchIssuedPending
// @Description: Scheduled Task
func ProcessFairLaunchIssuedPending(tx *gorm.DB) {
	processionResult, err := ProcessAllFairLaunchStateIssuedPendingInfoService(tx)
	if err != nil {
		return
	}
	if processionResult == nil || len(*processionResult) == 0 {
		return
	}
	PrintProcessionResult(processionResult)
}

// ProcessFairLaunchReservedSentPending
// @Description: Scheduled Task
func ProcessFairLaunchReservedSentPending(tx *gorm.DB) {
	processionResult, err := ProcessAllFairLaunchStateReservedSentPending(tx)
	if err != nil {
		return
	}
	if processionResult == nil || len(*processionResult) == 0 {
		return
	}
	PrintProcessionResult(processionResult)
}

// FairLaunchMint
// @Description: Scheduled Task
func FairLaunchMint(tx *gorm.DB) {
	processionResult, err := ProcessAllFairLaunchMintedInfos(tx)
	if err != nil {
		return
	}
	if processionResult == nil || len(*processionResult) == 0 {
		return
	}
	PrintProcessionResult(processionResult)
}

// SendFairLaunchAsset
// @Description: Scheduled Task
func SendFairLaunchAsset(tx *gorm.DB) {
	err := SendFairLaunchMintedAssetLocked(tx)
	if err != nil {
		err = utils.WriteToLogFile("./trade.SendFairLaunchAsset.log", "[TRADE.SFLA]", utils.ValueJsonString(err))
		if err != nil {
			utils.LogError("Write SendFairLaunchAsset err to log file", err)
		}
		btlLog.FairLaunchDebugLogger.Info("SendFairLaunchAsset: %v", err)
		return
	}
}

// RemoveMintedInventories
// @Description: Scheduled Task
func RemoveMintedInventories() {
	err := RemoveFairLaunchInventoryStateMintedInfos()
	if err != nil {
		btlLog.FairLaunchDebugLogger.Info("%v", err)
		return
	}
}

func GetAllFairLaunchInfos() (*[]models.FairLaunchInfo, error) {
	f := btldb.FairLaunchStore{DB: middleware.DB}
	var fairLaunchInfos []models.FairLaunchInfo
	err := f.DB.Find(&fairLaunchInfos).Error
	return &fairLaunchInfos, err
}

func GetFairLaunchInfo(id int) (*models.FairLaunchInfo, error) {
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.ReadFairLaunchInfo(uint(id))
}

func GetFairLaunchMintedInfo(id int) (*models.FairLaunchMintedInfo, error) {
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.ReadFairLaunchMintedInfo(uint(id))
}

func GetFairLaunchMintedInfoWhoseProcessNumberIsMoreThanTenThousand() (*[]models.FairLaunchMintedInfo, error) {
	return btldb.ReadFairLaunchMintedInfoWhoseProcessNumberIsMoreThanTenThousand()
}

func GetFairLaunchMintedInfosByFairLaunchId(fairLaunchId int) (*[]models.FairLaunchMintedInfo, error) {
	f := btldb.FairLaunchStore{DB: middleware.DB}
	var fairLaunchMintedInfos []models.FairLaunchMintedInfo
	//err := f.DB.Where("fair_launch_info_id = ?", int(uint(id))).Find(&fairLaunchMintedInfos).Error
	err := f.DB.Where(&models.FairLaunchMintedInfo{FairLaunchInfoID: int(uint(fairLaunchId))}).Find(&fairLaunchMintedInfos).Error
	return &fairLaunchMintedInfos, err
}

func SetFairLaunchInfo(fairLaunchInfo *models.FairLaunchInfo) error {
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.CreateFairLaunchInfo(fairLaunchInfo)
}

func SetFairLaunchMintedInfo(fairLaunchMintedInfo *models.FairLaunchMintedInfo) error {
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.CreateFairLaunchMintedInfo(fairLaunchMintedInfo)
}

// ProcessFairLaunchInfo
// @Description: Process fairLaunchInfo
func ProcessFairLaunchInfo(imageData string, name string, assetType int, amount int, reserved int, mintQuantity int, startTime int, endTime int, description string, feeRate int, userId int, username string) (*models.FairLaunchInfo, error) {
	if FeeRateSatPerKwToSatPerB(feeRate) > 500 {
		return nil, errors.New("fee rate exceeds max(500)")
	}
	// @dev: Validate fee rate
	feeRateResponse, err := UpdateAndGetFeeRateResponseTransformed()
	if err != nil {
		return nil, err
	}
	feeRateSatPerKw := feeRateResponse.SatPerKw.FastestFee
	// @dev: The allowable fee rate error is minus one sat per b
	if feeRate+FeeRateSatPerBToSatPerKw(1) < feeRateSatPerKw {
		return nil, errors.New("set fee rate not enough, it has changed")
	}
	lowest := feeRateResponse.SatPerKw.MinimumFee
	if feeRate < lowest {
		return nil, errors.New("set fee rate is less than lowest, it may not be confirmed forever")
	}
	err = ValidateStartAndEndTime(startTime, endTime)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "ValidateStartAndEndTime")
	}
	calculateSeparateAmount, err := AmountReservedAndMintQuantityToReservedTotalAndMintTotal(amount, reserved, mintQuantity)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "AmountReservedAndMintQuantityToReservedTotalAndMintTotal")
	}
	var fairLaunchInfo models.FairLaunchInfo
	// @dev: Setting fee rate need to bigger equal than fee rate now
	// @notice: sever do not need to check
	//err = ValidateFeeRate(feeRate)
	//if err != nil {
	//	return nil, utils.AppendErrorInfo(err, "ValidateFeeRate")
	//}
	setGasFee := GetIssuanceTransactionGasFee(feeRate)
	if !custodyAccount.IsAccountBalanceEnoughByUserId(uint(userId), uint64(setGasFee)) {
		return nil, errors.New("account balance not enough to pay issuance gas fee")
	}
	fairLaunchInfo = models.FairLaunchInfo{
		ImageData:              imageData,
		Name:                   name,
		AssetType:              taprpc.AssetType(assetType),
		Amount:                 amount,
		Reserved:               reserved,
		MintQuantity:           mintQuantity,
		StartTime:              startTime,
		EndTime:                endTime,
		Description:            description,
		FeeRate:                feeRate,
		SetGasFee:              setGasFee,
		SetTime:                utils.GetTimestamp(),
		ActualReserved:         calculateSeparateAmount.ActualReserved,
		ReserveTotal:           calculateSeparateAmount.ReserveTotal,
		MintNumber:             calculateSeparateAmount.MintNumber,
		IsFinalEnough:          calculateSeparateAmount.IsFinalEnough,
		FinalQuantity:          calculateSeparateAmount.FinalQuantity,
		MintTotal:              calculateSeparateAmount.MintTotal,
		ActualMintTotalPercent: calculateSeparateAmount.ActualMintTotalPercent,
		CalculationExpression:  calculateSeparateAmount.CalculationExpression,
		UserID:                 userId,
		Username:               username,
		State:                  models.FairLaunchStateNoPay,
	}
	return &fairLaunchInfo, nil
}

func ValidateFeeRate(feeRate int) (err error) {
	feeRateResponse, err := UpdateAndGetFeeRateResponseTransformed()
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateAndGetFeeRateResponseTransformed")
	}
	feeRateSatPerKw := feeRateResponse.SatPerKw.FastestFee
	//estimatedFeeRateSatPerKw, err := UpdateAndEstimateSmartFeeRateSatPerKw()
	if !(feeRate >= feeRateSatPerKw) {
		err = errors.New("setting fee rate need to bigger equal than mempool's recommended fastest fee rate now")
		return utils.AppendErrorInfo(err, "Got: "+strconv.Itoa(feeRate)+", Expected bigger equal than: "+strconv.Itoa(feeRateSatPerKw))
	}
	return nil
}

func ValidateMintedFeeRate(mintedNumber int, mintedFeeRateSatPerKw int) (err error) {
	feeRate, err := UpdateAndCalculateGasFeeRateByMempool(mintedNumber)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateAndCalculateGasFeeRateByMempool")
	}
	// @dev: maybe need to change comparison param
	calculatedFeeRateSatPerKw := feeRate.SatPerKw.FastestFee
	if !(mintedFeeRateSatPerKw >= calculatedFeeRateSatPerKw) {
		err = errors.New("setting minted calculated FeeRate SatPerKw need to bigger equal than calculated fee rate by mempool's recommended fastest fee rate now")
		return utils.AppendErrorInfo(err, "Got: "+strconv.Itoa(mintedFeeRateSatPerKw)+", Expected bigger equal than: "+strconv.Itoa(calculatedFeeRateSatPerKw))
	}
	return nil
}

func ValidateStartAndEndTime(startTime int, endTime int) error {
	now := utils.GetTimestamp()
	if !(startTime >= now-600) {
		return errors.New("start time must be greater than the current time(max time delay 600 seconds.)")
	}
	if !(endTime >= startTime+3600*2) {
		return errors.New("end time should be at least two hour after the start time")
	}
	if !(endTime <= now+3600*24*365) {
		return errors.New("end time cannot be more than one year from the current time")
	}
	return nil
}

// ProcessFairLaunchMintedInfo
// @Description: Process fairLaunchMintedInfo
func ProcessFairLaunchMintedInfo(fairLaunchInfoID int, mintedNumber int, mintedFeeRateSatPerKw int, addr string, userId int, username string) (*models.FairLaunchMintedInfo, error) {
	if FeeRateSatPerKwToSatPerB(mintedFeeRateSatPerKw) > 500 {
		return nil, errors.New("fee rate exceeds max(500)" + "; " + strconv.Itoa(mintedFeeRateSatPerKw))
	}
	// @dev: Validate fee rate
	feeRateResponse, err := UpdateAndCalculateGasFeeRateByMempool(mintedNumber)
	if err != nil {
		return nil, err
	}
	calculatedFeeRateSatPerKw := feeRateResponse.SatPerKw.FastestFee + FeeRateSatPerBToSatPerKw(2)
	// @dev: The allowable fee rate error is minus 100 sat/b
	// @notice: 2024-8-14 16:37:34 This check should be removed
	// @note: Use greater restrictions instead of removing check
	if mintedFeeRateSatPerKw+FeeRateSatPerBToSatPerKw(100) < calculatedFeeRateSatPerKw {
		return nil, errors.New("mint fee rate not enough, it has changed" + "; " + strconv.Itoa(mintedFeeRateSatPerKw))
	}
	if mintedFeeRateSatPerKw < FeeRateSatPerBToSatPerKw(3) {
		return nil, errors.New("mint fee rate not enough, it less than minimum value (3 sat/b)" + "; " + strconv.Itoa(mintedFeeRateSatPerKw))
	}
	numberToGasFeeRate, _ := NumberToGasFeeRate(mintedNumber)
	if mintedFeeRateSatPerKw < int(math.Ceil(float64(FeeRateSatPerBToSatPerKw(1))*numberToGasFeeRate)) {
		return nil, errors.New("mint fee rate not enough, it less than 1 multiple fee rate of minted number" + "; " + strconv.Itoa(mintedFeeRateSatPerKw))
	}
	lowest := feeRateResponse.SatPerKw.MinimumFee
	if mintedFeeRateSatPerKw < lowest {
		return nil, errors.New("set fee rate is less than lowest, it may not be confirmed forever" + "; " + strconv.Itoa(mintedFeeRateSatPerKw))
	}
	var fairLaunchMintedInfo models.FairLaunchMintedInfo
	isFairLaunchMintTimeRight, err := IsFairLaunchMintTimeRight(fairLaunchInfoID)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "IsFairLaunchMintTimeRight")
	}
	if !isFairLaunchMintTimeRight {
		err = errors.New("not valid mint time")
		return nil, err
	}
	decodedAddrInfo, err := api.GetDecodedAddrInfo(addr)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetDecodedAddrInfo")
	}
	var fairLaunchInfo *models.FairLaunchInfo
	fairLaunchInfo, err = GetFairLaunchInfo(fairLaunchInfoID)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetFairLaunchInfo")
	}
	decodedAddrAssetId := hex.EncodeToString(decodedAddrInfo.AssetId)
	if fairLaunchInfo.AssetID != decodedAddrAssetId {
		err = errors.New("decoded addr asset id is not equal fair launch info's asset id")
		return nil, err
	}
	// @dev: Setting fee rate need to bigger equal than calculated fee rate now
	// @notice: sever do not need to check
	//err = ValidateMintedFeeRate(mintedNumber,mintedFeeRateSatPerKw)
	//if err != nil {
	//	return nil, utils.AppendErrorInfo(err, "ValidateMintedFeeRate")
	//}
	// @dev: Update: 2024-8-22 10:27:40
	// @dev: Do not limit total minted
	//isValid, err := IsMintedNumberValid(userId, fairLaunchInfoID, mintedNumber)
	isValid, err := IsEachMintedNumberValid(mintedNumber)
	if err != nil || !isValid {
		return nil, utils.AppendErrorInfo(err, "Is Minted Number Valid")
	}
	mintedGasFee := GetMintedTransactionGasFee(mintedFeeRateSatPerKw)
	if !custodyAccount.IsAccountBalanceEnoughByUserId(uint(userId), uint64(mintedGasFee)) {
		return nil, errors.New("account balance not enough to pay minted gas fee")
	}
	fairLaunchMintedInfo = models.FairLaunchMintedInfo{
		FairLaunchInfoID:      fairLaunchInfoID,
		MintedNumber:          mintedNumber,
		MintedFeeRateSatPerKw: mintedFeeRateSatPerKw,
		MintedGasFee:          mintedGasFee,
		EncodedAddr:           addr,
		UserID:                userId,
		Username:              username,
		AssetID:               hex.EncodeToString(decodedAddrInfo.AssetId),
		AssetName:             fairLaunchInfo.Name,
		AssetType:             int(decodedAddrInfo.AssetType),
		AddrAmount:            int(decodedAddrInfo.Amount),
		ScriptKey:             hex.EncodeToString(decodedAddrInfo.ScriptKey),
		InternalKey:           hex.EncodeToString(decodedAddrInfo.InternalKey),
		TaprootOutputKey:      hex.EncodeToString(decodedAddrInfo.TaprootOutputKey),
		ProofCourierAddr:      decodedAddrInfo.ProofCourierAddr,
		MintedSetTime:         utils.GetTimestamp(),
		State:                 models.FairLaunchMintedStateNoPay,
	}
	return &fairLaunchMintedInfo, nil
}

type CalculateSeparateAmount struct {
	Amount                 int
	Reserved               int
	ActualReserved         float64
	ReserveTotal           int
	MintQuantity           int
	MintNumber             int
	IsFinalEnough          bool
	FinalQuantity          int
	MintTotal              int
	ActualMintTotalPercent float64
	CalculationExpression  string
}

// AmountReservedAndMintQuantityToReservedTotalAndMintTotal
// @Description: return Calculated result struct
func AmountReservedAndMintQuantityToReservedTotalAndMintTotal(amount int, reserved int, mintQuantity int) (*CalculateSeparateAmount, error) {
	if amount <= 0 || reserved <= 0 || mintQuantity <= 0 {
		return nil, errors.New("amount reserved and mint amount must be greater than zero")
	}
	if reserved > 99 {
		return nil, errors.New("reserved amount must be less equal than 99")
	}
	if amount <= mintQuantity {
		return nil, errors.New("amount must be greater than mint quantity")
	}
	reservedTotal := int(math.Ceil(float64(amount) * float64(reserved) / 100))
	mintTotal := amount - reservedTotal
	remainder := mintTotal % mintQuantity
	var finalQuantity int
	var isFinalEnough bool
	if remainder == 0 {
		isFinalEnough = true
		finalQuantity = mintQuantity
	} else {
		isFinalEnough = false
		finalQuantity = remainder
	}
	if mintTotal <= 0 || mintTotal < mintQuantity {
		return nil, errors.New("insufficient mint total amount")
	}
	reservedTotal = amount - mintTotal
	if reservedTotal <= 0 {
		return nil, errors.New("reserved amount is less equal than zero")
	}

	mintNumber := int(math.Ceil(float64(mintTotal) / float64(mintQuantity)))
	if mintNumber <= 0 {
		return nil, errors.New("mint number is less equal than zero")
	}
	actualReserved := float64(reservedTotal) * 100 / float64(amount)
	actualReserved = utils.RoundToDecimalPlace(actualReserved, 8)
	actualMintTotalPercent := 100 - actualReserved
	calculatedSeparateAmount := CalculateSeparateAmount{
		Amount:                 amount,
		Reserved:               reserved,
		ActualReserved:         actualReserved,
		ReserveTotal:           reservedTotal,
		MintQuantity:           mintQuantity,
		MintNumber:             mintNumber,
		IsFinalEnough:          isFinalEnough,
		FinalQuantity:          finalQuantity,
		MintTotal:              mintTotal,
		ActualMintTotalPercent: actualMintTotalPercent,
	}
	var err error
	calculatedSeparateAmount.CalculationExpression, err = CalculationExpressionBySeparateAmount(&calculatedSeparateAmount)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "CalculationExpressionBySeparateAmount")
	}
	return &calculatedSeparateAmount, nil
}

// CalculationExpressionBySeparateAmount
// @Description: Generate Calculation Expression By Separate Amount
func CalculationExpressionBySeparateAmount(calculateSeparateAmount *CalculateSeparateAmount) (string, error) {
	calculated := calculateSeparateAmount.ReserveTotal + calculateSeparateAmount.MintQuantity*(calculateSeparateAmount.MintNumber-1) + calculateSeparateAmount.FinalQuantity
	if reflect.DeepEqual(calculated, calculateSeparateAmount.Amount) {
		return fmt.Sprintf("%d+%d*%d+%d=%d", calculateSeparateAmount.ReserveTotal, calculateSeparateAmount.MintQuantity, calculateSeparateAmount.MintNumber-1, calculateSeparateAmount.FinalQuantity, calculated), nil
	}
	return "", errors.New("calculated result is not equal amount")
}

// CreateInventoryInfoByFairLaunchInfo
// @Description: Create Inventory Info By FairLaunchInfo
func CreateInventoryInfoByFairLaunchInfo(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) error {
	var FairLaunchInventoryInfos []models.FairLaunchInventoryInfo
	items := fairLaunchInfo.MintNumber - 1
	for ; items > 0; items -= 1 {
		FairLaunchInventoryInfos = append(FairLaunchInventoryInfos, models.FairLaunchInventoryInfo{
			FairLaunchInfoID: int(fairLaunchInfo.ID),
			Quantity:         fairLaunchInfo.MintQuantity,
			State:            models.FairLaunchInventoryStateOpen,
		})
	}
	FairLaunchInventoryInfos = append(FairLaunchInventoryInfos, models.FairLaunchInventoryInfo{
		FairLaunchInfoID: int(fairLaunchInfo.ID),
		Quantity:         fairLaunchInfo.FinalQuantity,
	})
	return CreateFairLaunchInventoryInfosBatchProcess(tx, &FairLaunchInventoryInfos)
}

// CreateAssetIssuanceInfoByFairLaunchInfo
// @Description: Create Asset Issuance Info By FairLaunchInfo
func CreateAssetIssuanceInfoByFairLaunchInfo(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) error {
	assetIssuance := models.AssetIssuance{
		AssetName:      fairLaunchInfo.Name,
		AssetId:        fairLaunchInfo.AssetID,
		AssetType:      fairLaunchInfo.AssetType,
		IssuanceUserId: fairLaunchInfo.UserID,
		IssuanceTime:   utils.GetTimestamp(),
		IsFairLaunch:   true,
		FairLaunchID:   int(fairLaunchInfo.ID),
		State:          models.AssetIssuanceStatePending,
	}
	a := btldb.AssetIssuanceStore{DB: middleware.DB}
	return a.CreateAssetIssuance(tx, &assetIssuance)
}

// GetAllInventoryInfoByFairLaunchInfoId
// @Description: Query all inventory by FairLaunchInfo id
func GetAllInventoryInfoByFairLaunchInfoId(fairLaunchInfoId int) (*[]models.FairLaunchInventoryInfo, error) {
	var fairLaunchInventoryInfos []models.FairLaunchInventoryInfo
	err := middleware.DB.Where("fair_launch_info_id = ?", fairLaunchInfoId).Find(&fairLaunchInventoryInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchInventoryInfos")
	}
	return &fairLaunchInventoryInfos, nil
}

// GetInventoryCouldBeMintedByFairLaunchInfoId
// @Description: Get all Inventory Could Be Minted By FairLaunchInfoId
func GetInventoryCouldBeMintedByFairLaunchInfoId(fairLaunchInfoId int) (*[]models.FairLaunchInventoryInfo, error) {
	var fairLaunchInventoryInfos []models.FairLaunchInventoryInfo
	err := middleware.DB.Where("fair_launch_info_id = ? AND is_minted = ? AND state = ?", fairLaunchInfoId, false, models.FairLaunchInventoryStateOpen).Find(&fairLaunchInventoryInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchInventoryInfos")
	}
	return &fairLaunchInventoryInfos, nil
}

func CalculateInventoryAmount(fairLaunchInventoryInfos *[]models.FairLaunchInventoryInfo) (amount int) {
	if fairLaunchInventoryInfos == nil {
		return 0
	}
	for _, fairLaunchInventoryInfo := range *(fairLaunchInventoryInfos) {
		amount += fairLaunchInventoryInfo.Quantity
	}
	return amount
}

type InventoryNumberAndAmount struct {
	Number int `json:"number"`
	Amount int `json:"amount"`
}

// GetNumberAndAmountOfInventoryCouldBeMinted
// @Description: call GetInventoryCouldBeMintedByFairLaunchInfoId
func GetNumberAndAmountOfInventoryCouldBeMinted(fairLaunchInfoId int) (*InventoryNumberAndAmount, error) {
	fairLaunchInventoryInfos, err := GetInventoryCouldBeMintedByFairLaunchInfoId(fairLaunchInfoId)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetInventoryCouldBeMintedByFairLaunchInfoId")
	}
	amount := CalculateInventoryAmount(fairLaunchInventoryInfos)
	return &InventoryNumberAndAmount{
		Number: len(*fairLaunchInventoryInfos),
		Amount: amount,
	}, err
}

func GetAmountOfInventoryCouldBeMintedByMintedNumber(fairLaunchInfoId int, mintedNumber int) (int, error) {
	fairLaunchInventoryInfos, err := GetInventoryCouldBeMintedByFairLaunchInfoId(fairLaunchInfoId)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "GetInventoryCouldBeMintedByFairLaunchInfoId")
	}
	if mintedNumber > len(*fairLaunchInventoryInfos) {
		err = errors.New("not enough inventory could be minted")
		return 0, err
	}
	fairLaunchInventoryInfoSlice := (*fairLaunchInventoryInfos)[0:mintedNumber]
	amount := CalculateInventoryAmount(&fairLaunchInventoryInfoSlice)
	return amount, nil
}

// IsMintAvailable
// @Description: Is Mint Available by fairLaunchInfoId and number
func IsMintAvailable(fairLaunchInfoId int, number int) bool {
	if !IsFairLaunchIssued(fairLaunchInfoId) {
		//err := errors.New("fairLaunch is not Issued")
		return false
	}
	inventoryNumberAndAmount, err := GetNumberAndAmountOfInventoryCouldBeMinted(fairLaunchInfoId)
	if err != nil {
		return false
	}
	return inventoryNumberAndAmount.Number >= number
}

// GetMintAmountByFairLaunchMintNumber
// @Description: Get Mint Amount By FairLaunch id and MintNumber
func GetMintAmountByFairLaunchMintNumber(fairLaunchInfoId int, number int) (amount int, err error) {
	if number <= 0 {
		err = errors.New("mint number must be greater than zero")
		return 0, err
	}
	fairLaunchInventoryInfos, err := GetInventoryCouldBeMintedByFairLaunchInfoId(fairLaunchInfoId)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "GetInventoryCouldBeMintedByFairLaunchInfoId")
	}
	allNum := len(*fairLaunchInventoryInfos)
	if allNum < number {
		err = errors.New("not enough mint amount")
		return 0, err
	}
	mintInventoryInfos := (*fairLaunchInventoryInfos)[:number]
	for _, inventory := range mintInventoryInfos {
		amount += inventory.Quantity
	}
	return amount, nil
}

// LockInventoryByFairLaunchMintedIdAndMintNumber
// @Description: Calculate MintAmount By id and MintNumber, then Update State, this function will lock inventory
func LockInventoryByFairLaunchMintedIdAndMintNumber(fairLaunchMintedInfoId int, number int) (*[]models.FairLaunchInventoryInfo, error) {
	if number <= 0 {
		err := errors.New("mint number must be greater than zero")
		return nil, err
	}
	fairLaunchMintedInfo, err := GetFairLaunchMintedInfo(fairLaunchMintedInfoId)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetFairLaunchMintedInfo")
	}
	fairLaunchInventoryInfos, err := GetInventoryCouldBeMintedByFairLaunchInfoId(fairLaunchMintedInfo.FairLaunchInfoID)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetInventoryCouldBeMintedByFairLaunchInfoId")
	}
	allNum := len(*fairLaunchInventoryInfos)
	if allNum < number {
		err = errors.New("not enough mint amount")
		return nil, err
	}
	mintInventoryInfos := (*fairLaunchInventoryInfos)[:number]
	//for _, inventory := range mintInventoryInfos {
	//	inventory.Status = models.StatusPending
	//}
	err = middleware.DB.Model(&mintInventoryInfos).Updates(map[string]any{"state": models.FairLaunchInventoryStateLocked, "fair_launch_minted_info_id": fairLaunchMintedInfoId}).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Updates mintInventoryInfos")
	}
	return &mintInventoryInfos, nil
}

// CalculateMintAmountByFairLaunchInventoryInfos
// @Description: Calculate MintAmount By FairLaunchInventoryInfos
func CalculateMintAmountByFairLaunchInventoryInfos(fairLaunchInventoryInfos *[]models.FairLaunchInventoryInfo) (amount int) {
	for _, inventory := range *fairLaunchInventoryInfos {
		amount += inventory.Quantity
	}
	return amount
}

// IsDuringMintTime
// @Description: timestamp now is between start and end
func IsDuringMintTime(start int, end int) bool {
	now := int(time.Now().Unix())
	return now >= start && now < end
}

// IsFairLaunchInfoMintTimeValid
// @Description: call IsDuringMintTime
func IsFairLaunchInfoMintTimeValid(fairLaunchInfo *models.FairLaunchInfo) bool {
	return IsDuringMintTime(fairLaunchInfo.StartTime, fairLaunchInfo.EndTime)
}

// IsFairLaunchMintTimeRight
// @Description: call GetFairLaunchInfo and IsFairLaunchInfoMintTimeValid
func IsFairLaunchMintTimeRight(fairLaunchInfoId int) (bool, error) {
	fairLaunchInfo, err := GetFairLaunchInfo(fairLaunchInfoId)
	if err != nil {
		return false, utils.AppendErrorInfo(err, "GetFairLaunchInfo")
	}
	return IsFairLaunchInfoMintTimeValid(fairLaunchInfo), nil
}

// AmountAndQuantityToNumber
// @Description: calculate Number by Amount And Quantity
func AmountAndQuantityToNumber(amount int, quantity int) int {
	return int(math.Ceil(float64(amount) / float64(quantity)))
}

// CreateInventoryAndAssetIssuanceInfoByFairLaunchInfo
// @Description: Update inventory and asset issuance
func CreateInventoryAndAssetIssuanceInfoByFairLaunchInfo(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) (err error) {
	err = CreateInventoryInfoByFairLaunchInfo(tx, fairLaunchInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateInventoryInfoByFairLaunchInfo")
	}
	err = CreateAssetIssuanceInfoByFairLaunchInfo(tx, fairLaunchInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateAssetIssuanceInfoByFairLaunchInfo")
	}
	return nil
}

// FairLaunchInfos

func GetAllFairLaunchInfoByState(state models.FairLaunchState) (fairLaunchInfos *[]models.FairLaunchInfo, err error) {
	_fairLaunchInfos := make([]models.FairLaunchInfo, 0)
	fairLaunchInfos = &(_fairLaunchInfos)
	err = middleware.DB.Where("state = ?", state).Find(fairLaunchInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchInfos")
	}
	return fairLaunchInfos, nil
}

func GetAllFairLaunchStateNoPayInfos() (fairLaunchInfos *[]models.FairLaunchInfo, err error) {
	return GetAllFairLaunchInfoByState(models.FairLaunchStateNoPay)
}

func GetAllFairLaunchStatePaidPendingInfos() (fairLaunchInfos *[]models.FairLaunchInfo, err error) {
	return GetAllFairLaunchInfoByState(models.FairLaunchStatePaidPending)
}

func GetAllFairLaunchStatePaidNoIssueInfos() (fairLaunchInfos *[]models.FairLaunchInfo, err error) {
	return GetAllFairLaunchInfoByState(models.FairLaunchStatePaidNoIssue)
}

func GetAllFairLaunchStateIssuedPendingInfos() (fairLaunchInfos *[]models.FairLaunchInfo, err error) {
	return GetAllFairLaunchInfoByState(models.FairLaunchStateIssuedPending)
}

func GetAllFairLaunchStateIssuedInfos() (fairLaunchInfos *[]models.FairLaunchInfo, err error) {
	return GetAllFairLaunchInfoByState(models.FairLaunchStateIssued)
}

func GetAllFairLaunchStateReservedSentPending() (fairLaunchInfos *[]models.FairLaunchInfo, err error) {
	return GetAllFairLaunchInfoByState(models.FairLaunchStateReservedSentPending)
}

func GetAllFairLaunchStateReservedSent() (fairLaunchInfos *[]models.FairLaunchInfo, err error) {
	return GetAllFairLaunchInfoByState(models.FairLaunchStateReservedSent)
}

func GetAllValidFairLaunchInfos() (fairLaunchInfos *[]models.FairLaunchInfo, err error) {
	_fairLaunchInfos := make([]models.FairLaunchInfo, 0)
	fairLaunchInfos = &_fairLaunchInfos
	err = middleware.DB.Find(fairLaunchInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchInfos")
	}
	return fairLaunchInfos, nil
}

type ProcessionResult struct {
	Id int `json:"id"`
	models.JsonResult
}

func ProcessAllFairLaunchInfos(tx *gorm.DB) (*[]ProcessionResult, error) {
	var processionResults []ProcessionResult
	allFairLaunchInfos, err := GetAllValidFairLaunchInfos()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetAllValidFairLaunchInfos")
	}
	for _, fairLaunchInfo := range *allFairLaunchInfos {
		if fairLaunchInfo.State == models.FairLaunchStateNoPay {
			err = ProcessFairLaunchStateNoPayInfoService(tx, &fairLaunchInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		} else if fairLaunchInfo.State == models.FairLaunchStatePaidPending {
			err = ProcessFairLaunchStatePaidPendingInfoService(tx, &fairLaunchInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		} else if fairLaunchInfo.State == models.FairLaunchStatePaidNoIssue {
			err = ProcessFairLaunchStatePaidNoIssueInfoService(tx, &fairLaunchInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		} else if fairLaunchInfo.State == models.FairLaunchStateIssuedPending {
			err = ProcessFairLaunchStateIssuedPendingInfoService(tx, &fairLaunchInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		} else if fairLaunchInfo.State == models.FairLaunchStateReservedSentPending {
			err = ProcessFairLaunchStateReservedSentPending(tx, &fairLaunchInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		}
	}
	if processionResults == nil || len(processionResults) == 0 {
		err = errors.New("procession results null")
		return nil, err
	}
	return &processionResults, nil
}

// FairLaunchInfos Procession by state

// ProcessAllFairLaunchStateNoPayInfoService
// @dev: Short time intervals
func ProcessAllFairLaunchStateNoPayInfoService(tx *gorm.DB) (*[]ProcessionResult, error) {
	var processionResults []ProcessionResult
	fairLaunchInfos, err := GetAllFairLaunchStateNoPayInfos()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetAllFairLaunchStateNoPayInfos")
	}
	for _, fairLaunchInfo := range *fairLaunchInfos {
		{
			err = IncreaseFairLaunchInfoProcessNumber(tx, &fairLaunchInfo)
			if err != nil {
				// @dev: Do nothing
			}
			err = ProcessFairLaunchStateNoPayInfoService(tx, &fairLaunchInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		}
	}
	if processionResults == nil || len(processionResults) == 0 {
		err = errors.New("procession results null")
		return nil, err
	}
	return &processionResults, nil
}

// ProcessAllFairLaunchStatePaidPendingInfoService
// @dev: Short time intervals
func ProcessAllFairLaunchStatePaidPendingInfoService(tx *gorm.DB) (*[]ProcessionResult, error) {
	var processionResults []ProcessionResult
	fairLaunchInfos, err := GetAllFairLaunchStatePaidPendingInfos()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetAllFairLaunchStatePaidPendingInfos")
	}
	for _, fairLaunchInfo := range *fairLaunchInfos {
		{
			err = IncreaseFairLaunchInfoProcessNumber(tx, &fairLaunchInfo)
			if err != nil {
				// @dev: Do nothing
			}
			err = ProcessFairLaunchStatePaidPendingInfoService(tx, &fairLaunchInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		}
	}
	if processionResults == nil || len(processionResults) == 0 {
		err = errors.New("procession results null")
		return nil, err
	}
	return &processionResults, nil
}

// ProcessAllFairLaunchStatePaidNoIssueInfoService
// @notice: This spends utxos，need to update
func ProcessAllFairLaunchStatePaidNoIssueInfoService(tx *gorm.DB) (*[]ProcessionResult, error) {
	var processionResults []ProcessionResult
	fairLaunchInfos, err := GetAllFairLaunchStatePaidNoIssueInfos()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetAllFairLaunchStatePaidNoIssueInfos")
	}
	for _, fairLaunchInfo := range *fairLaunchInfos {
		{
			err = IncreaseFairLaunchInfoProcessNumber(tx, &fairLaunchInfo)
			if err != nil {
				// @dev: Do nothing
			}
			err = ProcessFairLaunchStatePaidNoIssueInfoService(tx, &fairLaunchInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		}
	}
	if processionResults == nil || len(processionResults) == 0 {
		err = errors.New("procession results null")
		return nil, err
	}
	return &processionResults, nil
}

// ProcessAllFairLaunchStateIssuedPendingInfoService
// @dev: Short time intervals
func ProcessAllFairLaunchStateIssuedPendingInfoService(tx *gorm.DB) (*[]ProcessionResult, error) {
	var processionResults []ProcessionResult
	fairLaunchInfos, err := GetAllFairLaunchStateIssuedPendingInfos()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetAllFairLaunchStateIssuedPendingInfos")
	}
	for _, fairLaunchInfo := range *fairLaunchInfos {
		{
			err = IncreaseFairLaunchInfoProcessNumber(tx, &fairLaunchInfo)
			if err != nil {
				// @dev: Do nothing
			}
			err = ProcessFairLaunchStateIssuedPendingInfoService(tx, &fairLaunchInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		}
	}
	if processionResults == nil || len(processionResults) == 0 {
		err = errors.New("procession results null")
		return nil, err
	}
	return &processionResults, nil
}

// ProcessAllFairLaunchStateReservedSentPending
// @dev: Short time intervals
func ProcessAllFairLaunchStateReservedSentPending(tx *gorm.DB) (*[]ProcessionResult, error) {
	var processionResults []ProcessionResult
	fairLaunchInfos, err := GetAllFairLaunchStateReservedSentPending()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetAllFairLaunchStateReservedSentPending")
	}
	for _, fairLaunchInfo := range *fairLaunchInfos {
		{
			err = IncreaseFairLaunchInfoProcessNumber(tx, &fairLaunchInfo)
			if err != nil {
				// @dev: Do nothing
			}
			err = ProcessFairLaunchStateReservedSentPending(tx, &fairLaunchInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		}
	}
	if processionResults == nil || len(processionResults) == 0 {
		err = errors.New("procession results null")
		return nil, err
	}
	return &processionResults, nil
}

func UpdateFairLaunchInfoPaidId(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo, paidId int) (err error) {
	fairLaunchInfo.IssuanceFeePaidID = paidId
	fairLaunchInfo.PayMethod = models.FeePaymentMethodCustodyAccount
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(tx, fairLaunchInfo)
}

func ChangeFairLaunchInfoState(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo, state models.FairLaunchState) (err error) {
	fairLaunchInfo.State = state
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(tx, fairLaunchInfo)
}

func ClearFairLaunchInfoProcessNumber(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) (err error) {
	fairLaunchInfo.ProcessNumber = 0
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(tx, fairLaunchInfo)
}

func IncreaseFairLaunchInfoProcessNumber(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) (err error) {
	fairLaunchInfo.ProcessNumber += 1
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(tx, fairLaunchInfo)
}

func ChangeFairLaunchInfoStateAndUpdatePaidSuccessTime(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) (err error) {
	fairLaunchInfo.State = models.FairLaunchStatePaidNoIssue
	fairLaunchInfo.PaidSuccessTime = utils.GetTimestamp()
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(tx, fairLaunchInfo)
}

func UpdateFairLaunchInfoBatchKeyAndBatchState(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo, batchKey string, batchState string) (err error) {
	fairLaunchInfo.BatchKey = batchKey
	fairLaunchInfo.BatchState = batchState
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(tx, fairLaunchInfo)
}

func UpdateFairLaunchInfoBatchTxidAndAssetId(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo, batchTxidAnchor string, batchState string, assetId string) (err error) {
	fairLaunchInfo.BatchTxidAnchor = batchTxidAnchor
	fairLaunchInfo.BatchState = batchState
	fairLaunchInfo.AssetID = assetId
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(tx, fairLaunchInfo)
}

func FairLaunchTapdMint(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) (err error) {
	// @dev: 1.taprpc MintAsset
	var isCollectible bool
	if fairLaunchInfo.AssetType == taprpc.AssetType_COLLECTIBLE {
		isCollectible = true
	}
	newMeta := api.NewMetaWithImageStr(fairLaunchInfo.Description, fairLaunchInfo.ImageData)
	mintResponse, err := api.MintAssetAndGetResponse(fairLaunchInfo.Name, isCollectible, newMeta, fairLaunchInfo.Amount, false)
	if err != nil {
		return utils.AppendErrorInfo(err, "MintAssetAndGetResponse")
	}
	// @dev: 2.update batchKey and batchState
	batchKey := hex.EncodeToString(mintResponse.GetPendingBatch().GetBatchKey())
	batchState := mintResponse.GetPendingBatch().GetState().String()
	err = UpdateFairLaunchInfoBatchKeyAndBatchState(tx, fairLaunchInfo, batchKey, batchState)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateFairLaunchInfoBatchKeyAndBatchState")
	}
	return nil
}

func FairLaunchTapdMintFinalize(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) (err error) {
	// @dev: FeeRate maybe need to choose
	finalizeResponse, err := api.FinalizeBatchAndGetResponse(fairLaunchInfo.FeeRate)
	if err != nil {
		return utils.AppendErrorInfo(err, "FinalizeBatchAndGetResponse")
	}
	if hex.EncodeToString(finalizeResponse.GetBatch().GetBatchKey()) != fairLaunchInfo.BatchKey {
		err = errors.New("finalize batch key is not equal mint batch key")
		return err
	}
	batchTxidAnchor := finalizeResponse.GetBatch().GetBatchTxid()
	batchState := finalizeResponse.GetBatch().GetState().String()
	assetId, err := api.BatchTxidAnchorToAssetId(batchTxidAnchor)
	if err != nil {
		return utils.AppendErrorInfo(err, "BatchTxidAnchorToAssetId")
	}
	// @dev: Record paid fee
	err = CreateFairLaunchIncomeOfServerPayIssuanceFinalizeFee(tx, int(fairLaunchInfo.ID), batchTxidAnchor)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateFairLaunchIncomeOfServerPayIssuanceFinalizeFee")
	}
	err = UpdateFairLaunchInfoBatchTxidAndAssetId(tx, fairLaunchInfo, batchTxidAnchor, batchState, assetId)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateFairLaunchInfoBatchTxidAndAssetId")
	}
	return nil
}

func GetTransactionConfirmedNumber(txid string) (mumConfirmations int, err error) {
	response, err := api.GetListChainTransactions()
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "GetListChainTransactions")
	}
	for _, transaction := range *response {
		if txid == transaction.TxHash {
			return transaction.NumConfirmations, nil
		}
	}
	err = errors.New("did not match transaction hash")
	return 0, err
}

func IsTransactionConfirmed(txid string) bool {
	if txid == "" {
		//err := errors.New("empty transaction hash")
		return false
	}
	mumConfirmations, err := GetTransactionConfirmedNumber(txid)
	if err != nil {
		return false
	}
	return mumConfirmations > 0
}

func UpdateFairLaunchInfoReservedCouldMintAndState(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) (err error) {
	fairLaunchInfo.ReservedCouldMint = true
	fairLaunchInfo.State = models.FairLaunchStateIssued
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(tx, fairLaunchInfo)
}

func GetFairLaunchInfoState(fairLaunchId int) (fairLaunchState models.FairLaunchState, err error) {
	var fairLaunchInfo *models.FairLaunchInfo
	fairLaunchInfo, err = GetFairLaunchInfo(fairLaunchId)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "GetFairLaunchInfo")
	}
	return fairLaunchInfo.State, nil
}

func IsFairLaunchIssued(fairLaunchId int) bool {
	state, err := GetFairLaunchInfoState(fairLaunchId)
	if err != nil {
		return false
	}
	return state >= models.FairLaunchStateIssued
}

func UpdateFairLaunchInfoStateAndIssuanceTime(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) (err error) {
	fairLaunchInfo.State = models.FairLaunchStateIssuedPending
	fairLaunchInfo.IssuanceTime = utils.GetTimestamp()
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(tx, fairLaunchInfo)
}

func IsWalletBalanceEnough(value int) bool {
	response, err := api.WalletBalanceAndGetResponse()
	if err != nil {
		return false
	}
	return response.ConfirmedBalance >= int64(value)
}

func IsAssetBalanceEnough(assetId string, amount int) bool {
	response, err := api.ListBalancesAndGetResponse(true)
	if err != nil {
		return false
	}
	for k, v := range (*response).AssetBalances {
		if k == assetId {
			return v.Balance >= uint64(amount)
		}
	}
	// @dev: No asset id
	return false
}

// IsAssetUtxoEnough
// @dev: If amount is negative, return false as error
func IsAssetUtxoEnough(assetId string, amount int) bool {
	// @dev: Amount can not be negative number
	if amount < 0 {
		return false
	}
	response, err := api.ListUtxosAndGetResponse()
	if err != nil {
		return false
	}
	for _, managedUtxo := range response.ManagedUtxos {
		if managedUtxo == nil {
			// @dev: Utxo not found
			continue
		}
		assets := managedUtxo.Assets
		if assets == nil || len(assets) == 0 {
			// @dev: Anchor assets is null
			continue
		}
		// @dev: Find asset in assets
		for _, asset := range assets {
			if hex.EncodeToString(asset.AssetGenesis.AssetId) == assetId {
				return asset.Amount >= uint64(amount)
			}
		}
		// @dev: Did not find asset in assets
	}
	// @dev: Did not find asset in ManagedUtxos
	return false
}

func ProcessIssuedFairLaunchInfos(fairLaunchInfos *[]models.FairLaunchInfo) *[]models.FairLaunchInfo {
	var result []models.FairLaunchInfo
	for _, fairLaunchInfo := range *fairLaunchInfos {
		if IsFairLaunchInfoMintTimeValid(&fairLaunchInfo) {
			result = append(result, fairLaunchInfo)
		}
	}
	return &result
}

// FairLaunchInfos Procession

func ProcessFairLaunchStateNoPayInfoService(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) (err error) {
	// @dev: 1.pay fee
	payIssuanceFeeResult, err := PayIssuanceFee(fairLaunchInfo.UserID, fairLaunchInfo.FeeRate)
	if err != nil {
		return utils.AppendErrorInfo(err, "PayIssuanceFee")
	}
	// @dev: Record paid fee
	err = CreateFairLaunchIncomeOfUserPayIssuanceFee(tx, int(fairLaunchInfo.ID), payIssuanceFeeResult.PaidId, payIssuanceFeeResult.Fee, fairLaunchInfo.UserID, fairLaunchInfo.Username)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateFairLaunchIncomeOfUserPayIssuanceFee")
	}
	// @dev: 2.Store paidId
	err = UpdateFairLaunchInfoPaidId(tx, fairLaunchInfo, payIssuanceFeeResult.PaidId)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateFairLaunchInfoPaidId")
	}
	// @dev: 3.Change state
	err = ChangeFairLaunchInfoState(tx, fairLaunchInfo, models.FairLaunchStatePaidPending)
	if err != nil {
		return utils.AppendErrorInfo(err, "ChangeFairLaunchInfoState")
	}
	err = ClearFairLaunchInfoProcessNumber(tx, fairLaunchInfo)
	if err != nil {
		// @dev: Do nothing
	}
	return nil
}

func ProcessFairLaunchStatePaidPendingInfoService(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) (err error) {
	// @dev: 1.fee paid
	var issuanceFeePaid bool
	issuanceFeePaid, err = IsIssuanceFeePaid(fairLaunchInfo.IssuanceFeePaidID)
	if err != nil {
		if errors.Is(err, models.CustodyAccountPayInsideMissionFaild) {
			err = SetFairLaunchInfoFail(tx, fairLaunchInfo)
			if err != nil {
				return utils.AppendErrorInfo(err, "SetFairLaunchInfoFail")
			}
		}
	}
	if issuanceFeePaid {
		// @dev: Change state
		err = ChangeFairLaunchInfoStateAndUpdatePaidSuccessTime(tx, fairLaunchInfo)
		if err != nil {
			return utils.AppendErrorInfo(err, "ChangeFairLaunchInfoStateAndUpdatePaidSuccessTime")
		}
		return nil
	}
	// @dev: fee has not been paid
	err = ClearFairLaunchInfoProcessNumber(tx, fairLaunchInfo)
	if err != nil {
		// @dev: Do nothing
	}
	return nil
}

func ProcessFairLaunchStatePaidNoIssueInfoService(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) (err error) {
	// @dev: Check if confirmed balance enough
	if !IsWalletBalanceEnough(fairLaunchInfo.SetGasFee) {
		err = errors.New("lnd wallet balance is not enough")
		return err
	}
	// @dev: 1.tapd mint, add to batch, finalize
	err = FairLaunchTapdMint(tx, fairLaunchInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "FairLaunchTapdMint")
	}
	// @dev: Consider whether to use scheduled task to finalize
	// @dev: Do not mint assets in one batch, which will anchor all assets to one utxo
	err = FairLaunchTapdMintFinalize(tx, fairLaunchInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "FairLaunchTapdMintFinalize")
	}
	// @dev: 2.Update asset issuance table
	err = CreateAssetIssuanceInfoByFairLaunchInfo(tx, fairLaunchInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateAssetIssuanceInfoByFairLaunchInfo")
	}
	// @dev: 3.update MintedAndAvailableInfo
	err = CreateFairLaunchMintedAndAvailableInfoByFairLaunchInfo(tx, fairLaunchInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateFairLaunchMintedAndAvailableInfoByFairLaunchInfo")
	}
	// @dev: Update state and issuance time
	err = UpdateFairLaunchInfoStateAndIssuanceTime(tx, fairLaunchInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateFairLaunchInfoStateAndIssuanceTime")
	}
	err = ClearFairLaunchInfoProcessNumber(tx, fairLaunchInfo)
	if err != nil {
		// @dev: Do nothing
	}
	return nil
}

func ProcessFairLaunchStateIssuedPendingInfoService(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) (err error) {
	// @dev: 1.Is Transaction Confirmed
	if IsTransactionConfirmed(fairLaunchInfo.BatchTxidAnchor) {
		// @dev: Update FairLaunchInfo ReservedCouldMint And Change State
		err = UpdateFairLaunchInfoReservedCouldMintAndState(tx, fairLaunchInfo)
		if err != nil {
			return utils.AppendErrorInfo(err, "UpdateFairLaunchInfoReservedCouldMintAndState")
		}
		// @dev: Update Asset Issuance
		var a = btldb.AssetIssuanceStore{DB: middleware.DB}
		var assetIssuance *models.AssetIssuance
		assetIssuance, err = a.ReadAssetIssuanceByFairLaunchId(fairLaunchInfo.ID)
		if err != nil {
			return utils.AppendErrorInfo(err, "ReadAssetIssuanceByFairLaunchId")
		}
		assetIssuance.State = models.AssetIssuanceStateIssued
		err = a.UpdateAssetIssuance(tx, assetIssuance)
		if err != nil {
			return utils.AppendErrorInfo(err, "UpdateAssetIssuance")
		}
		// TODO: 	After the issued assets are confirmed on-chain,
		// 			the server directly inserts the proof into the local universe
		// @dev: Maybe do not need to process here
		return nil
	}

	//// @dev: Split asset
	//// TODO: Insert proof to universe before do this operation
	//{
	//	// @dev: Create tow asset addrs
	//	var splitAssetAddrOne string
	//	var splitAssetAddrTwo string
	//	assetId := fairLaunchInfo.AssetID
	//	oneThirdAmount := fairLaunchInfo.Amount / 3
	//	splitAssetAddrOne, err = api.NewAddrAndGetStringResponse(assetId, oneThirdAmount)
	//	if err != nil {
	//		return err
	//	}
	//	splitAssetAddrTwo, err = api.NewAddrAndGetStringResponse(assetId, oneThirdAmount)
	//	if err != nil {
	//		return err
	//	}
	//	// @dev: Get fee rate
	//	var feeRate *FeeRateResponseTransformed
	//	feeRate, err = UpdateAndGetFeeRateResponseTransformed()
	//	if err != nil {
	//		return err
	//	}
	//	feeRateSatPerKw := feeRate.SatPerKw.FastestFee
	//	_, err = api.SendAssetAddrSliceAndGetResponse([]string{splitAssetAddrOne, splitAssetAddrTwo}, feeRateSatPerKw)
	//	if err != nil {
	//		return err
	//	}
	//}

	//@dev: Transaction has not been Confirmed
	err = ClearFairLaunchInfoProcessNumber(tx, fairLaunchInfo)
	if err != nil {
		// @dev: Do nothing
	}
	return nil
}

func ProcessFairLaunchStateReservedSentPending(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) (err error) {
	// @dev: 1.Is Transaction Confirmed
	if IsTransactionConfirmed(fairLaunchInfo.ReservedSentAnchorOutpointTxid) {
		// @dev: Change FairLaunchInfo State
		err = ChangeFairLaunchInfoState(tx, fairLaunchInfo, models.FairLaunchStateReservedSent)
		if err != nil {
			return utils.AppendErrorInfo(err, "ChangeFairLaunchInfoState")
		}
		return nil
	}
	// @dev: Transaction has not been Confirmed
	err = ClearFairLaunchInfoProcessNumber(tx, fairLaunchInfo)
	if err != nil {
		// @dev: Do nothing
	}
	return nil
}

// FairLaunchMintedInfos

func GetAllFairLaunchMintedInfoByState(state models.FairLaunchMintedState) (fairLaunchMintedInfos *[]models.FairLaunchMintedInfo, err error) {
	_fairLaunchMintedInfos := make([]models.FairLaunchMintedInfo, 0)
	fairLaunchMintedInfos = &(_fairLaunchMintedInfos)
	err = middleware.DB.Where("state = ?", state).Find(fairLaunchMintedInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchMintedInfos")
	}
	return fairLaunchMintedInfos, nil
}

func GetAllFairLaunchMintedStateNoPayInfo() (fairLaunchMintedInfos *[]models.FairLaunchMintedInfo, err error) {
	return GetAllFairLaunchMintedInfoByState(models.FairLaunchMintedStateNoPay)
}

func GetAllFairLaunchMintedStatePaidPendingInfo() (fairLaunchMintedInfos *[]models.FairLaunchMintedInfo, err error) {
	return GetAllFairLaunchMintedInfoByState(models.FairLaunchMintedStatePaidPending)
}

func GetAllFairLaunchMintedStatePaidNoSendInfo() (fairLaunchMintedInfos *[]models.FairLaunchMintedInfo, err error) {
	return GetAllFairLaunchMintedInfoByState(models.FairLaunchMintedStatePaidNoSend)
}

func GetAllFairLaunchMintedStateSentPendingInfo() (fairLaunchMintedInfos *[]models.FairLaunchMintedInfo, err error) {
	return GetAllFairLaunchMintedInfoByState(models.FairLaunchMintedStateSentPending)
}

func GetAllFairLaunchMintedStateSentInfo() (fairLaunchMintedInfos *[]models.FairLaunchMintedInfo, err error) {
	return GetAllFairLaunchMintedInfoByState(models.FairLaunchMintedStateSent)
}

func GetAllValidFairLaunchMintedInfos() (fairLaunchMintedInfos *[]models.FairLaunchMintedInfo, err error) {
	_fairLaunchMintedInfos := make([]models.FairLaunchMintedInfo, 0)
	fairLaunchMintedInfos = &(_fairLaunchMintedInfos)
	err = middleware.DB.Order("minted_set_time").Order("paid_success_time").Find(fairLaunchMintedInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchMintedInfos")
	}
	return fairLaunchMintedInfos, nil
}

func ProcessAllFairLaunchMintedInfos(tx *gorm.DB) (*[]ProcessionResult, error) {
	var processionResults []ProcessionResult
	allFairLaunchMintedInfos, err := GetAllValidFairLaunchMintedInfos()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetAllValidFairLaunchMintedInfos")
	}
	for _, fairLaunchMintedInfo := range *allFairLaunchMintedInfos {
		if fairLaunchMintedInfo.State == models.FairLaunchMintedStateNoPay {
			err = IncreaseFairLaunchMintedInfoProcessNumber(tx, &fairLaunchMintedInfo)
			if err != nil {
				// @dev: Do nothing
			}
			err = ProcessFairLaunchMintedStateNoPayInfo(tx, &fairLaunchMintedInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchMintedInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchMintedInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		} else if fairLaunchMintedInfo.State == models.FairLaunchMintedStatePaidPending {
			err = IncreaseFairLaunchMintedInfoProcessNumber(tx, &fairLaunchMintedInfo)
			if err != nil {
				// @dev: Do nothing
			}
			err = ProcessFairLaunchMintedStatePaidPendingInfo(tx, &fairLaunchMintedInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchMintedInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchMintedInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		} else if fairLaunchMintedInfo.State == models.FairLaunchMintedStatePaidNoSend {
			err = IncreaseFairLaunchMintedInfoProcessNumber(tx, &fairLaunchMintedInfo)
			if err != nil {
				// @dev: Do nothing
			}
			err = ProcessFairLaunchMintedStatePaidNoSendInfo(tx, &fairLaunchMintedInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchMintedInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchMintedInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		} else if fairLaunchMintedInfo.State == models.FairLaunchMintedStateSentPending {
			err = IncreaseFairLaunchMintedInfoProcessNumber(tx, &fairLaunchMintedInfo)
			if err != nil {
				// @dev: Do nothing
			}
			err = ProcessFairLaunchMintedStateSentPendingInfo(tx, &fairLaunchMintedInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchMintedInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					Id: int(fairLaunchMintedInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		}
	}
	if processionResults == nil || len(processionResults) == 0 {
		err = errors.New("procession results null")
		return nil, err
	}
	return &processionResults, nil
}

func UpdateFairLaunchMintedInfoPaidId(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo, paidId int) (err error) {
	fairLaunchMintedInfo.MintFeePaidID = paidId
	fairLaunchMintedInfo.PayMethod = models.FeePaymentMethodCustodyAccount
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchMintedInfo(tx, fairLaunchMintedInfo)
}

func ChangeFairLaunchMintedInfoState(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo, state models.FairLaunchMintedState) (err error) {
	fairLaunchMintedInfo.State = state
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchMintedInfo(tx, fairLaunchMintedInfo)
}

func ChangeFairLaunchMintedInfoStateAndUpdatePaidSuccessTime(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
	fairLaunchMintedInfo.State = models.FairLaunchMintedStatePaidNoSend
	fairLaunchMintedInfo.PaidSuccessTime = utils.GetTimestamp()
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchMintedInfo(tx, fairLaunchMintedInfo)
}

func LockInventoryByFairLaunchMintedInfo(fairLaunchMintedInfo *models.FairLaunchMintedInfo) (lockedInventory *[]models.FairLaunchInventoryInfo, err error) {
	lockedInventory, err = LockInventoryByFairLaunchMintedIdAndMintNumber(int(fairLaunchMintedInfo.ID), fairLaunchMintedInfo.MintedNumber)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "LockInventoryByFairLaunchMintedIdAndMintNumber")
	}
	return lockedInventory, nil
}

func GetAllUnsentFairLaunchMintedInfos() (fairLaunchMintedInfos *[]models.FairLaunchMintedInfo, err error) {
	_fairLaunchMintedInfos := make([]models.FairLaunchMintedInfo, 0)
	fairLaunchMintedInfos = &(_fairLaunchMintedInfos)
	err = middleware.DB.Where("state = ? AND is_addr_sent = ?", models.FairLaunchMintedStateSentPending, false).Find(fairLaunchMintedInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchMintedInfos")
	}
	return fairLaunchMintedInfos, nil
}

// UpdateFairLaunchMintedInfosIsAddrSent
// Deprecated
func UpdateFairLaunchMintedInfosIsAddrSent(fairLaunchMintedInfos *[]models.FairLaunchMintedInfo, isAddrSent bool) (err error) {
	return middleware.DB.Model(&fairLaunchMintedInfos).Update("is_addr_sent", isAddrSent).Error
}

func SendAssetResponseScriptKeyAndInternalKeyToOutpoint(sendAssetResponse *taprpc.SendAssetResponse, scriptKey string, internalKey string) (outpoint string, err error) {
	for _, output := range sendAssetResponse.Transfer.Outputs {
		outputScriptKey := hex.EncodeToString(output.ScriptKey)
		outputAnchorInternalKey := hex.EncodeToString(output.Anchor.InternalKey)
		if outputScriptKey == scriptKey && outputAnchorInternalKey == internalKey {
			return output.Anchor.Outpoint, nil
		}
	}
	err = errors.New("can not find anchor outpoint value")
	return "", err
}

// UpdateFairLaunchMintedInfosBySendAssetResponse
// @dev: Updated outpoint and is_addr_sent
func UpdateFairLaunchMintedInfosBySendAssetResponse(tx *gorm.DB, fairLaunchMintedInfos *[]models.FairLaunchMintedInfo, sendAssetResponse *taprpc.SendAssetResponse) (err error) {
	var fairLaunchMintedInfosUpdated []models.FairLaunchMintedInfo
	// @dev: Deprecate anchor tx hash
	_ = hex.EncodeToString(sendAssetResponse.Transfer.AnchorTxHash)
	for _, fairLaunchMintedInfo := range *fairLaunchMintedInfos {
		scriptKey := fairLaunchMintedInfo.ScriptKey
		internalKey := fairLaunchMintedInfo.InternalKey
		var outpoint string
		outpoint, err = SendAssetResponseScriptKeyAndInternalKeyToOutpoint(sendAssetResponse, scriptKey, internalKey)
		if err != nil {
			return utils.AppendErrorInfo(err, "SendAssetResponseScriptKeyAndInternalKeyToOutpoint")
		}
		fairLaunchMintedInfo.OutpointTxHash, _ = utils.GetTransactionAndIndexByOutpoint(outpoint)
		// @dev: Update outpoint and isAddrSent
		fairLaunchMintedInfo.Outpoint = outpoint
		fairLaunchMintedInfo.IsAddrSent = true
		var address string
		address, err = api.GetListChainTransactionsOutpointAddress(outpoint)
		if err != nil {
			return utils.AppendErrorInfo(err, "GetListChainTransactionsOutpointAddress")
		}
		fairLaunchMintedInfo.Address = address
		fairLaunchMintedInfo.SendAssetTime = utils.GetTimestamp()
		fairLaunchMintedInfosUpdated = append(fairLaunchMintedInfosUpdated, fairLaunchMintedInfo)
	}
	return tx.Save(&fairLaunchMintedInfosUpdated).Error
}

func FairLaunchMintedInfosIdToString(fairLaunchMintedInfos *[]models.FairLaunchMintedInfo) string {
	var ids string
	if fairLaunchMintedInfos == nil {
		fairLaunchMintedInfos = &[]models.FairLaunchMintedInfo{}
	}
	for _, fairLaunchMintedInfo := range *fairLaunchMintedInfos {
		ids += strconv.Itoa(int(fairLaunchMintedInfo.ID)) + ";"
	}
	return ids
}

// SendFairLaunchMintedAssetLocked
// @dev: Trigger after ProcessFairLaunchMintedStatePaidNoSendInfo
func SendFairLaunchMintedAssetLocked(tx *gorm.DB) error {
	// @dev: all unsent
	unsentFairLaunchMintedInfos, err := GetAllUnsentFairLaunchMintedInfos()
	ids := "(id:" + FairLaunchMintedInfosIdToString(unsentFairLaunchMintedInfos) + ")"
	if err != nil {
		return utils.AppendErrorInfo(err, "GetAllUnsentFairLaunchMintedInfos"+ids)
	}
	assetIdToFairLaunchInfoId := make(map[string]int)
	assetIdToAddrs := make(map[string][]string)
	assetIdToAmount := make(map[string]int)
	assetIdToGasFeeTotal := make(map[string]int)
	assetIdToGasFeeRateTotal := make(map[string]int)
	// @dev: addr Slice
	for _, fairLaunchMintedInfo := range *unsentFairLaunchMintedInfos {
		assetId := fairLaunchMintedInfo.AssetID
		if assetIdToAddrs[assetId] == nil || len(assetIdToAddrs[assetId]) == 0 {
			assetIdToAddrs[assetId] = make([]string, 0)
		}
		assetIdToAddrs[assetId] = append(assetIdToAddrs[assetId], fairLaunchMintedInfo.EncodedAddr)
		assetIdToAmount[assetId] += fairLaunchMintedInfo.AddrAmount
		assetIdToGasFeeTotal[assetId] += fairLaunchMintedInfo.MintedGasFee
		assetIdToGasFeeRateTotal[assetId] += fairLaunchMintedInfo.MintedFeeRateSatPerKw
		assetIdToFairLaunchInfoId[assetId] = fairLaunchMintedInfo.FairLaunchInfoID
	}
	assetIdToGasFeeRateAverage := make(map[string]int)
	for assetId, feeRateTotal := range assetIdToGasFeeTotal {
		assetIdToGasFeeRateAverage[assetId] = int(math.Ceil(float64(feeRateTotal) / float64(len(assetIdToAddrs[assetId]))))
	}
	var feeRate *FeeRateResponseTransformed
	var response *taprpc.SendAssetResponse
	for assetId, addrs := range assetIdToAddrs {
		// @dev: Check if confirmed btc balance enough
		if !IsWalletBalanceEnough(assetIdToGasFeeTotal[assetId]) {
			err = errors.New("lnd wallet balance is not enough" + ids)
			return err
		}
		// @dev: Check if asset balance enough
		if !IsAssetBalanceEnough(assetId, assetIdToAmount[assetId]) {
			err = errors.New("tapd asset balance is not enough" + ids)
			return err
		}
		// @dev: Check if asset utxo is enough
		if !IsAssetUtxoEnough(assetId, assetIdToAmount[assetId]) {
			err = errors.New("tapd asset utxo is not enough" + ids)
			return err
		}
		feeRate, err = UpdateAndGetFeeRateResponseTransformed()
		if err != nil {
			return utils.AppendErrorInfo(err, "UpdateAndGetFeeRateResponseTransformed"+ids)
		}
		// @dev: Append fee of 2 sat per b
		feeRateSatPerKw := feeRate.SatPerKw.FastestFee + FeeRateSatPerBToSatPerKw(2)
		// TODO;
		// @dev: Make sure fee rate is less than average plus 100 sat/b
		// @notice: 2024-8-14 16:37:34 This check should be removed
		// @note: Use greater restrictions instead of removing check
		if assetIdToGasFeeRateAverage[assetId]+FeeRateSatPerBToSatPerKw(100) < feeRateSatPerKw {
			return errors.New("too high fee rate to send minted asset now" + ids)
		}
		if len(addrs) == 0 {
			//err = errors.New("length of addrs slice is zero, can't send assets and update")
			//return err
			continue
		}
		// @dev: Send Asset
		response, err = api.SendAssetAddrSliceAndGetResponse(addrs, feeRateSatPerKw)
		if err != nil {
			return utils.AppendErrorInfo(err, "SendAssetAddrSliceAndGetResponse"+ids)
		}
		// @dev: Record paid fee
		outpoint := response.Transfer.Outputs[0].Anchor.Outpoint
		txid, _ := utils.OutpointToTransactionAndIndex(outpoint)
		addrsStr := AddrsToString(addrs)
		err = CreateFairLaunchIncomeOfServerPaySendAssetFee(tx, assetId, assetIdToFairLaunchInfoId[assetId], txid, addrsStr)
		if err != nil {
			return utils.AppendErrorInfo(err, "CreateFairLaunchIncomeOfServerPaySendAssetFee"+ids)
		}
		// @dev: Update minted info
		err = UpdateFairLaunchMintedInfosBySendAssetResponse(tx, unsentFairLaunchMintedInfos, response)
		if err != nil {
			return utils.AppendErrorInfo(err, "UpdateFairLaunchMintedInfosBySendAssetResponse"+ids)
		}
	}
	return nil
}

func AddrsToString(addrs []string) string {
	var addrsStr string
	for i, addr := range addrs {
		if i == 0 {
			addrsStr = addr
		} else {
			addrsStr += ", " + addr
		}
	}
	return addrsStr
}

func GetAllLockedInventoryByFairLaunchMintedInfo(fairLaunchMintedInfo *models.FairLaunchMintedInfo) (*[]models.FairLaunchInventoryInfo, error) {
	var fairLaunchInventoryInfos []models.FairLaunchInventoryInfo
	err := middleware.DB.Where("state = ? AND fair_launch_minted_info_id = ?", models.FairLaunchInventoryStateLocked, fairLaunchMintedInfo.ID).Find(&fairLaunchInventoryInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchInventoryInfos")
	}
	return &fairLaunchInventoryInfos, nil
}

func UpdateLockedInventoryByFairLaunchMintedInfo(fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
	fairLaunchMintedInfos, err := GetAllLockedInventoryByFairLaunchMintedInfo(fairLaunchMintedInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "GetAllLockedInventoryByFairLaunchMintedInfo")
	}
	// @dev: Update
	err = middleware.DB.Model(&fairLaunchMintedInfos).Updates(map[string]any{"is_minted": true, "state": models.FairLaunchInventoryStateMinted}).Error
	if err != nil {
		return utils.AppendErrorInfo(err, "Updates fairLaunchMintedInfos")
	}
	return nil
}

func UpdateMintedNumberAndIsMintAllOfFairLaunchInfoByFairLaunchMintedInfo(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
	fairLaunchInfoId := fairLaunchMintedInfo.FairLaunchInfoID
	fairLaunchInfo, err := GetFairLaunchInfo(fairLaunchInfoId)
	if err != nil {
		return utils.AppendErrorInfo(err, "GetFairLaunchInfo")
	}
	var isMintAll bool
	if fairLaunchInfo.MintedNumber+fairLaunchMintedInfo.MintedNumber >= fairLaunchInfo.MintNumber {
		isMintAll = true
	}
	fairLaunchInfo.MintedNumber += fairLaunchMintedInfo.MintedNumber
	fairLaunchInfo.IsMintAll = isMintAll
	return tx.Save(fairLaunchInfo).Error
}

func ClearFairLaunchMintedInfoProcessNumber(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
	fairLaunchMintedInfo.ProcessNumber = 0
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchMintedInfo(tx, fairLaunchMintedInfo)
}

func IncreaseFairLaunchMintedInfoProcessNumber(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
	fairLaunchMintedInfo.ProcessNumber += 1
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchMintedInfo(tx, fairLaunchMintedInfo)
}

// ProcessSentButNotUpdatedMintedInfo
// @Description: If procession was interrupted, attempt to continue processing
func ProcessSentButNotUpdatedMintedInfo(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
	var transfer *taprpc.AssetTransfer
	scriptKey := fairLaunchMintedInfo.ScriptKey
	if scriptKey == "" {
		err = errors.New("scriptKey is empty")
		return err
	}
	transfer, err = GetAssetTransferByScriptKey(fairLaunchMintedInfo.ScriptKey)
	if err != nil {
		return utils.AppendErrorInfo(err, "Get Asset Transfer By ScriptKey")
	}
	var outpoint string
	outpoint, err = GetOutpointByTransferAndScriptKey(transfer, scriptKey)
	if err != nil {
		return utils.AppendErrorInfo(err, "Get Outpoint By Transfer And ScriptKey")
	}
	var txHash string
	var address string
	// @dev: get tx and address
	txHash, _ = utils.GetTransactionAndIndexByOutpoint(outpoint)
	address, err = api.GetListChainTransactionsOutpointAddress(outpoint)
	if err != nil {
		return utils.AppendErrorInfo(err, "Get ListChainTransactions Outpoint Address")
	}
	// @dev: Update db and return
	fairLaunchMintedInfo.OutpointTxHash = txHash
	fairLaunchMintedInfo.Outpoint = outpoint
	fairLaunchMintedInfo.IsAddrSent = true
	fairLaunchMintedInfo.Address = address
	fairLaunchMintedInfo.SendAssetTime = int(transfer.TransferTimestamp)
	return tx.Save(&fairLaunchMintedInfo).Error
}

// FairLaunchMintedInfos Procession

func ProcessFairLaunchMintedStateNoPayInfo(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
	// @dev: 1.pay fee
	payMintedFeeResult, err := PayMintFee(fairLaunchMintedInfo.UserID, fairLaunchMintedInfo.MintedFeeRateSatPerKw)
	if err != nil {
		return nil
	}
	// @dev: Record paid fee
	err = CreateFairLaunchIncomeOfUserPayMintedFee(tx, fairLaunchMintedInfo.AssetID, fairLaunchMintedInfo.FairLaunchInfoID, int(fairLaunchMintedInfo.ID), payMintedFeeResult.PaidId, payMintedFeeResult.Fee, fairLaunchMintedInfo.UserID, fairLaunchMintedInfo.Username)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateFairLaunchIncomeOfUserPayMintedFee")
	}
	// @dev: 2.Store paidId
	err = UpdateFairLaunchMintedInfoPaidId(tx, fairLaunchMintedInfo, payMintedFeeResult.PaidId)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateFairLaunchMintedInfoPaidId")
	}
	// @dev: 3.Change state
	err = ChangeFairLaunchMintedInfoState(tx, fairLaunchMintedInfo, models.FairLaunchMintedStatePaidPending)
	if err != nil {
		return utils.AppendErrorInfo(err, "ChangeFairLaunchMintedInfoState")
	}
	err = ClearFairLaunchMintedInfoProcessNumber(tx, fairLaunchMintedInfo)
	if err != nil {
		// @dev: Do nothing
	}
	return nil
}

func ProcessFairLaunchMintedStatePaidPendingInfo(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
	// @dev: 1.fee paid
	var isMintFeePaid bool
	isMintFeePaid, err = IsMintFeePaid(fairLaunchMintedInfo.MintFeePaidID)
	if err != nil {
		if errors.Is(err, models.CustodyAccountPayInsideMissionFaild) {
			// test
			err = SetFairLaunchMintedInfoFail(tx, fairLaunchMintedInfo)
			if err != nil {
				return utils.AppendErrorInfo(err, "SetFairLaunchMintedInfoFail")
			}
		}
	}
	if isMintFeePaid {
		// @dev: Change state
		err = ChangeFairLaunchMintedInfoStateAndUpdatePaidSuccessTime(tx, fairLaunchMintedInfo)
		if err != nil {
			return utils.AppendErrorInfo(err, "ChangeFairLaunchMintedInfoStateAndUpdatePaidSuccessTime")
		}
		return nil
	}
	// @dev: fee has not been paid
	err = ClearFairLaunchMintedInfoProcessNumber(tx, fairLaunchMintedInfo)
	if err != nil {
		// @dev: Do nothing
	}
	return nil
}

func ProcessFairLaunchMintedStatePaidNoSendInfo(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
	// @dev: 1. Update MintedAndAvailableInfo
	err = UpdateFairLaunchMintedAndAvailableInfoByFairLaunchMintedInfo(tx, fairLaunchMintedInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateFairLaunchMintedAndAvailableInfoByFairLaunchMintedInfo")
	}
	// @dev: Change state
	err = ChangeFairLaunchMintedInfoState(tx, fairLaunchMintedInfo, models.FairLaunchMintedStateSentPending)
	if err != nil {
		return utils.AppendErrorInfo(err, "ChangeFairLaunchMintedInfoState")
	}
	err = ClearFairLaunchMintedInfoProcessNumber(tx, fairLaunchMintedInfo)
	if err != nil {
		// @dev: Do nothing
	}
	return nil
}

func ProcessFairLaunchMintedStateSentPendingInfo(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
	//@ dev: Process sent but not updated
	err = ProcessSentButNotUpdatedMintedInfo(tx, fairLaunchMintedInfo)
	if err != nil {
		// @dev: Do not return
	}
	if fairLaunchMintedInfo.OutpointTxHash == "" {
		err = errors.New("no outpoint of transaction hash generated, asset may has not been sent")
		return err
	}
	// @dev: 1.Is Transaction Confirmed
	if IsTransactionConfirmed(fairLaunchMintedInfo.OutpointTxHash) {
		// @dev: Change state
		err = ChangeFairLaunchMintedInfoState(tx, fairLaunchMintedInfo, models.FairLaunchMintedStateSent)
		if err != nil {
			return utils.AppendErrorInfo(err, "ChangeFairLaunchMintedInfoState")
		}
		// @dev: Update MintedNumber and IsMintAll
		err = UpdateMintedNumberAndIsMintAllOfFairLaunchInfoByFairLaunchMintedInfo(tx, fairLaunchMintedInfo)
		if err != nil {
			return utils.AppendErrorInfo(err, "UpdateMintedNumberAndIsMintAllOfFairLaunchInfoByFairLaunchMintedInfo")
		}
		// @dev: Do not update fairLaunchMintedInfo here
		// @dev: Update minted user
		f := btldb.FairLaunchStore{DB: middleware.DB}
		err = f.CreateFairLaunchMintedUserInfo(tx, &models.FairLaunchMintedUserInfo{
			UserID:                 fairLaunchMintedInfo.UserID,
			FairLaunchMintedInfoID: int(fairLaunchMintedInfo.ID),
			FairLaunchInfoID:       fairLaunchMintedInfo.FairLaunchInfoID,
			MintedNumber:           fairLaunchMintedInfo.MintedNumber,
		})
		if err != nil {
			return utils.AppendErrorInfo(err, "CreateFairLaunchMintedUserInfo")
		}
		err = btldb.CreateBalance(middleware.DB, &models.Balance{
			AccountId:   custodyAccount.AdminUserInfo.Account.ID,
			BillType:    models.BillTypeAssetMintedSend,
			Away:        models.AWAY_OUT,
			Amount:      float64(fairLaunchMintedInfo.AddrAmount),
			Unit:        models.UNIT_ASSET_NORMAL,
			AssetId:     &(fairLaunchMintedInfo.AssetID),
			Invoice:     &(fairLaunchMintedInfo.EncodedAddr),
			PaymentHash: &(fairLaunchMintedInfo.OutpointTxHash),
			State:       models.STATE_SUCCESS,
			TypeExt: &models.BalanceTypeExt{
				Type: models.BTExtFirLaunch,
			},
		})
		if err != nil {
			return utils.AppendErrorInfo(err, "CreateBalance")
		}
		return nil
	}
	// @dev: Transaction has not been Confirmed
	err = ClearFairLaunchMintedInfoProcessNumber(tx, fairLaunchMintedInfo)
	if err != nil {
		// @dev: Do nothing
	}
	return nil
}

func SendFairLaunchReserved(fairLaunchInfo *models.FairLaunchInfo, addr string) (response *taprpc.SendAssetResponse, err error) {
	if addr == "" {
		err = errors.New("addr is null string")
		return nil, err
	}
	decodedAddrInfo, err := api.GetDecodedAddrInfo(addr)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetDecodedAddrInfo")
	}
	if int(decodedAddrInfo.Amount) != fairLaunchInfo.ReserveTotal {
		err = errors.New("wrong addr amount value")
		return nil, err
	}
	// @dev: Send
	addrSlice := []string{addr}
	// @dev: Use same fee rate of issuance instead
	feeRateSatPerKw := fairLaunchInfo.FeeRate
	response, err = api.SendAssetAddrSliceAndGetResponse(addrSlice, feeRateSatPerKw)
	if err != nil {
		//FairLaunchDebugLogger.Info("Send Asset AddrSlice And Get Response %v", err)
		return nil, utils.AppendErrorInfo(err, "SendAssetAddrSliceAndGetResponse")
	}
	return response, nil
}

func GetIssuedFairLaunchInfos() (*[]models.FairLaunchInfo, error) {
	return btldb.ReadIssuedFairLaunchInfos()
}

func GetIssuedAndTimeValidFairLaunchInfos() (*[]models.FairLaunchInfo, error) {
	return btldb.ReadIssuedAndTimeValidFairLaunchInfos()
}

func GetOwnFairLaunchInfosByUserId(id int) (*[]models.FairLaunchInfo, error) {
	var fairLaunchInfos []models.FairLaunchInfo
	err := middleware.DB.Where("user_id = ?", id).Order("set_time").Find(&fairLaunchInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchInfos")
	}
	return &fairLaunchInfos, nil
}

func GetOwnFairLaunchInfosByUserIdIssued(id int) (*[]models.FairLaunchInfo, error) {
	var fairLaunchInfos []models.FairLaunchInfo
	err := middleware.DB.Where(" user_id = ? AND state = ?", id, models.FairLaunchStateIssued).Order("set_time").Find(&fairLaunchInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchInfos")
	}
	return &fairLaunchInfos, nil
}

type FairLaunchInfoSimplified struct {
	ID                    int                    `json:"id"`
	Name                  string                 `json:"name"`
	ReserveTotal          int                    `json:"reserve_total"`
	CalculationExpression string                 `json:"calculation_expression"`
	AssetID               string                 `json:"asset_id"`
	State                 models.FairLaunchState `json:"state"`
}

func FairLaunchInfoToFairLaunchInfoSimplified(fairLaunchInfo models.FairLaunchInfo) FairLaunchInfoSimplified {
	return FairLaunchInfoSimplified{
		ID:                    int(fairLaunchInfo.ID),
		Name:                  fairLaunchInfo.Name,
		ReserveTotal:          fairLaunchInfo.ReserveTotal,
		CalculationExpression: fairLaunchInfo.CalculationExpression,
		AssetID:               fairLaunchInfo.AssetID,
		State:                 fairLaunchInfo.State,
	}
}

func FairLaunchInfoSliceToFairLaunchInfoSimplifiedSlice(airLaunchInfos *[]models.FairLaunchInfo) *[]FairLaunchInfoSimplified {
	var fairLaunchInfoSimplifiedSlice []FairLaunchInfoSimplified
	if airLaunchInfos == nil {
		return &fairLaunchInfoSimplifiedSlice
	}
	for _, fairLaunchInfo := range *airLaunchInfos {
		fairLaunchInfoSimplifiedSlice = append(fairLaunchInfoSimplifiedSlice, FairLaunchInfoToFairLaunchInfoSimplified(fairLaunchInfo))
	}
	return &fairLaunchInfoSimplifiedSlice
}

func GetFairLaunchInfoSimplifiedByUserIdIssued(id int) (*[]FairLaunchInfoSimplified, error) {
	fairLaunchInfos, err := GetOwnFairLaunchInfosByUserIdIssued(id)
	if err != nil {
		return nil, err
	}
	return FairLaunchInfoSliceToFairLaunchInfoSimplifiedSlice(fairLaunchInfos), nil
}

func GetOwnFairLaunchMintedInfosByUserId(id int) (*[]models.FairLaunchMintedInfo, error) {
	var fairLaunchMintedInfos []models.FairLaunchMintedInfo
	err := middleware.DB.Where(" user_id = ?", id).Order("minted_set_time").Find(&fairLaunchMintedInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchMintedInfos")
	}
	return &fairLaunchMintedInfos, nil
}

func GetOwnFairLaunchMintedNumberByUserIdAndFairLaunchId(userId int, fairLaunchId int) (int, error) {
	var fairLaunchMintedInfos []models.FairLaunchMintedInfo
	err := middleware.DB.Where(" user_id = ? AND fair_launch_info_id = ?", userId, fairLaunchId).Order("minted_set_time").Find(&fairLaunchMintedInfos).Error
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "Find fairLaunchMintedInfos")
	}
	var result int
	for _, fairLaunchMintedInfo := range fairLaunchMintedInfos {
		result += fairLaunchMintedInfo.MintedNumber
	}
	return result, nil
}

func IsMintedNumberValid(userId int, fairLaunchInfoId int, mintedNumber int) (bool, error) {
	recordNumber, err := GetOwnFairLaunchMintedNumberByUserIdAndFairLaunchId(userId, fairLaunchInfoId)
	if err != nil {
		return false, utils.AppendErrorInfo(err, "Get Own Fair Launch Minted Number By UserId And FairLaunchId")
	}
	if recordNumber+mintedNumber > models.MintMaxNumber {
		err = errors.New("Reach max mint number, available: " + strconv.Itoa(models.MintMaxNumber-recordNumber))
		return false, err
	}
	return true, nil
}

func IsEachMintedNumberValid(mintedNumber int) (bool, error) {
	if mintedNumber > models.MintMaxNumber {
		err := errors.New("Reach max mint number, available: " + strconv.Itoa(models.MintMaxNumber))
		return false, err
	}
	return true, nil
}

func ProcessSendFairLaunchReservedResponse(response *taprpc.SendAssetResponse) string {
	op := response.Transfer.Outputs[0].Anchor.Outpoint
	return op
}

func UpdateFairLaunchInfoIsReservedSent(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo, outpoint string) (err error) {
	fairLaunchInfo.IsReservedSent = true
	txid, _ := utils.OutpointToTransactionAndIndex(outpoint)
	fairLaunchInfo.ReservedSentAnchorOutpointTxid = txid
	fairLaunchInfo.ReservedSentAnchorOutpoint = outpoint
	fairLaunchInfo.State = models.FairLaunchStateReservedSentPending
	f := btldb.FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(tx, fairLaunchInfo)
}

func GetFairLaunchInfoByAssetId(assetId string) (*models.FairLaunchInfo, error) {
	var fairLaunchInfo models.FairLaunchInfo
	err := middleware.DB.Where("asset_id = ?", assetId).First(&fairLaunchInfo).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "First fairLaunchInfo")
	}
	return &fairLaunchInfo, nil
}

func GetAssetTransferByScriptKey(scriptKey string) (*taprpc.AssetTransfer, error) {
	response, err := api.ListTransfersAndGetResponse()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "List Transfers And Get Response")
	}
	for _, transfer := range response.Transfers {
		for _, output := range transfer.Outputs {
			if scriptKey == hex.EncodeToString(output.ScriptKey) {
				return transfer, nil
			}
		}
	}
	err = errors.New("scriptKey not found")
	return nil, err
}

func GetOutpointByTransferAndScriptKey(transfer *taprpc.AssetTransfer, scriptKey string) (string, error) {
	for _, output := range transfer.Outputs {
		if scriptKey == hex.EncodeToString(output.ScriptKey) {
			return output.Anchor.Outpoint, nil
		}
	}
	err := errors.New("scriptKey not found")
	return "", err
}

func DeleteFairLaunchInventoryInfosByState(fairLaunchInventoryState models.FairLaunchInventoryState) error {
	return middleware.DB.Where("state = ?", fairLaunchInventoryState).Delete(&models.FairLaunchInventoryInfo{}).Error
}

func RemoveFairLaunchInventoryStateMintedInfos() error {
	return DeleteFairLaunchInventoryInfosByState(models.FairLaunchInventoryStateMinted)
}

func GetClosedFairLaunchInfo() (*[]models.FairLaunchInfo, error) {
	return btldb.ReadClosedFairLaunchInfo()
}

func GetNotStartedFairLaunchInfo() (*[]models.FairLaunchInfo, error) {
	return btldb.ReadNotStartedFairLaunchInfo()
}

func FairLaunchInfosToAssetIdMapMintedRate(fairLaunchInfos *[]models.FairLaunchInfo) *map[string]float64 {
	if fairLaunchInfos == nil {
		return nil
	}
	assetIdMapMintedRate := make(map[string]float64)
	for _, fairLaunchInfo := range *fairLaunchInfos {
		if fairLaunchInfo.AssetID == "" {
			continue
		}
		assetIdMapMintedRate[fairLaunchInfo.AssetID] = float64(fairLaunchInfo.MintedNumber) / float64(fairLaunchInfo.MintNumber)
	}

	return &assetIdMapMintedRate
}

type FairLaunchAssetIdAndMintedRate struct {
	AssetId    string  `json:"asset_id"`
	MintedRate float64 `json:"minted_rate"`
}

func FairLaunchAssetIdAndMintedRateSortByMintedRate(fairLaunchAssetIdAndMintedRates *[]FairLaunchAssetIdAndMintedRate) {
	if fairLaunchAssetIdAndMintedRates == nil {
		return
	}
	sort.Slice(*fairLaunchAssetIdAndMintedRates, func(i, j int) bool {
		return (*fairLaunchAssetIdAndMintedRates)[i].MintedRate > (*fairLaunchAssetIdAndMintedRates)[j].MintedRate
	})
}

func AssetIdMapMintedRateToFairLaunchAssetIdAndMintedRates(assetIdMapMintedRate *map[string]float64) *[]FairLaunchAssetIdAndMintedRate {
	if assetIdMapMintedRate == nil {
		return nil
	}
	fairLaunchAssetIdAndMintedRates := make([]FairLaunchAssetIdAndMintedRate, len(*assetIdMapMintedRate))
	for assetId, mintedRate := range *assetIdMapMintedRate {
		fairLaunchAssetIdAndMintedRates = append(fairLaunchAssetIdAndMintedRates, FairLaunchAssetIdAndMintedRate{
			AssetId:    assetId,
			MintedRate: mintedRate,
		})
	}
	return &fairLaunchAssetIdAndMintedRates
}

func FairLaunchInfosToAssetIdMapFairLaunchInfo(fairLaunchInfos *[]models.FairLaunchInfo) *map[string]*models.FairLaunchInfo {
	if fairLaunchInfos == nil {
		return nil
	}
	assetIdMapFairLaunchInfo := make(map[string]*models.FairLaunchInfo)
	for _, fairLaunchInfo := range *fairLaunchInfos {
		if fairLaunchInfo.AssetID == "" {
			continue
		}
		assetIdMapFairLaunchInfo[fairLaunchInfo.AssetID] = &fairLaunchInfo
	}
	return &assetIdMapFairLaunchInfo
}

// GetSortedFairLaunchInfosByMintedRate
// @Description: Sot by minted rate
func GetSortedFairLaunchInfosByMintedRate(fairLaunchInfos *[]models.FairLaunchInfo) *[]models.FairLaunchInfo {
	if fairLaunchInfos == nil {
		return nil
	}
	var sortedFairLaunchInfos []models.FairLaunchInfo
	assetIdMapMintedRate := FairLaunchInfosToAssetIdMapMintedRate(fairLaunchInfos)
	fairLaunchAssetIdAndMintedRates := AssetIdMapMintedRateToFairLaunchAssetIdAndMintedRates(assetIdMapMintedRate)
	FairLaunchAssetIdAndMintedRateSortByMintedRate(fairLaunchAssetIdAndMintedRates)
	assetIdMapFairLaunchInfo := FairLaunchInfosToAssetIdMapFairLaunchInfo(fairLaunchInfos)
	for _, fairLaunchAssetIdAndMintedRate := range *fairLaunchAssetIdAndMintedRates {
		fairLaunchInfo, ok := (*assetIdMapFairLaunchInfo)[fairLaunchAssetIdAndMintedRate.AssetId]
		if !ok {
			continue
		}
		sortedFairLaunchInfos = append(sortedFairLaunchInfos, *fairLaunchInfo)
	}
	return &sortedFairLaunchInfos
}

func FairLaunchFollowsToFairLaunchInfoIdSlice(fairLaunchFollows *[]models.FairLaunchFollow) *[]int {
	if fairLaunchFollows == nil {
		return nil
	}
	var fairLaunchInfoIds []int
	for _, fairLaunchFollow := range *fairLaunchFollows {
		fairLaunchInfoIds = append(fairLaunchInfoIds, fairLaunchFollow.FairLaunchInfoId)
	}
	return &fairLaunchInfoIds
}

func GetFairLaunchInfosByIds(fairLaunchInfoIds *[]int) (*[]models.FairLaunchInfo, error) {
	if fairLaunchInfoIds == nil || len(*fairLaunchInfoIds) == 0 {
		return &[]models.FairLaunchInfo{}, nil
	}
	return btldb.GetFairLaunchInfosByIds(fairLaunchInfoIds)
}

func GetFollowedFairLaunchInfo(userId int) (*[]models.FairLaunchInfo, error) {
	fairLaunchFollows, err := GetFairLaunchFollowsByUserId(userId)
	if err != nil {
		return nil, err
	}
	fairLaunchInfoIds := FairLaunchFollowsToFairLaunchInfoIdSlice(fairLaunchFollows)
	fairLaunchInfos, err := GetFairLaunchInfosByIds(fairLaunchInfoIds)
	if err != nil {
		return nil, err
	}
	return fairLaunchInfos, nil
}

type FairLaunchMintFeeInfo struct {
	FeeRateOfMintNumber *FeeRateResponseTransformed `json:"fee_rate_of_mint_number"`
	FeeOfMintNumber     int                         `json:"fee_of_mint_number"`
}

type FairLaunchPlusInfo struct {
	FairLaunchInfo             *models.FairLaunchInfo        `json:"fair_launch_info"`
	HolderNumber               int                           `json:"holder_number"`
	FairLaunchMintNumberMapFee map[int]FairLaunchMintFeeInfo `json:"fair_launch_mint_number_map_fee"`
}

func UpdateAndGetAllCalculateGasFee() (map[int]FairLaunchMintFeeInfo, error) {
	UpdateFeeRateByMempool()
	fairLaunchMintNumberMapFee := make(map[int]FairLaunchMintFeeInfo)
	for number := 1; number < 11; number++ {
		feeRate, err := CalculateGasFeeRateByMempool(number)
		if err != nil {
			return nil, err
		}
		calculatedFeeRateSatPerKw := feeRate.SatPerKw.FastestFee + FeeRateSatPerBToSatPerKw(2)
		fee := GetMintedTransactionGasFee(calculatedFeeRateSatPerKw)
		fairLaunchMintNumberMapFee[number] = FairLaunchMintFeeInfo{
			FeeRateOfMintNumber: feeRate,
			FeeOfMintNumber:     fee,
		}
	}
	return fairLaunchMintNumberMapFee, nil
}

func ProcessToFairLaunchPlusInfo(fairLaunchInfo *models.FairLaunchInfo, holderNumber int) *FairLaunchPlusInfo {
	fairLaunchMintNumberMapFee, err := UpdateAndGetAllCalculateGasFee()
	if err != nil {
		return nil
	}
	return &FairLaunchPlusInfo{
		FairLaunchInfo:             fairLaunchInfo,
		HolderNumber:               holderNumber,
		FairLaunchMintNumberMapFee: fairLaunchMintNumberMapFee,
	}
}

func SplitFairLaunchInventoryInfos(fairLaunchInventoryInfos *[]models.FairLaunchInventoryInfo) *[][]models.FairLaunchInventoryInfo {
	if fairLaunchInventoryInfos == nil {
		return nil
	}
	var fairLaunchInventoryInfoSlilces [][]models.FairLaunchInventoryInfo
	totalLength := len(*fairLaunchInventoryInfos)
	twoDimensionalSliceLength := int(math.Ceil(float64(totalLength) / 1000))
	for i := 0; i < twoDimensionalSliceLength; i++ {
		var inventory []models.FairLaunchInventoryInfo
		if i != twoDimensionalSliceLength-1 {
			inventory = (*fairLaunchInventoryInfos)[i*1000 : i*1000+1000]
		} else {
			inventory = (*fairLaunchInventoryInfos)[i*1000 : totalLength]
		}
		fairLaunchInventoryInfoSlilces = append(fairLaunchInventoryInfoSlilces, inventory)
	}
	return &fairLaunchInventoryInfoSlilces
}

func SetFairLaunchInventoryInfos(tx *gorm.DB, fairLaunchInventoryInfos *[]models.FairLaunchInventoryInfo) error {
	return btldb.CreateFairLaunchInventoryInfos(tx, fairLaunchInventoryInfos)
}

func CreateFairLaunchInventoryInfosBatchProcess(tx *gorm.DB, fairLaunchInventoryInfos *[]models.FairLaunchInventoryInfo) error {
	splitFairLaunchInventoryInfos := SplitFairLaunchInventoryInfos(fairLaunchInventoryInfos)
	var err error
	for _, inventories := range *splitFairLaunchInventoryInfos {
		err = SetFairLaunchInventoryInfos(tx, &inventories)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetNotIssuedFairLaunchInfos() (*[]models.FairLaunchInfo, error) {
	return btldb.ReadNotIssuedFairLaunchInfos()
}

func GetNotSentFairLaunchMintedInfos() (*[]models.FairLaunchMintedInfo, error) {
	return btldb.ReadNotSentFairLaunchMintedInfos()
}

type ProcessingFairLaunchIssuanceAndMint struct {
	FairLaunchInfos       *[]models.FairLaunchInfo       `json:"fair_launch_infos"`
	FairLaunchMintedInfos *[]models.FairLaunchMintedInfo `json:"fair_launch_minted_infos"`
}

func GetProcessingFairLaunchIssuanceAndMint() (*ProcessingFairLaunchIssuanceAndMint, error) {
	FairLaunchInfos, err := GetNotIssuedFairLaunchInfos()
	if err != nil {
		return nil, err
	}
	FairLaunchMintedInfos, err := GetNotSentFairLaunchMintedInfos()
	if err != nil {
		return nil, err
	}
	return &ProcessingFairLaunchIssuanceAndMint{
		FairLaunchInfos:       FairLaunchInfos,
		FairLaunchMintedInfos: FairLaunchMintedInfos,
	}, nil
}

func CheckIsFairLaunchIssuanceAndMintProcessing() error {
	processingFairLaunchIssuanceAndMint, err := GetProcessingFairLaunchIssuanceAndMint()
	if err != nil {
		return err
	}
	if processingFairLaunchIssuanceAndMint == nil {
		return nil
	}
	fairLaunchInfos := processingFairLaunchIssuanceAndMint.FairLaunchInfos
	fairLaunchMintedInfos := processingFairLaunchIssuanceAndMint.FairLaunchMintedInfos
	if !(fairLaunchInfos == nil && fairLaunchMintedInfos == nil) {
		err = errors.New("processing issuance and mint exists")
		var issuanceNum int
		var mintNum int
		if fairLaunchInfos != nil {
			issuanceNum = len(*fairLaunchInfos)
		}
		if fairLaunchMintedInfos != nil {
			mintNum = len(*fairLaunchMintedInfos)
		}
		if issuanceNum == 0 && mintNum == 0 {
			return nil
		}
		info := fmt.Sprintf("issuance: %d, mint: %d", issuanceNum, mintNum)
		return utils.AppendErrorInfo(err, info)
	}
	return nil
}

func CreateFairLaunchMintedAndAvailableInfo(tx *gorm.DB, fairLaunchMintedAndAvailableInfo *models.FairLaunchMintedAndAvailableInfo) error {
	return btldb.CreateFairLaunchMintedAndAvailableInfo(tx, fairLaunchMintedAndAvailableInfo)
}

func UpdateFairLaunchMintedAndAvailableInfo(tx *gorm.DB, fairLaunchMintedAndAvailableInfo *models.FairLaunchMintedAndAvailableInfo) error {
	return btldb.UpdateFairLaunchMintedAndAvailableInfo(tx, fairLaunchMintedAndAvailableInfo)
}

func CreateFairLaunchMintedAndAvailableInfoByFairLaunchInfo(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) error {
	return CreateFairLaunchMintedAndAvailableInfo(tx, &models.FairLaunchMintedAndAvailableInfo{
		FairLaunchInfoID:      int(fairLaunchInfo.ID),
		MintedNumber:          0,
		MintedAmount:          0,
		AvailableNumber:       fairLaunchInfo.MintNumber,
		AvailableAmount:       fairLaunchInfo.MintTotal,
		ReserveTotal:          fairLaunchInfo.ReserveTotal,
		MintTotal:             fairLaunchInfo.MintTotal,
		MintNumber:            fairLaunchInfo.MintNumber,
		MintQuantity:          fairLaunchInfo.MintQuantity,
		FinalQuantity:         fairLaunchInfo.FinalQuantity,
		CalculationExpression: fairLaunchInfo.CalculationExpression,
	})
}

func GetFairLaunchMintedAndAvailableInfoByFairLaunchInfoId(fairLaunchInfoId int) (*models.FairLaunchMintedAndAvailableInfo, error) {
	return btldb.ReadFairLaunchMintedAndAvailableInfoByFairLaunchInfoId(fairLaunchInfoId)
}

func UpdateFairLaunchMintedAndAvailableInfoByFairLaunchMintedInfo(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo) error {
	if fairLaunchMintedInfo == nil || fairLaunchMintedInfo.MintedNumber == 0 {
		return errors.New("invalid fair launch minted info or minted number is zero")
	}
	mintedAndAvailableInfo, err := GetFairLaunchMintedAndAvailableInfoByFairLaunchInfoId(fairLaunchMintedInfo.FairLaunchInfoID)
	if err != nil {
		return err
	}
	number := fairLaunchMintedInfo.MintedNumber
	// TODO: only for test
	fmt.Println(number)
	var amount int
	if number > mintedAndAvailableInfo.AvailableNumber {
		return errors.New("minted number " + strconv.Itoa(number) + " exceeds available number")
	} else if number == mintedAndAvailableInfo.AvailableNumber {
		amount = (number-1)*mintedAndAvailableInfo.MintQuantity + mintedAndAvailableInfo.FinalQuantity
	} else {
		amount = number * mintedAndAvailableInfo.MintQuantity
	}
	// TODO: only for test
	fmt.Println(number)
	if amount != fairLaunchMintedInfo.AddrAmount {
		return errors.New("minted amount " + strconv.Itoa(amount) + " is not equal minted info's addr amount " + strconv.Itoa(fairLaunchMintedInfo.AddrAmount))
	}
	mintedAndAvailableInfo.MintedNumber += number
	if mintedAndAvailableInfo.MintedNumber > mintedAndAvailableInfo.MintNumber {
		return errors.New("minted number " + strconv.Itoa(mintedAndAvailableInfo.MintedNumber) + " exceeds mint number")
	}
	mintedAndAvailableInfo.MintedAmount += amount
	if mintedAndAvailableInfo.MintedAmount > mintedAndAvailableInfo.MintTotal {
		return errors.New("minted amount " + strconv.Itoa(mintedAndAvailableInfo.MintedAmount) + " exceeds mint total")
	}
	mintedAndAvailableInfo.AvailableNumber -= number
	if mintedAndAvailableInfo.AvailableNumber < 0 {
		return errors.New("available number " + strconv.Itoa(mintedAndAvailableInfo.AvailableNumber) + " is less than zero")
	}
	mintedAndAvailableInfo.AvailableAmount -= amount
	if mintedAndAvailableInfo.AvailableAmount < 0 {
		return errors.New("available amount " + strconv.Itoa(mintedAndAvailableInfo.AvailableAmount) + " is less than zero")
	}
	// TODO: only for test
	fmt.Println(utils.ValueJsonString(mintedAndAvailableInfo))
	return btldb.UpdateFairLaunchMintedAndAvailableInfo(tx, mintedAndAvailableInfo)
}

func GetAmountCouldBeMintByMintedNumber(fairLaunchInfoID int, mintedNumber int) (int, error) {
	if mintedNumber == 0 {
		return 0, errors.New("invalid minted number(0)")
	}
	mintedAndAvailableInfo, err := GetFairLaunchMintedAndAvailableInfoByFairLaunchInfoId(fairLaunchInfoID)
	if err != nil {
		return 0, err
	}
	var amount int
	if mintedNumber > mintedAndAvailableInfo.AvailableNumber {
		return 0, errors.New("minted number " + strconv.Itoa(mintedNumber) + " exceeds available number")
	} else if mintedNumber == mintedAndAvailableInfo.AvailableNumber {
		amount = (mintedNumber-1)*mintedAndAvailableInfo.MintQuantity + mintedAndAvailableInfo.FinalQuantity
	} else {
		amount = mintedNumber * mintedAndAvailableInfo.MintQuantity
	}
	return amount, nil
}

type NumberAndAmountCouldBeMint struct {
	Number int `json:"number"`
	Amount int `json:"amount"`
}

func GetNumberAndAmountCouldBeMint(fairLaunchInfoID int) (*NumberAndAmountCouldBeMint, error) {
	mintedAndAvailableInfo, err := GetFairLaunchMintedAndAvailableInfoByFairLaunchInfoId(fairLaunchInfoID)
	if err != nil {
		return nil, err
	}
	return &NumberAndAmountCouldBeMint{
		Number: mintedAndAvailableInfo.AvailableNumber,
		Amount: mintedAndAvailableInfo.AvailableAmount,
	}, nil
}

func GetAllFairLaunchInventoryInfo() (*[]models.FairLaunchInventoryInfo, error) {
	var fairLaunchInventoryInfos []models.FairLaunchInventoryInfo
	err := middleware.DB.Find(&fairLaunchInventoryInfos).Error
	if err != nil {
		return nil, err
	}
	return &fairLaunchInventoryInfos, nil
}

// TODO: Fix
func FairLaunchInventoryToMintedAndAvailableInfo(tx *gorm.DB) (*[]models.FairLaunchMintedAndAvailableInfo, error) {
	inventory, err := GetAllFairLaunchInventoryInfo()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetAllFairLaunchInventoryInfo")
	}
	fairLaunchInfoIdMapExists := make(map[int]bool)
	for _, item := range *inventory {
		fairLaunchInfoIdMapExists[item.FairLaunchInfoID] = true
	}
	for id := range fairLaunchInfoIdMapExists {
		fmt.Print(id)
		_, err = GetFairLaunchMintedAndAvailableInfoByFairLaunchInfoId(id)
		if err != nil {
			var fairLaunchInfo *models.FairLaunchInfo
			fairLaunchInfo, err = GetFairLaunchInfo(id)
			if err != nil {
				return nil, utils.AppendErrorInfo(err, "GetFairLaunchInfo")
			}
			err = CreateFairLaunchMintedAndAvailableInfoByFairLaunchInfo(tx, fairLaunchInfo)
			if err != nil {
				return nil, utils.AppendErrorInfo(err, "CreateFairLaunchMintedAndAvailableInfoByFairLaunchInfo")
			}
			fmt.Println("Created")
		}
	}
	for _, item := range *inventory {
		if item.State == models.FairLaunchInventoryStateOpen && item.FairLaunchMintedInfoID == 0 {
			continue
		}
		fmt.Println(item.State, item.FairLaunchMintedInfoID)
		var fairLaunchMintedInfo *models.FairLaunchMintedInfo
		fairLaunchMintedInfo, err = GetFairLaunchMintedInfo(item.FairLaunchMintedInfoID)
		err = UpdateFairLaunchMintedAndAvailableInfoByFairLaunchMintedInfo(tx, fairLaunchMintedInfo)
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "UpdateFairLaunchMintedAndAvailableInfoByFairLaunchMintedInfo")
		}
	}
	var mintedAndAvailableInfos []models.FairLaunchMintedAndAvailableInfo
	for id := range fairLaunchInfoIdMapExists {
		var mintedAndAvailableInfo *models.FairLaunchMintedAndAvailableInfo
		mintedAndAvailableInfo, err = GetFairLaunchMintedAndAvailableInfoByFairLaunchInfoId(id)
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "GetFairLaunchMintedAndAvailableInfoByFairLaunchInfoId")
		}
		mintedAndAvailableInfos = append(mintedAndAvailableInfos, *mintedAndAvailableInfo)
	}
	return &mintedAndAvailableInfos, nil
}

func DeleteFairLaunchInfo(fairLaunchInfoId uint) error {
	return btldb.DeleteFairLaunchInfo(fairLaunchInfoId)
}

func UpdateFairLaunchInfo(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) error {
	return btldb.UpdateFairLaunchInfo(tx, fairLaunchInfo)
}

func SetFairLaunchInfoFail(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) error {
	fairLaunchInfo.State = models.FairLaunchStateFail
	return UpdateFairLaunchInfo(tx, fairLaunchInfo)
}

func RemoveFairLaunchInfo(fairLaunchInfo *models.FairLaunchInfo) error {
	return DeleteFairLaunchInfo(fairLaunchInfo.ID)
}

func DeleteFairLaunchMintedInfo(fairLaunchMintedInfoId uint) error {
	return btldb.DeleteFairLaunchMintedInfo(fairLaunchMintedInfoId)
}

func RemoveFairLaunchMintedInfo(fairLaunchMintedInfo *models.FairLaunchMintedInfo) error {
	return DeleteFairLaunchMintedInfo(fairLaunchMintedInfo.ID)
}

func UpdateFairLaunchMintedInfo(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo) error {
	return btldb.UpdateFairLaunchMintedInfo(tx, fairLaunchMintedInfo)
}

func SetFairLaunchMintedInfoFail(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo) error {
	fairLaunchMintedInfo.State = models.FairLaunchMintedStateFail
	return UpdateFairLaunchMintedInfo(tx, fairLaunchMintedInfo)
}

func CancelAndRefundFairLaunchMintedInfo(tx *gorm.DB, fairLaunchMintedInfoId int) (BackAmountMissionId int, err error) {
	var fairLaunchMintedInfo *models.FairLaunchMintedInfo
	fairLaunchMintedInfo, err = GetFairLaunchMintedInfo(fairLaunchMintedInfoId)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "GetFairLaunchMintedInfo")
	}
	err = SetFairLaunchMintedInfoFail(tx, fairLaunchMintedInfo)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "SetFairLaunchMintedInfoFail")
	}
	mintFeePaidId := fairLaunchMintedInfo.MintFeePaidID
	missionId, err := custodyAccount.BackAmount(uint(mintFeePaidId))
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "BackAmount")
	}
	return int(missionId), nil
}

// TODO: This function maybe need to update.
//
//	Consider if use scheduled task

// Depreciated
func RefundBlockFairLaunchMintedInfos(tx *gorm.DB) (missionIds []int, err error) {
	fairLaunchMintedInfos, err := GetFairLaunchMintedInfoWhoseProcessNumberIsMoreThanTenThousand()
	if err != nil {
		return
	}
	if fairLaunchMintedInfos == nil || len(*fairLaunchMintedInfos) == 0 {
		return
	}
	var missionId int
	for _, fairLaunchMintedInfo := range *fairLaunchMintedInfos {
		missionId, err = CancelAndRefundFairLaunchMintedInfo(tx, int(fairLaunchMintedInfo.ID))
		if err != nil {
			return
		}
		missionIds = append(missionIds, missionId)
	}
	return
}

func GetFairLaunchMintedInfosWhoseUsernameIsNull() (*[]models.FairLaunchMintedInfo, error) {
	return btldb.ReadFairLaunchMintedInfosWhoseUsernameIsNull()
}

func UpdateFairLaunchMintedInfosWhoseUsernameIsNull() error {
	fairLaunchMintedInfos, err := GetFairLaunchMintedInfosWhoseUsernameIsNull()
	if err != nil {
		return utils.AppendErrorInfo(err, "GetFairLaunchMintedInfosWhoseUsernameIsNull")
	}
	for i, fairLaunchMintedInfo := range *fairLaunchMintedInfos {
		var username string
		username, err = IdToName(fairLaunchMintedInfo.UserID)
		if err != nil {
			continue
		}
		(*fairLaunchMintedInfos)[i].Username = username
	}
	return btldb.UpdateFairLaunchMintedInfos(fairLaunchMintedInfos)
}

func GetUserFirstFairLaunchMintedInfoByUserId(userId int) (*models.FairLaunchMintedInfo, error) {
	return btldb.ReadUserFirstFairLaunchMintedInfoByUserId(userId)
}

func GetUserFirstFairLaunchMintedInfoByUserIdAndAssetId(userId int, assetId string) (*models.FairLaunchMintedInfo, error) {
	return btldb.ReadUserFirstFairLaunchMintedInfoByUserIdAndAssetId(userId, assetId)
}

func GetUserFirstFairLaunchMintedInfoByUsernameAndAssetId(username string, assetId string) (*models.FairLaunchMintedInfo, error) {
	return btldb.ReadUserFirstFairLaunchMintedInfoByUsernameAndAssetId(username, assetId)
}

func GetUserFirstFairLaunchMintedInfosByUserIdSlice(userIdSlice []int) (*map[int]models.FairLaunchMintedInfo, error) {
	userIdMapFairLaunchMintedInfo := make(map[int]models.FairLaunchMintedInfo)
	for _, userId := range userIdSlice {
		if _, ok := userIdMapFairLaunchMintedInfo[userId]; ok {
			continue
		}
		fairLaunchMintedInfo, err := GetUserFirstFairLaunchMintedInfoByUserId(userId)
		if err != nil {
			continue
		}
		userIdMapFairLaunchMintedInfo[userId] = *fairLaunchMintedInfo
	}
	return &userIdMapFairLaunchMintedInfo, nil
}

// GetUserFirstFairLaunchMintedInfosByUserIdSliceAndAssetId
// @Description: Get user first fair launch minted infos by user id slice and asset id
func GetUserFirstFairLaunchMintedInfosByUserIdSliceAndAssetId(userIdSlice []int, assetId string) (*map[int]models.FairLaunchMintedInfo, error) {
	userIdMapFairLaunchMintedInfo := make(map[int]models.FairLaunchMintedInfo)
	for _, userId := range userIdSlice {
		if _, ok := userIdMapFairLaunchMintedInfo[userId]; ok {
			continue
		}
		fairLaunchMintedInfo, err := GetUserFirstFairLaunchMintedInfoByUserIdAndAssetId(userId, assetId)
		if err != nil {
			continue
		}
		userIdMapFairLaunchMintedInfo[userId] = *fairLaunchMintedInfo
	}
	return &userIdMapFairLaunchMintedInfo, nil
}

func GetUserFirstFairLaunchMintedInfosByUsernameSliceAndAssetId(usernameSlice []string, assetId string) (*map[string]models.FairLaunchMintedInfo, error) {
	usernameMapFairLaunchMintedInfo := make(map[string]models.FairLaunchMintedInfo)
	for _, username := range usernameSlice {
		if _, ok := usernameMapFairLaunchMintedInfo[username]; ok {
			continue
		}
		fairLaunchMintedInfo, err := GetUserFirstFairLaunchMintedInfoByUsernameAndAssetId(username, assetId)
		if err != nil {
			//@dev: do not return
			continue
			//return nil, err
		}
		usernameMapFairLaunchMintedInfo[username] = *fairLaunchMintedInfo
	}
	return &usernameMapFairLaunchMintedInfo, nil
}

type FairLaunchMintedInfoSimplified struct {
	ID              uint                         `gorm:"primarykey"`
	MintedGasFee    int                          `json:"minted_gas_fee"`
	MintFeePaidID   int                          `json:"mint_fee_paid_id"`
	PaidSuccessTime int                          `json:"paid_success_time"`
	UserID          int                          `json:"user_id" gorm:"index"`
	Username        string                       `json:"username" gorm:"type:varchar(255)"`
	AssetID         string                       `json:"asset_id" gorm:"type:varchar(255)" gorm:"index"`
	AssetName       string                       `json:"asset_name" gorm:"type:varchar(255)"`
	State           models.FairLaunchMintedState `json:"state"`
}

func FairLaunchMintedInfoToFairLaunchMintedInfoSimplified(fairLaunchMintedInfo *models.FairLaunchMintedInfo) *FairLaunchMintedInfoSimplified {
	if fairLaunchMintedInfo == nil {
		return nil
	}
	return &FairLaunchMintedInfoSimplified{
		ID:              fairLaunchMintedInfo.ID,
		MintedGasFee:    fairLaunchMintedInfo.MintedGasFee,
		MintFeePaidID:   fairLaunchMintedInfo.MintFeePaidID,
		PaidSuccessTime: fairLaunchMintedInfo.PaidSuccessTime,
		UserID:          fairLaunchMintedInfo.UserID,
		Username:        fairLaunchMintedInfo.Username,
		AssetID:         fairLaunchMintedInfo.AssetID,
		AssetName:       fairLaunchMintedInfo.AssetName,
		State:           fairLaunchMintedInfo.State,
	}
}

// BackAmountForFairLaunchMintedInfos
// @Description: Refund mint fee to user
// @dev: Should write log to file
func BackAmountForFairLaunchMintedInfos(fairLaunchMintedInfos *[]models.FairLaunchMintedInfo) (*[]models.JsonResult, error) {
	var jsonResults []models.JsonResult
	for _, fairLaunchMintedInfo := range *fairLaunchMintedInfos {
		paidId := fairLaunchMintedInfo.MintFeePaidID
		backAmountId, err := custodyAccount.BackAmount(uint(paidId))
		if err != nil {
			jsonResults = append(jsonResults, models.JsonResult{
				Success: false,
				Error:   "fairLaunchMintedInfoId:" + strconv.Itoa(int(fairLaunchMintedInfo.ID)) + ";" + "paidId:" + strconv.Itoa(paidId) + ";" + (err.Error()),
				Code:    models.BackAmountErr,
				Data:    backAmountId,
			})
		} else {
			jsonResults = append(jsonResults, models.JsonResult{
				Success: true,
				Error:   "",
				Code:    models.SUCCESS,
				Data:    backAmountId,
			})
		}
	}
	return &jsonResults, nil
}

func UsernameMapFairLaunchMintedInfoToFairLaunchMintedInfos(usernameMapFairLaunchMintedInfo *map[string]models.FairLaunchMintedInfo) *[]models.FairLaunchMintedInfo {
	if usernameMapFairLaunchMintedInfo == nil {
		return nil
	}
	var fairLaunchMintedInfos []models.FairLaunchMintedInfo
	for _, fairLaunchMintedInfo := range *usernameMapFairLaunchMintedInfo {
		fairLaunchMintedInfos = append(fairLaunchMintedInfos, fairLaunchMintedInfo)
	}
	return &fairLaunchMintedInfos
}

type RefundResult struct {
	Results    *[]models.JsonResult `json:"results"`
	FailNumber int                  `json:"fail_number"`
}

// RefundUserFirstMintByUsernameAndAssetId
// @Description: Refund user first mint by username and asset id
func RefundUserFirstMintByUsernameAndAssetId(usernameSlice []string, assetId string) (*RefundResult, error) {
	usernameMapFairLaunchMintedInfo, err := GetUserFirstFairLaunchMintedInfosByUsernameSliceAndAssetId(usernameSlice, assetId)
	if err != nil {
		return nil, err
	}
	fairLaunchMintedInfos := UsernameMapFairLaunchMintedInfoToFairLaunchMintedInfos(usernameMapFairLaunchMintedInfo)
	var refundResult *[]models.JsonResult
	refundResult, err = BackAmountForFairLaunchMintedInfos(fairLaunchMintedInfos)
	if err != nil {
		return nil, err
	}
	var failNumber int
	for _, result := range *refundResult {
		if !result.Success {
			failNumber++
		}
	}
	return &RefundResult{
		Results:    refundResult,
		FailNumber: failNumber,
	}, nil
}

type RefundUserFirstMintRequest struct {
	Usernames []string `json:"usernames"`
	AssetId   string   `json:"asset_id"`
}
