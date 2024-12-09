package pAccount

type PAccountBalanceChange struct {
	Id            uint    `gorm:"primary_key"`
	PoolAccountId uint    `gorm:"index;column:pool_account_id;not null"`
	AssetId       string  `gorm:"column:asset_id;varchar(128);not null"`
	BillId        uint    `gorm:"column:bill_id;unique;"`
	Amount        float64 `gorm:"column:amount;type:decimal(15,2);"`
	FinalBalance  float64 `gorm:"column:final_balance;type:decimal(15,2);not null"`
	// 外键关联
	PoolAccount *PoolAccount `gorm:"foreignkey:PoolAccountId"`
}

func (PAccountBalanceChange) TableName() string {
	return "custody_pool_account_balance_change"
}
