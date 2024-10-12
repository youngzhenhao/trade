package custodyModels

import "gorm.io/gorm"

type LockBillExt struct {
	gorm.Model
	BillId     uint                  `gorm:"column:bill_id;type:bigint unsigned;index:idx_bill_id" json:"billId"`
	LockId     string                `gorm:"column:lock_id;type:varchar(100);not null;unique;index:idx_lock_id" json:"lockId"`
	PayAccType LockBillExtPayAccType `gorm:"column:pay_acc_type;type:tinyint unsigned;default:0" json:"payAccType"`
	PayAccId   uint                  `gorm:"column:pay_acc_id;type:bigint unsigned;" json:"payAccId"`
	RevAccId   uint                  `gorm:"column:rev_acc_id;type:bigint unsigned;" json:"revAccId"`
	AssetId    string                `gorm:"column:asset_id;default:00;varchar(100)" json:"assetId"`
	Amount     float64               `gorm:"type:decimal(15,2);column:amount" json:"amount"`
	Status     LockBillExtStatus     `gorm:"column:status;type:tinyint unsigned;default:0" json:"status"`
}

func (LockBillExt) TableName() string {
	return "user_lock_bill_ext"
}

type LockBillExtPayAccType int8

const (
	LockBillExtPayAccTypeLock   LockBillExtPayAccType = 0
	LockBillExtPayAccTypeUnlock LockBillExtPayAccType = 1
)

type LockBillExtStatus int8

const (
	LockBillExtStatusInit    LockBillExtStatus = 0
	LockBillExtStatusSuccess LockBillExtStatus = 1
)
