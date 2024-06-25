package services

import (
	"gorm.io/gorm"
	"trade/api"
	"trade/models"
	"trade/utils"
)

func ProcessAssetTransfer(userId int, assetTransferSetRequest *models.AssetTransferSetRequest) (*models.AssetTransfer, error) {
	assetInfo, err := api.GetAssetInfo(assetTransferSetRequest.AssetID)
	errorAppendInfo := utils.ErrorAppendInfo(err)
	if err != nil {
		return nil, errorAppendInfo(utils.ToLowerWords("GetAssetInfo"))
	}
	assetTransfer := models.AssetTransfer{
		Model:             gorm.Model{},
		AssetID:           assetTransferSetRequest.AssetID,
		AssetName:         assetInfo.Name,
		AssetType:         assetInfo.AssetType,
		AssetAddressFrom:  assetTransferSetRequest.AssetAddressFrom,
		AssetAddressTo:    assetTransferSetRequest.AssetAddressTo,
		Amount:            assetTransferSetRequest.Amount,
		TransferType:      assetTransferSetRequest.TransferType,
		Inputs:            assetTransferSetRequest.Inputs,
		Outputs:           assetTransferSetRequest.Outputs,
		UserID:            userId,
		TransactionID:     assetTransferSetRequest.TransactionID,
		TransferTimestamp: assetTransferSetRequest.TransferTimestamp,
		AnchorTxChainFees: assetTransferSetRequest.AnchorTxChainFees,
	}
	return &assetTransfer, nil
}

func GetAssetTransfersByUserId(userId int) (*[]models.AssetTransfer, error) {
	return ReadAssetTransfersByUserId(userId)
}
