package services

import (
	"errors"
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
	assetInfo, err := api.GetAssetInfoApi(assetId)
	if err != nil {
		// @dev: Do not return
		btlLog.PreSale.Error("api GetAssetInfoApi err:%v", err)
	} else {
		name = assetInfo.Name
		assetType = assetInfo.AssetType
		if assetInfo.GroupKey != nil {
			groupKey = *assetInfo.GroupKey
		}
		amount = int(assetInfo.Amount)
		if assetInfo.Meta != nil {
			meta = *assetInfo.Meta
		}
	}
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
	if !(startTime >= now-600) {
		return errors.New("startTime(" + strconv.Itoa(startTime) + ") must be greater than the current time (max time delay 600 seconds.)")
	}
	if endTime != 0 {
		if !(endTime >= startTime+3600*2) {
			return errors.New("end time should be at least two hour after the start time")
		}
		if !(endTime <= now+3600*24*365) {
			return errors.New("end time cannot be more than one year from the current time")
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
	if batchGroupSetRequest.StartTime == "" {
		return errors.New("StartTime is empty")
	}
	var startTime int
	start, err := utils.DateTimeStringToTime(batchGroupSetRequest.StartTime)
	if err != nil {
		return utils.AppendErrorInfo(err, "DateTimeStringToTime")
	}
	startTime = int(uint(start.Unix()))
	var endTime int
	if batchGroupSetRequest.EndTime != "" {
		end, err := utils.DateTimeStringToTime(batchGroupSetRequest.EndTime)
		if err != nil {
			return utils.AppendErrorInfo(err, "DateTimeStringToTime")
		}
		endTime = int(end.Unix())
	}
	if startTime == 0 {
		return errors.New("start time is invalid(" + strconv.Itoa(startTime) + ")")
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
