package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func GeAssetLocalMintByUserId(c *gin.Context) {
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
	assetLocalMints, err := services.GetAssetLocalMintsByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetLocalMintsByUserIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetLocalMints,
	})
}

func GetAssetLocalMintAssetId(c *gin.Context) {
	assetId := c.Param("asset_id")
	assetLocalMint, err := services.GetAssetLocalMintByAssetId(assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetLocalMintByAssetIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetLocalMint,
	})
}

func SetAssetLocalMint(c *gin.Context) {
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
	var assetLocalMintSetRequest models.AssetLocalMintSetRequest
	err = c.ShouldBindJSON(&assetLocalMintSetRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	assetLocalMint := services.ProcessAssetLocalMintSetRequest(userId, username, assetLocalMintSetRequest)
	err = services.SetAssetLocalMint(&assetLocalMint)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.SetAssetLocalMintErr,
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

func SetAssetLocalMints(c *gin.Context) {
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
	var assetLocalMintSetRequests []models.AssetLocalMintSetRequest
	err = c.ShouldBindJSON(&assetLocalMintSetRequests)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	assetLocalMints := services.ProcessAssetLocalMintSetRequests(userId, username, &assetLocalMintSetRequests)
	err = services.SetAssetLocalMints(assetLocalMints)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.SetAssetLocalMintsErr,
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

func GetAllAssetLocalMintSimplified(c *gin.Context) {
	assetLocalMintSimplified, err := services.GetAllAssetLocalMintSimplified()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAllAssetLocalMintSimplifiedErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetLocalMintSimplified,
	})
}
