package models

import "gorm.io/gorm"

// TODO
type AssetBalanceBackup struct {
	gorm.Model
	AssetId      string `json:"asset_id" gorm:"type:varchar(255);index"`
	IsSat        bool   `json:"is_sat" gorm:"index"`
	Balance      int    `json:"balance" gorm:"index"`
	Username     string `json:"username" gorm:"type:varchar(255);index"`
	BackupFileId uint   `json:"backup_file_id" gorm:"index"`
}

type AssetBalanceBackupSetRequest struct {
	AssetId string `json:"asset_id" gorm:"type:varchar(255);index"`
	IsSat   bool   `json:"is_sat" gorm:"index"`
	Balance int    `json:"balance" gorm:"index"`
}
