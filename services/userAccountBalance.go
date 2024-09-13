package services

import (
	"math"
	"trade/middleware"
	"trade/models"
)

func ReadUserAccountBalancesByAssetId(assetId string) (*[]models.AccountBalance, error) {
	var accountBalances []models.AccountBalance
	err := middleware.DB.Where("amount <> ? AND asset_id = ?", 0, assetId).Order("amount desc").Find(&accountBalances).Error
	return &accountBalances, err
}

func ReadUserAccountBalancesByAssetIdLimitAndOffset(assetId string, limit int, offset int) (*[]models.AccountBalance, error) {
	var accountBalances []models.AccountBalance
	err := middleware.DB.Where("amount <> ? AND asset_id = ?", 0, assetId).Order("amount desc").Limit(limit).Offset(offset).Find(&accountBalances).Error
	return &accountBalances, err
}

func ReadAllAccountAssetBalancesByAssetId(assetId string) (*[]models.AccountBalance, error) {
	return ReadUserAccountBalancesByAssetId(assetId)
}

func GetAllAccountAssetBalancesByAssetId(assetId string) (*[]models.AccountBalance, error) {
	return ReadAllAccountAssetBalancesByAssetId(assetId)
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

type GetAccountAssetBalanceLimitAndOffsetRequest struct {
	AssetId string `json:"asset_id"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}

type GetAccountAssetBalancePageNumberByPageSizeRequest struct {
	AssetId  string `json:"asset_id"`
	PageSize int    `json:"page_size"`
}

func GetAccountAssetBalanceExtendsLimitAndOffset(assetId string, limit int, offset int) (*[]AccountAssetBalanceExtend, error) {
	var accountAssetBalanceExtends []AccountAssetBalanceExtend
	accountBalances, err := ReadUserAccountBalancesByAssetIdLimitAndOffset(assetId, limit, offset)
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

func GetAccountAssetBalanceLength(assetId string) (int, error) {
	response, err := GetAllAccountAssetBalancesByAssetId(assetId)
	if err != nil {
		return 0, err
	}
	if response == nil || len(*(response)) == 0 {
		return 0, nil
	}
	return len(*response), nil
}

func GetAccountAssetBalancePageNumberByPageSize(assetId string, pageSize int) (int, error) {
	recordsNum, err := GetAccountAssetBalanceLength(assetId)
	if err != nil {
		return 0, err
	}
	return int(math.Ceil(float64(recordsNum) / float64(pageSize))), nil
}
