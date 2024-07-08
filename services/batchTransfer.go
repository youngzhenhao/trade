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

func IsBatchTransferChanged(addrReceiveEventByAddrEncoded *models.BatchTransfer, old *models.BatchTransfer) bool {
	if addrReceiveEventByAddrEncoded == nil || old == nil {
		return true
	}
	if addrReceiveEventByAddrEncoded.Encoded != old.Encoded {
		return true
	}
	if addrReceiveEventByAddrEncoded.AssetID != old.AssetID {
		return true
	}
	if addrReceiveEventByAddrEncoded.Amount != old.Amount {
		return true
	}
	if addrReceiveEventByAddrEncoded.ScriptKey != old.ScriptKey {
		return true
	}
	if addrReceiveEventByAddrEncoded.InternalKey != old.InternalKey {
		return true
	}
	if addrReceiveEventByAddrEncoded.TaprootOutputKey != old.TaprootOutputKey {
		return true
	}
	if addrReceiveEventByAddrEncoded.ProofCourierAddr != old.ProofCourierAddr {
		return true
	}
	if addrReceiveEventByAddrEncoded.Txid != old.Txid {
		return true
	}
	if addrReceiveEventByAddrEncoded.Index != old.Index {
		return true
	}
	if addrReceiveEventByAddrEncoded.TransferTimestamp != old.TransferTimestamp {
		return true
	}
	if addrReceiveEventByAddrEncoded.AnchorTxHash != old.AnchorTxHash {
		return true
	}
	if addrReceiveEventByAddrEncoded.AnchorTxHeightHint != old.AnchorTxHeightHint {
		return true
	}
	if addrReceiveEventByAddrEncoded.AnchorTxChainFees != old.AnchorTxChainFees {
		return true
	}
	if addrReceiveEventByAddrEncoded.DeviceID != old.DeviceID {
		return true
	}
	if addrReceiveEventByAddrEncoded.UserID != old.UserID {
		return true
	}
	return false
}

func CheckBatchTransferIfUpdate(addrReceiveEvent *models.BatchTransfer) (*models.BatchTransfer, error) {
	if addrReceiveEvent == nil {
		return nil, errors.New("nil addr receive event")
	}
	txid := addrReceiveEvent.Txid
	index := addrReceiveEvent.Index
	addrReceiveEventByAddrEncoded, err := ReadBatchTransferByTxidAndIndex(txid, index)
	if err != nil {
		return addrReceiveEvent, nil
	}
	if !IsBatchTransferChanged(addrReceiveEventByAddrEncoded, addrReceiveEvent) {
		return addrReceiveEventByAddrEncoded, nil
	}
	addrReceiveEventByAddrEncoded.Encoded = addrReceiveEvent.Encoded
	addrReceiveEventByAddrEncoded.AssetID = addrReceiveEvent.AssetID
	addrReceiveEventByAddrEncoded.Amount = addrReceiveEvent.Amount
	addrReceiveEventByAddrEncoded.ScriptKey = addrReceiveEvent.ScriptKey
	addrReceiveEventByAddrEncoded.InternalKey = addrReceiveEvent.InternalKey
	addrReceiveEventByAddrEncoded.TaprootOutputKey = addrReceiveEvent.TaprootOutputKey
	addrReceiveEventByAddrEncoded.ProofCourierAddr = addrReceiveEvent.ProofCourierAddr
	addrReceiveEventByAddrEncoded.Txid = addrReceiveEvent.Txid
	addrReceiveEventByAddrEncoded.Index = addrReceiveEvent.Index
	addrReceiveEventByAddrEncoded.TransferTimestamp = addrReceiveEvent.TransferTimestamp
	addrReceiveEventByAddrEncoded.AnchorTxHash = addrReceiveEvent.AnchorTxHash
	addrReceiveEventByAddrEncoded.AnchorTxHeightHint = addrReceiveEvent.AnchorTxHeightHint
	addrReceiveEventByAddrEncoded.AnchorTxChainFees = addrReceiveEvent.AnchorTxChainFees
	addrReceiveEventByAddrEncoded.DeviceID = addrReceiveEvent.DeviceID
	addrReceiveEventByAddrEncoded.UserID = addrReceiveEvent.UserID
	return addrReceiveEventByAddrEncoded, nil
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
