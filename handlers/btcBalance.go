package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
			Data:    nil,
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
			Data:    nil,
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

func GetBtcBalanceCount(c *gin.Context) {
	var count int64
	var err error
	count, err = services.GetBtcBalanceCount()
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.GetBtcBalanceCountErr.Code(),
			ErrMsg: err.Error(),
			Data:   0,
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   count,
	})
}

func GetBtcBalanceOrderLimitOffset(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")

	if limit == "" {
		err := errors.New("limit is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]services.BtcBalanceInfo{},
		})
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]services.BtcBalanceInfo{},
		})
		return
	}
	if limitInt < 0 {
		err := errors.New("limit is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]services.BtcBalanceInfo{},
		})
		return
	}
	if offset == "" {
		err := errors.New("offset is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]services.BtcBalanceInfo{},
		})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]services.BtcBalanceInfo{},
		})
		return
	}
	if offsetInt < 0 {
		err := errors.New("offset is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]services.BtcBalanceInfo{},
		})
		return
	}
	var btcBalanceInfos *[]services.BtcBalanceInfo

	btcBalanceInfos, err = services.GetBtcBalanceOrderLimitOffset(limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.GetBtcBalanceOrderLimitOffsetErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]services.BtcBalanceInfo{},
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   btcBalanceInfos,
	})
}
