package custodyAssets

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
	"trade/btlLog"
	"trade/models"
	"trade/services/btldb"
	rpc "trade/services/servicesrpc"
)

// AssetOutsideSever 资产外部支付转账服务
type AssetOutsideSever struct {
	Queue *AssetOutsideUniqueQueue
}

var OutsideSever AssetOutsideSever

func (s *AssetOutsideSever) Start(ctx context.Context) {
	// Start 启动服务
	s.Queue = NewOutsideUniqueQueue()
	s.LoadMission()
	go s.runServer(ctx)
}
func (s *AssetOutsideSever) runServer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(10 * time.Second)
			if s.Queue.isEmpty() {
				continue
			}
			//获取可用资产列表
			assets, err := rpc.ListAssets()
			if err != nil {
				btlLog.CUST.Error("rpc.ListAssets error:%v", err)
				continue
			}
			list := make(map[string]uint64)
			for _, asset := range assets.Assets {
				assetId := hex.EncodeToString(asset.AssetGenesis.AssetId)
				list[assetId] += asset.Amount
			}
			firstAssetID := ""
			for {
				if s.Queue.isEmpty() {
					break
				}
				//获取 一个外部支付任务
				mission := s.Queue.getNextPkg()
				if firstAssetID == "" {
					firstAssetID = mission.AssetID
				} else if firstAssetID == mission.AssetID {
					s.Queue.addNewPkg(mission)
					break
				}
				//检查是否可交易：LIST AND UTXO
				if list[mission.AssetID] < uint64(mission.TotalAmount) {
					s.Queue.addNewPkg(mission)
					break
				}
				balance, err := rpc.GetBalance()
				if err != nil {
					s.Queue.addNewPkg(mission)
					continue
				}
				if balance.AccountBalance["default"].ConfirmedBalance < int64(len(mission.AddrTarget)*1000) {
					s.Queue.addNewPkg(mission)
					continue
				}
				err = s.payToOutside(mission)
				if err == nil {
					btlLog.CUST.Info("payToOutside success: id=%v,amount=%v", mission.AssetID, mission.TotalAmount)
				}
				if err != nil {
					btlLog.CUST.Info("payToOutside Fail: id=%v,amount=%v", mission.AssetID, mission.TotalAmount)
					s.Queue.addNewPkg(mission)
				}
			}
		}
	}
}
func (s *AssetOutsideSever) payToOutside(mission *OutsideMission) error {
	var addr []string
	for _, a := range mission.AddrTarget {
		addr = append(addr, a.Mission.Address)
	}
	response, err := rpc.SendAssets(addr)
	if err != nil {
		btlLog.CUST.Error("rpc.SendAssets error:%v", err)
		return err
	}
	b := response.Transfer.AnchorTxHash
	for i := 0; i < len(b)/2; i++ {
		temp := b[i]
		b[i] = b[len(b)-i-1]
		b[len(b)-i-1] = temp
	}
	txId := hex.EncodeToString(b)
	tx := models.PayOutsideTx{
		TxHash:     txId,
		Timestamp:  response.Transfer.TransferTimestamp,
		HeightHint: response.Transfer.AnchorTxHeightHint,
		ChainFees:  response.Transfer.AnchorTxChainFees,
		InputsNum:  uint(len(response.Transfer.Inputs)),
		OutputsNum: uint(len(response.Transfer.Outputs)),
		Status:     models.PayOutsideStatusTXPending,
	}
	err = btldb.CreatePayOutsideTx(&tx)
	if err != nil {
		btlLog.CUST.Error("btldb.CreatePayOutsideTx error:%w", err)
	}
	for _, a := range mission.AddrTarget {
		a.Mission.TxHash = txId
		a.Mission.Status = models.PayOutsideStatusPaid
		err = btldb.UpdatePayOutside(a.Mission)
		if err != nil {
			btlLog.CUST.Error("btldb.UpdatePayOutside error:%w", err)
		}
		//更新Balance表
		balance, err := btldb.ReadBalance(a.Mission.BalanceId)
		if err != nil {
			continue
		}
		balance.State = models.STATE_SUCCESS
		balance.PaymentHash = &txId
		err = btldb.UpdateBalance(balance)
		if err != nil {
			btlLog.CUST.Error("payToOutside db error")
		}
	}
	return nil
}
func (s *AssetOutsideSever) LoadMission() {
	outsides, err := btldb.LoadPendingOutsides()
	if err != nil {
		return
	}
	for index, outside := range *outsides {
		m := OutsideMission{
			AddrTarget: []*target{
				{
					Mission: &(*outsides)[index],
				},
			},
			AssetID:     outside.AssetId,
			TotalAmount: int64(outside.Amount),
		}
		OutsideSever.Queue.addNewPkg(&m)
	}
}

// AssetOutsideUniqueQueue 构建一个外部支付任务队列
type AssetOutsideUniqueQueue struct {
	items   []*OutsideMission
	itemSet map[string]*OutsideMission
}

