package services

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"sync"
	"time"
	"trade/config"
	"trade/middleware"
	"trade/models"
	"trade/services/btldb"
	"trade/services/servicesrpc"
)

var ServerFee uint64 = 100

func GetServerFee() uint64 {
	return ServerFee
}

func PayServerFee(account *models.Account) error {
	acc, err := servicesrpc.AccountInfo(account.UserAccountCode)
	if err != nil {
		return err
	}
	// Change the escrow account balance
	_, err = servicesrpc.AccountUpdate(account.UserAccountCode, acc.CurrentBalance-int64(ServerFee), -1)
	if err != nil {
		return err
	}
	return nil
}

var CMutex sync.Mutex

// CreateCustodyAccount 创建托管账户并保持马卡龙文件
func CreateCustodyAccount(user *models.User) (*models.Account, error) {
	// Create a custody account based on user information
	account, macaroon, err := servicesrpc.AccountCreate(0, 0)
	if err != nil {
		CUST.Error(err.Error())
		return nil, err
	}
	// Build a macaroon storage path
	macaroonDir := config.GetConfig().ApiConfig.CustodyAccount.MacaroonDir
	if _, err = os.Stat(macaroonDir); os.IsNotExist(err) {
		err = os.MkdirAll(macaroonDir, os.ModePerm)
		if err != nil {
			CUST.Error(fmt.Sprintf("创建目标文件夹 %s 失败: %v\n", macaroonDir, err))
			return nil, err
		}
	}
	macaroonFile := filepath.Join(macaroonDir, account.Id+".macaroon")
	// Store macaroon information
	err = saveMacaroon(macaroon, macaroonFile)
	if err != nil {
		CUST.Error(err.Error())
		return nil, err
	}
	// Build an account object
	var accountModel models.Account
	accountModel.UserName = user.Username
	accountModel.UserId = user.ID
	accountModel.UserAccountCode = account.Id
	accountModel.Label = &account.Label
	accountModel.Status = 1
	// Write to the database
	CMutex.Lock()
	defer CMutex.Unlock()
	err = btldb.CreateAccount(&accountModel)
	if err != nil {
		CUST.Error(err.Error())
		return nil, err
	}
	// Return to the escrow account information
	return &accountModel, nil
}

// Update  托管账户更新
func UpdateCustodyAccount(account *models.Account, away models.BalanceAway, balance uint64, invoice string, HasServerFee bool) (uint, error) {
	var err error
	if account.UserAccountCode == "admin" {
		CUST.Info("admin receive %v amount", balance)
		return 0, nil
	}

	acc, err := servicesrpc.AccountInfo(account.UserAccountCode)
	if err != nil {
		return 0, err
	}
	var amount int64
	switch away {
	case models.AWAY_IN:
		amount = acc.CurrentBalance + int64(balance)
	case models.AWAY_OUT:
		amount = acc.CurrentBalance - int64(balance)
	default:
		return 0, fmt.Errorf("away error")
	}

	if amount < 0 {
		return 0, errors.New("balance not enough")
	}
	// Change the escrow account balance
	_, err = servicesrpc.AccountUpdate(account.UserAccountCode, amount, -1)
	if err != nil {
		return 0, err
	}

	// Build a database storage object
	ba := models.Balance{}
	ba.AccountId = account.ID
	ba.Amount = float64(balance)
	ba.Unit = models.UNIT_SATOSHIS
	ba.BillType = models.BILL_TYPE_PAYMENT
	ba.Away = away
	ba.State = models.STATE_SUCCESS
	ba.Invoice = nil
	ba.PaymentHash = nil
	if HasServerFee {
		err = PayServerFee(account)
		ba.ServerFee = GetServerFee()
	}
	if invoice != "" {
		ba.Invoice = &invoice
		i, _ := servicesrpc.InvoiceDecode(invoice)
		if i.PaymentHash != "" {
			ba.PaymentHash = &i.PaymentHash
		}
	}
	// Update the database
	CMutex.Lock()
	defer CMutex.Unlock()
	err = middleware.DB.Create(&ba).Error

	if err != nil {
		CUST.Error(err.Error())
		return 0, err
	}
	return ba.ID, nil
}

type ApplyRequest struct {
	Amount int64  `json:"amount"`
	Memo   string `json:"memo"`
}

