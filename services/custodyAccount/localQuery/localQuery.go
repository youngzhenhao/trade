package localQuery

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
	"trade/services/servicesrpc"
)

const (
	DefaultAccount = "default"
	LockedAccount  = "locked"
)

type BillQueryQuest struct {
	UserName      string  `json:"username"`
	Away          int     `json:"away"`
	AssetId       string  `json:"assetId"`
	AmountMin     float64 `json:"amountMin"`
	AmountMax     float64 `json:"amountMax"`
	TimeStart     string  `json:"timeStart"`
	TimeEnd       string  `json:"timeEnd"`
	IncludeFailed bool    `json:"includeFailed"`
	OnlyAward     bool    `json:"onlyAward"`
	Page          int     `json:"page"`
	PageSize      int     `json:"pageSize"`
}

type BillListWithUser struct {
	ID          uint                `gorm:"primarykey" json:"id"`
	UserName    string              `gorm:"column:user_name" json:"username"`
	BillType    models.BalanceType  `gorm:"column:bill_type;type:smallint" json:"billType"`
	Away        models.BalanceAway  `gorm:"column:away;type:smallint" json:"away"`
	Amount      float64             `gorm:"column:amount;type:decimal(10,2)" json:"amount"`
	ServerFee   uint64              `gorm:"column:server_fee;type:bigint unsigned" json:"serverFee"`
	AssetId     *string             `gorm:"column:asset_id;type:varchar(512);default:'00'" json:"assetId"`
	Invoice     *string             `gorm:"column:invoice;type:varchar(512)" json:"invoice"`
	PaymentHash *string             `gorm:"column:payment_hash;type:varchar(100)" json:"paymentHash"`
	State       models.BalanceState `gorm:"column:State;type:smallint" json:"State"`
	Time        time.Time           `gorm:"column:created_at" json:"time"`
}

func BillQuery(quest BillQueryQuest) (*[]BillListWithUser, int64, error) {
	billQuery := models.Balance{}
	var err error

	db := middleware.DB

	q := db.Where(&billQuery)
	if !quest.IncludeFailed {
		q = q.Where("state =?", 1)
	}

	if quest.UserName != "" {
		account := models.Account{
			UserName: quest.UserName,
		}
		err = db.Where(&account).First(&account).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, 0, errors.New("account not found")
			}
			return nil, 0, err
		}
		q = q.Where("account_id =?", account.ID)
	}

	switch quest.Away {
	case 0, 1:
		q = q.Where("away =?", quest.Away)
	default:
	}
	if quest.AssetId != "" {
		q = q.Where("bill_balance.asset_id =?", quest.AssetId)
	}
	if quest.AmountMin != 0 {
		q = q.Where("bill_balance.amount >=?", quest.AmountMin)
	}
	if quest.AmountMax != 0 {
		q = q.Where("bill_balance.amount <=?", quest.AmountMax)
	}
	if quest.TimeStart != "" {
		q = q.Where("bill_balance.created_at >=?", quest.TimeStart)
	}
	if quest.TimeEnd != "" {
		q = q.Where("bill_balance.created_at <=?", quest.TimeEnd)
	}
	if quest.OnlyAward {
		q = q.Where("bill_balance.bill_type =? or bill_balance.bill_type =?", 5, 6)
	}
	if quest.PageSize == 0 {
		quest.PageSize = 500
	}

	// 查询总记录数
	var total int64
	err = q.Model(&models.Balance{}).Count(&total).Error
	if err != nil || total == 0 {
		return nil, 0, err
	}

	var billListWithUser []BillListWithUser
	err = q.Table("bill_balance").
		Joins("LEFT JOIN user_account ON user_account.id = bill_balance.account_id").
		Limit(quest.PageSize).Offset((quest.Page) * quest.PageSize).
		Select("bill_balance.*, user_account.user_name").
		Order("bill_balance.created_at DESC").
		Scan(&billListWithUser).Error
	if err != nil {
		return nil, 0, err
	}
	return &billListWithUser, total, nil
}

type BalanceQueryQuest struct {
	UserName string `json:"userName"`
}

type BalanceQueryResp struct {
	AccountName string  `json:"accountName"`
	AssetId     string  `json:"assetId"`
	Balance     float64 `json:"balance"`
}

func BalanceQuery(quest BalanceQueryQuest) *[]BalanceQueryResp {
	db := middleware.DB
	var err error
	var balances []BalanceQueryResp

	if quest.UserName == "" {
		return &balances
	}
	account := models.Account{
		UserName: quest.UserName,
	}
	err = db.Where(&account).First(&account).Error
	if err == nil {
		info, _ := servicesrpc.AccountInfo(account.UserAccountCode)
		if info != nil && info.CurrentBalance > 0 {
			balances = append(balances, BalanceQueryResp{
				AccountName: DefaultAccount,
				AssetId:     "00",
				Balance:     float64(info.CurrentBalance),
			})
		}

		var accountBalances []models.AccountBalance
		_ = db.Where("account_id =?", account.ID).Find(&accountBalances).Error
		if len(accountBalances) != 0 {
			for _, balance := range accountBalances {
				balances = append(balances, BalanceQueryResp{
					AccountName: DefaultAccount,
					AssetId:     balance.AssetId,
					Balance:     balance.Amount,
				})
			}
		}
	}

	lockedAccount := custodyModels.LockAccount{
		UserName: quest.UserName,
	}
	err = db.Where(&lockedAccount).First(&lockedAccount).Error
	if err == nil {
		var lockedBalances []custodyModels.LockBalance
		err = db.Where("account_id =?", lockedAccount.ID).Find(&lockedBalances).Error
		if len(lockedBalances) != 0 {
			for _, balance := range lockedBalances {
				balances = append(balances, BalanceQueryResp{
					AccountName: LockedAccount,
					AssetId:     balance.AssetId,
					Balance:     balance.Amount,
				})
			}
		}
	}
	return &balances
}

type GetAssetListQuest struct {
	AssetId  string `json:"assetId"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}
type GetAssetListResp struct {
	AssetId  string  `json:"assetId" gorm:"column:asset_id"`
	UserName string  `json:"userName" gorm:"column:user_name"`
	Amount   float64 `json:"amount" gorm:"column:amount"`
}

func GetAssetList(quest GetAssetListQuest) (*[]GetAssetListResp, int64) {
	db := middleware.DB
	var assetList []GetAssetListResp
	if quest.AssetId == "" {
		return &assetList, 0
	}
	q := db.Where("asset_id =?", quest.AssetId)

	// 查询总记录数
	var total int64
	err := q.Model(&models.AccountBalance{}).Count(&total).Error
	if err != nil || total == 0 {
		return nil, 0
	}

	q.Table("user_account_balance").
		Joins("LEFT JOIN user_account ON user_account.id = user_account_balance.account_id").
		Limit(quest.PageSize).Offset((quest.Page) * quest.PageSize).
		Select("user_account_balance.*, user_account.user_name").
		Order("user_account_balance.amount DESC").
		Scan(&assetList)

	return &assetList, total
}
