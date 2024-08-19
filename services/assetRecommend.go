package services

import (
	"errors"
	"time"
	"trade/models"
	"trade/services/btldb"
	"trade/utils"
)

func GetAssetRecommendsByUserId(userId int) (*[]models.AssetRecommend, error) {
	return btldb.ReadAssetRecommendsByUserId(userId)
}

func GetAssetRecommendsByAssetId(assetId string) (*[]models.AssetRecommend, error) {
	return btldb.ReadAssetRecommendsByAssetId(assetId)
}

func GetAssetRecommendByUserIdAndAssetId(userId int, assetId string) (*models.AssetRecommend, error) {
	return btldb.ReadAssetRecommendByUserIdAndAssetId(userId, assetId)
}

func ProcessAssetRecommendSetRequest(userId int, username string, assetRecommendSetRequest models.AssetRecommendSetRequest) models.AssetRecommend {
	if assetRecommendSetRequest.RecommendUserId == 0 {
		assetRecommendSetRequest.RecommendUserId = userId
	}
	if assetRecommendSetRequest.RecommendUsername == "" {
		assetRecommendSetRequest.RecommendUsername = username
	}
	if assetRecommendSetRequest.RecommendTime == 0 {
		assetRecommendSetRequest.RecommendTime = utils.GetTimestamp()
	}
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
	return btldb.CreateAssetRecommend(assetRecommend)
}

func GetAllAssetRecommendsUpdatedAtDesc() (*[]models.AssetRecommend, error) {
	return btldb.ReadAllAssetRecommendsUpdatedAtDesc()
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

func IsAssetRecommendChanged(assetRecommendByInvoice *models.AssetRecommend, old *models.AssetRecommend) bool {
	if assetRecommendByInvoice == nil || old == nil {
		return true
	}
	if assetRecommendByInvoice.AssetId != old.AssetId {
		return true
	}
	if assetRecommendByInvoice.AssetFromAddr != old.AssetFromAddr {
		return true
	}
	if assetRecommendByInvoice.RecommendUserId != old.RecommendUserId {
		return true
	}
	if assetRecommendByInvoice.RecommendUsername != old.RecommendUsername {
		return true
	}
	//if assetRecommendByInvoice.RecommendTime != old.RecommendTime {
	//	return true
	//}
	if assetRecommendByInvoice.DeviceId != old.DeviceId {
		return true
	}
	if assetRecommendByInvoice.UserId != old.UserId {
		return true
	}
	if assetRecommendByInvoice.Username != old.Username {
		return true
	}
	return false
}

func CheckAssetRecommendIfUpdate(assetRecommend *models.AssetRecommend) (*models.AssetRecommend, error) {
	if assetRecommend == nil {
		return nil, errors.New("nil asset recommend")
	}
	assetRecommendByUserIdAndAssetId, err := GetAssetRecommendByUserIdAndAssetId(assetRecommend.UserId, assetRecommend.AssetId)
	if err != nil {
		return assetRecommend, nil
	}
	if !IsAssetRecommendChanged(assetRecommendByUserIdAndAssetId, assetRecommend) {
		return assetRecommendByUserIdAndAssetId, nil
	}
	assetRecommendByUserIdAndAssetId.AssetId = assetRecommend.AssetId
	assetRecommendByUserIdAndAssetId.AssetFromAddr = assetRecommend.AssetFromAddr
	assetRecommendByUserIdAndAssetId.RecommendUserId = assetRecommend.RecommendUserId
	assetRecommendByUserIdAndAssetId.RecommendUsername = assetRecommend.RecommendUsername
	// @dev: Do not update Recommend Time
	//assetRecommendByUserIdAndAssetId.RecommendTime = assetRecommend.RecommendTime
	assetRecommendByUserIdAndAssetId.DeviceId = assetRecommend.DeviceId
	assetRecommendByUserIdAndAssetId.UserId = assetRecommend.UserId
	return assetRecommendByUserIdAndAssetId, nil
}

func CreateOrUpdateAssetRecommend(assetRecommend *models.AssetRecommend) (err error) {
	var assetRecommendUpdate *models.AssetRecommend
	assetRecommendUpdate, err = CheckAssetRecommendIfUpdate(assetRecommend)
	return btldb.UpdateAssetRecommend(assetRecommendUpdate)
}
