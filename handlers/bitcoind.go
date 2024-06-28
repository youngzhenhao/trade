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

func GetAddressesByOutpointSliceInRegtest(c *gin.Context) {
	var outpointSlice struct {
		Outpoints []string `json:"outpoints"`
	}
	err := c.ShouldBindJSON(&outpointSlice)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    "",
		})
		return
	}
	addresses, err := api.GetAddressesByOutpointSlice(models.Regtest, outpointSlice.Outpoints)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAddressesByOutpointSliceErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    addresses,
	})
}
