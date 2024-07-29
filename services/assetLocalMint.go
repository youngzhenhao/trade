package services

import (
	"errors"
	"time"
	"trade/models"
)

func ProcessAssetLocalMintSetRequest(userId int, username string, assetLocalMintRequest models.AssetLocalMintSetRequest) models.AssetLocalMint {
	var assetLocalMint models.AssetLocalMint
	assetLocalMint = models.AssetLocalMint{
		AssetVersion:    assetLocalMintRequest.AssetVersion,
		AssetType:       assetLocalMintRequest.AssetMetaType,
		Name:            assetLocalMintRequest.Name,
		AssetMetaData:   assetLocalMintRequest.AssetMetaData,
		AssetMetaType:   assetLocalMintRequest.AssetMetaType,
		AssetMetaHash:   assetLocalMintRequest.AssetMetaHash,
		Amount:          assetLocalMintRequest.Amount,
		NewGroupedAsset: assetLocalMintRequest.NewGroupedAsset,
		GroupKey:        assetLocalMintRequest.GroupKey,
		GroupAnchor:     assetLocalMintRequest.GroupAnchor,
		GroupedAsset:    assetLocalMintRequest.GroupedAsset,
		BatchKey:        assetLocalMintRequest.BatchKey,
		BatchTxid:       assetLocalMintRequest.BatchTxid,
		AssetId:         assetLocalMintRequest.AssetId,
		DeviceId:        assetLocalMintRequest.DeviceId,
		UserId:          userId,
		Username:        username,
	}
	return assetLocalMint
}

func ProcessAssetLocalMintSetRequests(userId int, username string, assetLocalMintRequests *[]models.AssetLocalMintSetRequest) *[]models.AssetLocalMint {
	var assetLocalMints []models.AssetLocalMint
	for _, assetLocalMintRequest := range *assetLocalMintRequests {
		assetLocalMint := ProcessAssetLocalMintSetRequest(userId, username, assetLocalMintRequest)
		assetLocalMints = append(assetLocalMints, assetLocalMint)
	}
	return &assetLocalMints
}

func GetAssetLocalMintsByUserId(userId int) (*[]models.AssetLocalMint, error) {
	return ReadAssetLocalMintsByUserId(userId)
}

func GetAssetLocalMintByAssetId(assetId string) (*models.AssetLocalMint, error) {
	return ReadAssetLocalMintByAssetId(assetId)
}

func IsAssetLocalMintChanged(assetLocalMintByTxidAndIndex *models.AssetLocalMint, old *models.AssetLocalMint) bool {
	if assetLocalMintByTxidAndIndex == nil || old == nil {
		return true
	}
	if assetLocalMintByTxidAndIndex.AssetVersion != old.AssetVersion {
		return true
	}
	if assetLocalMintByTxidAndIndex.AssetType != old.AssetType {
		return true
	}
	if assetLocalMintByTxidAndIndex.Name != old.Name {
		return true
	}
	if assetLocalMintByTxidAndIndex.AssetMetaData != old.AssetMetaData {
		return true
	}
	if assetLocalMintByTxidAndIndex.AssetMetaType != old.AssetMetaType {
		return true
	}
	if assetLocalMintByTxidAndIndex.AssetMetaHash != old.AssetMetaHash {
		return true
	}
	if assetLocalMintByTxidAndIndex.Amount != old.Amount {
		return true
	}
	if assetLocalMintByTxidAndIndex.NewGroupedAsset != old.NewGroupedAsset {
		return true
	}
	if assetLocalMintByTxidAndIndex.GroupKey != old.GroupKey {
		return true
	}
	if assetLocalMintByTxidAndIndex.GroupAnchor != old.GroupAnchor {
		return true
	}
	if assetLocalMintByTxidAndIndex.GroupedAsset != old.GroupedAsset {
		return true
	}
	if assetLocalMintByTxidAndIndex.BatchKey != old.BatchKey {
		return true
	}
	if assetLocalMintByTxidAndIndex.BatchTxid != old.BatchTxid {
		return true
	}
	if assetLocalMintByTxidAndIndex.AssetId != old.AssetId {
		return true
	}
	if assetLocalMintByTxidAndIndex.DeviceId != old.DeviceId {
		return true
	}
	if assetLocalMintByTxidAndIndex.UserId != old.UserId {
		return true
	}
	if assetLocalMintByTxidAndIndex.Username != old.Username {
		return true
	}
	return false
}

