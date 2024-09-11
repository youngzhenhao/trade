package models

import "gorm.io/gorm"

type AccountAssetReceive struct {
	gorm.Model
	Timestamp int64         `gorm:"column:timestamp;type:bigint;" json:"timestamp"`
	OutPoint  string        `gorm:"column:out_point;type:varchar(100);unique;not null;" json:"outPoint"`
	InvoiceId uint          `gorm:"column:invoice_id;type:varchar(128);" json:"invoiceId"`
	Amount    float64       `gorm:"type:decimal(15,2);column:amount" json:"amount"`
	Status    AddressStatus `gorm:"column:status;type:smallint" json:"status"`
}

func (AccountAssetReceive) TableName() string {
	return "user_account_Asset_Receive"
}

type AddressStatus int16

const (
	AddressStatusCOMPLETED AddressStatus = 4
)
