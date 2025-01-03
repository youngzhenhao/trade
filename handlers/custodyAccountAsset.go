package handlers

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"trade/btlLog"
	"trade/models"
	"trade/services/btldb"
	"trade/services/custodyAccount/custodyBase"
	"trade/services/custodyAccount/custodyBase/custodyFee"
	"trade/services/custodyAccount/defaultAccount/custodyAssets"
	"trade/services/custodyAccount/defaultAccount/custodyBtc/mempool"
	rpc "trade/services/servicesrpc"
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
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	e, err := custodyAssets.NewAssetEvent(userName, apply.AssetId)
	if err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	req, err := e.ApplyPayReq(&custodyAssets.AssetAddressApplyRequest{
		Amount: int64(apply.Amount),
	})
	if err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	addr := struct {
		Address string `json:"addr"`
	}{
		Address: req.GetPayReq(),
	}
	c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "", addr))
}

type SendAssetRequest struct {
	Address string `json:"address"`
}

func SendAsset(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	apply := SendAssetRequest{}
	if err := c.ShouldBindJSON(&apply); err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	e, err := custodyAssets.NewAssetEvent(userName, "")
	if err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}

	e.UserInfo.PaymentMux.Lock()
	defer e.UserInfo.PaymentMux.Unlock()

	err = e.SendPayment(&custodyAssets.AssetPacket{
		PayReq: apply.Address,
	})
	if err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	result := struct {
		Success string `json:"success"`
	}{
		Success: "success",
	}
	c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "", result))
}

func SendToUserAsset(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	e, err := custodyAssets.NewAssetEvent(userName, "")
	if err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}

	e.UserInfo.PaymentMux.Lock()
	defer e.UserInfo.PaymentMux.Unlock()

	pay := struct {
		NpubKey string  `json:"npub_key"`
		AssetId string  `json:"asset_id"`
		Amount  float64 `json:"amount"`
	}{}
	if err = c.ShouldBindJSON(&pay); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error() + "请求参数错误"})
		return
	}

	err = e.SendPaymentToUser(pay.NpubKey, pay.Amount, pay.AssetId)
	if err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	result := struct {
		Success string `json:"success"`
	}{
		Success: "success",
	}
	c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "", result))
}

type QueryAssetRequest struct {
	AssetId string `json:"asset_id"`
}

func QueryAsset(c *gin.Context) {
	userName := c.MustGet("username").(string)
	apply := QueryAssetRequest{}
	if err := c.ShouldBindJSON(&apply); err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	e, err := custodyAssets.NewAssetEvent(userName, apply.AssetId)
	if err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	balance, err := e.GetBalance()
	if err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "", balance))
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

type AddressResponce struct {
	Address string  `json:"addr"`
	AssetId string  `json:"asset_id"`
	Amount  float64 `json:"amount"`
}

func QueryAddress(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	invoiceRequest := struct {
		AssetId string `json:"asset_id"`
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
	// 查询账户发票
	addr, err := e.QueryPayReq()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	var addrs []AddressResponce
	for _, v := range addr {
		addrs = append(addrs, AddressResponce{
			Address: v.Invoice,
			AssetId: v.AssetId,
			Amount:  v.Amount,
		})
	}
	c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "", addrs))
}

func QueryAddresses(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	e, err := custodyAssets.NewAssetEvent(userName, "")
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MakeJsonErrorResultForHttp(models.DefaultErr, "用户不存在", nil))
		return
	}
	// 查询账户发票
	addr, err := e.QueryPayReqs()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	var addrs []AddressResponce
	for _, v := range addr {
		addrs = append(addrs, AddressResponce{
			Address: v.Invoice,
			AssetId: v.AssetId,
			Amount:  v.Amount,
		})
	}
	c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "", addrs))
}

func QueryAssetPayment(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	invoiceRequest := struct {
		AssetId string `json:"asset_id"`
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
	p, err := e.GetTransactionHistoryByAsset()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	p.Sort()
	//p2, err := custodyAccount.LockPaymentToPaymentList(e.UserInfo, *e.AssetId)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	//result := custodyAccount.MergePaymentList(p, p2)

	c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "", p))
}

func QueryAssetPayments(c *gin.Context) {
	// 获取登录用户信息
	userName := c.MustGet("username").(string)
	e, err := custodyAssets.NewAssetEvent(userName, "")
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MakeJsonErrorResultForHttp(models.DefaultErr, "用户不存在", nil))
		return
	}
	invoiceRequest := custodyBase.PaymentRequest{}
	if err := c.ShouldBindJSON(&invoiceRequest); err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	if invoiceRequest.Page == 0 && invoiceRequest.PageSize == 0 {
		invoiceRequest.Page = 1
		invoiceRequest.PageSize = 1000
		invoiceRequest.Away = 5
	}
	// 查询账户发票
	payments, err := e.GetTransactionHistory(&invoiceRequest)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "", payments))
}

type AssetBalance struct {
	AssetId string `json:"assetId"`
	Amount  int64  `json:"amount"`
	Price   int64  `json:"prices"`
}

func DealBalance(b []custodyBase.Balance) *[]AssetBalance {
	baseURL := "http://api.nostr.microlinktoken.com/realtime/one_price"
	queryParams := url.Values{}
	t := make(map[string]int64)
	for _, v := range b {
		if v.AssetId == "00" {
			queryParams.Add("ids", "btc")
		} else {
			queryParams.Add("ids", v.AssetId)
		}
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
		if v.AssetsId == "btc" {
			v.AssetsId = "00"
		}
		list = append(list, AssetBalance{
			AssetId: v.AssetsId,
			Amount:  t[v.AssetsId],
			Price:   int64(v.Price),
			//Price: 0,
		})
	}
	return &list
}

type DecodeAddressRequest struct {
	Address string `json:"addr"`
}

func DecodeAddress(c *gin.Context) {
	query := DecodeAddressRequest{}
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, "请求参数错误", nil))
		return
	}
	if query.Address == "" {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, "请求参数错误", nil))
	}
	q, err := rpc.DecodeAddr(query.Address)
	if err != nil {
		c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.DefaultErr, "地址解析失败："+err.Error(), nil))
		return
	}
	AssetId := hex.EncodeToString(q.AssetId)
	result := struct {
		AssetId   string  `json:"AssetId"`
		AssetType string  `json:"assetType"`
		Amount    uint64  `json:"amount"`
		FeeRate   float64 `json:"feeRate"`
	}{
		AssetId:   AssetId,
		AssetType: q.AssetType.String(),
		Amount:    q.Amount,
	}
	//判断地址是否为本地发票
	_, err = btldb.GetInvoiceByReq(query.Address)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		result.FeeRate = 0
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		result.FeeRate = float64(mempool.GetCustodyAssetFee())
	} else {
		result.FeeRate = float64(custodyFee.AssetInsideFee)
	}
	c.JSON(http.StatusOK, models.MakeJsonErrorResultForHttp(models.SUCCESS, "", result))
}