// ApplyInvoice 使用指定账户申请一张发票
func ApplyInvoice(userId uint, account *models.Account, applyRequest *ApplyRequest) (*lnrpc.AddInvoiceResponse, error) {
	//获取马卡龙路径
	var macaroonFile string
	macaroonDir := config.GetConfig().ApiConfig.CustodyAccount.MacaroonDir

	if account.UserAccountCode == "admin" {
		macaroonFile = config.GetConfig().ApiConfig.Lnd.MacaroonPath
	} else {
		macaroonFile = filepath.Join(macaroonDir, account.UserAccountCode+".macaroon")
	}
	if macaroonFile == "" {
		CUST.Error("macaroon file not found")
		return nil, fmt.Errorf("macaroon file not found")
	}
	//调用Lit节点发票申请接口
	invoice, err := servicesrpc.InvoiceCreate(applyRequest.Amount, applyRequest.Memo, macaroonFile)
	if err != nil {
		CUST.Error(err.Error())
		return nil, err
	}
	//获取发票信息
	info, _ := servicesrpc.InvoiceFind(invoice.RHash)

	//构建invoice对象
	var invoiceModel models.Invoice
	invoiceModel.UserID = userId
	invoiceModel.Invoice = invoice.PaymentRequest
	invoiceModel.AccountID = &account.ID
	invoiceModel.Amount = float64(info.Value)

	invoiceModel.Status = models.InvoiceStatus(info.State)
	template := time.Unix(info.CreationDate, 0)
	invoiceModel.CreateDate = &template
	expiry := int(info.Expiry)
	invoiceModel.Expiry = &expiry

	//写入数据库
	CMutex.Lock()
	defer CMutex.Unlock()
	err = middleware.DB.Create(&invoiceModel).Error

	if err != nil {
		CUST.Error(err.Error())
		return invoice, err
	}
	return invoice, nil
}

type PayInvoiceRequest struct {
	Invoice  string `json:"invoice"`
	FeeLimit int64  `json:"feeLimit"`
}

