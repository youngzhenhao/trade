package services

import (
	"errors"
	"trade/models"
)

func ProcessBatchTransferSetRequest(userId int, batchTransferRequest *models.BatchTransferRequest) *models.BatchTransfer {
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
	}
	return &batchTransfer
}

func ProcessBatchTransfersSetRequest(userId int, batchTransfersRequest *[]models.BatchTransferRequest) *[]models.BatchTransfer {
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
		})
	}
	return &batchTransfers
}

func GetBatchTransfersByUserId(userId int) (*[]models.BatchTransfer, error) {
	return ReadBatchTransfersByUserId(userId)
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
	return false
}

func CheckBatchTransferIfUpdate(addrReceiveEvent *models.BatchTransfer) (*models.BatchTransfer, error) {
	if addrReceiveEvent == nil {
		return nil, errors.New("nil batch transfer")
	}
	txid := addrReceiveEvent.Txid
	index := addrReceiveEvent.Index
	batchTransferByTxidAndIndex, err := ReadBatchTransferByTxidAndIndex(txid, index)
	if err != nil {
		return addrReceiveEvent, nil
	}
	if !IsBatchTransferChanged(batchTransferByTxidAndIndex, addrReceiveEvent) {
		return batchTransferByTxidAndIndex, nil
	}
	batchTransferByTxidAndIndex.Encoded = addrReceiveEvent.Encoded
	batchTransferByTxidAndIndex.AssetID = addrReceiveEvent.AssetID
	batchTransferByTxidAndIndex.Amount = addrReceiveEvent.Amount
	batchTransferByTxidAndIndex.ScriptKey = addrReceiveEvent.ScriptKey
	batchTransferByTxidAndIndex.InternalKey = addrReceiveEvent.InternalKey
	batchTransferByTxidAndIndex.TaprootOutputKey = addrReceiveEvent.TaprootOutputKey
	batchTransferByTxidAndIndex.ProofCourierAddr = addrReceiveEvent.ProofCourierAddr
	batchTransferByTxidAndIndex.Txid = addrReceiveEvent.Txid
	batchTransferByTxidAndIndex.TxTotalAmount = addrReceiveEvent.TxTotalAmount
	batchTransferByTxidAndIndex.Index = addrReceiveEvent.Index
	batchTransferByTxidAndIndex.TransferTimestamp = addrReceiveEvent.TransferTimestamp
	batchTransferByTxidAndIndex.AnchorTxHash = addrReceiveEvent.AnchorTxHash
	batchTransferByTxidAndIndex.AnchorTxHeightHint = addrReceiveEvent.AnchorTxHeightHint
	batchTransferByTxidAndIndex.AnchorTxChainFees = addrReceiveEvent.AnchorTxChainFees
	batchTransferByTxidAndIndex.DeviceID = addrReceiveEvent.DeviceID
	batchTransferByTxidAndIndex.UserID = addrReceiveEvent.UserID
	return batchTransferByTxidAndIndex, nil
}

func CreateOrUpdateBatchTransfer(transfer *models.BatchTransfer) (err error) {
	var batchTransfer *models.BatchTransfer
	batchTransfer, err = CheckBatchTransferIfUpdate(transfer)
	return UpdateBatchTransfer(batchTransfer)
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
	return UpdateBatchTransfers(&batchTransfers)
}
