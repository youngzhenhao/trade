package services

import (
	"errors"
	"gorm.io/gorm"
)

var adminUserId uint = 1

func NewRecharge() {

}

// 托管账户划扣费用
func PayAmountToAdmin(payUserId uint, gasFee, serveFee uint64) (uint, error) {
	id, err := PayAmountInside(payUserId, adminUserId, gasFee, serveFee, "")
	if err != nil {
		CUST.Error("PayAmountToAdmin failed:%s", err)
		return 0, err
	}
	return id, nil
}

func CheckAdminAccount() bool {
	adminUser, err := ReadUserByUsername("admin")
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			CUST.Error("CheckAdminAccount failed:%s", err)
			return false
		}
		// 创建管理员USER
		adminUser.Username = "admin"
		adminUser.Password = "admin"
		adminUser.Status = 1
		err = CreateUser(adminUser)
		if err != nil {
			CUST.Error("create AdminUser failed:%s", err)
			return false
		}
	}

	adminAccount, err := ReadAccountByUserId(adminUser.ID)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			CUST.Error("CheckAdminAccount failed:%s", err)
			return false
		}
		// 创建管理员ACCOUNT
		adminAccount.UserId = adminUser.ID
		adminAccount.UserName = adminUser.Username
		adminAccount.UserAccountCode = "admin"
		adminAccount.Status = 1
		err = CreateAccount(adminAccount)
		if err != nil {
			CUST.Error("create AdminAccount failed:%s", err)
			return false
		}
	}
	adminUserId = adminUser.ID
	CUST.Info("admin user id:%d", adminUserId)
	return true
}
