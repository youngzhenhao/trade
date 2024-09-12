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

func LogFileUploadsToDeviceIdMapLogFileUploads(logFileUploads *[]models.LogFileUpload) *map[string]*[]models.LogFileUpload {
	if logFileUploads == nil {
		return nil
	}
	deviceIdMapLogFileUploads := make(map[string]*[]models.LogFileUpload)
	for _, logFileUpload := range *logFileUploads {
		deviceId := logFileUpload.DeviceId
		mapLogFileUploads, ok := deviceIdMapLogFileUploads[deviceId]
		if !ok {
			deviceIdMapLogFileUploads[deviceId] = &[]models.LogFileUpload{logFileUpload}
		} else {
			*mapLogFileUploads = append(*mapLogFileUploads, logFileUpload)
		}
	}
	return &deviceIdMapLogFileUploads
}
