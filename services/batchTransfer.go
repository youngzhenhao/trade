package services

import (
	"errors"
	"sort"
	"time"
	"trade/api"
	"trade/models"
	"trade/services/btldb"
)

func ProcessBatchTransferSetRequest(userId int, username string, batchTransferRequest *models.BatchTransferRequest) *models.BatchTransfer {
	var batchTransfer models.BatchTransfer
	batchTransfer = models.BatchTransfer{
		Encoded:            batchTransferRequest.Encoded,
		AssetID:            batchTransferRequest.AssetID,
		Amount:             batchTransferRequest.Amount,
		ScriptKey:          batchTransferRequest.ScriptKey,
		InternalKey:        batchTransferRequest.InternalKey,
		TaprootOutputKey:   batchTransferRequest.TaprootOutputKey,
		ProofCourierAddr:   batchTransferRequest.ProofCourierAddr,
		Txid:               batchTransferRequest.Txid,
		TxTotalAmount:      batchTransferRequest.TxTotalAmount,
		Index:              batchTransferRequest.Index,
		TransferTimestamp:  batchTransferRequest.TransferTimestamp,
		AnchorTxHash:       batchTransferRequest.AnchorTxHash,
		AnchorTxHeightHint: batchTransferRequest.AnchorTxHeightHint,
		AnchorTxChainFees:  batchTransferRequest.AnchorTxChainFees,
		DeviceID:           batchTransferRequest.DeviceID,
		UserID:             userId,
		Username:           username,
	}
	return &batchTransfer
}

func ProcessBatchTransfersSetRequest(userId int, username string, batchTransfersRequest *[]models.BatchTransferRequest) *[]models.BatchTransfer {
	var batchTransfers []models.BatchTransfer
	for _, batchTransferRequest := range *batchTransfersRequest {
		batchTransfers = append(batchTransfers, models.BatchTransfer{
			Encoded:            batchTransferRequest.Encoded,
			AssetID:            batchTransferRequest.AssetID,
			Amount:             batchTransferRequest.Amount,
			ScriptKey:          batchTransferRequest.ScriptKey,
			InternalKey:        batchTransferRequest.InternalKey,
			TaprootOutputKey:   batchTransferRequest.TaprootOutputKey,
			ProofCourierAddr:   batchTransferRequest.ProofCourierAddr,
			Txid:               batchTransferRequest.Txid,
			TxTotalAmount:      batchTransferRequest.TxTotalAmount,
			Index:              batchTransferRequest.Index,
			TransferTimestamp:  batchTransferRequest.TransferTimestamp,
			AnchorTxHash:       batchTransferRequest.AnchorTxHash,
			AnchorTxHeightHint: batchTransferRequest.AnchorTxHeightHint,
			AnchorTxChainFees:  batchTransferRequest.AnchorTxChainFees,
			DeviceID:           batchTransferRequest.DeviceID,
			UserID:             userId,
			Username:           username,
		})
	}
	return &batchTransfers
}

func GetBatchTransfersByUserId(userId int) (*[]models.BatchTransfer, error) {
	return btldb.ReadBatchTransfersByUserId(userId)
}

func IsBatchTransferChanged(batchTransferByTxidAndIndex *models.BatchTransfer, old *models.BatchTransfer) bool {
	if batchTransferByTxidAndIndex == nil || old == nil {
		return true
	}
	if batchTransferByTxidAndIndex.Encoded != old.Encoded {
		return true
	}
	if batchTransferByTxidAndIndex.AssetID != old.AssetID {
		return true
	}
	if batchTransferByTxidAndIndex.Amount != old.Amount {
		return true
	}
	if batchTransferByTxidAndIndex.ScriptKey != old.ScriptKey {
		return true
	}
	if batchTransferByTxidAndIndex.InternalKey != old.InternalKey {
		return true
	}
	if batchTransferByTxidAndIndex.TaprootOutputKey != old.TaprootOutputKey {
		return true
	}
	if batchTransferByTxidAndIndex.ProofCourierAddr != old.ProofCourierAddr {
		return true
	}
	if batchTransferByTxidAndIndex.Txid != old.Txid {
		return true
	}
	if batchTransferByTxidAndIndex.TxTotalAmount != old.TxTotalAmount {
		return true
	}
	if batchTransferByTxidAndIndex.Index != old.Index {
		return true
	}
	if batchTransferByTxidAndIndex.TransferTimestamp != old.TransferTimestamp {
		return true
	}
	if batchTransferByTxidAndIndex.AnchorTxHash != old.AnchorTxHash {
		return true
	}
	if batchTransferByTxidAndIndex.AnchorTxHeightHint != old.AnchorTxHeightHint {
		return true
	}
	if batchTransferByTxidAndIndex.AnchorTxChainFees != old.AnchorTxChainFees {
		return true
	}
	if batchTransferByTxidAndIndex.DeviceID != old.DeviceID {
		return true
	}
	if batchTransferByTxidAndIndex.UserID != old.UserID {
		return true
	}
	if batchTransferByTxidAndIndex.Username != old.Username {
		return true
	}
	return false
}

