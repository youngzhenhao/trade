package pAccount

type PAccountAssetId struct {
	Id            uint   `gorm:"primary_key"`
	PoolAccountId uint   `gorm:"index;column:pool_account_id;uniqueIndex:unique_pool_asset;not null"`
	AssetId       string `gorm:"column:asset_id;type:varchar(128);uniqueIndex:unique_pool_asset;not null"`
	// 外键关联
	PoolAccount *PoolAccount `gorm:"foreignkey:PoolAccountId"`
}

func (PAccountAssetId) TableName() string {
	return "custody_pool_account_assetId"
}
