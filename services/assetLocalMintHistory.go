package services

import (
	"errors"
	"time"
	"trade/models"
	"trade/services/btldb"
)

// TODO:

func ProcessAssetLocalMintHistorySetRequest(userId int, username string, assetLocalMintHistoryRequest models.AssetLocalMintHistorySetRequest) models.AssetLocalMintHistory {
	var assetLocalMintHistory models.AssetLocalMintHistory
	assetLocalMintHistory = models.AssetLocalMintHistory{
		// TODO:
		AssetId:  "",
		DeviceId: "",
		UserId:   0,
		Username: "",
		Status:   0,
	}
	return assetLocalMintHistory
}

func ProcessAssetLocalMintHistorySetRequests(userId int, username string, assetLocalMintHistoryRequests *[]models.AssetLocalMintHistorySetRequest) *[]models.AssetLocalMintHistory {
	var assetLocalMintHistories []models.AssetLocalMintHistory
	for _, assetLocalMintHistoryRequest := range *assetLocalMintHistoryRequests {
		assetLocalMintHistory := ProcessAssetLocalMintHistorySetRequest(userId, username, assetLocalMintHistoryRequest)
		assetLocalMintHistories = append(assetLocalMintHistories, assetLocalMintHistory)
	}
	return &assetLocalMintHistories
}

func GetAssetLocalMintHistoriesByUserId(userId int) (*[]models.AssetLocalMintHistory, error) {
	return btldb.ReadAssetLocalMintHistoriesByUserId(userId)
}

func GetAssetLocalMintHistoryByAssetId(assetId string) (*models.AssetLocalMintHistory, error) {
	return btldb.ReadAssetLocalMintHistoryByAssetId(assetId)
}

func IsAssetLocalMintHistoryChanged(assetLocalMintHistoryByTxidAndIndex *models.AssetLocalMintHistory, old *models.AssetLocalMintHistory) bool {
	if assetLocalMintHistoryByTxidAndIndex == nil || old == nil {
		return true
	}
	// TODO:

	if assetLocalMintHistoryByTxidAndIndex.AssetId != old.AssetId {
		return true
	}
	if assetLocalMintHistoryByTxidAndIndex.DeviceId != old.DeviceId {
		return true
	}
	if assetLocalMintHistoryByTxidAndIndex.UserId != old.UserId {
		return true
	}
	if assetLocalMintHistoryByTxidAndIndex.Username != old.Username {
		return true
	}
	return false
}

func CheckAssetLocalMintHistoryIfUpdate(assetLocalMintHistory *models.AssetLocalMintHistory) (*models.AssetLocalMintHistory, error) {
	if assetLocalMintHistory == nil {
		return nil, errors.New("nil asset local mint history")
	}
	assetLocalMintHistoryByAssetId, err := GetAssetLocalMintHistoryByAssetId(assetLocalMintHistory.AssetId)
	if err != nil {
		return assetLocalMintHistory, nil
	}
	if !IsAssetLocalMintHistoryChanged(assetLocalMintHistoryByAssetId, assetLocalMintHistory) {
		return assetLocalMintHistoryByAssetId, nil
	}
	// TODO:
	assetLocalMintHistoryByAssetId.AssetId = assetLocalMintHistory.AssetId
	assetLocalMintHistoryByAssetId.DeviceId = assetLocalMintHistory.DeviceId
	assetLocalMintHistoryByAssetId.UserId = assetLocalMintHistory.UserId
	assetLocalMintHistoryByAssetId.Username = assetLocalMintHistory.Username
	return assetLocalMintHistoryByAssetId, nil
}

func CreateOrUpdateAssetLocalMintHistory(transfer *models.AssetLocalMintHistory) (err error) {
	var assetLocalMintHistory *models.AssetLocalMintHistory
	assetLocalMintHistory, err = CheckAssetLocalMintHistoryIfUpdate(transfer)
	return btldb.UpdateAssetLocalMintHistory(assetLocalMintHistory)
}

func CreateOrUpdateAssetLocalMintHistories(transfers *[]models.AssetLocalMintHistory) (err error) {
	var assetLocalMintHistories []models.AssetLocalMintHistory
	var assetLocalMintHistory *models.AssetLocalMintHistory
	for _, transfer := range *transfers {
		assetLocalMintHistory, err = CheckAssetLocalMintHistoryIfUpdate(&transfer)
		if err != nil {
			return err
		}
		assetLocalMintHistories = append(assetLocalMintHistories, *assetLocalMintHistory)
	}
	return btldb.UpdateAssetLocalMintHistories(&assetLocalMintHistories)
}

// SetAssetLocalMintHistory
// @Description: Set asset local mint history
func SetAssetLocalMintHistory(assetLocalMintHistory *models.AssetLocalMintHistory) error {
	return btldb.CreateAssetLocalMintHistory(assetLocalMintHistory)
}

func SetAssetLocalMintHistories(assetLocalMintHistories *[]models.AssetLocalMintHistory) error {
	return btldb.CreateAssetLocalMintHistories(assetLocalMintHistories)
}

func GetAllAssetLocalMintHistoriesUpdatedAtDesc() (*[]models.AssetLocalMintHistory, error) {
	return btldb.ReadAllAssetLocalMintHistoriesUpdatedAtDesc()
}

type AssetLocalMintHistorySimplified struct {
	UpdatedAt time.Time `json:"updated_at"`
	AssetType string    `json:"asset_type" gorm:"type:varchar(255)"`
	Name      string    `json:"name" gorm:"type:varchar(255)"`
	// TODO:

	AssetId  string `json:"asset_id" gorm:"type:varchar(255)"`
	DeviceId string `json:"device_id" gorm:"type:varchar(255)"`
	Username string `json:"username" gorm:"type:varchar(255)"`
}

func AssetLocalMintHistoryToAssetLocalMintHistorySimplified(assetLocalMintHistory models.AssetLocalMintHistory) AssetLocalMintHistorySimplified {
	return AssetLocalMintHistorySimplified{
		UpdatedAt: assetLocalMintHistory.UpdatedAt,
		// TODO:

		AssetId:  assetLocalMintHistory.AssetId,
		DeviceId: assetLocalMintHistory.DeviceId,
		Username: assetLocalMintHistory.Username,
	}
}

func AssetLocalMintHistorySliceToAssetLocalMintHistorySimplifiedSlice(assetLocalMintHistories *[]models.AssetLocalMintHistory) *[]AssetLocalMintHistorySimplified {
	if assetLocalMintHistories == nil {
		return nil
	}
	var assetLocalMintHistorySimplified []AssetLocalMintHistorySimplified
	for _, assetLocalMintHistory := range *assetLocalMintHistories {
		assetLocalMintHistorySimplified = append(assetLocalMintHistorySimplified, AssetLocalMintHistoryToAssetLocalMintHistorySimplified(assetLocalMintHistory))
	}
	return &assetLocalMintHistorySimplified
}

func GetAllAssetLocalMintHistorySimplified() (*[]AssetLocalMintHistorySimplified, error) {
	allAssetLocalMintHistories, err := GetAllAssetLocalMintHistoriesUpdatedAtDesc()
	if err != nil {
		return nil, err
	}
	allAssetLocalMintHistorySimplified := AssetLocalMintHistorySliceToAssetLocalMintHistorySimplifiedSlice(allAssetLocalMintHistories)
	return allAssetLocalMintHistorySimplified, nil
}
