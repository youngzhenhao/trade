package pAccount

type PAccountBalance struct {
	Id            uint    `gorm:"primary_key"`
	PoolAccountId uint    `gorm:"index;column:pool_account_id;uniqueIndex:unique_pool_asset;not null"`
	AssetId       string  `gorm:"column:asset_id;type:varchar(128);uniqueIndex:unique_pool_asset;not null"`
	Balance       float64 `gorm:"column:balance;type:decimal(15,2);"`
	// 外键关联
	PoolAccount *PoolAccount `gorm:"foreignkey:PoolAccountId"`
}

func (PAccountBalance) TableName() string {
	return "custody_pool_account_balances"
}
