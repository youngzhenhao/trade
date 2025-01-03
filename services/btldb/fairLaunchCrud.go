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
	err := middleware.DB.Where("state >= ?", models.FairLaunchStateIssued).Order("set_time").Find(&fairLaunchInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchInfos")
	}
	return &fairLaunchInfos, nil
}

func ReadNotIssuedFairLaunchInfos() (*[]models.FairLaunchInfo, error) {
	var fairLaunchInfos []models.FairLaunchInfo
	err := middleware.DB.Where("state BETWEEN ? AND ?", models.FairLaunchStateNoPay, models.FairLaunchStateIssuedPending).Order("set_time").Find(&fairLaunchInfos).Error
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
	err := middleware.DB.Where("state >= ? AND start_time <= ? AND end_time >= ?", models.FairLaunchStateIssued, now, now).Order("minted_number desc, set_time").Find(&fairLaunchInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find fairLaunchInfos")
	}
	return &fairLaunchInfos, nil
}

func (f *FairLaunchStore) UpdateFairLaunchInfo(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) error {
	return tx.Save(fairLaunchInfo).Error
}

func UpdateFairLaunchInfo(tx *gorm.DB, fairLaunchInfo *models.FairLaunchInfo) error {
	return tx.Save(fairLaunchInfo).Error
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

func (f *FairLaunchStore) UpdateFairLaunchMintedInfo(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo) error {
	return tx.Save(fairLaunchMintedInfo).Error
}

func UpdateFairLaunchMintedInfo(tx *gorm.DB, fairLaunchMintedInfo *models.FairLaunchMintedInfo) error {
	return tx.Save(fairLaunchMintedInfo).Error
}

func UpdateFairLaunchMintedInfos(fairLaunchMintedInfos *[]models.FairLaunchMintedInfo) error {
	return middleware.DB.Save(fairLaunchMintedInfos).Error
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
	err := middleware.DB.Where("state BETWEEN ? AND ?", models.FairLaunchMintedStateNoPay, models.FairLaunchMintedStateSentPending).Order("minted_set_time").Find(&fairLaunchMintedInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Find Not Issued fairLaunchMintedInfos")
	}
	return &fairLaunchMintedInfos, nil
}

func ReadFairLaunchMintedInfoWhoseProcessNumberIsMoreThanTenThousand() (*[]models.FairLaunchMintedInfo, error) {
	var fairLaunchMintedInfos []models.FairLaunchMintedInfo
	err := middleware.DB.Where("state BETWEEN ? AND ? AND process_number > ?", models.FairLaunchMintedStateNoPay, models.FairLaunchMintedStateSentPending, 10000).Find(&fairLaunchMintedInfos).Error
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Read fairLaunchMintedInfos")
	}
	return &fairLaunchMintedInfos, nil
}

func ReadFairLaunchMintedInfosWhoseUsernameIsNull() (*[]models.FairLaunchMintedInfo, error) {
	var fairLaunchMintedInfos []models.FairLaunchMintedInfo
	err := middleware.DB.Where("username = ?", "").Find(&fairLaunchMintedInfos).Error
	if err != nil {
		return nil, err
	}
	return &fairLaunchMintedInfos, nil
}

func ReadUserFirstFairLaunchMintedInfoByUserId(userId int) (*models.FairLaunchMintedInfo, error) {
	var fairLaunchMintedInfo models.FairLaunchMintedInfo
	err := middleware.DB.Where("state > ? AND user_id = ?", models.FairLaunchMintedStateFail, userId).Order("minted_set_time").First(&fairLaunchMintedInfo).Error
	if err != nil {
		return nil, err
	}
	return &fairLaunchMintedInfo, nil
}

func ReadUserFirstFairLaunchMintedInfoByUserIdAndAssetId(userId int, assetId string) (*models.FairLaunchMintedInfo, error) {
	var fairLaunchMintedInfo models.FairLaunchMintedInfo
	err := middleware.DB.Where("state > ? AND user_id = ? AND asset_id = ?", models.FairLaunchMintedStateFail, userId, assetId).Order("minted_set_time").First(&fairLaunchMintedInfo).Error
	if err != nil {
		return nil, err
	}
	return &fairLaunchMintedInfo, nil
}

func ReadUserFirstFairLaunchMintedInfoByUsernameAndAssetId(username string, assetId string) (*models.FairLaunchMintedInfo, error) {
	var fairLaunchMintedInfo models.FairLaunchMintedInfo
	err := middleware.DB.Where("state > ? AND username = ? AND asset_id = ?", models.FairLaunchMintedStateFail, username, assetId).Order("minted_set_time").First(&fairLaunchMintedInfo).Error
	if err != nil {
		return nil, err
	}
	return &fairLaunchMintedInfo, nil
}

// FairLaunchInventoryInfo

func (f *FairLaunchStore) CreateFairLaunchInventoryInfo(fairLaunchInventoryInfo *models.FairLaunchInventoryInfo) error {
	return f.DB.Create(fairLaunchInventoryInfo).Error
}

func (f *FairLaunchStore) CreateFairLaunchInventoryInfos(fairLaunchInventoryInfos *[]models.FairLaunchInventoryInfo) error {
	return f.DB.Create(fairLaunchInventoryInfos).Error
}

func CreateFairLaunchInventoryInfos(tx *gorm.DB, fairLaunchInventoryInfos *[]models.FairLaunchInventoryInfo) error {
	return tx.Create(fairLaunchInventoryInfos).Error
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

func (f *FairLaunchStore) CreateFairLaunchMintedUserInfo(tx *gorm.DB, fairLaunchMintedUserInfo *models.FairLaunchMintedUserInfo) error {
	return tx.Create(fairLaunchMintedUserInfo).Error
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

func CreateFairLaunchMintedAndAvailableInfo(tx *gorm.DB, fairLaunchMintedAndAvailableInfo *models.FairLaunchMintedAndAvailableInfo) error {
	return tx.Create(fairLaunchMintedAndAvailableInfo).Error
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

func UpdateFairLaunchMintedAndAvailableInfo(tx *gorm.DB, fairLaunchMintedAndAvailableInfo *models.FairLaunchMintedAndAvailableInfo) error {
	return tx.Save(fairLaunchMintedAndAvailableInfo).Error
}

func DeleteFairLaunchMintedAndAvailableInfo(id uint) error {
	var fairLaunchMintedAndAvailableInfo models.FairLaunchMintedAndAvailableInfo
	return middleware.DB.Delete(&fairLaunchMintedAndAvailableInfo, id).Error
}
