package btc_channel

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
			var e error
			if invoice.State == lnrpc.Invoice_SETTLED {
				tx := middleware.DB.Begin()
				if tx.Error != nil {
					btlLog.CUST.Error("invoice server 创建事务失败")
					continue
				}
				if invoice.CreationDate < time.Now().Unix()-60*60*24*3 {
					tx.Rollback()
					continue
				}
				var i models.Invoice
				if e = tx.Where("invoice =?", invoice.PaymentRequest).First(&i).Error; e != nil {
					tx.Rollback()
					continue
				}
				if i.Status == 1 {
					continue
				}
				i.Status = 1
				if e = tx.Save(&i).Error; e != nil {
					tx.Rollback()
					continue
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
				if e = tx.Create(&ba).Error; e != nil {
					tx.Rollback()
					continue
				}
				if e = tx.Commit().Error; e != nil {
					btlLog.CUST.Error("invoice server Error")
				}
			}
			if invoice.State == lnrpc.Invoice_CANCELED {
				tx := middleware.DB.Begin()
				if tx.Error != nil {
					btlLog.CUST.Error("invoice server 创建事务失败")
					continue
				}
				var i models.Invoice
				if e = tx.Where("invoice =?", invoice.PaymentRequest).First(&i).Error; e != nil {
					tx.Rollback()
					continue
				}
				if i.Status == 3 {
					continue
				}
				i.Status = 2
				if e = tx.Save(&i).Error; e != nil {
					tx.Rollback()
					continue
				}
				if e = tx.Commit().Error; e != nil {
					btlLog.CUST.Error("invoice server Error")
				}
			}
		}
	}
}
