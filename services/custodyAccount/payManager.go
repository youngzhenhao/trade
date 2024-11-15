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
	"trade/services/custodyAccount/custodyAssets"
	cBase "trade/services/custodyAccount/custodyBase"
	"trade/services/custodyAccount/custodyBtc"
	"trade/services/custodyAccount/lockPayment"
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
	// Start the custody account service
	custodyBtc.BtcSever.Start(ctx)
	custodyBtc.InvoiceServer.Start(ctx)
	custodyAssets.OutsideSever.Start(ctx)
	custodyAssets.InSideSever.Start(ctx)
	custodyAssets.AddressServer.Start(ctx)
	//Check the custody service status

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

func BackAmount(payInsideId uint) (uint, error) {
	missionId, err := CreateBackFeeMission(payInsideId)
	if err != nil {
		return 0, err
	}
	return missionId, nil
}

func CheckBackFeeMission(missionId uint) bool {
	return checkBackFeeMissionById(missionId)
}

// CheckPayInsideStatus
// 检查内部转账任务状态是否成功
func CheckPayInsideStatus(id uint) (bool, error) {
	p, err := btldb.ReadPayInside(id)
	if err != nil {
		return false, err
	}
	switch p.Status {
	case models.PayInsideStatusSuccess:
		return true, nil
	case models.PayInsideStatusFailed:
		return false, models.CustodyAccountPayInsideMissionFaild
	default:
		return false, models.CustodyAccountPayInsideMissionPending
	}
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
