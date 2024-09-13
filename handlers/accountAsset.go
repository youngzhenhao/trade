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

func GetAccountAssetBalanceLimitAndOffset(c *gin.Context) {
	var getAccountAssetBalanceLimitAndOffsetRequest services.GetAccountAssetBalanceLimitAndOffsetRequest
	err := c.ShouldBindJSON(&getAccountAssetBalanceLimitAndOffsetRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	assetId := getAccountAssetBalanceLimitAndOffsetRequest.AssetId
	limit := getAccountAssetBalanceLimitAndOffsetRequest.Limit
	offset := getAccountAssetBalanceLimitAndOffsetRequest.Offset
	accountAssetBalances, err := services.GetAccountAssetBalanceExtendsLimitAndOffset(assetId, limit, offset)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAccountAssetBalancesLimitAndOffsetErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    accountAssetBalances,
	})
}

func GetAccountAssetBalancePageNumberByPageSize(c *gin.Context) {
	var GetAccountAssetBalancePageNumberByPageSizeRequest services.GetAccountAssetBalancePageNumberByPageSizeRequest
	err := c.ShouldBindJSON(&GetAccountAssetBalancePageNumberByPageSizeRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	pageSize := GetAccountAssetBalancePageNumberByPageSizeRequest.PageSize
	assetId := GetAccountAssetBalancePageNumberByPageSizeRequest.AssetId
	if pageSize <= 0 || assetId == "" {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   errors.New("invalid asset id or page size").Error(),
			Code:    models.GetAccountAssetBalancePageNumberByPageSizeRequestInvalidErr,
			Data:    nil,
		})
		return
	}
	pageNumber, err := services.GetAccountAssetBalancePageNumberByPageSize(assetId, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAccountAssetBalancePageNumberByPageSizeErr,
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

func GetAccountAssetTransferLimitAndOffset(c *gin.Context) {
	var getAccountAssetTransferLimitAndOffsetRequest services.GetAccountAssetTransferLimitAndOffsetRequest
	err := c.ShouldBindJSON(&getAccountAssetTransferLimitAndOffsetRequest)
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
	accountAssetTransfers, err := services.GetAccountAssetTransfersLimitAndOffset(assetId, limit, offset)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAccountAssetTransfersLimitAndOffsetErr,
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
