package services

import (
	"trade/models"
	"trade/services/btldb"
)

func CreateLogFileUpload(logFileUpload *models.LogFileUpload) error {
	return btldb.CreateLogFileUpload(logFileUpload)
}

func GetFileUpload(id uint) (*models.LogFileUpload, error) {
	return btldb.ReadLogFileUpload(id)
}

func GetAllLogFiles() (*[]models.LogFileUpload, error) {
	return btldb.ReadAllLogFileUploads()
}
