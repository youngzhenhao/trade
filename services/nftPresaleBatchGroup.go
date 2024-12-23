package services

import (
	"encoding/hex"
	"errors"
	"math"
	"strconv"
	"trade/api"
	"trade/btlLog"
	"trade/config"
	"trade/middleware"
	"trade/models"
	"trade/services/btldb"
	"trade/utils"
)

func CreateNftPresaleBatchGroup(nftPresaleBatchGroup *models.NftPresaleBatchGroup) error {
	return btldb.CreateNftPresaleBatchGroup(nftPresaleBatchGroup)
}

func CreateNftPresaleBatchGroups(nftPresaleBatchGroups *[]models.NftPresaleBatchGroup) error {
	return btldb.CreateNftPresaleBatchGroups(nftPresaleBatchGroups)
}

func ReadNftPresaleBatchGroup(id uint) (*models.NftPresaleBatchGroup, error) {
	return btldb.ReadNftPresaleBatchGroup(id)
}

func ReadNftPresaleBatchGroupByGroupKey(groupKey string) (*models.NftPresaleBatchGroup, error) {
	return btldb.ReadNftPresaleBatchGroupByGroupKey(groupKey)
}

func ReadAllNftPresaleBatchGroups() (*[]models.NftPresaleBatchGroup, error) {
	return btldb.ReadAllNftPresaleBatchGroups()
}

func ReadSellingNftPresaleBatchGroups() (*[]models.NftPresaleBatchGroup, error) {
	return btldb.ReadSellingNftPresaleBatchGroups()
}

func ReadNotStartNftPresaleBatchGroups() (*[]models.NftPresaleBatchGroup, error) {
	return btldb.ReadNotStartNftPresaleBatchGroups()
}

func ReadEndNftPresaleBatchGroups() (*[]models.NftPresaleBatchGroup, error) {
	return btldb.ReadEndNftPresaleBatchGroups()
}

func GetAllNftPresaleBatchGroup() (*[]models.NftPresaleBatchGroup, error) {
	return ReadAllNftPresaleBatchGroups()
}

func GetSellingNftPresaleBatchGroups() (*[]models.NftPresaleBatchGroup, error) {
	return ReadSellingNftPresaleBatchGroups()
}

func GetNotStartNftPresaleBatchGroups() (*[]models.NftPresaleBatchGroup, error) {
	return ReadNotStartNftPresaleBatchGroups()
}

func GetEndNftPresaleBatchGroups() (*[]models.NftPresaleBatchGroup, error) {
	return ReadEndNftPresaleBatchGroups()
}

func UpdateNftPresaleBatchGroup(nftPresaleBatchGroup *models.NftPresaleBatchGroup) error {
	return btldb.UpdateNftPresaleBatchGroup(nftPresaleBatchGroup)
}

func UpdateNftPresaleBatchGroups(nftPresaleBatchGroups *[]models.NftPresaleBatchGroup) error {
	return btldb.UpdateNftPresaleBatchGroups(nftPresaleBatchGroups)
}

func DeleteNftPresaleBatchGroup(id uint) error {
	return btldb.DeleteNftPresaleBatchGroup(id)
}

func ProcessLaunchRequestNftPresale(nftPresaleSetRequest *models.NftPresaleSetRequest, startTime uint, endTime uint, info string) (*models.NftPresale, error) {
	if startTime == 0 {
		return nil, errors.New("startTime is 0")
	}
	var assetId string
	assetId = nftPresaleSetRequest.AssetId
	if assetId == "" {
		return nil, errors.New("nftPresaleSetRequest.AssetId(" + assetId + ") is null")
	}
	var name string
	var assetType string
	var groupKey string
	var amount int
	var meta string
	_asset, err := api.GetIncludeLeasedAssetById(assetId)
	if err != nil {
		btlLog.PreSale.Error("api GetIncludeLeasedAssetById err:%v", err)
		return nil, err
	} else {
		if _asset.Amount > math.MaxInt {
			return nil, errors.New("amount(" + strconv.FormatUint(_asset.Amount, 10) + ") is too large")
		}
		amount = int(_asset.Amount)
		if _asset.AssetGenesis != nil {
			name = _asset.AssetGenesis.Name
			assetType = _asset.AssetGenesis.AssetType.String()

		}
		if _asset.AssetGroup != nil {
			tweakedGroupKey := _asset.AssetGroup.TweakedGroupKey
			groupKey = hex.EncodeToString(tweakedGroupKey)
		}
	}
	assetMeta, err := api.FetchAssetMetaByAssetId(assetId)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "FetchAssetMetaByAssetId")
	}
	meta = assetMeta.Data

	groupKeyByAssetId, err := api.GetGroupKeyByAssetId(assetId)
	if err != nil {
		btlLog.PreSale.Error("api GetGroupKeyByAssetId err:%v", err)
	} else {
		groupKey = groupKeyByAssetId
	}
	return &models.NftPresale{
		AssetId:    assetId,
		Name:       name,
		AssetType:  assetType,
		Meta:       meta,
		GroupKey:   groupKey,
		Amount:     amount,
		Price:      nftPresaleSetRequest.Price,
		Info:       info,
		LaunchTime: utils.GetTimestamp(),
		StartTime:  int(startTime),
		EndTime:    int(endTime),
		State:      models.NftPresaleStateLaunched,
	}, nil
}

