package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func QueryFeeRate(c *gin.Context) {
	feeRate, err := services.GetMempoolFeeRate()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get FeeRate. " + err.Error(),
			Data:    "",
		})
		return
	}
	satPerKw := feeRate.SatPerKw.FastestFee
	satPerB := feeRate.SatPerB.FastestFee
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
			Data:    "",
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
			Data:    "",
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Data:    feeRate,
	})
}
