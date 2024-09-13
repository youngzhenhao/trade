package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func GetAccountAssetBalanceByAssetId(c *gin.Context) {
	assetId := c.Param("asset_id")
	accountAssetBalanceExtends, err := services.GetAccountAssetBalanceExtendsByAssetId(assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAccountAssetBalanceExtendsByAssetIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    accountAssetBalanceExtends,
	})
}

func GetAllAccountAssetTransferByAssetId(c *gin.Context) {
	assetId := c.Param("asset_id")
	accountAssetTransfers, err := services.GetAllAccountAssetTransfersByAssetId(assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAllAccountAssetTransfersByAssetIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    accountAssetTransfers,
	})
}

func GetAccountAssetTransferLimitAndOffset(c *gin.Context) {
	var getAccountAssetTransferLimitAndOffsetRequest services.GetAccountAssetTransferLimitAndOffsetRequest
	err := c.ShouldBindQuery(&getAccountAssetTransferLimitAndOffsetRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	assetId := getAccountAssetTransferLimitAndOffsetRequest.AssetId
	limit := getAccountAssetTransferLimitAndOffsetRequest.Limit
	offset := getAccountAssetTransferLimitAndOffsetRequest.Offset
	accountAssetTransfers, err := services.GetAllAccountAssetTransfersByAssetIdLimitAndOffset(assetId, limit, offset)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAllAccountAssetTransfersByAssetIdLimitAndOffsetErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    accountAssetTransfers,
	})
}

func GetAccountAssetTransferPageNumberByPageSize(c *gin.Context) {
	var GetAccountAssetTransferPageNumberByPageSizeRequest services.GetAccountAssetTransferPageNumberByPageSizeRequest
	err := c.ShouldBindJSON(&GetAccountAssetTransferPageNumberByPageSizeRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	pageSize := GetAccountAssetTransferPageNumberByPageSizeRequest.PageSize
	assetId := GetAccountAssetTransferPageNumberByPageSizeRequest.AssetId
	if pageSize <= 0 || assetId == "" {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   errors.New("invalid asset id or page size").Error(),
			Code:    models.GetAccountAssetTransferPageNumberByPageSizeRequestInvalidErr,
			Data:    nil,
		})
		return
	}
	pageNumber, err := services.GetAccountAssetTransferPageNumberByPageSize(assetId, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAccountAssetTransferPageNumberByPageSizeErr,
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

func GetAllAccountAssetBalanceSimplified(c *gin.Context) {
	//	 TODO
}
