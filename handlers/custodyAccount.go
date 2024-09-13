package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services/custodyAccount"
	"trade/services/custodyAccount/btc_channel"
	rpc "trade/services/servicesrpc"
)

// ApplyInvoice CustodyAccount开具发票
func ApplyInvoice(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	e, err := btc_channel.NewBtcChannelEvent(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error() + "用户不存在"})
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
	c.JSON(http.StatusOK, gin.H{"invoice": req.GetPayReq()})
}

func QueryInvoice(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	e, err := btc_channel.NewBtcChannelEvent(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error() + "用户不存在"})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error() + "用户不存在"})
		return
	}
	//获取支付发票请求
	pay := custodyAccount.PayInvoiceRequest{}
	if err := c.ShouldBindJSON(&pay); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error() + "请求参数错误"})
		return
	}
	a := btc_channel.BtcPacket{
		PayReq:   pay.Invoice,
		FeeLimit: pay.FeeLimit,
	}
	err2 := e.SendPayment(&a)
	if err2 != nil {
		c.JSON(http.StatusOK, gin.H{"error": "SendPayment error:" + err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"payment": "success"})
}

// QueryBalance CustodyAccount查询发票
func QueryBalance(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	e, err := btc_channel.NewBtcChannelEvent(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error() + "用户不存在"})
		return
	}
	getBalance, err := e.GetBalance()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error() + "用户不存在"})
		return
	}
	//获取交易查询请求
	query := custodyAccount.PaymentRequest{}
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if query.AssetId != "00" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "asset_id类型错误"})
		return
	}
	p, _ := e.GetTransactionHistory()
	if s, ok := p.(*btc_channel.BtcPaymentList); ok {
		c.JSON(http.StatusOK, gin.H{"payments": s.PaymentList})
	}
}

// DecodeInvoice  解析发票
func DecodeInvoice(c *gin.Context) {
	query := custodyAccount.DecodeInvoiceRequest{}
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, "请求参数错误", nil))
		return
	}
	q, err := rpc.InvoiceDecode(query.Invoice)
	if err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, "发票解析失败："+err.Error(), nil))
		return
	}
	result := struct {
		Amount    int64  `json:"amount"`
		Timestamp int64  `json:"timestamp"`
		Expiry    int64  `json:"expiry"`
		Memo      string `json:"memo"`
	}{
		Amount:    q.NumSatoshis,
		Timestamp: q.Timestamp,
		Expiry:    q.Expiry,
		Memo:      q.Description,
	}
	c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "", result))
}
