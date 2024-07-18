package services

import (
	"errors"
	"trade/models"
)

func ProcessAssetTransferProcessedSlice(userId int, assetTransferSetRequestSlice *[]models.AssetTransferProcessedSetRequest) (*[]models.AssetTransferProcessedDb, *[]models.AssetTransferProcessedInputDb, *[]models.AssetTransferProcessedOutputDb, error) {
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
	return ReadAssetTransferProcessedSliceByUserId(userId)
}

func CheckAssetTransferProcessedIfUpdate(assetTransferProcessed *models.AssetTransferProcessedDb) (*models.AssetTransferProcessedDb, error) {
	if assetTransferProcessed == nil {
		return nil, errors.New("nil asset transfer process")
	}
	assetTransferProcessedByTxid, err := ReadAssetTransferProcessedByTxid(assetTransferProcessed.Txid)
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
	return assetTransferProcessedByTxid, nil
}

func CheckAssetTransferProcessedInputIfUpdate(assetTransferProcessedInput *models.AssetTransferProcessedInputDb) (*models.AssetTransferProcessedInputDb, error) {
	if assetTransferProcessedInput == nil {
		return nil, errors.New("nil asset transfer process input")
	}
	assetTransferProcessedInputByTxidAndIndex, err := ReadAssetTransferProcessedInputByTxidAndIndex(assetTransferProcessedInput.Txid, assetTransferProcessedInput.Index)
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
	assetTransferProcessedOutputByTxidAndIndex, err := ReadAssetTransferProcessedOutputByTxidAndIndex(assetTransferProcessedOutput.Txid, assetTransferProcessedOutput.Index)
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
	return UpdateAssetTransferProcessedSlice(&assetTransferSlice)
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
	return UpdateAssetTransferProcessedInputSlice(&assetTransferInputSlice)
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
	return UpdateAssetTransferProcessedOutputSlice(&assetTransferOutputSlice)
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
	return ReadAssetTransferProcessedInputSliceByUserId(userId)
}

func GetAssetTransferProcessedOutputSliceByUserId(userId int) (*[]models.AssetTransferProcessedOutputDb, error) {
	return ReadAssetTransferProcessedOutputSliceByUserId(userId)
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
	return ReadAllAssetTransferProcessedSlice()
}

func GetAllAssetTransferProcessedInputSlice() (*[]models.AssetTransferProcessedInputDb, error) {
	return ReadAllAssetTransferProcessedInputSlice()
}

func GetAllAssetTransferProcessedOutputSlice() (*[]models.AssetTransferProcessedOutputDb, error) {
	return ReadAllAssetTransferProcessedOutputSlice()
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
	return ReadAssetTransferProcessedSliceByAssetId(assetId)
}

func GetAssetTransferProcessedInputSliceByAssetId(assetId string) (*[]models.AssetTransferProcessedInputDb, error) {
	return ReadAssetTransferProcessedInputSliceByAssetId(assetId)
}

func GetAssetTransferProcessedOutputSliceByAssetId(assetId string) (*[]models.AssetTransferProcessedOutputDb, error) {
	return ReadAssetTransferProcessedOutputSliceByAssetId(assetId)
}

// @dev: Use this
func GetAssetTransferProcessedSliceByAssetIdLimit(assetId string, limit int) (*[]models.AssetTransferProcessedDb, error) {
	return ReadAssetTransferProcessedSliceByAssetIdLimit(assetId, limit)
}

// Deprecated
func GetAssetTransferProcessedInputSliceByAssetIdLimit(assetId string, limit int) (*[]models.AssetTransferProcessedInputDb, error) {
	return ReadAssetTransferProcessedInputSliceByAssetIdLimit(assetId, limit)
}

// Deprecated
func GetAssetTransferProcessedOutputSliceByAssetIdLimit(assetId string, limit int) (*[]models.AssetTransferProcessedOutputDb, error) {
	return ReadAssetTransferProcessedOutputSliceByAssetIdLimit(assetId, limit)
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
	AssetId string `json:"asset_id"`
	Txid    string `json:"txid"`
	Amount  int    `json:"amount"`
	UserId  int    `json:"user_id"`
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
			AssetId: assetTransfer.AssetID,
			Txid:    assetTransfer.Txid,
			Amount:  GetTotalAmountOfOutputs(&(assetTransfer.Outputs)),
			UserId:  assetTransfer.UserID,
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