func NewOutsideUniqueQueue() *AssetOutsideUniqueQueue {
	return &AssetOutsideUniqueQueue{
		items:   []*OutsideMission{},
		itemSet: make(map[string]*OutsideMission),
	}
}
func (q *AssetOutsideUniqueQueue) addNewPkg(item *OutsideMission) bool {
	// addNewPkg 入队操作
	if i, exists := q.itemSet[item.AssetID]; exists {
		i.AddrTarget = append(i.AddrTarget, item.AddrTarget...)
		i.TotalAmount = i.TotalAmount + item.TotalAmount
		return true // 元素已存在，入队失败
	}
	q.items = append(q.items, item)
	q.itemSet[item.AssetID] = item
	return true
}
func (q *AssetOutsideUniqueQueue) getNextPkg() *OutsideMission {
	// 出队操作
	if len(q.items) == 0 {
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:]
	delete(q.itemSet, item.AssetID)
	return item
}
func (q *AssetOutsideUniqueQueue) isEmpty() bool {
	// 查看队列是否为空
	return len(q.items) == 0
}
func (q *AssetOutsideUniqueQueue) size() int {
	// 获取队列的长度
	return len(q.items)
}

// AssetInSideSever  TODO:  资产内部支付转账服务
type AssetInSideSever struct {
	Queue *AssetInsideUniqueQueue
}

var InSideSever AssetInSideSever

func (s *AssetInSideSever) Start(ctx context.Context) {
	// Start 启动服务
	s.Queue = NewInsideUniqueQueue()
	s.LoadMission()
	go s.runServer(ctx)
}
func (s *AssetInSideSever) runServer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if len(s.Queue.items) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}
			//取出队首元素
			mission := s.Queue.getNextPkg()
			if mission == nil {
				continue
			}
			//处理
			var err error
			err = s.payToInside(mission)
			if err != nil {
				btlLog.CUST.Info("payToInside fail: id=%v,amount=%v", mission.insideMission.AssetType, mission.insideMission.GasFee)
			}
			btlLog.CUST.Info("payToInside success: id=%v,amount=%v", mission.insideMission.AssetType, mission.insideMission.GasFee)
			select {
			case mission.err <- err:
			default:
			}
		}
	}
}
func (s *AssetInSideSever) NewMission(mission *isInsideMission) bool {
	return s.Queue.addNewPkg(mission)
}

func (s *AssetInSideSever) payToInside(mission *isInsideMission) error {
	switch mission.insideMission.PayType {
	case models.PayInsideByAddress:
		receiveAcc, _ := btldb.ReadAccountByUserId(mission.insideMission.ReceiveUserId)
		receiveBalance, err := btldb.GetAccountBalanceByGroup(receiveAcc.ID, mission.insideMission.AssetType)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ReadDbErr
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			receiveBalance = &models.AccountBalance{
				AccountID: receiveAcc.ID,
				AssetId:   mission.insideMission.AssetType,
				Amount:    float64(mission.insideMission.GasFee),
			}
		} else {
			receiveBalance.Amount += float64(mission.insideMission.GasFee)
		}
		err = btldb.UpdateAccountBalance(receiveBalance)
		if err != nil {
			return models.ReadDbErr
		}
		bill := models.Balance{
			AccountId:   receiveAcc.ID,
			BillType:    models.BillTypeAssetTransfer,
			Away:        models.AWAY_IN,
			Amount:      float64(mission.insideMission.GasFee),
			Unit:        models.UNIT_ASSET_NORMAL,
			ServerFee:   0,
			AssetId:     &mission.insideMission.AssetType,
			Invoice:     mission.insideMission.PayReq,
			PaymentHash: nil,
			State:       models.STATE_SUCCESS,
		}
		err = btldb.CreateBalance(&bill)
		if err != nil {
			return models.ReadDbErr
		}
		mission.insideMission.Status = models.PayInsideStatusSuccess
		err = btldb.UpdatePayInside(mission.insideMission)
		if err != nil {
			return models.ReadDbErr
		}
		return nil
	default:
		return fmt.Errorf("错误的内部转账类型:%v", mission.insideMission.PayType)
	}
}
func (s *AssetInSideSever) LoadMission() {
	//获取所有待处理任务
	params := btldb.QueryParams{
		"Status": models.PayInsideStatusPending,
	}
	a, err := btldb.GenericQuery(&models.PayInside{}, params)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return
	}
	//	处理转账任务
	for _, v := range a {
		if v.AssetType != "00" {
			i, err := btldb.GetInvoiceByReq(*v.PayReq)
			if err != nil {
				btlLog.CUST.Error("pollPayInsideMission find invoice error:%v", err)
				continue
			}
			mission := isInsideMission{
				isInside:      true,
				insideMission: v,
				insideInvoice: i,
			}
			//推送任务
			s.NewMission(&mission)
		}
	}
}

// AssetInsideUniqueQueue 构建一个内部支付任务队列
type AssetInsideUniqueQueue struct {
	items   []*isInsideMission
	itemSet map[uint]bool
}

func NewInsideUniqueQueue() *AssetInsideUniqueQueue {
	return &AssetInsideUniqueQueue{
		items:   []*isInsideMission{},
		itemSet: make(map[uint]bool),
	}
}
func (q *AssetInsideUniqueQueue) addNewPkg(item *isInsideMission) bool {
	// addNewPkg 入队操作
	if _, exists := q.itemSet[item.insideMission.ID]; exists {
		return false // 元素已存在，入队失败
	}
	q.items = append(q.items, item)
	q.itemSet[item.insideMission.ID] = true
	return true
}
func (q *AssetInsideUniqueQueue) getNextPkg() *isInsideMission {
	// 出队操作
	if len(q.items) == 0 {
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:]
	delete(q.itemSet, item.insideMission.ID)
	return item
}
func (q *AssetInsideUniqueQueue) isEmpty() bool {
	// 查看队列是否为空
	return len(q.items) == 0
}
func (q *AssetInsideUniqueQueue) size() int {
	// 获取队列的长度
	return len(q.items)
}
