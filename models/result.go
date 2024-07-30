package models

import (
	"encoding/json"
	"errors"
)

type JsonResult struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    any     `json:"data"`
}

type ErrCode int

var (
	SuccessErr = errors.New("").Error()
)

const (
	SUCCESS     ErrCode = 200
	DefaultErr  ErrCode = -1
	NameToIdErr ErrCode = iota + 500
	IdAtoiErr
	ShouldBindJsonErr
	SyncAssetIssuanceErr
	GetAssetInfoErr
	IsIdoParticipateTimeRightErr
	IsNotRightTime
	IdoIsNotPublished
	GetAllIdoPublishInfosErr
	GetOwnIdoPublishInfoErr
	GetOwnIdoParticipateInfoErr
	GetIdoParticipateInfoErr
	GetIdoParticipateInfosByAssetIdErr
	GetIdoPublishedInfosErr
	ProcessIdoPublishInfoErr
	ProcessIdoParticipateInfoErr
	SetIdoPublishInfoErr
	SetIdoParticipateInfoErr
	GetBtcBalanceByUsernameErr
	CreateOrUpdateBtcBalanceErr
	ProcessAssetTransferErr
	CreateAssetTransferErr
	GetAssetTransfersByUserIdErr
	GetAddressByOutpointErr
	GetAddressesByOutpointSliceErr
	DecodeRawTransactionSliceErr
	DecodeRawTransactionErr
	GetRawTransactionsByTxidsErr
	GenerateBlocksErr
	FaucetTransferBtcErr
	CreateAssetTransferProcessedErr
	GetAssetTransferProcessedSliceByUserIdErr
	GetAssetTransferCombinedSliceByUserIdErr
	CreateOrUpdateAssetTransferProcessedInputSliceErr
	CreateOrUpdateAssetTransferProcessedOutputSliceErr
	GetAddrReceiveEventsByUserIdErr
	CreateAddrReceiveEventsErr
	GetAddrReceiveEventsProcessedOriginByUserIdErr
	CreateOrUpdateBatchTransferErr
	GetBatchTransfersByUserIdErr
	CreateOrUpdateBatchTransfersErr
	GetAssetAddrsByUserIdErr
	CreateOrUpdateAssetAddrErr
	GetAssetLocksByUserIdErr
	CreateOrUpdateAssetLockErr
	GetAssetBalancesByUserIdErr
	CreateOrUpdateAssetBalanceErr
	CreateOrUpdateAssetBalancesErr
	GetAssetTransferCombinedSliceByAssetIdErr
	GetAssetAddrsByScriptKeyErr
	GetAssetBalancesByUserIdNonZeroErr
	GetAssetHolderNumberAssetBalanceErr
	GetAssetIdAndBalancesByAssetIdErr
	GetTimeByOutpointErr
	GetTimesByOutpointSliceErr
	ValidateAndGetProofFilePathErr
	IsLimitAndOffsetValidErr
	GetAssetBalanceByAssetIdNonZeroLengthErr
	GetAllUsernameAssetBalancesErr
	GetAssetAddrsByEncodedErr
	GetAssetBurnsByUserIdErr
	CreateAssetBurnErr
	UpdateUsernameByUserIdAllErr
	GetAssetBurnTotalErr
	GetAllUsernameAssetBalanceSimplifiedErr
	GetAllAssetAddrsErr
	GetAllAssetAddrSimplifiedErr
	GetAllAssetIdAndBalanceSimplifiedErr
	GetAllAssetIdAndBatchTransfersErr
	GetAllAddrReceiveSimplifiedErr
	GetAllAssetIdAndAddrReceiveSimplifiedErr
	GetAllAssetTransferCombinedSliceErr
	GetAllAssetTransferSimplifiedErr
	GetAllAssetIdAndAssetTransferCombinedSliceSimplifiedErr
	GetAssetLocalMintsByUserIdErr
	GetAssetLocalMintByAssetIdErr
	SetAssetLocalMintErr
	SetAssetLocalMintsErr
	GetAllAssetLocalMintSimplifiedErr
	UpdateUserIpByClientIpErr
	GetAllUserSimplifiedErr
	GetAllAssetBurnSimplifiedErr
	GetAssetRecommendsByUserIdErr
	GetAssetRecommendByAssetIdErr
	SetAssetRecommendErr
	GetAllAssetRecommendSimplifiedErr
)

func MakeJsonErrorResult(code ErrCode, errorString string, data any) string {
	jsonResult := JsonResult{
		Error: errorString,
		Code:  code,
		Data:  data,
	}
	if code == SUCCESS {
		jsonResult.Success = true
	} else {
		jsonResult.Success = false
	}
	jsonStr, err := json.Marshal(jsonResult)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return string(jsonStr)
}
