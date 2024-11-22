package btldb

import (
	"sync"
	"trade/middleware"
	"trade/models/custodyModels"
)

var custodyAssetMutex sync.Mutex

// GetAccountBalanceByGroup  retrieves an accountBalances by AccountId and AssetId
func GetAccountBalanceByGroup(AccountID uint, AssetID string) (*custodyModels.AccountBalance, error) {
	var accountBalance custodyModels.AccountBalance
	err := middleware.DB.Where("account_Id =? and asset_Id =?", AccountID, AssetID).First(&accountBalance).Error
	return &accountBalance, err
}

func UpdateAccountBalance(accountBalance *custodyModels.AccountBalance) error {
	custodyAssetMutex.Lock()
	defer custodyAssetMutex.Unlock()
	return middleware.DB.Save(accountBalance).Error
}
