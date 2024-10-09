package services

import (
	"trade/api"
	"trade/btlLog"
	"trade/models"
	"trade/services/btldb"
	"trade/utils"
)

func CreateNftPresale(nftPresale *models.NftPresale) error {
	return btldb.CreateNftPresale(nftPresale)
}

func CreateNftPresales(nftPresales *[]models.NftPresale) error {
	return btldb.CreateNftPresales(nftPresales)
}

func ReadNftPresale(id uint) (*models.NftPresale, error) {
	return btldb.ReadNftPresale(id)
}

func ReadNftPresalesByAssetId(assetId string) (*models.NftPresale, error) {
	return btldb.ReadNftPresalesByAssetId(assetId)
}

func ReadAllNftPresales() (*[]models.NftPresale, error) {
	return btldb.ReadAllNftPresales()
}

func ReadNftPresalesByNftPresaleState(nftPresaleState models.NftPresaleState) (*[]models.NftPresale, error) {
	return btldb.ReadNftPresalesByNftPresaleState(nftPresaleState)
}

func ReadNftPresalesBetweenNftPresaleState(stateStart models.NftPresaleState, stateEnd models.NftPresaleState) (*[]models.NftPresale, error) {
	return btldb.ReadNftPresalesBetweenNftPresaleState(stateStart, stateEnd)
}

func ReadNftPresalesByBuyerUserId(userId int) (*[]models.NftPresale, error) {
	return btldb.ReadNftPresalesByBuyerUserId(userId)
}

func UpdateNftPresale(nftPresale *models.NftPresale) error {
	return btldb.UpdateNftPresale(nftPresale)
}

func UpdateNftPresales(nftPresales *[]models.NftPresale) error {
	return btldb.UpdateNftPresales(nftPresales)
}

func DeleteNftPresale(id uint) error {
	return btldb.DeleteNftPresale(id)
}

func ProcessNftPresale(nftPresaleSetRequest *models.NftPresaleSetRequest) *models.NftPresale {

	var assetId string
	assetId = nftPresaleSetRequest.AssetId
	var name string
	var assetType string
	var groupKey string
	var amount int

	assetInfo, err := api.GetAssetInfo(assetId)
	if err != nil {
		// @dev: Do not return
		btlLog.PreSale.Error("api GetAssetInfo(AssetLeaves) err")
	} else {
		name = assetInfo.Name
		assetType = assetInfo.AssetType.String()
		groupKey = assetInfo.TweakedGroupKey
		amount = assetInfo.Amount
	}
	var meta string
	assetMeta, err := api.FetchAssetMetaByAssetId(assetId)
	if err != nil {
		// @dev: Do not return
		btlLog.PreSale.Error("api FetchAssetMetaByAssetId err")
	} else {
		meta = assetMeta.String()
	}

	return &models.NftPresale{
		AssetId:    assetId,
		Name:       name,
		AssetType:  assetType,
		Meta:       meta,
		GroupKey:   groupKey,
		Amount:     amount,
		Price:      nftPresaleSetRequest.Price,
		Info:       nftPresaleSetRequest.Info,
		LaunchTime: utils.GetTimestamp(),
		State:      models.NftPresaleStateLaunched,
	}
}
