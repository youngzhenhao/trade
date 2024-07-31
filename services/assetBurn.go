package services

import (
	"time"
	"trade/models"
	"trade/services/btldb"
)

func GetAssetBurnsByUserId(userId int) (*[]models.AssetBurn, error) {
	return btldb.ReadAssetBurnsByUserId(userId)
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
	return btldb.ReadAssetBurnsByAssetId(assetId)
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

type AssetBurnSimplified struct {
	UpdatedAt time.Time `json:"updated_at"`
	AssetId   string    `json:"asset_id" gorm:"type:varchar(255)"`
	Amount    int       `json:"amount"`
	DeviceId  string    `json:"device_id" gorm:"type:varchar(255)"`
	Username  string    `json:"username" gorm:"type:varchar(255)"`
}

func AssetBurnToAssetBurnSimplified(assetBurn models.AssetBurn) AssetBurnSimplified {
	return AssetBurnSimplified{
		UpdatedAt: assetBurn.UpdatedAt,
		AssetId:   assetBurn.AssetId,
		Amount:    assetBurn.Amount,
		DeviceId:  assetBurn.DeviceId,
		Username:  assetBurn.Username,
	}
}

func AssetBurnSliceToAssetBurnSimplifiedSlice(assetBurns *[]models.AssetBurn) *[]AssetBurnSimplified {
	if assetBurns == nil {
		return nil
	}
	var assetBurnSimplified []AssetBurnSimplified
	for _, assetBurn := range *assetBurns {
		assetBurnSimplified = append(assetBurnSimplified, AssetBurnToAssetBurnSimplified(assetBurn))
	}
	return &assetBurnSimplified
}

func GetAllAssetBurns() (*[]models.AssetBurn, error) {
	return btldb.ReadAllAssetBurns()
}

func GetAllAssetBurnsUpdatedAt() (*[]models.AssetBurn, error) {
	return btldb.ReadAllAssetBurnsUpdatedAt()
}

func GetAllAssetBurnSimplified() (*[]AssetBurnSimplified, error) {
	assetBurns, err := GetAllAssetBurnsUpdatedAt()
	if err != nil {
		return nil, err
	}
	return AssetBurnSliceToAssetBurnSimplifiedSlice(assetBurns), nil
}
