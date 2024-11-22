package custodyAssets

import (
	"encoding/hex"
	"fmt"
	"trade/btlLog"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
	"trade/services/custodyAccount/account"
	"trade/services/custodyAccount/custodyBase/custodyLimit"
	"trade/services/custodyAccount/defaultAccount/custodyBtc"
	rpc "trade/services/servicesrpc"
)

type invoiceInfo struct {
	Invoice string
	AssetId string
	Hash    *string
}

func RunInsideStep(usr *account.UserInfo, mission *custodyModels.AccountInsideMission) error {
	db := middleware.DB
	//获取usrInfo
	if usr == nil {
		var a models.Account
		if err := db.Where("id =?", mission.AccountId).First(&a).Error; err != nil {
			btlLog.CUST.Error("GetAccount error:%s", err)
			return err
		}
		usr, _ = account.GetUserInfo(a.UserName)
	}
	//获取发票信息
	invoice := models.Invoice{}
	if err := db.Where("id =?", mission.InvoiceId).First(&invoice).Error; err != nil {
		btlLog.CUST.Error("GetInvoice error:%s", err)
		return err
	}
	DecodeAddr, err := rpc.DecodeAddr(invoice.Invoice)
	if err != nil {
		btlLog.CUST.Error("发票解析失败", err)
		return err
	}
	id := hex.EncodeToString(DecodeAddr.AssetId)

	i := invoiceInfo{
		Invoice: invoice.Invoice,
		AssetId: id,
		Hash:    nil,
	}
	//run steps
	for {
		InsideSteps(usr, mission, i)
		custodyBtc.LogAIM(middleware.DB, mission)
		switch {
		case mission.State == custodyModels.AIMStateSuccess:
			return nil
		case mission.State == custodyModels.AIMStateDone:
			return fmt.Errorf(mission.Error)
		case mission.Retries >= 30:
			return nil
		}
	}
}
func InsideSteps(usr *account.UserInfo, mission *custodyModels.AccountInsideMission, i invoiceInfo) {
	var err error
	switch mission.State {
	case custodyModels.AIMStatePending:
		tx, back := middleware.GetTx()
		defer back()
		//创建BillBalance记录
		balance := getBillBalanceModel(usr, mission.Amount, i.AssetId, models.AWAY_OUT, i)
		if err = tx.Create(balance).Error; err != nil {
			btlLog.CUST.Error("CreateBillBalance error:%s", err)
			mission.Error = err.Error()
			mission.State = custodyModels.AIMStateDone
			return
		}
		mission.PayerBalanceId = balance.ID
		//创建扣款记录
		_, err = LessAssetBalance(tx, usr, mission.Amount, mission.PayerBalanceId, i.AssetId, custodyModels.ChangeTypeAssetPayLocal)
		if err != nil {
			btlLog.CUST.Error("LessBtcBalance error:%s", err)
			mission.Error = err.Error()
			mission.State = custodyModels.AIMStateDone
			return
		}
		//扣除手续费
		err = custodyBtc.PayFee(tx, usr, mission.Fee, mission.PayerBalanceId)
		if err != nil {
			btlLog.CUST.Error("PayFee error:%s", err)
			mission.Error = err.Error()
			mission.State = custodyModels.AIMStateDone
			return
		}
		balance.ServerFee = uint64(mission.Fee)

		//更新状态
		err = tx.Save(balance).Error
		if err != nil {
			btlLog.CUST.Error("SaveBillBalance error:%s", err)
			mission.Error = err.Error()
			mission.State = custodyModels.AIMStateDone
			return
		}
		mission.State = custodyModels.AIMStatePaid
		tx.Commit()

		go func() {
			//更新额度
			limitType := custodyModels.LimitType{
				AssetId:      i.AssetId,
				TransferType: custodyModels.LimitTransferTypeLocal,
			}
			err = custodyLimit.MinusLimit(middleware.DB, usr, &limitType, mission.Amount+mission.Fee)
			if err != nil {
				btlLog.CUST.Error("额度限制未正常更新:%s", err.Error())
				btlLog.CUST.Error("error PayInsideId:%v", mission.ID)
			}
		}()
		return

	case custodyModels.AIMStatePaid:
		//获取usrInfo
		var a models.Account
		if err = middleware.DB.Where("id =?", mission.ReceiverId).First(&a).Error; err != nil {
			btlLog.CUST.Error("GetAccount error:%s", err)
			mission.Retries += 1
			mission.Error = err.Error()
			return
		}
		rusr, err := account.GetUserInfo(a.UserName)
		if err != nil {
			btlLog.CUST.Error("GetUserInfo error:%s", err)
			mission.Retries += 1
			mission.Error = err.Error()
			return
		}
		//创建BillBalance记录
		tx, back := middleware.GetTx()
		defer back()
		rBalance := getBillBalanceModel(rusr, mission.Amount, i.AssetId, models.AWAY_IN, i)
		if err = tx.Create(rBalance).Error; err != nil {
			btlLog.CUST.Error("CreateBillBalance error:%s", err)
			mission.Retries += 1
			mission.Error = err.Error()
			return
		}
		mission.ReceiverBalanceId = rBalance.ID
		//创建收款记录
		_, err = AddAssetBalance(tx, rusr, mission.Amount, rBalance.ID, i.AssetId, custodyModels.ChangeTypeAssetReceiveLocal)
		if err != nil {
			mission.Retries += 1
			mission.Error = err.Error()
			return
		}
		mission.State = custodyModels.AIMStateSuccess
		tx.Commit()
		return

	}
}

func getBillBalanceModel(usr *account.UserInfo, amount float64, assetId string, away models.BalanceAway, invoice invoiceInfo) *models.Balance {
	ba := models.Balance{}
	ba.AccountId = usr.Account.ID
	ba.Amount = amount
	ba.AssetId = &assetId
	ba.Unit = models.UNIT_ASSET_NORMAL
	ba.BillType = models.BillTypeAssetTransfer
	ba.Away = away
	ba.Invoice = &invoice.Invoice
	ba.PaymentHash = invoice.Hash
	ba.State = models.STATE_SUCCESS
	ba.TypeExt = &models.BalanceTypeExt{Type: models.BTExtLocal}
	return &ba
}
