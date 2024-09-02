package btc_channel

import (
	"errors"
	"trade/btlLog"
	"trade/models"
	"trade/services/btldb"
	rpc "trade/services/servicesrpc"
)

func PutInAward(account *models.Account, _ string, amount int, memo *string) error {
	var err error
	acc, err := rpc.AccountInfo(account.UserAccountCode)
	if err != nil {
		return err
	}
	if amount < 0 {
		return errors.New("award amount is error")
	}
	newBalance := acc.CurrentBalance + int64(amount)
	// Change the escrow account balance
	_, err = rpc.AccountUpdate(account.UserAccountCode, newBalance, -1)
	// Build a database storage object
	ba := models.Balance{}
	ba.AccountId = account.ID
	ba.Amount = float64(amount)
	ba.Unit = models.UNIT_SATOSHIS
	ba.BillType = models.BillTypeAwardSat
	ba.Away = models.AWAY_IN
	if err != nil {
		ba.State = models.STATE_FAILED
	} else {
		ba.State = models.STATE_SUCCESS
	}
	invoiceType := "award"
	ba.Invoice = nil
	ba.PaymentHash = nil
	ba.ServerFee = 0
	ba.Invoice = &invoiceType
	// Update the database
	dbErr := btldb.CreateBalance(&ba)
	if dbErr != nil {
		btlLog.CUST.Error(dbErr.Error())
	}
	err = btldb.CreateAward(&models.AccountAward{
		AccountID: account.ID,
		AssetId:   "00",
		Amount:    float64(amount),
		Memo:      memo,
	})
	if err != nil {
		btlLog.CUST.Error(err.Error())
	}
	return nil
}
