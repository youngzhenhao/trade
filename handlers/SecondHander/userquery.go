package SecondHander

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/btlLog"
	"trade/services/custodyAccount/localQuery"
)

func QueryUserInfo(c *gin.Context) {
	var creds localQuery.UserInfoRep
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: err.Error(), Data: nil})
		return
	}
	if len(creds.Username) <= 5 {
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: "username is error", Data: nil})
		return
	}
	userInfo, err := localQuery.GetUserInfo(creds.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: userInfo})
}

func BlockUser(c *gin.Context) {
	var creds localQuery.BlockUserReq
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: err.Error(), Data: nil})
		return
	}

	if len(creds.Username) <= 5 {
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: "username is error", Data: nil})
		return
	}
	if creds.Memo == "" {
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: "memo is error", Data: nil})
		return
	}

	err := localQuery.BlockUser(creds.Username, creds.Memo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: nil})
}

func UnBlockUser(c *gin.Context) {
	var creds localQuery.UnblockUserReq
	if err := c.ShouldBindJSON(&creds); err != nil {
		btlLog.CUST.Error("%v", err)
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: err.Error(), Data: nil})
		return
	}
	if len(creds.Username) <= 5 {
		c.JSON(http.StatusBadRequest, Result{Errno: 400, ErrMsg: "username is error", Data: nil})
		return
	}
	err := localQuery.UnblockUser(creds.Username, creds.Memo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{Errno: 500, ErrMsg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, Result{Errno: 0, ErrMsg: "", Data: nil})
}
