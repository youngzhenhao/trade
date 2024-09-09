package services

import (
	"trade/models"
	"trade/services/btldb"
)

func CreateLogFileUpload(logFileUpload *models.LogFileUpload) error {
	return btldb.CreateLogFileUpload(logFileUpload)
}
