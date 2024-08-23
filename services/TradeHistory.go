package services

import (
	"trade/middleware"
	"trade/models"
)

func Create(tradeHistory *models.TradeHistory) error {
	return middleware.DB.Create(tradeHistory).Error
}

func FindByID(id uint) (*models.TradeHistory, error) {
	var tradeHistory models.TradeHistory
	if err := middleware.DB.First(&tradeHistory, id).Error; err != nil {
		return nil, err
	}
	return &tradeHistory, nil
}
func Update(tradeHistory *models.TradeHistory) error {
	return middleware.DB.Save(tradeHistory).Error
}

func Delete(id uint) error {
	return middleware.DB.Delete(&models.TradeHistory{}, id).Error
}

func GetAggregatedTradeData(interval string) ([]models.AggregatedTradeData, error) {
	var results []models.AggregatedTradeData
	query := `
        SELECT
            date_trunc(?, trade_time) AS period,
            AVG(price) AS avg_price,
            SUM(quantity) AS total_quantity,
            COUNT(*) AS trade_count
        FROM
            trade_history
        GROUP BY
            date_trunc(?, trade_time)
        ORDER BY
            period;
    `
	err := middleware.DB.Raw(query, interval, interval).Scan(&results).Error
	return results, err
}
