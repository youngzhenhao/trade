package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services/custodyAccount"
	"trade/services/custodyAccount/custodyAssets"
)

func QueryLockedPayments(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	invoiceRequest := struct {
		AssetId string `json:"asset_id"`
		Page    int    `json:"page"`
		Size    int    `json:"size"`
		Away    int    `json:"away"`
	}{}
	if err := c.ShouldBindJSON(&invoiceRequest); err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	e, err := custodyAssets.NewAssetEvent(userName, invoiceRequest.AssetId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MakeJsonErrorResultForHttp(models.DefaultErr, "用户不存在", nil))
		return
	}
	// 查询账户交易记录
	p, err := custodyAccount.LockPaymentToPaymentList(e.UserInfo, invoiceRequest.AssetId, invoiceRequest.Page, invoiceRequest.Size, invoiceRequest.Away)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	p.Sort()
	// 返回结果
	c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "", p))
}
