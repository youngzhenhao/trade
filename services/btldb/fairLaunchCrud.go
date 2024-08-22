package btldb

import (
	"gorm.io/gorm"
	"trade/middleware"
	"trade/models"
	"trade/utils"
)

type FairLaunchStore struct {
	DB *gorm.DB
}

// FairLaunchInfo

func (f *FairLaunchStore) CreateFairLaunchInfo(fairLaunchInfo *models.FairLaunchInfo) error {
	return f.DB.Create(fairLaunchInfo).Error
}

func (f *FairLaunchStore) ReadFairLaunchInfo(id uint) (*models.FairLaunchInfo, error) {
	var fairLaunchInfo models.FairLaunchInfo
	err := f.DB.First(&fairLaunchInfo, id).Error
	return &fairLaunchInfo, err
}

func ReadClosedFairLaunchInfo() (*[]models.FairLaunchInfo, error) {
	var fairLaunchInfos []models.FairLaunchInfo
	now := utils.GetTimestamp()
	err := middleware.DB.Where("end_time < ?", now).Order("set_time").Find(&fairLaunchInfos).Error
	return &fairLaunchInfos, err
}

func ReadNotStartedFairLaunchInfo() (*[]models.FairLaunchInfo, error) {
	var fairLaunchInfos []models.FairLaunchInfo
	now := utils.GetTimestamp()
	err := middleware.DB.Where("start_time > ?", now).Order("set_time").Find(&fairLaunchInfos).Error
	return &fairLaunchInfos, err
}

func ReadIssuedFairLaunchInfos() (*[]models.FairLaunchInfo, error) {
	var fairLaunchInfos []models.FairLaunchInfo
	// @dev: Add more condition
	// @dev: Remove `AND is_mint_all = ?`
	err := middleware.DB.Where("status = ? AND state >= ?", models.StatusNormal, models.FairLaunchStateIssued).Order("set_time").Find(&fairLaunchInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchInfos")
	}
	return &fairLaunchInfos, nil
}

func ReadNotIssuedFairLaunchInfos() (*[]models.FairLaunchInfo, error) {
	var fairLaunchInfos []models.FairLaunchInfo
	err := middleware.DB.Where("status = ? AND state BETWEEN ? AND ?", models.StatusNormal, models.FairLaunchStateNoPay, models.FairLaunchStateIssuedPending).Order("set_time").Find(&fairLaunchInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find Not Issued fairLaunchInfos")
	}
	return &fairLaunchInfos, nil
}

// ReadIssuedAndTimeValidFairLaunchInfos
// @Description: Order by MintedNumber desc, SetTime
func ReadIssuedAndTimeValidFairLaunchInfos() (*[]models.FairLaunchInfo, error) {
	var fairLaunchInfos []models.FairLaunchInfo
	now := utils.GetTimestamp()
	err := middleware.DB.Where("state >= ? AND start_time <= ? AND end_time >= ? AND status = ?", models.FairLaunchStateIssued, now, now, models.StatusNormal).Order("minted_number desc, set_time").Find(&fairLaunchInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchInfos")
	}
	return &fairLaunchInfos, nil
}

func (f *FairLaunchStore) UpdateFairLaunchInfo(fairLaunchInfo *models.FairLaunchInfo) error {
	return f.DB.Save(fairLaunchInfo).Error
}

func UpdateFairLaunchInfo(fairLaunchInfo *models.FairLaunchInfo) error {
	return middleware.DB.Save(fairLaunchInfo).Error
}

func (f *FairLaunchStore) DeleteFairLaunchInfo(id uint) error {
	var fairLaunchInfo models.FairLaunchInfo
	return f.DB.Delete(&fairLaunchInfo, id).Error
}

func DeleteFairLaunchInfo(id uint) error {
	var fairLaunchInfo models.FairLaunchInfo
	return middleware.DB.Delete(&fairLaunchInfo, id).Error
}

// FairLaunchMintedInfo

func (f *FairLaunchStore) CreateFairLaunchMintedInfo(fairLaunchMintedInfo *models.FairLaunchMintedInfo) error {
	return f.DB.Create(fairLaunchMintedInfo).Error
}

func (f *FairLaunchStore) ReadFairLaunchMintedInfo(id uint) (*models.FairLaunchMintedInfo, error) {
	var fairLaunchMintedInfo models.FairLaunchMintedInfo
	err := f.DB.First(&fairLaunchMintedInfo, id).Error
	return &fairLaunchMintedInfo, err
}

func (f *FairLaunchStore) UpdateFairLaunchMintedInfo(fairLaunchMintedInfo *models.FairLaunchMintedInfo) error {
	return f.DB.Save(fairLaunchMintedInfo).Error
}

func UpdateFairLaunchMintedInfo(fairLaunchMintedInfo *models.FairLaunchMintedInfo) error {
	return middleware.DB.Save(fairLaunchMintedInfo).Error
}

func (f *FairLaunchStore) DeleteFairLaunchMintedInfo(id uint) error {
	var fairLaunchMintedInfo models.FairLaunchMintedInfo
	return f.DB.Delete(&fairLaunchMintedInfo, id).Error
}

func DeleteFairLaunchMintedInfo(id uint) error {
	var fairLaunchMintedInfo models.FairLaunchMintedInfo
	return middleware.DB.Delete(&fairLaunchMintedInfo, id).Error
}

