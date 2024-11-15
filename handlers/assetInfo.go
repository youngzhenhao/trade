package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/api"
	"trade/models"
)

func GetAssetNames(c *gin.Context) {
	_ = c.MustGet("username").(string)
	var assetIds []string
	err := c.ShouldBindJSON(&assetIds)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	assetIdAndNames, err := api.GetAssetsName(assetIds)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetsNameErr,
			Data:    &[]api.AssetIdAndName{},
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetIdAndNames,
	})
}
