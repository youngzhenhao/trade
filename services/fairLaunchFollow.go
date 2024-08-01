package services

import (
	"errors"
	"time"
	"trade/models"
	"trade/services/btldb"
)

func ProcessFairLaunchFollowSetRequest(userId int, username string, fairLaunchFollowRequest models.FairLaunchFollowSetRequest) models.FairLaunchFollow {
	var fairLaunchFollow models.FairLaunchFollow
	fairLaunchFollow = models.FairLaunchFollow{
		FairLaunchInfoId: fairLaunchFollowRequest.FairLaunchInfoId,
		AssetId:          fairLaunchFollowRequest.AssetId,
		DeviceId:         fairLaunchFollowRequest.DeviceId,
		UserId:           userId,
		Username:         username,
	}
	return fairLaunchFollow
}

func ProcessFairLaunchFollowSetRequests(userId int, username string, fairLaunchFollowRequests *[]models.FairLaunchFollowSetRequest) *[]models.FairLaunchFollow {
	var fairLaunchFollows []models.FairLaunchFollow
	for _, fairLaunchFollowRequest := range *fairLaunchFollowRequests {
		fairLaunchFollow := ProcessFairLaunchFollowSetRequest(userId, username, fairLaunchFollowRequest)
		fairLaunchFollows = append(fairLaunchFollows, fairLaunchFollow)
	}
	return &fairLaunchFollows
}

func GetFairLaunchFollowsByUserId(userId int) (*[]models.FairLaunchFollow, error) {
	return btldb.ReadFairLaunchFollowsByUserId(userId)
}

func GetFairLaunchFollowByUserIdAndAssetId(userId int, assetId string) (*models.FairLaunchFollow, error) {
	return btldb.ReadFairLaunchFollowByUserIdAndAssetId(userId, assetId)
}

func IsFairLaunchFollowed(userId int, assetId string) bool {
	fairLaunchFollow, err := GetFairLaunchFollowByUserIdAndAssetId(userId, assetId)
	if err != nil || fairLaunchFollow == nil {
		return false
	}
	return true
}

func GetFairLaunchFollowByAssetId(assetId string) (*models.FairLaunchFollow, error) {
	return btldb.ReadFairLaunchFollowByAssetId(assetId)
}

func IsFairLaunchFollowChanged(fairLaunchFollowByTxidAndIndex *models.FairLaunchFollow, old *models.FairLaunchFollow) bool {
	if fairLaunchFollowByTxidAndIndex == nil || old == nil {
		return true
	}
	if fairLaunchFollowByTxidAndIndex.FairLaunchInfoId != old.FairLaunchInfoId {
		return true
	}
	if fairLaunchFollowByTxidAndIndex.AssetId != old.AssetId {
		return true
	}
	if fairLaunchFollowByTxidAndIndex.DeviceId != old.DeviceId {
		return true
	}
	if fairLaunchFollowByTxidAndIndex.UserId != old.UserId {
		return true
	}
	if fairLaunchFollowByTxidAndIndex.Username != old.Username {
		return true
	}
	return false
}

func CheckFairLaunchFollowIfUpdate(fairLaunchFollow *models.FairLaunchFollow) (*models.FairLaunchFollow, error) {
	if fairLaunchFollow == nil {
		return nil, errors.New("nil fair launch follow")
	}
	fairLaunchFollowByAssetId, err := GetFairLaunchFollowByAssetId(fairLaunchFollow.AssetId)
	if err != nil {
		return fairLaunchFollow, nil
	}
	if !IsFairLaunchFollowChanged(fairLaunchFollowByAssetId, fairLaunchFollow) {
		return fairLaunchFollowByAssetId, nil
	}
	fairLaunchFollowByAssetId.FairLaunchInfoId = fairLaunchFollow.FairLaunchInfoId
	fairLaunchFollowByAssetId.AssetId = fairLaunchFollow.AssetId
	fairLaunchFollowByAssetId.DeviceId = fairLaunchFollow.DeviceId
	fairLaunchFollowByAssetId.UserId = fairLaunchFollow.UserId
	fairLaunchFollowByAssetId.Username = fairLaunchFollow.Username
	return fairLaunchFollowByAssetId, nil
}