func CheckBatchTransferIfUpdate(batchTransfer *models.BatchTransfer) (*models.BatchTransfer, error) {
	if batchTransfer == nil {
		return nil, errors.New("nil batch transfer")
	}
	txid := batchTransfer.Txid
	index := batchTransfer.Index
	batchTransferByTxidAndIndex, err := btldb.ReadBatchTransferByTxidAndIndex(txid, index)
	if err != nil {
		return batchTransfer, nil
	}
	if !IsBatchTransferChanged(batchTransferByTxidAndIndex, batchTransfer) {
		return batchTransferByTxidAndIndex, nil
	}
	batchTransferByTxidAndIndex.Encoded = batchTransfer.Encoded
	batchTransferByTxidAndIndex.AssetID = batchTransfer.AssetID
	batchTransferByTxidAndIndex.Amount = batchTransfer.Amount
	batchTransferByTxidAndIndex.ScriptKey = batchTransfer.ScriptKey
	batchTransferByTxidAndIndex.InternalKey = batchTransfer.InternalKey
	batchTransferByTxidAndIndex.TaprootOutputKey = batchTransfer.TaprootOutputKey
	batchTransferByTxidAndIndex.ProofCourierAddr = batchTransfer.ProofCourierAddr
	batchTransferByTxidAndIndex.Txid = batchTransfer.Txid
	batchTransferByTxidAndIndex.TxTotalAmount = batchTransfer.TxTotalAmount
	batchTransferByTxidAndIndex.Index = batchTransfer.Index
	batchTransferByTxidAndIndex.TransferTimestamp = batchTransfer.TransferTimestamp
	batchTransferByTxidAndIndex.AnchorTxHash = batchTransfer.AnchorTxHash
	batchTransferByTxidAndIndex.AnchorTxHeightHint = batchTransfer.AnchorTxHeightHint
	batchTransferByTxidAndIndex.AnchorTxChainFees = batchTransfer.AnchorTxChainFees
	batchTransferByTxidAndIndex.DeviceID = batchTransfer.DeviceID
	batchTransferByTxidAndIndex.UserID = batchTransfer.UserID
	batchTransferByTxidAndIndex.Username = batchTransfer.Username
	return batchTransferByTxidAndIndex, nil
}

func CreateOrUpdateBatchTransfer(transfer *models.BatchTransfer) (err error) {
	var batchTransfer *models.BatchTransfer
	batchTransfer, err = CheckBatchTransferIfUpdate(transfer)
	return btldb.UpdateBatchTransfer(batchTransfer)
}

func CreateOrUpdateBatchTransfers(transfers *[]models.BatchTransfer) (err error) {
	var batchTransfers []models.BatchTransfer
	var batchTransfer *models.BatchTransfer
	for _, transfer := range *transfers {
		batchTransfer, err = CheckBatchTransferIfUpdate(&transfer)
		if err != nil {
			return err
		}
		batchTransfers = append(batchTransfers, *batchTransfer)
	}
	return btldb.UpdateBatchTransfers(&batchTransfers)
}

func GetAllBatchTransfersUpdatedAtDesc() (*[]models.BatchTransfer, error) {
	return btldb.ReadAllBatchTransfersUpdatedAtDesc()
}

type BatchTransferSimplified struct {
	UpdatedAt         time.Time `json:"updated_at"`
	Encoded           string    `json:"encoded"`
	AssetID           string    `json:"asset_id" gorm:"type:varchar(255)"`
	Amount            int       `json:"amount"`
	ScriptKey         string    `json:"script_key" gorm:"type:varchar(255)"`
	Txid              string    `json:"txid" gorm:"type:varchar(255)"`
	TxTotalAmount     int       `json:"tx_total_amount"`
	Index             int       `json:"index"`
	TransferTimestamp int       `json:"transfer_timestamp"`
	DeviceID          string    `json:"device_id" gorm:"type:varchar(255)"`
	UserID            int       `json:"user_id"`
	Username          string    `json:"username" gorm:"type:varchar(255)"`
}

