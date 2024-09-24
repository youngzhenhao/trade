package custodyAssets

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"gorm.io/gorm"
	"trade/btlLog"
	"trade/models"
	"trade/services/btldb"
	caccount "trade/services/custodyAccount/account"
	"trade/services/custodyAccount/btc_channel"
	cBase "trade/services/custodyAccount/custodyBase"
	rpc "trade/services/servicesrpc"
)

// AssetAddressApplyRequest 资产接收地址申请请求结构体
type AssetAddressApplyRequest struct {
	Amount int64
}

func (req *AssetAddressApplyRequest) GetPayReqAmount() int64 {
	return req.Amount
}

// AssetApplyAddress 资产接收地址申请请求的结构体
type AssetApplyAddress struct {
	Addr   *taprpc.Addr
	Amount int64
}

func (a *AssetApplyAddress) GetAmount() int64 {
	return a.Amount
}
func (a *AssetApplyAddress) GetPayReq() string {
	return a.Addr.Encoded
}

// AssetPacket 支付包结构体
type AssetPacket struct {
	PayReq          string
	DecodePayReq    *taprpc.Addr
	isInsideMission *isInsideMission
	err             chan error
}

func (p *AssetPacket) VerifyPayReq(userinfo *caccount.UserInfo) error {
	ServerFee := btc_channel.AssetOutsideFee
	//验证是否为本地发票
	i, err := btldb.GetInvoiceByReq(p.PayReq)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		btlLog.CUST.Error("验证本地发票失败", err)
		return models.ReadDbErr
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		p.isInsideMission = nil
	} else {
		p.isInsideMission = &isInsideMission{
			isInside:      true,
			insideInvoice: i,
		}
		ServerFee = btc_channel.AssetInsideFee
	}
	//TODO:验证网络

	//解码地址
	p.DecodePayReq, err = rpc.DecodeAddr(p.PayReq)
	if err != nil {
		btlLog.CUST.Error("地址解析失败", err)
		return fmt.Errorf("%w(pay_request=%s)", cBase.DecodeAddressFail, p.PayReq)
	}
	//TODO:验证地址版本

	assetId := hex.EncodeToString(p.DecodePayReq.AssetId)
	//验证资产金额
	balance, err := btldb.GetAccountBalanceByGroup(userinfo.Account.ID, assetId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return cBase.NotEnoughAssetFunds
		}
		btlLog.CUST.Error("获取账户余额失败", err)
		return models.ReadDbErr
	}
	if balance.Amount < float64(p.DecodePayReq.Amount) {
		return cBase.NotEnoughAssetFunds
	}
	//验证托管账户余额
	useAbleBalance, err := rpc.AccountInfo(userinfo.Account.UserAccountCode)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return cBase.GetBalanceErr
	}
	if useAbleBalance.CurrentBalance < int64(ServerFee) {
		return cBase.NotEnoughFeeFunds
	}
	return nil
}

type PaymentResponse struct {
	Timestamp int64               `json:"timestamp"`
	BillType  models.BalanceType  `json:"bill_type"`
	Away      models.BalanceAway  `json:"away"`
	Invoice   *string             `json:"addr"`
	Amount    float64             `json:"amount"`
	AssetId   *string             `json:"asset_id"`
	State     models.BalanceState `json:"state"`
	Fee       uint64              `json:"fee"`
}

type BtcPaymentList struct {
	PaymentList []PaymentResponse `json:"payments"`
}

func (r *BtcPaymentList) GetTxString() string {
	return ""
}

// isInsideMission 内部任务结构体
type isInsideMission struct {
	isInside      bool
	insideInvoice *models.Invoice
	insideMission *models.PayInside
	err           chan error
}

// OutsideMission 外部事件结构体
type OutsideMission struct {
	AddrTarget       []*target
	AssetID          string
	TotalAmount      int64
	RollBackNumber   int64
	MinPaymentNumber int64
}

type target struct {
	Mission *models.PayOutside
}
