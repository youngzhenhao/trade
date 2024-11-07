package custodyModels

import "gorm.io/gorm"

type LockBill struct {
	gorm.Model
	AccountID uint         `gorm:"column:account_id;type:bigint unsigned;index:idx_account_id" json:"accountId"`
	LockId    string       `gorm:"column:lock_id;type:varchar(100);not null;unique;index:idx_lock_id" json:"lockId"`
	BillType  LockBillType `gorm:"column:bill_type;type:smallint" json:"billType"`
	AssetId   string       `gorm:"column:asset_id;default:00;varchar(100)" json:"assetId"`
	Amount    float64      `gorm:"type:decimal(15,2);column:amount" json:"amount"`
}

func (LockBill) TableName() string {
	return "user_lock_bill"
}

type LockBillType int8

const (
	LockBillTypeLock LockBillType = iota
	LockBillTypeTransferByLockAsset
	LockBillTypeUnlock
	LockBillTypeTransferByUnlockAsset

	LockBillTypeAward LockBillType = 5

	LockErr = 66
)
