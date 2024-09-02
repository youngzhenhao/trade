package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateAssetLocalMint(assetLocalMint *models.AssetLocalMint) error {
	return middleware.DB.Create(assetLocalMint).Error
}

func CreateAssetLocalMints(assetLocalMints *[]models.AssetLocalMint) error {
	return middleware.DB.Create(assetLocalMints).Error
}

func ReadAllAssetLocalMints() (*[]models.AssetLocalMint, error) {
	var assetLocalMints []models.AssetLocalMint
	err := middleware.DB.Find(&assetLocalMints).Error
	return &assetLocalMints, err
}

func ReadAllAssetLocalMintsUpdatedAtDesc() (*[]models.AssetLocalMint, error) {
	var assetLocalMints []models.AssetLocalMint
	err := middleware.DB.Order("updated_at desc").Find(&assetLocalMints).Error
	return &assetLocalMints, err
}

func ReadAssetLocalMint(id uint) (*models.AssetLocalMint, error) {
	var assetLocalMint models.AssetLocalMint
	err := middleware.DB.First(&assetLocalMint, id).Error
	return &assetLocalMint, err
}

func ReadAssetLocalMintsByUserId(userId int) (*[]models.AssetLocalMint, error) {
	var assetLocalMints []models.AssetLocalMint
	err := middleware.DB.Where("user_id = ?", userId).Find(&assetLocalMints).Error
	return &assetLocalMints, err
}

func ReadAssetLocalMintByAssetId(assetId string) (*models.AssetLocalMint, error) {
	var assetLocalMint models.AssetLocalMint
	err := middleware.DB.Where("asset_id = ?", assetId).First(&assetLocalMint).Error
	return &assetLocalMint, err
}

func UpdateAssetLocalMint(assetLocalMint *models.AssetLocalMint) error {
	return middleware.DB.Save(assetLocalMint).Error
}

func UpdateAssetLocalMints(assetLocalMints *[]models.AssetLocalMint) error {
	return middleware.DB.Save(assetLocalMints).Error
}

func DeleteAssetLocalMint(id uint) error {
	var assetLocalMint models.AssetLocalMint
	return middleware.DB.Delete(&assetLocalMint, id).Error
}
