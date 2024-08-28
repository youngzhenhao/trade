package btldb

import (
	"gorm.io/gorm"
	"sync"
	"trade/middleware"
	"trade/models"
)

var invoiceMutex sync.Mutex

// CreateInvoice creates a new invoice record
func CreateInvoice(invoice *models.Invoice) error {
	invoiceMutex.Lock()
	defer invoiceMutex.Unlock()
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
	invoiceMutex.Lock()
	defer invoiceMutex.Unlock()
	return db.Save(invoice).Error
}

// DeleteInvoice soft deletes an invoice by Id
func DeleteInvoice(id uint) error {
	invoiceMutex.Lock()
	defer invoiceMutex.Unlock()
	var invoice models.Invoice
	return middleware.DB.Delete(&invoice, id).Error
}
