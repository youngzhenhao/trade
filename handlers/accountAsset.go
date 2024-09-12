package handlers

import (
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

func GetAllAccountAssetBalanceSimplified(c *gin.Context) {
	//	 TODO
}
