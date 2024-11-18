package custodyAssets

import (
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"sync"
	"trade/btlLog"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
)

var (
	ServerBusy    = errors.New("seriver is busy, please try again later")
	NoAwardType   = fmt.Errorf("no award type")
	AssetIdLock   = fmt.Errorf("award is lock")
	NoEnoughAward = fmt.Errorf("not enough award")
)
var (
	AwardLock = sync.Mutex{}
)

func PutInAward(account *models.Account, AssetId string, amount int, memo *string, lockedId string) (*models.AccountAward, error) {
	tx, back := middleware.GetTx()
	if tx == nil {
		return nil, ServerBusy
	}
	defer back()
	// Check if the asset is award type
	var in models.AwardInventory
	err := tx.Where("asset_Id =? ", AssetId).First(&in).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		btlLog.CUST.Error("err:%v", err)
		return nil, ServerBusy
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, NoAwardType
	}
	if in.Status != models.AwardInventoryAble {
		return nil, AssetIdLock
	}
	if in.Amount < float64(amount) {
		return nil, NoEnoughAward
	}

	AwardLock.Lock()
	defer AwardLock.Unlock()

	// Update the award inventory
	in.Amount -= float64(amount)
	err = tx.Save(&in).Error
	if err != nil {
		btlLog.CUST.Error("err:%v", err)
		return nil, ServerBusy
	}

	// Update the account balance
	var receiveBalance custodyModels.AccountBalance
	err = tx.Where("account_Id =? and asset_Id =?", account.ID, AssetId).First(&receiveBalance).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		btlLog.CUST.Error("err:%v", err)
		return nil, ServerBusy
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		r := custodyModels.AccountBalance{
			AccountID: account.ID,
			AssetId:   AssetId,
			Amount:    float64(amount),
		}
		err = tx.Save(&r).Error
		if err != nil {
			btlLog.CUST.Error("err:%v", err)
			return nil, ServerBusy
		}
	} else {
		receiveBalance.Amount += float64(amount)
		err = tx.Save(&receiveBalance).Error
		if err != nil {
			btlLog.CUST.Error("err:%v", err)
			return nil, ServerBusy
		}
	}

	// Build a database balance
	ba := models.Balance{
		TypeExt: &models.BalanceTypeExt{
			Type: models.BTExtAward,
		},
	}
	ba.AccountId = account.ID
	ba.Amount = float64(amount)
	ba.Unit = models.UNIT_ASSET_NORMAL
	ba.BillType = models.BillTypeAwardAsset
	ba.Away = models.AWAY_IN
	ba.AssetId = &AssetId
	ba.State = models.STATE_SUCCESS
	invoiceType := "award"
	ba.Invoice = nil
	ba.PaymentHash = memo
	ba.ServerFee = 0
	ba.Invoice = &invoiceType
	// Update the database
	dbErr := tx.Create(&ba).Error
	if dbErr != nil {
		btlLog.CUST.Error(dbErr.Error())
		return nil, ServerBusy
	}
	// Build a database AccountAward
	award := models.AccountAward{
		AccountID: account.ID,
		AssetId:   AssetId,
		Amount:    float64(amount),
		Memo:      memo,
	}
	err = tx.Create(&award).Error
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, ServerBusy
	}

	//Build a database AwardIdempotent
	Idempotent := models.AccountAwardIdempotent{
		AwardId:    award.ID,
		Idempotent: lockedId,
	}
	if err = tx.Create(&Idempotent).Error; err != nil {
		var mySQLErr *mysql.MySQLError
		if errors.As(err, &mySQLErr) {
			if mySQLErr.Number == 1062 {
				return nil, errors.New("RepeatedLockId")
			}
		}
		return nil, errors.New("ServiceError")
	}
	// Build a database AccountAwardExt
	awardExt := models.AccountAwardExt{
		BalanceId: ba.ID,
		AwardId:   award.ID,
	}
	err = tx.Create(&awardExt).Error
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, ServerBusy
	}
	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		btlLog.CUST.Error("award failed,not commit:%v", err)
		return nil, ServerBusy
	}
	return &award, nil
}
