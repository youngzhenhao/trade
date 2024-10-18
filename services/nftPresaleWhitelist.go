package services

import (
	"errors"
	"strconv"
	"trade/btlLog"
	"trade/models"
	"trade/services/btldb"
	"trade/utils"
)

func CreateNftPresaleWhitelist(nftPresaleWhitelist *models.NftPresaleWhitelist) error {
	return btldb.CreateNftPresaleWhitelist(nftPresaleWhitelist)
}

func CreateNftPresaleWhitelists(nftPresaleWhitelists *[]models.NftPresaleWhitelist) error {
	return btldb.CreateNftPresaleWhitelists(nftPresaleWhitelists)
}

func ReadNftPresaleWhitelist(id uint) (*models.NftPresaleWhitelist, error) {
	return btldb.ReadNftPresaleWhitelist(id)
}

func ReadNftPresaleWhitelistsByAssetId(assetId string) (*[]models.NftPresaleWhitelist, error) {
	return btldb.ReadNftPresaleWhitelistsByAssetId(assetId)
}

func ReadNftPresaleWhitelistsByBatchGroupId(batchGroupId int) (*[]models.NftPresaleWhitelist, error) {
	return btldb.ReadNftPresaleWhitelistsByBatchGroupId(batchGroupId)
}

func ReadNftPresaleWhitelistsByAssetIdOrBatchGroupId(assetId string, batchGroupId int) (*[]models.NftPresaleWhitelist, error) {
	return btldb.ReadNftPresaleWhitelistsByAssetIdOrBatchGroupId(assetId, batchGroupId)
}

func ReadAllNftPresaleWhitelists() (*[]models.NftPresaleWhitelist, error) {
	return btldb.ReadAllNftPresaleWhitelists()
}

func UpdateNftPresaleWhitelist(nftPresaleWhitelist *models.NftPresaleWhitelist) error {
	return btldb.UpdateNftPresaleWhitelist(nftPresaleWhitelist)
}

func UpdateNftPresaleWhitelists(nftPresaleWhitelists *[]models.NftPresaleWhitelist) error {
	return btldb.UpdateNftPresaleWhitelists(nftPresaleWhitelists)
}

func DeleteNftPresaleWhitelist(id uint) error {
	return btldb.DeleteNftPresaleWhitelist(id)
}

// @dev: Get

func GetNftPresaleWhitelistsByAssetId(assetId string) (*[]models.NftPresaleWhitelist, error) {
	return ReadNftPresaleWhitelistsByAssetId(assetId)
}

func GetNftPresaleWhitelistsByBatchGroupId(batchGroupId int) (*[]models.NftPresaleWhitelist, error) {
	return ReadNftPresaleWhitelistsByBatchGroupId(batchGroupId)
}

func GetNftPresaleWhitelistsByAssetIdOrBatchGroupId(assetId string, batchGroupId int) (*[]models.NftPresaleWhitelist, error) {
	return ReadNftPresaleWhitelistsByAssetIdOrBatchGroupId(assetId, batchGroupId)
}

// GetNftPresaleWhitelistsOnlyByAssetId
// @Description: Get nftPresale Whitelists
func GetNftPresaleWhitelistsOnlyByAssetId(assetId string) (*[]string, error) {
	nftPresale, err := GetNftPresaleByAssetId(assetId)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetNftPresaleByAssetId")
	}
	batchGroupId := nftPresale.BatchGroupId
	nftPresaleWhitelists, err := GetNftPresaleWhitelistsByAssetIdOrBatchGroupId(assetId, batchGroupId)
	if err != nil {
		nftPresaleWhitelists = &[]models.NftPresaleWhitelist{}
		btlLog.PreSale.Info("GetNftPresaleWhitelistsByAssetIdOrBatchGroupId err:%v", err)
	}
	var whitelists = make([]string, 0)
	usernameMap := make(map[string]bool)
	for _, batchGroupWhitelist := range *nftPresaleWhitelists {
		usernameMap[batchGroupWhitelist.Username] = true
	}
	for username := range usernameMap {
		whitelists = append(whitelists, username)
	}
	return &whitelists, nil
}

