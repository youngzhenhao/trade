package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"trade/btlLog"
	"trade/models"
	"trade/services/btldb"
	"trade/services/custodyAccount"
	"trade/services/custodyAccount/assets"
	"trade/services/custodyAccount/custodyBase"
)

// CreateCustodyAccount 创建托管账户
func CreateCustodyAccount(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)

	// 校验登录用户信息
	user, err := btldb.ReadUserByUsername(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}
	// 判断用户是否已经创建账户
	_, err = btldb.ReadAccountByUserId(user.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//创建账户
		cstAccount, err := custodyAccount.CreateCustodyAccount(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"accountModel": cstAccount})
		return
	}
	if err != nil {
		btlLog.CUST.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{"error": "get account info error"})
	}
	c.JSON(http.StatusOK, gin.H{"error": "用户已存在"})

}

// ApplyInvoice CustodyAccount开具发票
func ApplyInvoice(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	user, err := btldb.ReadUserByUsername(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	account, err := btldb.ReadAccountByUserId(user.ID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "service is default"})
			return
		}
		account, err = custodyAccount.CreateCustodyAccount(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建托管账户失败"})
			return
		}
	}

	//TODO 判断申请金额是否超过通道余额,检查申请内容是否合法
	apply := custodyAccount.ApplyRequest{}
	if err = c.ShouldBindJSON(&apply); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if apply.Amount <= 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发票信息不合法"})
		return
	}
	//生成一张发票
	invoiceRequest, err := custodyAccount.ApplyInvoice(user.ID, account, &apply)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"invoice": invoiceRequest.PaymentRequest})
}

func QueryInvoice(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	user, err := btldb.ReadUserByUsername(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	invoiceRequest := struct {
		AssetId string `json:"asset_id"`
	}{}
	if err := c.ShouldBindJSON(&invoiceRequest); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Request is erro"})
		return
	}

	// 查询账户发票
	invoices, err := custodyAccount.QueryInvoiceByUserId(user.ID, invoiceRequest.AssetId)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "service error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"invoices": invoices})
}

// PayInvoice CustodyAccount付款发票
func PayInvoice(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	user, err := btldb.ReadUserByUsername(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}
	// 选择托管账户
	account, err := btldb.ReadAccountByUserId(user.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		btlLog.CUST.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询账户信息失败"})
		return
	}
	if account.UserAccountCode == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "未找到账户信息"})
		return
	}

	//获取支付发票请求
	pay := custodyAccount.PayInvoiceRequest{}
	if err := c.ShouldBindJSON(&pay); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 支付发票
	_, err = custodyAccount.PayInvoice(account, &pay, true)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"payment": "success"})
}

// QueryBalance CustodyAccount查询发票
func QueryBalance(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	user, err := btldb.ReadUserByUsername(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	// 查询账户余额
	balance, err := custodyAccount.QueryAccountBalanceByUserId(user.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// QueryPayment  查询支付记录
func QueryPayment(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	user, err := btldb.ReadUserByUsername(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	//获取交易查询请求
	query := custodyAccount.PaymentRequest{}
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 查询交易记录
	payments, err := custodyAccount.QueryPaymentByUserId(user.ID, query.AssetId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"payments": payments})

}

// LookupInvoice 查询发票状态
func LookupInvoice(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	user, err := btldb.ReadUserByUsername(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}
	// 选择托管账户
	account, err := btldb.ReadAccountByUserId(user.ID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if account.UserAccountCode == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "未找到账户信息"})
		return
	}

	//获取发票查询请求
	lookup := custodyAccount.LookupInvoiceRequest{}
	if err := c.ShouldBindJSON(&lookup); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// 查询发票状态
	invoice, err := custodyAccount.LookupInvoice(&lookup)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"invoice": invoice})
}

// LookupPayment 查看支付记录
func LookupPayment(c *gin.Context) {}

type AssetBalance struct {
	AssetId string  `json:"assetId"`
	Amount  int64   `json:"amount"`
	Price   float64 `json:"prices"`
}

func QueryAssets(c *gin.Context) {
	userName := c.MustGet("username").(string)
	e, err := assets.NewAssetEvent(userName, "")
	if err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	balance, err := e.GetBalances()
	if err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	request := DealBalance(balance)
	if request == nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "", balance))
	} else {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "", request))
	}
}

func DealBalance(b []custodyBase.Balance) *[]AssetBalance {
	baseURL := "http://api.nostr.microlinktoken.com/realtime/one_price"
	queryParams := url.Values{}
	t := make(map[string]int64)
	for _, v := range b {
		queryParams.Add("ids", v.AssetId)
		queryParams.Add("numbers", strconv.FormatInt(v.Amount, 10))
		t[v.AssetId] = v.Amount
	}
	reqURL := baseURL + "?" + queryParams.Encode()
	resp, err := http.Get(reqURL)
	if err != nil {
		btlLog.CUST.Error("Error making request:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		btlLog.CUST.Error("Error reading response body:", err)
		return nil
	}
	fmt.Println(string(body))
	type temp struct {
		AssetsId string  `json:"id"`
		Price    float64 `json:"price"`
	}
	type List struct {
		List []temp `json:"list"`
	}
	r := struct {
		Success bool           `json:"success"`
		Error   string         `json:"error"`
		Code    models.ErrCode `json:"code"`
		Data    List           `json:"list"`
	}{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil
	}
	fmt.Println(r)
	var list []AssetBalance
	for _, v := range r.Data.List {
		list = append(list, AssetBalance{
			AssetId: v.AssetsId,
			Amount:  t[v.AssetsId],
			Price:   v.Price,
			//Price: 3000,
		})
	}
	return &list
}
