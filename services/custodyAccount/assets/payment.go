package assets

import (
	"fmt"
	"time"
	"trade/btlLog"
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
		if s.Queue.isEmpty() {
			continue
		}
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
			//TODO：创建交易：CREATE TX
			err = s.payToOutside(mission)
			//返回错误信息
			for index, _ := range mission.err {
				select {
				case mission.err[index] <- err:
				default:
				}
			}
		}
	}
}
func (s *AssetOutsideSever) payToOutside(mission *OutsideMission) error {
	var addr []string
	for _, a := range mission.AddrTarget {
		addr = append(addr, a.Addr)
	}
	assets, err := rpc.SendAssets(addr)
	if err != nil {
		btlLog.CUST.Error("rpc.SendAssets error:%w", err)
		return err
	}
	//TODO：处理交易结果
	fmt.Println(assets)
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
	return &AssetOutsideUniqueQueue{}
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

// AssetInSideSever 资产内部支付转账服务
type AssetInSideSever struct {
	Queue *AssetOutsideUniqueQueue
}

var InSideSever AssetInSideSever

func (s *AssetInSideSever) Start() {
	// Start 启动服务
	s.Queue = NewOutsideUniqueQueue()
	s.LoadMission()
	go s.runServer()
}
func (s *AssetInSideSever) runServer() {
	for {
		time.Sleep(60 * time.Second)

	}
}
func (s *AssetInSideSever) payToOutside(mission *OutsideMission) error {

	return nil
}
func (s *AssetInSideSever) LoadMission() {

}

// AssetInsideUniqueQueue 构建一个内部支付任务队列
type AssetInsideUniqueQueue struct {
	items   []*OutsideMission
	itemSet map[string]*OutsideMission
}

func NewInsideUniqueQueue() *AssetInsideUniqueQueue {
	return &AssetInsideUniqueQueue{}
}
func (q *AssetInsideUniqueQueue) addNewPkg(item *OutsideMission) bool {
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
func (q *AssetInsideUniqueQueue) getNextPkg() *OutsideMission {
	// 出队操作
	if len(q.items) == 0 {
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:]
	delete(q.itemSet, item.AssetID)
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
