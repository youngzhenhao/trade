package SecondHander

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/btlLog"
	"trade/services/custodyAccount/localQuery"
)

func QueryBills(c *gin.Context) {
	var creds localQuery.BillQueryQuest
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a, err := localQuery.BillQuery(creds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

func QueryBalance(c *gin.Context) {
	var creds localQuery.BalanceQueryQuest
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a := localQuery.BalanceQuery(creds)
	c.JSON(http.StatusOK, a)
}

func GetBalanceList(c *gin.Context) {
	var creds localQuery.GetAssetListQuest
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a := localQuery.GetAssetList(creds)
	c.JSON(http.StatusOK, a)
}
