package services

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
	"trade/api"
	"trade/middleware"
	"trade/models"
	"trade/utils"
)

func PrintProcessionResult(processionResult *[]ProcessionResult) {
	for _, result := range *processionResult {
		if !result.Success {
			FairLaunchDebugLogger.Info("%v", utils.ValueJsonString(result.Error))
		}
	}
}

// FairLaunchIssuance
// @Description: Scheduled Task
func FairLaunchIssuance() {
	processionResult, err := ProcessAllFairLaunchInfos()
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
func FairLaunchMint() {
	processionResult, err := ProcessAllFairLaunchMintedInfos()
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
func SendFairLaunchAsset() {
	err := SendFairLaunchMintedAssetLocked()
	if err != nil {
		return
	}
}

func GetAllFairLaunchInfos() (*[]models.FairLaunchInfo, error) {
	f := FairLaunchStore{DB: middleware.DB}
	var fairLaunchInfos []models.FairLaunchInfo
	err := f.DB.Find(&fairLaunchInfos).Error
	return &fairLaunchInfos, err
}

func GetFairLaunchInfo(id int) (*models.FairLaunchInfo, error) {
	f := FairLaunchStore{DB: middleware.DB}
	return f.ReadFairLaunchInfo(uint(id))
}

func GetFairLaunchMintedInfo(id int) (*models.FairLaunchMintedInfo, error) {
	f := FairLaunchStore{DB: middleware.DB}
	return f.ReadFairLaunchMintedInfo(uint(id))
}

func GetFairLaunchMintedInfosByFairLaunchId(fairLaunchId int) (*[]models.FairLaunchMintedInfo, error) {
	f := FairLaunchStore{DB: middleware.DB}
	var fairLaunchMintedInfos []models.FairLaunchMintedInfo
	//err := f.DB.Where("fair_launch_info_id = ?", int(uint(id))).Find(&fairLaunchMintedInfos).Error
	err := f.DB.Where(&models.FairLaunchMintedInfo{FairLaunchInfoID: int(uint(fairLaunchId))}).Find(&fairLaunchMintedInfos).Error
	return &fairLaunchMintedInfos, err
}

func SetFairLaunchInfo(fairLaunchInfo *models.FairLaunchInfo) error {
	f := FairLaunchStore{DB: middleware.DB}
	return f.CreateFairLaunchInfo(fairLaunchInfo)
}

func SetFairLaunchMintedInfo(fairLaunchMintedInfo *models.FairLaunchMintedInfo) error {
	f := FairLaunchStore{DB: middleware.DB}
	return f.CreateFairLaunchMintedInfo(fairLaunchMintedInfo)
}

// ProcessFairLaunchInfo
// @Description: Process fairLaunchInfo
func ProcessFairLaunchInfo(imageData string, name string, assetType int, amount int, reserved int, mintQuantity int, startTime int, endTime int, description string, feeRate int, userId int) (*models.FairLaunchInfo, error) {
	err := ValidateStartAndEndTime(startTime, endTime)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "ValidateStartAndEndTime")
	}
	calculateSeparateAmount, err := AmountReservedAndMintQuantityToReservedTotalAndMintTotal(amount, reserved, mintQuantity)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "AmountReservedAndMintQuantityToReservedTotalAndMintTotal")
	}
	var fairLaunchInfo models.FairLaunchInfo
	// @dev: Setting fee rate need to bigger equal than fee rate now
	err = ValidateFeeRate(feeRate)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "ValidateFeeRate")
	}
	setGasFee := GetTransactionFee(feeRate)
	if !IsAccountBalanceEnoughByUserId(uint(userId), uint64(setGasFee)) {
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
		return utils.AppendErrorInfo(err, "got: "+strconv.Itoa(feeRate)+", expected: >="+strconv.Itoa(feeRateSatPerKw))
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
func ProcessFairLaunchMintedInfo(fairLaunchInfoID int, mintedNumber int, mintedFeeRateSatPerKw int, addr string, userId int) (*models.FairLaunchMintedInfo, error) {
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
	feeRate, err := UpdateAndCalculateGasFeeRateByMempool(mintedNumber)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "UpdateAndCalculateGasFeeRateByMempool")
	}
	// TODO: maybe need to change comparison param
	calculatedFeeRateSatPerKw := feeRate.SatPerKw.FastestFee
	if !(mintedFeeRateSatPerKw >= calculatedFeeRateSatPerKw) {
		err = errors.New("setting minted calculated FeeRate SatPerKw need to bigger equal than calculated fee rate by mempool's recommended fastest fee rate now")
		return nil, utils.AppendErrorInfo(err, "got: "+strconv.Itoa(mintedFeeRateSatPerKw)+", expected: >="+strconv.Itoa(calculatedFeeRateSatPerKw))
	}
	mintedGasFee := GetMintedTransactionGasFee(mintedFeeRateSatPerKw)
	if !IsAccountBalanceEnoughByUserId(uint(userId), uint64(mintedGasFee)) {
		return nil, errors.New("account balance not enough to pay minted gas fee")
	}
	fairLaunchMintedInfo = models.FairLaunchMintedInfo{
		FairLaunchInfoID:      fairLaunchInfoID,
		MintedNumber:          mintedNumber,
		MintedFeeRateSatPerKw: mintedFeeRateSatPerKw,
		MintedGasFee:          mintedGasFee,
		EncodedAddr:           addr,
		UserID:                userId,
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
func CreateInventoryInfoByFairLaunchInfo(fairLaunchInfo *models.FairLaunchInfo) error {
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
	f := FairLaunchStore{DB: middleware.DB}
	return f.CreateFairLaunchInventoryInfos(&FairLaunchInventoryInfos)
}

// CreateAssetIssuanceInfoByFairLaunchInfo
// @Description: Create Asset Issuance Info By FairLaunchInfo
func CreateAssetIssuanceInfoByFairLaunchInfo(fairLaunchInfo *models.FairLaunchInfo) error {
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
	a := AssetIssuanceStore{DB: middleware.DB}
	return a.CreateAssetIssuance(&assetIssuance)
}

// GetAllInventoryInfoByFairLaunchInfoId
// @Description: Query all inventory by FairLaunchInfo id
func GetAllInventoryInfoByFairLaunchInfoId(fairLaunchInfoId int) (*[]models.FairLaunchInventoryInfo, error) {
	var fairLaunchInventoryInfos []models.FairLaunchInventoryInfo
	err := middleware.DB.Where("fair_launch_info_id = ? AND status = ?", fairLaunchInfoId, models.StatusNormal).Find(&fairLaunchInventoryInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchInventoryInfos")
	}
	return &fairLaunchInventoryInfos, nil
}

// GetInventoryCouldBeMintedByFairLaunchInfoId
// @Description: Get all Inventory Could Be Minted By FairLaunchInfoId
func GetInventoryCouldBeMintedByFairLaunchInfoId(fairLaunchInfoId int) (*[]models.FairLaunchInventoryInfo, error) {
	var fairLaunchInventoryInfos []models.FairLaunchInventoryInfo
	err := middleware.DB.Where("fair_launch_info_id = ? AND status = ? AND is_minted = ? AND state = ?", fairLaunchInfoId, models.StatusNormal, false, models.FairLaunchInventoryStateOpen).Find(&fairLaunchInventoryInfos).Error
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
func CreateInventoryAndAssetIssuanceInfoByFairLaunchInfo(fairLaunchInfo *models.FairLaunchInfo) (err error) {
	err = CreateInventoryInfoByFairLaunchInfo(fairLaunchInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateInventoryInfoByFairLaunchInfo")
	}
	err = CreateAssetIssuanceInfoByFairLaunchInfo(fairLaunchInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateAssetIssuanceInfoByFairLaunchInfo")
	}
	return nil
}

// FairLaunchInfos

func GetAllFairLaunchInfoByState(state models.FairLaunchState) (fairLaunchInfos *[]models.FairLaunchInfo, err error) {
	_fairLaunchInfos := make([]models.FairLaunchInfo, 0)
	fairLaunchInfos = &(_fairLaunchInfos)
	err = middleware.DB.Where("status = ? AND state = ?", models.StatusNormal, state).Find(fairLaunchInfos).Error
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
	err = middleware.DB.Where("status = ?", models.StatusNormal).Find(fairLaunchInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchInfos")
	}
	return fairLaunchInfos, nil
}

type ProcessionResult struct {
	id int
	models.JsonResult
}

func ProcessAllFairLaunchInfos() (*[]ProcessionResult, error) {
	var processionResults []ProcessionResult
	allFairLaunchInfos, err := GetAllValidFairLaunchInfos()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetAllValidFairLaunchInfos")
	}
	for _, fairLaunchInfo := range *allFairLaunchInfos {
		if fairLaunchInfo.State == models.FairLaunchStateNoPay {
			err = ProcessFairLaunchStateNoPayInfoService(&fairLaunchInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		} else if fairLaunchInfo.State == models.FairLaunchStatePaidPending {
			err = ProcessFairLaunchStatePaidPendingInfoService(&fairLaunchInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		} else if fairLaunchInfo.State == models.FairLaunchStatePaidNoIssue {
			err = ProcessFairLaunchStatePaidNoIssueInfoService(&fairLaunchInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		} else if fairLaunchInfo.State == models.FairLaunchStateIssuedPending {
			err = ProcessFairLaunchStateIssuedPendingInfoService(&fairLaunchInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		} else if fairLaunchInfo.State == models.FairLaunchStateReservedSentPending {
			err = ProcessFairLaunchStateReservedSentPending(&fairLaunchInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchInfo.ID),
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

func UpdateFairLaunchInfoPaidId(fairLaunchInfo *models.FairLaunchInfo, paidId int) (err error) {
	fairLaunchInfo.IssuanceFeePaidID = paidId
	fairLaunchInfo.PayMethod = models.FeePaymentMethodCustodyAccount
	f := FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(fairLaunchInfo)
}

func ChangeFairLaunchInfoState(fairLaunchInfo *models.FairLaunchInfo, state models.FairLaunchState) (err error) {
	fairLaunchInfo.State = state
	f := FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(fairLaunchInfo)
}

func ChangeFairLaunchInfoStateAndUpdatePaidSuccessTime(fairLaunchInfo *models.FairLaunchInfo) (err error) {
	fairLaunchInfo.State = models.FairLaunchStatePaidNoIssue
	fairLaunchInfo.PaidSuccessTime = utils.GetTimestamp()
	f := FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(fairLaunchInfo)
}

func UpdateFairLaunchInfoBatchKeyAndBatchState(fairLaunchInfo *models.FairLaunchInfo, batchKey string, batchState string) (err error) {
	fairLaunchInfo.BatchKey = batchKey
	fairLaunchInfo.BatchState = batchState
	f := FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(fairLaunchInfo)
}

func UpdateFairLaunchInfoBatchTxidAndAssetId(fairLaunchInfo *models.FairLaunchInfo, batchTxidAnchor string, batchState string, assetId string) (err error) {
	fairLaunchInfo.BatchTxidAnchor = batchTxidAnchor
	fairLaunchInfo.BatchState = batchState
	fairLaunchInfo.AssetID = assetId
	f := FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(fairLaunchInfo)
}

func FairLaunchTapdMint(fairLaunchInfo *models.FairLaunchInfo) (err error) {
	// @dev: 1.taprpc MintAsset
	var isCollectible bool
	if fairLaunchInfo.AssetType == taprpc.AssetType_COLLECTIBLE {
		isCollectible = true
	}
	newMeta := api.NewMeta(fairLaunchInfo.Description, fairLaunchInfo.ImageData)
	mintResponse, err := api.MintAssetAndGetResponse(fairLaunchInfo.Name, isCollectible, newMeta, fairLaunchInfo.Amount, false)
	if err != nil {
		return utils.AppendErrorInfo(err, "MintAssetAndGetResponse")
	}
	// @dev: 2.update batchKey and batchState
	batchKey := hex.EncodeToString(mintResponse.GetPendingBatch().GetBatchKey())
	batchState := mintResponse.GetPendingBatch().GetState().String()
	err = UpdateFairLaunchInfoBatchKeyAndBatchState(fairLaunchInfo, batchKey, batchState)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateFairLaunchInfoBatchKeyAndBatchState")
	}
	return nil
}

func FairLaunchTapdMintFinalize(fairLaunchInfo *models.FairLaunchInfo) (err error) {
	//TODO: FeeRate maybe need to choose
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
	err = UpdateFairLaunchInfoBatchTxidAndAssetId(fairLaunchInfo, batchTxidAnchor, batchState, assetId)
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

func UpdateFairLaunchInfoReservedCouldMintAndState(fairLaunchInfo *models.FairLaunchInfo) (err error) {
	fairLaunchInfo.ReservedCouldMint = true
	fairLaunchInfo.State = models.FairLaunchStateIssued
	f := FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(fairLaunchInfo)
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

func UpdateFairLaunchInfoStateAndIssuanceTime(fairLaunchInfo *models.FairLaunchInfo) (err error) {
	fairLaunchInfo.State = models.FairLaunchStateIssuedPending
	fairLaunchInfo.IssuanceTime = utils.GetTimestamp()
	f := FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(fairLaunchInfo)
}

// FairLaunchInfos Procession

func ProcessFairLaunchStateNoPayInfoService(fairLaunchInfo *models.FairLaunchInfo) (err error) {
	// @dev: 1.pay fee
	paidId, err := PayIssuanceFee(fairLaunchInfo.UserID, fairLaunchInfo.FeeRate)
	if err != nil {
		return utils.AppendErrorInfo(err, "PayIssuanceFee")
	}
	// @dev: 2.Store paidId
	err = UpdateFairLaunchInfoPaidId(fairLaunchInfo, paidId)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateFairLaunchInfoPaidId")
	}
	// @dev: 3.Change state
	err = ChangeFairLaunchInfoState(fairLaunchInfo, models.FairLaunchStatePaidPending)
	if err != nil {
		return utils.AppendErrorInfo(err, "ChangeFairLaunchInfoState")
	}
	return nil
}

func ProcessFairLaunchStatePaidPendingInfoService(fairLaunchInfo *models.FairLaunchInfo) (err error) {
	// @dev: 1.fee paid
	if IsIssuanceFeePaid(fairLaunchInfo.IssuanceFeePaidID) {
		// @dev: Change state
		err = ChangeFairLaunchInfoStateAndUpdatePaidSuccessTime(fairLaunchInfo)
		if err != nil {
			return utils.AppendErrorInfo(err, "ChangeFairLaunchInfoStateAndUpdatePaidSuccessTime")
		}
		return nil
	}
	// @dev: fee has not been paid
	return nil
}

func ProcessFairLaunchStatePaidNoIssueInfoService(fairLaunchInfo *models.FairLaunchInfo) (err error) {
	// @dev: 1.tapd mint, add to batch, finalize
	err = FairLaunchTapdMint(fairLaunchInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "FairLaunchTapdMint")
	}
	// @dev: Consider whether to use scheduled task to finalize
	err = FairLaunchTapdMintFinalize(fairLaunchInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "FairLaunchTapdMintFinalize")
	}
	// @dev: 2.Update asset issuance table
	err = CreateAssetIssuanceInfoByFairLaunchInfo(fairLaunchInfo)
	// @dev: 3.update inventory
	err = CreateInventoryInfoByFairLaunchInfo(fairLaunchInfo)
	// @dev: Update state and issuance time
	err = UpdateFairLaunchInfoStateAndIssuanceTime(fairLaunchInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateFairLaunchInfoStateAndIssuanceTime")
	}
	return nil
}

func ProcessFairLaunchStateIssuedPendingInfoService(fairLaunchInfo *models.FairLaunchInfo) (err error) {
	// @dev: 1.Is Transaction Confirmed
	if IsTransactionConfirmed(fairLaunchInfo.BatchTxidAnchor) {
		// @dev: Update FairLaunchInfo ReservedCouldMint And Change State
		err = UpdateFairLaunchInfoReservedCouldMintAndState(fairLaunchInfo)
		if err != nil {
			return utils.AppendErrorInfo(err, "UpdateFairLaunchInfoReservedCouldMintAndState")
		}
		// @dev: Update Asset Issuance
		var a = AssetIssuanceStore{DB: middleware.DB}
		var assetIssuance *models.AssetIssuance
		assetIssuance, err = a.ReadAssetIssuanceByFairLaunchId(fairLaunchInfo.ID)
		if err != nil {
			return utils.AppendErrorInfo(err, "ReadAssetIssuanceByFairLaunchId")
		}
		assetIssuance.State = models.AssetIssuanceStateIssued
		err = a.UpdateAssetIssuance(assetIssuance)
		if err != nil {
			return utils.AppendErrorInfo(err, "UpdateAssetIssuance")
		}
		return nil
	}
	// @dev: Transaction has not been Confirmed
	return nil
}

func ProcessFairLaunchStateReservedSentPending(fairLaunchInfo *models.FairLaunchInfo) (err error) {
	// @dev: 1.Is Transaction Confirmed
	if IsTransactionConfirmed(fairLaunchInfo.ReservedSentAnchorOutpointTxid) {
		// @dev: Change FairLaunchInfo State
		err = ChangeFairLaunchInfoState(fairLaunchInfo, models.FairLaunchStateReservedSent)
		if err != nil {
			return utils.AppendErrorInfo(err, "ChangeFairLaunchInfoState")
		}
		return nil
	}
	// @dev: Transaction has not been Confirmed
	return nil
}

// FairLaunchMintedInfos

func GetAllFairLaunchMintedInfoByState(state models.FairLaunchMintedState) (fairLaunchMintedInfos *[]models.FairLaunchMintedInfo, err error) {
	_fairLaunchMintedInfos := make([]models.FairLaunchMintedInfo, 0)
	fairLaunchMintedInfos = &(_fairLaunchMintedInfos)
	err = middleware.DB.Where("status = ? AND state = ?", models.StatusNormal, state).Find(fairLaunchMintedInfos).Error
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
	err = middleware.DB.Order("minted_set_time").Order("paid_success_time").Where("status = ?", models.StatusNormal).Find(fairLaunchMintedInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchMintedInfos")
	}
	return fairLaunchMintedInfos, nil
}

func ProcessAllFairLaunchMintedInfos() (*[]ProcessionResult, error) {
	var processionResults []ProcessionResult
	allFairLaunchMintedInfos, err := GetAllValidFairLaunchMintedInfos()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetAllValidFairLaunchMintedInfos")
	}
	for _, fairLaunchMintedInfo := range *allFairLaunchMintedInfos {
		if fairLaunchMintedInfo.State == models.FairLaunchMintedStateNoPay {
			err = ProcessFairLaunchMintedStateNoPayInfo(&fairLaunchMintedInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchMintedInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchMintedInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		} else if fairLaunchMintedInfo.State == models.FairLaunchMintedStatePaidPending {
			err = ProcessFairLaunchMintedStatePaidPendingInfo(&fairLaunchMintedInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchMintedInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchMintedInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		} else if fairLaunchMintedInfo.State == models.FairLaunchMintedStatePaidNoSend {
			err = ProcessFairLaunchMintedStatePaidNoSendInfo(&fairLaunchMintedInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchMintedInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchMintedInfo.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		} else if fairLaunchMintedInfo.State == models.FairLaunchMintedStateSentPending {
			err = ProcessFairLaunchMintedStateSentPendingInfo(&fairLaunchMintedInfo)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchMintedInfo.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					id: int(fairLaunchMintedInfo.ID),
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

func UpdateFairLaunchMintedInfoPaidId(fairLaunchMintedInfo *models.FairLaunchMintedInfo, paidId int) (err error) {
	fairLaunchMintedInfo.MintFeePaidID = paidId
	fairLaunchMintedInfo.PayMethod = models.FeePaymentMethodCustodyAccount
	f := FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchMintedInfo(fairLaunchMintedInfo)
}

func ChangeFairLaunchMintedInfoState(fairLaunchMintedInfo *models.FairLaunchMintedInfo, state models.FairLaunchMintedState) (err error) {
	fairLaunchMintedInfo.State = state
	f := FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchMintedInfo(fairLaunchMintedInfo)
}

func ChangeFairLaunchMintedInfoStateAndUpdatePaidSuccessTime(fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
	fairLaunchMintedInfo.State = models.FairLaunchMintedStatePaidNoSend
	fairLaunchMintedInfo.PaidSuccessTime = utils.GetTimestamp()
	f := FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchMintedInfo(fairLaunchMintedInfo)
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
	err = middleware.DB.Where("status = ? AND state = ? AND is_addr_sent = ?", models.StatusNormal, models.FairLaunchMintedStateSentPending, false).Find(fairLaunchMintedInfos).Error
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

// GetTransactionAndIndexByOutpoint
// @dev: Split outpoint
func GetTransactionAndIndexByOutpoint(outpoint string) (transaction string, index string) {
	result := strings.Split(outpoint, ":")
	return result[0], result[1]
}

func GetListChainTransactionsOutpointAddress(outpoint string) (address string, err error) {
	response, err := api.GetListChainTransactions()
	if err != nil {
		return "", utils.AppendErrorInfo(err, "GetListChainTransactions")
	}
	tx, indexStr := GetTransactionAndIndexByOutpoint(outpoint)
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "", utils.AppendErrorInfo(err, "GetTransactionAndIndexByOutpoint")
	}
	for _, transaction := range *response {
		if transaction.TxHash == tx {
			return transaction.DestAddresses[index], nil
		}
	}
	err = errors.New("did not match transaction outpoint")
	return "", err
}

// UpdateFairLaunchMintedInfosBySendAssetResponse
// @dev: Updated outpoint and is_addr_sent
func UpdateFairLaunchMintedInfosBySendAssetResponse(fairLaunchMintedInfos *[]models.FairLaunchMintedInfo, sendAssetResponse *taprpc.SendAssetResponse) (err error) {
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
		fairLaunchMintedInfo.OutpointTxHash, _ = GetTransactionAndIndexByOutpoint(outpoint)
		// @dev: Update outpoint and isAddrSent
		fairLaunchMintedInfo.Outpoint = outpoint
		fairLaunchMintedInfo.IsAddrSent = true
		var address string
		address, err = GetListChainTransactionsOutpointAddress(outpoint)
		if err != nil {
			return utils.AppendErrorInfo(err, "GetListChainTransactionsOutpointAddress")
		}
		fairLaunchMintedInfo.Address = address
		fairLaunchMintedInfo.SendAssetTime = utils.GetTimestamp()
		fairLaunchMintedInfosUpdated = append(fairLaunchMintedInfosUpdated, fairLaunchMintedInfo)
	}
	return middleware.DB.Save(&fairLaunchMintedInfosUpdated).Error
}

// SendFairLaunchMintedAssetLocked
// @dev: Trigger after ProcessFairLaunchMintedStatePaidNoSendInfo
func SendFairLaunchMintedAssetLocked() (err error) {
	// @dev: all unsent
	unsentFairLaunchMintedInfos, err := GetAllUnsentFairLaunchMintedInfos()
	if err != nil {
		return utils.AppendErrorInfo(err, "GetAllUnsentFairLaunchMintedInfos")
	}
	// @dev: addr Slice
	var addrSlice []string
	for _, fairLaunchMintedInfo := range *unsentFairLaunchMintedInfos {
		addrSlice = append(addrSlice, fairLaunchMintedInfo.EncodedAddr)
	}
	feeRate, err := UpdateAndGetFeeRateResponseTransformed()
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateAndGetFeeRateResponseTransformed")
	}
	feeRateSatPerKw := feeRate.SatPerKw.FastestFee
	if len(addrSlice) == 0 {
		err = errors.New("length of addr slice is zero, can't send assets and update")
		return err
	}
	// @dev: Send Asset
	response, err := api.SendAssetAddrSliceAndGetResponse(addrSlice, feeRateSatPerKw)
	if err != nil {
		return utils.AppendErrorInfo(err, "SendAssetAddrSliceAndGetResponse")
	}
	// @dev: Update minted info
	err = UpdateFairLaunchMintedInfosBySendAssetResponse(unsentFairLaunchMintedInfos, response)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateFairLaunchMintedInfosBySendAssetResponse")
	}
	return nil
}

func GetAllLockedInventoryByFairLaunchMintedInfo(fairLaunchMintedInfo *models.FairLaunchMintedInfo) (*[]models.FairLaunchInventoryInfo, error) {
	var fairLaunchInventoryInfos []models.FairLaunchInventoryInfo
	err := middleware.DB.Where("status = ? AND state = ? AND fair_launch_minted_info_id = ?", models.StatusNormal, models.FairLaunchInventoryStateLocked, fairLaunchMintedInfo.ID).Find(&fairLaunchInventoryInfos).Error
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

func UpdateMintedNumberAndIsMintAllOfFairLaunchInfoByFairLaunchMintedInfo(fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
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
	return middleware.DB.Save(fairLaunchInfo).Error
}

// FairLaunchMintedInfos Procession

func ProcessFairLaunchMintedStateNoPayInfo(fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
	// @dev: 1.pay fee
	paidId, err := PayMintFee(fairLaunchMintedInfo.UserID, fairLaunchMintedInfo.MintedFeeRateSatPerKw)
	if err != nil {
		return nil
	}
	// @dev: 2.Store paidId
	err = UpdateFairLaunchMintedInfoPaidId(fairLaunchMintedInfo, paidId)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateFairLaunchMintedInfoPaidId")
	}
	// @dev: 3.Change state
	err = ChangeFairLaunchMintedInfoState(fairLaunchMintedInfo, models.FairLaunchMintedStatePaidPending)
	if err != nil {
		return utils.AppendErrorInfo(err, "ChangeFairLaunchMintedInfoState")
	}
	return nil
}

func ProcessFairLaunchMintedStatePaidPendingInfo(fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
	// @dev: 1.fee paid
	if IsMintFeePaid(fairLaunchMintedInfo.MintFeePaidID) {
		// @dev: Change state
		err = ChangeFairLaunchMintedInfoStateAndUpdatePaidSuccessTime(fairLaunchMintedInfo)
		if err != nil {
			return utils.AppendErrorInfo(err, "ChangeFairLaunchMintedInfoStateAndUpdatePaidSuccessTime")
		}
		return nil
	}
	// @dev: fee has not been paid
	return nil
}

func ProcessFairLaunchMintedStatePaidNoSendInfo(fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
	// @dev: Locked Inventory
	lockedInventory, err := LockInventoryByFairLaunchMintedInfo(fairLaunchMintedInfo)
	if err != nil {
		return utils.AppendErrorInfo(err, "LockInventoryByFairLaunchMintedInfo")
	}
	// @dev: Calculate mint amount
	calculatedMintAmount := CalculateMintAmountByFairLaunchInventoryInfos(lockedInventory)
	if calculatedMintAmount != fairLaunchMintedInfo.AddrAmount {
		err = errors.New("calculated amount is not equal fairLaunchMintedInfo's addr amount")
		return err
	}
	// @dev: Change state
	err = ChangeFairLaunchMintedInfoState(fairLaunchMintedInfo, models.FairLaunchMintedStateSentPending)
	if err != nil {
		return utils.AppendErrorInfo(err, "ChangeFairLaunchMintedInfoState")
	}
	return nil
}

func ProcessFairLaunchMintedStateSentPendingInfo(fairLaunchMintedInfo *models.FairLaunchMintedInfo) (err error) {
	if fairLaunchMintedInfo.OutpointTxHash == "" {
		err = errors.New("no outpoint of transaction hash generated, asset may has not been sent")
		return err
	}
	// @dev: 1.Is Transaction Confirmed
	if IsTransactionConfirmed(fairLaunchMintedInfo.OutpointTxHash) {
		// @dev: Change state
		err = ChangeFairLaunchMintedInfoState(fairLaunchMintedInfo, models.FairLaunchMintedStateSent)
		if err != nil {
			return utils.AppendErrorInfo(err, "ChangeFairLaunchMintedInfoState")
		}
		// @dev: Update MintedNumber and IsMintAll
		err = UpdateMintedNumberAndIsMintAllOfFairLaunchInfoByFairLaunchMintedInfo(fairLaunchMintedInfo)
		if err != nil {
			return utils.AppendErrorInfo(err, "UpdateMintedNumberAndIsMintAllOfFairLaunchInfoByFairLaunchMintedInfo")
		}
		// @dev: Update Inventory
		err = UpdateLockedInventoryByFairLaunchMintedInfo(fairLaunchMintedInfo)
		if err != nil {
			return utils.AppendErrorInfo(err, "UpdateLockedInventoryByFairLaunchMintedInfo")
		}
		// @dev: Update minted user
		f := FairLaunchStore{DB: middleware.DB}
		err = f.CreateFairLaunchMintedUserInfo(&models.FairLaunchMintedUserInfo{
			UserID:                 fairLaunchMintedInfo.UserID,
			FairLaunchMintedInfoID: int(fairLaunchMintedInfo.ID),
			FairLaunchInfoID:       fairLaunchMintedInfo.FairLaunchInfoID,
			MintedNumber:           fairLaunchMintedInfo.MintedNumber,
		})
		if err != nil {
			return utils.AppendErrorInfo(err, "CreateFairLaunchMintedUserInfo")
		}
		err = CreateBalance(&models.Balance{
			AccountId:   AdminAccountId,
			BillType:    models.BILL_TYPE_ASSET_MINTED_SEND,
			Away:        models.AWAY_OUT,
			Amount:      float64(fairLaunchMintedInfo.AddrAmount),
			Unit:        models.UNIT_ASSET_NORMAL,
			Invoice:     &(fairLaunchMintedInfo.EncodedAddr),
			PaymentHash: &(fairLaunchMintedInfo.OutpointTxHash),
			State:       models.STATE_SUCCESS,
		})
		if err != nil {
			return utils.AppendErrorInfo(err, "CreateBalance")
		}
		return nil
	}
	// @dev: Transaction has not been Confirmed
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
	var fairLaunchInfos []models.FairLaunchInfo
	// @dev: Add more condition
	err := middleware.DB.Where("status = ? AND is_mint_all = ? AND state >= ?", models.StatusNormal, false, models.FairLaunchStateIssued).Find(&fairLaunchInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchInfos")
	}
	return &fairLaunchInfos, nil
}

func GetOwnFairLaunchInfosByUserId(id int) (*[]models.FairLaunchInfo, error) {
	var fairLaunchInfos []models.FairLaunchInfo
	err := middleware.DB.Where("status = ? AND user_id = ?", models.StatusNormal, id).Find(&fairLaunchInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchInfos")
	}
	return &fairLaunchInfos, nil
}

func GetOwnFairLaunchMintedInfosByUserId(id int) (*[]models.FairLaunchMintedInfo, error) {
	var fairLaunchMintedInfos []models.FairLaunchMintedInfo
	err := middleware.DB.Where("status = ? AND user_id = ?", models.StatusNormal, id).Find(&fairLaunchMintedInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchMintedInfos")
	}
	return &fairLaunchMintedInfos, nil
}

func ProcessSendFairLaunchReservedResponse(response *taprpc.SendAssetResponse) (txid string) {
	txid, _ = GetTransactionAndIndexByOutpoint(response.Transfer.Outputs[0].Anchor.Outpoint)
	return txid
}

func UpdateFairLaunchInfoIsReservedSent(fairLaunchInfo *models.FairLaunchInfo, txid string) (err error) {
	fairLaunchInfo.IsReservedSent = true
	fairLaunchInfo.ReservedSentAnchorOutpointTxid = txid
	fairLaunchInfo.State = models.FairLaunchStateReservedSentPending
	f := FairLaunchStore{DB: middleware.DB}
	return f.UpdateFairLaunchInfo(fairLaunchInfo)
}

func GetFairLaunchInfoByAssetId(assetId string) (*models.FairLaunchInfo, error) {
	var fairLaunchInfo models.FairLaunchInfo
	err := middleware.DB.Where("asset_id = ?", assetId).First(&fairLaunchInfo).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "First fairLaunchInfo")
	}
	return &fairLaunchInfo, nil
}
