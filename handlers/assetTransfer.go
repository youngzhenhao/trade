package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func SetAssetTransfer(c *gin.Context) {
	username := c.MustGet("username").(string)
	userId, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	var assetTransferProcessedSetRequestSlice []models.AssetTransferProcessedSetRequest
	err = c.ShouldBindJSON(&assetTransferProcessedSetRequestSlice)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	var assetTransferProcessedSlice *[]models.AssetTransferProcessedDb
	var assetTransferProcessedInputsSlice *[]models.AssetTransferProcessedInputDb
	var assetTransferProcessedOutputsSlice *[]models.AssetTransferProcessedOutputDb
	assetTransferProcessedSlice, assetTransferProcessedInputsSlice, assetTransferProcessedOutputsSlice, err = services.ProcessAssetTransferProcessedSlice(userId, &assetTransferProcessedSetRequestSlice)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ProcessAssetTransferErr,
			Data:    nil,
		})
		return
	}
	err = services.CreateOrUpdateAssetTransferProcessedSlice(assetTransferProcessedSlice)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.CreateAssetTransferProcessedErr,
			Data:    nil,
		})
		return
	}
	// @dev: Store inputs and outputs in db
	err = services.CreateOrUpdateAssetTransferProcessedInputSlice(assetTransferProcessedInputsSlice)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.CreateOrUpdateAssetTransferProcessedInputSliceErr,
			Data:    nil,
		})
		return
	}
	err = services.CreateOrUpdateAssetTransferProcessedOutputSlice(assetTransferProcessedOutputsSlice)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.CreateOrUpdateAssetTransferProcessedOutputSliceErr,
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

func GetAssetTransfer(c *gin.Context) {
	username := c.MustGet("username").(string)
	userId, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	assetTransferCombinedSlice, err := services.GetAssetTransferCombinedSliceByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetTransferCombinedSliceByUserIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    assetTransferCombinedSlice,
	})
}

func GetAssetTransferTxids(c *gin.Context) {
	username := c.MustGet("username").(string)
	userId, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	txids, err := services.GetAssetTransferTxidsByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetTransferProcessedSliceByUserIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    txids,
	})
}
