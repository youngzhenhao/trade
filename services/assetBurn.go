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