func ReadNotSentFairLaunchMintedInfos() (*[]models.FairLaunchMintedInfo, error) {
	var fairLaunchMintedInfos []models.FairLaunchMintedInfo
	err := middleware.DB.Where("status = ? AND state BETWEEN ? AND ?", models.StatusNormal, models.FairLaunchMintedStateNoPay, models.FairLaunchMintedStateSentPending).Order("minted_set_time").Find(&fairLaunchMintedInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find Not Issued fairLaunchMintedInfos")
	}
	return &fairLaunchMintedInfos, nil
}

// FairLaunchInventoryInfo

func (f *FairLaunchStore) CreateFairLaunchInventoryInfo(fairLaunchInventoryInfo *models.FairLaunchInventoryInfo) error {
	return f.DB.Create(fairLaunchInventoryInfo).Error
}

func (f *FairLaunchStore) CreateFairLaunchInventoryInfos(fairLaunchInventoryInfos *[]models.FairLaunchInventoryInfo) error {
	return f.DB.Create(fairLaunchInventoryInfos).Error
}

func CreateFairLaunchInventoryInfos(fairLaunchInventoryInfos *[]models.FairLaunchInventoryInfo) error {
	return middleware.DB.Create(fairLaunchInventoryInfos).Error
}

func (f *FairLaunchStore) ReadFairLaunchInventoryInfo(id uint) (*models.FairLaunchInventoryInfo, error) {
	var fairLaunchInventoryInfo models.FairLaunchInventoryInfo
	err := f.DB.First(&fairLaunchInventoryInfo, id).Error
	return &fairLaunchInventoryInfo, err
}

func (f *FairLaunchStore) UpdateFairLaunchInventoryInfo(fairLaunchInventoryInfo *models.FairLaunchInventoryInfo) error {
	return f.DB.Save(fairLaunchInventoryInfo).Error
}

func (f *FairLaunchStore) DeleteFairLaunchInventoryInfo(id uint) error {
	var fairLaunchInventoryInfo models.FairLaunchInventoryInfo
	return f.DB.Delete(&fairLaunchInventoryInfo, id).Error
}

// FairLaunchMintedUserInfo

func (f *FairLaunchStore) CreateFairLaunchMintedUserInfo(fairLaunchMintedUserInfo *models.FairLaunchMintedUserInfo) error {
	return f.DB.Create(fairLaunchMintedUserInfo).Error
}

func (f *FairLaunchStore) ReadFairLaunchMintedUserInfo(id uint) (*models.FairLaunchMintedUserInfo, error) {
	var fairLaunchMintedUserInfo models.FairLaunchMintedUserInfo
	err := f.DB.First(&fairLaunchMintedUserInfo, id).Error
	return &fairLaunchMintedUserInfo, err
}

func (f *FairLaunchStore) UpdateFairLaunchMintedUserInfo(fairLaunchMintedUserInfo *models.FairLaunchMintedUserInfo) error {
	return f.DB.Save(fairLaunchMintedUserInfo).Error
}

func (f *FairLaunchStore) DeleteFairLaunchMintedUserInfo(id uint) error {
	var fairLaunchMintedUserInfo models.FairLaunchMintedUserInfo
	return f.DB.Delete(&fairLaunchMintedUserInfo, id).Error
}

func GetFairLaunchInfosByIds(fairLaunchInfoIds *[]int) (*[]models.FairLaunchInfo, error) {
	var fairLaunchInfos []models.FairLaunchInfo
	err := middleware.DB.Where(fairLaunchInfoIds).Find(&fairLaunchInfos).Error
	return &fairLaunchInfos, err
}

// FairLaunchMintedAndAvailableInfo

func CreateFairLaunchMintedAndAvailableInfo(fairLaunchMintedAndAvailableInfo *models.FairLaunchMintedAndAvailableInfo) error {
	return middleware.DB.Create(fairLaunchMintedAndAvailableInfo).Error
}

func CreateFairLaunchMintedAndAvailableInfos(fairLaunchMintedAndAvailableInfos *[]models.FairLaunchMintedAndAvailableInfo) error {
	return middleware.DB.Create(fairLaunchMintedAndAvailableInfos).Error
}

func ReadFairLaunchMintedAndAvailableInfo(id uint) (*models.FairLaunchMintedAndAvailableInfo, error) {
	var fairLaunchMintedAndAvailableInfo models.FairLaunchMintedAndAvailableInfo
	err := middleware.DB.First(&fairLaunchMintedAndAvailableInfo, id).Error
	return &fairLaunchMintedAndAvailableInfo, err
}

func ReadFairLaunchMintedAndAvailableInfoByFairLaunchInfoId(fairLaunchInfoId int) (*models.FairLaunchMintedAndAvailableInfo, error) {
	var fairLaunchMintedAndAvailableInfo models.FairLaunchMintedAndAvailableInfo
	err := middleware.DB.Where("fair_launch_info_id = ?", fairLaunchInfoId).First(&fairLaunchMintedAndAvailableInfo).Error
	return &fairLaunchMintedAndAvailableInfo, err
}

func UpdateFairLaunchMintedAndAvailableInfo(fairLaunchMintedAndAvailableInfo *models.FairLaunchMintedAndAvailableInfo) error {
	return middleware.DB.Save(fairLaunchMintedAndAvailableInfo).Error
}

func DeleteFairLaunchMintedAndAvailableInfo(id uint) error {
	var fairLaunchMintedAndAvailableInfo models.FairLaunchMintedAndAvailableInfo
	return middleware.DB.Delete(&fairLaunchMintedAndAvailableInfo, id).Error
}