func CheckAssetLocalMintIfUpdate(assetLocalMint *models.AssetLocalMint) (*models.AssetLocalMint, error) {
	if assetLocalMint == nil {
		return nil, errors.New("nil asset local mint")
	}
	assetLocalMintByAssetId, err := GetAssetLocalMintByAssetId(assetLocalMint.AssetId)
	if err != nil {
		return assetLocalMint, nil
	}
	if !IsAssetLocalMintChanged(assetLocalMintByAssetId, assetLocalMint) {
		return assetLocalMintByAssetId, nil
	}
	assetLocalMintByAssetId.AssetVersion = assetLocalMint.AssetVersion
	assetLocalMintByAssetId.AssetType = assetLocalMint.AssetType
	assetLocalMintByAssetId.Name = assetLocalMint.Name
	assetLocalMintByAssetId.AssetMetaData = assetLocalMint.AssetMetaData
	assetLocalMintByAssetId.AssetMetaType = assetLocalMint.AssetMetaType
	assetLocalMintByAssetId.AssetMetaHash = assetLocalMint.AssetMetaHash
	assetLocalMintByAssetId.Amount = assetLocalMint.Amount
	assetLocalMintByAssetId.NewGroupedAsset = assetLocalMint.NewGroupedAsset
	assetLocalMintByAssetId.GroupKey = assetLocalMint.GroupKey
	assetLocalMintByAssetId.GroupAnchor = assetLocalMint.GroupAnchor
	assetLocalMintByAssetId.GroupedAsset = assetLocalMint.GroupedAsset
	assetLocalMintByAssetId.BatchKey = assetLocalMint.BatchKey
	assetLocalMintByAssetId.BatchTxid = assetLocalMint.BatchTxid
	assetLocalMintByAssetId.AssetId = assetLocalMint.AssetId
	assetLocalMintByAssetId.DeviceId = assetLocalMint.DeviceId
	assetLocalMintByAssetId.UserId = assetLocalMint.UserId
	assetLocalMintByAssetId.Username = assetLocalMint.Username
	return assetLocalMintByAssetId, nil
}

func CreateOrUpdateAssetLocalMint(transfer *models.AssetLocalMint) (err error) {
	var assetLocalMint *models.AssetLocalMint
	assetLocalMint, err = CheckAssetLocalMintIfUpdate(transfer)
	return UpdateAssetLocalMint(assetLocalMint)
}

func CreateOrUpdateAssetLocalMints(transfers *[]models.AssetLocalMint) (err error) {
	var assetLocalMints []models.AssetLocalMint
	var assetLocalMint *models.AssetLocalMint
	for _, transfer := range *transfers {
		assetLocalMint, err = CheckAssetLocalMintIfUpdate(&transfer)
		if err != nil {
			return err
		}
		assetLocalMints = append(assetLocalMints, *assetLocalMint)
	}
	return UpdateAssetLocalMints(&assetLocalMints)
}

// SetAssetLocalMint
// @Description: Set asset local mint
func SetAssetLocalMint(assetLocalMint *models.AssetLocalMint) error {
	return CreateAssetLocalMint(assetLocalMint)
}

func SetAssetLocalMints(assetLocalMints *[]models.AssetLocalMint) error {
	return CreateAssetLocalMints(assetLocalMints)
}

func GetAllAssetLocalMintsUpdatedAtDesc() (*[]models.AssetLocalMint, error) {
	return ReadAllAssetLocalMintsUpdatedAtDesc()
}

type AssetLocalMintSimplified struct {
	UpdatedAt       time.Time `json:"updated_at"`
	AssetType       string    `json:"asset_type" gorm:"type:varchar(255)"`
	Name            string    `json:"name" gorm:"type:varchar(255)"`
	AssetMetaHash   string    `json:"asset_meta_hash" gorm:"type:varchar(255)"`
	Amount          int       `json:"amount"`
	NewGroupedAsset bool      `json:"new_grouped_asset"`
	GroupKey        string    `json:"group_key" gorm:"type:varchar(255)"`
	GroupedAsset    bool      `json:"grouped_asset"`
	BatchTxid       string    `json:"batch_txid" gorm:"type:varchar(255)"`
	AssetId         string    `json:"asset_id" gorm:"type:varchar(255)"`
	DeviceId        string    `json:"device_id" gorm:"type:varchar(255)"`
	Username        string    `json:"username" gorm:"type:varchar(255)"`
}

func AssetLocalMintToAssetLocalMintSimplified(assetLocalMint models.AssetLocalMint) AssetLocalMintSimplified {
	return AssetLocalMintSimplified{
		UpdatedAt:       assetLocalMint.UpdatedAt,
		AssetType:       assetLocalMint.AssetType,
		Name:            assetLocalMint.Name,
		AssetMetaHash:   assetLocalMint.AssetMetaHash,
		Amount:          assetLocalMint.Amount,
		NewGroupedAsset: assetLocalMint.NewGroupedAsset,
		GroupKey:        assetLocalMint.GroupKey,
		GroupedAsset:    assetLocalMint.GroupedAsset,
		BatchTxid:       assetLocalMint.BatchTxid,
		AssetId:         assetLocalMint.AssetId,
		DeviceId:        assetLocalMint.DeviceId,
		Username:        assetLocalMint.Username,
	}
}

func AssetLocalMintSliceToAssetLocalMintSimplifiedSlice(assetLocalMints *[]models.AssetLocalMint) *[]AssetLocalMintSimplified {
	if assetLocalMints == nil {
		return nil
	}
	var assetLocalMintSimplified []AssetLocalMintSimplified
	for _, assetLocalMint := range *assetLocalMints {
		assetLocalMintSimplified = append(assetLocalMintSimplified, AssetLocalMintToAssetLocalMintSimplified(assetLocalMint))
	}
	return &assetLocalMintSimplified
}

func GetAllAssetLocalMintSimplified() (*[]AssetLocalMintSimplified, error) {
	allAssetLocalMints, err := GetAllAssetLocalMintsUpdatedAtDesc()
	if err != nil {
		return nil, err
	}
	allAssetLocalMintSimplified := AssetLocalMintSliceToAssetLocalMintSimplifiedSlice(allAssetLocalMints)
	return allAssetLocalMintSimplified, nil
}
