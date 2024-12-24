package custodyBtc

import (
	"context"
	"encoding/hex"
	"github.com/lightningnetwork/lnd/lnrpc"
	"strconv"
	"time"
	"trade/btlLog"
	"trade/config"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
	caccount "trade/services/custodyAccount/account"
	"trade/services/custodyAccount/defaultAccount/custodyBtc/mempool"
	"trade/services/servicesrpc"
	"trade/utils"
)

type SubscribeInvoiceServer struct {
}

var InvoiceServer SubscribeInvoiceServer

func (s *SubscribeInvoiceServer) Start(ctx context.Context) {
	go s.runServer(ctx)
}
func (s *SubscribeInvoiceServer) runServer(ctx context.Context) {
	lndconf := config.GetConfig().ApiConfig.Lnd
	grpcHost := lndconf.Host + ":" + strconv.Itoa(lndconf.Port)
	tlsCertPath := lndconf.TlsCertPath
	macaroonPath := lndconf.MacaroonPath
	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.InvoiceSubscription{
		AddIndex: 1,
	}
	stream, err := client.SubscribeInvoices(ctx, request)
	if err != nil {
		return
	}
	for {
		invoice, err := stream.Recv()
		if err != nil {
			return
		}
		if invoice != nil {
			if invoice.State == lnrpc.Invoice_SETTLED {
				dealSettledInvoice(invoice)
				go subscriptionReceiveBtcBalance(float64(invoice.Value))
			} else if invoice.State == lnrpc.Invoice_CANCELED {
				DealCanceledInvoice(invoice)
			}
		}
	}
}
func dealSettledInvoice(invoice *lnrpc.Invoice) {
	tx := middleware.DB.Begin()
	defer tx.Rollback()
	if tx.Error != nil {
		btlLog.CUST.Error("invoice server 创建事务失败")
		return
	}
	if invoice.CreationDate < time.Now().Unix()-60*60*24*3 {
		return
	}
	var i models.Invoice
	var err error
	if err = tx.Where("invoice =? and status = 0", invoice.PaymentRequest).First(&i).Error; err != nil {
		return
	}
	i.Status = 1
	if err = tx.Save(&i).Error; err != nil {
		return
	}
	ba := models.Balance{}
	ba.AccountId = *i.AccountID
	ba.Amount = i.Amount
	ba.Unit = models.UNIT_SATOSHIS
	ba.BillType = models.BillTypeRecharge
	ba.Away = models.AWAY_IN
	ba.State = models.STATE_SUCCESS
	ba.Invoice = &i.Invoice
	hash := hex.EncodeToString(invoice.RHash)
	ba.PaymentHash = &hash
	ba.TypeExt = &models.BalanceTypeExt{Type: models.BTExtOnChannel}
	if err = tx.Create(&ba).Error; err != nil {
		return
	}
	// 余额变动
	UserInfo, err := caccount.GetUserInfoById(i.UserID)
	if err != nil {
		return
	}
	_, err = AddBtcBalance(tx, UserInfo, i.Amount, ba.ID, custodyModels.ChangeTypeBtcReceiveOutside)
	if err != nil {
		return
	}
	if err = tx.Commit().Error; err != nil {
		btlLog.CUST.Error("invoice server Error %s", err.Error())
	}
}
func DealCanceledInvoice(invoice *lnrpc.Invoice) {
	tx := middleware.DB.Begin()
	defer tx.Rollback()

	if tx.Error != nil {
		btlLog.CUST.Error("invoice server 创建事务失败")
		return
	}
	var i models.Invoice
	var err error
	if err = tx.Where("invoice =? and status = 0", invoice.PaymentRequest).First(&i).Error; err != nil {
		return
	}
	i.Status = 2
	if err = tx.Save(&i).Error; err != nil {
		return
	}
	if err = tx.Commit().Error; err != nil {
		btlLog.CUST.Error("invoice server Error")
	}
}

func subscriptionReceiveBtcBalance(amount float64) {
	if config.GetLoadConfig().NetWork == "regtest" {
		return
	}
	d := mempool.NewDingding()
	channels, err := servicesrpc.GetChannelInfo()
	if err != nil {
		btlLog.CUST.Error("GetChannelInfo error:%s", err)
		return
	}
	var balances []float64
	for _, c := range channels {
		if c.LocalBalance >= 0 {
			balances = append(balances, float64(c.LocalBalance))
		}
	}
	err = d.ReceiveBtcChannel(amount, balances)
	if err != nil {
		btlLog.CUST.Error("SendBtcPayOutChange error:%s", err)
	}
}
