package custodyAssets

import (
	"context"
	"errors"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"gorm.io/gorm"
	"strconv"
	"time"
	"trade/btlLog"
	"trade/config"
	"trade/middleware"
	"trade/models"
	"trade/services/btldb"
	"trade/utils"
)

type SubscribeAddressServer struct {
	running bool
	Cancel  context.CancelFunc
}

var AddressServer SubscribeAddressServer

func (s *SubscribeAddressServer) Start(ctx context.Context) {
	loadEvent()
	go s.runServer(ctx)
}
func (s *SubscribeAddressServer) runServer(ctx context.Context) {
	tapdconf := config.GetConfig().ApiConfig.Tapd
	grpcHost := tapdconf.Host + ":" + strconv.Itoa(tapdconf.Port)
	tlsCertPath := tapdconf.TlsCertPath
	macaroonPath := tapdconf.MacaroonPath
	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()
	client := taprpc.NewTaprootAssetsClient(conn)
	request := &taprpc.SubscribeReceiveEventsRequest{}
	stream, err := client.SubscribeReceiveEvents(ctx, request)
	if err != nil {
		btlLog.CUST.Info("AddressServer start error")
		return
	}
	btlLog.CUST.Info("AddressServer start")
	s.running = true
	for {
		event, err := stream.Recv()
		if err != nil {
			btlLog.CUST.Info("AddressServer stop:%v", err.Error())
			s.running = false
			return
		}
		if event != nil {
			btlLog.CUST.Info("%v", event)
			tx := middleware.DB.Begin()
			if tx.Error != nil {
				btlLog.CUST.Error("address server 创建事务失败")
				continue
			}
			err = dealEvent(tx, event)
			if err != nil {
				btlLog.CUST.Error("address even deal error:%v", event.Outpoint)
				tx.Rollback()
				tx = nil
				continue
			}
			tx.Commit()
		}
	}
}

func dealEvent(tx *gorm.DB, event *taprpc.ReceiveEvent) error {
	var err error
	switch event.Status {
	case taprpc.AddrEventStatus_ADDR_EVENT_STATUS_COMPLETED:
		var i models.Invoice
		if err = tx.Where("invoice =?", event.Address.Encoded).First(&i).Error; err != nil {
			btlLog.CUST.Error(err.Error())
			return err
		}
		var r models.AccountAssetReceive
		if err = tx.Where("out_point =?", event.Outpoint).First(&r).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			btlLog.CUST.Error(err.Error())
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Timestamp = event.Timestamp
			r.OutPoint = event.Outpoint
			r.InvoiceId = i.ID
			r.Amount = float64(event.Address.Amount)
			r.Status = models.AddressStatusCOMPLETED
			if err = tx.Create(&r).Error; err != nil {
				btlLog.CUST.Error(err.Error())
				return err
			}
		} else {
			r.Status = models.AddressStatusCOMPLETED
			if err = tx.Save(&r).Error; err != nil {
				btlLog.CUST.Error(err.Error())
				return err
			}
		}
		b := models.Balance{
			AccountId:   *i.AccountID,
			BillType:    models.BillTypeAssetTransfer,
			Away:        models.AWAY_IN,
			Amount:      r.Amount,
			Unit:        models.UNIT_ASSET_NORMAL,
			ServerFee:   0,
			AssetId:     &i.AssetId,
			Invoice:     &i.Invoice,
			PaymentHash: &event.Outpoint,
			State:       models.STATE_SUCCESS,
			TypeExt: &models.BalanceTypeExt{
				Type: models.BTExtOnChannel,
			},
		}
		if err = tx.Create(&b).Error; err != nil {
			btlLog.CUST.Error(err.Error())
			return err
		}

		receiveBalance, err := btldb.GetAccountBalanceByGroup(*i.AccountID, i.AssetId)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			btlLog.CUST.Error("err:%v", models.ReadDbErr)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newBalance := models.AccountBalance{
				AccountID: *i.AccountID,
				AssetId:   i.AssetId,
				Amount:    r.Amount,
			}
			if err = tx.Create(&newBalance).Error; err != nil {
				btlLog.CUST.Error(err.Error())
				return err
			}
		} else {
			receiveBalance.Amount += r.Amount
			if err = tx.Save(&receiveBalance).Error; err != nil {
				btlLog.CUST.Error(err.Error())
				return err
			}
		}
		return nil

	default:
		var i models.Invoice
		if err = tx.Where("invoice =?", event.Address.Encoded).First(&i).Error; err != nil {
			btlLog.CUST.Error(err.Error())
			return err
		}

		var r models.AccountAssetReceive
		if err = tx.Where("out_point = ?", event.Outpoint).First(&r).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			btlLog.CUST.Error(err.Error())
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Timestamp = event.Timestamp
			r.OutPoint = event.Outpoint
			r.InvoiceId = i.ID
			r.Amount = float64(event.Address.Amount)
			r.Status = models.AddressStatus(event.Status)
			if err = tx.Create(&r).Error; err != nil {
				btlLog.CUST.Error(err.Error())
				return err
			}
		} else {
			r.Status = models.AddressStatus(event.Status)
			if err = tx.Save(&r).Error; err != nil {
				btlLog.CUST.Error(err.Error())
				return err
			}
		}
		return nil
	}
}

func loadEvent() {
	tapdconf := config.GetConfig().ApiConfig.Tapd
	grpcHost := tapdconf.Host + ":" + strconv.Itoa(tapdconf.Port)
	tlsCertPath := tapdconf.TlsCertPath
	macaroonPath := tapdconf.MacaroonPath
	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()
	client := taprpc.NewTaprootAssetsClient(conn)
	request := &taprpc.AddrReceivesRequest{}
	response, err := client.AddrReceives(context.Background(), request)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return
	}
	btlLog.CUST.Info("%v", uint64(time.Now().Unix()-3600*24))
	for i := len(response.Events) - 1; i >= 0; i-- {
		if response.Events[i].CreationTimeUnixSeconds < uint64(time.Now().Unix()-3600*24) {
			break
		}
		tx := middleware.DB.Begin()
		if tx.Error != nil {
			btlLog.CUST.Error("address server 创建事务失败")
			continue
		}
		var r models.AccountAssetReceive
		if err = tx.Where("out_point =?", response.Events[i].Outpoint).Where("status = 4").First(&r).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			btlLog.CUST.Error(err.Error())
			tx = nil
			continue
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			event := &taprpc.ReceiveEvent{
				Timestamp:          int64(response.Events[i].CreationTimeUnixSeconds) * 1000000,
				Address:            response.Events[i].Addr,
				Outpoint:           response.Events[i].Outpoint,
				Status:             response.Events[i].Status,
				ConfirmationHeight: response.Events[i].ConfirmationHeight,
			}
			err = dealEvent(tx, event)
			if err != nil {
				btlLog.CUST.Error("address even deal error:%v", event.Outpoint)
				tx.Rollback()
				tx = nil
				continue
			}
			tx.Commit()
		}
	}
}
