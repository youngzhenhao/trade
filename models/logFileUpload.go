package models

import "gorm.io/gorm"

type LogFileUpload struct {
	gorm.Model
	DeviceId       string `json:"device_id"`
	OriginFileName string `json:"origin_file_name"`
	FileSavePath   string `json:"file_save_path"`
}
