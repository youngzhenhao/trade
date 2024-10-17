package services

import (
	"errors"
	"trade/api"
	"trade/btlLog"
	"trade/config"
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

func ProcessLaunchRequestNftPresale(nftPresaleSetRequest *models.NftPresaleSetRequest, nftPresaleBatchGroup *models.NftPresaleBatchGroup) (*models.NftPresale, error) {
	if nftPresaleBatchGroup == nil {
		return nil, errors.New("batchGroupSetRequest is nil")
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
		LaunchTime: utils.GetTimestamp(),
		State:      models.NftPresaleStateLaunched,
	}, nil
}

func ProcessLaunchRequestNftPresales(nftPresaleSetRequests *[]models.NftPresaleSetRequest, nftPresaleBatchGroup *models.NftPresaleBatchGroup) (*[]models.NftPresale, error) {
	if nftPresaleSetRequests == nil {
		return nil, errors.New("nftPresaleSetRequests is nil")
	}
	var nftPresales []models.NftPresale
	for _, nftPresaleSetRequest := range *nftPresaleSetRequests {
		nftPresale, err := ProcessLaunchRequestNftPresale(&nftPresaleSetRequest, nftPresaleBatchGroup)
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "ProcessLaunchRequestNftPresale")
		}
		nftPresales = append(nftPresales, *(nftPresale))
	}
	return &nftPresales, nil
}

func ProcessNftPresaleBatchGroupLaunchRequest(nftPresaleBatchGroupLaunchRequest *models.NftPresaleBatchGroupLaunchRequest) error {
	if nftPresaleBatchGroupLaunchRequest == nil {
		return errors.New("NftPresaleBatchGroupLaunchRequest is nil")
	}

	// @dev: Value copy
	batchGroupSetRequest := nftPresaleBatchGroupLaunchRequest.BatchGroupSetRequest
	if batchGroupSetRequest.GroupKey == "" {
		return errors.New("GroupKey is empty")
	}
	var groupName string
	network, err := api.NetworkStringToNetwork(config.GetLoadConfig().NetWork)
	if err != nil {
		return utils.AppendErrorInfo(err, "NetworkStringToNetwork")
	} else {
		groupName, err = GetGroupNameByGroupKey(network, batchGroupSetRequest.GroupKey)
		if err != nil {
			return utils.AppendErrorInfo(err, "GetGroupNameByGroupKey")
		}
	}
	if batchGroupSetRequest.StartTime == "" {
		return errors.New("StartTime is empty")
	}
	var startTime int
	start, err := utils.DateTimeStringToTime(batchGroupSetRequest.StartTime)
	if err != nil {
		return utils.AppendErrorInfo(err, "DateTimeStringToTime")
	}
	startTime = int(start.Unix())
	var endTime int
	if batchGroupSetRequest.EndTime != "" {
		end, err := utils.DateTimeStringToTime(batchGroupSetRequest.EndTime)
		if err != nil {
			return utils.AppendErrorInfo(err, "DateTimeStringToTime")
		}
		endTime = int(end.Unix())
	}
	// @dev: Pointer
	nftPresaleSetRequests := nftPresaleBatchGroupLaunchRequest.NftPresaleSetRequests
	if nftPresaleSetRequests == nil {
		return errors.New("NftPresales is nil")
	} else if len(*nftPresaleSetRequests) == 0 {
		return errors.New("NftPresales is empty (length is 0)")
	}
	var nftPresales []models.NftPresaleSetRequest

	_ = nftPresales
	// TODO

	var nftPresaleBatchGroup = models.NftPresaleBatchGroup{
		GroupKey:     batchGroupSetRequest.GroupKey,
		GroupName:    groupName,
		SoldNumber:   0,
		Supply:       0,
		LowestPrice:  0,
		HighestPrice: 0,
		StartTime:    startTime,
		EndTime:      endTime,
		Info:         batchGroupSetRequest.Info,
	}
	// TODO: Modify ProcessLaunchRequestNftPresales's logic and params
	ProcessLaunchRequestNftPresales(nftPresaleSetRequests, &nftPresaleBatchGroup)

	// TODO: save data in db by tx
	return nil
}