type ProcessLaunchRequestNftPresalesResult struct {
	NftPresales  *[]models.NftPresale `json:"nft_presales"`
	Supply       int                  `json:"supply"`
	LowestPrice  int                  `json:"lowest_price"`
	HighestPrice int                  `json:"highest_price"`
}

func ProcessLaunchRequestNftPresales(nftPresaleSetRequests *[]models.NftPresaleSetRequest, startTime uint, endTime uint, info string, groupKey string) (*ProcessLaunchRequestNftPresalesResult, error) {
	if nftPresaleSetRequests == nil {
		return nil, errors.New("nftPresaleSetRequests is nil")
	}
	if len(*nftPresaleSetRequests) == 0 {
		return nil, errors.New("nftPresaleSetRequests length is 0")
	}
	if startTime == 0 {
		return nil, errors.New("startTime is 0")
	}
	var nftPresales []models.NftPresale
	var supply uint
	price := (*nftPresaleSetRequests)[0].Price
	var lowestPrice = uint(price)
	var highestPrice = uint(price)
	for _, nftPresaleSetRequest := range *nftPresaleSetRequests {
		nftPresale, err := ProcessLaunchRequestNftPresale(&nftPresaleSetRequest, startTime, endTime, info)
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "ProcessLaunchRequestNftPresale")
		}
		if nftPresale.GroupKey != groupKey {
			return nil, errors.New("nftPresale.GroupKey(" + nftPresale.GroupKey + ") is not equal groupKey(" + groupKey + ")")
		}
		nftPresales = append(nftPresales, *(nftPresale))
		supply += 1
		if lowestPrice > uint(nftPresale.Price) {
			lowestPrice = uint(nftPresale.Price)
		}
		if highestPrice < uint(nftPresale.Price) {
			highestPrice = uint(nftPresale.Price)
		}
	}
	var result = ProcessLaunchRequestNftPresalesResult{
		NftPresales:  &nftPresales,
		Supply:       int(supply),
		LowestPrice:  int(lowestPrice),
		HighestPrice: int(highestPrice),
	}
	return &result, nil
}

func ValidateStartAndEndTimeForNftPresale(startTime int, endTime int) error {
	now := utils.GetTimestamp()
	if !(startTime >= now-3600*24) {
		return errors.New("startTime(" + strconv.Itoa(startTime) + ") must be greater than the current time (max time delay a day i.e. 3600*24 seconds.)")
	}
	if endTime != 0 {
		if !(endTime >= startTime+3600*2) {
			return errors.New("end time should be at least two hour after the start time")
		}
		if !(endTime <= now+3600*24*365*2) {
			return errors.New("end time cannot be more than two year from the current time")
		}
	}
	return nil
}

