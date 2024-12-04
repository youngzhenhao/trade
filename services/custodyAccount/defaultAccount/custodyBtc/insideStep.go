package custodyBtc

import (
	"encoding/hex"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"trade/btlLog"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
	"trade/services/custodyAccount/account"
	"trade/services/custodyAccount/custodyBase/custodyLimit"
	"trade/services/custodyAccount/custodyBase/custodyPayTN"
	rpc "trade/services/servicesrpc"
)

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
	DecodePayReq, err := rpc.InvoiceDecode(invoice.Invoice)
	if err != nil {
		btlLog.CUST.Error("发票解析失败", err)
		return err
	}
	i := invoiceInfo{
		Invoice: invoice.Invoice,
		Hash:    DecodePayReq.PaymentHash,
	}
	//run steps
	for {
		InsideSteps(usr, mission, i)
		LogAIM(middleware.DB, mission)
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

func RunInsidePTNStep(usr *account.UserInfo, receiveUsr *account.UserInfo, mission *custodyModels.AccountInsideMission) error {
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
	if receiveUsr == nil {
		var a models.Account
		if err := db.Where("id =?", mission.ReceiverId).First(&a).Error; err != nil {
			btlLog.CUST.Error("GetAccount error:%s", err)
			return err
		}
		receiveUsr, _ = account.GetUserInfo(a.UserName)
	}
	//获取发票信息
	PTN := custodyPayTN.PayToNpubKey{
		NpubKey: receiveUsr.User.Username,
		Amount:  mission.Amount,
		AssetId: mission.AssetId,
		Time:    mission.CreatedAt.Unix(),
		Vision:  0,
	}
	invoice, _ := PTN.Encode()
	h, _ := custodyPayTN.HashEncodedString(invoice)

	i := invoiceInfo{
		Invoice: invoice,
		Hash:    h,
	}
	//run steps
	for {
		InsideSteps(usr, mission, i)
		LogAIM(middleware.DB, mission)
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
		balance := getBillBalanceModel(usr, mission.Amount, models.AWAY_OUT, i)
		if err = tx.Create(balance).Error; err != nil {
			btlLog.CUST.Error("CreateBillBalance error:%s", err)
			mission.Error = err.Error()
			mission.State = custodyModels.AIMStateDone
			return
		}
		mission.PayerBalanceId = balance.ID
		//创建扣款记录
		_, err = LessBtcBalance(tx, usr, mission.Amount, mission.PayerBalanceId, custodyModels.ChangeTypeBtcPayLocal)
		if err != nil {
			btlLog.CUST.Error("LessBtcBalance error:%s", err)
			mission.Error = err.Error()
			mission.State = custodyModels.AIMStateDone
			return
		}
		//扣除手续费
		err = PayFee(tx, usr, mission.Fee, mission.PayerBalanceId)
		if err != nil {
			btlLog.CUST.Error("PayFee error:%s", err)
			mission.Error = err.Error()
			mission.State = custodyModels.AIMStateDone
			return
		}
		balance.ServerFee = uint64(mission.Fee)
		//签收发票
		err = tx.Model(&models.Invoice{}).
			Where("id =?", mission.InvoiceId).
			Updates(&models.Invoice{Status: models.InvoiceStatusLocal}).Error
		if err != nil {
			btlLog.CUST.Error("UpdateInvoice error:%s", err)
			mission.Error = err.Error()
			mission.State = custodyModels.AIMStateDone
			return
		}
		//更新状态
		balance.State = models.STATE_SUCCESS
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
				AssetId:      "00",
				TransferType: custodyModels.LimitTransferTypeLocal,
			}
			err = custodyLimit.MinusLimit(middleware.DB, usr, &limitType, mission.Amount+mission.Fee)
			if err != nil {
				btlLog.CUST.Error("额度限制未正常更新:%s", err.Error())
				btlLog.CUST.Error("error PayInsideId:%v", mission.ID)
			}
			//取消发票
			if strings.HasPrefix(i.Invoice, "PTN") {
				//PTN发票不需要取消
				return
			}
			h, _ := hex.DecodeString(i.Hash)
			err = rpc.InvoiceCancel(h)
			if err != nil {
				btlLog.CUST.Error("取消发票失败 %s", i.Hash)
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
		rBalance := getBillBalanceModel(rusr, mission.Amount, models.AWAY_IN, i)
		if err = tx.Create(rBalance).Error; err != nil {
			btlLog.CUST.Error("CreateBillBalance error:%s", err)
			mission.Retries += 1
			mission.Error = err.Error()
			return
		}
		mission.ReceiverBalanceId = rBalance.ID
		//创建收款记录
		_, err = AddBtcBalance(tx, rusr, mission.Amount, rBalance.ID, custodyModels.ChangeTypeBtcReceiveLocal)
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

type invoiceInfo struct {
	Invoice string
	Hash    string
}

func getBillBalanceModel(usr *account.UserInfo, amount float64, away models.BalanceAway, invoice invoiceInfo) *models.Balance {
	ba := models.Balance{}
	ba.AccountId = usr.Account.ID
	ba.Amount = amount
	ba.Unit = models.UNIT_SATOSHIS
	ba.BillType = models.BillTypePayment
	ba.Away = away
	ba.Invoice = &invoice.Invoice
	ba.PaymentHash = &invoice.Hash
	ba.State = models.STATE_SUCCESS
	ba.TypeExt = &models.BalanceTypeExt{Type: models.BTExtLocal}
	return &ba
}

func LogAIM(tx *gorm.DB, mission *custodyModels.AccountInsideMission) {
	tx.Save(mission)
}

func LoadAIMMission() {
	var missions []custodyModels.AccountInsideMission
	middleware.DB.Where("type = 'btc' AND (state =? OR state =?)", custodyModels.AIMStatePending, custodyModels.AIMStatePaid).Find(&missions)
	for _, m := range missions {
		if m.InvoiceId == 0 {

		} else {
			_ = RunInsideStep(nil, &m)
		}

	}
}
