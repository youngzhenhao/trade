package assets

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"trade/btlLog"
	"trade/config"
	"trade/models"
	"trade/services/assetsyncinfo"
	"trade/services/btldb"
	caccount "trade/services/custodyAccount/account"
	cBase "trade/services/custodyAccount/custodyBase"
	rpc "trade/services/servicesrpc"
)

type AssetEvent struct {
	UserInfo *caccount.UserInfo
	AssetId  *string
}

func NewAssetEvent(UserName string, AssetId string) (*AssetEvent, error) {
	var (
		e   AssetEvent
		err error
	)
	e.UserInfo, err = caccount.GetUserInfo(UserName)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, caccount.CustodyAccountGetErr
	}
	e.AssetId = &AssetId
	return &e, nil
}

func (e *AssetEvent) GetBalance() ([]cBase.Balance, error) {
	balance, err := btldb.GetAccountBalanceByGroup(e.UserInfo.Account.ID, *e.AssetId)
	if err != nil {
		return nil, models.ReadDbErr
	}
	balances := []cBase.Balance{
		{
			AssetId: balance.AssetId,
			Amount:  int64(balance.Amount),
		},
	}
	return balances, nil
}
func (e *AssetEvent) GetBalances() ([]cBase.Balance, error) {
	temp, err := btldb.GetAccountBalanceByAccountId(e.UserInfo.Account.ID)
	if err != nil {
		return nil, models.ReadDbErr
	}
	var balances []cBase.Balance
	for _, b := range *temp {
		balances = append(balances, cBase.Balance{
			AssetId: b.AssetId,
			Amount:  int64(b.Amount),
		})
	}
	return balances, nil
}
func (e *AssetEvent) GetCustodyAssetPermission(assetId, universe string) (*models.AssetSyncInfo, error) {
	r := assetsyncinfo.SyncInfoRequest{
		Id:       assetId,
		Universe: universe,
	}
	s, err := assetsyncinfo.GetAssetSyncInfo(&r)
	if err != nil {
		return nil, err
	}
	if s.AssetType == models.AssetTypeNFT {
		return nil, fmt.Errorf("NFT custody is not supported")
	}
	return s, nil
}

var ApplyAddrMutex sync.Mutex

var CreateAddrErr = errors.New("CreateAddrErr")

func (e *AssetEvent) ApplyPayReq(Request cBase.PayReqApplyRequest) (cBase.PayReqApplyResponse, error) {
	var applyRequest *AssetAddressApplyRequest
	var ok bool
	if applyRequest, ok = Request.(*AssetAddressApplyRequest); !ok {
		return nil, errors.New("invalid apply request")
	}
	universe := config.GetConfig().ApiConfig.Tapd.UniverseHost
	//调用Lit节点发票申请接口
	addr, err := rpc.NewAddr(*e.AssetId, int(applyRequest.Amount), universe)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, fmt.Errorf("%w: %s", CreateAddrErr, err.Error())
	}
	template := time.Now()
	expiry := 0
	//构建invoice对象
	var invoiceModel models.Invoice
	invoiceModel.UserID = e.UserInfo.User.ID
	invoiceModel.Invoice = addr.Encoded
	invoiceModel.AccountID = &e.UserInfo.Account.ID
	invoiceModel.AssetId = *e.AssetId
	invoiceModel.Amount = float64(addr.Amount)
	invoiceModel.Status = models.InvoiceStatusIsTaproot
	invoiceModel.CreateDate = &template
	invoiceModel.Expiry = &expiry
	//写入数据库
	err = btldb.CreateInvoice(&invoiceModel)
	if err != nil {
		btlLog.CUST.Error(err.Error(), models.ReadDbErr)
		return nil, models.ReadDbErr
	}
	return &AssetApplyAddress{
		Addr:   addr,
		Amount: applyRequest.Amount,
	}, nil
}

func (e *AssetEvent) SendPayment(payRequest cBase.PayPacket) error {
	var bt *AssetPacket
	var ok bool
	if bt, ok = payRequest.(*AssetPacket); !ok {
		return errors.New("invalid pay request")
	}
	bt.err = make(chan error, 1)
	defer close(bt.err)

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
		//超时处理
		return cBase.TimeoutErr
	case err = <-bt.err:
		//错误处理
		return err
	}

}

func (e *AssetEvent) payToInside(bt *AssetPacket) {

}
func (e *AssetEvent) payToOutside(bt *AssetPacket) {
	m := OutsideMission{
		AddrTarget: []*target{
			{
				Addr: bt.DecodePayReq.Encoded,
			},
		},
		AssetID:     string(bt.DecodePayReq.AssetId),
		TotalAmount: int64(bt.DecodePayReq.Amount),
		err:         []chan error{bt.err},
	}
	BtcSever.Queue.addNewPkg(&m)
}
func (e *AssetEvent) GetTransactionHistory() {

}