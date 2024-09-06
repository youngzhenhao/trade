package assets

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"trade/btlLog"
	"trade/middleware"
	"trade/models"
	"trade/services/btldb"
)

func PutInAward(account *models.Account, AssetId string, amount int, memo *string) error {
	var in models.AwardInventory
	err := middleware.DB.Where("asset_Id =? ", AssetId).First(&in).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		btlLog.CUST.Error("err:%v", models.ReadDbErr)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("no award type")
	}
	if in.Status != models.AwardInventoryAble {
		return fmt.Errorf("award is lock")
	}
	if in.Amount < float64(amount) {
		return fmt.Errorf("not enough award")
	}

	receiveBalance, err := btldb.GetAccountBalanceByGroup(account.ID, AssetId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		btlLog.CUST.Error("err:%v", models.ReadDbErr)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		r := models.AccountBalance{
			AccountID: account.ID,
			AssetId:   AssetId,
			Amount:    float64(amount),
		}
		err = btldb.UpdateAccountBalance(&r)
		if err != nil {
			btlLog.CUST.Error("err:%v", models.ReadDbErr)
		}
	} else {
		receiveBalance.Amount += float64(amount)
		err = btldb.UpdateAccountBalance(receiveBalance)
		if err != nil {
			btlLog.CUST.Error("err:%v", models.ReadDbErr)
		}
	}

	// Build a database storage object
	ba := models.Balance{}
	ba.AccountId = account.ID
	ba.Amount = float64(amount)
	ba.Unit = models.UNIT_ASSET_NORMAL
	ba.BillType = models.BillTypeAwardAsset
	ba.Away = models.AWAY_IN
	ba.AssetId = &AssetId
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
		AssetId:   AssetId,
		Amount:    float64(amount),
		Memo:      memo,
	})
	if err != nil {
		btlLog.CUST.Error(err.Error())
	}
	return nil
}