// PayInvoice 使用指定账户支付发票
func PayInvoice(account *models.Account, PayInvoiceRequest *PayInvoiceRequest, HasServerFee bool) (bool, error) {
	//检查数据库中是否有该发票的记录
	a, err := GenericQueryByObject(&models.Balance{
		Invoice: &PayInvoiceRequest.Invoice,
	})
	if err != nil {
		CUST.Error(err.Error())
		return false, err
	}
	if len(a) > 0 {
		for _, v := range a {
			if v.State == models.STATE_SUCCESS {
				CUST.Info("该发票已支付")
				return false, fmt.Errorf("该发票已支付")
			}
			if v.State == models.STATE_UNKNOW {
				CUST.Info("该发票支付状态未知")
				return false, fmt.Errorf("该发票支付状态未知")
			}
		}
	}
	// 判断账户余额是否足够
	info, err := servicesrpc.InvoiceDecode(PayInvoiceRequest.Invoice)
	if err != nil {
		//CUST.Error("发票解析失败")
		return false, fmt.Errorf("发票解析失败(pay_request=%s)", PayInvoiceRequest.Invoice)
	}
	userBalance, err := servicesrpc.AccountInfo(account.UserAccountCode)
	if err != nil {
		CUST.Error("查询账户余额失败")
		return false, fmt.Errorf("查询账户余额失败（UserId=%d）", account.UserId)
	}
	if HasServerFee {
		if info.NumSatoshis+int64(GetServerFee()) > userBalance.CurrentBalance {
			CUST.Info("账户余额不足")
			return false, fmt.Errorf("账户余额不足（UserId=%d）", account.UserId)
		}
	} else {
		if info.NumSatoshis > userBalance.CurrentBalance {
			CUST.Info("账户余额不足")
			return false, fmt.Errorf("账户余额不足（UserId=%d）", account.UserId)
		}
	}

	//判断是否为节点内部转账
	i, err := GetInvoiceByReq(PayInvoiceRequest.Invoice)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		CUST.Error("数据库错误")
		return false, fmt.Errorf("数据库错误")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		_, err = PayAmountInside(account.UserId, i.UserID, uint64(info.NumSatoshis), 0, PayInvoiceRequest.Invoice, HasServerFee)
		if err != nil {
			CUST.Error("转账失败")
			return false, fmt.Errorf("转账失败")
		}
		i.Status = models.InvoiceStatusLocal
		err = UpdateInvoice(middleware.DB, i)
		if err != nil {
			CUST.Error("更新发票状态失败, invoice_id:%v", i.ID)
		}
		//更改发票状态
		h, _ := hex.DecodeString(info.PaymentHash)
		err = servicesrpc.InvoiceCancel(h)
		if err != nil {
			CUST.Error("取消发票失败")
		}
		CUST.Info("inside transfer success: paymenthash=%v,amount=%v", info.PaymentHash, info.NumSatoshis)
		return true, nil
	}
	//获取马卡龙路径
	var macaroonFile string
	macaroonDir := config.GetConfig().ApiConfig.CustodyAccount.MacaroonDir

	if account.UserAccountCode == "admin" {
		macaroonFile = config.GetConfig().ApiConfig.Lnd.MacaroonPath
	} else {
		macaroonFile = filepath.Join(macaroonDir, account.UserAccountCode+".macaroon")
	}
	if macaroonFile == "" {
		CUST.Error("macaroon file not found")
		return false, fmt.Errorf("macaroon file not found")
	}

	payment, err := servicesrpc.InvoicePay(macaroonFile, PayInvoiceRequest.Invoice, info.NumSatoshis, PayInvoiceRequest.FeeLimit)
	if err != nil {
		CUST.Error("pay invoice fail")
		return false, fmt.Errorf("pay invoice fail")
	}

	var balanceModel models.Balance
	balanceModel.AccountId = account.ID
	balanceModel.BillType = models.BILL_TYPE_PAYMENT
	balanceModel.Away = models.AWAY_OUT
	balanceModel.Amount = float64(payment.ValueSat)
	balanceModel.Unit = models.UNIT_SATOSHIS
	balanceModel.Invoice = &payment.PaymentRequest
	balanceModel.PaymentHash = &payment.PaymentHash
	if HasServerFee {
		err = PayServerFee(account)
		balanceModel.ServerFee = GetServerFee() + uint64(payment.FeeSat)
	}
	if payment.Status == lnrpc.Payment_SUCCEEDED {
		balanceModel.State = models.STATE_SUCCESS
	} else if payment.Status == lnrpc.Payment_FAILED {
		balanceModel.State = models.STATE_FAILED
	} else {
		balanceModel.State = models.STATE_UNKNOW
	}

	CMutex.Lock()
	defer CMutex.Unlock()
	err = middleware.DB.Create(&balanceModel).Error

	if err != nil {
		CUST.Error(err.Error())
		return false, err
	}
	return true, nil
}

// QueryAccountBalanceByUserId 通过用户ID查询账户余额
func QueryAccountBalanceByUserId(userId uint) (uint64, error) {
	// 查询账户
	account, err := btldb.ReadAccountByUserId(userId)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	// 查询账户余额
	userBalance, err := servicesrpc.AccountInfo(account.UserAccountCode)
	if err != nil {
		CUST.Error("Query failed: %s", err)
		return 0, err
	}
	return uint64(userBalance.CurrentBalance), nil
}

// QueryAccountBalanceByUsername 通过用户名查询账户余额
func QueryAccountBalanceByUsername(username string) (uint64, error) {
	user, err := btldb.ReadUserByUsername(username)
	if err != nil {
		return 0, err
	}
	return QueryAccountBalanceByUserId(user.ID)
}

// IsAccountBalanceEnoughByUserId  判断账户余额是否足够
func IsAccountBalanceEnoughByUserId(userId uint, value uint64) bool {
	balance, err := QueryAccountBalanceByUserId(userId)
	if err != nil {
		return false
	}
	return balance >= value
}

type InvoiceResponce struct {
	Invoice string               `json:"invoice"`
	AssetId string               `json:"asset_id"`
	Amount  int64                `json:"amount"`
	Status  models.InvoiceStatus `json:"status"`
}

