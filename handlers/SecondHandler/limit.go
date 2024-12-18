package SecondHandler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/btlLog"
	"trade/services/custodyAccount/custodyBase/custodyLimit"
)

func GetUserLimitHandler(c *gin.Context) {
	var creds = struct {
		Username  string `json:"username"`
		LimitType string `json:"limitType"`
		PageNum   int    `json:"pageNum"`
		PageSize  int    `json:"pageSize"`
	}{}
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	if creds.PageNum <= 0 {
		creds.PageNum = 1
	}
	if creds.PageSize <= 0 {
		creds.PageSize = 10
	}
	total, res, err := custodyLimit.GetUserLimit(creds.Username, creds.LimitType, creds.PageNum, creds.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	limits := struct {
		Count int64                     `json:"count"`
		Bills *[]custodyLimit.UserLimit `json:"limits"`
	}{
		Count: total,
		Bills: res,
	}
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: limits})
}

func SetUserLimitLevelHandler(c *gin.Context) {
	var creds = struct {
		Username  string `json:"username"`
		LimitType string `json:"limitType"`
		Level     int    `json:"level"`
	}{}
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	err := custodyLimit.SetUserLimitLevel(creds.Username, creds.LimitType, creds.Level)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: nil})
}

func SetUserTodayLimitHandler(c *gin.Context) {
	var creds = struct {
		Username     string `json:"username"`
		LimitType    string `json:"limitType"`
		UsefulAmount int    `json:"todayUsefulAmount"`
		UsefulCount  int    `json:"todayUsefulCount"`
	}{}
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}

	err := custodyLimit.SetUserTodayLimit(creds.Username, creds.LimitType, creds.UsefulAmount, creds.UsefulCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: nil})
}

func GetLimitTypeHandler(c *gin.Context) {
	var creds = struct {
		PageNum  int `json:"pageNum"`
		PageSize int `json:"pageSize"`
	}{}
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}

	types, count, err := custodyLimit.GetLimitTypes(creds.PageNum, creds.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	typesResults := struct {
		Count int64                      `json:"count"`
		Types *[]custodyLimit.LimitTypes `json:"types"`
	}{
		Count: count,
		Types: types,
	}

	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: typesResults})
}

func CreateOrUpdateLimitTypeHandle(c *gin.Context) {
	var creds = struct {
		AssetId      string `json:"assetId"`
		TransferType int    `json:"transferType"`
		LimitName    string `json:"limitName"`
	}{}
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}

	err := custodyLimit.CreateOrUpdateLimitType(creds.AssetId, creds.TransferType, creds.LimitName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: nil})
}

func GetLimitTypeLevelsHandle(c *gin.Context) {
	var creds = struct {
		LimitName string `json:"limitName"`
		PageNum   int    `json:"pageNum"`
		PageSize  int    `json:"pageSize"`
	}{}
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}

	levels, err := custodyLimit.GetLimitTypeLevels(creds.LimitName, creds.PageNum, creds.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: levels})
}

func CreateOrUpdateLimitTypeLevelHandle(c *gin.Context) {
	var creds = struct {
		LimitName string `json:"limitName"`
		Level     int    `json:"level"`
		Amount    int    `json:"amount"`
		Count     int    `json:"count"`
	}{}
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}

	err := custodyLimit.CreateOrUpdateLimitTypeLevel(creds.LimitName, creds.Level, creds.Amount, creds.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: nil})
}
