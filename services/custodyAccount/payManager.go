package custodyAccount

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"trade/btlLog"
	"trade/config"
	"trade/models"
	"trade/services/btldb"
	"trade/services/custodyAccount/btc_channel"
	"trade/services/custodyAccount/custodyAssets"
)

var (
	AdminUserId    uint = 1
	AdminAccount   *models.Account
	AdminAccountId uint = 1
)

func CustodyStart(ctx context.Context, cfg *config.Config) bool {
	// Check the admin account
	if !checkAdminAccount() {
		log.Println("Admin account is not set")
		return false
	}
	// Check the custody account MacaroonDir
	if cfg.ApiConfig.CustodyAccount.MacaroonDir == "" {
		log.Println("Custody account MacaroonDir is not set")
		return false
	}
	fmt.Println("Custody account MacaroonDir is set:", cfg.ApiConfig.CustodyAccount.MacaroonDir)
	// Start the custody account service
	btc_channel.BtcSever.Start(ctx)
	btc_channel.InvoiceServer.Start(ctx)
	custodyAssets.OutsideSever.Start(ctx)
	custodyAssets.InSideSever.Start()
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
		adminUser.Status = 1
		err = btldb.CreateUser(adminUser)
		if err != nil {
			btlLog.CUST.Error("create AdminUser failed:%s", err)
			return false
		}
	}
	adminAccount, err := btldb.ReadAccountByUserId(adminUser.ID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			btlLog.CUST.Error("CheckAdminAccount failed:%s", err)
			return false
		}
		// 创建管理员ACCOUNT
		adminAccount.UserId = adminUser.ID
		adminAccount.UserName = adminUser.Username
		adminAccount.UserAccountCode = "admin"
		adminAccount.Status = models.AccountStatusEnable
		err = btldb.CreateAccount(adminAccount)
		if err != nil {
			btlLog.CUST.Error("create AdminAccount failed:%s", err)
			return false
		}
	}
	AdminUserId = adminUser.ID
	AdminAccountId = adminAccount.ID
	AdminAccount = adminAccount
	btlLog.CUST.Info("admin user id:%d", AdminUserId)
	btlLog.CUST.Info("admin account id:%d", AdminAccount.ID)
	return true
}

// 托管账户划扣费用
func PayAmountToAdmin(payUserId uint, gasFee uint64) (uint, error) {
	e, err := btc_channel.NewBtcChannelEventByUserId(payUserId)
	if err != nil {
		btlLog.CUST.Error("PayAmountToAdmin failed:%s", err)
		return 0, err
	}
	id, err := btc_channel.PayFirLunchFee(e, gasFee)
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

// CheckPayInsideStatus 检查内部转账任务状态是否成功
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

// IsAccountBalanceEnoughByUserId  判断账户余额是否足够
func IsAccountBalanceEnoughByUserId(userId uint, value uint64) bool {
	e, err := btc_channel.NewBtcChannelEventByUserId(userId)
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

type ApplyRequest struct {
	Amount int64  `json:"amount"`
	Memo   string `json:"memo"`
}
type PayInvoiceRequest struct {
	Invoice  string `json:"invoice"`
	FeeLimit int64  `json:"feeLimit"`
}
type PaymentRequest struct {
	AssetId string `json:"asset_id"`
}
type DecodeInvoiceRequest struct {
	Invoice string `json:"invoice"`
}
