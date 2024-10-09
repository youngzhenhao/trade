package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"trade/models"
	"trade/services"
)

func GetAssetManagedUtxoByUserId(c *gin.Context) {
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
	assetManagedUtxos, err := services.GetAssetManagedUtxosByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetManagedUtxosByUserIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetManagedUtxos,
	})
}

func GetAssetManagedUtxoAssetId(c *gin.Context) {
	assetId := c.Param("asset_id")
	assetManagedUtxos, err := services.GetAssetManagedUtxosByAssetId(assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetManagedUtxoByAssetIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetManagedUtxos,
	})
}

func GetAssetManagedUtxoIdsByUserId(c *gin.Context) {
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
	assetManagedUtxos, err := services.GetAssetManagedUtxosByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetManagedUtxosByUserIdErr,
			Data:    nil,
		})
		return
	}
	assetManagedUtxoIds := services.AssetManagedUtxosToAssetManagedUtxoIds(assetManagedUtxos)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetManagedUtxoIds,
	})
}

func GetAssetManagedUtxoAssetIdsByUserId(c *gin.Context) {
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
	assetManagedUtxos, err := services.GetAssetManagedUtxosByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetManagedUtxosByUserIdErr,
			Data:    nil,
		})
		return
	}
	assetIds := services.AssetManagedUtxosToAssetIds(assetManagedUtxos)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetIds,
	})
}

func SetAssetManagedUtxos(c *gin.Context) {
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
	var assetManagedUtxoSetRequests []models.AssetManagedUtxoSetRequest
	err = c.ShouldBindJSON(&assetManagedUtxoSetRequests)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	assetManagedUtxos := services.ProcessAssetManagedUtxoSetRequests(userId, username, &assetManagedUtxoSetRequests)
	err = services.CreateOrUpdateAssetManagedUtxos(userId, assetManagedUtxos)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.SetAssetManagedUtxosErr,
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

func RemoveAssetManagedUtxos(c *gin.Context) {
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
	var assetManagedUtxoIds []int
	err = c.ShouldBindJSON(&assetManagedUtxoIds)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	err = services.ValidateUserIdAndAssetManagedUtxoIds(userId, &assetManagedUtxoIds)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ValidateUserIdAndAssetManagedUtxoIdsErr,
			Data:    nil,
		})
		return
	}
	err = services.RemoveAssetManagedUtxoByIds(&assetManagedUtxoIds)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.SetAssetManagedUtxosErr,
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

func GetAllAssetManagedUtxoSimplified(c *gin.Context) {
	assetManagedUtxoSimplified, err := services.GetAllAssetManagedUtxoSimplified()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAllAssetManagedUtxoSimplifiedErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetManagedUtxoSimplified,
	})
}

func GetAssetManagedUtxoLimitAndOffset(c *gin.Context) {
	var getAssetManagedUtxoLimitAndOffsetRequest services.GetAssetManagedUtxoLimitAndOffsetRequest
	err := c.ShouldBindJSON(&getAssetManagedUtxoLimitAndOffsetRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	assetId := getAssetManagedUtxoLimitAndOffsetRequest.AssetId
	limit := getAssetManagedUtxoLimitAndOffsetRequest.Limit
	offset := getAssetManagedUtxoLimitAndOffsetRequest.Offset

	{
		// @dev: total page number
		number, err := services.GetAssetManagedUtxoPageNumberByPageSize(assetId, limit)
		// @dev: limit is pageSize
		pageNumber := offset/limit + 1
		if pageNumber > number {
			err = errors.New("page number must be greater than max value " + strconv.Itoa(number))
			c.JSON(http.StatusOK, models.JsonResult{
				Success: false,
				Error:   err.Error(),
				Code:    models.PageNumberExceedsTotalNumberErr,
				Data:    nil,
			})
			return
		}
	}

	assetManagedUtxo, err := services.GetAssetManagedUtxoLimitAndOffset(assetId, limit, offset)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetManagedUtxoLimitAndOffsetErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    assetManagedUtxo,
	})
}

func GetAssetManagedUtxoPageNumberByPageSize(c *gin.Context) {
	var getAssetManagedUtxoPageNumberByPageSizeRequest services.GetAssetManagedUtxoPageNumberByPageSizeRequest
	err := c.ShouldBindJSON(&getAssetManagedUtxoPageNumberByPageSizeRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	pageSize := getAssetManagedUtxoPageNumberByPageSizeRequest.PageSize
	assetId := getAssetManagedUtxoPageNumberByPageSizeRequest.AssetId
	if pageSize <= 0 || assetId == "" {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   errors.New("invalid asset id or page size").Error(),
			Code:    models.GetAssetManagedUtxoPageNumberByPageSizeRequestInvalidErr,
			Data:    nil,
		})
		return
	}
	pageNumber, err := services.GetAssetManagedUtxoPageNumberByPageSize(assetId, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetManagedUtxoPageNumberByPageSizeErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    pageNumber,
	})
}
