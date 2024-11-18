package btldb

import (
	"sync"
	"trade/middleware"
	"trade/models/custodyModels"
)

var custodyAssetMutex sync.Mutex

// CreateAccountBalance creates a new AccountBalance record
func CreateAccountBalance(accountBalance *custodyModels.AccountBalance) error {
	custodyAssetMutex.Lock()
	defer custodyAssetMutex.Unlock()
	return middleware.DB.Create(accountBalance).Error
}

// GetAccountBalance GetInvoice AccountBalance an invoice by Id
func GetAccountBalance(id uint) (*custodyModels.AccountBalance, error) {
	var accountBalance custodyModels.AccountBalance
	err := middleware.DB.First(&accountBalance, id).Error
	return &accountBalance, err
}

// GetAccountBalanceByAccountId  retrieves an accountBalances by AccountId
func GetAccountBalanceByAccountId(AccountID uint) (*[]custodyModels.AccountBalance, error) {
	var accountBalances []custodyModels.AccountBalance
	err := middleware.DB.Where("account_Id =?", AccountID).Find(&accountBalances).Error
	return &accountBalances, err
}

// GetAccountBalanceByGroup  retrieves an accountBalances by AccountId and AssetId
func GetAccountBalanceByGroup(AccountID uint, AssetID string) (*custodyModels.AccountBalance, error) {
	var accountBalance custodyModels.AccountBalance
	err := middleware.DB.Where("account_Id =? and asset_Id =?", AccountID, AssetID).First(&accountBalance).Error
	return &accountBalance, err
}

// UpdateAccountBalance updates an existing UpdateAccountBalance
func UpdateAccountBalance(accountBalance *custodyModels.AccountBalance) error {
	custodyAssetMutex.Lock()
	defer custodyAssetMutex.Unlock()
	return middleware.DB.Save(accountBalance).Error
}

// DeleteAccountBalance soft deletes an invoice by I'd
func DeleteAccountBalance(id uint) error {
	custodyAssetMutex.Lock()
	defer custodyAssetMutex.Unlock()
	var accountBalance custodyModels.AccountBalance
	return middleware.DB.Delete(&accountBalance, id).Error
}