type AssetIdAndBatchTransferSimplified struct {
	AssetId       string                     `json:"asset_id"`
	AssetName     string                     `json:"asset_name"`
	BatchTransfer *[]BatchTransferSimplified `json:"batch_transfer"`
}

func BatchTransferToBatchTransferSimplified(batchTransfer models.BatchTransfer) BatchTransferSimplified {
	return BatchTransferSimplified{
		UpdatedAt:         batchTransfer.UpdatedAt,
		Encoded:           batchTransfer.Encoded,
		AssetID:           batchTransfer.AssetID,
		Amount:            batchTransfer.Amount,
		ScriptKey:         batchTransfer.ScriptKey,
		Txid:              batchTransfer.Txid,
		TxTotalAmount:     batchTransfer.TxTotalAmount,
		Index:             batchTransfer.Index,
		TransferTimestamp: batchTransfer.TransferTimestamp,
		DeviceID:          batchTransfer.DeviceID,
		UserID:            batchTransfer.UserID,
		Username:          batchTransfer.Username,
	}
}

func BatchTransferSliceToBatchTransferSimplifiedSlice(batchTransfers *[]models.BatchTransfer) *[]BatchTransferSimplified {
	if batchTransfers == nil {
		return nil
	}
	var batchTransferSimplified []BatchTransferSimplified
	for _, batchTransfer := range *batchTransfers {
		batchTransferSimplified = append(batchTransferSimplified, BatchTransferToBatchTransferSimplified(batchTransfer))
	}
	return &batchTransferSimplified
}

func BatchTransfersToAssetIdMapBatchTransfers(batchTransfers *[]models.BatchTransfer) *map[string]*[]models.BatchTransfer {
	assetIdMapBatchTransfer := make(map[string]*[]models.BatchTransfer)
	if batchTransfers == nil {
		return &assetIdMapBatchTransfer
	}
	for _, batchTransfer := range *batchTransfers {
		transfers, ok := assetIdMapBatchTransfer[batchTransfer.AssetID]
		if !ok {
			assetIdMapBatchTransfer[batchTransfer.AssetID] = &[]models.BatchTransfer{batchTransfer}
		} else {
			*transfers = append(*transfers, batchTransfer)
		}
	}
	return &assetIdMapBatchTransfer
}

func AssetIdMapBatchTransfersToAssetIdSliceSort(assetIdMapBatchTransfers *map[string]*[]models.BatchTransfer) []string {
	var assetIdSlice []string
	for assetId, _ := range *assetIdMapBatchTransfers {
		assetIdSlice = append(assetIdSlice, assetId)
	}
	sort.Strings(assetIdSlice)
	return assetIdSlice
}

func BatchTransfersToAssetIdAndBatchTransferSimplifiedSort(batchTransfers *[]models.BatchTransfer) *[]AssetIdAndBatchTransferSimplified {
	var assetIdAndBatchTransfers []AssetIdAndBatchTransferSimplified
	assetIdMapBatchTransfers := BatchTransfersToAssetIdMapBatchTransfers(batchTransfers)
	assetIdSlice := AssetIdMapBatchTransfersToAssetIdSliceSort(assetIdMapBatchTransfers)
	for _, assetId := range assetIdSlice {
		assetIdAndBatchTransfers = append(assetIdAndBatchTransfers, AssetIdAndBatchTransferSimplified{
			AssetId:       assetId,
			AssetName:     api.GetAssetNameByAssetId(assetId),
			BatchTransfer: BatchTransferSliceToBatchTransferSimplifiedSlice((*assetIdMapBatchTransfers)[assetId]),
		})
	}
	return &assetIdAndBatchTransfers
}

func GetAllAssetIdAndBatchTransferSimplified() (*[]AssetIdAndBatchTransferSimplified, error) {
	allBatchTransfers, err := GetAllBatchTransfersUpdatedAtDesc()
	if err != nil {
		return nil, err
	}
	assetIdAndBatchTransfers := BatchTransfersToAssetIdAndBatchTransferSimplifiedSort(allBatchTransfers)
	return assetIdAndBatchTransfers, nil
}
