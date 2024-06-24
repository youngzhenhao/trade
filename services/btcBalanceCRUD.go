package services

import (
	"trade/middleware"
	"trade/models"
)

func CreateBtcBalance(btcBalance *models.BtcBalance) error {
	return middleware.DB.Create(btcBalance).Error
}

func CreateBtcBalances(btcBalances *[]models.BtcBalance) error {
	return middleware.DB.Create(btcBalances).Error
}

func ReadAllBtcBalances() (*[]models.BtcBalance, error) {
	var btcBalances []models.BtcBalance
	err := middleware.DB.Find(&btcBalances).Error
	return &btcBalances, err
}

func ReadBtcBalance(id uint) (*models.BtcBalance, error) {
	var btcBalance models.BtcBalance
	err := middleware.DB.First(&btcBalance, id).Error
	return &btcBalance, err
}

func ReadBtcBalanceByUsername(username string) (*models.BtcBalance, error) {
	var btcBalance models.BtcBalance
	err := middleware.DB.Where("username = ? AND status = ?", username, 1).First(&btcBalance).Error
	return &btcBalance, err
}

func UpdateBtcBalance(btcBalance *models.BtcBalance) error {
	return middleware.DB.Save(btcBalance).Error
}

func UpdateBtcBalances(btcBalances *[]models.BtcBalance) error {
	return middleware.DB.Save(btcBalances).Error
}

func DeleteBtcBalance(id uint) error {
	var btcBalance models.BtcBalance
	return middleware.DB.Delete(&btcBalance, id).Error
}
