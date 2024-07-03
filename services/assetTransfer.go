package services

import (
	"trade/models"
)

func ProcessAssetTransferProcessedSlice(userId int, assetTransferSetRequestSlice *[]models.AssetTransferProcessedSetRequest) (*[]models.AssetTransferProcessed, error) {
	var assetTransferProcessedSlice []models.AssetTransferProcessed
	for _, assetTransferSetRequest := range *assetTransferSetRequestSlice {
		assetTransferProcessedSlice = append(assetTransferProcessedSlice, models.AssetTransferProcessed{
			Txid:               assetTransferSetRequest.Txid,
			AssetID:            assetTransferSetRequest.AssetID,
			TransferTimestamp:  assetTransferSetRequest.TransferTimestamp,
			AnchorTxHash:       assetTransferSetRequest.AnchorTxHash,
			AnchorTxHeightHint: assetTransferSetRequest.AnchorTxHeightHint,
			AnchorTxChainFees:  assetTransferSetRequest.AnchorTxChainFees,
			Inputs:             assetTransferSetRequest.Inputs,
			Outputs:            assetTransferSetRequest.Outputs,
			UserID:             userId,
		})
	}
	return &assetTransferProcessedSlice, nil
}

func GetAssetTransferProcessedSliceByUserId(userId int) (*[]models.AssetTransferProcessed, error) {
	return ReadAssetTransferProcessedSliceByUserId(userId)
}

func CreateAssetTransferProcessedIfNotExistOrUpdate(assetTransferProcessed *models.AssetTransferProcessed) (err error) {
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

func IsAssetTransferProcessedChanged(assetTransferProcessed *models.AssetTransferProcessed, old *models.AssetTransferProcessed) bool {
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
	if len(old.Inputs) != len(assetTransferProcessed.Inputs) {
		return true
	}
	// @dev: Only check slice length
	if len(old.Outputs) != len(assetTransferProcessed.Outputs) {
		return true
	}
	if old.UserID != assetTransferProcessed.UserID {
		return true
	}
	return false
}

func CreateOrUpdateAssetTransferProcessedSlice(assetTransferProcessedSlice *[]models.AssetTransferProcessed) (err error) {
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
