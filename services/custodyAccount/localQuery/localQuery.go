package localQuery

import (
	"errors"
	"fmt"
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
	UserName      string   `json:"username"`
	Away          int      `json:"away"`
	AssetId       string   `json:"assetId"`
	Invoice       string   `json:"invoice"`
	PaymentHash   string   `json:"hash"`
	AmountMin     float64  `json:"amountMin"`
	AmountMax     float64  `json:"amountMax"`
	ServerFeeMin  uint64   `json:"feeMin"`
	ServerFeeMax  uint64   `json:"feeMax"`
	TimeStart     string   `json:"timeStart"`
	TimeEnd       string   `json:"timeEnd"`
	IncludeFailed bool     `json:"includeFailed"`
	Tags          []string `json:"tags"`
	Page          int      `json:"page"`
	PageSize      int      `json:"pageSize"`
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
	Type        uint                `gorm:"column:type" json:"type"`
}

func BillQuery(quest BillQueryQuest) (*[]BillListWithUser, int64, error) {
	billQuery := models.Balance{}
	var err error

	q := middleware.DB
	q = q.Where(&billQuery)
	if !quest.IncludeFailed {
		q = q.Where("state =?", 1)
	}

	if quest.UserName != "" {
		account := models.Account{
			UserName: quest.UserName,
		}
		err = middleware.DB.Where(&account).First(&account).Error
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
	if quest.ServerFeeMin != 0 {
		q = q.Where("bill_balance.server_fee >=?", quest.ServerFeeMin)
	}
	if quest.ServerFeeMax != 0 {
		q = q.Where("bill_balance.server_fee <=?", quest.ServerFeeMax)
	}
	if quest.Invoice != "" {
		q = q.Where("bill_balance.invoice =?", quest.Invoice)
	}
	if quest.PaymentHash != "" {
		q = q.Where("bill_balance.payment_hash =?", quest.PaymentHash)
	}
	if quest.TimeStart != "" {
		q = q.Where("bill_balance.created_at >=?", quest.TimeStart)
	}
	if quest.TimeEnd != "" {
		q = q.Where("bill_balance.created_at <=?", quest.TimeEnd)
	}
	if quest.PageSize == 0 {
		quest.PageSize = 500
	}
	if quest.Tags != nil && len(quest.Tags) != 0 {
		for _, tag := range quest.Tags {
			q = q.Where("bill_balance_type_ext.type =?", models.ToBalanceTypeExtList(tag))
		}
	}
	q = q.Table("bill_balance").
		Joins("LEFT JOIN user_account ON user_account.id = bill_balance.account_id").
		Joins("LEFT JOIN bill_balance_type_ext ON bill_balance.id = bill_balance_type_ext.balance_id")

	// 查询总记录数
	var total int64
	err = q.Model(&models.Balance{}).Count(&total).Error
	if err != nil || total == 0 {
		return nil, 0, err
	}
	var billListWithUser []BillListWithUser
	err = q.Limit(quest.PageSize).Offset((quest.Page) * quest.PageSize).
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

type TotalBillListQuest struct {
	AssetId   string `json:"assetId"`
	TimeStart string `json:"timeStart"`
	TimeEnd   string `json:"timeEnd"`
	OderBy    uint   `json:"orderBy"`
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
}
type TotalBillListResp struct {
	UserName       string  `json:"userName" gorm:"column:user_name"`
	AssetId        string  `json:"assetId" gorm:"column:asset_id"`
	SumAwayEnter   float64 `json:"sumAwayEnter" gorm:"column:sum_away_enter"`
	CountAwayEnter int     `json:"countAwayEnter" gorm:"column:count_away_enter"`
	SumAwayOut     float64 `json:"sumAwayOut" gorm:"column:sum_away_out"`
	CountAwayOut   int     `json:"countAwayOut" gorm:"column:count_away_out"`
	NetIncome      float64 `json:"netIncome" gorm:"column:netIncome"`
}

func TotalBillList(quest *TotalBillListQuest) ([]TotalBillListResp, int64, error) {
	db := middleware.DB
	var err error
	q := db.Select("user_account.user_name," +
		"asset_id," +
		"SUM(CASE WHEN away = 0 THEN amount ELSE 0 END) AS sum_away_enter," +
		"count(CASE WHEN away = 0 THEN amount ELSE 0 END) as count_away_enter," +
		"SUM(CASE WHEN away = 1 THEN amount ELSE 0 END) AS sum_away_out," +
		"count(CASE WHEN away = 1 THEN amount ELSE 0 END) as count_away_out," +
		"SUM(CASE WHEN away = 0 THEN amount ELSE 0 END) - SUM(CASE WHEN away = 1 THEN amount ELSE 0 END) AS netIncome")
	q = q.Table("bill_balance")
	q = q.Joins("left JOIN  user_account on bill_balance.account_id = user_account.id")
	q.Where("bill_balance.state = ?", 1)

	if quest.TimeStart != "" {
		q = q.Where("bill_balance.created_at >=?", quest.TimeStart)
	}
	if quest.TimeEnd != "" {
		q = q.Where("bill_balance.created_at <=?", quest.TimeEnd)
	}
	if quest.AssetId == "" {
		return nil, 0, errors.New("must have assetId")
	}
	q.Where("bill_balance.asset_id = ?", quest.AssetId)
	q.Group("account_id,asset_id")

	var oder string
	switch quest.OderBy {
	case 0:
		oder = "sum_away_enter"
	case 1:
		oder = "count_away_enter"
	case 2:
		oder = "sum_away_out"
	case 3:
		oder = "count_away_out"
	case 4:
		oder = "netIncome"
	default:
		oder = "sum_away_enter"
	}
	var count int64
	err = q.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	q.Order(fmt.Sprintf("ABS(%s) desc", oder))
	q.Limit(quest.PageSize).Offset((quest.Page) * quest.PageSize)
	var total []TotalBillListResp
	err = q.Scan(&total).Error
	if err != nil {
		return total, 0, err
	}
	return total, count, nil
}
