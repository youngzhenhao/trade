package btldb

import (
	"gorm.io/gorm"
	"trade/models"
)

type AssetIssuanceStore struct {
	DB *gorm.DB
}

// AssetIssuance

func (a *AssetIssuanceStore) CreateAssetIssuance(tx *gorm.DB, assetIssuance *models.AssetIssuance) error {
	return tx.Create(assetIssuance).Error
}

func (a *AssetIssuanceStore) ReadAssetIssuance(id uint) (*models.AssetIssuance, error) {
	var assetIssuance models.AssetIssuance
	err := a.DB.First(&assetIssuance, id).Error
	return &assetIssuance, err
}

func (a *AssetIssuanceStore) ReadAssetIssuanceByFairLaunchId(fairLaunchId uint) (*models.AssetIssuance, error) {
	var assetIssuance models.AssetIssuance
	err := a.DB.Where("is_fair_launch = ? AND fair_launch_id = ?", true, fairLaunchId).First(&assetIssuance).Error
	return &assetIssuance, err
}

func (a *AssetIssuanceStore) UpdateAssetIssuance(tx *gorm.DB, assetIssuance *models.AssetIssuance) error {
	return tx.Save(assetIssuance).Error
}

func (a *AssetIssuanceStore) DeleteAssetIssuance(id uint) error {
	var assetIssuance models.AssetIssuance
	return a.DB.Delete(&assetIssuance, id).Error
}
