package btldb

import (
	"trade/middleware"
	"trade/models"
)

func CreateLogFileUpload(logFileUpload *models.LogFileUpload) error {
	return middleware.DB.Create(logFileUpload).Error
}

func CreateLogFileUploadAndGetId(logFileUpload *models.LogFileUpload) (uint, error) {
	err := middleware.DB.Create(logFileUpload).Error
	if err != nil {
		return 0, err
	}
	return logFileUpload.ID, nil
}

func CreateLogFileUploads(logFileUploads *[]models.LogFileUpload) error {
	return middleware.DB.Create(logFileUploads).Error
}

func ReadLogFileUpload(id uint) (*models.LogFileUpload, error) {
	var logFileUpload models.LogFileUpload
	err := middleware.DB.First(&logFileUpload, id).Error
	return &logFileUpload, err
}

func ReadAllLogFileUploads() (*[]models.LogFileUpload, error) {
	var logFileUploads []models.LogFileUpload
	err := middleware.DB.Order("updated_at desc").Find(&logFileUploads).Error
	return &logFileUploads, err
}

func ReadLogFileUploadsByDeviceId(deviceId string) (*[]models.LogFileUpload, error) {
	var logFileUploads []models.LogFileUpload
	err := middleware.DB.Where("device_id = ?", deviceId).Find(&logFileUploads).Error
	return &logFileUploads, err
}

func UpdateLogFileUpload(logFileUpload *models.LogFileUpload) error {
	return middleware.DB.Save(logFileUpload).Error
}

func UpdateLogFileUploads(logFileUploads *[]models.LogFileUpload) error {
	return middleware.DB.Save(logFileUploads).Error
}

func DeleteLogFileUpload(id uint) error {
	var logFileUpload models.LogFileUpload
	return middleware.DB.Delete(&logFileUpload, id).Error
}
