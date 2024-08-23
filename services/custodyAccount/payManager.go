package custodyAccount

import (
	"errors"
	"gorm.io/gorm"
	"trade/btlLog"
	"trade/models"
	"trade/services/btldb"
)

var (
	AdminUserId    uint = 1
	AdminAccount   *models.Account
	AdminAccountId uint = 1
)

// 托管账户划扣费用
func PayAmountToAdmin(payUserId uint, gasFee uint64) (uint, error) {
	id, err := CreatePayInsideMission(payUserId, AdminUserId, gasFee, 0, "00")
	if err != nil {
		btlLog.CUST.Error("PayAmountToAdmin failed:%s", err)
		return 0, err
	}
	return id, nil
}
func BackAmount(payUserId uint) (uint, error) {
	missionId, err := CreateBackFeeMission(payUserId)
	if err != nil {
		return 0, err
	}
	return missionId, nil
}
func CheckBackFeeMission(missionId uint) bool {
	return checkBackFeeMissionById(missionId)
}

func CheckAdminAccount() bool {
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

	adminAccount, err := btldb.ReadAccount(adminUser.ID)

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
	return true
}
