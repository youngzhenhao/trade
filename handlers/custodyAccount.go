package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trade/services"
)

// CreateCustodyAccount 创建托管账户
func CreateCustodyAccount(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)

	// 校验登录用户信息
	user, err := services.ReadUserByUsername(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}
	// 判断用户是否已经创建账户
	accounts, _ := services.ReadAccountByUserIds(user.ID)
	if len(accounts) > 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "托管账户已存在"})
		return
	}
	//创建账户
	cstAccount, err := services.CreateCustodyAccount(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accountModel": cstAccount})
}

// ApplyInvoice CustodyAccount开具发票
func ApplyInvoice(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	user, err := services.ReadUserByUsername(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	account, err := services.ReadAccountByUserId(user.ID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "service is default"})
			return
		}
		account, err = services.CreateCustodyAccount(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建托管账户失败"})
			return
		}
	}

	//TODO 判断申请金额是否超过通道余额,检查申请内容是否合法
	apply := services.ApplyRequest{}
	if err = c.ShouldBindJSON(&apply); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	if apply.Amount <= 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发票信息不合法"})
		return
	}
	//生成一张发票
	invoiceRequest, err := services.ApplyInvoice(user.ID, account, &apply)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"invoice": invoiceRequest.PaymentRequest})
}

func QueryInvoice(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	user, err := services.ReadUserByUsername(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	invoiceRequest := struct {
		AssetId string `json:"asset_id"`
	}{}
	if err := c.ShouldBindJSON(&invoiceRequest); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Request is erro"})
	}

	// 查询账户发票
	invoices, err := services.QueryInvoiceByUserId(user.ID, invoiceRequest.AssetId)
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
	user, err := services.ReadUserByUsername(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}
	// 选择托管账户
	account, err := services.ReadAccountByUserId(user.ID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if account.UserAccountCode == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "未找到账户信息"})
		return
	}

	//获取支付发票请求
	pay := services.PayInvoiceRequest{}
	if err := c.ShouldBindJSON(&pay); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	// 支付发票
	_, err = services.PayInvoice(account, &pay)
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
	user, err := services.ReadUserByUsername(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	// 查询账户余额
	balance, err := services.QueryAccountBalanceByUserId(user.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// QueryPayment  查询支付记录
func QueryPayment(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	user, err := services.ReadUserByUsername(userName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	//获取交易查询请求
	query := services.PaymentRequest{}
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	// 查询交易记录
	payments, err := services.QueryPaymentByUserId(user.ID, query.AssetId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"payments": payments})

}
