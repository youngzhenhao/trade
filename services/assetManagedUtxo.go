package services

import (
	"errors"
	"math"
	"time"
	"trade/models"
	"trade/services/btldb"
)

func ProcessAssetManagedUtxoSetRequest(userId int, username string, assetManagedUtxoRequest models.AssetManagedUtxoSetRequest) models.AssetManagedUtxo {
	var assetManagedUtxo models.AssetManagedUtxo
	assetManagedUtxo = models.AssetManagedUtxo{
		Op:                          assetManagedUtxoRequest.Op,
		OutPoint:                    assetManagedUtxoRequest.OutPoint,
		Time:                        assetManagedUtxoRequest.Time,
		AmtSat:                      assetManagedUtxoRequest.AmtSat,
		InternalKey:                 assetManagedUtxoRequest.InternalKey,
		TaprootAssetRoot:            assetManagedUtxoRequest.TaprootAssetRoot,
		MerkleRoot:                  assetManagedUtxoRequest.MerkleRoot,
		Version:                     assetManagedUtxoRequest.Version,
		AssetGenesisPoint:           assetManagedUtxoRequest.AssetGenesisPoint,
		AssetGenesisName:            assetManagedUtxoRequest.AssetGenesisName,
		AssetGenesisMetaHash:        assetManagedUtxoRequest.AssetGenesisMetaHash,
		AssetGenesisAssetID:         assetManagedUtxoRequest.AssetGenesisAssetID,
		AssetGenesisAssetType:       assetManagedUtxoRequest.AssetGenesisAssetType,
		AssetGenesisOutputIndex:     assetManagedUtxoRequest.AssetGenesisOutputIndex,
		AssetGenesisVersion:         assetManagedUtxoRequest.AssetGenesisVersion,
		Amount:                      assetManagedUtxoRequest.Amount,
		LockTime:                    assetManagedUtxoRequest.LockTime,
		RelativeLockTime:            assetManagedUtxoRequest.RelativeLockTime,
		ScriptVersion:               assetManagedUtxoRequest.ScriptVersion,
		ScriptKey:                   assetManagedUtxoRequest.ScriptKey,
		ScriptKeyIsLocal:            assetManagedUtxoRequest.ScriptKeyIsLocal,
		AssetGroupRawGroupKey:       assetManagedUtxoRequest.AssetGroupRawGroupKey,
		AssetGroupTweakedGroupKey:   assetManagedUtxoRequest.AssetGroupTweakedGroupKey,
		AssetGroupAssetWitness:      assetManagedUtxoRequest.AssetGroupAssetWitness,
		ChainAnchorTx:               assetManagedUtxoRequest.ChainAnchorTx,
		ChainAnchorBlockHash:        assetManagedUtxoRequest.ChainAnchorBlockHash,
		ChainAnchorOutpoint:         assetManagedUtxoRequest.ChainAnchorOutpoint,
		ChainAnchorInternalKey:      assetManagedUtxoRequest.ChainAnchorInternalKey,
		ChainAnchorMerkleRoot:       assetManagedUtxoRequest.ChainAnchorMerkleRoot,
		ChainAnchorTapscriptSibling: assetManagedUtxoRequest.ChainAnchorTapscriptSibling,
		ChainAnchorBlockHeight:      assetManagedUtxoRequest.ChainAnchorBlockHeight,
		IsSpent:                     assetManagedUtxoRequest.IsSpent,
		LeaseOwner:                  assetManagedUtxoRequest.LeaseOwner,
		LeaseExpiry:                 assetManagedUtxoRequest.LeaseExpiry,
		IsBurn:                      assetManagedUtxoRequest.IsBurn,
		DeviceId:                    assetManagedUtxoRequest.DeviceId,
		UserId:                      userId,
		Username:                    username,
	}
	return assetManagedUtxo
}

func ProcessAssetManagedUtxoSetRequests(userId int, username string, assetManagedUtxoRequests *[]models.AssetManagedUtxoSetRequest) *[]models.AssetManagedUtxo {
	var assetManagedUtxos []models.AssetManagedUtxo
	for _, assetManagedUtxoRequest := range *assetManagedUtxoRequests {
		assetManagedUtxo := ProcessAssetManagedUtxoSetRequest(userId, username, assetManagedUtxoRequest)
		assetManagedUtxos = append(assetManagedUtxos, assetManagedUtxo)
	}
	return &assetManagedUtxos
}

