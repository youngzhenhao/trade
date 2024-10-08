package services

import (
	"trade/models"
	"trade/services/btldb"
)

func ProcessNftTransfer(userId int, username string, nftTransferSetRequest *models.NftTransferSetRequest) *models.NftTransfer {
	return &models.NftTransfer{
		Txid:     nftTransferSetRequest.Txid,
		AssetId:  nftTransferSetRequest.AssetId,
		Time:     nftTransferSetRequest.Time,
		FromAddr: nftTransferSetRequest.FromAddr,
		ToAddr:   nftTransferSetRequest.ToAddr,
		FromInfo: nftTransferSetRequest.FromInfo,
		ToInfo:   nftTransferSetRequest.ToInfo,
		DeviceId: nftTransferSetRequest.DeviceId,
		UserId:   userId,
		Username: username,
	}
}

func CreateNftTransfer(nftTransfer *models.NftTransfer) error {
	return btldb.CreateNftTransfer(nftTransfer)
}

func ReadNftTransfersByAssetId(assetId string) (*[]models.NftTransfer, error) {
	return btldb.ReadNftTransfersByAssetId(assetId)
}

func GetNftTransferByAssetId(assetId string) (*[]models.NftTransfer, error) {
	return ReadNftTransfersByAssetId(assetId)
}
