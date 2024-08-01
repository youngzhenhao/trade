package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"trade/models"
	"trade/services"
)

func QueryFeeRate(c *gin.Context) {
	feeRate, err := services.GetMempoolFeeRate()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get FeeRate. " + err.Error(),
			Data:    nil,
		})
		return
	}
	// @dev: Add 1 rat per b fee rate for recommend fee rate
	satPerKw := feeRate.SatPerKw.FastestFee + services.FeeRateSatPerBToSatPerKw(1)
	satPerB := feeRate.SatPerB.FastestFee + 1
	// @dev: Sat per b has been self-incremented
	btcPerKb := services.FeeRateSatPerBToBtcPerKb(satPerB)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Data: services.FeeRateResponse{
			SatPerKw: satPerKw,
			SatPerB:  satPerB,
			BtcPerKb: btcPerKb,
		},
	})
}

func QueryAllFeeRate(c *gin.Context) {
	allFeeRate, err := services.GetAllFeeRateInfos()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get FeeRate. " + err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Data:    allFeeRate,
	})
}

func QueryRecommendedFeeRate(c *gin.Context) {
	feeRate, err := services.GetMempoolFeeRate()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get FeeRate. " + err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Data:    feeRate,
	})
}

func QueryFairLaunchIssuanceFee(c *gin.Context) {
	feeRate := c.Query("fee_rate")
	feeRateSatPerKw, err := strconv.Atoi(feeRate)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.FeeRateAtoiErr,
			Data:    nil,
		})
		return
	}
	if feeRateSatPerKw < services.FeeRateSatPerBToSatPerKw(1) {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   errors.New("fee rate too low").Error(),
			Code:    models.FeeRateInvalidErr,
			Data:    nil,
		})
		return
	}
	fee := services.GetIssuanceTransactionGasFee(feeRateSatPerKw)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: false,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    fee,
	})
}

func QueryFairLaunchMintFee(c *gin.Context) {
	feeRate := c.Query("fee_rate")
	feeRateSatPerKw, err := strconv.Atoi(feeRate)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.FeeRateAtoiErr,
			Data:    nil,
		})
		return
	}
	if feeRateSatPerKw < services.FeeRateSatPerBToSatPerKw(1) {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   errors.New("fee rate too low").Error(),
			Code:    models.FeeRateInvalidErr,
			Data:    nil,
		})
		return
	}
	fee := services.GetMintedTransactionGasFee(feeRateSatPerKw)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: false,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    fee,
	})
}
