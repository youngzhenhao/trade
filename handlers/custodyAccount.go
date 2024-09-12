package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services/custodyAccount"
	"trade/services/custodyAccount/btc_channel"
)

// CreateCustodyAccount 创建托管账户
func CreateCustodyAccount(c *gin.Context) {
	//TODO
}

// ApplyInvoice CustodyAccount开具发票
func ApplyInvoice(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	e, err := btc_channel.NewBtcChannelEvent(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MakeJsonErrorResultForHttp(models.DefaultErr, "用户不存在", nil))
		return
	}
	apply := custodyAccount.ApplyRequest{}
	if err = c.ShouldBindJSON(&apply); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	a := btc_channel.BtcApplyInvoiceRequest{
		Amount: apply.Amount,
		Memo:   apply.Memo,
	}
	req, err := e.ApplyPayReq(&a)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"invoice": req})
}

func QueryInvoice(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	e, err := btc_channel.NewBtcChannelEvent(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MakeJsonErrorResultForHttp(models.DefaultErr, "用户不存在", nil))
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
	invoices, err := e.QueryPayReq()
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
	e, err := btc_channel.NewBtcChannelEvent(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MakeJsonErrorResultForHttp(models.DefaultErr, "用户不存在", nil))
		return
	}
	//获取支付发票请求
	pay := custodyAccount.PayInvoiceRequest{}
	if err := c.ShouldBindJSON(&pay); err != nil {
		c.JSON(http.StatusUnauthorized, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	a := btc_channel.BtcPacket{
		PayReq:   pay.Invoice,
		FeeLimit: pay.FeeLimit,
	}
	err2 := e.SendPayment(&a)
	if err2 != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err2.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "支付成功", "Success"))
}

// QueryBalance CustodyAccount查询发票
func QueryBalance(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	e, err := btc_channel.NewBtcChannelEvent(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MakeJsonErrorResultForHttp(models.DefaultErr, "用户不存在", nil))
		return
	}
	getBalance, err := e.GetBalance()
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": getBalance[0].Amount})
}

// QueryPayment  查询支付记录
func QueryPayment(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	e, err := btc_channel.NewBtcChannelEvent(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MakeJsonErrorResultForHttp(models.DefaultErr, "用户不存在", nil))
		return
	}

	//获取交易查询请求
	query := custodyAccount.PaymentRequest{}
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if query.AssetId != "00" {
		return
	}
	p, _ := e.GetTransactionHistory()
	if s, ok := p.(*btc_channel.BtcPaymentList); ok {
		c.JSON(http.StatusOK, gin.H{"payments": s})
	}
}

// LookupInvoice 查询发票状态
func LookupInvoice(c *gin.Context) {
	//todo
}

// LookupPayment 查看支付记录
func LookupPayment(c *gin.Context) {}
