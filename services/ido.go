package services

import (
	"encoding/hex"
	"strconv"
	"time"
	"trade/api"
	"trade/middleware"
	"trade/models"
	"trade/services/btldb"
	"trade/services/custodyAccount"
	"trade/utils"
)

func GetOwnIdoPublishInfosByUserId(UserId int) (*[]models.IdoPublishInfo, error) {
	return btldb.ReadIdoPublishInfosByUserId(UserId)
}

func GetOwnIdoParticipateInfosByUserId(UserId int) (*[]models.IdoParticipateInfo, error) {
	return btldb.ReadIdoParticipateInfosByUserId(UserId)
}

func GetIdoPublishInfo(id int) (*models.IdoPublishInfo, error) {
	return btldb.ReadIdoPublishInfo(uint(id))
}

func GetIdoPublishInfosByAssetId(assetId string) (*[]models.IdoPublishInfo, error) {
	return btldb.ReadIdoPublishInfosByAssetId(assetId)
}

func GetIdoParticipateInfo(id int) (*models.IdoParticipateInfo, error) {
	return btldb.ReadIdoParticipateInfo(uint(id))
}

func GetIdoPublishedInfos() (*[]models.IdoPublishInfo, error) {
	var idoPublishInfos []models.IdoPublishInfo
	err := middleware.DB.Where("state >= ?", models.IdoParticipateStateSent).Order("set_time").Find(&idoPublishInfos).Error
	errorAppendInfo := utils.ErrorAppendInfo(err)
	if err != nil {
		return nil, errorAppendInfo(utils.ToLowerWords("IdoPublishInfoFind"))
	}
	return &idoPublishInfos, nil
}

func ProcessIdoPublishedInfos(idoPublishInfos *[]models.IdoPublishInfo) *[]models.IdoPublishInfo {
	var result []models.IdoPublishInfo
	for _, idoPublishInfo := range *idoPublishInfos {
		if IsIdoPublishInfoTimeValid(&idoPublishInfo) {
			result = append(result, idoPublishInfo)
		}
	}
	return &result
}

func IsIdoPublishInfoTimeValid(idoPublishInfo *models.IdoPublishInfo) bool {
	return IsDuringParticipateTime(idoPublishInfo.StartTime, idoPublishInfo.EndTime)
}

func IsDuringParticipateTime(start int, end int) bool {
	now := int(time.Now().Unix())
	return now >= start && now < end
}

func GetAllIdoPublishInfos() (*[]models.IdoPublishInfo, error) {
	return btldb.ReadAllIdoPublishInfos()
}

func ProcessIdoPublishInfo(userId int, assetID string, totalAmount int, minimumQuantity int, unitPrice int, startTime int, endTime int, feeRate int) (*models.IdoPublishInfo, error) {
	err := ValidateStartAndEndTime(startTime, endTime)
	errorAppendInfo := utils.ErrorAppendInfo(err)
	if err != nil {
		return nil, errorAppendInfo("validate start and end time")
	}
	var idoPublishInfo models.IdoPublishInfo
	// TODO: need to calculate fee
	setGasFee := GetIdoPublishTransactionGasFee(feeRate)
	if !custodyAccount.IsAccountBalanceEnoughByUserId(uint(userId), uint64(setGasFee)) {
		return nil, errorAppendInfo("account balance not enough to pay publish gas fee")
	}
	assetInfo, err := api.GetAssetInfo(assetID)
	if err != nil {
		return nil, errorAppendInfo(utils.ToLowerWords("GetAssetInfo"))
	}
	idoPublishInfo = models.IdoPublishInfo{
		AssetID:         assetID,
		AssetName:       assetInfo.Name,
		AssetType:       assetInfo.AssetType,
		TotalAmount:     totalAmount,
		MinimumQuantity: minimumQuantity,
		UnitPrice:       unitPrice,
		StartTime:       startTime,
		EndTime:         endTime,
		FeeRate:         feeRate,
		GasFee:          setGasFee,
		SetTime:         utils.GetTimestamp(),
		UserID:          userId,
		State:           models.IdoPublishStateNoPay,
	}
	return &idoPublishInfo, nil
}

func SetIdoPublishInfo(idoPublishInfo *models.IdoPublishInfo) error {
	return btldb.CreateIdoPublishInfo(idoPublishInfo)
}

func IsIdoParticipateTimeValid(idoPublishInfo *models.IdoPublishInfo) bool {
	return IsDuringParticipateTime(idoPublishInfo.StartTime, idoPublishInfo.EndTime)
}

func IsIdoParticipateTimeRight(idoParticipateInfoId int) (bool, error) {
	idoPublishInfo, err := GetIdoPublishInfo(idoParticipateInfoId)
	errorAppendInfo := utils.ErrorAppendInfo(err)
	if err != nil {
		return false, errorAppendInfo("GetIdoPublishInfo")
	}
	return IsIdoParticipateTimeValid(idoPublishInfo), nil
}

func IsIdoPublished(idoPublishInfoId int) bool {
	state, err := GetIdoPublishInfoState(idoPublishInfoId)
	if err != nil {
		return false
	}
	return state >= models.IdoPublishStatePublished
}

