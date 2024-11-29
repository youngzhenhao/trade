package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/api"
	"trade/models"
	"trade/services"
)

func GetAssetMetaImage(c *gin.Context) {
	var assetIds []string
	err := c.ShouldBindJSON(&assetIds)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   0,
		})
		return
	}
	datas, errs := services.GetAssetMetaImage(assetIds)
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data: gin.H{
			"datas": datas,
			"errs":  errs,
		},
	})
}

func GetGroupFirstImageDataInMainnet(c *gin.Context) {
	groupKey := c.Query("group_key")
	imageData, err := api.GetGroupFirstImageData(models.Mainnet, groupKey)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.GetGroupFirstImageDataErr.Code(),
			ErrMsg: err.Error(),
			Data:   "",
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   imageData,
	})
}

func GetGroupFirstImageDataInTestnet(c *gin.Context) {
	groupKey := c.Query("group_key")
	imageData, err := api.GetGroupFirstImageData(models.Testnet, groupKey)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.GetGroupFirstImageDataErr.Code(),
			ErrMsg: err.Error(),
			Data:   "",
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   imageData,
	})
}

func GetGroupFirstImageDataInRegtest(c *gin.Context) {
	groupKey := c.Query("group_key")
	imageData, err := api.GetGroupFirstImageData(models.Regtest, groupKey)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.GetGroupFirstImageDataErr.Code(),
			ErrMsg: err.Error(),
			Data:   "",
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   imageData,
	})
}
