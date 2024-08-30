package assets

// AssetOutsideSever 资产外部支付转账服务
type AssetOutsideSever struct {
	Queue *AssetOutsideUniqueQueue
}

var BtcSever AssetOutsideSever

func (m *AssetOutsideSever) Start() {
	// Start 启动服务
	m.Queue = NewUniqueQueue()
	m.LoadMission()
	go m.runServer()
}

func (m *AssetOutsideSever) runServer() {
	for {
		//TODO：获取可用资产列表
		//TODO:获取 一个外部支付任务
		//TODO：检查是否可交易：LIST AND UTXO
		//TODO：创建交易：CREATE TX
		//TODO：返回错误信息
	}
}
func (m *AssetOutsideSever) payToOutside(mission *OutsideMission) error {
	return nil
}
func (m *AssetOutsideSever) LoadMission() {}

// AssetOutsideUniqueQueue 构建一个外部支付任务队列
type AssetOutsideUniqueQueue struct {
	items   []*OutsideMission
	itemSet map[string]*OutsideMission
}

func NewUniqueQueue() *AssetOutsideUniqueQueue {
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
