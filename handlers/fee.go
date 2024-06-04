package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
	"trade/utils"
)

func QueryFeeRate(c *gin.Context) {
	feeRate, err := services.GetFeeRate()
	if err != nil {
		utils.LogError("Get FeeRate.", err)
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

func QueryAllFeeRate(c *gin.Context) {
	allFeeRate, err := services.GetFeeRate()
	if err != nil {
		utils.LogError("Get FeeRate.", err)
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
