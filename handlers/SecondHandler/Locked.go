package SecondHandler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trade/btlLog"
	"trade/middleware"
	"trade/models"
	"trade/services/custodyAccount/lockPayment"
)

type GetBalanceRequest struct {
	Npubkey string `json:"npubkey"`
	AssetId string `json:"assetId"`
}
type GetBalanceResponse struct {
	UnlockedBalance float64 `json:"unlockedBalance"`
	LockedBalance   float64 `json:"lockedBalance"`
	LockedId        string  `json:"lockedId"`
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

	err, unlockedBalance, lockedBalance, tag1 := lockPayment.GetBalance(creds.Npubkey, creds.AssetId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	res.UnlockedBalance = unlockedBalance
	res.LockedBalance = lockedBalance
	res.Tag1Balance = tag1
	c.JSON(http.StatusOK, res)
}

type LockRequest struct {
	Npubkey  string  `json:"npubkey"`
	LockedId string  `json:"lockedId"`
	AssetId  string  `json:"assetId"`
	Amount   float64 `json:"amount"`
	Tag      int     `json:"tag"`
}
type LockResponse struct {
	Error string `json:"error"`
}

func Lock(c *gin.Context) {
	var creds LockRequest
	var res LockResponse
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		res.Error = err.Error()
		c.JSON(http.StatusBadRequest, &res)
		return
	}
	//TODO Verification request
	err := lockPayment.Lock(creds.Npubkey, creds.LockedId, creds.AssetId, creds.Amount, creds.Tag)
	if err != nil {
		res.Error = err.Error()
		c.JSON(http.StatusInternalServerError, &res)
		return
	}
	c.JSON(http.StatusOK, &res)
}

type UnlockRequest struct {
	Npubkey  string  `json:"npubkey"`
	LockedId string  `json:"lockedId"`
	AssetId  string  `json:"assetId"`
	Amount   float64 `json:"amount"`
	Tag      int     `json:"tag"`
}
type UnlockResponse struct {
	Error string `json:"error"`
}

func Unlock(c *gin.Context) {
	var creds UnlockRequest
	var res UnlockResponse
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		res.Error = err.Error()
		c.JSON(http.StatusBadRequest, &res)
		return
	}
	//TODO Verification request
	err := lockPayment.Unlock(creds.Npubkey, creds.LockedId, creds.AssetId, creds.Amount, creds.Tag)
	if err != nil {
		res.Error = err.Error()
		c.JSON(http.StatusInternalServerError, &res)
		return
	}
	c.JSON(http.StatusOK, &res)
}

type PayByLockedRequest struct {
	LockedId        string  `json:"lockedId"`
	PayerNpubkey    string  `json:"payerNpubkey"`
	ReceiverNpubkey string  `json:"receiverNpubkey"`
	AssetId         string  `json:"assetId"`
	Amount          float64 `json:"amount"`
	PayType         int8    `json:"payType"`
	Tag             int     `json:"tag"`
}
type PayType int8

const (
	PayTypeLock PayType = iota
	PayTypeUnlock
)

type PayByLockedResponse struct {
	TxId      string `json:"txId"`
	ErrorCode int    `json:"code"`
}

func PayAsset(c *gin.Context) {
	var creds PayByLockedRequest
	var res PayByLockedResponse
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		res.ErrorCode = lockPayment.GetErrorCode(lockPayment.BadRequest)
		c.JSON(http.StatusBadRequest, &res)
		return
	}
	//TODO Verification request
	var err error
	if creds.PayType == int8(PayTypeLock) {
		err = lockPayment.TransferByLock(creds.LockedId, creds.PayerNpubkey, creds.ReceiverNpubkey, creds.AssetId, creds.Amount, creds.Tag)
	} else if creds.PayType == int8(PayTypeUnlock) {
		if creds.Tag != 0 {
			res.ErrorCode = lockPayment.GetErrorCode(lockPayment.BadRequest)
			c.JSON(http.StatusBadRequest, &res)
			return
		}
		err = lockPayment.TransferByUnlock(creds.LockedId, creds.PayerNpubkey, creds.ReceiverNpubkey, creds.AssetId, creds.Amount)
	}
	if err != nil {
		res.ErrorCode = lockPayment.GetErrorCode(err)
		c.JSON(http.StatusInternalServerError, &res)
		return
	}
	res.TxId = creds.LockedId
	c.JSON(http.StatusOK, &res)
}

type CheckUserStatusRequest struct {
	Npubkey string `json:"npubkey"`
}

type CheckUserStatusResponse struct {
	Status int8   `json:"status"`
	Error  string `json:"error"`
}

func CheckUserStatus(c *gin.Context) {
	var creds CheckUserStatusRequest
	var res CheckUserStatusResponse
	if err := c.ShouldBindJSON(&creds); err != nil {
		res.Error = fmt.Sprintf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, &res)
		return
	}
	db := middleware.DB
	var user models.User
	err := db.Where("user_name = ?", creds.Npubkey).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		res.Error = fmt.Sprintf("Server error: %v", err)
		c.JSON(http.StatusInternalServerError, &res)
		return
	}
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		res.Status = -1
		res.Error = "User not found"
	case user.Status != 0:
		res.Status = 0
		res.Error = "User locked"
	default:
		res.Status = 1
	}
	c.JSON(http.StatusOK, &res)
}