func GetAssetManagedUtxosByUserId(userId int) (*[]models.AssetManagedUtxo, error) {
	return btldb.ReadAssetManagedUtxosByUserId(userId)
}

func GetAssetManagedUtxosByAssetId(assetId string) (*[]models.AssetManagedUtxo, error) {
	return btldb.ReadAssetManagedUtxosByAssetId(assetId)
}

func GetAssetManagedUtxoByUserIdAndAssetId(userId int, assetId string) (*models.AssetManagedUtxo, error) {
	return btldb.ReadAssetManagedUtxoByUserIdAndAssetId(userId, assetId)
}

func IsAssetManagedUtxoChanged(assetManagedUtxoByTxidAndIndex *models.AssetManagedUtxo, old *models.AssetManagedUtxo) bool {
	if assetManagedUtxoByTxidAndIndex == nil || old == nil {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.Op != old.Op {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.OutPoint != old.OutPoint {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.Time != old.Time {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.AmtSat != old.AmtSat {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.InternalKey != old.InternalKey {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.TaprootAssetRoot != old.TaprootAssetRoot {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.MerkleRoot != old.MerkleRoot {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.Version != old.Version {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.AssetGenesisPoint != old.AssetGenesisPoint {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.AssetGenesisName != old.AssetGenesisName {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.AssetGenesisMetaHash != old.AssetGenesisMetaHash {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.AssetGenesisAssetID != old.AssetGenesisAssetID {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.AssetGenesisAssetType != old.AssetGenesisAssetType {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.AssetGenesisOutputIndex != old.AssetGenesisOutputIndex {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.AssetGenesisVersion != old.AssetGenesisVersion {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.Amount != old.Amount {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.LockTime != old.LockTime {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.RelativeLockTime != old.RelativeLockTime {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.ScriptVersion != old.ScriptVersion {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.ScriptKey != old.ScriptKey {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.ScriptKeyIsLocal != old.ScriptKeyIsLocal {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.AssetGroupRawGroupKey != old.AssetGroupRawGroupKey {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.AssetGroupTweakedGroupKey != old.AssetGroupTweakedGroupKey {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.AssetGroupAssetWitness != old.AssetGroupAssetWitness {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.ChainAnchorTx != old.ChainAnchorTx {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.ChainAnchorBlockHash != old.ChainAnchorBlockHash {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.ChainAnchorOutpoint != old.ChainAnchorOutpoint {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.ChainAnchorInternalKey != old.ChainAnchorInternalKey {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.ChainAnchorMerkleRoot != old.ChainAnchorMerkleRoot {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.ChainAnchorTapscriptSibling != old.ChainAnchorTapscriptSibling {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.ChainAnchorBlockHeight != old.ChainAnchorBlockHeight {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.IsSpent != old.IsSpent {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.LeaseOwner != old.LeaseOwner {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.LeaseExpiry != old.LeaseExpiry {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.IsBurn != old.IsBurn {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.DeviceId != old.DeviceId {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.UserId != old.UserId {
		return true
	}
	if assetManagedUtxoByTxidAndIndex.Username != old.Username {
		return true
	}
	return false
}

func CheckAssetManagedUtxoIfUpdate(userId int, assetManagedUtxo *models.AssetManagedUtxo) (*models.AssetManagedUtxo, error) {
	if assetManagedUtxo == nil {
		return nil, errors.New("nil asset local mint history")
	}
	assetManagedUtxoByUserIdAndAssetId, err := GetAssetManagedUtxoByUserIdAndAssetId(userId, assetManagedUtxo.AssetGenesisAssetID)
	if err != nil {
		return assetManagedUtxo, nil
	}
	if !IsAssetManagedUtxoChanged(assetManagedUtxoByUserIdAndAssetId, assetManagedUtxo) {
		return assetManagedUtxoByUserIdAndAssetId, nil
	}
	assetManagedUtxoByUserIdAndAssetId.Op = assetManagedUtxo.Op
	assetManagedUtxoByUserIdAndAssetId.OutPoint = assetManagedUtxo.OutPoint
	assetManagedUtxoByUserIdAndAssetId.Time = assetManagedUtxo.Time
	assetManagedUtxoByUserIdAndAssetId.AmtSat = assetManagedUtxo.AmtSat
	assetManagedUtxoByUserIdAndAssetId.InternalKey = assetManagedUtxo.InternalKey
	assetManagedUtxoByUserIdAndAssetId.TaprootAssetRoot = assetManagedUtxo.TaprootAssetRoot
	assetManagedUtxoByUserIdAndAssetId.MerkleRoot = assetManagedUtxo.MerkleRoot
	assetManagedUtxoByUserIdAndAssetId.Version = assetManagedUtxo.Version
	assetManagedUtxoByUserIdAndAssetId.AssetGenesisPoint = assetManagedUtxo.AssetGenesisPoint
	assetManagedUtxoByUserIdAndAssetId.AssetGenesisName = assetManagedUtxo.AssetGenesisName
	assetManagedUtxoByUserIdAndAssetId.AssetGenesisMetaHash = assetManagedUtxo.AssetGenesisMetaHash
	assetManagedUtxoByUserIdAndAssetId.AssetGenesisAssetID = assetManagedUtxo.AssetGenesisAssetID
	assetManagedUtxoByUserIdAndAssetId.AssetGenesisAssetType = assetManagedUtxo.AssetGenesisAssetType
	assetManagedUtxoByUserIdAndAssetId.AssetGenesisOutputIndex = assetManagedUtxo.AssetGenesisOutputIndex
	assetManagedUtxoByUserIdAndAssetId.AssetGenesisVersion = assetManagedUtxo.AssetGenesisVersion
	assetManagedUtxoByUserIdAndAssetId.Amount = assetManagedUtxo.Amount
	assetManagedUtxoByUserIdAndAssetId.LockTime = assetManagedUtxo.LockTime
	assetManagedUtxoByUserIdAndAssetId.RelativeLockTime = assetManagedUtxo.RelativeLockTime
	assetManagedUtxoByUserIdAndAssetId.ScriptVersion = assetManagedUtxo.ScriptVersion
	assetManagedUtxoByUserIdAndAssetId.ScriptKey = assetManagedUtxo.ScriptKey
	assetManagedUtxoByUserIdAndAssetId.ScriptKeyIsLocal = assetManagedUtxo.ScriptKeyIsLocal
	assetManagedUtxoByUserIdAndAssetId.AssetGroupRawGroupKey = assetManagedUtxo.AssetGroupRawGroupKey
	assetManagedUtxoByUserIdAndAssetId.AssetGroupTweakedGroupKey = assetManagedUtxo.AssetGroupTweakedGroupKey
	assetManagedUtxoByUserIdAndAssetId.AssetGroupAssetWitness = assetManagedUtxo.AssetGroupAssetWitness
	assetManagedUtxoByUserIdAndAssetId.ChainAnchorTx = assetManagedUtxo.ChainAnchorTx
	assetManagedUtxoByUserIdAndAssetId.ChainAnchorBlockHash = assetManagedUtxo.ChainAnchorBlockHash
	assetManagedUtxoByUserIdAndAssetId.ChainAnchorOutpoint = assetManagedUtxo.ChainAnchorOutpoint
	assetManagedUtxoByUserIdAndAssetId.ChainAnchorInternalKey = assetManagedUtxo.ChainAnchorInternalKey
	assetManagedUtxoByUserIdAndAssetId.ChainAnchorMerkleRoot = assetManagedUtxo.ChainAnchorMerkleRoot
	assetManagedUtxoByUserIdAndAssetId.ChainAnchorTapscriptSibling = assetManagedUtxo.ChainAnchorTapscriptSibling
	assetManagedUtxoByUserIdAndAssetId.ChainAnchorBlockHeight = assetManagedUtxo.ChainAnchorBlockHeight
	assetManagedUtxoByUserIdAndAssetId.IsSpent = assetManagedUtxo.IsSpent
	assetManagedUtxoByUserIdAndAssetId.LeaseOwner = assetManagedUtxo.LeaseOwner
	assetManagedUtxoByUserIdAndAssetId.LeaseExpiry = assetManagedUtxo.LeaseExpiry
	assetManagedUtxoByUserIdAndAssetId.IsBurn = assetManagedUtxo.IsBurn
	assetManagedUtxoByUserIdAndAssetId.DeviceId = assetManagedUtxo.DeviceId
	assetManagedUtxoByUserIdAndAssetId.UserId = assetManagedUtxo.UserId
	assetManagedUtxoByUserIdAndAssetId.Username = assetManagedUtxo.Username
	return assetManagedUtxoByUserIdAndAssetId, nil
}

func CreateOrUpdateAssetManagedUtxo(userId int, transfer *models.AssetManagedUtxo) (err error) {
	var assetManagedUtxo *models.AssetManagedUtxo
	assetManagedUtxo, err = CheckAssetManagedUtxoIfUpdate(userId, transfer)
	return btldb.UpdateAssetManagedUtxo(assetManagedUtxo)
}

// CreateOrUpdateAssetManagedUtxos
// @Description: create or update asset managed utxos
func CreateOrUpdateAssetManagedUtxos(userId int, transfers *[]models.AssetManagedUtxo) (err error) {
	if transfers == nil || len(*transfers) == 0 {
		return nil
	}
	var assetManagedUtxos []models.AssetManagedUtxo
	var assetManagedUtxo *models.AssetManagedUtxo
	for _, transfer := range *transfers {
		assetManagedUtxo, err = CheckAssetManagedUtxoIfUpdate(userId, &transfer)
		if err != nil {
			return err
		}
		assetManagedUtxos = append(assetManagedUtxos, *assetManagedUtxo)
	}
	return btldb.UpdateAssetManagedUtxos(&assetManagedUtxos)
}

func SetAssetManagedUtxo(assetManagedUtxo *models.AssetManagedUtxo) error {
	return btldb.CreateAssetManagedUtxo(assetManagedUtxo)
}

func SetAssetManagedUtxos(assetManagedUtxos *[]models.AssetManagedUtxo) error {
	return btldb.CreateAssetManagedUtxos(assetManagedUtxos)
}

func GetAllAssetManagedUtxosUpdatedAtDesc() (*[]models.AssetManagedUtxo, error) {
	return btldb.ReadAllAssetManagedUtxosUpdatedAtDesc()
}

type AssetManagedUtxoSimplified struct {
	UpdatedAt             time.Time `json:"updated_at"`
	OutPoint              string    `json:"out_point"`
	Time                  int       `json:"time"`
	AmtSat                int       `json:"amt_sat"`
	AssetGenesisPoint     string    `json:"asset_genesis_point"`
	AssetGenesisName      string    `json:"asset_genesis_name"`
	AssetGenesisMetaHash  string    `json:"asset_genesis_meta_hash"`
	AssetGenesisAssetID   string    `json:"asset_genesis_asset_id"`
	AssetGenesisAssetType string    `json:"asset_genesis_asset_type"`
	Amount                int       `json:"amount"`
	LockTime              int       `json:"lock_time"`
	RelativeLockTime      int       `json:"relative_lock_time"`
	ScriptKey             string    `json:"script_key"`
	AssetGroupRawGroupKey string    `json:"asset_group_raw_group_key"`
	ChainAnchorOutpoint   string    `json:"chain_anchor_outpoint"`
	IsSpent               bool      `json:"is_spent"`
	IsBurn                bool      `json:"is_burn"`
	DeviceId              string    `json:"device_id"`
	Username              string    `json:"username"`
}

func AssetManagedUtxoToAssetManagedUtxoSimplified(assetManagedUtxo models.AssetManagedUtxo) AssetManagedUtxoSimplified {
	return AssetManagedUtxoSimplified{
		UpdatedAt:             assetManagedUtxo.UpdatedAt,
		OutPoint:              assetManagedUtxo.OutPoint,
		Time:                  assetManagedUtxo.Time,
		AmtSat:                assetManagedUtxo.AmtSat,
		AssetGenesisPoint:     assetManagedUtxo.AssetGenesisPoint,
		AssetGenesisName:      assetManagedUtxo.AssetGenesisName,
		AssetGenesisMetaHash:  assetManagedUtxo.AssetGenesisMetaHash,
		AssetGenesisAssetID:   assetManagedUtxo.AssetGenesisAssetID,
		AssetGenesisAssetType: assetManagedUtxo.AssetGenesisAssetType,
		Amount:                assetManagedUtxo.Amount,
		LockTime:              assetManagedUtxo.LockTime,
		RelativeLockTime:      assetManagedUtxo.RelativeLockTime,
		ScriptKey:             assetManagedUtxo.ScriptKey,
		AssetGroupRawGroupKey: assetManagedUtxo.AssetGroupRawGroupKey,
		ChainAnchorOutpoint:   assetManagedUtxo.ChainAnchorOutpoint,
		IsSpent:               assetManagedUtxo.IsSpent,
		IsBurn:                assetManagedUtxo.IsBurn,
		DeviceId:              assetManagedUtxo.DeviceId,
		Username:              assetManagedUtxo.Username,
	}
}

func AssetManagedUtxoSliceToAssetManagedUtxoSimplifiedSlice(assetManagedUtxos *[]models.AssetManagedUtxo) *[]AssetManagedUtxoSimplified {
	if assetManagedUtxos == nil {
		return nil
	}
	var assetManagedUtxoSimplified []AssetManagedUtxoSimplified
	for _, assetManagedUtxo := range *assetManagedUtxos {
		assetManagedUtxoSimplified = append(assetManagedUtxoSimplified, AssetManagedUtxoToAssetManagedUtxoSimplified(assetManagedUtxo))
	}
	return &assetManagedUtxoSimplified
}

func GetAllAssetManagedUtxoSimplified() (*[]AssetManagedUtxoSimplified, error) {
	allAssetManagedUtxos, err := GetAllAssetManagedUtxosUpdatedAtDesc()
	if err != nil {
		return nil, err
	}
	allAssetManagedUtxoSimplified := AssetManagedUtxoSliceToAssetManagedUtxoSimplifiedSlice(allAssetManagedUtxos)
	return allAssetManagedUtxoSimplified, nil
}

func RemoveAssetManagedUtxoByIds(assetManagedUtxoIds *[]int) error {
	return btldb.DeleteAssetManagedUtxoByIds(assetManagedUtxoIds)
}

func AssetManagedUtxosToAssetManagedUtxoIds(assetManagedUtxos *[]models.AssetManagedUtxo) *[]int {
	if assetManagedUtxos == nil {
		return nil
	}
	var assetManagedUtxoIds []int
	for _, assetManagedUtxo := range *assetManagedUtxos {
		assetManagedUtxoIds = append(assetManagedUtxoIds, int(assetManagedUtxo.ID))
	}
	return &assetManagedUtxoIds
}

func AssetManagedUtxosToAssetIds(assetManagedUtxos *[]models.AssetManagedUtxo) *[]string {
	if assetManagedUtxos == nil {
		return nil
	}
	var assetIds []string
	for _, assetManagedUtxo := range *assetManagedUtxos {
		assetIds = append(assetIds, assetManagedUtxo.AssetGenesisAssetID)
	}
	return &assetIds
}

func GetAssetManagedUtxosByIds(assetManagedUtxoIds *[]int) (*[]models.AssetManagedUtxo, error) {
	return btldb.ReadAssetManagedUtxosByIds(assetManagedUtxoIds)
}

func ValidateUserIdAndAssetManagedUtxoIds(userId int, assetManagedUtxoIds *[]int) error {
	assetManagedUtxos, err := GetAssetManagedUtxosByIds(assetManagedUtxoIds)
	if err != nil {
		return err
	}
	for _, assetManagedUtxo := range *assetManagedUtxos {
		if assetManagedUtxo.UserId != userId {
			return errors.New("user id does not match")
		}
	}
	return nil
}

type GetAssetManagedUtxoLimitAndOffsetRequest struct {
	AssetId string `json:"asset_id"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}

type GetAssetManagedUtxoPageNumberByPageSizeRequest struct {
	AssetId  string `json:"asset_id"`
	PageSize int    `json:"page_size"`
}

func GetAssetManagedUtxoByAssetIdLimitAndOffset(assetId string, limit int, offset int) (*[]models.AssetManagedUtxo, error) {
	return btldb.ReadAssetManagedUtxosByAssetIdLimitAndOffset(assetId, limit, offset)
}

func GetAssetManagedUtxoLimitAndOffset(assetId string, limit int, offset int) (*[]models.AssetManagedUtxo, error) {
	return GetAssetManagedUtxoByAssetIdLimitAndOffset(assetId, limit, offset)
}

func GetAssetManagedUtxoPageNumberByPageSize(assetId string, pageSize int) (int, error) {
	recordsNum, err := GetAssetManagedUtxoLength(assetId)
	if err != nil {
		return 0, err
	}
	return int(math.Ceil(float64(recordsNum) / float64(pageSize))), nil
}

func GetAllAssetManagedUtxosByAssetId(assetId string) (*[]models.AssetManagedUtxo, error) {
	return btldb.ReadAssetManagedUtxosByAssetId(assetId)
}

func GetAssetManagedUtxoLength(assetId string) (int, error) {
	response, err := GetAllAssetManagedUtxosByAssetId(assetId)
	if err != nil {
		return 0, err
	}
	if response == nil || len(*(response)) == 0 {
		return 0, nil
	}
	return len(*response), nil
}
