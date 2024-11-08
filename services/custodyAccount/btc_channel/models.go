package btc_channel

import (
	"errors"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"gorm.io/gorm"
	"trade/btlLog"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
	"trade/services/btldb"
	caccount "trade/services/custodyAccount/account"
	cBase "trade/services/custodyAccount/custodyBase"
	"trade/services/custodyAccount/custodyBase/custodyFee"
	"trade/services/custodyAccount/custodyBase/custodyLimit"
	"trade/services/custodyAccount/custodyBase/custodyRpc"
	rpc "trade/services/servicesrpc"
)

// BtcApplyInvoice 申请发票返回的结构体
type BtcApplyInvoice struct {
	LnInvoice *lnrpc.AddInvoiceResponse
	Amount    int64
}

func (in *BtcApplyInvoice) GetAmount() int64 {
	return in.Amount
}
func (in *BtcApplyInvoice) GetPayReq() string {
	return in.LnInvoice.PaymentRequest
}

// BtcApplyInvoiceRequest 发票申请请求结构体
type BtcApplyInvoiceRequest struct {
	Amount int64
	Memo   string
}

func (req *BtcApplyInvoiceRequest) GetPayReqAmount() int64 {
	return req.Amount
}

type BtcPacketErr error

var (
	NotSufficientFunds BtcPacketErr = errors.New("not sufficient funds")
	DecodeInvoiceFail  BtcPacketErr = errors.New("decode invoice fail")
)

// BtcPacket 支付包结构体
type BtcPacket struct {
	PayReq          string
	FeeLimit        int64
	DecodePayReq    *lnrpc.PayReq
	isInsideMission *isInsideMission
	err             chan error
}

func (p *BtcPacket) VerifyPayReq(userinfo *caccount.UserInfo) error {
	ServerFee := custodyFee.ChannelBtcServiceFee
	//验证是否为本地发票
	i, err := btldb.GetInvoiceByReq(p.PayReq)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		btlLog.CUST.Error("验证本地发票失败", err)
		return models.ReadDbErr
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		p.isInsideMission = nil
	} else {
		if i.Status != models.InvoiceStatusPending {
			return fmt.Errorf("发票已被使用")
		}
		p.isInsideMission = &isInsideMission{
			isInside:      true,
			insideInvoice: i,
		}
		ServerFee = custodyFee.ChannelBtcInsideServiceFee
	}
	//解码发票
	p.DecodePayReq, err = rpc.InvoiceDecode(p.PayReq)
	if err != nil {
		btlLog.CUST.Error("发票解析失败", err)
		return fmt.Errorf("(pay_request=%s)", "发票解析失败：", p.PayReq)
	}
	endAmount := p.DecodePayReq.NumSatoshis + p.FeeLimit + int64(ServerFee)

	//验证限额
	limitType := custodyModels.LimitType{
		AssetId:      "00",
		TransferType: custodyModels.LimitTransferTypeLocal,
	}
	if p.isInsideMission == nil {
		limitType.TransferType = custodyModels.LimitTransferTypeOutside
	}
	err = custodyLimit.CheckLimit(middleware.DB, userinfo, &limitType, float64(endAmount))
	if err != nil {
		return err
	}

	//验证金额
	usableBalance, err := custodyRpc.GetAccountInfo(userinfo)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return cBase.GetbalanceErr
	}
	if usableBalance.CurrentBalance < endAmount {
		return NotSufficientFunds
	}

	return nil
}

// isInsideMission 内部任务结构体
type isInsideMission struct {
	isInside      bool
	insideInvoice *models.Invoice
	insideMission *models.PayInside
	err           chan error
}
type InvoiceResponce struct {
	Invoice string               `json:"invoice"`
	AssetId string               `json:"asset_id"`
	Amount  int64                `json:"amount"`
	Status  models.InvoiceStatus `json:"status"`
}
