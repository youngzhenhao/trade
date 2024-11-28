package services

import (
	"math"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
	"trade/utils"
)

func ReadUserAccountBalancesByAssetId(assetId string) (*[]custodyModels.AccountBalance, error) {
	var accountBalances []custodyModels.AccountBalance
	err := middleware.DB.Where("amount <> ? AND asset_id = ?", 0, assetId).Order("amount desc").Find(&accountBalances).Error
	return &accountBalances, err
}

func ReadAccountBalancesByAssetId(assetId string) (*[]custodyModels.AccountBalance, error) {
	var accountBalances []custodyModels.AccountBalance
	err := middleware.DB.Where("asset_id = ?", assetId).Find(&accountBalances).Error
	return &accountBalances, err
}

func ReadUserAccountBalancesByAssetIdLimitAndOffset(assetId string, limit int, offset int) (*[]custodyModels.AccountBalance, error) {
	var accountBalances []custodyModels.AccountBalance
	err := middleware.DB.Where("amount <> ? AND asset_id = ?", 0, assetId).Order("amount desc").Limit(limit).Offset(offset).Find(&accountBalances).Error
	return &accountBalances, err
}

func ReadAllAccountAssetBalancesByAssetId(assetId string) (*[]custodyModels.AccountBalance, error) {
	return ReadUserAccountBalancesByAssetId(assetId)
}

func GetAllAccountAssetBalancesByAssetId(assetId string) (*[]custodyModels.AccountBalance, error) {
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
		return nil, utils.AppendErrorInfo(err, "ReadUserAccountAccountId")
	}
	return &UserIdAndUsername{
		UserId:   int(account.UserId),
		Username: account.UserName,
	}, nil
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
		return nil, utils.AppendErrorInfo(err, "ReadUserAccountBalancesByAssetId")
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
		return nil, utils.AppendErrorInfo(err, "ReadUserAccountBalancesByAssetIdLimitAndOffset")
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

func GetAccountAssetBalanceLength(assetId string) (int64, error) {
	var count int64
	err := middleware.DB.
		Model(&custodyModels.AccountBalance{}).
		Where("amount <> ? AND asset_id = ?", 0, assetId).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetAccountAssetBalancePageNumberByPageSize(assetId string, pageSize int) (int, error) {
	recordsNum, err := GetAccountAssetBalanceLength(assetId)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "GetAccountAssetBalanceLength")
	}
	return int(math.Ceil(float64(recordsNum) / float64(pageSize))), nil
}

func GetAccountAssetBalanceUserHoldTotalAmount(assetId string) (int, error) {
	response, err := ReadAccountBalancesByAssetId(assetId)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "ReadAccountBalancesByAssetId")
	}
	if response == nil || len(*(response)) == 0 {
		return 0, nil
	}
	var totalAmount int
	for _, accountBalance := range *response {
		totalAmount += int(math.Floor(accountBalance.Amount))
	}
	return totalAmount, nil
}
