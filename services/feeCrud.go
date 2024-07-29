package services

import (
	"gorm.io/gorm"
	"trade/models"
)

type FeeRateInfoStore struct {
	DB *gorm.DB
}

// FeeRateInfo

func (f *FeeRateInfoStore) CreateFeeRateInfo(feeRateInfo *models.FeeRateInfo) error {
	return f.DB.Create(feeRateInfo).Error
}

func (f *FeeRateInfoStore) ReadFeeRateInfo(id uint) (*models.FeeRateInfo, error) {
	var feeRateInfo models.FeeRateInfo
	err := f.DB.First(&feeRateInfo, id).Error
	return &feeRateInfo, err
}

func (f *FeeRateInfoStore) UpdateFeeRateInfo(feeRateInfo *models.FeeRateInfo) error {
	return f.DB.Save(feeRateInfo).Error
}

func (f *FeeRateInfoStore) DeleteFeeRateInfo(id uint) error {
	var feeRateInfo models.FeeRateInfo
	return f.DB.Delete(&feeRateInfo, id).Error
}
