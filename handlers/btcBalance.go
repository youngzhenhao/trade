package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func GetBtcBalance(c *gin.Context) {
	username := c.MustGet("username").(string)
	btcBalance, err := services.GetBtcBalanceByUsername(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetBtcBalanceByUsernameErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    btcBalance,
	})
}

func SetBtcBalance(c *gin.Context) {
	username := c.MustGet("username").(string)
	var btcBalanceSetRequest models.BtcBalanceSetRequest
	err := c.ShouldBindJSON(&btcBalanceSetRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    "",
		})
		return
	}
	btcBalance := models.BtcBalance{
		Username:           username,
		TotalBalance:       btcBalanceSetRequest.TotalBalance,
		ConfirmedBalance:   btcBalanceSetRequest.ConfirmedBalance,
		UnconfirmedBalance: btcBalanceSetRequest.UnconfirmedBalance,
		LockedBalance:      btcBalanceSetRequest.LockedBalance,
	}
	err = services.CreateOrUpdateBtcBalance(&btcBalance)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.CreateOrUpdateBtcBalanceErr,
			Data:    "",
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    nil,
	})
}
