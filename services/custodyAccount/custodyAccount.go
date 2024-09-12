package custodyAccount

//
//type PaymentResponse struct {
//	Timestamp int64               `json:"timestamp"`
//	BillType  models.BalanceType  `json:"bill_type"`
//	Away      models.BalanceAway  `json:"away"`
//	Invoice   *string             `json:"invoice"`
//	Amount    float64             `json:"amount"`
//	AssetId   *string             `json:"asset_id"`
//	State     models.BalanceState `json:"state"`
//	Fee       uint64              `json:"fee"`
//}
//
//// QueryPaymentByUserId 查询用户支付记录
//func QueryPaymentByUserId(userId uint, assetId string) ([]PaymentResponse, error) {
//	accountId, err := btldb.ReadAccountByUserId(userId)
//	if err != nil {
//		return nil, fmt.Errorf("not find account info")
//	}
//	params := btldb.QueryParams{
//		"AccountId": accountId.ID,
//		"AssetId":   assetId,
//	}
//	a, err := btldb.GenericQuery(&models.Balance{}, params)
//	if err != nil {
//		btlLog.CUST.Error(err.Error())
//		return nil, fmt.Errorf("query payment error")
//	}
//	var results []PaymentResponse
//	if len(a) > 0 {
//		for i := len(a) - 1; i >= 0; i-- {
//			if a[i].State == models.STATE_FAILED {
//				continue
//			}
//			v := a[i]
//			r := PaymentResponse{}
//			r.Timestamp = v.CreatedAt.Unix()
//			r.BillType = v.BillType
//			r.Away = v.Away
//			r.Invoice = v.Invoice
//			r.Amount = v.Amount
//			btcAssetId := "00"
//			r.AssetId = &btcAssetId
//			r.State = v.State
//			r.Fee = v.ServerFee
//			results = append(results, r)
//		}
//	}
//	return results, nil
//}
//
//// saveMacaroon 保存macaroon字节切片到指定文件
//func saveMacaroon(macaroon []byte, macaroonFile string) error {
//	file, err := os.OpenFile(macaroonFile, os.O_RDWR|os.O_CREATE, 0644)
//	if err != nil {
//		panic(err)
//	}
//	defer func(file *os.File) {
//		err := file.Close()
//		if err != nil {
//			fmt.Println(err)
//		}
//	}(file)
//
//	// 将字节切片写入指定位置
//	data := macaroon
//	_, err = file.Write(data)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//var PaymentMutex sync.Mutex
//
//// PollPayment 遍历所有未确认的发票，轮询支付状态
//func PollPayment() {
//
//	PaymentMutex.Lock()
//	defer PaymentMutex.Unlock()
//	//查询数据库，获取所有未确认的支付
//	params := btldb.QueryParams{
//		"State": models.STATE_UNKNOW,
//	}
//	a, err := btldb.GenericQuery(&models.Balance{}, params)
//	if err != nil {
//		btlLog.CUST.Error(err.Error())
//		return
//	}
//	if len(a) > 0 {
//		for _, v := range a {
//			if v.Invoice == nil || *v.AssetId != "00" || v.PaymentHash == nil {
//				continue
//			}
//			temp, err := rpc.PaymentTrack(*v.PaymentHash)
//			if err != nil {
//				btlLog.CUST.Warning(err.Error())
//				continue
//			}
//			if temp.Status == lnrpc.Payment_SUCCEEDED {
//				v.State = models.STATE_SUCCESS
//				err = middleware.DB.Save(&v).Error
//				if err != nil {
//					btlLog.CUST.Warning(err.Error())
//				}
//			} else if temp.Status == lnrpc.Payment_FAILED {
//				v.State = models.STATE_FAILED
//				err = middleware.DB.Save(&v).Error
//				if err != nil {
//					btlLog.CUST.Warning(err.Error())
//				}
//			}
//		}
//	}
//}
//
//var InvoiceMutex sync.Mutex
//
//// PollInvoice 遍历所有未支付的发票，轮询发票状态
//func PollInvoice() {
//	InvoiceMutex.Lock()
//	defer InvoiceMutex.Unlock()
//	//查询数据库，获取所有未支付的发票
//	params := btldb.QueryParams{
//		"Status": models.InvoiceStatusPending,
//	}
//	a, err := btldb.GenericQuery(&models.Invoice{}, params)
//	if err != nil {
//		btlLog.CUST.Error(err.Error())
//		return
//	}
//	if len(a) > 0 {
//		for _, v := range a {
//			invoice, err := rpc.InvoiceDecode(v.Invoice)
//			if err != nil {
//				btlLog.CUST.Warning(err.Error())
//				continue
//			}
//			rHash, err := hex.DecodeString(invoice.PaymentHash)
//			if err != nil {
//				btlLog.CUST.Warning(err.Error())
//				continue
//			}
//			temp, err := rpc.InvoiceFind(rHash)
//			if err != nil {
//				btlLog.CUST.Warning(err.Error())
//				continue
//			}
//			if int16(temp.State) != int16(v.Status) {
//				v.Status = models.InvoiceStatus(temp.State)
//				if v.Status == models.InvoiceStatusSuccess {
//					ba := models.Balance{}
//					ba.AccountId = *v.AccountID
//					ba.Amount = v.Amount
//					ba.Unit = models.UNIT_SATOSHIS
//					ba.BillType = models.BillTypeRecharge
//					ba.Away = models.AWAY_IN
//					ba.State = models.STATE_SUCCESS
//					ba.Invoice = &v.Invoice
//					hash := hex.EncodeToString(rHash)
//					ba.PaymentHash = &hash
//					err = middleware.DB.Save(&ba).Error
//					if err != nil {
//						btlLog.CUST.Warning(err.Error())
//					}
//				}
//				err = middleware.DB.Save(&v).Error
//				if err != nil {
//					btlLog.CUST.Warning(err.Error())
//				}
//			}
//		}
//	}
//}
//
//// PayAmountInside 内部转账比特币
//func PayAmountInside(payUserId, receiveUserId uint, gasFee, serveFee uint64, invoice string, HasServerFee bool) (uint, error) {
//	amount := gasFee + serveFee
//	payAccount, err := btldb.ReadAccountByUserId(payUserId)
//	if err != nil {
//		btlLog.CUST.Error("ReadAccount error:%v", err)
//		return 0, err
//	}
//	outId, err := UpdateCustodyAccount(payAccount, models.AWAY_OUT, amount, invoice, HasServerFee)
//	if err != nil {
//		btlLog.CUST.Error("UpdateCustodyAccount error(payUserId:%v):%v", payUserId, err)
//		return 0, err
//	}
//
//	mark := func(Id uint, gasFee uint64, HasServerFee bool) {
//		var fee uint64
//		if HasServerFee {
//			fee = GetServerFee()
//		}
//		remark := fmt.Sprintf("gasFee:%v ,serverFee:%v ,local: true", gasFee, fee)
//		Ext := models.BalanceExt{
//			BalanceId:   Id,
//			BillExtDesc: &remark,
//		}
//		err = btldb.CreateBalanceExt(&Ext)
//		if err != nil {
//			btlLog.CUST.Error("CreateBalanceExt error:%v", err)
//		}
//	}
//	mark(outId, gasFee, HasServerFee)
//
//	receiveAccount, err := btldb.ReadAccountByUserId(receiveUserId)
//	if err != nil {
//		btlLog.CUST.Error("ReadAccount error:%v", err)
//		return 0, err
//	}
//	_, err = UpdateCustodyAccount(receiveAccount, models.AWAY_IN, amount, invoice, false)
//	if err != nil {
//		btlLog.CUST.Error("UpdateCustodyAccount error(receiveUserId:%v):%v", receiveUserId, err)
//		return 0, err
//	}
//	return outId, nil
//}
//
//// CreatePayInsideMission 创建内部转账任务
//func CreatePayInsideMission(payUserId, receiveUserId uint, gasFee, serveFee uint64, assetType string) (uint, error) {
//	//获取支付账户信息
//	payAccount, err := btldb.ReadAccountByUserId(payUserId)
//	if err != nil {
//		btlLog.CUST.Error("Not find pay account info(UserId=%v):%v", payUserId, err)
//		return 0, fmt.Errorf("not find pay account info")
//	}
//	//获取账户信息
//	acc, err := rpc.AccountInfo(payAccount.UserAccountCode)
//	if err != nil {
//		btlLog.CUST.Error("AccountInfo error(UserId=%v):%v", payUserId, err)
//		return 0, fmt.Errorf("AccountInfo error")
//	}
//
//	//检查账户余额是否足够
//	if assetType == "00" {
//		if acc.CurrentBalance < int64(gasFee) {
//			btlLog.CUST.Error("Account balance not enough(UserId=%v)", payUserId)
//			return 0, fmt.Errorf("account balance not enough")
//		}
//	} else {
//		return 0, fmt.Errorf("not support assetType")
//	}
//
//	//创建支付请求
//	var (
//		payReq         string
//		payType        models.PayInsideType
//		receiveAccount *models.Account
//	)
//	apply := ApplyRequest{
//		Amount: int64(gasFee),
//		Memo:   SetMemoSign(),
//	}
//
//	//检测目标账户是否合法
//	switch receiveUserId {
//	case AdminUserId:
//		// 获取管理员账户信息
//		receiveAccount = AdminAccount
//		payType = models.PayInsideToAdmin
//	default:
//		//获取非管理员账户信息
//		receiveAccount, err = btldb.ReadAccountByUserId(receiveUserId)
//		if err != nil {
//			btlLog.CUST.Error("Not find receive account info(UserId=%v):%v", receiveUserId, err)
//			return 0, fmt.Errorf("not find receive account info")
//		}
//		payType = models.PayInsideByInvoice
//	}
//	//创建发票
//	invoice, err := ApplyInvoice(receiveUserId, receiveAccount, &apply)
//	if err != nil {
//		return 0, fmt.Errorf("apply userid = %v invoice error:%v", receiveUserId, err)
//	}
//	payReq = invoice.PaymentRequest
//	//创建转账任务
//	payInside := models.PayInside{
//		PayUserId:     payUserId,
//		GasFee:        gasFee,
//		ServeFee:      serveFee,
//		ReceiveUserId: receiveUserId,
//		PayType:       payType,
//		AssetType:     assetType,
//		PayReq:        &payReq,
//		Status:        models.PayInsideStatusPending,
//	}
//	//写入数据库
//	err = btldb.CreatePayInside(&payInside)
//	if err != nil {
//		btlLog.CUST.Error("CreatePayInside error:%v", err)
//		return 0, err
//	}
//	return payInside.ID, nil
//}
//
//var PayInsideMutex sync.Mutex
//
//// QueryPayInsideMission 处理内部转账任务
//func PollPayInsideMission() {
//	PayInsideMutex.Lock()
//	defer PayInsideMutex.Unlock()
//	//获取所有待处理任务
//	params := btldb.QueryParams{
//		"Status": models.PayInsideStatusPending,
//	}
//	a, err := btldb.GenericQuery(&models.PayInside{}, params)
//	if err != nil {
//		btlLog.CUST.Error(err.Error())
//		return
//	}
//	//	处理转账任务
//	for _, v := range a {
//		if v.AssetType == "00" {
//			//获取支付账户信息
//			payAccount, err := btldb.ReadAccountByUserId(v.PayUserId)
//			if err != nil {
//				btlLog.CUST.Error("pollPayInsideMission find pay account error:%v", err)
//				continue
//			}
//			//构建支付请求
//			payReq := PayInvoiceRequest{
//				Invoice:  *v.PayReq,
//				FeeLimit: 0,
//			}
//			//检测是否有服务费
//			HasServerFee := false
//			if v.ServeFee > 0 {
//				HasServerFee = true
//			}
//			//支付接口调用
//			balanceId, err := PayInvoice(payAccount, &payReq, HasServerFee)
//			if err != nil {
//				btlLog.CUST.Error("pollPayInsideMission:%v", err)
//				v.Status = models.PayInsideStatusFailed
//			} else {
//				v.Status = models.PayInsideStatusSuccess
//			}
//			//更新数据库状态
//			v.BalanceId = balanceId
//			err = btldb.UpdatePayInside(v)
//			if err != nil {
//				btlLog.CUST.Error("UpdatePayInside database error(id=%v):%v", v.ID, err)
//				continue
//			}
//		}
//	}
//}
//
//// CheckPayInsideStatus 检查内部转账任务状态是否成功
//func CheckPayInsideStatus(id uint) (bool, error) {
//	p, err := btldb.ReadPayInside(id)
//	if err != nil {
//		return false, err
//	}
//	switch p.Status {
//	case models.PayInsideStatusSuccess:
//		return true, nil
//	case models.PayInsideStatusFailed:
//		return false, models.CustodyAccountPayInsideMissionFaild
//	default:
//		return false, models.CustodyAccountPayInsideMissionPending
//	}
//}
//
//func SetMemoSign() string {
//	return "internal transfer"
//}
//
//type LookupInvoiceRequest struct {
//	InvoiceHash string `json:"invoice_hash"`
//}
//
//type LookupInvoiceResponse struct {
//	Invoice *lnrpc.Invoice `json:"invoice"`
//}
//
//// LookupInvoice 查询发票状态
//func LookupInvoice(req *LookupInvoiceRequest) (*LookupInvoiceResponse, error) {
//	//查看发票状态
//	balance := models.Balance{
//		Away:        0,
//		PaymentHash: &req.InvoiceHash,
//		State:       models.STATE_SUCCESS,
//	}
//	err := middleware.DB.Where(balance).First(&balance).Error
//	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		btlLog.CUST.Error("LookupInvoice database error:%v", err.Error())
//		return nil, err
//	}
//
//	//查找发票信息
//	invoiceHash, err := hex.DecodeString(req.InvoiceHash)
//	if err != nil {
//		btlLog.CUST.Error("Decode invoice hash error:%v", err.Error())
//		return nil, err
//	}
//	invoice, err := rpc.InvoiceFind(invoiceHash)
//	if err != nil {
//		btlLog.CUST.Error("FindInvoice error:%v", err.Error())
//		return nil, err
//	}
//	//返回结果
//	result := LookupInvoiceResponse{
//		Invoice: invoice,
//	}
//	return &result, nil
//}
