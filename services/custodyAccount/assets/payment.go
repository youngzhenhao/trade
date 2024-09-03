package assets

import (
	"encoding/hex"
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

func (s *AssetOutsideSever) Start() {
	// Start 启动服务
	s.Queue = NewOutsideUniqueQueue()
	s.LoadMission()
	go s.runServer()
}
func (s *AssetOutsideSever) runServer() {
	for {
		time.Sleep(60 * time.Second)
		//获取可用资产列表
		assets, err := rpc.ListAssets()
		if err != nil {
			btlLog.CUST.Error("rpc.ListAssets error:%w", err)
			continue
		}
		list := make(map[string]uint64)
		for _, asset := range assets.Assets {
			list[string(asset.AssetGenesis.AssetId)] += asset.Amount
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
				continue
			}
			balance, err := rpc.GetBalance()
			if err != nil {
				continue
			}
			if balance.AccountBalance["default"].ConfirmedBalance < int64(len(mission.AddrTarget)*1000) {
				continue
			}
			err = s.payToOutside(mission)
			if err == nil {
				btlLog.CUST.Info("payToOutside success: id=%v,amount=%v", mission.AssetID, mission.TotalAmount)
			}
			//返回错误信息
			for index, _ := range mission.err {
				select {
				case _, ok := <-mission.err[index]:
					if !ok {
						continue
					} else {
						mission.err[index] <- err
					}
				default:
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
	txId := hex.EncodeToString(response.Transfer.AnchorTxHash)
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
		//todo：扣除手续费

		//更新Balance表
		balance, err := btldb.ReadBalance(a.Mission.BalanceId)
		if err != nil {
			continue
		}
		balance.State = models.STATE_SUCCESS
		err = btldb.UpdateBalance(balance)
		if err != nil {
			btlLog.CUST.Error("payToOutside db error")
		}
	}
	return nil
}
func (s *AssetOutsideSever) LoadMission() {

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
		i.err = append(i.err, item.err...)
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

// TODO: AssetInSideSever 资产内部支付转账服务
type AssetInSideSever struct {
	Queue *AssetInsideUniqueQueue
}

var InSideSever AssetInSideSever

func (s *AssetInSideSever) Start() {
	// Start 启动服务
	s.Queue = NewInsideUniqueQueue()
	s.LoadMission()
	go s.runServer()
}
func (s *AssetInSideSever) runServer() {
	for {
		if len(s.Queue.items) == 0 {
			time.Sleep(5 * time.Second)
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
			mission.insideMission.Status = models.PayInsideStatusFailed
		} else {
			mission.insideMission.Status = models.PayInsideStatusSuccess
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
func (s *AssetInSideSever) payToInside(mission *isInsideMission) error {

	return nil
}
func (s *AssetInSideSever) LoadMission() {

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
