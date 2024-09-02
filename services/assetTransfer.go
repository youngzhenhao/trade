package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
	"time"
	"trade/api"
	"trade/middleware"
	"trade/models"
	"trade/services/btldb"
)

func ProcessAssetTransferProcessedSlice(userId int, username string, assetTransferSetRequestSlice *[]models.AssetTransferProcessedSetRequest) (*[]models.AssetTransferProcessedDb, *[]models.AssetTransferProcessedInputDb, *[]models.AssetTransferProcessedOutputDb, error) {
	var assetTransferProcessedSlice []models.AssetTransferProcessedDb
	var assetTransferProcessedInputsSlice []models.AssetTransferProcessedInputDb
	var assetTransferProcessedOutputsSlice []models.AssetTransferProcessedOutputDb
	for _, assetTransferSetRequest := range *assetTransferSetRequestSlice {
		txid := assetTransferSetRequest.Txid
		assetTransferProcessedSlice = append(assetTransferProcessedSlice, models.AssetTransferProcessedDb{
			Txid:               txid,
			AssetID:            assetTransferSetRequest.AssetID,
			TransferTimestamp:  assetTransferSetRequest.TransferTimestamp,
			AnchorTxHash:       assetTransferSetRequest.AnchorTxHash,
			AnchorTxHeightHint: assetTransferSetRequest.AnchorTxHeightHint,
			AnchorTxChainFees:  assetTransferSetRequest.AnchorTxChainFees,
			Inputs:             len(assetTransferSetRequest.Inputs),
			Outputs:            len(assetTransferSetRequest.Outputs),
			DeviceID:           assetTransferSetRequest.DeviceID,
			UserID:             userId,
			Username:           username,
		})
		for index, input := range assetTransferSetRequest.Inputs {
			assetTransferProcessedInputsSlice = append(assetTransferProcessedInputsSlice, models.AssetTransferProcessedInputDb{
				Txid:        txid,
				AssetID:     assetTransferSetRequest.AssetID,
				Index:       index,
				Address:     input.Address,
				Amount:      input.Amount,
				AnchorPoint: input.AnchorPoint,
				ScriptKey:   input.ScriptKey,
				UserID:      userId,
			})
		}
		for index, output := range assetTransferSetRequest.Outputs {
			assetTransferProcessedOutputsSlice = append(assetTransferProcessedOutputsSlice, models.AssetTransferProcessedOutputDb{
				Txid:                   txid,
				AssetID:                assetTransferSetRequest.AssetID,
				Index:                  index,
				Address:                output.Address,
				Amount:                 output.Amount,
				AnchorOutpoint:         output.AnchorOutpoint,
				AnchorValue:            output.AnchorValue,
				AnchorInternalKey:      output.AnchorInternalKey,
				AnchorTaprootAssetRoot: output.AnchorTaprootAssetRoot,
				AnchorMerkleRoot:       output.AnchorMerkleRoot,
				AnchorTapscriptSibling: output.AnchorTapscriptSibling,
				AnchorNumPassiveAssets: output.AnchorNumPassiveAssets,
				ScriptKey:              output.AnchorInternalKey,
				ScriptKeyIsLocal:       output.ScriptKeyIsLocal,
				NewProofBlob:           output.NewProofBlob,
				SplitCommitRootHash:    output.SplitCommitRootHash,
				OutputType:             output.OutputType,
				AssetVersion:           output.AssetVersion,
				UserID:                 userId,
			})
		}
	}
	return &assetTransferProcessedSlice, &assetTransferProcessedInputsSlice, &assetTransferProcessedOutputsSlice, nil
}

func GetAssetTransferProcessedSliceByUserId(userId int) (*[]models.AssetTransferProcessedDb, error) {
	return btldb.ReadAssetTransferProcessedSliceByUserId(userId)
}

func CheckAssetTransferProcessedIfUpdate(assetTransferProcessed *models.AssetTransferProcessedDb) (*models.AssetTransferProcessedDb, error) {
	if assetTransferProcessed == nil {
		return nil, errors.New("nil asset transfer process")
	}
	assetTransferProcessedByTxid, err := btldb.ReadAssetTransferProcessedByTxid(assetTransferProcessed.Txid)
	if err != nil {
		return assetTransferProcessed, nil
	}
	if !IsAssetTransferProcessedChanged(assetTransferProcessed, assetTransferProcessedByTxid) {
		return assetTransferProcessedByTxid, nil
	}
	assetTransferProcessedByTxid.Txid = assetTransferProcessed.Txid
	assetTransferProcessedByTxid.AssetID = assetTransferProcessed.AssetID
	assetTransferProcessedByTxid.TransferTimestamp = assetTransferProcessed.TransferTimestamp
	assetTransferProcessedByTxid.AnchorTxHash = assetTransferProcessed.AnchorTxHash
	assetTransferProcessedByTxid.AnchorTxHeightHint = assetTransferProcessed.AnchorTxHeightHint
	assetTransferProcessedByTxid.AnchorTxChainFees = assetTransferProcessed.AnchorTxChainFees
	assetTransferProcessedByTxid.Inputs = assetTransferProcessed.Inputs
	assetTransferProcessedByTxid.Outputs = assetTransferProcessed.Outputs
	assetTransferProcessedByTxid.DeviceID = assetTransferProcessed.DeviceID
	assetTransferProcessedByTxid.UserID = assetTransferProcessed.UserID
	assetTransferProcessedByTxid.Username = assetTransferProcessed.Username
	return assetTransferProcessedByTxid, nil
}

func CheckAssetTransferProcessedInputIfUpdate(assetTransferProcessedInput *models.AssetTransferProcessedInputDb) (*models.AssetTransferProcessedInputDb, error) {
	if assetTransferProcessedInput == nil {
		return nil, errors.New("nil asset transfer process input")
	}
	assetTransferProcessedInputByTxidAndIndex, err := btldb.ReadAssetTransferProcessedInputByTxidAndIndex(assetTransferProcessedInput.Txid, assetTransferProcessedInput.Index)
	if err != nil {
		return assetTransferProcessedInput, nil
	}
	if !IsAssetTransferProcessedInputChanged(assetTransferProcessedInput, assetTransferProcessedInputByTxidAndIndex) {
		return assetTransferProcessedInputByTxidAndIndex, nil
	}
	assetTransferProcessedInputByTxidAndIndex.Txid = assetTransferProcessedInput.Txid
	assetTransferProcessedInputByTxidAndIndex.AssetID = assetTransferProcessedInput.AssetID
	assetTransferProcessedInputByTxidAndIndex.Index = assetTransferProcessedInput.Index
	assetTransferProcessedInputByTxidAndIndex.Address = assetTransferProcessedInput.Address
	assetTransferProcessedInputByTxidAndIndex.Amount = assetTransferProcessedInput.Amount
	assetTransferProcessedInputByTxidAndIndex.AnchorPoint = assetTransferProcessedInput.AnchorPoint
	assetTransferProcessedInputByTxidAndIndex.ScriptKey = assetTransferProcessedInput.ScriptKey
	assetTransferProcessedInputByTxidAndIndex.UserID = assetTransferProcessedInput.UserID
	return assetTransferProcessedInputByTxidAndIndex, nil
}

