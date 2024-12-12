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
	SuccessErr = SUCCESS.Error()
)

// Err type:Normal
const (
	SUCCESS    ErrCode = 200
	DefaultErr ErrCode = -1
	ReadDbErr  ErrCode = 4001
)

func (e ErrCode) Code() int {
	return int(e)
}

// Err type:Unknown
const (
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
	GetAssetRecommendByUserIdAndAssetIdErr
	GetAllFairLaunchInfosErr
	FairLaunchInfoIdInvalidErr
	GetFairLaunchInfoErr
	FairLaunchMintedInfoIdInvalidErr
	GetFairLaunchMintedInfosByFairLaunchIdErr
	ProcessFairLaunchInfoErr
	SetFairLaunchInfoErr
	IsFairLaunchMintTimeRightErr
	IsTimeRightErr
	IsFairLaunchIssued
	ProcessFairLaunchMintedInfoErr
	SetFairLaunchMintedInfoErr
	GetInventoryCouldBeMintedByFairLaunchInfoIdErr
	UpdateAndCalculateGasFeeRateByMempoolErr
	GetNumberAndAmountOfInventoryCouldBeMintedErr
	GetFairLaunchInfoByAssetIdErr
	InvalidUserIdErr
	SendFairLaunchReservedErr
	UpdateFairLaunchInfoIsReservedSentErr
	GetIssuedFairLaunchInfosErr
	GetOwnFairLaunchInfosByUserIdErr
	GetOwnFairLaunchMintedInfosByUserIdErr
	GetFairLaunchInfoSimplifiedByUserIdIssuedErr
	GetClosedFairLaunchInfoErr
	GetNotStartedFairLaunchInfoErr
	GetAllUsernameAndAssetIdAssetAddrsErr
	FeeRateAtoiErr
	FeeRateInvalidErr
	GetFairLaunchFollowsByUserIdErr
	SetFollowFairLaunchInfoErr
	SetUnfollowFairLaunchInfoErr
	GetAllFairLaunchFollowSimplifiedErr
	GetFollowedFairLaunchInfoErr
	IsFairLaunchInfoIdAndAssetIdValidErr
	FairLaunchInfoAssetIdInvalidErr
	GetAssetLocalMintHistoriesByUserIdErr
	GetAssetLocalMintHistoryByAssetIdErr
	SetAssetLocalMintHistoryErr
	SetAssetLocalMintHistoriesErr
	GetAllAssetLocalMintHistorySimplifiedErr
	GetAssetManagedUtxosByUserIdErr
	GetAssetManagedUtxoByAssetIdErr
	SetAssetManagedUtxosErr
	GetAllAssetManagedUtxoSimplifiedErr
	ValidateUserIdAndAssetManagedUtxoIdsErr
	GetAmountCouldBeMintByMintedNumberErr
	CreateFairLaunchIncomeOfServerPaySendReservedFeeErr
	GetAssetBalanceByUserIdAndAssetIdErr
	GetAssetTransferByTxidErr
	FormFileErr
	DeviceIdIsNullErr
	OsGetPwdErr
	SaveUploadedFileErr
	CreateLogFileUploadErr
	FileSizeTooLargeErr
	GetAccountAssetBalanceExtendsByAssetIdErr
	BackAmountErr
	GetAllLogFilesErr
	GetFileUploadErr
	OsOpenFileErr
	IoCopyFIleErr
	GetAllAccountAssetTransfersByAssetIdErr
	RefundUserFirstMintByUsernameAndAssetIdErr
	GetAssetHolderBalancePageNumberRequestInvalidErr
	GetAssetHolderBalancePageNumberByPageSizeErr
	GetAccountAssetTransfersLimitAndOffsetErr
	GetAccountAssetTransferPageNumberByPageSizeRequestInvalidErr
	GetAccountAssetTransferPageNumberByPageSizeErr
	GetAccountAssetBalancesLimitAndOffsetErr
	GetAccountAssetBalancePageNumberByPageSizeRequestInvalidErr
	GetAccountAssetBalancePageNumberByPageSizeErr
	GetAssetManagedUtxoLimitAndOffsetErr
	GetAssetManagedUtxoPageNumberByPageSizeRequestInvalidErr
	GetAssetManagedUtxoPageNumberByPageSizeErr
	InvalidTweakedGroupKeyErr
	SetAssetGroupErr
	GetAssetGroupErr
	CreateNftTransferErr
	GetNftTransferByAssetIdErr
	PageNumberExceedsTotalNumberErr
	CreateNftPresaleErr
	CreateNftPresalesErr
	GetNftPresalesByAssetIdErr
	GetLaunchedNftPresalesErr
	GetNftPresalesByBuyerUserIdErr
	BuyNftPresaleErr
	GetNftPresaleByGroupKeyErr
	FetchAssetMetaErr
	GetUserDataErr
	GetUserDataYamlErr
	ReSetFailOrCanceledNftPresaleErr
	GetAccountAssetBalanceUserHoldTotalAmountErr
	ProcessNftPresaleBatchGroupLaunchRequestAndCreateErr
	AddWhitelistsByRequestsErr
	GetBatchGroupsErr
	GetNftPresaleByBatchGroupIdErr
	GetBlockchainInfoErr
	CreateOrUpdateAssetListsErr
	GetAssetListsByUserIdNonZeroErr
	IsAssetListRecordExistErr
	GetUserStatsYamlErr
	GetUserStatsErr
	StatsUserInfoToCsvErr
	TooManyQueryParamsErr
	DateFormatErr
	GetSpecifiedDateUserStatsErr
	InvalidQueryParamErr
	GetUserActiveRecordErr
	GetActiveUserCountBetweenErr
	GetDateLoginCountErr
	PageNumberOutOfRangeErr
	NegativeValueErr
	GetDateIpLoginRecordErr
	GetDateIpLoginRecordCountErr
	GetDateIpLoginCountErr
	GetNewUserCountErr
	GetAssetsNameErr
	DateIpLoginRecordToCsvErr
	GetNewUserRecordAllErr
	GetDateIpLoginRecordAllErr
	GetBackRewardsErr
	RedisGetVerifyErr
	RedisSetRandErr
	AssetBalanceBackupErr
	InvalidHashLengthErr
	UpdateAssetBalanceBackupErr
	GetLatestAssetBalanceHistoriesErr
	CreateAssetBalanceHistoriesErr
	GetGroupFirstImageDataErr
	ProcessPoolAddLiquidityBatchRequestErr
	PoolCreateErr
	ProcessPoolRemoveLiquidityBatchRequestErr
	ProcessPoolSwapExactTokenForTokenNoPathBatchRequestErr
	ProcessPoolSwapTokenForExactTokenNoPathBatchRequestErr
	QueryPoolInfoErr
	AtoiErr
	QueryShareRecordsErr
	QueryUserShareRecordsErr
	QueryShareRecordsCountErr
	QueryUserShareRecordsCountErr
	QuerySwapRecordsCountErr
	QueryUserSwapRecordsCountErr
	QuerySwapRecordsErr
	QueryUserSwapRecordsErr
	UsernameEmptyErr
	QueryUserLpAwardBalanceErr
	QueryUserWithdrawAwardRecordsCountErr
	QueryUserWithdrawAwardRecordsErr
	LimitEmptyErr
	OffsetEmptyErr
	LimitLessThanZeroErr
	OffsetLessThanZeroErr
	UsernameNotMatchErr
	CalcAddLiquidityErr
	CalcRemoveLiquidityErr
	CalcSwapExactTokenForTokenNoPathErr
	CalcSwapTokenForExactTokenNoPathErr
	QueryAddLiquidityBatchCountErr
	QueryAddLiquidityBatchErr
	QueryRemoveLiquidityBatchCountErr
	QueryRemoveLiquidityBatchErr
	QuerySwapExactTokenForTokenNoPathBatchCountErr
	QuerySwapExactTokenForTokenNoPathBatchErr
	QuerySwapTokenForExactTokenNoPathBatchCountErr
	QuerySwapTokenForExactTokenNoPathBatchErr
	QueryWithdrawAwardBatchCountErr
	QueryWithdrawAwardBatchErr
)

