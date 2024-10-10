package SecondHander

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/btlLog"
	"trade/services/custodyAccount/lockPayment"
)

type GetBalanceRequest struct {
	Npubkey string `json:"npubkey"`
	AssetId string `json:"assetId"`
}
type GetBalanceResponse struct {
	UnlockedBalance float64 `json:"unlockedBalance"`
	LockedBalance   float64 `json:"lockedBalance"`
}

func GetBalance(c *gin.Context) {
	var creds GetBalanceRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err, unlockedBalance, lockedBalance := lockPayment.GetBalance(creds.Npubkey, creds.AssetId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var res GetBalanceResponse
	res.UnlockedBalance = unlockedBalance
	res.LockedBalance = lockedBalance
	c.JSON(http.StatusOK, gin.H{"error": res})
}

type LockRequest struct {
	Npubkey  string  `json:"npubkey"`
	LockedId string  `json:"lockedId"`
	AssetId  string  `json:"assetId"`
	Amount   float64 `json:"amount"`
}
type LockResponse struct {
	Error error `json:"error"`
}

func Lock(c *gin.Context) {
	var creds LockRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO Verification request
	err := lockPayment.Lock(creds.Npubkey, creds.LockedId, creds.AssetId, creds.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": nil})
}

type UnlockRequest struct {
	Npubkey  string  `json:"npubkey"`
	LockedId string  `json:"lockedId"`
	AssetId  string  `json:"assetId"`
	Amount   float64 `json:"amount"`
}
type UnlockResponse struct {
	Error error `json:"error"`
}

func Unlock(c *gin.Context) {
	var creds UnlockRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO Verification request
	err := lockPayment.Unlock(creds.Npubkey, creds.LockedId, creds.AssetId, creds.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": nil})

}

type PayByLockedRequest struct {
	LockedId        string  `json:"lockedId"`
	PayerNpubkey    string  `json:"payerNpubkey"`
	ReceiverNpubkey string  `json:"receiverNpubkey"`
	AssetId         string  `json:"assetId"`
	Amount          float64 `json:"amount"`
	PayType         int8    `json:"payType"`
}
type PayType int8

const (
	PayTypeLock PayType = iota
	PayTypeUnlock
)

type PayByLockedResponse struct {
	TxId string `json:"txId"`
}

func PayAsset(c *gin.Context) {
	var creds PayByLockedRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO Verification request
	var err error
	if creds.PayType == int8(PayTypeLock) {
		err = lockPayment.TransferByLock(creds.LockedId, creds.PayerNpubkey, creds.ReceiverNpubkey, creds.AssetId, creds.Amount)
	} else if creds.PayType == int8(PayTypeUnlock) {
		err = lockPayment.TransferByUnlock(creds.LockedId, creds.PayerNpubkey, creds.ReceiverNpubkey, creds.AssetId, creds.Amount)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": nil})
}

// TODO: add more api
func GetLockedBalanceList(c *gin.Context) {}
