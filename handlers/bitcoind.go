package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/api"
	"trade/models"
)

func GetAddressByOutpointInMainnet(c *gin.Context) {
	outpoint := c.Param("op")
	address, err := api.GetAddressByOutpoint(models.Mainnet, outpoint)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAddressByOutpointErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    address,
	})
}

func GetAddressByOutpointInTestnet(c *gin.Context) {
	outpoint := c.Param("op")
	address, err := api.GetAddressByOutpoint(models.Testnet, outpoint)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAddressByOutpointErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    address,
	})
}

func GetAddressByOutpointInRegtest(c *gin.Context) {
	outpoint := c.Param("op")
	address, err := api.GetAddressByOutpoint(models.Regtest, outpoint)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAddressByOutpointErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    address,
	})
}
