package SecondHander

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services/custodyAccount/assets"
	"trade/services/custodyAccount/btc_channel"
)

type AwardRequest struct {
	Username string `json:"username"`
	Amount   int    `json:"amount"`
	Memo     string `json:"memo"`
}

func PutInSatoshiAward(c *gin.Context) {
	var creds AwardRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	e, err := btc_channel.NewBtcChannelEvent(creds.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = btc_channel.PutInAward(e.UserInfo.Account, "", creds.Amount, &creds.Memo)
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

type AwardAssetRequest struct {
	Username string `json:"username"`
	AssetId  string `json:"assetId"`
	Amount   int    `json:"amount"`
	Memo     string `json:"memo"`
}

func PutAssetAward(c *gin.Context) {
	var creds AwardAssetRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	e, err := assets.NewAssetEvent(creds.Username, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = assets.PutInAward(e.UserInfo.Account, creds.AssetId, creds.Amount, &creds.Memo)
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
