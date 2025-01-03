package services

import (
	"errors"
	"math"
	"strings"
	"trade/middleware"
	"trade/models"
)

type AccountAssetTransfer struct {
	BillBalanceId int    `json:"bill_balance_id"`
	AccountId     int    `json:"account_id"`
	Username      string `json:"username"`
	BillType      string `json:"bill_type"`
	Away          string `json:"away"`
	Amount        int    `json:"amount"`
	ServerFee     int    `json:"server_fee"`
	AssetId       string `json:"asset_id"`
	Invoice       string `json:"invoice"`
	Outpoint      string `json:"outpoint"`
	Time          int    `json:"time"`
}

func BillBalanceToAccountAssetTransfer(billBalance *models.Balance, username string) *AccountAssetTransfer {
	if billBalance == nil {
		return nil
	}
	var assetId string
	if billBalance.AssetId != nil {
		assetId = *billBalance.AssetId
	}
	var invoice string
	if billBalance.Invoice != nil {
		invoice = *billBalance.Invoice
	}
	var outpoint string
	if billBalance.PaymentHash != nil && *billBalance.PaymentHash != "" && strings.Contains(*billBalance.PaymentHash, ":") {
		outpoint = *billBalance.PaymentHash
	}
	return &AccountAssetTransfer{
		BillBalanceId: int(billBalance.ID),
		AccountId:     int(billBalance.AccountId),
		Username:      username,
		BillType:      billBalance.BillType.String(),
		Away:          billBalance.Away.String(),
		Amount:        int(billBalance.Amount),
		ServerFee:     int(billBalance.ServerFee),
		AssetId:       assetId,
		Invoice:       invoice,
		Outpoint:      outpoint,
		Time:          int(billBalance.CreatedAt.Unix()),
	}
}

func BillBalancesToAccountAssetTransfers(billBalances *[]models.Balance) *[]AccountAssetTransfer {
	if billBalances == nil {
		return nil
	}
	var accountAssetTransfer []AccountAssetTransfer
	for _, billBalance := range *billBalances {
		var usernameByAccountId string
		userIdAndUsername, err := GetUserIdAndUsernameByAccountId(billBalance.AccountId)
		if err != nil {
			continue
		} else {
			usernameByAccountId = userIdAndUsername.Username
		}
		accountAssetTransfer = append(accountAssetTransfer, *BillBalanceToAccountAssetTransfer(&billBalance, usernameByAccountId))
	}
	return &accountAssetTransfer
}

// @dev: Use this now
func GetAllAccountAssetTransfersByBillBalanceAssetTransferAndAwardAsset(assetId string) (*[]AccountAssetTransfer, error) {
	billBalances, err := ReadBillBalanceAssetTransferAndAwardAssetByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	accountAssetTransfers := BillBalancesToAccountAssetTransfers(billBalances)
	return accountAssetTransfers, nil
}

func GetAccountAssetTransfersByBillBalanceAssetTransferAndAwardAssetLimitAndOffset(assetId string, limit int, offset int) (*[]AccountAssetTransfer, error) {
	billBalances, err := ReadBillBalanceAssetTransferAndAwardAssetByAssetIdLimitAndOffset(assetId, limit, offset)
	if err != nil {
		return nil, err
	}
	accountAssetTransfers := BillBalancesToAccountAssetTransfers(billBalances)
	return accountAssetTransfers, nil
}

func GetAllAccountAssetTransfersByBillBalanceAssetTransfer(assetId string) (*[]AccountAssetTransfer, error) {
	billBalances, err := ReadBillBalanceAssetTransferByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	accountAssetTransfers := BillBalancesToAccountAssetTransfers(billBalances)
	return accountAssetTransfers, nil
}

// GetAllAccountAssetTransfersByAssetId
// @Description: Get all account asset transfers by asset id
func GetAllAccountAssetTransfersByAssetId(assetId string) (*[]AccountAssetTransfer, error) {
	if assetId == "00" {
		return nil, errors.New("invalid asset id")
	}
	return GetAllAccountAssetTransfersByBillBalanceAssetTransferAndAwardAsset(assetId)
}

// GetAccountAssetTransfersLimitAndOffset
// @Description: Get all account asset transfers by asset id limit and offset
func GetAccountAssetTransfersLimitAndOffset(assetId string, limit int, offset int) (*[]AccountAssetTransfer, error) {
	if assetId == "00" {
		return nil, errors.New("invalid asset id")
	}
	return GetAccountAssetTransfersByBillBalanceAssetTransferAndAwardAssetLimitAndOffset(assetId, limit, offset)
}

type GetAccountAssetTransferLimitAndOffsetRequest struct {
	AssetId string `json:"asset_id"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}

type GetAccountAssetTransferPageNumberByPageSizeRequest struct {
	AssetId  string `json:"asset_id"`
	PageSize int    `json:"page_size"`
}

func GetAccountAssetTransferLength(assetId string) (int64, error) {
	var count int64
	// Query bill balance
	err := middleware.DB.
		Model(&models.Balance{}).
		Where("amount <> ? AND bill_type IN ? AND asset_id = ?", 0, []models.BalanceType{models.BillTypeAssetTransfer, models.BillTypeAwardAsset}, assetId).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetAccountAssetTransferPageNumberByPageSize(assetId string, pageSize int) (int, error) {
	recordsNum, err := GetAccountAssetTransferLength(assetId)
	if err != nil {
		return 0, err
	}
	return int(math.Ceil(float64(recordsNum) / float64(pageSize))), nil
}
