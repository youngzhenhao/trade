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

func GetAddressesByOutpointSliceInMainnet(c *gin.Context) {
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
	addresses, err := api.GetAddressesByOutpointSlice(models.Mainnet, outpointSlice.Outpoints)
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

func GetAddressesByOutpointSliceInTestnet(c *gin.Context) {
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
	addresses, err := api.GetAddressesByOutpointSlice(models.Testnet, outpointSlice.Outpoints)
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

func GetTransactionsByOutpointSliceInMainnet(c *gin.Context) {
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
	transactions, err := api.GetTransactionsByOutpointSlice(models.Mainnet, outpointSlice.Outpoints)
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
		Data:    transactions,
	})
}

func GetTransactionsByOutpointSliceInTestnet(c *gin.Context) {
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
	transactions, err := api.GetTransactionsByOutpointSlice(models.Testnet, outpointSlice.Outpoints)
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
		Data:    transactions,
	})
}

func GetTransactionsByOutpointSliceInRegtest(c *gin.Context) {
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
	transactions, err := api.GetTransactionsByOutpointSlice(models.Regtest, outpointSlice.Outpoints)
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
		Data:    transactions,
	})
}

func GetTransactionByOutpointInMainnet(c *gin.Context) {
	outpoint := c.Param("op")
	transaction, err := api.GetTransactionByOutpoint(models.Mainnet, outpoint)
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
		Data:    transaction,
	})
}

func GetTransactionByOutpointInTestnet(c *gin.Context) {
	outpoint := c.Param("op")
	transaction, err := api.GetTransactionByOutpoint(models.Testnet, outpoint)
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
		Data:    transaction,
	})
}

func GetTransactionByOutpointInRegtest(c *gin.Context) {
	outpoint := c.Param("op")
	transaction, err := api.GetTransactionByOutpoint(models.Regtest, outpoint)
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
		Data:    transaction,
	})
}