func CheckAssetTransferProcessedOutputIfUpdate(assetTransferProcessedOutput *models.AssetTransferProcessedOutputDb) (*models.AssetTransferProcessedOutputDb, error) {
	if assetTransferProcessedOutput == nil {
		return nil, errors.New("nil asset transfer process input")
	}
	assetTransferProcessedOutputByTxidAndIndex, err := btldb.ReadAssetTransferProcessedOutputByTxidAndIndex(assetTransferProcessedOutput.Txid, assetTransferProcessedOutput.Index)
	if err != nil {
		return assetTransferProcessedOutput, nil
	}
	if !IsAssetTransferProcessedOutputChanged(assetTransferProcessedOutput, assetTransferProcessedOutputByTxidAndIndex) {
		return assetTransferProcessedOutputByTxidAndIndex, nil
	}
	assetTransferProcessedOutputByTxidAndIndex.Txid = assetTransferProcessedOutput.Txid
	assetTransferProcessedOutputByTxidAndIndex.AssetID = assetTransferProcessedOutput.AssetID
	assetTransferProcessedOutputByTxidAndIndex.Index = assetTransferProcessedOutput.Index
	assetTransferProcessedOutputByTxidAndIndex.Address = assetTransferProcessedOutput.Address
	assetTransferProcessedOutputByTxidAndIndex.Amount = assetTransferProcessedOutput.Amount
	assetTransferProcessedOutputByTxidAndIndex.AnchorOutpoint = assetTransferProcessedOutput.AnchorOutpoint
	assetTransferProcessedOutputByTxidAndIndex.AnchorValue = assetTransferProcessedOutput.AnchorValue
	assetTransferProcessedOutputByTxidAndIndex.AnchorInternalKey = assetTransferProcessedOutput.AnchorInternalKey
	assetTransferProcessedOutputByTxidAndIndex.AnchorTaprootAssetRoot = assetTransferProcessedOutput.AnchorTaprootAssetRoot
	assetTransferProcessedOutputByTxidAndIndex.AnchorMerkleRoot = assetTransferProcessedOutput.AnchorMerkleRoot
	assetTransferProcessedOutputByTxidAndIndex.AnchorTapscriptSibling = assetTransferProcessedOutput.AnchorTapscriptSibling
	assetTransferProcessedOutputByTxidAndIndex.AnchorNumPassiveAssets = assetTransferProcessedOutput.AnchorNumPassiveAssets
	assetTransferProcessedOutputByTxidAndIndex.ScriptKey = assetTransferProcessedOutput.ScriptKey
	assetTransferProcessedOutputByTxidAndIndex.ScriptKeyIsLocal = assetTransferProcessedOutput.ScriptKeyIsLocal
	assetTransferProcessedOutputByTxidAndIndex.NewProofBlob = assetTransferProcessedOutput.NewProofBlob
	assetTransferProcessedOutputByTxidAndIndex.SplitCommitRootHash = assetTransferProcessedOutput.SplitCommitRootHash
	assetTransferProcessedOutputByTxidAndIndex.OutputType = assetTransferProcessedOutput.OutputType
	assetTransferProcessedOutputByTxidAndIndex.AssetVersion = assetTransferProcessedOutput.AssetVersion
	assetTransferProcessedOutputByTxidAndIndex.UserID = assetTransferProcessedOutput.UserID
	return assetTransferProcessedOutputByTxidAndIndex, nil
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
	if old.DeviceID != assetTransferProcessed.DeviceID {
		return true
	}
	if old.UserID != assetTransferProcessed.UserID {
		return true
	}
	if old.Username != assetTransferProcessed.Username {
		return true
	}
	return false
}

func IsAssetTransferProcessedInputChanged(assetTransferProcessedInput *models.AssetTransferProcessedInputDb, old *models.AssetTransferProcessedInputDb) bool {
	if assetTransferProcessedInput == nil || old == nil {
		return true
	}
	if old.Txid != assetTransferProcessedInput.Txid {
		return true
	}
	if old.Index != assetTransferProcessedInput.Index {
		return true
	}
	if old.AssetID != assetTransferProcessedInput.AssetID {
		return true
	}
	if old.Address != assetTransferProcessedInput.Address {
		return true
	}
	if old.Amount != assetTransferProcessedInput.Amount {
		return true
	}
	if old.AnchorPoint != assetTransferProcessedInput.AnchorPoint {
		return true
	}
	if old.ScriptKey != assetTransferProcessedInput.ScriptKey {
		return true
	}
	if old.UserID != assetTransferProcessedInput.UserID {
		return true
	}
	return false
}

func IsAssetTransferProcessedOutputChanged(assetTransferProcessedOutput *models.AssetTransferProcessedOutputDb, old *models.AssetTransferProcessedOutputDb) bool {
	if assetTransferProcessedOutput == nil || old == nil {
		return true
	}
	if old.Txid != assetTransferProcessedOutput.Txid {
		return true
	}
	if old.AssetID != assetTransferProcessedOutput.AssetID {
		return true
	}
	if old.Index != assetTransferProcessedOutput.Index {
		return true
	}
	if old.Address != assetTransferProcessedOutput.Address {
		return true
	}
	if old.Amount != assetTransferProcessedOutput.Amount {
		return true
	}
	if old.AnchorOutpoint != assetTransferProcessedOutput.AnchorOutpoint {
		return true
	}
	if old.AnchorValue != assetTransferProcessedOutput.AnchorValue {
		return true
	}
	if old.AnchorInternalKey != assetTransferProcessedOutput.AnchorInternalKey {
		return true
	}
	if old.AnchorTaprootAssetRoot != assetTransferProcessedOutput.AnchorTaprootAssetRoot {
		return true
	}
	if old.AnchorMerkleRoot != assetTransferProcessedOutput.AnchorMerkleRoot {
		return true
	}
	if old.AnchorTapscriptSibling != assetTransferProcessedOutput.AnchorTapscriptSibling {
		return true
	}
	if old.AnchorNumPassiveAssets != assetTransferProcessedOutput.AnchorNumPassiveAssets {
		return true
	}
	if old.ScriptKey != assetTransferProcessedOutput.ScriptKey {
		return true
	}
	if old.ScriptKeyIsLocal != assetTransferProcessedOutput.ScriptKeyIsLocal {
		return true
	}
	if old.NewProofBlob != assetTransferProcessedOutput.NewProofBlob {
		return true
	}
	if old.SplitCommitRootHash != assetTransferProcessedOutput.SplitCommitRootHash {
		return true
	}
	if old.OutputType != assetTransferProcessedOutput.OutputType {
		return true
	}
	if old.AssetVersion != assetTransferProcessedOutput.AssetVersion {
		return true
	}
	if old.UserID != assetTransferProcessedOutput.UserID {
		return true
	}
	return false
}

