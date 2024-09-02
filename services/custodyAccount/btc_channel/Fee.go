package btc_channel

import (
	"fmt"
	"time"
	"trade/btlLog"
	"trade/config"
	"trade/models"
	"trade/services/btldb"
	caccount "trade/services/custodyAccount/account"
	rpc "trade/services/servicesrpc"
)

var (
	ChannelBtcServiceFee = uint64(100)
	AssetServiceFee      = 100
)

func PayServerFee(account *models.Account, fee uint64) error {
	acc, err := rpc.AccountInfo(account.UserAccountCode)
	if err != nil {
		return err
	}
	// Change the escrow account balance
	_, err = rpc.AccountUpdate(account.UserAccountCode, acc.CurrentBalance-int64(fee), -1)
	if err != nil {
		return err
	}
	return nil
}

func PayFirLunchFee(e *BtcChannelEvent, gasFee uint64) (uint, error) {
	//获取账户信息
	acc, err := rpc.AccountInfo(e.UserInfo.Account.UserAccountCode)
	if err != nil {
		btlLog.CUST.Error("AccountInfo error(UserId=%v):%v", e.UserInfo.User.ID, err)
		return 0, fmt.Errorf("AccountInfo error")
	}
	//检查账户余额是否足够
	if acc.CurrentBalance < int64(gasFee) {
		btlLog.CUST.Error("Account balance not enough(UserId=%v)", e.UserInfo.User.ID)
		return 0, fmt.Errorf("account balance not enough")
	}
	//管理员账户
	macaroonFile := config.GetConfig().ApiConfig.Lnd.MacaroonPath
	if macaroonFile == "" {
		btlLog.CUST.Error(caccount.MacaroonFindErr.Error())
		return 0, caccount.MacaroonFindErr
	}
	//调用Lit节点发票申请接口
	res, err := rpc.InvoiceCreate(int64(gasFee), SetMemoSign(), macaroonFile)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return 0, fmt.Errorf("%w: %s", CreateInvoiceErr, err.Error())
	}

	//构建invoice对象
	var invoiceModel models.Invoice
	var adminId uint = 1
	invoiceModel.UserID = adminId
	invoiceModel.AccountID = &adminId

	invoiceModel.Invoice = res.PaymentRequest
	invoiceModel.Amount = float64(gasFee)

	invoiceModel.Status = models.InvoiceStatusPending
	template := time.Now()
	invoiceModel.CreateDate = &template
	expiry := 86400
	invoiceModel.Expiry = &expiry
	//写入数据库
	err = btldb.CreateInvoice(&invoiceModel)
	if err != nil {
		btlLog.CUST.Error(err.Error(), models.ReadDbErr)
		return 0, models.ReadDbErr
	}

	//创建内部转账任务
	payInside := models.PayInside{
		PayUserId:     e.UserInfo.User.ID,
		GasFee:        gasFee,
		ServeFee:      0,
		ReceiveUserId: 1,
		PayType:       models.PayInsideToAdmin,
		AssetType:     "00",
		PayReq:        &res.PaymentRequest,
		Status:        models.PayInsideStatusPending,
	}
	//写入数据库
	err = btldb.CreatePayInside(&payInside)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return 0, err
	}
	//递交给内部转账服务
	var mission isInsideMission
	mission.isInside = true
	mission.insideInvoice = &invoiceModel
	mission.insideMission = &payInside
	BtcSever.NewMission(&mission)
	return payInside.ID, nil
}

// CheckFirLunchFee 检查内部转账任务状态是否成功
func CheckFirLunchFee(id uint) (bool, error) {
	p, err := btldb.ReadPayInside(id)
	if err != nil {
		return false, err
	}
	switch p.Status {
	case models.PayInsideStatusSuccess:
		return true, nil
	case models.PayInsideStatusFailed:
		return false, models.CustodyAccountPayInsideMissionFaild
	default:
		return false, models.CustodyAccountPayInsideMissionPending
	}
}
func SetMemoSign() string {
	return "internal transfer"
}
