package btldb

import (
	"trade/middleware"
	"trade/models"
)

// CreateDeviceManager creates a new DeviceManager record
func CreateDeviceManager(device *models.DeviceManager) error {
	return middleware.DB.Create(device).Error
}

// ReadDeviceManagerByID retrieves a DeviceManager by Id
func ReadDeviceManagerByID(id uint) (*models.DeviceManager, error) {
	var device models.DeviceManager
	err := middleware.DB.First(&device, id).Error
	return &device, err
}

// ReadAllDeviceManagers retrieves all DeviceManager records
func ReadAllDeviceManagers() (*[]models.DeviceManager, error) {
	var devices []models.DeviceManager
	err := middleware.DB.Find(&devices).Error
	return &devices, err
}

// ReadDeviceManagerByNpubKey retrieves a DeviceManager by npub_key
func ReadDeviceManagerByNpubKey(npubKey string) (*models.DeviceManager, error) {
	var device models.DeviceManager
	err := middleware.DB.Where("npub_key = ?", npubKey).First(&device).Error
	return &device, err
}

// UpdateDeviceManager updates an existing DeviceManager record
func UpdateDeviceManager(device *models.DeviceManager) error {
	return middleware.DB.Save(device).Error
}

// DeleteDeviceManager soft deletes a DeviceManager record by Id
func DeleteDeviceManager(id uint) error {
	var device models.DeviceManager
	return middleware.DB.Delete(&device, id).Error
}
