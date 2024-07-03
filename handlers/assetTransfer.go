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
	var assetTransferProcessedSlice *[]models.AssetTransferProcessed
	assetTransferProcessedSlice, err = services.ProcessAssetTransferProcessedSlice(userId, &assetTransferProcessedSetRequestSlice)
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
	assetTransferProcessedSlice, err := services.GetAssetTransferProcessedSliceByUserId(userId)
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
		Data:    assetTransferProcessedSlice,
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
