package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/btlLog"
	"trade/models"
	"trade/services/custodyAccount"
	"trade/services/custodyAccount/lockPayment"
)

type GetBalanceRequest struct {
	AssetId string `json:"assetId"`
}
type GetBalanceResponse struct {
	TotalBalance    float64 `json:"totalBalance"`
	UnlockedBalance float64 `json:"unlockedBalance"`
	LockedBalance   float64 `json:"lockedBalance"`
	Tag1Balance     float64 `json:"tag1Balance"`
}

func GetBalance(c *gin.Context) {
	var creds GetBalanceRequest
	var res GetBalanceResponse
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, res)
		return
	}
	userName := c.MustGet("username").(string)
	err, unlockedBalance, lockedBalance, tag1 := lockPayment.GetBalance(userName, creds.AssetId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	res.UnlockedBalance = unlockedBalance
	res.LockedBalance = lockedBalance
	res.TotalBalance = unlockedBalance + lockedBalance
	res.Tag1Balance = tag1
	c.JSON(http.StatusOK, res)
}

func GetAssetBalanceList(c *gin.Context) {
	userName := c.MustGet("username").(string)
	list, err := custodyAccount.GetAssetBalanceList(userName)
	if err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, fmt.Sprintf("GetAssetBalanceList failed: %v", err.Error()), nil))
		return
	}
	if list == nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, fmt.Sprintf("GetAssetBalanceList failed"), nil))
		return
	}
	request := DealBalance(*list)
	if request == nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "", list))
	} else {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "", request))
	}
}
