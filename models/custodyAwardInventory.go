package models

import "gorm.io/gorm"

type AwardInventory struct {
	gorm.Model
	AssetId     string               `gorm:"column:asset_id;type:varchar(128);unique" json:"assetId"`
	Amount      float64              `gorm:"column:amount;type:decimal(15,2)" json:"amount"`
	TotalAmount float64              `gorm:"column:total_amount;type:decimal(15,2)" json:"totalAmount"`
	Status      AwardInventoryStatus `gorm:"column:status;type:int" json:"status"`
}

func (AwardInventory) TableName() string {
	return "user_account_award_inventory"
}

type AwardInventoryStatus uint

const AwardInventoryAble AwardInventoryStatus = 1
