package other

import (
	"fmt"
	"gorm.io/gorm"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
	caccount "trade/services/custodyAccount/account"
	"trade/services/custodyAccount/defaultAccount/custodyBtc"
)

func PayFirLunchFee(e *custodyBtc.BtcChannelEvent, gasFee uint64) (uint, error) {
	tx, back := middleware.GetTx()
	defer back()
	var err error
	assetId := custodyBtc.BtcId
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
	id, err := custodyBtc.LessBtcBalance(tx, e.UserInfo, float64(gasFee), balance.ID, custodyModels.ChangeFirLunchFee)
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

func BackFirLunchFee(id uint) (uint, error) {
	db := middleware.DB
	//检查退费id是否合法
	var err error
	p := custodyModels.AccountBalanceChange{}
	err = db.Where("id =?", id).First(&p).Error
	if err != nil {
		return 0, fmt.Errorf("get Fee id  %d failed, err: %v", id, err)
	}
	if p.ChangeType != custodyModels.ChangeFirLunchFee {
		return 0, fmt.Errorf("Fee id  %d is not FirLunchFee", id)
	}
	//获取用户信息
	var account models.Account
	err = db.Model(models.Account{}).Where("id =?", p.AccountId).First(&account).Error
	if err != nil {
		return 0, fmt.Errorf("get user info failed, err: %v", err)
	}
	usr, err := caccount.GetUserInfo(account.UserName)
	if err != nil {
		return 0, fmt.Errorf("get user info failed, err: %v", err)
	}
	//退费
	tx, back := middleware.GetTx()
	defer back()

	assetId := custodyBtc.BtcId
	invoice := "backFee"
	balance := models.Balance{
		AccountId:   usr.Account.ID,
		BillType:    models.BILL_TYPE_BACK_FEE,
		Away:        models.AWAY_IN,
		Amount:      p.ChangeAmount,
		Unit:        models.UNIT_SATOSHIS,
		ServerFee:   0,
		AssetId:     &assetId,
		Invoice:     &invoice,
		PaymentHash: nil,
		State:       1,
		TypeExt: &models.BalanceTypeExt{
			Type: models.BTExtBackFee,
		},
	}
	if err = tx.Create(&balance).Error; err != nil {
		return 0, err
	}

	backFeeMission := models.BackFee{
		Model:         gorm.Model{},
		PayInsideId:   p.ID,
		BackBalanceId: balance.ID,
		Status:        models.BackFeeStatePaid,
	}
	if err = tx.Create(&backFeeMission).Error; err != nil {
		return 0, err
	}

	feeId, err := custodyBtc.AddBtcBalance(tx, usr, p.ChangeAmount, balance.ID, custodyModels.ChangeTypeBackFee)
	if err != nil {
		return 0, err
	}
	tx.Commit()

	return feeId, nil
}
