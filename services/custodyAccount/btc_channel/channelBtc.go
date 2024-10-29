package btc_channel

import (
	"context"
	"errors"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"gorm.io/gorm"
	"path/filepath"
	"time"
	"trade/btlLog"
	"trade/config"
	"trade/models"
	"trade/services/btldb"
	caccount "trade/services/custodyAccount/account"
	cBase "trade/services/custodyAccount/custodyBase"
	"trade/services/custodyAccount/custodyBase/custodyFee"
	"trade/services/custodyAccount/custodyBase/custodyRpc"
	rpc "trade/services/servicesrpc"
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
		btlLog.CUST.Warning("%s,UserName:%s", err.Error(), UserName)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %s", caccount.CustodyAccountGetErr, "userName不存在")
		}
		return nil, caccount.CustodyAccountGetErr
	}
	btlLog.CUST.Info("UserName:%s", UserName)
	return &e, nil
}

func NewBtcChannelEventByUserId(UserId uint) (*BtcChannelEvent, error) {
	var (
		e   BtcChannelEvent
		err error
	)
	e.UserInfo, err = caccount.GetUserInfoById(UserId)
	if err != nil {
		btlLog.CUST.Error("%s,UserName:%s", err.Error(), UserId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %s", caccount.CustodyAccountGetErr, "userName不存在")
		}
		return nil, caccount.CustodyAccountGetErr
	}
	btlLog.CUST.Info("UserName:%s", UserId)
	return &e, nil
}

//获取余额

func (e *BtcChannelEvent) GetBalance() ([]cBase.Balance, error) {
	acc, err := custodyRpc.GetAccountInfo(e.UserInfo)
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

func (e *BtcChannelEvent) QueryPayReq() ([]InvoiceResponce, error) {
	params := btldb.QueryParams{
		"UserID":  e.UserInfo.User.ID,
		"AssetId": "00",
	}
	a, err := btldb.GenericQuery(&models.Invoice{}, params)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, err
	}
	if len(a) > 0 {
		var invoices []InvoiceResponce
		for j := len(a) - 1; j >= 0; j-- {
			var i InvoiceResponce
			i.Invoice = a[j].Invoice
			i.AssetId = a[j].AssetId
			i.Amount = int64(a[j].Amount)
			i.Status = a[j].Status
			invoices = append(invoices, i)
		}
		return invoices, nil
	}
	return nil, nil
}

func (e *BtcChannelEvent) SendPayment(payRequest cBase.PayPacket) error {
	var bt *BtcPacket
	var ok bool
	if bt, ok = payRequest.(*BtcPacket); !ok {
		return errors.New("invalid pay request")
	}
	bt.err = make(chan error, 1)
	//defer close(bt.err)

	err := bt.VerifyPayReq(e.UserInfo)
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
	ctx, cancel := context.WithTimeout(context.Background(), cBase.Timeout)
	defer cancel()
	select {
	case <-ctx.Done():
		go func(c chan error) {
			err := <-c
			if err != nil {
				btlLog.CUST.Error("btc sendPayment timeout:%s", err.Error())
			}
			close(c)
		}(bt.err)
		//超时处理
		return cBase.TimeoutErr
	case err := <-bt.err:
		//错误处理
		return err
	}
}

