package account

import (
	"errors"
	"sync"
	"trade/btlLog"
	"trade/models"
	"trade/services/btldb"
)

type AccError error

var (
	CustodyAccountCreateErr AccError = errors.New("创建托管账户失败")
	CustodyAccountUpdateErr AccError = errors.New("更新托管账户失败")
	CustodyAccountGetErr    AccError = errors.New("获取托管账户失败")
	CustodyAccountDeleteErr AccError = errors.New("删除托管账户失败")
	MacaroonSaveErr         AccError = errors.New("保存macaroon文件失败")
	MacaroonFindErr         AccError = errors.New("找不到macaroon文件")
)

var CMutex sync.Mutex

// CreateAccount 创建托管账户并存储马卡龙文件
func CreateAccount(user *models.User) (*models.Account, error) {

	// Build an account object
	var accountModel models.Account
	accountModel.UserName = user.Username
	accountModel.UserId = user.ID
	accountModel.UserAccountCode = ""
	accountModel.Label = nil
	accountModel.Status = models.AccountStatusEnable
	// Write to the database
	CMutex.Lock()
	defer CMutex.Unlock()
	err := btldb.CreateAccount(&accountModel)
	if err != nil {
		btlLog.CACC.Error(err.Error())
		return nil, err
	}
	// Return to the escrow account information
	return &accountModel, nil
}

// GetAccountByUserName 获取托管账户信息
func GetAccountByUserName(username string) (*models.Account, error) {
	return btldb.ReadAccountByName(username)
}