// GetNftPresaleWhitelistsByNftPresale
// @Description: Get nftPresale Whitelists
func GetNftPresaleWhitelistsByNftPresale(nftPresale *models.NftPresale) (*[]string, error) {
	if nftPresale == nil {
		return nil, errors.New("nftPresale is nil")
	}
	assetId := nftPresale.AssetId
	batchGroupId := nftPresale.BatchGroupId
	nftPresaleWhitelists, err := GetNftPresaleWhitelistsByAssetIdOrBatchGroupId(assetId, batchGroupId)
	if err != nil {
		nftPresaleWhitelists = &[]models.NftPresaleWhitelist{}
		btlLog.PreSale.Info("GetNftPresaleWhitelistsByAssetIdOrBatchGroupId err:%v", err)
	}
	var whitelists = make([]string, 0)
	usernameMap := make(map[string]bool)
	for _, batchGroupWhitelist := range *nftPresaleWhitelists {
		usernameMap[batchGroupWhitelist.Username] = true
	}
	for username := range usernameMap {
		whitelists = append(whitelists, username)
	}
	return &whitelists, nil
}

// @dev: Process

func ProcessNftPresaleWhitelistSetRequest(nftPresaleWhitelistSetRequest *models.NftPresaleWhitelistSetRequest) (*models.NftPresaleWhitelist, error) {
	if nftPresaleWhitelistSetRequest == nil {
		return nil, errors.New("nftPresaleWhitelistSetRequest is nil")
	}
	username := nftPresaleWhitelistSetRequest.Username
	if username == "" {
		return nil, errors.New("username is empty")
	}
	userId, err := NameToId(username)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "NameToId")
	}
	whitelistType := nftPresaleWhitelistSetRequest.WhitelistType
	if whitelistType == models.WhitelistTypeAsset {
		assetId := nftPresaleWhitelistSetRequest.AssetId
		if assetId == "" {
			return nil, errors.New("assetId is empty")
		}
		return &models.NftPresaleWhitelist{
			WhitelistType: whitelistType,
			AssetId:       assetId,
			BatchGroupId:  0,
			UserId:        userId,
			Username:      username,
		}, nil
	} else if whitelistType == models.WhitelistTypeGroupBatch {
		batchGroupId := nftPresaleWhitelistSetRequest.BatchGroupId
		if batchGroupId == 0 {
			return nil, errors.New("batchGroupId is 0")
		}
		return &models.NftPresaleWhitelist{
			WhitelistType: whitelistType,
			AssetId:       "",
			BatchGroupId:  batchGroupId,
			UserId:        userId,
			Username:      username,
		}, nil
	} else {
		return nil, errors.New("whitelistType(" + strconv.Itoa(int(whitelistType)) + ") is invalid")
	}
}

func ProcessNftPresaleWhitelistSetRequests(nftPresaleWhitelistSetRequests *[]models.NftPresaleWhitelistSetRequest) (*[]models.NftPresaleWhitelist, error) {
	if nftPresaleWhitelistSetRequests == nil {
		return nil, errors.New("nftPresaleWhitelistSetRequests is nil")
	}
	var nftPresaleWhitelists []models.NftPresaleWhitelist
	for _, nftPresaleWhitelistSetRequest := range *nftPresaleWhitelistSetRequests {
		nftPresaleWhitelist, err := ProcessNftPresaleWhitelistSetRequest(&nftPresaleWhitelistSetRequest)
		if err != nil {
			return nil, utils.AppendErrorInfo(err, "ProcessNftPresaleWhitelistSetRequest")
		}
		nftPresaleWhitelists = append(nftPresaleWhitelists, *nftPresaleWhitelist)
	}
	return &nftPresaleWhitelists, nil
}

// @dev: Add whitelist

func AddWhitelistByRequest(nftPresaleWhitelistSetRequest *models.NftPresaleWhitelistSetRequest) error {
	if nftPresaleWhitelistSetRequest == nil {
		return errors.New("nftPresaleWhitelistSetRequest is nil")
	}
	nftPresaleWhitelist, err := ProcessNftPresaleWhitelistSetRequest(nftPresaleWhitelistSetRequest)
	if err != nil {
		return utils.AppendErrorInfo(err, "ProcessNftPresaleWhitelistSetRequest")
	}
	err = CreateNftPresaleWhitelist(nftPresaleWhitelist)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateNftPresaleWhitelist")
	}
	return nil
}

func AddWhitelistsByRequests(nftPresaleWhitelistSetRequests *[]models.NftPresaleWhitelistSetRequest) error {
	if nftPresaleWhitelistSetRequests == nil {
		return errors.New("nftPresaleWhitelistSetRequests is nil")
	}
	nftPresaleWhitelists, err := ProcessNftPresaleWhitelistSetRequests(nftPresaleWhitelistSetRequests)
	if err != nil {
		return utils.AppendErrorInfo(err, "ProcessNftPresaleWhitelistSetRequest")
	}
	err = CreateNftPresaleWhitelists(nftPresaleWhitelists)
	if err != nil {
		return utils.AppendErrorInfo(err, "CreateNftPresaleWhitelists")
	}
	return nil
}