func CreateOrUpdateAssetTransferProcessedSlice(assetTransferProcessedSlice *[]models.AssetTransferProcessedDb) (err error) {
	var assetTransferSlice []models.AssetTransferProcessedDb
	var transfer *models.AssetTransferProcessedDb
	for _, assetTransferProcessed := range *assetTransferProcessedSlice {
		transfer, err = CheckAssetTransferProcessedIfUpdate(&assetTransferProcessed)
		if err != nil {
			return err
		}
		assetTransferSlice = append(assetTransferSlice, *transfer)
	}
	return btldb.UpdateAssetTransferProcessedSlice(&assetTransferSlice)
}

func CreateOrUpdateAssetTransferProcessedInputSlice(assetTransferProcessedInputSlice *[]models.AssetTransferProcessedInputDb) (err error) {
	var assetTransferInputSlice []models.AssetTransferProcessedInputDb
	var input *models.AssetTransferProcessedInputDb
	for _, assetTransferProcessedInput := range *assetTransferProcessedInputSlice {
		input, err = CheckAssetTransferProcessedInputIfUpdate(&assetTransferProcessedInput)
		if err != nil {
			return err
		}
		assetTransferInputSlice = append(assetTransferInputSlice, *input)
	}
	return btldb.UpdateAssetTransferProcessedInputSlice(&assetTransferInputSlice)
}

func CreateOrUpdateAssetTransferProcessedOutputSlice(assetTransferProcessedOutputSlice *[]models.AssetTransferProcessedOutputDb) (err error) {
	var assetTransferOutputSlice []models.AssetTransferProcessedOutputDb
	var input *models.AssetTransferProcessedOutputDb
	for _, assetTransferProcessedOutput := range *assetTransferProcessedOutputSlice {
		input, err = CheckAssetTransferProcessedOutputIfUpdate(&assetTransferProcessedOutput)
		if err != nil {
			return err
		}
		assetTransferOutputSlice = append(assetTransferOutputSlice, *input)
	}
	return btldb.UpdateAssetTransferProcessedOutputSlice(&assetTransferOutputSlice)
}

func GetTxidsByAssetTransfers(transfers *[]models.AssetTransferProcessedDb) []string {
	var txids []string
	for _, transfer := range *transfers {
		txids = append(txids, transfer.Txid)
	}
	return txids
}

func GetAssetTransferTxidsByUserId(userId int) ([]string, error) {
	transfers, err := GetAssetTransferProcessedSliceByUserId(userId)
	if err != nil {
		return nil, err
	}
	return GetTxidsByAssetTransfers(transfers), nil
}

func GetAssetTransferProcessedInputSliceByUserId(userId int) (*[]models.AssetTransferProcessedInputDb, error) {
	return btldb.ReadAssetTransferProcessedInputSliceByUserId(userId)
}

func GetAssetTransferProcessedOutputSliceByUserId(userId int) (*[]models.AssetTransferProcessedOutputDb, error) {
	return btldb.ReadAssetTransferProcessedOutputSliceByUserId(userId)
}

func GetInputsByTxidWithTransfersInputs(transfersInputs *[]models.AssetTransferProcessedInputDb, inputLength int, txid string) ([]models.AssetTransferProcessedInput, error) {
	result := make([]models.AssetTransferProcessedInput, inputLength)
	for _, input := range *transfersInputs {
		if input.Txid == txid && input.Index < inputLength {
			result[input.Index] = models.AssetTransferProcessedInput{
				Address:     input.Address,
				Amount:      input.Amount,
				AnchorPoint: input.AnchorPoint,
				ScriptKey:   input.ScriptKey,
			}
		}
	}
	return result, nil
}

func GetOutputsByTxidWithTransfersOutputs(transfersOutputs *[]models.AssetTransferProcessedOutputDb, outputLength int, txid string) ([]models.AssetTransferProcessedOutput, error) {
	result := make([]models.AssetTransferProcessedOutput, outputLength)
	for _, output := range *transfersOutputs {
		if output.Txid == txid && output.Index < outputLength {
			result[output.Index] = models.AssetTransferProcessedOutput{
				Address:                output.Address,
				Amount:                 output.Amount,
				AnchorOutpoint:         output.AnchorOutpoint,
				AnchorValue:            output.AnchorValue,
				AnchorInternalKey:      output.AnchorInternalKey,
				AnchorTaprootAssetRoot: output.AnchorTaprootAssetRoot,
				AnchorMerkleRoot:       output.AnchorMerkleRoot,
				AnchorTapscriptSibling: output.AnchorTapscriptSibling,
				AnchorNumPassiveAssets: output.AnchorNumPassiveAssets,
				ScriptKey:              output.ScriptKey,
				ScriptKeyIsLocal:       output.ScriptKeyIsLocal,
				NewProofBlob:           output.NewProofBlob,
				SplitCommitRootHash:    output.SplitCommitRootHash,
				OutputType:             output.OutputType,
				AssetVersion:           output.AssetVersion,
			}
		}
	}
	return result, nil
}

func CombineAssetTransfers(transfers *[]models.AssetTransferProcessedDb, transfersInputs *[]models.AssetTransferProcessedInputDb, transfersOutputs *[]models.AssetTransferProcessedOutputDb) (*[]models.AssetTransferProcessedCombined, error) {
	var err error
	var transferCombinedSlice []models.AssetTransferProcessedCombined
	for _, transfer := range *transfers {
		var transferCombined models.AssetTransferProcessedCombined
		inputs := make([]models.AssetTransferProcessedInput, transfer.Inputs)
		inputs, err = GetInputsByTxidWithTransfersInputs(transfersInputs, transfer.Inputs, transfer.Txid)
		if err != nil {
			return nil, err
		}
		outputs := make([]models.AssetTransferProcessedOutput, transfer.Outputs)
		outputs, err = GetOutputsByTxidWithTransfersOutputs(transfersOutputs, transfer.Outputs, transfer.Txid)
		transferCombined = models.AssetTransferProcessedCombined{
			Model:              transfer.Model,
			Txid:               transfer.Txid,
			AssetID:            transfer.AssetID,
			TransferTimestamp:  transfer.TransferTimestamp,
			AnchorTxHash:       transfer.AnchorTxHash,
			AnchorTxHeightHint: transfer.AnchorTxHeightHint,
			AnchorTxChainFees:  transfer.AnchorTxChainFees,
			Inputs:             inputs,
			Outputs:            outputs,
			DeviceID:           transfer.DeviceID,
			UserID:             transfer.UserID,
			Username:           transfer.Username,
			Status:             transfer.Status,
		}
		transferCombinedSlice = append(transferCombinedSlice, transferCombined)
	}
	return &transferCombinedSlice, nil
}

