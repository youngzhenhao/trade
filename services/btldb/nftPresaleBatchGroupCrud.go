package btldb

import (
	"trade/middleware"
	"trade/models"
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
