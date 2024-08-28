package account

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"sync"
	"trade/btlLog"
	"trade/config"
	"trade/models"
	"trade/services/btldb"
	rpc "trade/services/servicesrpc"
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
	// Create a custody account based on user information
	account, macaroon, err := rpc.AccountCreate(0, 0)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, err
	}
	// Build a macaroon storage path
	macaroonDir := config.GetConfig().ApiConfig.CustodyAccount.MacaroonDir
	if _, err = os.Stat(macaroonDir); os.IsNotExist(err) {
		err = os.MkdirAll(macaroonDir, os.ModePerm)
		if err != nil {
			btlLog.CUST.Error(fmt.Sprintf("创建目标文件夹 %s 失败: %v\n", macaroonDir, err))
			return nil, err
		}
	}
	macaroonFile := filepath.Join(macaroonDir, account.Id+".macaroon")
	// Store macaroon information
	err = saveMacaroon(macaroon, macaroonFile)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, err
	}
	// Build an account object
	var accountModel models.Account
	accountModel.UserName = user.Username
	accountModel.UserId = user.ID
	accountModel.UserAccountCode = account.Id
	accountModel.Label = &account.Label
	accountModel.Status = models.AccountStatusEnable
	// Write to the database
	CMutex.Lock()
	defer CMutex.Unlock()
	err = btldb.CreateAccount(&accountModel)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, err
	}
	// Return to the escrow account information
	return &accountModel, nil
}

// saveMacaroon 保存macaroon字节切片到指定文件
func saveMacaroon(macaroon []byte, macaroonFile string) error {
	file, err := os.OpenFile(macaroonFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	// 将字节切片写入指定位置
	data := macaroon
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// GetAccountByUserName 获取托管账户信息
func GetAccountByUserName(username string) (*models.Account, error) {
	return btldb.ReadAccountByName(username)
}

type UserInfo struct {
	User    *models.User
	Account *models.Account
}

// GetUserInfo 获取用户信息
func GetUserInfo(username string) (*UserInfo, error) {
	user, err := btldb.ReadUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", models.ReadDbErr, err)
	}
	account, err := GetAccountByUserName(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", models.ReadDbErr, err)
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果账户不存在，则创建账户
		newAccount, err := CreateAccount(user)
		if err != nil {
			return nil, CustodyAccountCreateErr
		}
		return &UserInfo{
			User:    user,
			Account: newAccount,
		}, nil
	}
	return &UserInfo{
		User:    user,
		Account: account,
	}, nil
}

// GetUserInfoById 获取用户信息
func GetUserInfoById(userId uint) (*UserInfo, error) {
	user, err := btldb.ReadUser(userId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", models.ReadDbErr, err)
	}
	account, err := GetAccountByUserName(user.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", models.ReadDbErr, err)
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果账户不存在，则创建账户
		newAccount, err := CreateAccount(user)
		if err != nil {
			return nil, CustodyAccountCreateErr
		}
		return &UserInfo{
			User:    user,
			Account: newAccount,
		}, nil
	}
	return &UserInfo{
		User:    user,
		Account: account,
	}, nil
}
