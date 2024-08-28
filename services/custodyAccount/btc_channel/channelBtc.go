package btc_channel

import (
	"context"
	"errors"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"path/filepath"
	"sync"
	"time"
	"trade/btlLog"
	"trade/config"
	"trade/models"
	"trade/services/btldb"
	caccount "trade/services/custodyAccount/account"
	cBase "trade/services/custodyAccount/custodyBase"
	rpc "trade/services/servicesrpc"
)

const (
	timeout   = 20 * time.Second
	serverFee = 0
)

type BtcChannel error

var (
	GetbalanceErr = errors.New("GetbalanceErr")
	TimeoutErr    = errors.New("TimeoutErr")
)

//构建用户的BTC通道事件

type BtcChannelEvent struct {
	UserInfo *caccount.UserInfo
}

func NewBtcChannelEvent(UserName string) (*BtcChannelEvent, error) {
	var (
		e   BtcChannelEvent
		err error
	)
	e.UserInfo, err = caccount.GetUserInfo(UserName)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, caccount.CustodyAccountGetErr
	}
	return &e, nil
}

//获取余额

func (e *BtcChannelEvent) GetBalance() ([]cBase.Balance, error) {
	acc, err := rpc.AccountInfo(e.UserInfo.Account.UserAccountCode)
	if err != nil {
		return nil, err
	}
	balances := []cBase.Balance{
		{
			AssetId: "00",
			Amount:  acc.CurrentBalance,
		},
	}
	return balances, nil
}

//请求发票

var ApplyInvoiceMutex sync.Mutex

var CreateInvoiceErr = errors.New("CreateInvoiceErr")

func (e *BtcChannelEvent) ApplyPayReq(Request cBase.PayReqApplyRequest) (cBase.PayReqApplyResponse, error) {
	var applyRequest *BtcApplyInvoiceRequest
	var ok bool
	if applyRequest, ok = Request.(*BtcApplyInvoiceRequest); !ok {
		return nil, errors.New("invalid apply request")
	}
	//获取马卡龙路径
	var macaroonFile string
	macaroonDir := config.GetConfig().ApiConfig.CustodyAccount.MacaroonDir
	if e.UserInfo.Account.UserAccountCode == "admin" {
		macaroonFile = config.GetConfig().ApiConfig.Lnd.MacaroonPath
	} else {
		macaroonFile = filepath.Join(macaroonDir, e.UserInfo.Account.UserAccountCode+".macaroon")
	}
	if macaroonFile == "" {
		btlLog.CUST.Error(caccount.MacaroonFindErr.Error())
		return nil, caccount.MacaroonFindErr
	}
	//调用Lit节点发票申请接口
	invoice, err := rpc.InvoiceCreate(applyRequest.Amount, applyRequest.Memo, macaroonFile)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, fmt.Errorf("%w: %s", CreateInvoiceErr, err.Error())
	}
	//获取发票信息
	info, _ := rpc.InvoiceFind(invoice.RHash)

	//构建invoice对象
	var invoiceModel models.Invoice
	invoiceModel.UserID = e.UserInfo.User.ID
	invoiceModel.Invoice = invoice.PaymentRequest
	invoiceModel.AccountID = &e.UserInfo.Account.ID
	invoiceModel.Amount = float64(info.Value)

	invoiceModel.Status = models.InvoiceStatus(info.State)
	template := time.Unix(info.CreationDate, 0)
	invoiceModel.CreateDate = &template
	expiry := int(info.Expiry)
	invoiceModel.Expiry = &expiry
	//写入数据库
	err = btldb.CreateInvoice(&invoiceModel)
	if err != nil {
		btlLog.CUST.Error(err.Error(), models.ReadDbErr)
		return nil, models.ReadDbErr
	}
	return &BtcApplyInvoice{
		LnInvoice: invoice,
		Amount:    applyRequest.Amount,
	}, nil
}

