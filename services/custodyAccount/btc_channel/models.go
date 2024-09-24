package btc_channel

import (
	"errors"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"gorm.io/gorm"
	"trade/btlLog"
	"trade/models"
	"trade/services/btldb"
	caccount "trade/services/custodyAccount/account"
	cBase "trade/services/custodyAccount/custodyBase"
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
	ServerFee := ChannelBtcServiceFee
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
		ServerFee = ChannelBtcInsideServiceFee
	}
	//解码发票
	p.DecodePayReq, err = rpc.InvoiceDecode(p.PayReq)
	if err != nil {
		btlLog.CUST.Error("发票解析失败", err)
		return fmt.Errorf("(pay_request=%s)", "发票解析失败：", p.PayReq)
	}
	//TODO:限额，暂时做单次限额，后续改为每日限额
	if (p.DecodePayReq.NumSatoshis + p.FeeLimit + int64(ServerFee)) > 500000 {
		btlLog.CUST.Error("amount>500000,超过当前转账限制")
		return fmt.Errorf("amount>500000,超过当前转账限制")
	}
	//验证金额
	useableBalance, err := rpc.AccountInfo(userinfo.Account.UserAccountCode)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return cBase.GetbalanceErr
	}
	if useableBalance.CurrentBalance < (p.DecodePayReq.NumSatoshis + p.FeeLimit + int64(ServerFee)) {
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

type BtcPaymentList struct {
	PaymentList []PaymentResponse `json:"payments"`
}

func (r *BtcPaymentList) GetTxString() string {
	return ""
}

type PaymentResponse struct {
	Timestamp int64               `json:"timestamp"`
	BillType  models.BalanceType  `json:"bill_type"`
	Away      models.BalanceAway  `json:"away"`
	Invoice   *string             `json:"invoice"`
	Amount    float64             `json:"amount"`
	AssetId   *string             `json:"asset_id"`
	State     models.BalanceState `json:"state"`
	Fee       uint64              `json:"fee"`
}

type InvoiceResponce struct {
	Invoice string               `json:"invoice"`
	AssetId string               `json:"asset_id"`
	Amount  int64                `json:"amount"`
	Status  models.InvoiceStatus `json:"status"`
}
