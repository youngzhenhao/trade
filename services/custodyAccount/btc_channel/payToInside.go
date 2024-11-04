package btc_channel

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"
	"trade/btlLog"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
	"trade/services/btldb"
	caccount "trade/services/custodyAccount/account"
	"trade/services/custodyAccount/custodyBase/custodyFee"
	"trade/services/custodyAccount/custodyBase/custodyLimit"
	"trade/services/custodyAccount/custodyBase/custodyRpc"
	rpc "trade/services/servicesrpc"
)

// BTCPayInsideSever btc支付内部转账服务
type BTCPayInsideSever struct {
	Queue *BTCPayInsideUniqueQueue
}

var BtcSever BTCPayInsideSever

func (m *BTCPayInsideSever) Start(ctx context.Context) {
	// Start 启动服务
	m.Queue = NewUniqueQueue()
	m.LoadMission()
	go m.runServer(ctx)
}
func (m *BTCPayInsideSever) runServer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("收到结束信号，退出循环")
			return

		default:
			if len(m.Queue.items) == 0 {
				time.Sleep(5 * time.Second)
				continue
			}
			//取出队首元素
			mission := m.Queue.getNextPkg()
			if mission == nil {
				continue
			}
			//处理
			var err error
			if mission.insideInvoice.Status != models.InvoiceStatusPending {
				err = fmt.Errorf("invoice is close")
			} else {
				err = m.payToInside(mission)
			}
			if err != nil {
				mission.insideMission.Status = models.PayInsideStatusFailed
			} else {
				mission.insideMission.Status = models.PayInsideStatusSuccess
				CloseInvoice(mission.insideInvoice)
				btlLog.CUST.Info("inside transfer success: id=%v,amount=%v", mission.insideMission.ID, mission.insideMission.GasFee)
			}
			select {
			case mission.err <- err:
			default:
			}
			err = btldb.UpdatePayInside(mission.insideMission)
			if err != nil {
				btlLog.CUST.Error("更新内部转账记录失败, mission_id:%v，error:%v", mission.insideMission.ID, err)
			}
		}
	}
}
func (m *BTCPayInsideSever) NewMission(mission *isInsideMission) bool {
	return m.Queue.addNewPkg(mission)
}
func (m *BTCPayInsideSever) payToInside(mission *isInsideMission) error {
	//todo 内部转账逻辑需要重写
	if mission.insideInvoice.Status != models.InvoiceStatusPending {
		return nil
	}
	var payToAdmin bool
	fee := custodyFee.ChannelBtcInsideServiceFee
	switch mission.insideMission.PayType {
	case models.PayInsideToAdmin, models.FairLunchFee, models.ChannelBTCFee, models.ChannelBTCOutSideFee:
		payToAdmin = true
		fee = 0
	default:
	}
	amount := mission.insideMission.GasFee + mission.insideMission.ServeFee
	//变更付款方账户
	payAcc, err := caccount.GetUserInfoById(mission.insideMission.PayUserId)
	if err != nil {
		btlLog.CUST.Error("获取账户信息失败, mission_id:%v，error:%v", mission.insideMission.ID, err)
		return fmt.Errorf("获取账户信息失败")
	}
	balanceId, err := updateCustodyAccount(payAcc, models.AWAY_OUT, amount, mission.insideInvoice.Invoice, fee)
	if err != nil {
		btlLog.CUST.Error("内部付款方账户更新失败, mission_id:%v，error:%v", mission.insideMission.ID, err)
		return fmt.Errorf("付款失败")
	}
	//限制额度更新
	limitType := custodyModels.LimitType{
		AssetId:      "00",
		TransferType: custodyModels.LimitTransferTypeLocal,
	}
	err = custodyLimit.MinusLimit(middleware.DB, payAcc, &limitType, float64(amount+fee))
	if err != nil {
		btlLog.CUST.Error("额度限制未正常更新:%s", err.Error())
		btlLog.CUST.Error("error PayInsideId:%v", mission.insideMission.ID)
	}

	mission.insideMission.BalanceId = balanceId
	//变更收款方账户，如果是内部转账给管理员，则跳过
	if !payToAdmin {
		revAcc, err := caccount.GetUserInfoById(mission.insideMission.ReceiveUserId)
		if err != nil {
			btlLog.CUST.Error("获取收款账户信息失败, mission_id:%v，error:%v", mission.insideMission.ReceiveUserId, err)
			return nil
		}
		_, err = updateCustodyAccount(revAcc, models.AWAY_IN, amount, mission.insideInvoice.Invoice, 0)
		if err != nil {
			btlLog.CUST.Error("内部付收款方账户更新失败, mission_id:%v，error:%v", mission.insideMission.ReceiveUserId, err)
			return nil
		}
	}
	return nil
}
func (m *BTCPayInsideSever) LoadMission() {
	//获取所有待处理任务
	params := btldb.QueryParams{
		"AssetType": "00",
		"Status":    models.PayInsideStatusPending,
	}
	a, err := btldb.GenericQuery(&models.PayInside{}, params)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return
	}
	//	处理转账任务
	for _, v := range a {
		if v.AssetType == "00" {
			i, err := btldb.GetInvoiceByReq(*v.PayReq)
			if err != nil {
				btlLog.CUST.Error("pollPayInsideMission find invoice error:%v", err)
				continue
			}
			mission := isInsideMission{
				isInside:      true,
				insideMission: v,
				insideInvoice: i,
				err:           make(chan error, 1),
			}
			go func(c chan error) {
				err := <-c
				if err != nil {
					btlLog.CUST.Error("btc sendPayment timeout:%s", err.Error())
				}
				close(c)
			}(mission.err)
			//推送任务
			m.NewMission(&mission)
		}
	}
}
func updateCustodyAccount(usr *caccount.UserInfo, away models.BalanceAway, balance uint64, invoice string, ServerFee uint64) (uint, error) {
	var err error
	var updateAway string
	switch away {
	case models.AWAY_IN:
		updateAway = custodyRpc.UpdateBalancePlus
	case models.AWAY_OUT:
		updateAway = custodyRpc.UpdateBalanceMinus
	default:
		return 0, fmt.Errorf("away error")
	}
	if balance <= 0 {
		return 0, fmt.Errorf("balance error")
	}
	// Change the escrow usr balance
	_, err = custodyRpc.UpdateBalance(usr, updateAway, int64(balance))

	// Build a database storage object
	ba := models.Balance{}
	ba.AccountId = usr.Account.ID
	ba.Amount = float64(balance)
	ba.Unit = models.UNIT_SATOSHIS
	ba.BillType = models.BillTypePayment
	ba.Away = away
	if err != nil {
		ba.State = models.STATE_FAILED
	} else {
		ba.State = models.STATE_SUCCESS
	}
	ba.Invoice = nil
	ba.PaymentHash = nil
	//	计算服务费
	ba.ServerFee = ServerFee
	if invoice != "" {
		i, _ := rpc.InvoiceDecode(invoice)
		if i.PaymentHash != "" {
			ba.PaymentHash = &i.PaymentHash
		}
	}
	ba.Invoice = &invoice
	// Update the database
	dbErr := btldb.CreateBalance(&ba)
	if dbErr != nil {
		btlLog.CUST.Error(dbErr.Error())
		return 0, nil
	}
	if ServerFee > 0 {
		err = custodyFee.PayServiceFeeSync(usr, ServerFee, ba.ID, models.ChannelBTCFee, "payToInside Fee")
	}
	return ba.ID, nil
}
func CloseInvoice(invoice *models.Invoice) {
	invoice.Status = models.InvoiceStatusLocal
	err := btldb.UpdateInvoice(middleware.DB, invoice)
	if err != nil {
		btlLog.CUST.Error("更新发票状态失败, invoice_id:%v", invoice.ID)
	}
	DecodePayReq, err := rpc.InvoiceDecode(invoice.Invoice)
	if err != nil {
		btlLog.CUST.Error("发票解析失败", err)
	}
	h, _ := hex.DecodeString(DecodePayReq.PaymentHash)
	err = rpc.InvoiceCancel(h)
	if err != nil {
		btlLog.CUST.Error("取消发票失败")
	}
}

// BTCPayInsideUniqueQueue 构建一个任务队列
type BTCPayInsideUniqueQueue struct {
	items   []*isInsideMission
	itemSet map[uint]bool
}

func NewUniqueQueue() *BTCPayInsideUniqueQueue {
	return &BTCPayInsideUniqueQueue{
		items:   []*isInsideMission{},
		itemSet: make(map[uint]bool),
	}
}
func (q *BTCPayInsideUniqueQueue) addNewPkg(item *isInsideMission) bool {
	// addNewPkg 入队操作
	if _, exists := q.itemSet[item.insideMission.ID]; exists {
		return false // 元素已存在，入队失败
	}
	q.items = append(q.items, item)
	q.itemSet[item.insideMission.ID] = true
	return true
}
func (q *BTCPayInsideUniqueQueue) getNextPkg() *isInsideMission {
	// 出队操作
	if len(q.items) == 0 {
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:]
	delete(q.itemSet, item.insideMission.ID)
	return item
}
func (q *BTCPayInsideUniqueQueue) isEmpty() bool {
	// 查看队列是否为空
	return len(q.items) == 0
}
func (q *BTCPayInsideUniqueQueue) size() int {
	// 获取队列的长度
	return len(q.items)
}