func CreateOrUpdateFairLaunchFollow(transfer *models.FairLaunchFollow) (err error) {
	var fairLaunchFollow *models.FairLaunchFollow
	fairLaunchFollow, err = CheckFairLaunchFollowIfUpdate(transfer)
	return btldb.UpdateFairLaunchFollow(fairLaunchFollow)
}

func CreateOrUpdateFairLaunchFollows(transfers *[]models.FairLaunchFollow) (err error) {
	var fairLaunchFollows []models.FairLaunchFollow
	var fairLaunchFollow *models.FairLaunchFollow
	for _, transfer := range *transfers {
		fairLaunchFollow, err = CheckFairLaunchFollowIfUpdate(&transfer)
		if err != nil {
			return err
		}
		fairLaunchFollows = append(fairLaunchFollows, *fairLaunchFollow)
	}
	return btldb.UpdateFairLaunchFollows(&fairLaunchFollows)
}

// SetFairLaunchFollow
// @dev: Set
func SetFairLaunchFollow(fairLaunchFollow *models.FairLaunchFollow) error {
	return btldb.CreateFairLaunchFollow(fairLaunchFollow)
}

func SetFairLaunchFollows(fairLaunchFollows *[]models.FairLaunchFollow) error {
	return btldb.CreateFairLaunchFollows(fairLaunchFollows)
}

func GetAllFairLaunchFollowsUpdatedAtDesc() (*[]models.FairLaunchFollow, error) {
	return btldb.ReadAllFairLaunchFollowsUpdatedAtDesc()
}

type FairLaunchFollowSimplified struct {
	UpdatedAt        time.Time `json:"updated_at"`
	FairLaunchInfoId int       `json:"fair_launch_info_id"`
	AssetId          string    `json:"asset_id" gorm:"type:varchar(255)"`
	DeviceId         string    `json:"device_id" gorm:"type:varchar(255)"`
	Username         string    `json:"username" gorm:"type:varchar(255)"`
}

func FairLaunchFollowToFairLaunchFollowSimplified(fairLaunchFollow models.FairLaunchFollow) FairLaunchFollowSimplified {
	return FairLaunchFollowSimplified{
		UpdatedAt:        fairLaunchFollow.UpdatedAt,
		FairLaunchInfoId: fairLaunchFollow.FairLaunchInfoId,
		AssetId:          fairLaunchFollow.AssetId,
		DeviceId:         fairLaunchFollow.DeviceId,
		Username:         fairLaunchFollow.Username,
	}
}

func FairLaunchFollowSliceToFairLaunchFollowSimplifiedSlice(fairLaunchFollows *[]models.FairLaunchFollow) *[]FairLaunchFollowSimplified {
	if fairLaunchFollows == nil {
		return nil
	}
	var fairLaunchFollowSimplified []FairLaunchFollowSimplified
	for _, fairLaunchFollow := range *fairLaunchFollows {
		fairLaunchFollowSimplified = append(fairLaunchFollowSimplified, FairLaunchFollowToFairLaunchFollowSimplified(fairLaunchFollow))
	}
	return &fairLaunchFollowSimplified
}

func GetAllFairLaunchFollowSimplified() (*[]FairLaunchFollowSimplified, error) {
	allFairLaunchFollows, err := GetAllFairLaunchFollowsUpdatedAtDesc()
	if err != nil {
		return nil, err
	}
	allFairLaunchFollowSimplified := FairLaunchFollowSliceToFairLaunchFollowSimplifiedSlice(allFairLaunchFollows)
	return allFairLaunchFollowSimplified, nil
}

func SetFollowFairLaunchInfo(fairLaunchFollow *models.FairLaunchFollow) error {
	if IsFairLaunchFollowed(fairLaunchFollow.UserId, fairLaunchFollow.AssetId) {
		return errors.New("already followed")
	}
	return SetFairLaunchFollow(fairLaunchFollow)
}

func SetUnfollowFairLaunchInfo(userId int, assetId string) error {
	if !IsFairLaunchFollowed(userId, assetId) {
		return errors.New("not followed yet")
	}
	fairLaunchFollow, err := GetFairLaunchFollowByUserIdAndAssetId(userId, assetId)
	if err != nil {
		return err
	}
	return btldb.DeleteFairLaunchFollow(fairLaunchFollow.ID)
}
