package services

import (
	"time"
	"trade/models"
)

func GetAssetRecommendsByUserId(userId int) (*[]models.AssetRecommend, error) {
	return ReadAssetRecommendsByUserId(userId)
}

func GetAssetRecommendsByAssetId(assetId string) (*[]models.AssetRecommend, error) {
	return ReadAssetRecommendsByAssetId(assetId)
}

func GetAssetRecommendByUserIdAndAssetId(userId int, assetId string) (*models.AssetRecommend, error) {
	return ReadAssetRecommendByUserIdAndAssetId(userId, assetId)
}

func ProcessAssetRecommendSetRequest(userId int, username string, assetRecommendSetRequest models.AssetRecommendSetRequest) models.AssetRecommend {
	return models.AssetRecommend{
		AssetId:           assetRecommendSetRequest.AssetId,
		AssetFromAddr:     assetRecommendSetRequest.AssetFromAddr,
		RecommendUserId:   assetRecommendSetRequest.RecommendUserId,
		RecommendUsername: assetRecommendSetRequest.RecommendUsername,
		RecommendTime:     assetRecommendSetRequest.RecommendTime,
		DeviceId:          assetRecommendSetRequest.DeviceId,
		UserId:            userId,
		Username:          username,
	}
}

func SetAssetRecommend(assetRecommend *models.AssetRecommend) error {
	return CreateAssetRecommend(assetRecommend)
}

func GetAllAssetRecommendsUpdatedAtDesc() (*[]models.AssetRecommend, error) {
	return ReadAllAssetRecommendsUpdatedAtDesc()
}

type AssetRecommendSimplified struct {
	UpdatedAt         time.Time `json:"updated_at"`
	AssetId           string    `json:"asset_id" gorm:"type:varchar(255)"`
	AssetFromAddr     string    `json:"asset_from_addr" gorm:"type:varchar(255)"`
	RecommendUsername string    `json:"recommend_username" gorm:"type:varchar(255)"`
	RecommendTime     int       `json:"recommend_time"`
	DeviceId          string    `json:"device_id" gorm:"type:varchar(255)"`
	Username          string    `json:"username" gorm:"type:varchar(255)"`
}

func AssetRecommendToAssetRecommendSimplified(assetRecommend models.AssetRecommend) AssetRecommendSimplified {
	return AssetRecommendSimplified{
		UpdatedAt:         assetRecommend.UpdatedAt,
		AssetId:           assetRecommend.AssetId,
		AssetFromAddr:     assetRecommend.AssetFromAddr,
		RecommendUsername: assetRecommend.RecommendUsername,
		RecommendTime:     assetRecommend.RecommendTime,
		DeviceId:          assetRecommend.DeviceId,
		Username:          assetRecommend.Username,
	}
}

func AssetRecommendSliceToAssetRecommendSimplifiedSlice(assetRecommends *[]models.AssetRecommend) *[]AssetRecommendSimplified {
	if assetRecommends == nil {
		return nil
	}
	var assetRecommendSimplified []AssetRecommendSimplified
	for _, assetRecommend := range *assetRecommends {
		assetRecommendSimplified = append(assetRecommendSimplified, AssetRecommendToAssetRecommendSimplified(assetRecommend))
	}
	return &assetRecommendSimplified
}

func GetAllAssetRecommendSimplified() (*[]AssetRecommendSimplified, error) {
	assetRecommends, err := GetAllAssetRecommendsUpdatedAtDesc()
	if err != nil {
		return nil, err
	}
	return AssetRecommendSliceToAssetRecommendSimplifiedSlice(assetRecommends), nil
}