func GetAssetTransferCombinedSliceByUserId(userId int) (*[]models.AssetTransferProcessedCombined, error) {
	var err error
	var transferCombinedSlice *[]models.AssetTransferProcessedCombined
	// @dev: 1.AssetTransferProcessedDb
	transfers, err := GetAssetTransferProcessedSliceByUserId(userId)
	if err != nil {
		return nil, err
	}
	// @dev: 2.Get all inputs by user id
	transfersInputs, err := GetAssetTransferProcessedInputSliceByUserId(userId)
	if err != nil {
		return nil, err
	}
	// @dev: 3.Get all outputs by user id
	transfersOutputs, err := GetAssetTransferProcessedOutputSliceByUserId(userId)
	if err != nil {
		return nil, err
	}
	// @dev: 4.Range and combine data
	transferCombinedSlice, err = CombineAssetTransfers(transfers, transfersInputs, transfersOutputs)
	if err != nil {
		return nil, err
	}
	return transferCombinedSlice, nil
}

func GetAllAssetTransferProcessedSlice() (*[]models.AssetTransferProcessedDb, error) {
	return btldb.ReadAllAssetTransferProcessedSlice()
}

func GetAllAssetTransferProcessedInputSlice() (*[]models.AssetTransferProcessedInputDb, error) {
	return btldb.ReadAllAssetTransferProcessedInputSlice()
}

func GetAllAssetTransferProcessedOutputSlice() (*[]models.AssetTransferProcessedOutputDb, error) {
	return btldb.ReadAllAssetTransferProcessedOutputSlice()
}

// GetAllAssetTransferCombinedSlice
// @Description: get all asset transfer combined slice
func GetAllAssetTransferCombinedSlice() (*[]models.AssetTransferProcessedCombined, error) {
	var err error
	var transferCombinedSlice *[]models.AssetTransferProcessedCombined
	// @dev: 1.Get all transfers
	allAssetTransfers, err := GetAllAssetTransferProcessedSlice()
	if err != nil {
		return nil, err
	}
	// @dev: 2.Get all inputs
	allAssetTransfersInputs, err := GetAllAssetTransferProcessedInputSlice()
	if err != nil {
		return nil, err
	}
	// @dev: 3.Get all outputs
	allAssetTransfersOutputs, err := GetAllAssetTransferProcessedOutputSlice()
	if err != nil {
		return nil, err
	}
	// @dev: 4.Range and combine data
	transferCombinedSlice, err = CombineAssetTransfers(allAssetTransfers, allAssetTransfersInputs, allAssetTransfersOutputs)
	if err != nil {
		return nil, err
	}
	return transferCombinedSlice, nil
}

func GetAssetTransferProcessedSliceByAssetId(assetId string) (*[]models.AssetTransferProcessedDb, error) {
	return btldb.ReadAssetTransferProcessedSliceByAssetId(assetId)
}

func GetAssetTransferProcessedInputSliceByAssetId(assetId string) (*[]models.AssetTransferProcessedInputDb, error) {
	return btldb.ReadAssetTransferProcessedInputSliceByAssetId(assetId)
}

func GetAssetTransferProcessedOutputSliceByAssetId(assetId string) (*[]models.AssetTransferProcessedOutputDb, error) {
	return btldb.ReadAssetTransferProcessedOutputSliceByAssetId(assetId)
}

// @dev: Use this
func GetAssetTransferProcessedSliceByAssetIdLimit(assetId string, limit int) (*[]models.AssetTransferProcessedDb, error) {
	return btldb.ReadAssetTransferProcessedSliceByAssetIdLimit(assetId, limit)
}

// Deprecated
func GetAssetTransferProcessedInputSliceByAssetIdLimit(assetId string, limit int) (*[]models.AssetTransferProcessedInputDb, error) {
	return btldb.ReadAssetTransferProcessedInputSliceByAssetIdLimit(assetId, limit)
}

// Deprecated
func GetAssetTransferProcessedOutputSliceByAssetIdLimit(assetId string, limit int) (*[]models.AssetTransferProcessedOutputDb, error) {
	return btldb.ReadAssetTransferProcessedOutputSliceByAssetIdLimit(assetId, limit)
}

func GetAssetTransferCombinedSliceByAssetId(assetId string) (*[]models.AssetTransferProcessedCombined, error) {
	var err error
	var transferCombinedSlice *[]models.AssetTransferProcessedCombined
	transfers, err := GetAssetTransferProcessedSliceByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	transfersInputs, err := GetAssetTransferProcessedInputSliceByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	transfersOutputs, err := GetAssetTransferProcessedOutputSliceByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	transferCombinedSlice, err = CombineAssetTransfers(transfers, transfersInputs, transfersOutputs)
	if err != nil {
		return nil, err
	}
	return transferCombinedSlice, nil
}

func GetAssetTransferCombinedSliceByAssetIdLimit(assetId string, limit int) (*[]models.AssetTransferProcessedCombined, error) {
	var err error
	var transferCombinedSlice *[]models.AssetTransferProcessedCombined
	// @dev: Use limit only here
	transfers, err := GetAssetTransferProcessedSliceByAssetIdLimit(assetId, limit)
	if err != nil {
		return nil, err
	}
	transfersInputs, err := GetAssetTransferProcessedInputSliceByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	transfersOutputs, err := GetAssetTransferProcessedOutputSliceByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	transferCombinedSlice, err = CombineAssetTransfers(transfers, transfersInputs, transfersOutputs)
	if err != nil {
		return nil, err
	}
	return transferCombinedSlice, nil
}

type AssetIdAndUserAssetTransferAmount struct {
	AssetId                  string                     `json:"asset_id"`
	UserAssetTransferAmounts *[]UserAssetTransferAmount `json:"user_asset_transfer_amounts"`
}

