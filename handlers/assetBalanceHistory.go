package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func GetLatestAssetBalanceHistories(c *gin.Context) {
	username := c.MustGet("username").(string)
	records, err := services.GetLatestAssetBalanceHistories(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetLatestAssetBalanceHistoriesErr,
			Data:    new([]models.AssetBalanceHistoryRecord),
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    records,
	})
}

func CreateAssetBalanceHistories(c *gin.Context) {
	username := c.MustGet("username").(string)
	var assetBalanceHistorySetRequests []models.AssetBalanceHistorySetRequest
	err := c.ShouldBindJSON(&assetBalanceHistorySetRequests)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	err = services.CreateAssetBalanceHistories(username, &assetBalanceHistorySetRequests)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.CreateAssetBalanceHistoriesErr,
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
