package btc_channel

import (
	"errors"
	"trade/btlLog"
	"trade/middleware"
	"trade/models"
	rpc "trade/services/servicesrpc"
)

func PutInAward(account *models.Account, _ string, amount int, memo *string) (*models.AccountAward, error) {
	var err error
	tx, back := middleware.GetTx()
	defer back()
	// Build a database Balance
	ba := models.Balance{}
	ba.AccountId = account.ID
	ba.Amount = float64(amount)
	ba.Unit = models.UNIT_SATOSHIS
	ba.BillType = models.BillTypeAwardSat
	ba.Away = models.AWAY_IN
	ba.State = models.STATE_SUCCESS
	invoiceType := "award"
	ba.Invoice = nil
	ba.PaymentHash = memo
	ba.ServerFee = 0
	ba.Invoice = &invoiceType
	if err = tx.Create(&ba).Error; err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, err
	}
	// Build a database AccountAward
	award := models.AccountAward{
		AccountID: account.ID,
		AssetId:   "00",
		Amount:    float64(amount),
		Memo:      memo,
	}
	if err = tx.Create(&award).Error; err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, err
	}
	// Build a database  AccountAwardExt
	awardExt := models.AccountAwardExt{
		BalanceId: ba.ID,
		AwardId:   award.ID,
	}
	if err = tx.Create(&awardExt).Error; err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, err
	}
	// Update the escrow account balance
	acc, err := rpc.AccountInfo(account.UserAccountCode)
	if err != nil {
		return nil, err
	}
	if amount < 0 || amount > 1000000 {
		return nil, errors.New("award amount is error")
	}
	newBalance := acc.CurrentBalance + int64(amount)
	if newBalance < 0 || newBalance > 100000000 {
		return nil, errors.New("amount is error(<0 or >100000000)")
	}

	// Change the escrow account balance
	_, err = rpc.AccountUpdate(account.UserAccountCode, newBalance, -1)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, err
	}
	tx.Commit()
	return &award, nil
}
