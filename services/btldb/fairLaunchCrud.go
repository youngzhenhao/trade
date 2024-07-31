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
	err := middleware.DB.Where("end_time < ?", now).Find(&fairLaunchInfos).Error
	return &fairLaunchInfos, err
}

func ReadNotStartedFairLaunchInfo() (*[]models.FairLaunchInfo, error) {
	var fairLaunchInfos []models.FairLaunchInfo
	now := utils.GetTimestamp()
	err := middleware.DB.Where("start_time > ?", now).Find(&fairLaunchInfos).Error
	return &fairLaunchInfos, err
}

func (f *FairLaunchStore) UpdateFairLaunchInfo(fairLaunchInfo *models.FairLaunchInfo) error {
	return f.DB.Save(fairLaunchInfo).Error
}

func (f *FairLaunchStore) DeleteFairLaunchInfo(id uint) error {
	var fairLaunchInfo models.FairLaunchInfo
	return f.DB.Delete(&fairLaunchInfo, id).Error
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

func (f *FairLaunchStore) DeleteFairLaunchMintedInfo(id uint) error {
	var fairLaunchMintedInfo models.FairLaunchMintedInfo
	return f.DB.Delete(&fairLaunchMintedInfo, id).Error
}

// FairLaunchInventoryInfo

func (f *FairLaunchStore) CreateFairLaunchInventoryInfo(fairLaunchInventoryInfo *models.FairLaunchInventoryInfo) error {
	return f.DB.Create(fairLaunchInventoryInfo).Error
}

func (f *FairLaunchStore) CreateFairLaunchInventoryInfos(fairLaunchInventoryInfos *[]models.FairLaunchInventoryInfo) error {
	return f.DB.Create(fairLaunchInventoryInfos).Error
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
