package services

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/lightninglabs/lightning-terminal/litrpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"sync"
	"time"
	"trade/config"
	"trade/middleware"
	"trade/models"
	"trade/services/servicesrpc"
)

var mutex sync.Mutex

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
	mutex.Lock()
	defer mutex.Unlock()
	err = CreateAccount(&accountModel)
	if err != nil {
		CUST.Error(err.Error())
		return nil, err
	}
	// Return to the escrow account information
	return &accountModel, nil
}

// Update  托管账户更新
func UpdateCustodyAccount(account *models.Account, away models.BalanceAway, balance uint64, invoice string) (uint, error) {
	var err error
	if account.UserAccountCode != "admin" {
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
	if invoice != "" {
		ba.Invoice = &invoice
		i, _ := DecodeInvoice(invoice)
		if i.PaymentHash != "" {
			ba.PaymentHash = &i.PaymentHash
		}
	}
	// Update the database
	mutex.Lock()
	defer mutex.Unlock()
	err = middleware.DB.Create(&ba).Error
	if err != nil {
		CUST.Error(err.Error())
		return 0, err
	}
	return ba.ID, nil
}

func PayAmountInside(payUserId, receiveUserId uint, gasFee, serveFee uint64, invoice string) (uint, error) {
	amount := gasFee + serveFee
	payAccount, err := ReadAccountByUserId(payUserId)
	if err != nil {
		CUST.Error("ReadAccountByUserId error:%v", err)
		return 0, err
	}
	outId, err := UpdateCustodyAccount(payAccount, models.AWAY_OUT, amount, invoice)
	if err != nil {
		CUST.Error("UpdateCustodyAccount error(payUserId:%v):%v", payUserId, err)
		return 0, err
	}
	remark := fmt.Sprintf("gasFee:%v ,serverFee:%v", gasFee, serveFee)
	Ext := models.BalanceExt{
		BalanceId:   outId,
		BillExtDesc: &remark,
	}
	err = CreateBalanceExt(&Ext)
	if err != nil {
		CUST.Error("CreateBalanceExt error:%v", err)
	}

	receiveAccount, err := ReadAccountByUserId(receiveUserId)
	if err != nil {
		CUST.Error("ReadAccountByUserId error:%v", err)
		return 0, err
	}
	Id, err := UpdateCustodyAccount(receiveAccount, models.AWAY_IN, amount, invoice)
	if err != nil {
		CUST.Error("UpdateCustodyAccount error(receiveUserId:%v):%v", receiveUserId, err)
		return 0, err
	}
	return Id, nil
}

// QueryCustodyAccount  托管账户查询
func QueryCustodyAccount(accountCode string) (*litrpc.Account, error) {
	return servicesrpc.AccountInfo(accountCode)
}

// DeleteCustodianAccount 托管账户删除
func DeleteCustodianAccount() error {
	//TODO: 获取托管账户ID
	id := "test"
	//删除Lit节点托管账户
	err := servicesrpc.AccountRemove(id)
	//TODO: 更新数据库相关信息

	//TODO: 返回删除结果

	return err
}

type ApplyRequest struct {
	Amount int64  `json:"amount"`
	Memo   string `json:"memo"`
}

// ApplyInvoice 使用指定账户申请一张发票
func ApplyInvoice(user *models.User, account *models.Account, applyRequest *ApplyRequest) (*lnrpc.AddInvoiceResponse, error) {
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
	info, _ := FindInvoice(invoice.RHash)

	//构建invoice对象
	var invoiceModel models.Invoice
	invoiceModel.UserID = user.ID
	invoiceModel.Invoice = invoice.PaymentRequest
	invoiceModel.AccountID = &account.ID
	invoiceModel.Amount = float64(info.Value)

	invoiceModel.Status = int16(info.State)
	template := time.Unix(info.CreationDate, 0)
	invoiceModel.CreateDate = &template
	expiry := int(info.Expiry)
	invoiceModel.Expiry = &expiry

	//写入数据库
	mutex.Lock()
	defer mutex.Unlock()
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
func PayInvoice(account *models.Account, PayInvoiceRequest *PayInvoiceRequest) (bool, error) {
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
	info, err := DecodeInvoice(PayInvoiceRequest.Invoice)
	if err != nil {
		CUST.Error("发票解析失败")
		return false, fmt.Errorf("发票解析失败")
	}

	userBalance, err := QueryCustodyAccount(account.UserAccountCode)
	if err != nil {
		CUST.Error("查询账户余额失败")
		return false, fmt.Errorf("查询账户余额失败")
	}
	if info.NumSatoshis > userBalance.CurrentBalance {
		CUST.Info("账户余额不足")
		return false, fmt.Errorf("账户余额不足")
	}

	//判断是否为节点内部转账
	i, err := GetInvoiceByReq(PayInvoiceRequest.Invoice)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		CUST.Error("数据库错误")
		return false, fmt.Errorf("数据库错误")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		_, err = PayAmountInside(account.UserId, i.UserID, uint64(info.NumSatoshis), 0, PayInvoiceRequest.Invoice)
		if err != nil {
			CUST.Error("转账失败")
			return false, fmt.Errorf("转账失败")
		}
		i.Status = 1
		err = UpdateInvoice(middleware.DB, i)
		if err != nil {
			CUST.Error("更新发票状态失败, invoice_id:%v", i.ID)
		}
		//更改发票状态
		h, _ := hex.DecodeString(info.PaymentHash)
		err = CancelInvoice(h)
		if err != nil {
			CUST.Error("取消发票失败")
		}
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

	payment, err := servicesrpc.InvoicePay(macaroonFile, PayInvoiceRequest.Invoice, PayInvoiceRequest.FeeLimit)
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
	if payment.Status == lnrpc.Payment_SUCCEEDED {
		balanceModel.State = models.STATE_SUCCESS
	} else if payment.Status == lnrpc.Payment_FAILED {
		balanceModel.State = models.STATE_FAILED
	} else {
		balanceModel.State = models.STATE_UNKNOW
	}
	mutex.Lock()
	defer mutex.Unlock()
	err = middleware.DB.Create(&balanceModel).Error
	if err != nil {
		CUST.Error(err.Error())
		return false, err
	}
	return true, nil
}

// QueryAccountBalanceByUserId 查询用户账户余额
func QueryAccountBalanceByUserId(userId uint) (uint64, error) {
	// 查询账户
	account, err := ReadAccountByUserId(userId)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	// 查询账户余额
	userBalance, err := QueryCustodyAccount(account.UserAccountCode)
	if err != nil {
		CUST.Error("Query failed: %s", err)
		return 0, err
	}
	return uint64(userBalance.CurrentBalance), nil
}

type InvoiceResponce struct {
	Invoice string `json:"invoice"`
	AssetId string `json:"asset_id"`
	Amount  int64  `json:"amount"`
	Status  int16  `json:"status"`
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
		for _, v := range a {
			var i InvoiceResponce
			i.Invoice = v.Invoice
			i.AssetId = v.AssetId
			i.Amount = int64(v.Amount)
			i.Status = v.Status
			invoices = append(invoices, i)
		}
		return invoices, nil
	}
	return nil, nil

}

// QueryPaymentByUserId 查询用户支付记录
func QueryPaymentByUserId(userId uint, assetId string) {

}

// DecodeInvoice  解析发票信息
func DecodeInvoice(invoice string) (*lnrpc.PayReq, error) {
	return servicesrpc.InvoiceDecode(invoice)
}

// FindInvoice 查询节点内部发票
func FindInvoice(rHash []byte) (*lnrpc.Invoice, error) {
	return servicesrpc.InvoiceFind(rHash)
}

// CancelInvoice 取消发票
func CancelInvoice(hash []byte) error {
	return servicesrpc.InvoiceCancel(hash)
}

// TrackPayment 跟踪支付状态
func TrackPayment(paymentHash string) (*lnrpc.Payment, error) {
	return servicesrpc.PaymentTrack(paymentHash)
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

// PollPayment 遍历所有未确认的发票，轮询支付状态
func pollPayment() {
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
			temp, err := TrackPayment(*v.PaymentHash)
			if err != nil {
				CUST.Warning(err.Error())
				continue
			}
			if temp.Status == lnrpc.Payment_SUCCEEDED {
				v.State = models.STATE_SUCCESS
				mutex.Lock()
				defer mutex.Unlock()
				err = middleware.DB.Save(&v).Error
				if err != nil {
					CUST.Warning(err.Error())
				}
			} else if temp.Status == lnrpc.Payment_FAILED {
				v.State = models.STATE_FAILED
				mutex.Lock()
				defer mutex.Unlock()
				err = middleware.DB.Save(&v).Error
				if err != nil {
					CUST.Warning(err.Error())
				}
			}

		}
	}
}

// PollInvoice 遍历所有未支付的发票，轮询发票状态
func pollInvoice() {
	//查询数据库，获取所有未支付的发票
	params := QueryParams{
		"Status": lnrpc.Invoice_OPEN,
	}
	a, err := GenericQuery(&models.Invoice{}, params)
	if err != nil {
		CUST.Error(err.Error())
		return
	}
	if len(a) > 0 {
		for _, v := range a {
			invoice, err := DecodeInvoice(v.Invoice)
			if err != nil {
				CUST.Warning(err.Error())
				continue
			}
			rHash, err := hex.DecodeString(invoice.PaymentHash)
			if err != nil {
				CUST.Warning(err.Error())
				continue
			}
			temp, err := FindInvoice(rHash)
			if err != nil {
				CUST.Warning(err.Error())
				continue
			}
			if int16(temp.State) != v.Status {
				v.Status = int16(temp.State)
				mutex.Lock()
				defer mutex.Unlock()
				if v.Status == int16(lnrpc.Invoice_SETTLED) {
					ba := models.Balance{}
					ba.AccountId = *v.AccountID
					ba.Amount = float64(v.Amount)
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