type AssetIdAndUserAssetTransferAmountMap struct {
	AssetId                    string       `json:"asset_id"`
	UserAssetTransferAmountMap *map[int]int `json:"user_asset_transfer_amount_map"`
}

type UserAssetTransferAmount struct {
	UserId              int `json:"user_id"`
	AssetTransferAmount int `json:"asset_transfer_amount"`
}

type AssetTransfer struct {
	AssetId  string `json:"asset_id"`
	Txid     string `json:"txid"`
	Amount   int    `json:"amount"`
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
}

type AssetIdAndTransfer struct {
	AssetId        string           `json:"asset_id"`
	AssetTransfers *[]AssetTransfer `json:"asset_receives"`
}

func GetTotalAmountOfOutputs(assetTransferProcessedOutputs *[]models.AssetTransferProcessedOutput) int {
	var totalAmountOfOutputs int
	for _, output := range *assetTransferProcessedOutputs {
		totalAmountOfOutputs += output.Amount
	}
	return totalAmountOfOutputs
}

func AssetTransferCombinedSliceToAssetTransfers(assetTransferProcessedCombined *[]models.AssetTransferProcessedCombined) *[]AssetTransfer {
	var assetTransfers []AssetTransfer
	for _, assetTransfer := range *assetTransferProcessedCombined {
		assetTransfers = append(assetTransfers, AssetTransfer{
			AssetId:  assetTransfer.AssetID,
			Txid:     assetTransfer.Txid,
			Amount:   GetTotalAmountOfOutputs(&(assetTransfer.Outputs)),
			UserId:   assetTransfer.UserID,
			Username: assetTransfer.Username,
		})
	}
	return &assetTransfers
}

// GetAllAssetTransfers
// @Description: Get all asset transfer
func GetAllAssetTransfers() (*[]AssetTransfer, error) {
	allAssetTransferCombined, err := GetAllAssetTransferCombinedSlice()
	if err != nil {
		return nil, err
	}
	assetTransfers := AssetTransferCombinedSliceToAssetTransfers(allAssetTransferCombined)
	return assetTransfers, nil
}

func AssetTransfersToAssetIdMapAssetTransfers(assetTransfers *[]AssetTransfer) *map[string]*[]AssetTransfer {
	AssetIdMapAssetTransfers := make(map[string]*[]AssetTransfer)
	for _, assetTransfer := range *assetTransfers {
		receives, ok := AssetIdMapAssetTransfers[assetTransfer.AssetId]
		if !ok {
			AssetIdMapAssetTransfers[assetTransfer.AssetId] = &[]AssetTransfer{assetTransfer}
		} else {
			*receives = append(*receives, assetTransfer)
		}
	}
	return &AssetIdMapAssetTransfers
}

func AssetIdMapAssetTransfersToAssetIdAndTransfers(AssetIdMapAssetTransfers *map[string]*[]AssetTransfer) *[]AssetIdAndTransfer {
	var assetIdAndTransfers []AssetIdAndTransfer
	for assetId, assetTransfers := range *AssetIdMapAssetTransfers {
		assetIdAndTransfers = append(assetIdAndTransfers, AssetIdAndTransfer{
			AssetId:        assetId,
			AssetTransfers: assetTransfers,
		})
	}
	return &assetIdAndTransfers
}

func AssetTransfersToAssetIdAndTransfers(assetTransfers *[]AssetTransfer) *[]AssetIdAndTransfer {
	assetIdMapAssetTransfers := AssetTransfersToAssetIdMapAssetTransfers(assetTransfers)
	assetIdAndTransfers := AssetIdMapAssetTransfersToAssetIdAndTransfers(assetIdMapAssetTransfers)
	return assetIdAndTransfers
}

func AssetTransfersToUserMapAssetTransferAmount(assetTransfers *[]AssetTransfer) *map[int]int {
	userMapAssetTransferAmount := make(map[int]int)
	for _, assetTransfer := range *assetTransfers {
		balances, ok := userMapAssetTransferAmount[assetTransfer.UserId]
		if !ok || balances == 0 {
			userMapAssetTransferAmount[assetTransfer.UserId] = assetTransfer.Amount
		} else {
			userMapAssetTransferAmount[assetTransfer.UserId] += assetTransfer.Amount
		}
	}
	return &userMapAssetTransferAmount
}

func UserMapAssetTransferAmountToUserAssetTransferAmount(userMapAssetTransferAmount *map[int]int) *[]UserAssetTransferAmount {
	var userAssetTransferAmount []UserAssetTransferAmount
	for userId, receiveAmount := range *userMapAssetTransferAmount {
		userAssetTransferAmount = append(userAssetTransferAmount, UserAssetTransferAmount{
			UserId:              userId,
			AssetTransferAmount: receiveAmount,
		})
	}
	return &userAssetTransferAmount
}

func AssetTransfersToUserAssetTransferAmount(assetTransfers *[]AssetTransfer) *[]UserAssetTransferAmount {
	userMapAssetTransfers := AssetTransfersToUserMapAssetTransferAmount(assetTransfers)
	userAssetTransferAmount := UserMapAssetTransferAmountToUserAssetTransferAmount(userMapAssetTransfers)
	return userAssetTransferAmount
}

func AssetTransfersToUserAssetTransferAmountMap(assetTransfers *[]AssetTransfer) *map[int]int {
	userMapAssetTransfers := AssetTransfersToUserMapAssetTransferAmount(assetTransfers)
	return userMapAssetTransfers
}

// GetAllAssetIdAndUserAssetTransferAmount
// @Description: asset transfers to user asset transfer amount
func GetAllAssetIdAndUserAssetTransferAmount() (*[]AssetIdAndUserAssetTransferAmount, error) {
	var assetIdAndUserAssetTransferAmount []AssetIdAndUserAssetTransferAmount
	allAssetTransfers, err := GetAllAssetTransfers()
	if err != nil {
		return nil, err
	}
	assetIdAndTransfers := AssetTransfersToAssetIdAndTransfers(allAssetTransfers)
	for _, assetIdAndTransfer := range *assetIdAndTransfers {
		userAssetTransferAmount := AssetTransfersToUserAssetTransferAmount(assetIdAndTransfer.AssetTransfers)
		assetIdAndUserAssetTransferAmount = append(assetIdAndUserAssetTransferAmount, AssetIdAndUserAssetTransferAmount{
			AssetId:                  assetIdAndTransfer.AssetId,
			UserAssetTransferAmounts: userAssetTransferAmount,
		})
	}
	return &assetIdAndUserAssetTransferAmount, nil
}