// CreateBatchGroupAndNftPresales
// @Description: create batchGroup and nftPresales
func CreateBatchGroupAndNftPresales(nftPresaleBatchGroup *models.NftPresaleBatchGroup, nftPresales *[]models.NftPresale) error {
	var err error
	if nftPresaleBatchGroup == nil {
		return errors.New("nftPresaleBatchGroup is nil")
	}
	if nftPresales == nil {
		return errors.New("NftPresales is nil")
	}
	tx := middleware.DB.Begin()
	// @dev: 1. Create nftPresaleBatchGroup
	err = tx.Create(nftPresaleBatchGroup).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// @dev: Get ID
	batchGroupId := nftPresaleBatchGroup.ID
	if batchGroupId == 0 {
		tx.Rollback()
		return errors.New("batchGroupId is 0")
	}
	for i := range *nftPresales {
		(*nftPresales)[i].BatchGroupId = int(batchGroupId)
	}
	// @dev: 2. Create nftPresales
	err = tx.Create(nftPresales).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// ProcessNftPresaleBatchGroupLaunchRequestAndCreate
// @Description: Process nftPresaleBatchGroupLaunchRequest and then create records in db
func ProcessNftPresaleBatchGroupLaunchRequestAndCreate(nftPresaleBatchGroupLaunchRequest *models.NftPresaleBatchGroupLaunchRequest) error {
	if nftPresaleBatchGroupLaunchRequest == nil {
		return errors.New("NftPresaleBatchGroupLaunchRequest is nil")
	}
	// @dev: Value copy
	batchGroupSetRequest := nftPresaleBatchGroupLaunchRequest.BatchGroupSetRequest
	// @dev: GroupKey and GroupName
	groupKey := batchGroupSetRequest.GroupKey
	if groupKey == "" {
		return errors.New("GroupKey is empty")
	}
	var groupName string
	// @dev: network
	network, err := api.NetworkStringToNetwork(config.GetLoadConfig().NetWork)
	if err != nil {
		return utils.AppendErrorInfo(err, "NetworkStringToNetwork")
	} else {
		groupName, err = GetGroupNameByGroupKey(network, groupKey)
		if err != nil {
			return utils.AppendErrorInfo(err, "GetGroupNameByGroupKey")
		}
	}
	// @dev: startTime and endTime
	{
		//if batchGroupSetRequest.StartTime == "" {
		//	return errors.New("StartTime is empty")
		//}
		//var startTime int
		//start, err := utils.DateTimeStringToTime(batchGroupSetRequest.StartTime)
		//if err != nil {
		//	return utils.AppendErrorInfo(err, "DateTimeStringToTime")
		//}
		//startTime = int(uint(start.Unix()))
		//var endTime int
		//if batchGroupSetRequest.EndTime != "" {
		//	end, err := utils.DateTimeStringToTime(batchGroupSetRequest.EndTime)
		//	if err != nil {
		//		return utils.AppendErrorInfo(err, "DateTimeStringToTime")
		//	}
		//	endTime = int(end.Unix())
		//}
	}
	startTime := batchGroupSetRequest.StartTime
	endTime := batchGroupSetRequest.EndTime
	if startTime == 0 {
		// @dev: Do not return error, but set to now timestamp
		startTime = utils.GetTimestamp()
		err = errors.New("start time is invalid(" + strconv.Itoa(startTime) + ")")
		btlLog.PreSale.Info("%v, it has been set to now", err)
	}
	// @dev: Validate startTime and endTime
	err = ValidateStartAndEndTimeForNftPresale(startTime, endTime)
	if err != nil {
		return utils.AppendErrorInfo(err, "ValidateStartAndEndTimeForNftPresale")
	}
	// @dev: Info
	info := batchGroupSetRequest.Info
	// @dev: Pointer
	nftPresaleSetRequests := nftPresaleBatchGroupLaunchRequest.NftPresaleSetRequests
	if nftPresaleSetRequests == nil {
		return errors.New("NftPresales is nil")
	} else if len(*nftPresaleSetRequests) == 0 {
		return errors.New("NftPresales is empty (length is 0)")
	}
	// @dev: Process NftPresales, result includes NftPresales, Supply, LowestPrice, HighestPrice
	processedResult, err := ProcessLaunchRequestNftPresales(nftPresaleSetRequests, uint(startTime), uint(endTime), info, groupKey)
	if err != nil {
		return utils.AppendErrorInfo(err, "ProcessLaunchRequestNftPresales")
	}
	if processedResult.Supply == 0 {
		return errors.New("processedResult.Supply is 0")
	}
	if processedResult.LowestPrice == 0 {
		return errors.New("processedResult.LowestPrice is 0")
	}
	if processedResult.HighestPrice == 0 {
		return errors.New("processedResult.HighestPrice is 0")
	}
	if processedResult.NftPresales == nil {
		return errors.New("processedResult.NftPresales is nil")
	}
	// @dev: NftPresaleBatchGroup
	var nftPresaleBatchGroup = models.NftPresaleBatchGroup{
		GroupKey:     groupKey,
		GroupName:    groupName,
		SoldNumber:   0,
		Supply:       processedResult.Supply,
		LowestPrice:  processedResult.LowestPrice,
		HighestPrice: processedResult.HighestPrice,
		StartTime:    startTime,
		EndTime:      endTime,
		Info:         info,
	}
	nftPresales := processedResult.NftPresales
	// @dev: Create records in db
	err = CreateBatchGroupAndNftPresales(&nftPresaleBatchGroup, nftPresales)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateBatchGroupAndNftPresales")
	}
	return nil
}

