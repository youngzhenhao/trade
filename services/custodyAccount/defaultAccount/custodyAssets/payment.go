package custodyAssets

//// AssetOutsideSever 资产外部支付转账服务
//type AssetOutsideSever struct {
//	Queue *AssetOutsideUniqueQueue
//}
//
//var OutsideSever AssetOutsideSever
//
//func (s *AssetOutsideSever) Start(ctx context.Context) {
//	// Start 启动服务
//	s.Queue = NewOutsideUniqueQueue()
//	s.LoadMission()
//	go s.runServer(ctx)
//}
//func (s *AssetOutsideSever) runServer(ctx context.Context) {
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		default:
//			time.Sleep(10 * time.Second)
//			if s.Queue.isEmpty() {
//				continue
//			}
//			//获取可用资产列表
//			assets, err := rpc.ListAssets()
//			if err != nil {
//				btlLog.CUST.Error("rpc.ListAssets error:%v", err)
//				continue
//			}
//			list := make(map[string]uint64)
//			for _, asset := range assets.Assets {
//				assetId := hex.EncodeToString(asset.AssetGenesis.AssetId)
//				list[assetId] += asset.Amount
//			}
//			firstAssetID := ""
//			for {
//				if s.Queue.isEmpty() {
//					break
//				}
//				//获取 一个外部支付任务
//				mission := s.Queue.getNextPkg()
//				if firstAssetID == "" {
//					firstAssetID = mission.AssetID
//				} else if firstAssetID == mission.AssetID {
//					s.Queue.addNewPkg(mission)
//					break
//				}
//				//检查是否可交易：LIST AND UTXO
//				if list[mission.AssetID] < uint64(mission.TotalAmount) {
//					s.Queue.addNewPkg(mission)
//					break
//				}
//				balance, err := rpc.GetBalance()
//				if err != nil {
//					s.Queue.addNewPkg(mission)
//					continue
//				}
//				if balance.AccountBalance["default"].ConfirmedBalance < int64(len(mission.AddrTarget)*1000) {
//					s.Queue.addNewPkg(mission)
//					continue
//				}
//				err = s.payToOutside(mission)
//				if err == nil {
//					btlLog.CUST.Info("payToOutside success: id=%v,amount=%v", mission.AssetID, mission.TotalAmount)
//				}
//				if err != nil {
//					btlLog.CUST.Info("payToOutside Fail: id=%v,amount=%v", mission.AssetID, mission.TotalAmount)
//					s.Queue.addNewPkg(mission)
//				}
//			}
//		}
//	}
//}
//func (s *AssetOutsideSever) payToOutside(mission *OutsideMission) error {
//	var addr []string
//	for _, a := range mission.AddrTarget {
//		addr = append(addr, a.Mission.Address)
//	}
//	response, err := rpc.SendAssets(addr)
//	if err != nil {
//		btlLog.CUST.Error("rpc.SendAssets error:%v", err)
//		return err
//	}
//	b := response.Transfer.AnchorTxHash
//	for i := 0; i < len(b)/2; i++ {
//		temp := b[i]
//		b[i] = b[len(b)-i-1]
//		b[len(b)-i-1] = temp
//	}
//	txId := hex.EncodeToString(b)
//	tx := custodyModels.PayOutsideTx{
//		TxHash:     txId,
//		Timestamp:  response.Transfer.TransferTimestamp,
//		HeightHint: response.Transfer.AnchorTxHeightHint,
//		ChainFees:  response.Transfer.AnchorTxChainFees,
//		InputsNum:  uint(len(response.Transfer.Inputs)),
//		OutputsNum: uint(len(response.Transfer.Outputs)),
//		Status:     custodyModels.PayOutsideStatusTXPending,
//	}
//	err = btldb.CreatePayOutsideTx(&tx)
//	if err != nil {
//		btlLog.CUST.Error("btldb.CreatePayOutsideTx error:%w", err)
//	}
//	for _, a := range mission.AddrTarget {
//		a.Mission.TxHash = txId
//		a.Mission.Status = custodyModels.PayOutsideStatusPaid
//		err = btldb.UpdatePayOutside(a.Mission)
//		if err != nil {
//			btlLog.CUST.Error("btldb.UpdatePayOutside error:%w", err)
//		}
//		//更新Balance表
//		balance, err := btldb.ReadBalance(a.Mission.BalanceId)
//		if err != nil {
//			continue
//		}
//		balance.State = models.STATE_SUCCESS
//		balance.PaymentHash = &txId
//		db := middleware.DB
//		err = btldb.UpdateBalance(db, balance)
//		if err != nil {
//			btlLog.CUST.Error("payToOutside db error")
//		}
//	}
//	return nil
//}
//func (s *AssetOutsideSever) LoadMission() {
//	outsides, err := btldb.LoadPendingOutsides()
//	if err != nil {
//		return
//	}
//	for index, outside := range *outsides {
//		m := OutsideMission{
//			AddrTarget: []*target{
//				{
//					Mission: &(*outsides)[index],
//				},
//			},
//			AssetID:     outside.AssetId,
//			TotalAmount: int64(outside.Amount),
//		}
//		OutsideSever.Queue.addNewPkg(&m)
//	}
//}
//
//// AssetOutsideUniqueQueue 构建一个外部支付任务队列
//type AssetOutsideUniqueQueue struct {
//	items   []*OutsideMission
//	itemSet map[string]*OutsideMission
//}
//
//func NewOutsideUniqueQueue() *AssetOutsideUniqueQueue {
//	return &AssetOutsideUniqueQueue{
//		items:   []*OutsideMission{},
//		itemSet: make(map[string]*OutsideMission),
//	}
//}
//func (q *AssetOutsideUniqueQueue) addNewPkg(item *OutsideMission) bool {
//	// addNewPkg 入队操作
//	if i, exists := q.itemSet[item.AssetID]; exists {
//		i.AddrTarget = append(i.AddrTarget, item.AddrTarget...)
//		i.TotalAmount = i.TotalAmount + item.TotalAmount
//		return true // 元素已存在，入队失败
//	}
//	q.items = append(q.items, item)
//	q.itemSet[item.AssetID] = item
//	return true
//}
//func (q *AssetOutsideUniqueQueue) getNextPkg() *OutsideMission {
//	// 出队操作
//	if len(q.items) == 0 {
//		return nil
//	}
//	item := q.items[0]
//	q.items = q.items[1:]
//	delete(q.itemSet, item.AssetID)
//	return item
//}
//func (q *AssetOutsideUniqueQueue) isEmpty() bool {
//	// 查看队列是否为空
//	return len(q.items) == 0
//}
//func (q *AssetOutsideUniqueQueue) size() int {
//	// 获取队列的长度
//	return len(q.items)
//}
