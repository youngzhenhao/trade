package SecondHander

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services/custodyAccount/btc_channel"
)

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

type AwardRequest struct {
	Username string `json:"username"`
	Amount   int    `json:"amount"`
	Memo     string `json:"memo"`
}