// @dev: Use map
func GetAllAssetIdAndUserAssetTransferAmountMap() (*[]AssetIdAndUserAssetTransferAmountMap, error) {
	var assetIdAndUserAssetTransferAmount []AssetIdAndUserAssetTransferAmountMap
	allAssetTransfers, err := GetAllAssetTransfers()
	if err != nil {
		return nil, err
	}
	assetIdAndTransfers := AssetTransfersToAssetIdAndTransfers(allAssetTransfers)
	for _, assetIdAndTransfer := range *assetIdAndTransfers {
		userAssetTransferAmountMap := AssetTransfersToUserAssetTransferAmountMap(assetIdAndTransfer.AssetTransfers)
		assetIdAndUserAssetTransferAmount = append(assetIdAndUserAssetTransferAmount, AssetIdAndUserAssetTransferAmountMap{
			AssetId:                    assetIdAndTransfer.AssetId,
			UserAssetTransferAmountMap: userAssetTransferAmountMap,
		})
	}
	return &assetIdAndUserAssetTransferAmount, nil
}

func AssetTransfersToAddressAmountMap(allAssetTransfers *[]models.AssetTransferProcessedCombined) *map[string]*AssetIdAndAmount {
	addressAmountMap := make(map[string]*AssetIdAndAmount)
	for _, assetTransfer := range *allAssetTransfers {

		for _, input := range assetTransfer.Inputs {
			_, ok := addressAmountMap[input.Address]
			if !ok {
				addressAmountMap[input.Address] = &AssetIdAndAmount{
					AssetId: assetTransfer.AssetID,
				}
			}
			if (*(addressAmountMap[input.Address])).AssetId == assetTransfer.AssetID {
				(*(addressAmountMap[input.Address])).Amount -= input.Amount
			}
		}
		for _, output := range assetTransfer.Outputs {
			_, ok := addressAmountMap[output.Address]
			if !ok {
				addressAmountMap[output.Address] = &AssetIdAndAmount{
					AssetId: assetTransfer.AssetID,
				}
			}
			if (*(addressAmountMap[output.Address])).AssetId == assetTransfer.AssetID {
				(*(addressAmountMap[output.Address])).Amount += output.Amount
			}
		}
	}
	return &addressAmountMap
}

type AssetIdAndAmount struct {
	AssetId string `json:"asset_id"`
	Amount  int    `json:"amount"`
}

// AllAssetTransferCombinedToAddressAmountMap
// @Description: all asset transfer combined to address amount map
func AllAssetTransferCombinedToAddressAmountMap() (*map[string]*AssetIdAndAmount, error) {
	allAssetTransfers, err := GetAllAssetTransferCombinedSlice()
	if err != nil {
		return nil, err
	}
	addressAmountMap := AssetTransfersToAddressAmountMap(allAssetTransfers)
	return addressAmountMap, nil
}

func SetAssetTransfer(transfers *[]models.AssetTransferProcessedSetRequest) error {
	userByte := sha256.Sum256([]byte(AdminUploadUserName))
	username := hex.EncodeToString(userByte[:])
	userId, err := NameToId(username)
	if err != nil {
		// @dev: Admin upload user does not exist
		password, _ := hashPassword(username)
		if password == "" {
			password = username
		}
		err = btldb.CreateUser(&models.User{
			Username: username,
			Password: password,
		})
		if err != nil {
			return err
		}
		userId, err = NameToId(username)
		if err != nil {
			return err
		}
	}
	var assetTransferProcessedSlice *[]models.AssetTransferProcessedDb
	var assetTransferProcessedInputsSlice *[]models.AssetTransferProcessedInputDb
	var assetTransferProcessedOutputsSlice *[]models.AssetTransferProcessedOutputDb
	assetTransferProcessedSlice, assetTransferProcessedInputsSlice, assetTransferProcessedOutputsSlice, err = ProcessAssetTransferProcessedSlice(userId, username, transfers)
	if err != nil {
		return err
	}
	err = CreateOrUpdateAssetTransferProcessedSlice(assetTransferProcessedSlice)
	if err != nil {
		return err
	}
	// @dev: Store inputs and outputs in db
	err = CreateOrUpdateAssetTransferProcessedInputSlice(assetTransferProcessedInputsSlice)
	if err != nil {
		return err
	}
	err = CreateOrUpdateAssetTransferProcessedOutputSlice(assetTransferProcessedOutputsSlice)
	if err != nil {
		return err
	}
	return nil
}

// ListAndSetAssetTransfers
// @Description: List and set asset transfers
// @dev: Use config's network
func ListAndSetAssetTransfers(network models.Network, deviceId string) error {
	transfers, err := api.ListTransfersAndGetProcessedResponse(network, deviceId)
	if err != nil {
		return err
	}
	if transfers == nil || len(*transfers) == 0 {
		return nil
	}
	err = SetAssetTransfer(transfers)
	if err != nil {
		return nil
	}
	return nil
}

type AssetTransferProcessedInputSimplified struct {
	Address     string `json:"address" gorm:"type:varchar(255)"`
	Amount      int    `json:"amount"`
	AnchorPoint string `json:"anchor_point" gorm:"type:varchar(255)"`
	ScriptKey   string `json:"script_key" gorm:"type:varchar(255)"`
}

type AssetTransferProcessedOutputSimplified struct {
	Address        string `json:"address" gorm:"type:varchar(255)"`
	Amount         int    `json:"amount"`
	AnchorOutpoint string `json:"anchor_outpoint" gorm:"type:varchar(255)"`
	AnchorValue    int    `json:"anchor_value"`
	ScriptKey      string `json:"script_key" gorm:"type:varchar(255)"`
}

type AssetTransferProcessedCombinedSimplified struct {
	UpdatedAt         time.Time                                `json:"updated_at"`
	Txid              string                                   `json:"txid" gorm:"type:varchar(255)"`
	AssetID           string                                   `json:"asset_id" gorm:"type:varchar(255)"`
	TransferTimestamp int                                      `json:"transfer_timestamp"`
	Inputs            []AssetTransferProcessedInputSimplified  `json:"inputs"`
	Outputs           []AssetTransferProcessedOutputSimplified `json:"outputs"`
	DeviceID          string                                   `json:"device_id" gorm:"type:varchar(255)"`
	UserID            int                                      `json:"user_id"`
	Username          string                                   `json:"username" gorm:"type:varchar(255)"`
}

func AssetTransferProcessedInputToAssetTransferProcessedInputSimplified(input models.AssetTransferProcessedInput) AssetTransferProcessedInputSimplified {
	return AssetTransferProcessedInputSimplified{
		Address:     input.Address,
		Amount:      input.Amount,
		AnchorPoint: input.AnchorPoint,
		ScriptKey:   input.ScriptKey,
	}
}

