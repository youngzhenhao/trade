package custodyAccount

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"trade/btlLog"
	"trade/config"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
	"trade/services/btldb"
	"trade/services/custodyAccount/account"
	cBase "trade/services/custodyAccount/custodyBase"
	"trade/services/custodyAccount/defaultAccount/custodyAssets"
	"trade/services/custodyAccount/defaultAccount/custodyBtc"
	"trade/services/custodyAccount/lockPayment"
	"trade/services/servicesrpc"
)

var (
	AdminUserInfo *account.UserInfo
)

type ApplyRequest struct {
	Amount int64  `json:"amount"`
	Memo   string `json:"memo"`
}

type PayInvoiceRequest struct {
	Invoice  string `json:"invoice"`
	FeeLimit int64  `json:"feeLimit"`
}

type PaymentRequest struct {
	AssetId  string `json:"asset_id"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type DecodeInvoiceRequest struct {
	Invoice string `json:"invoice"`
}

func CustodyStart(ctx context.Context, cfg *config.Config) bool {
	// Check the admin account
	if !checkAdminAccount() {
		btlLog.CUST.Error("Admin account is not set")
		return false
	}
	// Check the custody account MacaroonDir
	if cfg.ApiConfig.CustodyAccount.MacaroonDir == "" {
		log.Println("Custody account MacaroonDir is not set")
		return false
	}
	fmt.Println("Custody account MacaroonDir is set:", cfg.ApiConfig.CustodyAccount.MacaroonDir)
	{
		//收款地址监听
		custodyBtc.InvoiceServer.Start(ctx)
		// 加载pending mission
		custodyBtc.LoadAOMMission()
		custodyBtc.LoadAIMMission()
	}

	{
		//收款地址监听
		custodyAssets.AddressServer.Start(ctx)
		//asset 转账监听
		//custodyAssets.OutsideSever.Start(ctx)
		// 加载pending mission
		custodyAssets.LoadAIMMission()
	}

	return true
}

func checkAdminAccount() bool {
	adminUser, err := btldb.ReadUserByUsername("admin")
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			btlLog.CUST.Error("CheckAdminAccount failed:%s", err)
			return false
		}
		// 创建管理员USER
		adminUser.Username = "admin"
		adminUser.Password = "admin"
		err = btldb.CreateUser(adminUser)
		if err != nil {
			btlLog.CUST.Error("create AdminUser failed:%s", err)
			return false
		}
	}
	//err = AutoMargeBalance()
	//if err != nil {
	//	btlLog.CUST.Error("AutoMargeBalance failed:%s", err)
	//	return false
	//}
	adminAccount, err := account.GetUserInfo("admin")
	if err != nil {
		btlLog.CUST.Error("CheckAdminAccount failed:%s", err)
		return false
	}
	if adminAccount.Account.UserAccountCode == "admin" {
		btlLog.CUST.Error("admin user is old : admin")
		return false
	}
	AdminUserInfo = adminAccount

	btlLog.CUST.Info("admin user id:%d", AdminUserInfo.User.ID)
	btlLog.CUST.Info("admin account id:%d", AdminUserInfo.Account.ID)
	btlLog.CUST.Info("admin account lit id:%s", AdminUserInfo.Account.UserAccountCode)
	btlLog.CUST.Info("admin lockAccount id:%d", AdminUserInfo.LockAccount.ID)
	return true
}

// PayAmountToAdmin
// 托管账户划扣费用
func PayAmountToAdmin(payUserId uint, gasFee uint64) (uint, error) {
	e, err := custodyBtc.NewBtcChannelEventByUserId(payUserId)
	if err != nil {
		btlLog.CUST.Error("PayAmountToAdmin failed:%s", err)
		return 0, err
	}
	id, err := custodyBtc.PayFirLunchFee(e, gasFee)
	if err != nil {
		btlLog.CUST.Error("PayAmountToAdmin failed:%s", err)
		return 0, err
	}
	return id, nil
}

// CheckPayInsideStatus
// 检查内部转账任务状态是否成功
func CheckPayInsideStatus(id uint) (bool, error) {
	return custodyBtc.CheckFirLunchFee(id)
}

func BackAmount(payInsideId uint) (uint, error) {
	return 0, fmt.Errorf("not support")
}

func CheckBackFeeMission(missionId uint) bool {
	return false
}

// IsAccountBalanceEnoughByUserId
// 判断账户余额是否足够
func IsAccountBalanceEnoughByUserId(userId uint, value uint64) bool {
	e, err := custodyBtc.NewBtcChannelEventByUserId(userId)
	if err != nil {
		btlLog.CUST.Error("PayAmountToAdmin failed:%s", err)
		return false
	}
	balance, err := e.GetBalance()
	if err != nil {
		return false
	}

	return balance[0].Amount >= int64(value)
}

// GetAccountBalance
// @Description: Get account balance
func GetAccountBalance(userId uint) (int64, error) {
	e, err := custodyBtc.NewBtcChannelEventByUserId(userId)
	if err != nil {
		return 0, err
	}
	balance, err := e.GetBalance()
	if err != nil {
		return 0, err
	}
	return balance[0].Amount, nil
}

func LockPaymentToPaymentList(usr *account.UserInfo, assetId string, pageNum, pageSize, away int) (*cBase.PaymentList, error) {
	btc, err := lockPayment.ListTransferBTC(usr, assetId, pageNum, pageSize, away)
	if err != nil {
		return nil, err
	}
	db := middleware.DB
	var list cBase.PaymentList
	for i := range btc {
		v := btc[i]
		r := cBase.PaymentResponse{}
		r.Timestamp = v.CreatedAt.Unix()

		switch v.BillType {
		case custodyModels.LockBillTypeLock:
			r.Away = models.AWAY_IN
		case custodyModels.LockBillTypeAward:
			r.Away = models.AWAY_IN
			var awardExt models.AccountAwardExt
			db.Where("balance_id =? and account_type =1", v.ID).First(&awardExt)
			var award models.AccountAward
			db.Where("id =?", awardExt.AwardId).First(&award)
			v.LockId = cBase.GetAwardType(*award.Memo)

		default:
			r.Away = models.AWAY_OUT
		}
		r.BillType = models.LockedTransfer
		r.Invoice = &v.LockId
		r.Address = &v.LockId
		r.Target = &v.LockId
		r.PaymentHash = &v.LockId
		r.Amount = v.Amount
		r.AssetId = &v.AssetId
		r.State = models.STATE_SUCCESS
		r.Fee = 0
		list.PaymentList = append(list.PaymentList, r)
	}
	return &list, nil
}

func AutoMargeBalance() error {
	accounts, err := servicesrpc.ListAccounts()
	if err != nil {
		return err
	}
	for _, acc := range accounts {
		if acc.CurrentBalance > 100 {
			db := middleware.DB
			var a models.Account
			err := db.Where("user_account_code =?", acc.Id).First(&a).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			var balance custodyModels.AccountBtcBalance
			err = db.Where("account_id =?", a.ID).First(&balance).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				balance.AccountId = a.ID
				balance.Amount = float64(acc.CurrentBalance)
				db.Create(&balance)
			}
		}
	}
	return nil
}
