package services

import (
	"math"
	"trade/middleware"
	"trade/models"
)

func ReadUserAccountBalancesByAssetId(assetId string) (*[]models.AccountBalance, error) {
	var accountBalances []models.AccountBalance
	err := middleware.DB.Where("amount <> ? AND asset_id = ?", 0, assetId).Find(&accountBalances).Error
	return &accountBalances, err
}

type UserIdAndUsername struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
}

func ReadUserAccountAccountId(accountId uint) (*models.Account, error) {
	var account models.Account
	err := middleware.DB.First(&account, accountId).Error
	return &account, err
}

func GetUserIdAndUsernameByAccountId(accountId uint) (*UserIdAndUsername, error) {
	account, err := ReadUserAccountAccountId(accountId)
	if err != nil {
		return nil, err
	}
	return &UserIdAndUsername{
		UserId:   int(account.UserId),
		Username: account.UserName,
	}, err
}

type AccountAssetBalanceExtend struct {
	AccountID uint   ` json:"account_id"`
	AssetId   string ` json:"asset_id"`
	Amount    int    ` json:"amount"`
	UserID    int    ` json:"user_id"`
	Username  string ` json:"username"`
}

func GetAccountAssetBalanceExtendsByAssetId(assetId string) (*[]AccountAssetBalanceExtend, error) {
	var accountAssetBalanceExtends []AccountAssetBalanceExtend
	accountBalances, err := ReadUserAccountBalancesByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	accountIdMapUserIdAndUsername := make(map[uint]*UserIdAndUsername)
	for _, accountBalance := range *accountBalances {
		accountId := accountBalance.AccountID
		userIdAndUsername, ok := accountIdMapUserIdAndUsername[accountId]
		if !ok {
			userIdAndUsername, err = GetUserIdAndUsernameByAccountId(accountBalance.AccountID)
			if err != nil {
				continue
			}
			accountIdMapUserIdAndUsername[accountBalance.AccountID] = userIdAndUsername
			accountAssetBalanceExtends = append(accountAssetBalanceExtends, AccountAssetBalanceExtend{
				AccountID: accountId,
				AssetId:   accountBalance.AssetId,
				Amount:    int(math.Floor(accountBalance.Amount)),
				UserID:    userIdAndUsername.UserId,
				Username:  userIdAndUsername.Username,
			})
		}
	}
	return &accountAssetBalanceExtends, nil
}
