package services

import (
	"trade/middleware"
	"trade/models"
)

// CreatePayInside creates a new PayInside record
func CreatePayInside(payInside *models.PayInside) error {
	return middleware.DB.Create(payInside).Error
}

// ReadPayInside retrieves a PayInside by Id
func ReadPayInside(id uint) (*models.PayInside, error) {
	var payInside models.PayInside
	err := middleware.DB.First(&payInside, id).Error
	return &payInside, err
}

// UpdatePayInside updates an existing PayInside
func UpdatePayInside(payInside *models.PayInside) error {
	return middleware.DB.Save(payInside).Error
}

// DeletePayInside soft deletes a PayInside by Id
func DeletePayInside(id uint) error {
	var payInside models.PayInside
	return middleware.DB.Delete(&payInside, id).Error
}
