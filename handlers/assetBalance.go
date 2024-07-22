package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func GetAssetBalance(c *gin.Context) {
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
	assetBalances, err := services.GetAssetBalancesByUserIdNonZero(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetBalancesByUserIdNonZeroErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetBalances,
	})
}

func SetAssetBalance(c *gin.Context) {
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
	var assetBalanceSetRequest models.AssetBalanceSetRequest
	err = c.ShouldBindJSON(&assetBalanceSetRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	assetBalance := services.ProcessAssetBalanceSetRequest(userId, username, &assetBalanceSetRequest)
	err = services.CreateOrUpdateAssetBalance(assetBalance, userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.CreateOrUpdateAssetBalanceErr,
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

func SetAssetBalances(c *gin.Context) {
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
	var assetBalanceSetRequestSlice []models.AssetBalanceSetRequest
	err = c.ShouldBindJSON(&assetBalanceSetRequestSlice)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	assetBalances := services.ProcessAssetBalanceSetRequestSlice(userId, username, &assetBalanceSetRequestSlice)
	err = services.CreateOrUpdateAssetBalances(assetBalances, userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.CreateOrUpdateAssetBalancesErr,
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

func GetAssetHolderNumber(c *gin.Context) {
	assetId := c.Param("asset_id")
	holderNumber, err := services.GetAssetHolderNumberAssetBalance(assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetHolderNumberAssetBalanceErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    holderNumber,
	})
}

func GetAssetHolderBalance(c *gin.Context) {
	assetId := c.Param("asset_id")
	holderBalances, err := services.GetAssetIdAndBalancesByAssetId(assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetIdAndBalancesByAssetIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    holderBalances,
	})
}

func GetAssetHolderBalanceLimitAndOffset(c *gin.Context) {
	var assetIdLimitOffset models.AssetHolderBalanceLimitAndOffsetRequest
	err := c.ShouldBindJSON(&assetIdLimitOffset)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	assetId := assetIdLimitOffset.AssetId
	limit := assetIdLimitOffset.Limit
	offset := assetIdLimitOffset.Offset
	isValid, err := services.IsLimitAndOffsetValid(assetId, limit, offset)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.IsLimitAndOffsetValidErr,
			Data:    nil,
		})
		return
	}
	if !isValid {
		err = errors.New("records number is less equal than offset")
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.IsLimitAndOffsetValidErr,
			Data:    nil,
		})
		return
	}
	holderBalances, err := services.GetAssetIdAndBalancesByAssetIdLimitAndOffset(assetId, limit, offset)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetIdAndBalancesByAssetIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    holderBalances,
	})
}

func GetAssetHolderBalanceRecordsNumber(c *gin.Context) {
	assetId := c.Param("asset_id")
	// @dev: Query total records number
	recordsNum, err := services.GetAssetBalanceByAssetIdNonZeroLength(assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetBalanceByAssetIdNonZeroLengthErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    recordsNum,
	})
}

func GetAssetHolderUsernameBalanceAll(c *gin.Context) {
	usernameBalances, err := services.GetAllUsernameAssetBalances()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAllUsernameAssetBalancesErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    usernameBalances,
	})
}