// QueryInvoiceByUserId 查询用户发票
func QueryInvoiceByUserId(userId uint, assetId string) ([]InvoiceResponce, error) {
	params := QueryParams{
		"UserID":  userId,
		"AssetId": assetId,
	}
	a, err := GenericQuery(&models.Invoice{}, params)
	if err != nil {
		CUST.Error(err.Error())
		return nil, err
	}
	if len(a) > 0 {
		var invoices []InvoiceResponce
		for j := len(a) - 1; j >= 0; j-- {
			var i InvoiceResponce
			i.Invoice = a[j].Invoice
			i.AssetId = a[j].AssetId
			i.Amount = int64(a[j].Amount)
			i.Status = a[j].Status
			invoices = append(invoices, i)
		}
		return invoices, nil
	}
	return nil, nil
}

type PaymentRequest struct {
	AssetId string `json:"asset_id"`
}

type PaymentResponse struct {
	Timestamp int64               `json:"timestamp"`
	BillType  models.BalanceType  `json:"bill_type"`
	Away      models.BalanceAway  `json:"away"`
	Invoice   *string             `json:"invoice"`
	Amount    float64             `json:"amount"`
	AssetId   *string             `json:"asset_id"`
	State     models.BalanceState `json:"state"`
	Fee       uint64              `json:"fee"`
}

// QueryPaymentByUserId 查询用户支付记录
func QueryPaymentByUserId(userId uint, assetId string) ([]PaymentResponse, error) {
	accountId, err := btldb.ReadAccountByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("not find account info")
	}
	params := QueryParams{
		"AccountId": accountId.ID,
		"AssetId":   assetId,
	}
	a, err := GenericQuery(&models.Balance{}, params)
	if err != nil {
		CUST.Error(err.Error())
		return nil, fmt.Errorf("query payment error")
	}
	var results []PaymentResponse
	if len(a) > 0 {
		for i := len(a) - 1; i >= 0; i-- {
			v := a[i]
			r := PaymentResponse{}
			r.Timestamp = v.CreatedAt.Unix()
			r.BillType = v.BillType
			r.Away = v.Away
			r.Invoice = v.Invoice
			r.Amount = v.Amount
			btcAssetId := "00"
			r.AssetId = &btcAssetId
			r.State = v.State
			r.Fee = v.ServerFee
			results = append(results, r)
		}
	}
	return results, nil
}

