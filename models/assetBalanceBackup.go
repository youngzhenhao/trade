package models

import "gorm.io/gorm"

type AssetBalanceBackup struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(255);uniqueIndex"`
	Hash     string `json:"hash" gorm:"type:varchar(255);index"`
}

type AssetBalanceBackupSetRequest struct {
	Hash string `json:"hash" gorm:"type:varchar(255);index"`
}
