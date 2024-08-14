package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateFairLaunchIncome(fairLaunchIncome *models.FairLaunchIncome) error {
	return middleware.DB.Create(fairLaunchIncome).Error
}

func CreateFairLaunchIncomes(fairLaunchIncomes *[]models.FairLaunchIncome) error {
	return middleware.DB.Create(fairLaunchIncomes).Error
}

func ReadAllFairLaunchIncomes() (*[]models.FairLaunchIncome, error) {
	var fairLaunchIncomes []models.FairLaunchIncome
	err := middleware.DB.Find(&fairLaunchIncomes).Error
	return &fairLaunchIncomes, err
}

func ReadAllFairLaunchIncomesUpdatedAtDesc() (*[]models.FairLaunchIncome, error) {
	var fairLaunchIncomes []models.FairLaunchIncome
	err := middleware.DB.Order("updated_at desc").Find(&fairLaunchIncomes).Error
	return &fairLaunchIncomes, err
}

func ReadFairLaunchIncomesWhoseTxidIsNotNullAndSatAmountIsZero() (*[]models.FairLaunchIncome, error) {
	var fairLaunchIncomes []models.FairLaunchIncome
	err := middleware.DB.Where("txid <> ? AND sat_amount = ?", "", 0).Find(&fairLaunchIncomes).Error
	return &fairLaunchIncomes, err
}

func ReadFairLaunchIncome(id uint) (*models.FairLaunchIncome, error) {
	var fairLaunchIncome models.FairLaunchIncome
	err := middleware.DB.First(&fairLaunchIncome, id).Error
	return &fairLaunchIncome, err
}

// TODO: Useless
func ReadFairLaunchIncomesByUserId(userId int) (*[]models.FairLaunchIncome, error) {
	var fairLaunchIncomes []models.FairLaunchIncome
	err := middleware.DB.Where("user_id = ?", userId).Find(&fairLaunchIncomes).Error
	return &fairLaunchIncomes, err
}

// TODO: Useless
func ReadFairLaunchIncomesByUserIdUpdatedAtDesc(userId int) (*[]models.FairLaunchIncome, error) {
	var fairLaunchIncomes []models.FairLaunchIncome
	err := middleware.DB.Where("user_id = ?", userId).Order("updated_at desc").Find(&fairLaunchIncomes).Error
	return &fairLaunchIncomes, err
}

func ReadFairLaunchIncomesByAssetId(assetId string) (*[]models.FairLaunchIncome, error) {
	var fairLaunchIncomes []models.FairLaunchIncome
	err := middleware.DB.Where("asset_id = ?", assetId).Find(&fairLaunchIncomes).Error
	return &fairLaunchIncomes, err
}

func ReadFairLaunchIncomesByUserIdAndAssetId(userId int, assetId string) (*[]models.FairLaunchIncome, error) {
	var fairLaunchIncomes []models.FairLaunchIncome
	err := middleware.DB.Where("user_id = ? AND asset_id = ?", userId, assetId).Find(&fairLaunchIncomes).Error
	return &fairLaunchIncomes, err
}

func UpdateFairLaunchIncome(fairLaunchIncome *models.FairLaunchIncome) error {
	return middleware.DB.Save(fairLaunchIncome).Error
}

func UpdateFairLaunchIncomes(fairLaunchIncomes *[]models.FairLaunchIncome) error {
	if fairLaunchIncomes == nil {
		return nil
	}
	return middleware.DB.Save(fairLaunchIncomes).Error
}

func DeleteFairLaunchIncome(id uint) error {
	var fairLaunchIncome models.FairLaunchIncome
	return middleware.DB.Delete(&fairLaunchIncome, id).Error
}