func (e *BtcChannelEvent) payToInside(bt *BtcPacket) {
	//创建内部转账任务
	payInside := models.PayInside{
		PayUserId:     e.UserInfo.User.ID,
		GasFee:        uint64(bt.DecodePayReq.NumSatoshis),
		ServeFee:      0,
		ReceiveUserId: bt.isInsideMission.insideInvoice.UserID,
		PayType:       models.PayInsideByInvoice,
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
	bt.isInsideMission.err = bt.err
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
		bt.err <- fmt.Errorf("account is abnormal")
		return
	}
	var balanceModel models.Balance
	balanceModel.State = models.STATE_UNKNOW
	balanceModel.AccountId = e.UserInfo.Account.ID
	balanceModel.BillType = models.BillTypePayment
	balanceModel.Away = models.AWAY_OUT
	balanceModel.Amount = float64(bt.DecodePayReq.NumSatoshis)
	balanceModel.Unit = models.UNIT_SATOSHIS
	balanceModel.Invoice = &bt.PayReq
	balanceModel.PaymentHash = &bt.DecodePayReq.PaymentHash
	balanceModel.State = models.STATE_UNKNOW
	err := btldb.CreateBalance(&balanceModel)
	if err != nil {
		btlLog.CUST.Error(err.Error())
	}

	payment, err := custodyRpc.PayBtcInvoice(e.UserInfo, macaroonFile, bt.PayReq, bt.DecodePayReq.NumSatoshis, bt.FeeLimit)
	if err != nil {
		btlLog.CUST.Error("pay invoice fail %s", err.Error())
		bt.err <- err
		return
	}
	track, err := rpc.PaymentTrack(payment.PaymentHash)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		bt.err <- fmt.Errorf("payment outside unknown,Please contact the administrator")
		return
	}
	switch track.Status {
	case lnrpc.Payment_SUCCEEDED:
		err = custodyFee.PayServiceFeeSync(e.UserInfo, custodyFee.ChannelBtcServiceFee, balanceModel.ID, models.ChannelBTCOutSideFee, "payToOutside Fee")
		if err != nil {
			btlLog.CUST.Error(err.Error())
		}
		balanceModel.ServerFee = custodyFee.ChannelBtcServiceFee + uint64(track.FeeSat)
		balanceModel.State = models.STATE_SUCCESS
		bt.err <- nil
		btlLog.CUST.Info("payment outside success balanceId:%v,amount:%v,%v", balanceModel.ID, balanceModel.Amount)
	case lnrpc.Payment_FAILED:
		btlLog.CUST.Error("payment outside failed balanceId:%v,amount:%v,%v", balanceModel.ID, balanceModel.Amount)
		btlLog.CUST.Error(track.FailureReason.String())
		bt.err <- fmt.Errorf(track.FailureReason.String())
		balanceModel.State = models.STATE_FAILED
	default:
		btlLog.CUST.Error("payment outside unknown balanceId:%v,amount:%v,%v", balanceModel.ID, balanceModel.Amount)
		balanceModel.State = models.STATE_UNKNOW
		bt.err <- fmt.Errorf("payment outside unknown,Please contact the administrator")
		return
	}
	err = btldb.UpdateBalance(&balanceModel)
	if err != nil {
		btlLog.CUST.Error(err.Error())
	}
}

func (e *BtcChannelEvent) GetTransactionHistory() (*cBase.PaymentList, error) {
	params := btldb.QueryParams{
		"AccountId": e.UserInfo.Account.ID,
		"AssetId":   "00",
	}
	a, err := btldb.GenericQuery(&models.Balance{}, params)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, fmt.Errorf("query payment error")
	}
	var results cBase.PaymentList
	if len(a) > 0 {
		for i := len(a) - 1; i >= 0; i-- {
			v := a[i]
			r := cBase.PaymentResponse{}
			r.Timestamp = v.CreatedAt.Unix()
			r.BillType = v.BillType
			r.Away = v.Away
			r.Invoice = v.Invoice
			r.Address = v.Invoice
			r.Target = v.Invoice
			r.PaymentHash = v.PaymentHash
			if *v.Invoice == "award" && v.PaymentHash != nil {
				awardType := cBase.GetAwardType(*v.PaymentHash)
				r.Invoice = &awardType
				r.Address = &awardType
				r.Target = &awardType
			}
			r.Amount = v.Amount
			btcAssetId := "00"
			r.AssetId = &btcAssetId
			r.State = v.State
			r.Fee = v.ServerFee
			results.PaymentList = append(results.PaymentList, r)
		}
	}
	return &results, nil
}
