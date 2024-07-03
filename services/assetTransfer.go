package services

import (
	"trade/models"
)

func ProcessAssetTransferProcessedSlice(userId int, assetTransferSetRequestSlice *[]models.AssetTransferProcessedSetRequest) (*[]models.AssetTransferProcessedDb, error) {
	var assetTransferProcessedSlice []models.AssetTransferProcessedDb
	for _, assetTransferSetRequest := range *assetTransferSetRequestSlice {
		assetTransferProcessedSlice = append(assetTransferProcessedSlice, models.AssetTransferProcessedDb{
			Txid:               assetTransferSetRequest.Txid,
			AssetID:            assetTransferSetRequest.AssetID,
			TransferTimestamp:  assetTransferSetRequest.TransferTimestamp,
			AnchorTxHash:       assetTransferSetRequest.AnchorTxHash,
			AnchorTxHeightHint: assetTransferSetRequest.AnchorTxHeightHint,
			AnchorTxChainFees:  assetTransferSetRequest.AnchorTxChainFees,
			Inputs:             len(assetTransferSetRequest.Inputs),
			Outputs:            len(assetTransferSetRequest.Outputs),
			UserID:             userId,
		})
	}
	return &assetTransferProcessedSlice, nil
}

func GetAssetTransferProcessedSliceByUserId(userId int) (*[]models.AssetTransferProcessedDb, error) {
	return ReadAssetTransferProcessedSliceByUserId(userId)
}

func CreateAssetTransferProcessedIfNotExistOrUpdate(assetTransferProcessed *models.AssetTransferProcessedDb) (err error) {
	assetTransferProcessedByTxid, err := ReadAssetTransferProcessedByTxid(assetTransferProcessed.Txid)
	if err != nil {
		err = CreateAssetTransferProcessed(assetTransferProcessed)
		if err != nil {
			return err
		}
		return nil
	}
	if !IsAssetTransferProcessedChanged(assetTransferProcessed, assetTransferProcessedByTxid) {
		return nil
	}
	assetTransferProcessed.Txid = assetTransferProcessedByTxid.Txid
	assetTransferProcessed.AssetID = assetTransferProcessedByTxid.AssetID
	assetTransferProcessed.TransferTimestamp = assetTransferProcessedByTxid.TransferTimestamp
	assetTransferProcessed.AnchorTxHash = assetTransferProcessedByTxid.AnchorTxHash
	assetTransferProcessed.AnchorTxHeightHint = assetTransferProcessedByTxid.AnchorTxHeightHint
	assetTransferProcessed.AnchorTxChainFees = assetTransferProcessedByTxid.AnchorTxChainFees
	assetTransferProcessed.Inputs = assetTransferProcessedByTxid.Inputs
	assetTransferProcessed.Outputs = assetTransferProcessedByTxid.Outputs
	assetTransferProcessed.UserID = assetTransferProcessedByTxid.UserID
	return UpdateAssetTransferProcessed(assetTransferProcessed)
}

func IsAssetTransferProcessedChanged(assetTransferProcessed *models.AssetTransferProcessedDb, old *models.AssetTransferProcessedDb) bool {
	if assetTransferProcessed == nil || old == nil {
		return true
	}
	if old.Txid != assetTransferProcessed.Txid {
		return true
	}
	if old.AssetID != assetTransferProcessed.AssetID {
		return true
	}
	if old.TransferTimestamp != assetTransferProcessed.TransferTimestamp {
		return true
	}
	if old.AnchorTxHash != assetTransferProcessed.AnchorTxHash {
		return true
	}
	if old.AnchorTxHeightHint != assetTransferProcessed.AnchorTxHeightHint {
		return true
	}
	if old.AnchorTxChainFees != assetTransferProcessed.AnchorTxChainFees {
		return true
	}
	// @dev: Only check slice length
	if (old.Inputs) != (assetTransferProcessed.Inputs) {
		return true
	}
	// @dev: Only check slice length
	if (old.Outputs) != (assetTransferProcessed.Outputs) {
		return true
	}
	if old.UserID != assetTransferProcessed.UserID {
		return true
	}
	return false
}

func CreateOrUpdateAssetTransferProcessedSlice(assetTransferProcessedSlice *[]models.AssetTransferProcessedDb) (err error) {
	for _, assetTransferProcessed := range *assetTransferProcessedSlice {
		err = CreateAssetTransferProcessedIfNotExistOrUpdate(&assetTransferProcessed)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetAssetTransferTxidsByUserId(userId int) ([]string, error) {
	transfers, err := GetAssetTransferProcessedSliceByUserId(userId)
	if err != nil {
		return nil, err
	}
	var txids []string
	for _, transfer := range *transfers {
		txids = append(txids, transfer.Txid)
	}
	return txids, nil
}
