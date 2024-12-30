package custodyModels

import "gorm.io/gorm"

type AccountBalanceChange struct {
	gorm.Model
	AccountId    uint       `gorm:"column:account_id;type:bigint unsigned;" json:"accountId"`
	AssetId      string     `gorm:"column:asset_id;type:varchar(128);" json:"assetId"`
	ChangeAmount float64    `gorm:"type:decimal(15,2);column:amount" json:"amount"`
	Away         ChangeAway `gorm:"column:away;type:tinyint unsigned" json:"away"`
	FinalBalance float64    `gorm:"type:decimal(15,2);column:final_balance" json:"finalBalance"`
	BalanceId    uint       `gorm:"column:balance_id;type:bigint unsigned;index:idx_balance_id" json:"balanceId"`
	ChangeType   ChangeType `gorm:"column:change_type;type:varchar(128)" json:"changeType"`
}

func (AccountBalanceChange) TableName() string {
	return "user_account_changes"
}

type ChangeAway uint

const (
	ChangeAwayAdd  ChangeAway = 0
	ChangeAwayLess ChangeAway = 1
)

type ChangeType string

const (
	ChangeTypeFault             = "fault"
	ChangeTypeBtcPayOutside     = "pay_outside_btc"
	ChangeTypeBtcReceiveOutside = "receive_outside_btc"
	ChangeTypeBtcFee            = "btc_fee"
	ChangeFirLunchFee           = "fir_lunch_fee"
	ChangeTypeBtcPayLocal       = "pay_local_btc"
	ChangeTypeBtcReceiveLocal   = "receive_local_btc"
	ChangeTypeBackFee           = "back_fee"
	ChangeTypeAward             = "award"

	ChangeTypePayToPoolAccount       = "pay_to_pool_account"
	ChangeTypeReceiveFromPoolAccount = "receive_from_pool_account"

	ChangeTypeAssetPayOutside    = "pay_outside_asset"
	ChangTypeAssetReceiveOutside = "receive_outside_asset"

	ChangeTypeAssetPayLocal     = "pay_local_asset"
	ChangeTypeAssetReceiveLocal = "receive_local_asset"

	ChangeTypeLock           = "lock"
	ChangeTypeUnlock         = "unlock"
	ChangeTypeLockedTransfer = "locked_transfer"

	ClearLimitUser = "clear_limit_user"
)