func AssetTransferProcessedInputSliceToAssetTransferProcessedInputSimplifiedSlice(inputs []models.AssetTransferProcessedInput) []AssetTransferProcessedInputSimplified {
	var assetTransferProcessedInputSimplified []AssetTransferProcessedInputSimplified
	for _, input := range inputs {
		assetTransferProcessedInputSimplified = append(assetTransferProcessedInputSimplified, AssetTransferProcessedInputToAssetTransferProcessedInputSimplified(input))
	}
	return assetTransferProcessedInputSimplified
}

func AssetTransferProcessedOutputToAssetTransferProcessedOutputSimplified(output models.AssetTransferProcessedOutput) AssetTransferProcessedOutputSimplified {
	return AssetTransferProcessedOutputSimplified{
		Address:        output.Address,
		Amount:         output.Amount,
		AnchorOutpoint: output.AnchorOutpoint,
		AnchorValue:    output.AnchorValue,
		ScriptKey:      output.ScriptKey,
	}
}

func AssetTransferProcessedOutputSliceToAssetTransferProcessedOutputSimplifiedSlice(outputs []models.AssetTransferProcessedOutput) []AssetTransferProcessedOutputSimplified {
	var assetTransferProcessedOutputSimplified []AssetTransferProcessedOutputSimplified
	for _, output := range outputs {
		assetTransferProcessedOutputSimplified = append(assetTransferProcessedOutputSimplified, AssetTransferProcessedOutputToAssetTransferProcessedOutputSimplified(output))
	}
	return assetTransferProcessedOutputSimplified
}

func AssetTransferProcessedCombinedToAssetTransferProcessedCombinedSimplified(assetTransferProcessedCombined models.AssetTransferProcessedCombined) AssetTransferProcessedCombinedSimplified {
	return AssetTransferProcessedCombinedSimplified{
		UpdatedAt:         assetTransferProcessedCombined.UpdatedAt,
		Txid:              assetTransferProcessedCombined.Txid,
		AssetID:           assetTransferProcessedCombined.AssetID,
		TransferTimestamp: assetTransferProcessedCombined.TransferTimestamp,
		Inputs:            AssetTransferProcessedInputSliceToAssetTransferProcessedInputSimplifiedSlice(assetTransferProcessedCombined.Inputs),
		Outputs:           AssetTransferProcessedOutputSliceToAssetTransferProcessedOutputSimplifiedSlice(assetTransferProcessedCombined.Outputs),
		DeviceID:          assetTransferProcessedCombined.DeviceID,
		UserID:            assetTransferProcessedCombined.UserID,
		Username:          assetTransferProcessedCombined.Username,
	}
}

func AssetTransferProcessedCombinedSliceToAssetTransferProcessedCombinedSimplifiedSlice(assetTransferProcessedCombinedSlice *[]models.AssetTransferProcessedCombined) *[]AssetTransferProcessedCombinedSimplified {
	if assetTransferProcessedCombinedSlice == nil {
		return nil
	}
	var assetTransferProcessedCombinedSimplified []AssetTransferProcessedCombinedSimplified
	for _, assetTransfer := range *assetTransferProcessedCombinedSlice {
		assetTransferProcessedCombinedSimplified = append(assetTransferProcessedCombinedSimplified, AssetTransferProcessedCombinedToAssetTransferProcessedCombinedSimplified(assetTransfer))
	}
	return &assetTransferProcessedCombinedSimplified
}

// GetAllAssetTransferCombinedSliceSimplified
// @Description: Get all asset transfer combined slice simplified
func GetAllAssetTransferCombinedSliceSimplified() (*[]AssetTransferProcessedCombinedSimplified, error) {
	allAssetTransfer, err := GetAllAssetTransferCombinedSlice()
	if err != nil {
		return nil, err
	}
	return AssetTransferProcessedCombinedSliceToAssetTransferProcessedCombinedSimplifiedSlice(allAssetTransfer), nil
}

type AssetIdAndAssetTransferCombinedSliceSimplified struct {
	AssetId        string                                      `json:"asset_id"`
	AssetName      string                                      `json:"asset_name"`
	AssetTransfers *[]AssetTransferProcessedCombinedSimplified `json:"asset_transfers"`
}

func AssetTransferProcessedCombinedSimplifiedSliceToAssetIdMapAssetTransferProcessedCombinedSimplified(assetTransferProcessedCombinedSimplifiedSlice *[]AssetTransferProcessedCombinedSimplified) *map[string]*[]AssetTransferProcessedCombinedSimplified {
	if assetTransferProcessedCombinedSimplifiedSlice == nil {
		return nil
	}
	assetIdMapAssetTransferProcessedCombinedSimplified := make(map[string]*[]AssetTransferProcessedCombinedSimplified)
	for _, assetTransferProcessedCombinedSimplified := range *assetTransferProcessedCombinedSimplifiedSlice {
		assetTransfers, ok := assetIdMapAssetTransferProcessedCombinedSimplified[assetTransferProcessedCombinedSimplified.AssetID]
		if !ok {
			assetIdMapAssetTransferProcessedCombinedSimplified[assetTransferProcessedCombinedSimplified.AssetID] = &[]AssetTransferProcessedCombinedSimplified{assetTransferProcessedCombinedSimplified}
		} else {
			*assetTransfers = append(*assetTransfers, assetTransferProcessedCombinedSimplified)
		}
	}
	return &assetIdMapAssetTransferProcessedCombinedSimplified
}

func AssetIdMapAssetTransferProcessedCombinedSimplifiedToAssetIdSlice(assetIdMapAssetTransferProcessedCombinedSimplified *map[string]*[]AssetTransferProcessedCombinedSimplified) []string {
	var assetIdSlice []string
	if assetIdMapAssetTransferProcessedCombinedSimplified == nil {
		return assetIdSlice
	}
	for assetId, _ := range *assetIdMapAssetTransferProcessedCombinedSimplified {
		assetIdSlice = append(assetIdSlice, assetId)
	}
	return assetIdSlice
}

