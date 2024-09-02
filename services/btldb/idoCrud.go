package btldb

import (
	"trade/middleware"
	"trade/models"
)

// IdoPublishInfo

func CreateIdoPublishInfo(idoPublishInfo *models.IdoPublishInfo) error {
	return middleware.DB.Create(idoPublishInfo).Error
}

func CreateIdoPublishInfos(idoPublishInfos *[]models.IdoPublishInfo) error {
	return middleware.DB.Create(idoPublishInfos).Error
}

func ReadAllIdoPublishInfos() (*[]models.IdoPublishInfo, error) {
	var idoPublishInfos []models.IdoPublishInfo
	err := middleware.DB.Find(&idoPublishInfos).Error
	return &idoPublishInfos, err
}

func ReadIdoPublishInfo(id uint) (*models.IdoPublishInfo, error) {
	var idoPublishInfo models.IdoPublishInfo
	err := middleware.DB.First(&idoPublishInfo, id).Error
	return &idoPublishInfo, err
}

func ReadIdoPublishInfosByAssetId(assetId string) (*[]models.IdoPublishInfo, error) {
	var idoPublishInfos []models.IdoPublishInfo
	err := middleware.DB.Where("asset_id = ?", assetId).Find(&idoPublishInfos).Error
	return &idoPublishInfos, err
}

func ReadIdoPublishInfosByUserId(userId int) (*[]models.IdoPublishInfo, error) {
	var idoPublishInfos []models.IdoPublishInfo
	err := middleware.DB.Where("user_id = ?", userId).Find(&idoPublishInfos).Error
	return &idoPublishInfos, err
}

func UpdateIdoPublishInfo(idoPublishInfo *models.IdoPublishInfo) error {
	return middleware.DB.Save(idoPublishInfo).Error
}

func UpdateIdoPublishInfos(idoPublishInfos *[]models.IdoPublishInfo) error {
	return middleware.DB.Save(idoPublishInfos).Error
}

func DeleteIdoPublishInfo(id uint) error {
	var idoPublishInfo models.IdoPublishInfo
	return middleware.DB.Delete(&idoPublishInfo, id).Error
}

// IdoParticipateInfo

func CreateIdoParticipateInfo(idoParticipateInfo *models.IdoParticipateInfo) error {
	return middleware.DB.Create(idoParticipateInfo).Error
}

func CreateIdoParticipateInfos(idoParticipateInfos *[]models.IdoParticipateInfo) error {
	return middleware.DB.Create(idoParticipateInfos).Error
}

func ReadIdoParticipateInfo(id uint) (*models.IdoParticipateInfo, error) {
	var idoParticipateInfo models.IdoParticipateInfo
	err := middleware.DB.First(&idoParticipateInfo, id).Error
	return &idoParticipateInfo, err
}

func ReadIdoParticipateInfosByAssetId(assetId string) (*[]models.IdoParticipateInfo, error) {
	var idoParticipateInfos []models.IdoParticipateInfo
	err := middleware.DB.Where("asset_id = ?", assetId).Find(&idoParticipateInfos).Error
	return &idoParticipateInfos, err
}

func ReadIdoParticipateInfosByUserId(userId int) (*[]models.IdoParticipateInfo, error) {
	var idoParticipateInfos []models.IdoParticipateInfo
	err := middleware.DB.Where("user_id = ?", userId).Find(&idoParticipateInfos).Error
	return &idoParticipateInfos, err
}

func UpdateIdoParticipateInfo(idoParticipateInfo *models.IdoParticipateInfo) error {
	return middleware.DB.Save(idoParticipateInfo).Error
}

func UpdateIdoParticipateInfos(idoParticipateInfos *[]models.IdoParticipateInfo) error {
	return middleware.DB.Save(idoParticipateInfos).Error
}

func DeleteIdoParticipateInfo(id uint) error {
	var idoParticipateInfo models.IdoParticipateInfo
	return middleware.DB.Delete(&idoParticipateInfo, id).Error
}

// IdoParticipateUserInfo

func CreateParticipateIdoUserInfo(participateIdoUserInfo *models.IdoParticipateUserInfo) error {
	return middleware.DB.Create(participateIdoUserInfo).Error
}

func CreateParticipateIdoUserInfos(participateIdoUserInfos *[]models.IdoParticipateUserInfo) error {
	return middleware.DB.Create(participateIdoUserInfos).Error
}

func ReadParticipateIdoUserInfo(id uint) (*models.IdoParticipateUserInfo, error) {
	var participateIdoUserInfo models.IdoParticipateUserInfo
	err := middleware.DB.First(&participateIdoUserInfo, id).Error
	return &participateIdoUserInfo, err
}

func ReadParticipateIdoUserInfosByAssetId(assetId string) (*[]models.IdoParticipateUserInfo, error) {
	var participateIdoUserInfos []models.IdoParticipateUserInfo
	err := middleware.DB.Where("asset_id = ?", assetId).Find(&participateIdoUserInfos).Error
	return &participateIdoUserInfos, err
}

func ReadParticipateIdoUserInfosByUserId(userId int) (*[]models.IdoParticipateUserInfo, error) {
	var participateIdoUserInfos []models.IdoParticipateUserInfo
	err := middleware.DB.Where("user_id = ?", userId).Find(&participateIdoUserInfos).Error
	return &participateIdoUserInfos, err
}

func UpdateParticipateIdoUserInfo(participateIdoUserInfo *models.IdoParticipateUserInfo) error {
	return middleware.DB.Save(participateIdoUserInfo).Error
}

func UpdateParticipateIdoUserInfos(participateIdoUserInfos *[]models.IdoParticipateUserInfo) error {
	return middleware.DB.Save(participateIdoUserInfos).Error
}

func DeleteParticipateIdoUserInfo(id uint) error {
	var participateIdoUserInfo models.IdoParticipateUserInfo
	return middleware.DB.Delete(&participateIdoUserInfo, id).Error
}