// Err type:CustodyAccount
const (
	_ ErrCode = iota + 1000
	_
	_
	_

	//CcustodyAccountPayInsideMissionSuccess
	CustodyAccountPayInsideMissionFaild
	CustodyAccountPayInsideMissionPending
)

func (e ErrCode) Error() string {
	switch {
	case errors.Is(e, SUCCESS):
		return ""
	case errors.Is(e, DefaultErr):
		return "error"
	case errors.Is(e, CustodyAccountPayInsideMissionFaild):
		return "custody account pay inside mission faild"
	case errors.Is(e, CustodyAccountPayInsideMissionPending):
		return "custody account pay inside mission pending"
	case errors.Is(e, ReadDbErr):
		return "get server data error"

	default:
		return ""
	}
}

func MakeJsonErrorResult(code ErrCode, errorString string, data any) string {
	jsr := JsonResult{
		Error: errorString,
		Code:  code,
		Data:  data,
	}
	if errors.Is(code, SUCCESS) {
		jsr.Success = true
	} else {
		jsr.Success = false
	}
	jstr, err := json.Marshal(jsr)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return string(jstr)
}
func MakeJsonErrorResultForHttp(code ErrCode, errorString string, data any) *JsonResult {
	jsr := JsonResult{
		Error: errorString,
		Code:  code,
		Data:  data,
	}
	if errors.Is(code, SUCCESS) {
		jsr.Success = true
	} else {
		jsr.Success = false
	}
	return &jsr
}
