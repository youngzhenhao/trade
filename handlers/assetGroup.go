package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
	"trade/utils"
)

func GetGroupFirstAssetMeta(c *gin.Context) {
	tweakedGroupKey := c.Param("group_key")
	if tweakedGroupKey == "" && len(tweakedGroupKey) != 66 || !utils.IsHexString(tweakedGroupKey) {
		err := errors.New("invalid tweaked group key")
		c.JSON(http.StatusOK, models.JsonResult{
			Success: true,
			Error:   err.Error(),
			Code:    models.InvalidTweakedGroupKeyErr,
			Data:    nil,
		})
		return
	}
	assetGroup, err := services.GetAssetGroup(tweakedGroupKey)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetGroupErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetGroup.FirstAssetMeta,
	})
}

func GetGroupFirstAssetId(c *gin.Context) {
	tweakedGroupKey := c.Param("group_key")
	if tweakedGroupKey == "" && len(tweakedGroupKey) != 66 || !utils.IsHexString(tweakedGroupKey) {
		err := errors.New("invalid tweaked group key")
		c.JSON(http.StatusOK, models.JsonResult{
			Success: true,
			Error:   err.Error(),
			Code:    models.InvalidTweakedGroupKeyErr,
			Data:    nil,
		})
		return
	}
	assetGroup, err := services.GetAssetGroup(tweakedGroupKey)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetGroupErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetGroup.FirstAssetId,
	})
}

func SetGroupFirstAssetMeta(c *gin.Context) {
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
	var assetGroupSetRequest models.AssetGroupSetRequest
	err = c.ShouldBindJSON(&assetGroupSetRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	err = services.SetAssetGroup(userId, username, &assetGroupSetRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.SetAssetGroupErr,
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
