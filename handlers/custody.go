package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/btlLog"
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
	err, unlockedBalance, lockedBalance := lockPayment.GetBalance(userName, creds.AssetId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	res.UnlockedBalance = unlockedBalance
	res.LockedBalance = lockedBalance
	res.TotalBalance = unlockedBalance + lockedBalance
	c.JSON(http.StatusOK, res)
}

func GetAssetBalanceList(c *gin.Context) {
	userName := c.MustGet("username").(string)
	list := custodyAccount.GetAssetBalanceList(userName)
	c.JSON(http.StatusOK, list)
}
