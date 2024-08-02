package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func GetAssetLocalMintHistoryByUserId(c *gin.Context) {
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
	assetLocalMintHistories, err := services.GetAssetLocalMintHistoriesByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetLocalMintHistoriesByUserIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetLocalMintHistories,
	})
}

func GetAssetLocalMintHistoryAssetId(c *gin.Context) {
	assetId := c.Param("asset_id")
	assetLocalMintHistory, err := services.GetAssetLocalMintHistoryByAssetId(assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetLocalMintHistoryByAssetIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetLocalMintHistory,
	})
}

func SetAssetLocalMintHistory(c *gin.Context) {
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
	var assetLocalMintHistorySetRequest models.AssetLocalMintHistorySetRequest
	err = c.ShouldBindJSON(&assetLocalMintHistorySetRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	assetLocalMintHistory := services.ProcessAssetLocalMintHistorySetRequest(userId, username, assetLocalMintHistorySetRequest)
	err = services.CreateOrUpdateAssetLocalMintHistory(&assetLocalMintHistory)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.SetAssetLocalMintHistoryErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    nil,
	})
}

func SetAssetLocalMintHistories(c *gin.Context) {
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
	var assetLocalMintHistorySetRequests []models.AssetLocalMintHistorySetRequest
	err = c.ShouldBindJSON(&assetLocalMintHistorySetRequests)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	assetLocalMintHistories := services.ProcessAssetLocalMintHistorySetRequests(userId, username, &assetLocalMintHistorySetRequests)
	err = services.CreateOrUpdateAssetLocalMintHistories(assetLocalMintHistories)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.SetAssetLocalMintHistoriesErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    nil,
	})
}

func GetAllAssetLocalMintHistorySimplified(c *gin.Context) {
	assetLocalMintHistorySimplified, err := services.GetAllAssetLocalMintHistorySimplified()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAllAssetLocalMintHistorySimplifiedErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetLocalMintHistorySimplified,
	})
}
