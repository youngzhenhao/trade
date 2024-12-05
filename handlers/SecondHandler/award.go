package SecondHandler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/btlLog"
	"trade/models"
	"trade/services/custodyAccount/account"
	"trade/services/custodyAccount/defaultAccount/custodyAssets"
	"trade/services/custodyAccount/defaultAccount/custodyBtc"
	"trade/services/custodyAccount/lockPayment"
)

type AwardRequest struct {
	Username    string `json:"username"`
	Amount      int    `json:"amount"`
	Memo        string `json:"memo"`
	LockedId    string `json:"lockedId"`
	AccountType string `json:"accountType"`
}

func PutInSatoshiAward(c *gin.Context) {
	var creds AwardRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	btlLog.CUST.Info("%v", creds)
	e, err := custodyBtc.NewBtcChannelEvent(creds.Username)
	if err != nil {
		if !errors.Is(err, account.ErrUserLocked) {
			btlLog.CUST.Info("%v", creds)
			btlLog.CUST.Info("%v", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch creds.AccountType {
	case "default":
		award, err := custodyBtc.PutInAward(e.UserInfo, "", creds.Amount, &creds.Memo, creds.LockedId)
		if err != nil {
			btlLog.CUST.Error("%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		btlLog.CUST.Info("Success PutInSatoshiAward %v,%v, %v", e.UserInfo.User.Username, creds.Amount, award.ID)
		c.JSON(http.StatusOK, models.JsonResult{
			Success: true,
			Error:   "",
			Code:    models.SUCCESS,
			Data:    award.ID,
		})
		return
	case "locked":
		award, err := lockPayment.PutInAwardLockBTC(e.UserInfo, float64(creds.Amount), &creds.Memo, creds.LockedId)
		if err != nil {
			btlLog.CUST.Error("%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		btlLog.CUST.Info("Success PutInSatoshiAward %v,%v, %v", e.UserInfo.User.Username, creds.Amount, award.ID)
		c.JSON(http.StatusOK, models.JsonResult{
			Success: true,
			Error:   "",
			Code:    models.SUCCESS,
			Data:    award.ID,
		})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account type"})
		return
	}
}

type AwardAssetRequest struct {
	Username    string `json:"username"`
	AssetId     string `json:"assetId"`
	Amount      int    `json:"amount"`
	Memo        string `json:"memo"`
	LockedId    string `json:"lockedId"`
	AccountType string `json:"accountType"`
}

func PutAssetAward(c *gin.Context) {
	var creds AwardAssetRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	e, err := custodyAssets.NewAssetEvent(creds.Username, "")
	if err != nil {
		if !errors.Is(err, account.ErrUserLocked) {
			btlLog.CUST.Info("%v", creds)
			btlLog.CUST.Info("%v", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	switch creds.AccountType {
	case "default":
		award, err := custodyAssets.PutInAward(e.UserInfo, creds.AssetId, creds.Amount, &creds.Memo, creds.LockedId)
		if err != nil {
			btlLog.CUST.Error("%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		btlLog.CUST.Info("Success PutAssetAward %v, %v, %v, %v", e.UserInfo.Account, creds.AssetId, creds.Amount, award.ID)
		c.JSON(http.StatusOK, models.JsonResult{
			Success: true,
			Error:   "",
			Code:    models.SUCCESS,
			Data:    award.ID,
		})
		return
	case "locked":
		award, err := lockPayment.PutInAwardLockAsset(e.UserInfo, creds.AssetId, float64(creds.Amount), &creds.Memo, creds.LockedId)
		if err != nil {
			btlLog.CUST.Error("%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		btlLog.CUST.Info("Success PutAssetAward %v, %v, %v, %v", e.UserInfo.Account, creds.AssetId, creds.Amount, award.ID)
		c.JSON(http.StatusOK, models.JsonResult{
			Success: true,
			Error:   "",
			Code:    models.SUCCESS,
			Data:    award.ID,
		})
		return

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account type"})
		return
	}
}
