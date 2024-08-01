package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateFairLaunchFollow(fairLaunchFollow *models.FairLaunchFollow) error {
	return middleware.DB.Create(fairLaunchFollow).Error
}

func CreateFairLaunchFollows(fairLaunchFollows *[]models.FairLaunchFollow) error {
	return middleware.DB.Create(fairLaunchFollows).Error
}

func ReadAllFairLaunchFollows() (*[]models.FairLaunchFollow, error) {
	var fairLaunchFollows []models.FairLaunchFollow
	err := middleware.DB.Find(&fairLaunchFollows).Error
	return &fairLaunchFollows, err
}

func ReadAllFairLaunchFollowsUpdatedAtDesc() (*[]models.FairLaunchFollow, error) {
	var fairLaunchFollows []models.FairLaunchFollow
	err := middleware.DB.Order("updated_at desc").Find(&fairLaunchFollows).Error
	return &fairLaunchFollows, err
}

func ReadFairLaunchFollow(id uint) (*models.FairLaunchFollow, error) {
	var fairLaunchFollow models.FairLaunchFollow
	err := middleware.DB.First(&fairLaunchFollow, id).Error
	return &fairLaunchFollow, err
}

func ReadFairLaunchFollowsByUserId(userId int) (*[]models.FairLaunchFollow, error) {
	var fairLaunchFollows []models.FairLaunchFollow
	err := middleware.DB.Where("user_id = ? AND status = ?", userId, 1).Find(&fairLaunchFollows).Error
	return &fairLaunchFollows, err
}

func ReadFairLaunchFollowsByUserIdUpdatedAtDesc(userId int) (*[]models.FairLaunchFollow, error) {
	var fairLaunchFollows []models.FairLaunchFollow
	err := middleware.DB.Where("user_id = ? AND status = ?", userId, 1).Order("updated_at desc").Find(&fairLaunchFollows).Error
	return &fairLaunchFollows, err
}

func ReadFairLaunchFollowByAssetId(assetId string) (*models.FairLaunchFollow, error) {
	var fairLaunchFollow models.FairLaunchFollow
	err := middleware.DB.Where("asset_id = ? AND status = ?", assetId, 1).First(&fairLaunchFollow).Error
	return &fairLaunchFollow, err
}

func ReadFairLaunchFollowByUserIdAndAssetId(userId int, assetId string) (*models.FairLaunchFollow, error) {
	var fairLaunchFollow models.FairLaunchFollow
	err := middleware.DB.Where("user_id = ? AND asset_id = ? AND status = ?", userId, assetId, 1).First(&fairLaunchFollow).Error
	return &fairLaunchFollow, err
}

func UpdateFairLaunchFollow(fairLaunchFollow *models.FairLaunchFollow) error {
	return middleware.DB.Save(fairLaunchFollow).Error
}

func UpdateFairLaunchFollows(fairLaunchFollows *[]models.FairLaunchFollow) error {
	return middleware.DB.Save(fairLaunchFollows).Error
}

func DeleteFairLaunchFollow(id uint) error {
	var fairLaunchFollow models.FairLaunchFollow
	return middleware.DB.Delete(&fairLaunchFollow, id).Error
}
