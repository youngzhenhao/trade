package services

import (
	"trade/middleware"
	"trade/models"
)

// GetLatestAssetBalanceHistories
// sub query
func GetLatestAssetBalanceHistories(username string) (*[]models.AssetBalanceHistoryRecord, error) {
	var records []models.AssetBalanceHistoryRecord

	// SELECT *
	// FROM asset_balance_histories
	// WHERE id IN (
	//     SELECT MAX(id)
	//     FROM asset_balance_histories
	//     WHERE username = ?
	//     GROUP BY asset_id
	// );

	subQuery := middleware.DB.Model(&models.AssetBalanceHistory{}).
		Select("MAX(id)").
		Where("username = ?", username).
		Group("asset_id")

	err := middleware.DB.
		Table("asset_balance_histories").
		Select("id, asset_id, balance, username").
		Where("id IN (?)", subQuery).
		Scan(&records).Error

	return &records, err
}

func CreateAssetBalanceHistories(username string, requests *[]models.AssetBalanceHistorySetRequest) error {
	var records []models.AssetBalanceHistory
	if requests == nil || len(*requests) == 0 {
		return nil
	}
	for _, request := range *requests {
		records = append(records, models.AssetBalanceHistory{
			AssetId:  request.AssetId,
			Balance:  request.Balance,
			Username: username,
		})
	}
	return middleware.DB.Create(&records).Error
}
