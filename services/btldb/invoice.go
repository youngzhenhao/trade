package btldb

import (
	"gorm.io/gorm"
	"trade/middleware"
	"trade/models"
)

// CreateInvoice creates a new invoice record
func CreateInvoice(invoice *models.Invoice) error {
	return middleware.DB.Create(invoice).Error
}

// GetInvoice retrieves an invoice by Id
func GetInvoice(id uint) (*models.Invoice, error) {
	var invoice models.Invoice
	err := middleware.DB.First(&invoice, id).Error
	return &invoice, err
}

// GetInvoiceByReq retrieves an invoice by request
func GetInvoiceByReq(invoiceReq string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := middleware.DB.Where("invoice =?", invoiceReq).First(&invoice).Error
	return &invoice, err
}

// UpdateInvoice updates an existing invoice
func UpdateInvoice(db *gorm.DB, invoice *models.Invoice) error {
	return db.Save(invoice).Error
}

// DeleteInvoice soft deletes an invoice by Id
func DeleteInvoice(id uint) error {
	var invoice models.Invoice
	return middleware.DB.Delete(&invoice, id).Error
}
