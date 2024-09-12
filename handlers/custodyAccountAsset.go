package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"trade/btlLog"
	"trade/models"
	"trade/services/custodyAccount/custodyAssets"
	"trade/services/custodyAccount/custodyBase"
)

type ApplyAddressRequest struct {
	Amount  float64 `json:"amount"`
	AssetId string  `json:"asset_id"`
}

func ApplyAddress(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	apply := ApplyAddressRequest{}
	if err := c.ShouldBindJSON(&apply); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	e, err := custodyAssets.NewAssetEvent(userName, apply.AssetId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	req, err := e.ApplyPayReq(&custodyAssets.AssetAddressApplyRequest{
		Amount: int64(apply.Amount),
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"req": req})
}

type SendAssetRequest struct {
	Address string `json:"address"`
}

func SendAsset(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	apply := SendAssetRequest{}
	if err := c.ShouldBindJSON(&apply); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	e, err := custodyAssets.NewAssetEvent(userName, "")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	err = e.SendPayment(&custodyAssets.AssetPacket{
		PayReq: apply.Address,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

type QueryAssetRequest struct {
	AssetId string `json:"asset_id"`
}

func QueryAsset(c *gin.Context) {
	userName := c.MustGet("username").(string)
	apply := QueryAssetRequest{}
	if err := c.ShouldBindJSON(&apply); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	e, err := custodyAssets.NewAssetEvent(userName, apply.AssetId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	balance, err := e.GetBalance()
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func QueryAssets(c *gin.Context) {
	userName := c.MustGet("username").(string)
	e, err := custodyAssets.NewAssetEvent(userName, "")
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

type AwardAssetRequest struct {
	Username string `json:"username"`
	AssetId  string `json:"assetId"`
	Amount   int    `json:"amount"`
	Memo     string `json:"memo"`
}

func Award(c *gin.Context) {
	var creds AwardAssetRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	e, err := custodyAssets.NewAssetEvent(creds.Username, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = custodyAssets.PutInAward(e.UserInfo.Account, creds.AssetId, creds.Amount, &creds.Memo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    nil,
	})
}

func QueryAddress(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	invoiceRequest := struct {
		AssetId string `json:"asset_id"`
	}{}
	if err := c.ShouldBindJSON(&invoiceRequest); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Request is erro"})
		return
	}
	e, err := custodyAssets.NewAssetEvent(userName, invoiceRequest.AssetId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MakeJsonErrorResultForHttp(models.DefaultErr, "用户不存在", nil))
		return
	}
	// 查询账户发票
	invoices, err := e.QueryPayReq()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "service error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"addr": invoices})
}

func QueryAddresses(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	e, err := custodyAssets.NewAssetEvent(userName, "")
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
	invoices, err := e.QueryPayReqs()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "service error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"addrs": invoices})
}

func QueryAssetPayment(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	invoiceRequest := struct {
		AssetId string `json:"asset_id"`
	}{}
	if err := c.ShouldBindJSON(&invoiceRequest); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Request is erro"})
		return
	}
	e, err := custodyAssets.NewAssetEvent(userName, invoiceRequest.AssetId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MakeJsonErrorResultForHttp(models.DefaultErr, "用户不存在", nil))
		return
	}
	// 查询账户发票
	invoices, err := e.GetTransactionHistoryByAsset()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "service error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"addr": invoices})
}

func QueryAssetPayments(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	e, err := custodyAssets.NewAssetEvent(userName, "")
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
	invoices, err := e.GetTransactionHistory()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "service error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"addrs": invoices})
}

type AssetBalance struct {
	AssetId string  `json:"assetId"`
	Amount  int64   `json:"amount"`
	Price   float64 `json:"prices"`
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			btlLog.CUST.Error("Error closing response body:", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		btlLog.CUST.Error("Error reading response body:", err)
		return nil
	}
	type temp struct {
		AssetsId string  `json:"id"`
		Price    float64 `json:"price"`
	}
	type List struct {
		List []temp `json:"list"`
	}
	r := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		Code    int    `json:"code"`
		Data    List   `json:"data"`
	}{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil
	}
	var list []AssetBalance
	for _, v := range r.Data.List {
		list = append(list, AssetBalance{
			AssetId: v.AssetsId,
			Amount:  t[v.AssetsId],
			//Price:   v.Price,
			Price: 0,
		})
	}
	return &list
}
