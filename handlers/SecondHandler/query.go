package SecondHandler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/btlLog"
	"trade/services/custodyAccount/localQuery"
)

type Result struct {
	Errno  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func QueryBills(c *gin.Context) {
	var creds localQuery.BillQueryQuest
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: err.Error(), Data: nil})
		return
	}
	if creds.Page == 0 {
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: "Page must be greater than 0", Data: nil})
		return
	}
	creds.Page = creds.Page - 1
	a, count, err := localQuery.BillQuery(creds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	Bill := struct {
		Count int64                          `json:"count"`
		Bills *[]localQuery.BillListWithUser `json:"bills"`
	}{
		Count: count,
		Bills: a,
	}
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: Bill})
}

func QueryBalance(c *gin.Context) {
	var creds localQuery.BalanceQueryQuest
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: err.Error(), Data: nil})
		return
	}
	a := localQuery.BalanceQuery(creds)
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: a})
}

func QueryBalancesChange(c *gin.Context) {
	var creds localQuery.BalanceChangeQuest
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: err.Error(), Data: nil})
		return
	}
	if creds.AssetId == "" || creds.UserName == "" {
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: "AssetId and UserName must not be empty", Data: nil})
		return
	}
	if creds.Page != 0 {
		creds.Page = creds.Page - 1
	}
	count, resp, err := localQuery.BalancesChange(creds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	a := struct {
		Count int64                           `json:"count"`
		List  *[]localQuery.BalanceChangeResp `json:"list"`
	}{
		Count: count,
		List:  resp,
	}
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: a})
}

func GetBalanceList(c *gin.Context) {
	var creds localQuery.GetAssetListQuest
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: err.Error(), Data: nil})
		return
	}
	if creds.Page == 0 {
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: "Page must be greater than 0", Data: nil})
		return
	}
	creds.Page = creds.Page - 1

	a, count, total := localQuery.GetAssetList(creds)
	list := struct {
		Count int64                          `json:"count"`
		List  *[]localQuery.GetAssetListResp `json:"list"`
		Total float64                        `json:"total"`
	}{
		Count: count,
		List:  a,
		Total: total,
	}
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: list})
}

func TotalBillList(c *gin.Context) {
	var creds localQuery.TotalBillListQuest
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: err.Error(), Data: nil})
		return
	}
	if creds.Page == 0 {
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: "Page must be greater than 0", Data: nil})
		return
	}
	creds.Page = creds.Page - 1
	a, count, err := localQuery.TotalBillList(&creds)
	list := struct {
		Count int64                           `json:"count"`
		List  *[]localQuery.TotalBillListResp `json:"list"`
	}{
		Count: count,
		List:  &a,
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
	}
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: list})
}

func QueryLockedBills(c *gin.Context) {
	var creds localQuery.LockedBillsQueryQuest
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: err.Error(), Data: nil})
		return
	}
	if creds.Page == 0 {
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: "Page must be greater than 0", Data: nil})
		return
	}
	creds.Page = creds.Page - 1
	bills, count, err := localQuery.LockedBillsQuery(creds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	Bill := struct {
		Count int64                              `json:"count"`
		Bills *[]localQuery.LockedBillsQueryResp `json:"bills"`
	}{
		Count: count,
		Bills: bills,
	}
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: Bill})
}

func QueryChannelAssetInfo(c *gin.Context) {
	var creds localQuery.ChannelQueryQuest
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: err.Error(), Data: nil})
		return
	}
	if creds.AssetId == "" {
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: "AssetId must not be empty", Data: nil})
		return
	}
	assetIdInfo, err := localQuery.QueryChannelAssetInfo(&creds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: assetIdInfo})
}