func GetIdoPublishInfoState(idoPublishInfoId int) (idoPublishInfoState models.IdoPublishState, err error) {
	var idoPublishInfo *models.IdoPublishInfo
	idoPublishInfo, err = GetIdoPublishInfo(idoPublishInfoId)
	errorAppendInfo := utils.ErrorAppendInfo(err)
	if err != nil {
		return 0, errorAppendInfo(utils.ToLowerWords("GetIdoPublishInfo"))
	}
	return idoPublishInfo.State, nil
}

func SetIdoParticipateInfo(idoParticipateInfo *models.IdoParticipateInfo) error {
	return btldb.CreateIdoParticipateInfo(idoParticipateInfo)
}

func ProcessIdoParticipateInfo(userId int, idoPublishInfoID int, boughtAmount int, feeRate int, encodedAddr string) (*models.IdoParticipateInfo, error) {
	var idoParticipateInfo models.IdoParticipateInfo
	isTimeRight, err := IsIdoParticipateTimeRight(idoPublishInfoID)
	errorAppendInfo := utils.ErrorAppendInfo(err)
	if err != nil {
		return nil, errorAppendInfo(utils.ToLowerWords("IsIdoParticipateTimeRight"))
	}
	if !isTimeRight {
		return nil, errorAppendInfo("it is not right participate time now")
	}
	decodedAddrInfo, err := api.GetDecodedAddrInfo(encodedAddr)
	if err != nil {
		return nil, errorAppendInfo(utils.ToLowerWords("GetDecodedAddrInfo"))
	}
	var idoPublishInfo *models.IdoPublishInfo
	idoPublishInfo, err = GetIdoPublishInfo(idoPublishInfoID)
	if err != nil {
		return nil, errorAppendInfo(utils.ToLowerWords("GetIdoPublishInfo"))
	}
	decodedAddrAssetId := hex.EncodeToString(decodedAddrInfo.AssetId)
	if idoPublishInfo.AssetID != decodedAddrAssetId {
		return nil, errorAppendInfo("decoded addr asset id is not equal ido publish info's asset id")
	}
	isValid, err := IsIdoParticipateAmountValid(idoPublishInfoID, boughtAmount)
	if err != nil || !isValid {
		return nil, errorAppendInfo("is participate amount valid")
	}
	// TODO: need to calculate fee
	participateGasFee := GetIdoParticipateTransactionGasFee(feeRate)
	if !custodyAccount.IsAccountBalanceEnoughByUserId(uint(userId), uint64(participateGasFee)) {
		return nil, errorAppendInfo("account balance not enough to pay minted gas fee")
	}
	idoParticipateInfo = models.IdoParticipateInfo{
		IdoPublishInfoID: idoPublishInfoID,
		AssetID:          idoPublishInfo.AssetID,
		AssetName:        idoPublishInfo.AssetName,
		AssetType:        idoPublishInfo.AssetType,
		BoughtAmount:     boughtAmount,
		FeeRate:          feeRate,
		GasFee:           participateGasFee,
		SetTime:          utils.GetTimestamp(),
		UserID:           userId,
		EncodedAddr:      encodedAddr,
		ScriptKey:        hex.EncodeToString(decodedAddrInfo.ScriptKey),
		InternalKey:      hex.EncodeToString(decodedAddrInfo.InternalKey),
		TaprootOutputKey: hex.EncodeToString(decodedAddrInfo.TaprootOutputKey),
		ProofCourierAddr: decodedAddrInfo.ProofCourierAddr,
		State:            models.IdoParticipateStateNoPay,
	}
	return &idoParticipateInfo, nil
}

func IsIdoParticipateAmountValid(idoPublishInfoID int, boughtAmount int) (isValid bool, err error) {
	errorAppendInfo := utils.ErrorAppendInfo(err)
	idoPublishInfo, err := GetIdoPublishInfo(idoPublishInfoID)
	if err != nil {
		return false, errorAppendInfo(utils.ToLowerWords("GetIdoPublishInfo"))
	}
	participateAmount, err := GetParticipateAmountByIdoPublishInfoId(idoPublishInfoID)
	if err != nil {
		return false, errorAppendInfo(utils.ToLowerWords("GetParticipateAmountByIdoPublishInfoId"))
	}
	if !(participateAmount > idoPublishInfo.MinimumQuantity) {
		info := "boughtAmount is not bigger than MinimumQuantity"
		return false, errorAppendInfo(utils.ToLowerWords(info))
	}

	if boughtAmount+participateAmount > idoPublishInfo.TotalAmount {
		info := "Reach max participate amount, available: " + strconv.Itoa(idoPublishInfo.TotalAmount-participateAmount)
		return false, errorAppendInfo(info)
	}
	return true, nil
}

func GetParticipateAmountByIdoPublishInfoId(idoPublishInfoId int) (int, error) {
	idoPublishInfo, err := GetIdoPublishInfo(idoPublishInfoId)
	errorAppendInfo := utils.ErrorAppendInfo(err)
	if err != nil {
		return 0, errorAppendInfo(utils.ToLowerWords("GetIdoPublishInfo"))
	}
	return idoPublishInfo.ParticipateAmount, nil
}

//TODO: need to process ido publish and participate infos
