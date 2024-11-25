package models

import "gorm.io/gorm"

type AssetBalanceHistory struct {
	gorm.Model
	AssetId  string `json:"asset_id" gorm:"type:varchar(255);index"`
	Balance  int    `json:"balance" gorm:"index"`
	Username string `json:"username" gorm:"type:varchar(255);index"`
}

type AssetBalanceHistorySetRequest struct {
	AssetId string `json:"asset_id" gorm:"type:varchar(255);index"`
	Balance int    `json:"balance" gorm:"index"`
}

type AssetBalanceHistoryRecord struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	AssetId  string `json:"asset_id" gorm:"type:varchar(255);index"`
	Balance  int    `json:"balance" gorm:"index"`
	Username string `json:"username" gorm:"type:varchar(255);index"`
}
