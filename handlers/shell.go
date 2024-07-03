package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func GenerateBlockOne(c *gin.Context) {
	out, err := services.GenerateBlocks(1)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GenerateBlocksErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    out,
	})
}

func FaucetTransferOneTenthBtc(c *gin.Context) {
	address := c.Param("address")
	out, err := services.FaucetTransferBtc(address, 0.1)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.FaucetTransferBtcErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    out,
	})
}

func FaucetTransferOneHundredthBtc(c *gin.Context) {
	address := c.Param("address")
	out, err := services.FaucetTransferBtc(address, 0.01)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.FaucetTransferBtcErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    out,
	})
}
