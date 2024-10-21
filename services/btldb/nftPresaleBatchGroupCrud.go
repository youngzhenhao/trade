package btldb

import (
	"trade/middleware"
	"trade/models"
	"trade/utils"
)

func CreateNftPresaleBatchGroup(nftPresaleBatchGroup *models.NftPresaleBatchGroup) error {
	return middleware.DB.Create(nftPresaleBatchGroup).Error
}

func CreateNftPresaleBatchGroups(nftPresaleBatchGroups *[]models.NftPresaleBatchGroup) error {
	return middleware.DB.Create(nftPresaleBatchGroups).Error
}

func ReadNftPresaleBatchGroup(id uint) (*models.NftPresaleBatchGroup, error) {
	var nftPresaleBatchGroup models.NftPresaleBatchGroup
	err := middleware.DB.First(&nftPresaleBatchGroup, id).Error
	return &nftPresaleBatchGroup, err
}

func ReadNftPresaleBatchGroupByGroupKey(groupKey string) (*models.NftPresaleBatchGroup, error) {
	var nftPresaleBatchGroup models.NftPresaleBatchGroup
	err := middleware.DB.Where("group_key = ?", groupKey).First(&nftPresaleBatchGroup).Error
	return &nftPresaleBatchGroup, err
}

func ReadAllNftPresaleBatchGroups() (*[]models.NftPresaleBatchGroup, error) {
	var nftPresaleBatchGroups []models.NftPresaleBatchGroup
	err := middleware.DB.Order("start_time desc").Find(&nftPresaleBatchGroups).Error
	return &nftPresaleBatchGroups, err
}

func ReadSellingNftPresaleBatchGroups() (*[]models.NftPresaleBatchGroup, error) {
	var nftPresaleBatchGroups []models.NftPresaleBatchGroup
	now := utils.GetTimestamp()
	err := middleware.DB.Where("start_time <= ? AND end_time = ?", now, 0).Or("start_time <= ? AND end_time >= ?", now, now).Order("start_time desc").Find(&nftPresaleBatchGroups).Error
	return &nftPresaleBatchGroups, err
}

func ReadNotStartNftPresaleBatchGroups() (*[]models.NftPresaleBatchGroup, error) {
	var nftPresaleBatchGroups []models.NftPresaleBatchGroup
	now := utils.GetTimestamp()
	err := middleware.DB.Where("start_time > ?", now).Order("start_time desc").Find(&nftPresaleBatchGroups).Error
	return &nftPresaleBatchGroups, err
}

func ReadEndNftPresaleBatchGroups() (*[]models.NftPresaleBatchGroup, error) {
	var nftPresaleBatchGroups []models.NftPresaleBatchGroup
	now := utils.GetTimestamp()
	err := middleware.DB.Where("end_time BETWEEN ? AND ?", 1, now).Order("start_time desc").Find(&nftPresaleBatchGroups).Error
	return &nftPresaleBatchGroups, err
}

func UpdateNftPresaleBatchGroup(nftPresaleBatchGroup *models.NftPresaleBatchGroup) error {
	return middleware.DB.Save(nftPresaleBatchGroup).Error
}

func UpdateNftPresaleBatchGroups(nftPresaleBatchGroups *[]models.NftPresaleBatchGroup) error {
	return middleware.DB.Save(nftPresaleBatchGroups).Error
}

func DeleteNftPresaleBatchGroup(id uint) error {
	var nftPresaleBatchGroup models.NftPresaleBatchGroup
	return middleware.DB.Delete(&nftPresaleBatchGroup, id).Error
}