// GetBatchGroups
// @Description: Get batch groups
func GetBatchGroups(state models.NftPresaleBatchGroupState) (*[]models.NftPresaleBatchGroup, error) {
	var nftPresaleBatchGroups *[]models.NftPresaleBatchGroup
	var err error
	switch state {
	case models.NftPresaleBatchGroupStateAll:
		nftPresaleBatchGroups, err = GetAllNftPresaleBatchGroup()
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "GetAllNftPresaleBatchGroup")
		}
	case models.NftPresaleBatchGroupStateSelling:
		nftPresaleBatchGroups, err = GetSellingNftPresaleBatchGroups()
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "GetSellingNftPresaleBatchGroups")
		}
	case models.NftPresaleBatchGroupStateNotStart:
		nftPresaleBatchGroups, err = GetNotStartNftPresaleBatchGroups()
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "GetNotStartNftPresaleBatchGroups")
		}
	case models.NftPresaleBatchGroupStateEnd:
		nftPresaleBatchGroups, err = GetEndNftPresaleBatchGroups()
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "GetEndNftPresaleBatchGroups")
		}
	default:
		nftPresaleBatchGroups, err = GetAllNftPresaleBatchGroup()
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "GetAllNftPresaleBatchGroup")
		}
	}
	if nftPresaleBatchGroups == nil {
		nftPresaleBatchGroups = &[]models.NftPresaleBatchGroup{}
	}
	return nftPresaleBatchGroups, nil
}

// @dev: Simplify

func NftPresaleBatchGroupToNftPresaleBatchGroupSimplified(nftPresaleBatchGroup *models.NftPresaleBatchGroup) *models.NftPresaleBatchGroupSimplified {
	if nftPresaleBatchGroup == nil {
		return nil
	}
	assetId, err := GetFirstAssetIdByBatchGroupId(int(nftPresaleBatchGroup.ID))
	if err != nil {
		btlLog.PreSale.Error("%v", err)
	}
	return &models.NftPresaleBatchGroupSimplified{
		ID:           nftPresaleBatchGroup.ID,
		UpdatedAt:    nftPresaleBatchGroup.UpdatedAt,
		GroupKey:     nftPresaleBatchGroup.GroupKey,
		GroupName:    nftPresaleBatchGroup.GroupName,
		SoldNumber:   nftPresaleBatchGroup.SoldNumber,
		Supply:       nftPresaleBatchGroup.Supply,
		LowestPrice:  nftPresaleBatchGroup.LowestPrice,
		HighestPrice: nftPresaleBatchGroup.HighestPrice,
		StartTime:    nftPresaleBatchGroup.StartTime,
		EndTime:      nftPresaleBatchGroup.EndTime,
		Info:         nftPresaleBatchGroup.Info,
		FirstAssetId: assetId,
	}
}

func NftPresaleBatchGroupSliceToNftPresaleBatchGroupSimplifiedSlice(nftPresaleBatchGroups *[]models.NftPresaleBatchGroup) *[]models.NftPresaleBatchGroupSimplified {
	if nftPresaleBatchGroups == nil {
		return nil
	}
	var nftPresaleBatchGroupSimplifiedSlice []models.NftPresaleBatchGroupSimplified
	for _, nftPresaleBatchGroup := range *nftPresaleBatchGroups {
		nftPresaleBatchGroupSimplifiedSlice = append(nftPresaleBatchGroupSimplifiedSlice, *(NftPresaleBatchGroupToNftPresaleBatchGroupSimplified(&nftPresaleBatchGroup)))
	}
	return &nftPresaleBatchGroupSimplifiedSlice
}
