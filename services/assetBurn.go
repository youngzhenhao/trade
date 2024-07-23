package services

import (
	"trade/models"
)

func GetAssetBurnsByUserId(userId int) (*[]models.AssetBurn, error) {
	return ReadAssetBurnsByUserId(userId)
}

func ProcessAssetBurnSetRequest(userId int, username string, assetBurnSetRequest *models.AssetBurnSetRequest) *models.AssetBurn {
	var assetBurn models.AssetBurn
	assetBurn = models.AssetBurn{
		AssetId:  assetBurnSetRequest.AssetId,
		Amount:   assetBurnSetRequest.Amount,
		DeviceId: assetBurnSetRequest.DeviceId,
		UserId:   userId,
		Username: username,
	}
	return &assetBurn
}

func GetAssetBurnsByAssetId(assetId string) (*[]models.AssetBurn, error) {
	return ReadAssetBurnsByAssetId(assetId)
}

type AssetBurnTotal struct {
	AssetId     string `json:"asset_id"`
	TotalAmount int    `json:"total_amount"`
}

func AssetBurnsToAssetBurnTotal(assetBurns *[]models.AssetBurn) *AssetBurnTotal {
	if assetBurns == nil || len(*assetBurns) == 0 {
		return nil
	}
	var assetBurnTotal AssetBurnTotal
	assetBurnTotal.AssetId = (*assetBurns)[0].AssetId
	for _, assetBurn := range *assetBurns {
		assetBurnTotal.TotalAmount += assetBurn.Amount
	}
	return &assetBurnTotal
}

// @dev: Get total burn amount by asset id
func GetAssetBurnTotal(assetId string) (*AssetBurnTotal, error) {
	assetBurns, err := GetAssetBurnsByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	assetBurnTotal := AssetBurnsToAssetBurnTotal(assetBurns)
	if assetBurnTotal == nil {
		return &AssetBurnTotal{
			AssetId:     assetId,
			TotalAmount: 0,
		}, nil
	}
	return assetBurnTotal, nil
}