// GetAllAssetIdAndAssetTransferCombinedSliceSimplified
// @Description: Get all asset id and asset transfer combined slice simplified
func GetAllAssetIdAndAssetTransferCombinedSliceSimplified() (*[]AssetIdAndAssetTransferCombinedSliceSimplified, error) {
	var allAssetIdAndAssetTransferCombinedSliceSimplified []AssetIdAndAssetTransferCombinedSliceSimplified
	allAssetTransferCombinedSliceSimplified, err := GetAllAssetTransferCombinedSliceSimplified()
	if err != nil {
		return nil, err
	}
	assetIdMapAssetTransferProcessedCombinedSimplified := AssetTransferProcessedCombinedSimplifiedSliceToAssetIdMapAssetTransferProcessedCombinedSimplified(allAssetTransferCombinedSliceSimplified)
	assetIdSlice := AssetIdMapAssetTransferProcessedCombinedSimplifiedToAssetIdSlice(assetIdMapAssetTransferProcessedCombinedSimplified)
	for _, assetId := range assetIdSlice {
		allAssetIdAndAssetTransferCombinedSliceSimplified = append(allAssetIdAndAssetTransferCombinedSliceSimplified, AssetIdAndAssetTransferCombinedSliceSimplified{
			AssetId:        assetId,
			AssetName:      api.GetAssetNameByAssetId(assetId),
			AssetTransfers: (*assetIdMapAssetTransferProcessedCombinedSimplified)[assetId],
		})
	}
	return &allAssetIdAndAssetTransferCombinedSliceSimplified, nil
}

func GetAssetTransferProcessedByTxid(txid string) (*[]models.AssetTransferProcessedDb, error) {
	return btldb.ReadAssetTransferProcessedSliceByTxid(txid)
}

func GetAssetTransferProcessedInputSliceByTxid(txid string) (*[]models.AssetTransferProcessedInputDb, error) {
	return btldb.ReadAssetTransferProcessedInputSliceByTxid(txid)
}

func GetAssetTransferProcessedOutputSliceByTxid(txid string) (*[]models.AssetTransferProcessedOutputDb, error) {
	return btldb.ReadAssetTransferProcessedOutputSliceByTxid(txid)
}

// GetAssetTransferByTxid
// @Description: Get asset transfer by txid
func GetAssetTransferByTxid(txid string) (*models.AssetTransferProcessedCombined, error) {
	assetTransferProcessed, err := GetAssetTransferProcessedByTxid(txid)
	if err != nil {
		return nil, err
	} else if assetTransferProcessed == nil || len(*assetTransferProcessed) == 0 {
		return nil, errors.New("assetTransferProcessed is nil or empty")
	}
	assetTransferProcessedInput, err := GetAssetTransferProcessedInputSliceByTxid(txid)
	if err != nil {
		return nil, err
	} else if assetTransferProcessedInput == nil || len(*assetTransferProcessedInput) == 0 {
		return nil, errors.New("input of assetTransferProcessed is nil or empty")
	}
	assetTransferProcessedOutput, err := GetAssetTransferProcessedOutputSliceByTxid(txid)
	if err != nil {
		return nil, err
	} else if assetTransferProcessedOutput == nil || len(*assetTransferProcessedOutput) == 0 {
		return nil, errors.New("output of assetTransferProcessed is nil or empty")
	}
	var transferCombined models.AssetTransferProcessedCombined
	var transferCombinedSlice *[]models.AssetTransferProcessedCombined
	transferCombinedSlice, err = CombineAssetTransfers(assetTransferProcessed, assetTransferProcessedInput, assetTransferProcessedOutput)
	if err != nil {
		return nil, err
	} else if transferCombinedSlice == nil || len(*transferCombinedSlice) == 0 {
		return nil, errors.New("transferCombinedSlice is nil or empty")
	}
	transferCombined = (*transferCombinedSlice)[0]
	return &transferCombined, nil
}

func GetAssetTransferProcessedOutputSliceWhoseAddressIsNull() (*[]models.AssetTransferProcessedOutputDb, error) {
	return btldb.ReadAssetTransferProcessedOutputSliceWhoseAddressIsNull()
}

func UpdateAssetTransferProcessedOutputSliceWhoseAddressIsNull(network models.Network) error {
	// @dev: Find asset transfers
	assetTransfers, err := GetAssetTransferProcessedOutputSliceWhoseAddressIsNull()
	if err != nil {
		return err
	}
	// @dev: Get outpoints
	var outpoints []string
	for _, assetTransfer := range *assetTransfers {
		outpoints = append(outpoints, assetTransfer.AnchorOutpoint)
	}
	// @dev: Get addresses by outpoints
	outpointMapAddress, err := api.GetAddressesByOutpointSlice(network, outpoints)
	if err != nil {
		return err
	}
	// @dev: Update address
	for i, assetTransfer := range *assetTransfers {
		address, ok := outpointMapAddress[assetTransfer.AnchorOutpoint]
		if !ok {
			continue
		}
		(*assetTransfers)[i].Address = address
	}
	return btldb.UpdateAssetTransferProcessedOutputSlice(assetTransfers)
}

func ExcludeAssetTransferProcessedSetRequestWhoseOutpointAddressIsNull(assetTransferProcessedSetRequests []models.AssetTransferProcessedSetRequest) []models.AssetTransferProcessedSetRequest {
	var assetTransferProcessedSetRequestsProcessed []models.AssetTransferProcessedSetRequest
	for _, assetTransferProcessedSetRequest := range assetTransferProcessedSetRequests {
		if assetTransferProcessedSetRequest.Outputs[0].Address == "" {
			continue
		} else {
			assetTransferProcessedSetRequestsProcessed = append(assetTransferProcessedSetRequestsProcessed, assetTransferProcessedSetRequest)
		}
	}
	return assetTransferProcessedSetRequestsProcessed
}

// DeleteAssetTransferTransactionByTxid
// @Description: Delete asset transfer transaction by txid
func DeleteAssetTransferTransactionByTxid(txid string) error {
	assetTransferProcessed, err := GetAssetTransferProcessedByTxid(txid)
	if err != nil {
		return err
	} else if assetTransferProcessed == nil || len(*assetTransferProcessed) == 0 {
		return errors.New("assetTransferProcessed is nil or empty")
	}
	assetTransferProcessedInput, err := GetAssetTransferProcessedInputSliceByTxid(txid)
	if err != nil {
		return err
	} else if assetTransferProcessedInput == nil || len(*assetTransferProcessedInput) == 0 {
		return errors.New("input of assetTransferProcessed is nil or empty")
	}
	assetTransferProcessedOutput, err := GetAssetTransferProcessedOutputSliceByTxid(txid)
	if err != nil {
		return err
	} else if assetTransferProcessedOutput == nil || len(*assetTransferProcessedOutput) == 0 {
		return errors.New("output of assetTransferProcessed is nil or empty")
	}
	return middleware.DB.Transaction(func(tx *gorm.DB) error {
		err := btldb.DeleteAssetTransferProcessedSlice(assetTransferProcessed)
		if err != nil {
			return err
		}
		err = btldb.DeleteAssetTransferProcessedInputSlice(assetTransferProcessedInput)
		if err != nil {
			return err
		}
		err = btldb.DeleteAssetTransferProcessedOutputSlice(assetTransferProcessedOutput)
		if err != nil {
			return err
		}
		return nil
	})
}
