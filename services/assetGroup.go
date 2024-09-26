package services

import (
	"trade/models"
	"trade/services/btldb"
)

func CreateAssetGroup(assetGroup *models.AssetGroup) error {
	return btldb.CreateAssetGroup(assetGroup)
}

func ReadAssetGroupByTweakedGroupKey(tweakedGroupKey string) (*models.AssetGroup, error) {
	return btldb.ReadAssetGroupByTweakedGroupKey(tweakedGroupKey)
}

func GetAssetGroup(tweakedGroupKey string) (*models.AssetGroup, error) {
	return ReadAssetGroupByTweakedGroupKey(tweakedGroupKey)
}

func ProcessAssetGroupSetRequest(userId int, username string, assetGroupSetRequest *models.AssetGroupSetRequest) *models.AssetGroup {
	if assetGroupSetRequest == nil {
		return nil
	}
	return &models.AssetGroup{
		TweakedGroupKey: assetGroupSetRequest.TweakedGroupKey,
		FirstAssetMeta:  assetGroupSetRequest.FirstAssetMeta,
		FirstAssetId:    assetGroupSetRequest.FirstAssetId,
		DeviceId:        assetGroupSetRequest.DeviceId,
		UserId:          userId,
		Username:        username,
	}
}

func SetAssetGroup(userId int, username string, assetGroupSetRequest *models.AssetGroupSetRequest) error {
	assetGroup := ProcessAssetGroupSetRequest(userId, username, assetGroupSetRequest)
	return btldb.CreateAssetGroup(assetGroup)
}
