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
	}{}
	var res *[]custodyLimit.UserLimit

	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, res)
		return
	}
	res, err := custodyLimit.GetUserLimit(creds.Username, creds.LimitType)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, res)
}
