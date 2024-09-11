package services

import (
	"errors"
	"strings"
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

func GetAllAccountAssetTransfersByBillBalanceAssetTransferAndAwardAsset(assetId string) (*[]AccountAssetTransfer, error) {
	billBalances, err := ReadBillBalanceAssetTransferAndAwardAssetByAssetId(assetId)
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
