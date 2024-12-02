package custodyBtc

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
	"trade/btlLog"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
	"trade/services/btldb"
	caccount "trade/services/custodyAccount/account"
	cBase "trade/services/custodyAccount/custodyBase"
	"trade/services/custodyAccount/custodyBase/custodyFee"
	"trade/services/custodyAccount/custodyBase/custodyLimit"
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
		return nil, fmt.Errorf("%w: %w", caccount.CustodyAccountGetErr, err)
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
		return nil, fmt.Errorf("%w: %w", caccount.CustodyAccountGetErr, err)
	}
	btlLog.CUST.Info("UserName:%s", UserId)
	return &e, nil
}

//获取余额

func (e *BtcChannelEvent) GetBalance() ([]cBase.Balance, error) {
	DB := middleware.DB
	balance := getBtcBalance(DB, e.UserInfo.Account.ID)
	balances := []cBase.Balance{
		{
			AssetId: "00",
			Amount:  int64(balance),
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
	//调用Lit节点发票申请接口
	invoice, err := rpc.InvoiceCreate(applyRequest.Amount, applyRequest.Memo)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, fmt.Errorf("%w: %s", CreateInvoiceErr, err.Error())
	}
	//TODO:取消Find接口
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

func (e *BtcChannelEvent) SendPaymentToUser(amount int64) error {
	var err error
	//验证余额，检查限额
	limitType := custodyModels.LimitType{
		AssetId:      "00",
		TransferType: custodyModels.LimitTransferTypeLocal,
	}
	err = custodyLimit.CheckLimit(middleware.DB, e.UserInfo, &limitType, float64(amount))
	if err != nil {
		return err
	}
	if !CheckBtcBalance(middleware.DB, e.UserInfo, float64(amount)) {
		return NotSufficientFunds
	}

	return nil
}

func (e *BtcChannelEvent) payToInside(bt *BtcPacket) {
	m := custodyModels.AccountInsideMission{
		AccountId:  e.UserInfo.Account.ID,
		AssetId:    BtcId,
		Type:       custodyModels.AIMTypeBtc,
		ReceiverId: *bt.isInsideMission.insideInvoice.AccountID,
		InvoiceId:  bt.isInsideMission.insideInvoice.ID,
		Amount:     float64(bt.DecodePayReq.NumSatoshis),
		Fee:        float64(custodyFee.ChannelBtcInsideServiceFee),
		FeeType:    BtcId,
		State:      custodyModels.AIMStatePending,
	}
	LogAIM(middleware.DB, &m)
	err := RunInsideStep(e.UserInfo, &m)
	bt.err <- err
}

func (e *BtcChannelEvent) payToOutside(bt *BtcPacket) {
	tx, back := middleware.GetTx()
	defer back()

	var balanceModel models.Balance
	balanceModel.AccountId = e.UserInfo.Account.ID
	balanceModel.BillType = models.BillTypePayment
	balanceModel.Away = models.AWAY_OUT
	balanceModel.Amount = float64(bt.DecodePayReq.NumSatoshis)
	balanceModel.Unit = models.UNIT_SATOSHIS
	balanceModel.Invoice = &bt.PayReq
	balanceModel.PaymentHash = &bt.DecodePayReq.PaymentHash
	balanceModel.State = models.STATE_UNKNOW
	balanceModel.TypeExt = &models.BalanceTypeExt{Type: models.BTExtOnChannel}
	err := btldb.CreateBalance(tx, &balanceModel)
	if err != nil {
		btlLog.CUST.Error(err.Error())
	}

	outsideMission := custodyModels.AccountOutsideMission{
		AccountId: e.UserInfo.Account.ID,
		AssetId:   "00",
		Type:      custodyModels.AOMTypeBtc,
		Target:    bt.PayReq,
		Hash:      bt.DecodePayReq.PaymentHash,
		Amount:    float64(bt.DecodePayReq.NumSatoshis),
		FeeLimit:  float64(bt.FeeLimit),
		BalanceId: balanceModel.ID,
		State:     custodyModels.AOMStatePending,
	}
	LogAOM(tx, &outsideMission)
	if err = tx.Commit().Error; err != nil {
		bt.err <- err
		return
	}
	err = RunOutsideSteps(e.UserInfo, &outsideMission)
	bt.err <- err
}

func (e *BtcChannelEvent) GetTransactionHistory(query *cBase.PaymentRequest) (*cBase.PaymentList, error) {

	if query.Page <= 0 {
		return nil, fmt.Errorf("page error")
	}

	db := middleware.DB
	var err error
	var a []models.Balance
	offset := (query.Page - 1) * query.PageSize
	q := db.Where("account_id = ? and asset_id = ?", e.UserInfo.Account.ID, "00")
	switch query.Away {
	case 0, 1:
		q = q.Where("away = ?", query.Away)
	default:
	}
	err = q.Order("created_at desc").
		Limit(query.PageSize).
		Offset(offset).
		Find(&a).Error
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