func (e *BtcChannelEvent) SendPayment(payRequest cBase.PayPacket) error {
	var bt *BtcPacket
	var ok bool
	if bt, ok = payRequest.(*BtcPacket); !ok {
		return errors.New("invalid pay request")
	}
	bt.err = make(chan error, 1)
	defer close(bt.err)
	useableBalance, err := e.GetBalance()
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return GetbalanceErr
	}
	err = bt.VerifyPayReq(useableBalance[0].Amount)
	if err != nil {
		return err
	}
	if bt.isInsideMission != nil {
		//发起本地转账
		bt.isInsideMission.err = bt.err
		go e.payToInside(bt)
	} else {
		//发起外部转账
		go e.payToOutside(bt)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	select {
	case <-ctx.Done():
		//超时处理
		return TimeoutErr
	case err = <-bt.err:
		//错误处理
		return err
	}
}

func (e *BtcChannelEvent) payToInside(bt *BtcPacket) {
	//创建内部转账任务
	payInside := models.PayInside{
		PayUserId:     e.UserInfo.User.ID,
		GasFee:        uint64(bt.DecodePayReq.NumSatoshis),
		ServeFee:      serverFee,
		ReceiveUserId: bt.isInsideMission.insideInvoice.UserID,
		PayType:       models.PayInsideByInvioce,
		AssetType:     "00",
		PayReq:        &bt.PayReq,
		Status:        models.PayInsideStatusPending,
	}
	//写入数据库
	err := btldb.CreatePayInside(&payInside)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		bt.err <- err
		return
	}
	//递交给内部转账服务
	bt.isInsideMission.insideMission = &payInside
	BtcSever.NewMission(bt.isInsideMission)

}

func (e *BtcChannelEvent) payToOutside(bt *BtcPacket) {
	var macaroonFile string
	macaroonDir := config.GetConfig().ApiConfig.CustodyAccount.MacaroonDir

	if e.UserInfo.Account.UserAccountCode == "admin" {
		macaroonFile = config.GetConfig().ApiConfig.Lnd.MacaroonPath
	} else {
		macaroonFile = filepath.Join(macaroonDir, e.UserInfo.Account.UserAccountCode+".macaroon")
	}
	if macaroonFile == "" {
		btlLog.CUST.Error("macaroon file not found")
		bt.err <- nil
		return
	}
	payment, err := rpc.InvoicePay(macaroonFile, bt.PayReq, bt.DecodePayReq.NumSatoshis, bt.FeeLimit)
	if err != nil {
		btlLog.CUST.Error("pay invoice fail")
		bt.err <- err
		return
	}

	bt.err <- nil
	var balanceModel models.Balance
	//TODO：扣除服务费
	//if HasServerFee {
	//err = PayServerFee(bt.account)
	//balanceModel.ServerFee = GetServerFee() + uint64(payment.FeeSat)
	//}
	switch payment.Status {
	case lnrpc.Payment_SUCCEEDED:
		balanceModel.State = models.STATE_SUCCESS
	case lnrpc.Payment_FAILED:
		balanceModel.State = models.STATE_FAILED
	default:
		balanceModel.State = models.STATE_UNKNOW
	}
	balanceModel.AccountId = e.UserInfo.Account.ID
	balanceModel.BillType = models.BILL_TYPE_PAYMENT
	balanceModel.Away = models.AWAY_OUT
	balanceModel.Amount = float64(payment.ValueSat)
	balanceModel.Unit = models.UNIT_SATOSHIS
	balanceModel.Invoice = &payment.PaymentRequest
	balanceModel.PaymentHash = &payment.PaymentHash

	err = btldb.CreateBalance(&balanceModel)
	if err != nil {
		btlLog.CUST.Error(err.Error())
	}
	btlLog.CUST.Info("payment outside success balanceId:%v,amount:%v,%v", balanceModel.ID, balanceModel.Amount)
}

func (e *BtcChannelEvent) GetTransactionHistory() {

}
