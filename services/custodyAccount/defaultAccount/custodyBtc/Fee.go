package custodyBtc

import (
	"gorm.io/gorm"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
	"trade/services/custodyAccount/account"
)

func PayFirLunchFee(e *BtcChannelEvent, gasFee uint64) (uint, error) {
	tx, back := middleware.GetTx()
	defer back()
	var err error
	assetId := BtcId
	invoice := "lnfirlunchFee"
	balance := models.Balance{
		AccountId:   e.UserInfo.Account.ID,
		BillType:    models.BillTypePayment,
		Away:        models.AWAY_OUT,
		Amount:      float64(gasFee),
		Unit:        models.UNIT_SATOSHIS,
		ServerFee:   0,
		AssetId:     &assetId,
		Invoice:     &invoice,
		PaymentHash: nil,
		State:       1,
		TypeExt: &models.BalanceTypeExt{
			Type: models.BTExtFirLaunch,
		},
	}
	if err = tx.Create(&balance).Error; err != nil {
		return 0, err
	}
	//扣除手续费
	id, err := LessBtcBalance(tx, e.UserInfo, float64(gasFee), balance.ID, custodyModels.ChangeFirLunchFee)
	if err != nil {
		return 0, err
	}
	tx.Commit()

	return id, nil
}

// CheckFirLunchFee 检查内部转账任务状态是否成功
func CheckFirLunchFee(id uint) (bool, error) {
	p := custodyModels.AccountBalanceChange{}
	err := middleware.DB.Where("id =?", id).First(&p).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func PayFee(Db *gorm.DB, usr *account.UserInfo, amount float64, balanceId uint) error {
	_, err := LessBtcBalance(Db, usr, amount, balanceId, custodyModels.ChangeTypeBtcFee)
	return err
}