// saveMacaroon 保存macaroon字节切片到指定文件
func saveMacaroon(macaroon []byte, macaroonFile string) error {
	file, err := os.OpenFile(macaroonFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	// 将字节切片写入指定位置
	data := macaroon
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

var PaymentMutex sync.Mutex

// PollPayment 遍历所有未确认的发票，轮询支付状态
func pollPayment() {

	PaymentMutex.Lock()
	defer PaymentMutex.Unlock()
	//查询数据库，获取所有未确认的支付
	params := QueryParams{
		"State": models.STATE_UNKNOW,
	}
	a, err := GenericQuery(&models.Balance{}, params)
	if err != nil {
		CUST.Error(err.Error())
		return
	}
	if len(a) > 0 {
		for _, v := range a {
			if v.Invoice == nil {
				continue
			}
			temp, err := servicesrpc.PaymentTrack(*v.PaymentHash)
			if err != nil {
				CUST.Warning(err.Error())
				continue
			}
			if temp.Status == lnrpc.Payment_SUCCEEDED {
				v.State = models.STATE_SUCCESS
				err = middleware.DB.Save(&v).Error
				if err != nil {
					CUST.Warning(err.Error())
				}
			} else if temp.Status == lnrpc.Payment_FAILED {
				v.State = models.STATE_FAILED
				err = middleware.DB.Save(&v).Error
				if err != nil {
					CUST.Warning(err.Error())
				}
			}

		}
	}
}

var InvoiceMutex sync.Mutex

// PollInvoice 遍历所有未支付的发票，轮询发票状态
func pollInvoice() {
	InvoiceMutex.Lock()
	defer InvoiceMutex.Unlock()
	//查询数据库，获取所有未支付的发票
	params := QueryParams{
		"Status": models.InvoiceStatusPending,
	}
	a, err := GenericQuery(&models.Invoice{}, params)
	if err != nil {
		CUST.Error(err.Error())
		return
	}
	if len(a) > 0 {
		for _, v := range a {
			invoice, err := servicesrpc.InvoiceDecode(v.Invoice)
			if err != nil {
				CUST.Warning(err.Error())
				continue
			}
			rHash, err := hex.DecodeString(invoice.PaymentHash)
			if err != nil {
				CUST.Warning(err.Error())
				continue
			}
			temp, err := servicesrpc.InvoiceFind(rHash)
			if err != nil {
				CUST.Warning(err.Error())
				continue
			}
			if int16(temp.State) != int16(v.Status) {
				v.Status = models.InvoiceStatus(temp.State)
				if v.Status == models.InvoiceStatusSuccess {
					ba := models.Balance{}
					ba.AccountId = *v.AccountID
					ba.Amount = v.Amount
					ba.Unit = models.UNIT_SATOSHIS
					ba.BillType = models.BILL_TYPE_RECHARGE
					ba.Away = models.AWAY_IN
					ba.State = models.STATE_SUCCESS
					ba.Invoice = &v.Invoice
					hash := hex.EncodeToString(rHash)
					ba.PaymentHash = &hash
					err = middleware.DB.Save(&ba).Error
					if err != nil {
						CUST.Warning(err.Error())
					}
				}
				err = middleware.DB.Save(&v).Error
				if err != nil {
					CUST.Warning(err.Error())
				}
			}
		}
	}
}

// PayAmountInside 内部转账比特币
func PayAmountInside(payUserId, receiveUserId uint, gasFee, serveFee uint64, invoice string, HasServerFee bool) (uint, error) {
	amount := gasFee + serveFee
	payAccount, err := btldb.ReadAccountByUserId(payUserId)
	if err != nil {
		CUST.Error("ReadAccountByUserId error:%v", err)
		return 0, err
	}
	outId, err := UpdateCustodyAccount(payAccount, models.AWAY_OUT, amount, invoice, HasServerFee)
	if err != nil {
		CUST.Error("UpdateCustodyAccount error(payUserId:%v):%v", payUserId, err)
		return 0, err
	}

	mark := func(Id uint, gasFee uint64, HasServerFee bool) {
		var fee uint64
		if HasServerFee {
			fee = GetServerFee()
		}
		remark := fmt.Sprintf("gasFee:%v ,serverFee:%v ,local: true", gasFee, fee)
		Ext := models.BalanceExt{
			BalanceId:   Id,
			BillExtDesc: &remark,
		}
		err = btldb.CreateBalanceExt(&Ext)
		if err != nil {
			CUST.Error("CreateBalanceExt error:%v", err)
		}
	}
	mark(outId, gasFee, HasServerFee)

	receiveAccount, err := btldb.ReadAccountByUserId(receiveUserId)
	if err != nil {
		CUST.Error("ReadAccountByUserId error:%v", err)
		return 0, err
	}
	_, err = UpdateCustodyAccount(receiveAccount, models.AWAY_IN, amount, invoice, false)
	if err != nil {
		CUST.Error("UpdateCustodyAccount error(receiveUserId:%v):%v", receiveUserId, err)
		return 0, err
	}
	return outId, nil
}

// CreatePayInsideMission 创建内部转账任务
func CreatePayInsideMission(payUserId, receiveUserId uint, gasFee, serveFee uint64, assetType string) (uint, error) {
	//获取支付账户信息
	payAccount, err := btldb.ReadAccountByUserId(payUserId)
	if err != nil {
		CUST.Error("Not find pay account info(UserId=%v):%v", payUserId, err)
		return 0, fmt.Errorf("not find pay account info")
	}
	//获取账户信息
	acc, err := servicesrpc.AccountInfo(payAccount.UserAccountCode)
	if err != nil {
		CUST.Error("AccountInfo error(UserId=%v):%v", payUserId, err)
		return 0, fmt.Errorf("AccountInfo error")
	}

	//检查账户余额是否足够
	if assetType == "00" {
		if acc.CurrentBalance < int64(gasFee) {
			CUST.Error("Account balance not enough(UserId=%v)", payUserId)
			return 0, fmt.Errorf("account balance not enough")
		}
	} else {
		return 0, fmt.Errorf("not support assetType")
	}

	//创建支付请求
	var (
		payReq         string
		payType        models.PayInsideType
		receiveAccount *models.Account
	)
	apply := ApplyRequest{
		Amount: int64(gasFee),
		Memo:   SetMemoSign(),
	}

	//检测目标账户是否合法
	switch receiveUserId {
	case AdminUserId:
		// 获取管理员账户信息
		receiveAccount = AdminAccount
		payType = models.PayInsideToAdmin
	default:
		//获取非管理员账户信息
		receiveAccount, err = btldb.ReadAccountByUserId(receiveUserId)
		if err != nil {
			CUST.Error("Not find receive account info(UserId=%v):%v", receiveUserId, err)
			return 0, fmt.Errorf("not find receive account info")
		}
		payType = models.PayInsideByInvioce
	}
	//创建发票
	invoice, err := ApplyInvoice(receiveUserId, receiveAccount, &apply)
	if err != nil {
		return 0, fmt.Errorf("apply userid = %v invoice error:%v", receiveUserId, err)
	}
	payReq = invoice.PaymentRequest
	//创建转账任务
	payInside := models.PayInside{
		PayUserId:     payUserId,
		GasFee:        gasFee,
		ServeFee:      serveFee,
		ReceiveUserId: receiveUserId,
		PayType:       payType,
		AssetType:     assetType,
		PayReq:        &payReq,
		Status:        models.PayInsideStatusPending,
	}
	//写入数据库
	err = btldb.CreatePayInside(&payInside)
	if err != nil {
		CUST.Error("CreatePayInside error:%v", err)
		return 0, err
	}
	return payInside.ID, nil
}

var PayInsideMutex sync.Mutex

// QueryPayInsideMission 处理内部转账任务
func pollPayInsideMission() {
	PayInsideMutex.Lock()
	defer PayInsideMutex.Unlock()
	//获取所有待处理任务
	params := QueryParams{
		"Status": models.PayInsideStatusPending,
	}
	a, err := GenericQuery(&models.PayInside{}, params)
	if err != nil {
		CUST.Error(err.Error())
		return
	}
	//	处理转账任务
	for _, v := range a {
		if v.AssetType == "00" {
			//获取支付账户信息
			payAccount, err := btldb.ReadAccountByUserId(v.PayUserId)
			if err != nil {
				CUST.Error("pollPayInsideMission find pay account error:%v", err)
				continue
			}
			//构建支付请求
			payReq := PayInvoiceRequest{
				Invoice:  *v.PayReq,
				FeeLimit: 0,
			}
			//检测是否有服务费
			HasServerFee := false
			if v.ServeFee > 0 {
				HasServerFee = true
			}
			//支付接口调用
			_, err = PayInvoice(payAccount, &payReq, HasServerFee)
			if err != nil {
				CUST.Error("pollPayInsideMission:%v", err)
				continue
			}
			//更新数据库状态
			v.Status = models.PayInsideStatusSuccess
			err = btldb.UpdatePayInside(v)
			if err != nil {
				CUST.Error("UpdatePayInside database error(id=%v):%v", v.ID, err)
				continue
			}
		}
	}
}

// CheckPayInsideStatus 检查内部转账任务状态是否成功
func CheckPayInsideStatus(id uint) (bool, error) {
	p, err := btldb.ReadPayInside(id)
	if err != nil {
		return false, err
	}
	if p.Status == models.PayInsideStatusSuccess {
		return true, nil
	}
	return false, nil
}

func SetMemoSign() string {
	return "internal transfer"
}

type LookupInvoiceRequest struct {
	InvoiceHash string `json:"invoice_hash"`
}

type LookupInvoiceResponse struct {
	Invoice *lnrpc.Invoice `json:"invoice"`
}

// LookupInvoice 查询发票状态
func LookupInvoice(req *LookupInvoiceRequest) (*LookupInvoiceResponse, error) {
	//查看发票状态
	balance := models.Balance{
		Away:        0,
		PaymentHash: &req.InvoiceHash,
		State:       models.STATE_SUCCESS,
	}
	err := middleware.DB.Where(balance).First(&balance).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		CUST.Error("LookupInvoice database error:%v", err.Error())
		return nil, err
	}

	//查找发票信息
	invoiceHash, err := hex.DecodeString(req.InvoiceHash)
	if err != nil {
		CUST.Error("Decode invoice hash error:%v", err.Error())
		return nil, err
	}
	invoice, err := servicesrpc.InvoiceFind(invoiceHash)
	if err != nil {
		CUST.Error("FindInvoice error:%v", err.Error())
		return nil, err
	}
	//返回结果
	result := LookupInvoiceResponse{
		Invoice: invoice,
	}
	return &result, nil
}
